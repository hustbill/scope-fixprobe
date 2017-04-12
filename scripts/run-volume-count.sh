#!/bin/sh

sudo docker run --rm -ti \
	--net=host \
	-v /var/run/scope/plugins:/var/run/scope/plugins \
	--name weaveworksplugins-scope-volume-count gallna/scope-volume-count:latest
