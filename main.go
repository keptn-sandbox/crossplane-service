package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	cloudevents "github.com/cloudevents/sdk-go/v2" // make sure to use v2 cloudevents here
	"github.com/kelseyhightower/envconfig"
	keptn "github.com/keptn/go-utils/pkg/lib/keptn"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
)

var keptnOptions = keptn.KeptnOpts{}

type envConfig struct {
	// Port on which to listen for cloudevents
	Port int `envconfig:"RCV_PORT" default:"8080"`
	// Path to which cloudevents are sent
	Path string `envconfig:"RCV_PATH" default:"/"`
	// Whether we are running locally (e.g., for testing) or on production
	Env string `envconfig:"ENV" default:"local"`
	// URL of the Keptn configuration service (this is where we can fetch files from the config repo)
	ConfigurationServiceUrl string `envconfig:"CONFIGURATION_SERVICE" default:""`
}

// ServiceName specifies the current services name (e.g., used as source when sending CloudEvents)
const ServiceName = "crossplane-service"

// path for the crossplane file that will be applied
const CrossPlaneFilename = "crossplane/cluster.yaml"

// EnvironemtsetupFinishedEventData is the name of an echo triggered event
const EnvironmentsetupEventTriggeredType = "sh.keptn.event.environment-setup.triggered"

// EnvironemtsetupFinishedEventData is the name of an echo started event
const EnvironmentsetupStartedEventType = "sh.keptn.event.environment-setup.started"

// EnvironemtsetupFinishedEventData is the name of an echo finished event
const EnvironmentsetupFinishedEventType = "sh.keptn.event.environment-setup.finished"

// EnvironemtsetupFinishedEventData is the data of an echo triggered event
type EnvironmentsetupTriggeredEventData struct {
	keptnv2.EventData
}

// EnvironemtsetupFinishedEventData is the data of an echo started event
type EnvironmentsetupStartedEventData struct {
	keptnv2.EventData
}

// EnvironemtsetupFinishedEventData is the data of an echo finished event
type EnvironmentsetupFinishedEventData struct {
	keptnv2.EventData
}

const EnvironmentTeardownTriggeredEventType = "sh.keptn.event.environment-teardown.triggered"
const EnvironmentTeardownStartedEventType = "sh.keptn.event.environment-teardown.started"
const EnvironmentTeardownFinishedEventType = "sh.keptn.event.environment-teardown.finished"

type EnvironmentTeardownTriggeredEventData struct {
	keptnv2.EventData
}
type EnvironmentTeardownStartedEventData struct {
	keptnv2.EventData
}
type EnvironmentTeardownFinishedEventData struct {
	keptnv2.EventData
}

/**
 * Parses a Keptn Cloud Event payload (data attribute)
 */
func parseKeptnCloudEventPayload(event cloudevents.Event, data interface{}) error {
	err := event.DataAs(data)
	if err != nil {
		log.Fatalf("Got Data Error: %s", err.Error())
		return err
	}
	return nil
}

/**
 * This method gets called when a new event is received from the Keptn Event Distributor
 * Depending on the Event Type will call the specific event handler functions, e.g: handleDeploymentFinishedEvent
 * See https://github.com/keptn/spec/blob/0.2.0-alpha/cloudevents.md for details on the payload
 */
