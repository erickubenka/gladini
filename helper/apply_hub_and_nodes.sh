#!/bin/sh
kubectl apply --filename=../selenium-hub-deployment.yaml
kubectl apply --filename=../selenium-node-chrome-deployment.yaml
