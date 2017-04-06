#!/bin/sh

scope launch --mode app --probe.docker=true --no-probe --app.window 8760h

curl -X POST -H "Content-Type: application/json" http://192.168.255.97:4040/api/report -d @/home/hzhang/git/scope/docker/demo.json
