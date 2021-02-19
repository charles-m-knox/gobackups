PHONY=build build-image build-dist
PWD=$(shell pwd)
IMG=gobackups-ui
TAG=latest

build:
	make build-image
	make build-dist

build-image:
	docker build \
		-t $(IMG):$(TAG) \
		.

build-dist:
	mkdir -p ./dist
	docker run \
		-it \
		--rm \
		-v "$(PWD)/dist:/site/dist" \
		$(IMG):$(TAG)
