check-env-gh-token:
	@if [ -z "$${GITHUB_TOKEN}" ] && ! grep -q '^GITHUB_TOKEN=' docker/.env 2>/dev/null; then \
		echo "Error: GITHUB_TOKEN is not set in environment or docker/.env file"; \
		exit 1; \
	else \
		echo "GITHUB_TOKEN is set."; \
	fi

docker-build-all: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build

DOCKER_TAG ?= latest

docker-build-release-local: check-env-gh-token
	docker build -f ./Dockerfile-local -t pellnetwork/intellix:${DOCKER_TAG} . \
	--build-arg GITHUB_TOKEN=$${GITHUB_TOKEN} \
	--build-arg PELLDVS_VERSION=v0.2.2

docker-build-contracts: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build hardhat

docker-build-pelldvs: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build pelldvs

docker-build-operator: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build operator

docker-all-up:
	@cd docker && docker compose up -d

docker-all-down:
	@cd docker && docker compose down -v

docker-all-status:
	@cd docker && docker compose ps -a

docker-hardhat-up:
	@cd docker && docker compose up hardhat -d

docker-hardhat-down:
	@cd docker && docker compose down hardhat -v

docker-hardhat-logs:
	@cd docker && docker compose logs hardhat -f

docker-hardhat-shell:
	@cd docker && docker compose exec -it hardhat bash

docker-hardhat-rerun:
	make docker-hardhat-down
	make docker-hardhat-up
	make docker-hardhat-logs

docker-emulator-up:
	@cd docker && docker compose up emulator -d

docker-emulator-down:
	@cd docker && docker compose down emulator -v

docker-emulator-logs:
	@cd docker && docker compose logs -f emulator

docker-emulator-shell:
	@cd docker && docker compose exec -it emulator bash

docker-emulator-rerun:
	make docker-emulator-down
	make docker-emulator-up
	make docker-emulator-logs

docker-dvs-up:
	@cd docker && docker compose up dvs -d

docker-dvs-down:
	@cd docker && docker compose down dvs -v

docker-dvs-logs:
	@cd docker && docker compose logs dvs -f

docker-dvs-shell:
	@cd docker && docker compose exec -it dvs bash

docker-dvs-rerun:
	make docker-dvs-down
	make docker-dvs-up
	make docker-dvs-logs

docker-gateway-up:
	@cd docker && docker compose up gateway -d

docker-gateway-down:
	@cd docker && docker compose down gateway -v

docker-gateway-logs:
	@cd docker && docker compose logs gateway -f

docker-gateway-shell:
	@cd docker && docker compose exec -it gateway bash

docker-gateway-rerun:
	make docker-gateway-down
	make docker-gateway-up
	make docker-gateway-logs

docker-operator-up:
	@cd docker && docker compose up operator -d

docker-operator-down:
	@cd docker && docker compose down operator -v

docker-operator-logs:
	@cd docker && docker compose logs operator -f

docker-operator-shell:
	@cd docker && docker compose exec -it operator bash

docker-operator-rerun:
	make docker-operator-down
	make docker-operator-up
	make docker-operator-logs

docker-test:
	@cd docker && docker compose run --rm test

docker-all-logs-in-ci:
	@cd docker && \
		docker compose logs emulator -n 50 && \
		echo -e "\n\n\t==================== emulator logs end \n\n" && \
		docker compose logs hardhat -n 50 && \
		echo -e "\n\n\t==================== hardhat logs end \n\n" && \
		docker compose logs dvs -n 50 && \
		echo -e "\n\n\t==================== dvs logs end \n\n" && \
		docker compose logs gateway -n 50 && \
		echo -e "\n\n\t==================== gateway logs end \n\n" && \
		docker compose logs operator -n 50 && \
		echo -e "\n\n\t==================== operator logs end \n\n"
