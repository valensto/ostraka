include .env
export

NETWORKS="$(shell docker network ls)"
VOLUMES="$(shell docker volume ls)"
SUCCESS=[ done "✅" ]

local:
## Run ostraka without docker Usage: make local
	@echo [ starting ostraka... ]
	go run ./cmd/ostraka/main.go

.PHONY: dev
## Start all containers Usage: make dev
dev: down
	@echo [ starting ostraka... ]
	docker compose up
	@echo $(SUCCESS)

.PHONY: test
## Run tests Usage: make test
test:
	@echo [ testing ostraka... ]
	@docker exec -it core sh -c "reflex -r '(\.go$|go\.mod)' -s sh /test.sh"
	@echo $(SUCCESS)

.PHONY: down
## Stop all containers Usage: make down
down:
	@echo [ teardown all containers... ]
	docker-compose down
	@echo $(SUCCESS)

.PHONY: clear
## Clear all containers and volumes Usage: make clear
clear: down
	@echo [ stop and clear all containers... ]
	docker image rm -f ostraka-core
	@echo $(SUCCESS)

.PHONY: basic-auth-pwd
## Generate basic auth password and user (default admin) Usage: make -e user=john basic-auth-pwd
basic-auth-pwd:
	@./scripts/htpasswd.sh $(user)

.PHONY: sse-example
## Run sse example Usage: make sse-example
sse-example: down
	@echo [ starting sse-example... ]
	@cp ./examples/sse/sse_order.yaml ./.ostraka/workflows
	docker compose up -d
	@sh ./scripts/check_uri.sh POST ${HOST}:${PORT}/webhook/orders
	@sh ./examples/sse/events.sh
	@echo $(SUCCESS)

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)