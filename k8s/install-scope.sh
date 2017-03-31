#!/bin/sh

sudo curl -L git.io/scope -o /usr/local/bin/scope
sudo chmod a+x /usr/local/bin/scope
scope launch u1604s-k8s-master u1604s-client1 u1604s-client2 u1604s-client3
