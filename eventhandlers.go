package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"gopkg.in/yaml.v2"
)

/**
* Here are all the handler functions for the individual event
* See https://github.com/keptn/spec/blob/0.8.0-alpha/cloudevents.md for details on the payload
**/

// GenericLogKeptnCloudEventHandler is a generic handler for Keptn Cloud Events that logs the CloudEvent
func GenericLogKeptnCloudEventHandler(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data interface{}) error {
	log.Printf("Handling %s Event: %s", incomingEvent.Type(), incomingEvent.Context.GetID())
	log.Printf("CloudEvent %T: %v", data, data)

	return nil
}

func HandleEnvironmentSetupTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *EnvironmentsetupTriggeredEventData) error {
	log.Printf("Handling environment-setup.triggered Event: %s", incomingEvent.Context.GetID())

	_, err := myKeptn.SendTaskStartedEvent(data, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	log.Printf("Looking for Crossplane cluster %s file in Keptn git repo...", CrossPlaneFilename)

	// load crossplane file
	keptnResourceContent, err := myKeptn.GetKeptnResource(CrossPlaneFilename)

	if err != nil {
		logMessage := fmt.Sprintf("No %s file found for service %s in stage %s in project %s", CrossPlaneFilename, data.Service, data.Stage, data.Project)
		log.Printf(logMessage)

		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}
	log.Printf("Crossplane file found.")

	// store crossplane file locally
	_ = os.Mkdir("crossplane", 0644)
	err = ioutil.WriteFile(CrossPlaneFilename, []byte(keptnResourceContent), 0644)
	if err != nil {
		logMessage := fmt.Sprintf("Could not store crossplane file locally: %s", err.Error())
		log.Printf(logMessage)
		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}
	log.Printf("Crossplane file stored locally.")

	log.Printf("Now applying crossplane file.")
	// now execute crossplane
	_, err = ExecuteCommand("kubectl", []string{"apply", "-f", CrossPlaneFilename})

	if err != nil {
		logMessage := fmt.Sprintf("Error while applying crossplane cluster manifest: %s", err.Error())
		log.Printf(logMessage)

		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}
	log.Printf("Crossplane file applied.")

	// waiting for cluster to be ready
	// wait for secret
	var secretName string
	secretDefaultName := "kubeconfig-keptn-crossplane"
	secretNamespace := "crossplane-system"
	// wait as usually the secret is not immediately available
	time.Sleep(10 * time.Second)
	for secretName != secretDefaultName {
		log.Printf("Checking availability of secret %s in namespace %s", secretDefaultName, secretNamespace)
		secretName, err = CheckAvailabilityOfSecret(secretDefaultName, secretNamespace)
		log.Printf("Retrieved secret name: %s", secretName)

		if err != nil {
			logMessage := fmt.Sprintf("Could not retrieve secret %s yet - waiting for 30 seconds", secretDefaultName)
			log.Printf(logMessage)

			// we will consider this a
			_, err = myKeptn.SendTaskStatusChangedEvent(&keptnv2.EventData{
				Message: logMessage,
			}, ServiceName)
			if err != nil {
				log.Printf("Error: %s", err)
			}
			// interval before we check the chaosengine status again
			time.Sleep(30 * time.Second)
		}

	}
	log.Printf("Secret found. Continuing...")

	// first getting the kubeconfig
	// kubectl get secrets cluster-details-keptn-crossplane -n default -o jsonpath={'.data.kubeconfig'} | base64 -d > kubeconfig
	kubeconfigEncoded, err := ExecuteCommand("kubectl", []string{"get", "secrets", secretName, "-n", secretNamespace, "-o", "jsonpath={.data.kubeconfig}"})
	if err != nil {
		logMessage := fmt.Sprintf("Error while getting kubeconfig: %s", err.Error())
		log.Printf(logMessage)
	}
	//log.Printf(kubeconfigEncoded)

	kubeconfig, err := base64.StdEncoding.DecodeString(kubeconfigEncoded)
	if err != nil {
		logMessage := fmt.Sprintf("Could not base64 decode the Kubeconfig: %s", err.Error())
		log.Printf(logMessage)
		return err
	}
	err = ioutil.WriteFile("kubeconfig", kubeconfig, 0644)
	if err != nil {
		logMessage := fmt.Sprintf("Could not store kubeconfig file locally: %s", err.Error())
		log.Printf(logMessage)
		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}

	if string(kubeconfig) == "" {
		logMessage := "KubeConfig is empty"
		log.Printf(logMessage)
	}

	// k get nodes --kubeconfig kubeconfig
	nodes, err := ExecuteCommand("kubectl", []string{"get", "nodes", "--kubeconfig", "kubeconfig"})
	if err != nil {
		logMessage := fmt.Sprintf("Error while getting kubeconfig: %s", err.Error())
		log.Printf(logMessage)
	}
	logMessage := nodes
	log.Printf(logMessage)

	_, err = myKeptn.SendTaskStatusChangedEvent(&keptnv2.EventData{
		Message: logMessage,
	}, ServiceName)
	if err != nil {
		log.Printf("Error: %s", err)
	}

	_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultPass,
	}, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task finished CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

func HandleEnvironmentTeardownTriggeredEvent(myKeptn *keptnv2.Keptn, incomingEvent cloudevents.Event, data *EnvironmentTeardownTriggeredEventData) error {
	log.Printf("Handling environment-teardown.triggered Event: %s", incomingEvent.Context.GetID())

	_, err := myKeptn.SendTaskStartedEvent(data, ServiceName)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task started CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	log.Printf("Looking for Crossplane cluster %s file in Keptn git repo...", CrossPlaneFilename)

	// load crossplane file
	keptnResourceContent, err := myKeptn.GetKeptnResource(CrossPlaneFilename)

	if err != nil {
		logMessage := fmt.Sprintf("No %s file found for service %s in stage %s in project %s", CrossPlaneFilename, data.Service, data.Stage, data.Project)
		log.Printf(logMessage)

		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}
	log.Printf("Crossplane file found.")

	// store crossplane file locally
	_ = os.Mkdir("crossplane", 0644)
	err = ioutil.WriteFile(CrossPlaneFilename, []byte(keptnResourceContent), 0644)
	if err != nil {
		logMessage := fmt.Sprintf("Could not store crossplane file locally: %s", err.Error())
		log.Printf(logMessage)
		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}
	log.Printf("Crossplane file stored locally.")

	log.Printf("Now starting to delete cluster based on crossplane file.")
	// now execute crossplane
	kubectlresult, err := ExecuteCommand("kubectl", []string{"delete", "-f", CrossPlaneFilename})
	log.Printf(kubectlresult)
	if err != nil {
		logMessage := fmt.Sprintf("Error while deleting crossplane cluster manifest: %s", err.Error())
		log.Printf(logMessage)

		_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
			Status:  keptnv2.StatusErrored,
			Result:  keptnv2.ResultFailed,
			Message: logMessage,
		}, ServiceName)

		return err
	}
	log.Printf("Crossplane cluster deleted.")

	_, err = myKeptn.SendTaskFinishedEvent(&keptnv2.EventData{
		Status: keptnv2.StatusSucceeded,
		Result: keptnv2.ResultPass,
	}, ServiceName)

	if err != nil {
		errMsg := fmt.Sprintf("Failed to send task finished CloudEvent (%s), aborting...", err.Error())
		log.Println(errMsg)
		return err
	}

	return nil
}

// ExecuteCommand exectues the command using the args
func ExecuteCommand(command string, args []string) (string, error) {
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("Error executing command %s %s: %s\n%s", command, strings.Join(args, " "), err.Error(), string(out))
	}
	return string(out), nil
}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func CheckAvailabilityOfSecret(secretname string, namespace string) (string, error) {
	secretname, err := ExecuteCommand("kubectl", []string{"get", "secrets", secretname, "-n", namespace, "-o", "jsonpath='{.metadata.name}'"})

	if err != nil {
		return "", err
	}

	return strings.Trim(secretname, "'"), nil
}

func (hv *helmValues) getHelmValues(filename string) *helmValues {

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, hv)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return hv
}
