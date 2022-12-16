.PHONY:

IS_LETS_GO_AGGREGATOR_RUNNING := $(shell docker ps --filter name=lets_go_aggregator --filter status=running -aq)
IS_LETS_GO_AGGREGATOR_EXITED := $(shell docker ps --filter name=lets_go_aggregator -aq)
IS_LETS_GO_AGGREGATOR := $(shell docker images --filter=reference="*/lets_go_aggregator" -aq)
IS_LETS_GO_WEBSERVICE_RUNNING := $(shell docker ps --filter name=lets_go_webservice --filter status=running -aq)
IS_LETS_GO_WEBSERVICE_EXITED := $(shell docker ps --filter name=lets_go_webservice -aq)
IS_LETS_GO_WEBSERVICE := $(shell docker images --filter=reference="*/lets_go_webservice" -aq)

up:
	docker-compose up --build --detach

down:
	docker-compose down

stop:
	docker-compose -f ./docker-compose-host.yml down

clear:

ifneq ($(strip $(IS_LETS_GO_AGGREGATOR_RUNNING)),)
	docker stop $(IS_LETS_GO_AGGREGATOR_RUNNING)
endif

ifneq ($(strip $(IS_LETS_GO_AGGREGATOR_EXITED)),)
	docker rm $(IS_LETS_GO_AGGREGATOR_EXITED)
endif

ifneq ($(strip $(IS_LETS_GO_AGGREGATOR)),)
	docker rmi $(IS_LETS_GO_AGGREGATOR)
endif

ifneq ($(strip $(IS_LETS_GO_WEBSERVICE_RUNNING)),)
	docker stop $(IS_LETS_GO_WEBSERVICE_RUNNING)
endif

ifneq ($(strip $(IS_LETS_GO_WEBSERVICE_EXITED)),)
	docker rm $(IS_LETS_GO_WEBSERVICE_EXITED)
endif

ifneq ($(strip $(IS_LETS_GO_WEBSERVICE)),)
	docker rmi $(IS_LETS_GO_WEBSERVICE)
endif

run: clear
	cat ./.env.ci
	docker-compose --env-file ./.env.ci -f ./docker-compose-host.yml up --detach

cover:
	go test -v -coverprofile cover.out ./...
	go tool cover -html cover.out -o cover.html

swagger_install:
	go install github.com/swaggo/swag/cmd/swag@latest

swagger_init:
	swag init --generalInfo main.go --dir ./cmd/webservice,./internal/webservice,./model --output ./cmd/webservice/docs
