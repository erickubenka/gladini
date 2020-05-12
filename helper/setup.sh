#!/bin/sh
# Grid Deployment and Service
kubectl create --filename=../selenium-hub-deployment.yaml
kubectl create --filename=../selenium-hub-service.yaml

# Node Deployment and AutoScaler
kubectl create --filename=../selenium-node-chrome-deployment.yaml

sleep 5

# Kubernetes Metrics Server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.6/components.yaml

sleep 5

# Prometheus Server and Adapter
helm init
helm install --name prometheus stable/prometheus
#helm install  --name prometheus stable/prometheus -f "../prometheus/values.yaml"
#helm install --name prometheus-adapter -f values.yaml stable/prometheus-adapter
helm install --name prometheus-adapter  stable/prometheus-adapter

# Selex Metrics
kubectl create --filename=../selex-deployment.yaml
kubectl create --filename=../selex-service.yaml