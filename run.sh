#!/usr/bin/env bash

cluster_name='seldon-test'

# Set up environment
kind create cluster --name $cluster_name
kubectl create namespace seldon-system

# Run test

# Tear down environment
#kind delete cluster --name $cluster_name

