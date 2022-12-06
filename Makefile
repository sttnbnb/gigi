include .env

DOCKER_FILE=./build/Dockerfile
ENV_FILE=.env
APP_NAME=gigi

docker/run: docker/rm docker/build
	docker run -d --name $(APP_NAME) --env-file=$(ENV_FILE) --restart=always gigi:latest

docker/start:
	docker start $(APP_NAME)

docker/stop:
	docker stop $(APP_NAME)

docker/rm:
	docker rm -f $(APP_NAME)

docker/build:
	docker build -t $(APP_NAME) -f $(DOCKER_FILE) .

docker/logs:
	docker logs -f $(APP_NAME)
