#!/bin/sh

# scope tunnel (30040 -> remote port 4040)
ssh -i ~/.ssh/k8s-deploy-key.pem  -L 30040:localhost:4040  ubuntu@35.166.123.55 -q -N
