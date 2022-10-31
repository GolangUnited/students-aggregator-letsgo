.PHONY:

build-aggregator-image:
	docker build -t aggregator -f ./Dockerfile.aggregator .

build-webserver-image:
	docker build -t web-service -f ./Dockerfile.webservice .

up-aggregator-container: build-aggregator-image
	docker run --name aggregator --detach --rm -p 3306:3306 aggregator	

up-webserver-container: build-webserver-image
	docker run --name webservice --detach --rm -p 80:8080 web-service

up-mongo-container:
	docker run --name mongodb --detach --rm -p 27017:27017 mongo:latest

up:
	docker-compose up --build --detach

down:
	docker-compose down
