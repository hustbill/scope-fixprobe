#!/bin/sh

sudo docker run --rm -it \
	  --privileged --net=host --pid=host \
	  -v /lib/modules:/lib/modules \
	  -v /usr/src:/usr/src \
	  -v /sys/kernel/debug/:/sys/kernel/debug/ \
	  -v /var/run/scope/plugins:/var/run/scope/plugins \
	  --name weaveworksplugins-scope-http-statistics weaveworksplugins/scope-http-statistics
