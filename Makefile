.PHONY:

build-webservice:
	go mod download && go mod tidy
	go build -o ./.bin/webservice ./cmd/webservice/main.go

run-webservice: build-webservice
	./.bin/webservice

build-webserver-image:
	docker build -t web-service -f ./Dockerfile.webservice .

start-webserver-container: build-webserver-image
	docker run -p 80:8080 --name webservice webservice -d --rm

run:
	docker-compose upbuil