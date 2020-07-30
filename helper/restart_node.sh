#!/bin/sh
kubectl rollout restart deployment selenium-node-chrome
kubectl rollout restart deployment selenium-node-firefox