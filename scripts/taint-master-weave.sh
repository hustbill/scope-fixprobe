#!/bin/bash

sudo cp /etc/kubernetes/admin.conf  ~/admin.conf
#kubectl create -f https://git.io/weave-kube
kubectl apply -f https://git.io/weave-kube-1.6
kubectl taint nodes --all node-role.kubernetes.io/master- # single node!