func processKeptnCloudEvent(ctx context.Context, event cloudevents.Event) error {
	// create keptn handler
	log.Printf("Initializing Keptn Handler")
	myKeptn, err := keptnv2.NewKeptn(&event, keptnOptions)
	if err != nil {
		return errors.New("Could not create Keptn Handler: " + err.Error())
	}

	log.Printf("gotEvent(%s): %s - %s", event.Type(), myKeptn.KeptnContext, event.Context.GetID())

	if err != nil {
		log.Printf("failed to parse incoming cloudevent: %v", err)
		return err
	}

	/**
	* CloudEvents types in Keptn 0.8.0 follow the following pattern:
	* - sh.keptn.event.${EVENTNAME}.triggered
	* - sh.keptn.event.${EVENTNAME}.started
	* - sh.keptn.event.${EVENTNAME}.status.changed
	* - sh.keptn.event.${EVENTNAME}.finished
	*
	* For convenience, types can be generated using the following methods:
	* - triggered:      keptnv2.GetTriggeredEventType(${EVENTNAME}) (e.g,. keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName))
	* - started:        keptnv2.GetStartedEventType(${EVENTNAME}) (e.g,. keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName))
	* - status.changed: keptnv2.GetStatusChangedEventType(${EVENTNAME}) (e.g,. keptnv2.GetStatusChangedEventType(keptnv2.DeploymentTaskName))
	* - finished:       keptnv2.GetFinishedEventType(${EVENTNAME}) (e.g,. keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName))
	*
	* Keptn reserves some Cloud Event types, please read up on that here: https://keptn.sh/docs/0.8.x/manage/shipyard/
	*
	* For those Cloud Events the keptn/go-utils library conveniently provides several data structures
	* and strings in github.com/keptn/go-utils/pkg/lib/v0_2_0, e.g.:
	* - deployment: DeploymentTaskName, DeploymentTriggeredEventData, DeploymentStartedEventData, DeploymentFinishedEventData
	* - test: TestTaskName, TestTriggeredEventData, TestStartedEventData, TestFinishedEventData
	* - ... (they all follow the same pattern)
	*
	*
	* In most cases you will be interested in processing .triggered events (e.g., sh.keptn.event.deployment.triggered),
	* which you an achieve as follows:
	* if event.type() == keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName) { ... }
	*
	* Processing the event payload can be achieved as follows:
	*
	* eventData := &keptnv2.DeploymentTriggeredEventData{}
	* parseKeptnCloudEventPayload(event, eventData)
	*
	* See https://github.com/keptn/spec/blob/0.2.0-alpha/cloudevents.md for more details of Keptn Cloud Events and their payload
	* Also, see https://github.com/keptn-sandbox/echo-service/blob/a90207bc119c0aca18368985c7bb80dea47309e9/pkg/events.go as an example how to create your own CloudEvents
	**/

	/**
	* The following code presents a very generic implementation of processing almost all possible
	* Cloud Events that are retrieved by this service.
	* Please follow the documentation provided above for more guidance on the different types.
	* Feel free to delete parts that you don't need.
	**/
	switch event.Type() {

	// -------------------------------------------------------
	// your custom cloud event, e.g., sh.keptn.your-event
	// see https://github.com/keptn-sandbox/echo-service/blob/a90207bc119c0aca18368985c7bb80dea47309e9/pkg/events.go
	// for an example on how to generate your own CloudEvents and structs
	case keptnv2.GetTriggeredEventType("environment-setup"): // sh.keptn.event.your-event.triggered
		log.Printf("Processing environment-setup.triggered Event")

		eventData := &EnvironmentsetupTriggeredEventData{}
		parseKeptnCloudEventPayload(event, eventData)

		return HandleEnvironmentSetupTriggeredEvent(myKeptn, event, eventData)

		break
	case keptnv2.GetStartedEventType(keptnv2.ConfigureMonitoringTaskName): // sh.keptn.event.your-event.started
		log.Printf("Processing your-event.started Event")
		// eventData := &keptnv2.YourEventStartedEventData{}
		// parseKeptnCloudEventPayload(event, eventData)

		// Just log this event
		// return GenericLogKeptnCloudEventHandler(myKeptn, event, eventData)

		break
	case keptnv2.GetFinishedEventType(keptnv2.ConfigureMonitoringTaskName): // sh.keptn.event.your-event.finished
		log.Printf("Processing your-event.finished Event")

		// eventData := &keptnv2.YourEventFinishedEventData{}
		// parseKeptnCloudEventPayload(event, eventData)

		// Just log this event
		// return GenericLogKeptnCloudEventHandler(myKeptn, event, eventData)

		break

		// environment teardown event
	case keptnv2.GetTriggeredEventType("environment-teardown"): // sh.keptn.event.your-event.triggered
		log.Printf("Processing environment-teardown.triggered Event")

		eventData := &EnvironmentTeardownTriggeredEventData{}
		parseKeptnCloudEventPayload(event, eventData)

		return HandleEnvironmentTeardownTriggeredEvent(myKeptn, event, eventData)

		break
	case keptnv2.GetStartedEventType(keptnv2.ConfigureMonitoringTaskName): // sh.keptn.event.your-event.started
		log.Printf("Processing your-event.started Event")
		// eventData := &keptnv2.YourEventStartedEventData{}
		// parseKeptnCloudEventPayload(event, eventData)

		// Just log this event
		// return GenericLogKeptnCloudEventHandler(myKeptn, event, eventData)

		break
	case keptnv2.GetFinishedEventType(keptnv2.ConfigureMonitoringTaskName): // sh.keptn.event.your-event.finished
		log.Printf("Processing your-event.finished Event")

		// eventData := &keptnv2.YourEventFinishedEventData{}
		// parseKeptnCloudEventPayload(event, eventData)

		// Just log this event
		// return GenericLogKeptnCloudEventHandler(myKeptn, event, eventData)

		break
	}

	// Unknown Event -> Throw Error!
	var errorMsg string
	errorMsg = fmt.Sprintf("Unhandled Keptn Cloud Event: %s", event.Type())

	log.Print(errorMsg)
	return errors.New(errorMsg)
}

/**
 * Usage: ./main
 * no args: starts listening for cloudnative events on localhost:port/path
 *
 * Environment Variables
 * env=runlocal   -> will fetch resources from local drive instead of configuration service
 */
func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	os.Exit(_main(os.Args[1:], env))
}

/**
 * Opens up a listener on localhost:port/path and passes incoming requets to gotEvent
 */
func _main(args []string, env envConfig) int {
	// configure keptn options
	if env.Env == "local" {
		log.Println("env=local: Running with local filesystem to fetch resources")
		keptnOptions.UseLocalFileSystem = true
	}

	keptnOptions.ConfigurationServiceURL = env.ConfigurationServiceUrl

	log.Println("Starting crossplane-service...")
	log.Printf("    on Port = %d; Path=%s", env.Port, env.Path)

	ctx := context.Background()
	ctx = cloudevents.WithEncodingStructured(ctx)

	log.Printf("Creating new http handler")

	// configure http server to receive cloudevents
	p, err := cloudevents.NewHTTP(cloudevents.WithPath(env.Path), cloudevents.WithPort(env.Port))

	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("Starting receiver")
	log.Fatal(c.StartReceiver(ctx, processKeptnCloudEvent))

	return 0
}
