#!/bin/sh
 export POD_NAME=$(kubectl get pods --namespace default -l "app.kubernetes.io/component=jenkins-master" -l "app.kubernetes.io/instance=ornery-moth" -o jsonpath="{.items[0].metadata.name}")
 kubectl --namespace default port-forward $POD_NAME 8080:8080