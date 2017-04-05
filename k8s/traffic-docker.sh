docker run --rm -it \
			 --net=host --pid=host --privileged \
			 -v /var/run:/var/run \
			 --name weaveworksplugins-scope-traffic-control weaveworksplugins/scope-traffic-control:latest
