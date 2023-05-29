include .env
export

NETWORKS="$(shell docker network ls)"
VOLUMES="$(shell docker volume ls)"
SUCCESS=[ done "\xE2\x9C\x94" ]

local:
	@echo [ starting ostraka... ]
	go run ./cmd/ostraka/main.go

.PHONY: dev
dev: down
	@echo [ starting ostraka... ]
	docker compose up
	@echo $(SUCCESS)

.PHONY: test
test:
	@echo [ testing ostraka... ]
	docker exec -it ostraka-service sh -c "reflex -r '(\.go$|go\.mod)' -s sh /test.sh"
	@echo $(SUCCESS)

.PHONY: down
down:
	@echo [ teardown all containers... ]
	export TARGET=dev
	docker-compose down
	@echo $(SUCCESS)

.PHONY: sse-example
sse-example:
	@echo [ starting sse-example... ]
	@docker compose up -d
	@sh ./scripts/check_uri.sh POST http://localhost:4000/webhook/orders
	@sh ./examples/sse/events.sh
	@echo $(SUCCESS)