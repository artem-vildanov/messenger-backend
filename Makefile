TAG=$(shell git rev-parse --short HEAD)

# network as a paramether
define RUN_MIGRATOR
	docker build \
		--build-arg ENV=$(2) \
		-f ./deployment/migrator.Dockerfile \
		-t messenger-migrator .
	docker run --rm \
		--name messenger-migrator \
		--network $(1) \
		messenger-migrator
endef

# build service image
build-service:
	docker build \
		-f ./deployment/service.Dockerfile \
		-t messenger-service:$(TAG) .

# run whole application
run: build-service
	IMAGE_TAG=$(TAG) docker stack deploy \
		-c ./deployment/docker-compose.yml \
		messenger-app

# hot reload of service
update-service: build-service
	docker service update \
		--image messenger-service:$(TAG) \
		messenger-app_messenger-service

# stop whole application
stop:
	docker stack rm messenger-app

# run migrator
migrate:
	$(call RUN_MIGRATOR,messenger-app_messenger-network,dev)

# build binaries with debug flags
build-debug:
	go build \
		-gcflags "all=-N -l" \
		-o ./bin/service ./cmd/service/main.go 

int-tests-ci:
	make int-tests
	make int-tests-cleanup

int-tests:
	docker compose \
		-f ./deployment/docker-compose.test.yml \
		-p messenger-app \
		up -d
	$(call RUN_MIGRATOR,messenger-app_messenger-network-test,test)
	docker exec \
		-w /app \
		messenger-app-messenger-service-test-1 \
		./bin/tests -test.v

int-tests-cleanup:
	docker compose \
		-f ./deployment/docker-compose.test.yml \
		-p messenger-app \
		down -v
	docker rm -f messenger-app-messenger-service-test-1
	docker image rm -f messenger-app-messenger-service-test

rerun-int-tests:
	docker rm -f messenger-app-messenger-service-test-1
	docker image rm -f messenger-app-messenger-service-test
	docker compose -f ./deployment/docker-compose.test.yml \
		-p messenger-app \
		up messenger-service-test -d
	docker exec \
		-w /app \
		messenger-app-messenger-service-test-1 \
		./bin/tests -test.v

run-concrete-int-test:
	docker rm -f messenger-app-messenger-service-test-1
	docker image rm -f messenger-app-messenger-service-test
	docker compose -f ./deployment/docker-compose.test.yml \
		-p messenger-app \
		up messenger-service-test -d
	docker exec \
		-w /app \
		messenger-app-messenger-service-test-1 \
		./bin/tests -test.v -test.run ^Test_WS$
