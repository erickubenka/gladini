#!/bin/sh
kubectl delete deployment selenium-node-firefox
kubectl delete deployment selenium-node-chrome
kubectl delete service selenium-hub
kubectl delete deployment selenium-hub
