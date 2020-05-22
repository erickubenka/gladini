#!/bin/sh

# Gladius
helm install gladius ../gladius

sleep 3

# Prometheus Server and Adapter
helm install prometheus stable/prometheus -f "../prometheus/values.yaml"
helm install prometheus-adapter stable/prometheus-adapter -f "../prometheus-adapter/values.yaml"