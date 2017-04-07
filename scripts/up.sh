#!/bin/bash

systemctl start kubelet
kubeadm init

#kubectl create -f https://git.io/weave-kube
#kubectl apply -f https://git.io/weave-kube-1.6
#kubectl taint nodes --all node-role.kubernetes.io/master- # single node!
