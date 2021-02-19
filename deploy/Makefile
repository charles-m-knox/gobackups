PHONY=build run clean clean-logs-example clean-logs clean-example-files clean-example
PWD=$(shell pwd)
IMG=gobackups-api
TAG=latest
CONTAINER=gobackups
BIND_INTERFACE=127.0.0.1
PORT=12302
EXAMPLE_COMPOSE=../example/docker-compose.example.yml
NORMAL_COMPOSE=./docker-compose.yml
EXAMPLE_PROJECT_NAME=gobackups-example
REAL_PROJECT_NAME=gobackups

build:
	cd ../src; make build
	cd ../ui; make build

run-example:
	# if the log file is not a valid json array,
	# currently the app exits on startup
	./createlog.sh ../example/log-example.json
	docker-compose -p $(EXAMPLE_PROJECT_NAME) -f $(EXAMPLE_COMPOSE) up -d
	docker-compose -p $(EXAMPLE_PROJECT_NAME) -f $(EXAMPLE_COMPOSE) logs -f

run:
	# if the log file is not a valid json array,
	# currently the app exits on startup
	./createlog.sh ./log.json
	docker-compose -p $(REAL_PROJECT_NAME) -f $(NORMAL_COMPOSE) up -d
	docker-compose -p $(REAL_PROJECT_NAME) -f $(NORMAL_COMPOSE) logs -f

	# docker rm -f $(CONTAINER) || true
	# echo "WARNING: This command will mount your entire host file system to"
	# echo "         /host inside the $(CONTAINER) container."
	# echo "         If you do not want this, please edit src/Makefile"
	# echo "         (Waiting 5 seconds before continuing)"
	# sleep 5
	# docker run \
	# 	-d \
	# 	-p "$(BIND_INTERFACE):$(PORT):$(PORT)" \
	# 	--restart=always \
	# 	--name "$(CONTAINER)" \
	# 	-v /:/host/:ro \
	# 	-it $(CONTAINER):latest
	# docker logs -f $(CONTAINER)

down-example:
	docker-compose -p $(EXAMPLE_PROJECT_NAME) -f $(EXAMPLE_COMPOSE) down

down:
	docker-compose -p $(REAL_PROJECT_NAME) -f $(NORMAL_COMPOSE) down

clean-logs:
	-rm log*.json

clean-logs-example:
	-rm ../example/log*.json

clean-example-files:
	-rm -rf ../example/target_dir/*.tar.gz
	-rm -rf ../example/target_dir/*.tar.gz.gpg
	-rm -rf ../example/target_dir/rsync_target/source_dir

clean-example: clean-logs-example clean-example-files

clean: down down-example clean-logs clean-example

logs:
	docker-compose -p $(REAL_PROJECT_NAME) -f $(NORMAL_COMPOSE) logs -f

logs-example:
	docker-compose -p $(EXAMPLE_PROJECT_NAME) -f $(EXAMPLE_COMPOSE) logs -f
