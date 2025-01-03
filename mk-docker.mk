check-env-gh-token:
	@if [ -z "$${GITHUB_TOKEN}" ] && ! grep -q '^GITHUB_TOKEN=' docker/.env 2>/dev/null; then \
		echo "Error: GITHUB_TOKEN is not set in environment or docker/.env file"; \
		exit 1; \
	else \
		echo "GITHUB_TOKEN is set."; \
	fi

docker-build: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build

docker-build-pelldvs: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build pelldvs

docker-build-operator: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build operator

docker-up:
	@cd docker && docker compose down -v && docker compose up -d

docker-logs:
	@cd docker && docker compose logs -f

docker-down:
	@cd docker && docker compose down -v

docker-test:
	@cd docker && docker compose run --rm test /root/scripts/test.sh

docker-contracts-logs-eth:
	@cd docker && docker compose logs -f eth

docker-contracts-logs-hardhat:
	@cd docker && docker compose logs -f hardhat

docker-contracts-bash-eth:
	@cd docker && docker compose exec eth sh

docker-contracts-bash-hardhat:
	@cd docker && docker compose exec hardhat bash
