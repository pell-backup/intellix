COMETBFT_BUILD_OPTIONS += badgerdb,boltdb,cleveldb,rocksdb

include ../../common.mk

all: docker runner

docker:
	make docker-contracts
	make docker-aggregator
	make docker-operator

DOCKER_CONTRACTS_ENV_FILE = ./dockers/contracts-price-dvs/.env
DOCKER_CONTRACTS_ENV_CONTENT = "ACCOUNT_SECRETKEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

ifeq ($(shell test -e $(DOCKER_CONTRACTS_ENV_FILE) && echo -n yes),yes)
    HAS_DOCKER_CONTRACTS_ENV_FILE=1
else
    HAS_DOCKER_CONTRACTS_ENV_FILE=0
endif

gen-env-file-for-dokcer-contracts:
	@echo "Generating env file for contracts"
	@if [ $($HAS_DOCKER_CONTRACTS_ENV_FILE)=0 ]; then echo "$(DOCKER_CONTRACTS_ENV_CONTENT)" > $(DOCKER_CONTRACTS_ENV_FILE); fi

docker-contracts-build:gen-env-file-for-dokcer-contracts
	@echo "Building contracts"
	cd ./dockers/contracts-price-dvs && docker compose build

docker-contracts-start:gen-env-file-for-dokcer-contracts
	@echo "Starting contracts"
	cd ./dockers/contracts-price-dvs && docker compose up -d

runner:
	go build -o build/runner ./runner

.PHONY: all node docker generator runner docker-contracts
