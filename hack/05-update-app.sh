#!/bin/bash


kubectl -n frank delete -f ../k8s/app.yaml

kubectl -n frank apply -f ../k8s/app.yaml

kubectl -n kube-system rollout restart deployment/coredns