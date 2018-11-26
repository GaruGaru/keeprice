BIN_NAME=keeprice
BIN_OUTPUT=dist/${BIN_NAME}

fmt:
	go fmt ./...

deps:
	go mod vendor
	go mod verify

build: fmt deps
	go build -o ${BIN_OUTPUT}


COMPOSE=docker/docker-compose.yml
DOCKER_IMAGE=garugaru/keeprice

docker-up:
	docker-compose -f ${COMPOSE} up

docker-upd:
	docker-compose -f ${COMPOSE} up -d

docker-down:
	docker-compose -f ${COMPOSE} down

docker-build:
	docker-compose -f ${COMPOSE} build

docker-push: docker-build
	docker push ${DOCKER_IMAGE}