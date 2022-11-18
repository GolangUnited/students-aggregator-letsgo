.PHONY:

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
	docker rmi $$(docker images --filter=reference="*/lets_go_aggregator" -aq) 
	docker rmi $$(docker images --filter=reference="*/lets_go_webservice" -aq) 

run: clear
	docker-compose -f ./docker-compose-host.yml up --detach
	
