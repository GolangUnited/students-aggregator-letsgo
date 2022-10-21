.PHONY:

build-webservice:
	go mod download && go mod tidy
	go build -o ./.bin/webservice ./cmd/webservice/main.go

run-webservice: build-webservice
	./.bin/webservice

build-webserver-image:
	docker build -t web-service -f ./Dockerfile.webservice .

up-webserver-container: build-webserver-image
	docker run --name web-service --detach --rm -p 80:8080 web-service

up:
	docker-compose up --build --detach

down:
	docker-compose down
