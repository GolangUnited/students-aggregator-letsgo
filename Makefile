.PHONY:

IS_LETS_GO_AGGREGATOR := $(shell docker images --filter=reference="*/lets_go_aggregator" -aq) 
IS_LETS_GO_WEBSERVICE := $(shell docker images --filter=reference="*/lets_go_webservice" -aq) 

up:
	docker-compose up --build --detach

down:
	docker-compose down

stop:
	docker-compose -f ./docker-compose-host.yml down

clear:
	docker ps --filter name=lets_go_aggregator --filter status=running -aq | xargs docker stop
	docker ps --filter name=lets_go_webservice --filter status=running -aq | xargs docker stop
	docker ps --filter name=lets_go_aggregator --filter status=exited -aq | xargs docker rm
	docker ps --filter name=lets_go_webservice --filter status=exited -aq | xargs docker rm

ifneq ($(strip $(IS_LETS_GO_AGGREGATOR)),)
	docker rmi $$(docker images --filter=reference="*/lets_go_aggregator" -aq)
endif

ifneq ($(strip $(IS_LETS_GO_WEBSERVICE)),)
	docker rmi $$(docker images --filter=reference="*/lets_go_webservice" -aq)
endif

run: clear
	docker-compose -f ./docker-compose-host.yml up --detach
	
