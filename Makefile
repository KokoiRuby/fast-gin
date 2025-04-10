PROJECT_NAME = fast-gin
DOCKER_COMPOSE = docker compose -f ./docker/docker-compose.yml -p $(PROJECT_NAME)# set project name

# .PHONY tells Make that 'build' is not a file target
.PHONY: build
build:
	$(DOCKER_COMPOSE) build app

# .PHONY tells Make that 'up' is not a file target
.PHONY: up
# Up and attach to 'app' service logs, 'build' as dependency
up: build
	$(DOCKER_COMPOSE) up --attach app

# .PHONY tells Make that 'down' is not a file target
.PHONY: down
# Remove volumes and locally built images
down:
	$(DOCKER_COMPOSE) down --volumes --rmi=local