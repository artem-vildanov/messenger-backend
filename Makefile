build-service:
	@docker build -f ./deployment/service.Dockerfile -t messenger-service .

run: build-service
	@docker stack deploy -c ./deployment/docker-compose.yml messenger-app

update-service: build-service
	@docker service update --image messenger-service messenger-app_messenger-service

stop:
	@docker stack rm messenger-app

migrate:
	@docker build -f ./deployment/migrator.Dockerfile -t messenger-migrator .
	@docker run --rm --name messenger-migrator --network messenger-app_messenger-network messenger-migrator

build-debug:
	@go build -gcflags "all=-N -l" -o ./bin/service ./cmd/service/main.go 
