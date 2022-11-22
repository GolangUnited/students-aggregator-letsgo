.PHONY:

IS_LETS_GO_AGGREGATOR := $(shell docker images --filter=reference="*/lets_go_aggregator" -aq) 
IS_LETS_GO_WEBSERVICE := $(shell docker images --filter=reference="*/lets_go_webservice" -aq) 
IS_LETS_GO_AGGREGATOR_RUNNING := $(shell docker ps --filter name=lets_go_aggregator --filter status=running -aq) 
IS_LETS_GO_AGGREGATOR_EXITED := $(shell docker ps --filter name=lets_go_aggregator --filter status=exited -aq) 
IS_LETS_GO_WEBSERVICE_RUNNING := $(shell docker ps --filter name=lets_go_webservice --filter status=running -aq) 
IS_LETS_GO_WEBSERVICE_EXITED := $(shell docker ps --filter name=lets_go_webservice --filter status=exited -aq) 

up:
	docker-compose up --build --detach

down:
	docker-compose down

stop:
	docker-compose -f ./docker-compose-host.yml down

clear:

ifneq ($(strip $(IS_LETS_GO_AGGREGATOR_RUNNING)),)
	docker rmi $$(IS_LETS_GO_AGGREGATOR_RUNNING)
endif

ifneq ($(strip $(IS_LETS_GO_AGGREGATOR_EXITED)),)
	docker rmi $$(IS_LETS_GO_AGGREGATOR_EXITED)
endif

ifneq ($(strip $(IS_LETS_GO_WEBSERVICE_RUNNING)),)
	docker rmi $$(IS_LETS_GO_WEBSERVICE_RUNNING)
endif

ifneq ($(strip $(IS_LETS_GO_WEBSERVICE_EXITED)),)
	docker rmi $$(IS_LETS_GO_WEBSERVICE_EXITED)
endif

ifneq ($(strip $(IS_LETS_GO_AGGREGATOR)),)
	docker rmi $$(IS_LETS_GO_AGGREGATOR)
endif

ifneq ($(strip $(IS_LETS_GO_WEBSERVICE)),)
	docker rmi $$(IS_LETS_GO_WEBSERVICE)
endif

run: clear
	docker-compose -f ./docker-compose-host.yml up --detach
	
