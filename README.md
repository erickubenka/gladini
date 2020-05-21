# Gladius - A Test Automation Platform in Kubernetes

## Introduction

## Prequisites

* Kubectl
* Helm 3
* Docker
* GoLang

## Selenium Grid
We will use the current alpha version of selenium to support bleeding edge features and get in touch with current bugs as soon as possible.

Selenium Grid is deployed in Kubernetes and exposed on `localhost:4444`.

!Important: We have to expose ports form `4442` to `4444` because Selenium Grid system will use `4442` and `4443` for internal communication.

### Starting Grid

1. Run `bash helpers/setup.sh`.

### Connecting to Grid

1. Get the current port of Selenium Hub in local Kubernetes cluster
````
kubectl describe service selenium-hub
````
2. Connect to Grid via `localhost:NODE_PORT`

## Prometheus

We will use the default kubernetes metrics server and Prometheus as well as Prometheus adapter for creating custom metrics. 
With this approach we want to enable a autoscaling selenium hub.

### Custom Metrics
With `selex` we create a simple custom metric exporter for Selenium grid by using the Selnium grid API provided under `http://localhost:30020/status`implemented in a simple `golang` HTTP Listener. Therefore we implemented some structs to parse the JSON response into objects. 

To build things up, see the provided `Dockerfile` in `selex/Dockerfile`.

Run these commands to build the Docker container locallly, if you want to:
````bash
cd selex/
docker build --tag selex:1.0 .
docker run --publish 8080:8080 --detach --name selex selex:1.0 
# Verify it works
docker stop selex
docker rm selex
````

Prometheus and Prometheus Adapter need some configuration to know about our metrics.
First of all we changed the scrape-configs in `prometheus/values.yaml` and add the custom metrics exporter `selex`.
````yaml
scrape_configs:
    - job_name: prometheus
    static_configs:
        - targets:
        - localhost:9090
    - job_name: selex
    scrape_interval: 10s
    static_configs:
        - targets:
        - selex:8080
````

As first step we have to enable custom and external metrics in Prometheus Adapert by chaning the values in `prometheus-adapater/values.yaml` in `rules.custom` and `rules.external` by removing the `[]` and commenting in the outcommented values. 
````yaml
rules:
  default: true
  custom: []
# [...]
  external: []
# [...]
````

### Links
#### Metrics Server
https://github.com/kubernetes-sigs/metrics-server

#### Prometheus
https://github.com/helm/charts/tree/master/stable/prometheus
https://github.com/helm/charts/tree/master/stable/prometheus-adapter

https://github.com/directxman12/k8s-prometheus-adapter

#### Autoscaling Guide
https://learnk8s.io/autoscaling-apps-kubernetes
https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/#autoscaling-on-multiple-metrics-and-custom-metrics

#### A selenium 3 exporter for prometheus
https://prometheus.io/docs/guides/go-application/
https://github.com/wakeful/selenium_grid_exporter



https://kubernetes.io/docs/reference/kubectl/docker-cli-to-kubectl/

##############################
##############################
Next steps:
https://github.com/zalando-incubator/kube-metrics-adapter
https://github.com/kubernetes-sigs/metrics-server/issues/131
https://github.com/DirectXMan12/k8s-prometheus-adapter/issues/164
https://blog.kloia.com/kubernetes-hpa-externalmetrics-prometheus-acb1d8a4ed50
https://itnext.io/horizontal-pod-autoscale-with-custom-metrics-8cb13e9d475
https://www.magalix.com/blog/the-adapter-pattern

Solved
Downscaled nodes are markes as offline in hub, but they simple should be removed....
curl -X DELETE http://localhost:30020/se/grid/distributor/node/83d72fdc-5c06-44fb-82ec-2c2f9bead7c2

Think about downscaling - do not downscale when session is present..


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
bash open_j.sh
````

#### ADMIN PASS
Oa3IDD66bJ
