#!/bin/sh
#Remove custom metrics at first
kubectl delete deployment selex
kubectl delete service selex

# Remove Prmoetheus at first
helm delete prometheus
helm delete prometheus-adapter 

helm del --purge prometheus
helm del --purge prometheus-adapter

# Remove Metrics Server
kubectl delete -f https://github.com/kubernetes-sigs/metrics-server/releases/download/v0.3.6/components.yaml

# Remove seleniumd Nodes and Grid
kubectl delete deployment selenium-node-chrome
kubectl delete deployment selenium-hub
kubectl delete service selenium-hub

