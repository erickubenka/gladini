#!/bin/sh
#Remove autoscaler
kubectl delete hpa selenium-node-chrome-hpa

# Remove Prmoetheus at first
helm delete prometheus
helm delete prometheus-adapter 

# Remove Metrics Server
#kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.6/components.yaml

#Remove custom metrics at first
kubectl delete deployment selex
kubectl delete service selex

# Remove seleniumd Nodes and Grid
kubectl delete deployment selenium-node-chrome
kubectl delete deployment selenium-hub
kubectl delete service selenium-hub

