#!/usr/bin/env bash

cluster_name='seldon-test'
infra_ns_name='seldon-system'
deploy_ns_name='seldon'

# Set up environment
kind create cluster --name $cluster_name
kubectl create namespace $infra_ns_name
kubectl create namespace $deploy_ns_name
helm install seldon-core seldon-core-operator \
  --namespace $infra_ns_name \
  --repo https://storage.googleapis.com/seldon-charts

# Run test

# Tear down environment
kind delete cluster --name $cluster_name

