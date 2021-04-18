#!/usr/bin/env bash

cluster_name='seldon-test'
ns_name='seldon-system'

# Set up environment
kind create cluster --name $cluster_name
kubectl create namespace $ns_name
helm install seldon-core seldon-core-operator \
  --namespace $ns_name \
  --repo https://storage.googleapis.com/seldon-charts

# Run test

# Tear down environment
kind delete cluster --name $cluster_name

