#!/bin/bash

systemctl start kubelet
kubeadm init
sudo cp /etc/kubernetes/admin.conf  ~/admin.conf
kubectl apply -f https://git.io/weave-kube-1.6
kubectl taint nodes --all node-role.kubernetes.io/master- # single node!
# kubectl create -f https://git.io/weave-kube  1.5
