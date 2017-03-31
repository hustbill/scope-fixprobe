#!/bin/sh

# zipkin tunnel (30001 -> remote port 30001)
ssh -i ~/.ssh/k8s-deploy-key.pem  -L 3001:localhost:80  -L 3002:localhost:9411  ubuntu@35.166.123.55 -q -N
