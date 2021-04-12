# Gladini - A scalable Test Automation Selenium Grid in Kubernetes

## Introduction
Gladini is a simple Selenium 4 grid built with Kubernetes. It supports basic autoscaling features based on external custom metrics  
and Prometheus. It is just a proof of concept, how to scale deployments with Selenium 4.

Gladini uses the current alpha version of Selenium 4 as writing this 4.0.0-beta3. There are plenty of features that worked on Selenium 3 grids, but will currently not work on Selenium 4 grids, therefore Gladini will sneak and trick around them until Selenium provides better API to interact with cluster.

## Prequisites

* Kubernetes
* Kubectl
* Helm 3
* Docker
* GoLang

## Components

### Underlying Selenium Grid
As mentioned in the introduction section, Gladini uses Selenium 4.  
Selenium Grid is deployed in Kubernetes and exposed on `4444` for all Kubernetes cluster communication.

> !!! - IMPORTANT: Node have to expose ports form `4442` to `4444` because Selenium Grid system will use `4442` and `4443` for internal communication.

### Prometheus
Gladini uses Prometheus as well as Prometheus adapter for providing metrics and especially custom metrics to the Kubernetes cluster. 
With the exposed custom metrics Gladini is able to fullfill autoscaling nodes and browser sessions.

### Metrics
Gladini ships with its own Prometheus metrics exporter, called `Selex`, a simple abbreviation for "Selenium exporter".  
Selex will use the Selenium Grid Hub API, provided under http://selenium-hub:NODE_Port/status, and is implemented in a simple `Golang` HTTP Listener. To parse the JSON into objects, Selex provide some structs. 

> Note: As writing this, the Selenium 4 API is hard to use, because sometimes the response-json structure will change, depending on state of the node. To workaround this Selex will use dynamically JSON-map-parsing.

If you want to extend Selex by modifying the Go files, you can build and test it locally.  
For these instructions please have a look into `Dockerfile` in `selex/Dockerfile`.

#### Register Selex on Prometheus
Prometheus and Prometheus Adapter need some configuration to know about our metrics.
First of all we changed the scrape-configs in `prometheus/values.yaml` and add the custom metrics exporter `selex`.
````yaml
scrape_configs:
    - job_name: prometheus
    static_configs:
        - targets:
        - localhost:9090
    - job_name: selex
    scrape_interval: 5s
    static_configs:
        - targets:
        - selex:8080
````

As first step we have to enable custom and external metrics in Prometheus Adapter by changing the values in `prometheus-adapater/values.yaml` in `rules.custom` and `rules.external` by removing the `[]` and commenting in the outcommented values. 
````yaml
rules:
  default: true
  custom: []
# [...]
  external: []
# [...]
````

For chrome autoscaling metrics Selex provides a metric called `selenium_grid_hub_chromeSessionsInUsePercent`, which indicates the current used chrome sessions in percentage to all available crhome sessions. But Selex provide some more metrics, listed below in appendix section.

````yaml
external:
- seriesQuery: '{__name__= "selenium_grid_hub_chromeSessionsInUsePercent"}'
  resources:
      template: <<.Resource>>
  name:
      matches: "selenium_grid_hub_chromeSessionsInUsePercent"
      as: "selenium_grid_hub_chromesessionsinusepercent"
  metricsQuery: sum(<<.Series>>)
````

> Note: See the little change between `external.seriesQuery.name.matches` and `external.seriesQuery.name.as`? It is all about lowercase, because Kubernetes `HorizontalPodAutoScaler` will query the custom and external metrics API with lowercase.

kubectl get --raw /apis/custom.metrics.k8s.io/v1beta1  
kubectl get --raw /apis/external.metrics.k8s.io/v1beta1

### How upscaling works
Gladini uses the provided metrics of Selex to configure a Kubernetes `HorizontalPodAutoscaler` for Selenium browser node deployments.  
The rest of this magical thing is simply part of Kubernetes.

````yaml
metrics:
- type: External
  external:
    metric:
      name: selenium_grid_hub_chromesessionsinusepercent
    target:
      type: Value
      value: 75
````

>Currently only Google Chrome is supported, so the HPA is just available for chrome nodes.

### How downscaling works
Downscaling works the same way as upscaling, but reversed. If the node HPA detects that all metrics are below target, it will remove nodes automatically.  
To provide a graceful shutdown of Selenium nodes, they need to unregister themselves on the Hub/Grid. To fullfill this, the nodes have to call the followings on shutdown.

