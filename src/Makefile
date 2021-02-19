PHONY=build run
PWD=$(shell pwd)
IMG=gobackups-api
TAG=latest
CONTAINER=gobackups
BIND_INTERFACE=127.0.0.1
PORT=12403

build:
	docker build \
		-t $(IMG):$(TAG) \
		.
run:
	docker rm -f $(CONTAINER) || true
	echo "WARNING: This command will mount your entire host file system to"
	echo "         /host inside the $(CONTAINER) container."
	echo "         If you do not want this, please edit src/Makefile"
	echo "         (Waiting 5 seconds before continuing)"
	sleep 5
	docker run \
		-d \
		-p "$(BIND_INTERFACE):$(PORT):$(PORT)" \
		--restart=always \
		--name "$(CONTAINER)" \
		-v /:/host/:ro \
		-it $(CONTAINER):latest
	docker logs -f $(CONTAINER)
