include .env
export

NETWORKS="$(shell docker network ls)"
VOLUMES="$(shell docker volume ls)"
SUCCESS=[ done "\xE2\x9C\x94" ]

.PHONY: all
all:
	@echo [ starting ostraka... ]
	docker-compose up
	@echo $(SUCCESS)

.PHONY: re
re: down
	@echo [ restarting ostraka... ]
	docker-compose up
	@echo $(SUCCESS)

.PHONY: down
down:
	@echo [ teardown all containers... ]
	docker-compose down
	@echo $(SUCCESS)