````bash
# This will store the UUID of the node into nodeid
nodeid=$(curl http://localhost:5555/status | grep id | awk '{print substr($2, 2, 36)}')

# Tell the Hub/Grid to remove this node.
curl -X DELETE http://selenium-hub:4444/se/grid/distributor/node/$nodeid
````

By using this option, you may encounter the problem, that nodes with workload got killed by automatic downscaling, because Kubernetes does not care about running sessions. To avoid this, node will check for workload, when terminitations starts and just wait up to an hour, by configuring the property `terminationGracePeriodSeconds` to value `3600`. 

````bash
while [ $(curl http://localhost:5555/status | jq '.value.node.sessions | length') -ge 1 ]; do sleep 5; echo 'Shutdown requested but session running'; done;
````

Complete teardown for node deployment have to look like this:
````yaml
containers:
- name: selenium-node-chrome
  # other things here...
  lifecycle:
    preStop:
      exec:
        command: [
          "/bin/sh", 
          "-c", 
          "while [ $(curl http://localhost:5555/status | jq '.value.node.sessions | length') -ge 1 ]; do sleep 5; echo 'Shutdown requested but session running'; done; \ 
            nodeid=$(curl http://localhost:5555/status | grep id | awk '{print substr($2, 2, 36)}'); \
            curl -X DELETE http://selenium-hub:4444/se/grid/distributor/node/$nodeid"
          ]
terminationGracePeriodSeconds: 3600
````

## Getting Started

### Starting Gladini

1. Checkout this project.
2. Run `bash helpers/setup.sh`.

These little helper will do the following things:
1. Startup a Selenium 4 Hub deployment
2. Expose the Hub on external port
3. Startup a Selenium 4 Chrome Node deployment
4. Startup a Selex container to expose custom metrics
5. Expose the metrics one external port
6. Installs and configures a Prometheus Helm chart
7. Installs and configures a Prometheus Adapter Helm chart
8. Creates a HorizontalPodAutoscaler for chrome deployment

### Connecting to Grid

#### Seleniumd Grid URL

Gladini will expose underlying Selenium Hub Container on port `30020`.  
1. Use `http://localhost:30020` as your selenium remote grid url.

#### Prometheus Graph UI
To experiment with Prometheus queries, you can run the `bash helpers/open_prometheus.sh` command.  
It will map port `9090` on your local machine.

Then simply call http://localhost:9090/graph

### Shutdown evertyhing

1. Run `bash helpers/teardown.sh`

## Links
#### General Kubernetes Metrics Server
https://github.com/kubernetes-sigs/metrics-server  

#### Custom Metrics Exporter
https://prometheus.io/docs/guides/go-application/  
https://github.com/wakeful/selenium_grid_exporter

#### Prometheus
https://github.com/helm/charts/tree/master/stable/prometheus  
https://github.com/helm/charts/tree/master/stable/prometheus-adapter  
https://github.com/directxman12/k8s-prometheus-adapter

https://github.com/DirectXMan12/k8s-prometheus-adapter/issues/164  
https://blog.kloia.com/kubernetes-hpa-externalmetrics-prometheus-acb1d8a4ed50  
https://itnext.io/horizontal-pod-autoscale-with-custom-metrics-8cb13e9d475  
https://medium.com/@zhimin.wen/custom-prometheus-metrics-for-apps-running-in-kubernetes-498d69ada7aa  

#### Autoscaling
https://learnk8s.io/autoscaling-apps-kubernetes  
https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-multiple-metrics-and-custom-metrics

#### GoCölient Library to talk with Kubectl
https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/  
https://github.com/kubernetes/client-go/tree/master/examples/create-update-delete-deployment  
https://github.com/kubernetes/client-go/tree/release-13.0  
https://github.com/kubernetes/sample-apiserver/blob/master/go.mod  
https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go

## Troubleshooting
When kubernetes liveness probe fails for selenium node pods, it will spawn the node again with same id
Solved by using `Selvidere`

## Jenkins

Jenkins will be started by a default helm chart. 
````
helm init
helm install /stable/jenkins
````

### Configure Jenkins
If you fire up the Jenkins pod at the first time, you have to follow the instructions to set a new password and access default UI of Jenkins. 

### Using Jenkins outside the Cluster
Per default Jenkins is configured to only be accessible inside the cluster.  
If you want to use the Jenkins UI outside of your cluster, just run the script.

````
bash open_jenkins.sh
````
