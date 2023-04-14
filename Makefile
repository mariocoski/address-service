SHELL := /bin/bash

# Define directory and file names
MIGRATIONS_DIR := liquibase/development/migrations/changelogs
SEEDS_DIR := liquibase/development/seeds/changelogs
MIGRATION_NAME := $(strip $(filter-out $@,$(MAKECMDGOALS))) 
TIMESTAMP := $(shell date +%Y%m%d%H%M%S)
FILE_NAME := $(TIMESTAMP)_$(MIGRATION_NAME).sql

# Define the migration template
define MIGRATION_TEMPLATE
--liquibase formatted sql 
--changeset user:${TIMESTAMP}-${MIGRATION_NAME} splitStatements:false

-- SQL statements here

--rollback 

endef
export MIGRATION_TEMPLATE
 
# "Usage: make migration MIGRATION_NAME=<migration_name>"
.PHONY: migration
migration:
	@echo "$$MIGRATION_TEMPLATE" > $(MIGRATIONS_DIR)/$(FILE_NAME)

.PHONY: migration-update
migration-update: ## apply migrate
	liquibase --changelog-file=liquibase/development/migrations/changelog.yaml update

# "Usage: make seed MIGRATION_NAME=<migration_name>"
.PHONY: seed
seed:
	@echo "$$MIGRATION_TEMPLATE" > $(SEEDS_DIR)/$(FILE_NAME)

.PHONY: seed-update
seed-update:
	liquibase --database-changelog-lock-table-name=seedschangeloglock --database-changelog-table-name=seedschangelog --changelog-file=liquibase/development/seeds/changelog.yaml update

.PHONY: build
build:
	go build -o bin/main cmd/http/main.go

.PHONY: run
run:
	go run cmd/http/main.go

.PHONY: install
install:
	sh ./scripts/install-liquibase-mac.sh
