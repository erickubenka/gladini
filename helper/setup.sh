#!/bin/sh
# Grid Deployment and Service
kubectl create --filename=../selenium-hub-deployment.yaml
kubectl create --filename=../selenium-hub-service.yaml

# Node Deployment and AutoScaler
kubectl create --filename=../selenium-node-chrome-deployment.yaml

# Selex Metrics
kubectl create --filename=../selex-deployment.yaml
kubectl create --filename=../selex-service.yaml

#sleep 5

# Kubernetes Metrics Server
#kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.6/components.yaml

sleep 3

# Prometheus Server and Adapter
helm install prometheus stable/prometheus -f "../prometheus/values.yaml"
helm install prometheus-adapter stable/prometheus-adapter -f "../prometheus-adapter/values.yaml"

# setup autoscaler
kubectl create --filename=../selenium-node-chrome-hpa.yaml