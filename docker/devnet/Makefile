default:
	@echo "hello"

build-pelldvs:
	docker compose -p intellix-testnet -f ../docker-compose.build.yml build pelldvs

build-operator:
	docker compose -p intellix-testnet -f ../docker-compose.build.yml build operator

down-all:
	make operator-down-operator;
	#make operator-down-dispatcher;
	make cosmos-down-abci;
	make dvs-down-all;

dvs-up-all:
	make dvs-up-dvs;
	make dvs-up-gateway;

dvs-up-dvs:
	docker compose --env-file .env -f docker-compose.dvs.yml up dvs -d

dvs-up-gateway:
	docker compose --env-file .env -f docker-compose.dvs.yml up gateway -d

dvs-logs-dvs:
	docker compose -f docker-compose.dvs.yml logs -f dvs -n 100

dvs-logs-gateway:
	docker compose -f docker-compose.dvs.yml logs -f gateway -n 100

dvs-down-all:
	@make dvs-down-dvs;
	@make dvs-down-gateway;

dvs-down-dvs:
	docker compose -f docker-compose.dvs.yml down -v dvs

dvs-down-gateway:
	docker compose -f docker-compose.dvs.yml down -v gateway

dvs-shell-dvs:
	docker compose --env-file .env -f docker-compose.dvs.yml exec -it dvs bash

dvs-shell-gateway:
	docker compose --env-file .env -f docker-compose.dvs.yml exec -it gateway bash

dvs-rm-dvs:
	docker compose --env-file .env -f docker-compose.dvs.yml rm dvs -f

dvs-rm-gateway:
	docker compose --env-file .env -f docker-compose.dvs.yml rm gateway -f

dvs-dvs-rerun:
	make dvs-down-dvs;
	make dvs-up-dvs;
	make dvs-logs-dvs;

dvs-gateway-rerun:
	make dvs-down-gateway;
	make dvs-up-gateway;
	make dvs-logs-gateway;

cosmos-up-abci:
	docker compose --env-file .env -f docker-compose.cosmos.yml up abci -d

cosmos-down-abci:
	docker compose -f docker-compose.cosmos.yml down abci -v

cosmos-logs-abci:
	docker compose --env-file .env -f docker-compose.cosmos.yml logs abci -n 100 -f

cosmos-abci-rerun:
	make cosmos-down-abci;
	make cosmos-up-abci;
	make cosmos-logs-abci;

# # Operator dispatcher
#operator-up-dispatcher:
#	docker compose --env-file .env -f docker-compose.yml up dispatcher -d
#
#operator-down-dispatcher:
#	docker compose -f docker-compose.yml down dispatcher -v
#
#operator-rm-dispatcher:
#	docker compose --env-file .env -f docker-compose.yml rm dispatcher -f
#
#operator-logs-dispatcher:
#	docker compose --env-file .env -f docker-compose.yml logs dispatcher -n 100 -f
#
#operator-dispatcher-rerun:
#	make operator-down-dispatcher;
#	make operator-up-dispatcher;
#	make operator-logs-dispatcher;


# Operator Operator
operator-up-operator:
	docker compose --env-file .env -f docker-compose.yml up operator -d

operator-down-operator:
	docker compose -f docker-compose.yml down operator -v

operator-rm-operator:
	docker compose --env-file .env -f docker-compose.yml rm operator -f

operator-logs-operator:
	docker compose --env-file .env -f docker-compose.yml logs operator -n 100 -f

operator-operator-rerun:
	make operator-down-operator;
	make operator-up-operator;
	make operator-logs-operator;

operator-up-operator01:
	docker compose --env-file .env -f docker-compose.yml up operator01 -d

operator-up-operator02:
	docker compose --env-file .env -f docker-compose.yml up operator02 -d

operator-up-all:
	@make operator-up-operator01;
	@make operator-up-operator02;

operator-down-operator01:
	docker compose -f docker-compose.yml down -v

operator-down-operator02:
	docker compose -f docker-compose.yml down -v

operator-down-all:
	@make operator-down-operator01;
	@make operator-down-operator02;

operator-logs-operator01:
	docker compose -f docker-compose.yml logs -f operator01 -n 100

operator-logs-operator02:
	docker compose -f docker-compose.yml logs -f operator02 -n 100

