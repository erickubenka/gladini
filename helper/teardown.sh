#!/bin/sh
# Remove Prometheus at first
helm delete prometheus
helm delete prometheus-adapter 

# Remove Metrics Server
helm uninstall gladius
