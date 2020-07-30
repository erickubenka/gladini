#!/bin/sh
kubectl create --filename=../selenium-hub-deployment.yaml
kubectl create --filename=../selenium-hub-service.yaml
kubectl create --filename=../selenium-node-chrome-deployment.yaml
kubectl create --filename=../selenium-node-firefox-deployment.yaml
