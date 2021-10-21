#!/usr/bin/env sh

kubectl create ns monitoring

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm install prometheus prometheus-community/prometheus --namespace monitoring

kubectl create ns keptn
kubectl apply -f  https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.7.0/deploy/service.yaml 



# Prometheus installed namespace
kubectl set env deployment/prometheus-service -n keptn --containers="prometheus-service" PROMETHEUS_NS="monitoring"

# Setup Prometheus Endpoint
kubectl set env deployment/prometheus-service -n keptn --containers="prometheus-service" PROMETHEUS_ENDPOINT="http://prometheus-server.monitoring.svc.cluster.local:80"

# Alert Manager installed namespace
kubectl set env deployment/prometheus-service -n keptn --containers="prometheus-service" ALERT_MANAGER_NS="monitoring"

kubectl apply -f https://raw.githubusercontent.com/keptn-contrib/prometheus-service/release-0.6.2/deploy/role.yaml -n monitoring

# send the keptn configure monitoring event to the prometheus-service

{
  "specversion": "1.0",
  "id": "c4d3a334-6cb9-4e8c-a372-7e0b45942f53",
  "source": "source-service",
  "type": "sh.keptn.event.monitoring.configure",
  "datacontenttype": "application/json",
  "data": {
    "type": "prometheus",
    "project": "crossplane",
    "service": "helloservice"
  },
  "shkeptncontext": "a3e5f16d-8888-4720-82c7-6995062905c1",
  "triggeredid": "3f9640b6-1d2a-4f11-95f5-23259f1d82d6"
}

# configure monitoring 
# there needs to be another configure monitoring for the lighthouse service to define prometheus as SLI provider 
# (maybe this can be done only once as a setup step??)

curl -v -H "Content-Type:application/cloudevents+json" -X POST --data @configure.json http://prometheus-service.keptn:8080

# keptn auth --endpoint=http://34.67.191.73.nip.io/api --api-token=wcw3YyJRBrS2OzziQKIZt3zzUSQMC4oSessgPhgPnwNIy
# keptn configure monitoring prometheus --project=crossplane --service=helloservice




