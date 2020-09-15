BIN = $(GOPATH)/bin
BASE = $(GOPATH)/src/$(PACKAGE)
PKGS = go list ./... | grep -v "^vendor/"

# Tools
## Testing library
GINKGO = $(BIN)/ginkgo
$(BIN)/ginkgo:
	go get -u github.com/onsi/ginkgo/ginkgo

## Migration tool
GOOSE = go run -tags='no_mysql no_sqlite3 no_mssql no_redshift' github.com/pressly/goose/cmd/goose

## Source linter
LINT = $(BIN)/golint
$(BIN)/golint:
	go get -u golang.org/x/lint/golint

## Combination linter
METALINT = $(BIN)/gometalinter.v2
$(BIN)/gometalinter.v2:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	$(METALINT) --install


.PHONY: installtools
installtools: | $(LINT) $(GINKGO)
	echo "Installing tools"

.PHONY: metalint
metalint: | $(METALINT)
	$(METALINT) ./... --vendor \
	--fast \
	--exclude="exported (function)|(var)|(method)|(type).*should have comment or be unexported" \
	--format="{{.Path.Abs}}:{{.Line}}:{{if .Col}}{{.Col}}{{end}}:{{.Severity}}: {{.Message}} ({{.Linter}})"

.PHONY: lint
lint:
	$(LINT) $$($(PKGS)) | grep -v -E "exported (function)|(var)|(method)|(type).*should have comment or be unexported"

#Test
TEST_DB = vulcanize_testing
TEST_CONNECT_STRING = postgresql://localhost:5432/$(TEST_DB)?sslmode=disable

.PHONY: test
test: | $(GINKGO) $(LINT)
	go vet ./...
	go fmt ./...
	dropdb --if-exists $(TEST_DB)
	createdb $(TEST_DB)
	psql -q $(TEST_DB) < test_data/vulcanize_schema.sql
	make migrate NAME=$(TEST_DB)
	make reset NAME=$(TEST_DB)
	make migrate NAME=$(TEST_DB)
	$(GINKGO) -r --skipPackage=integration_tests,integration

.PHONY: integrationtest
integrationtest: | $(GINKGO) $(LINT)
	go vet ./...
	go fmt ./...
	dropdb --if-exists $(TEST_DB)
	createdb $(TEST_DB)
	psql -q $(TEST_DB) < test_data/vulcanize_schema.sql
	make migrate NAME=$(TEST_DB)
	make reset NAME=$(TEST_DB)
	make migrate NAME=$(TEST_DB)
	$(GINKGO) -r transformers/integration_tests/

.PHONY: validatemigrationorder
validatemigrationorder: vdb-mcd-transformers
	./vdb-mcd-transformers checkMigrations

vdb-mcd-transformers:
	go build

# Build is really "clean/rebuild"
.PHONY: build
build:
	- rm vdb-mcd-transformers
	go fmt ./...
	go build

#Database
HOST_NAME = localhost
PORT = 5432
NAME =
CONNECT_STRING=postgresql://$(HOST_NAME):$(PORT)/$(NAME)?sslmode=disable

# Parameter checks
## Check that DB variables are provided
.PHONY: checkdbvars
checkdbvars:
	test -n "$(HOST_NAME)" # $$HOST_NAME
	test -n "$(PORT)" # $$PORT
	test -n "$(NAME)" # $$NAME
	@echo $(CONNECT_STRING)

## Check that the migration variable (id/timestamp) is provided
.PHONY: checkmigration
checkmigration:
	test -n "$(MIGRATION)" # $$MIGRATION

## Check that the migration name is provided
.PHONY: checkmigname
checkmigname:
	test -n "$(NAME)" # $$NAME

# Migration operations
## Rollback the last migration
.PHONY: rollback
rollback: checkdbvars
	cd db/migrations;\
	  $(GOOSE) -table "maker.goose_db_version" postgres "$(CONNECT_STRING)" down
	pg_dump -n 'maker' -n 'api' -O -s $(CONNECT_STRING) > db/schema.sql


## Rollback to a select migration (id/timestamp)
.PHONY: rollback_to
rollback_to: checkmigration checkdbvars
	cd db/migrations;\
	  $(GOOSE) -table "maker.goose_db_version" postgres "$(CONNECT_STRING)" down-to "$(MIGRATION)"

## Apply all migrations not already run
.PHONY: migrate
migrate: checkdbvars
	psql $(NAME) -c 'CREATE SCHEMA IF NOT EXISTS maker;'
	cd db/migrations;\
	  $(GOOSE) -table "maker.goose_db_version" postgres "$(CONNECT_STRING)" up
	pg_dump -n 'maker' -n 'api' -O -s $(CONNECT_STRING) > db/schema.sql

.PHONY: reset
reset: checkdbvars
	cd db/migrations/;\
		$(GOOSE) -table "maker.goose_db_version" postgres "$(CONNECT_STRING)" reset
	pg_dump -n 'maker' -n 'api' -O -s $(CONNECT_STRING) > db/schema.sql
	psql $(NAME) -c 'DROP SCHEMA maker CASCADE;'

## Create a new migration file
.PHONY: new_migration
new_migration: checkmigname
	cd db/migrations;\
	  $(GOOSE) create $(NAME) sql

## Check which migrations are applied at the moment
.PHONY: migration_status
migration_status: checkdbvars
	cd db/migrations;\
	  $(GOOSE) postgres "$(CONNECT_STRING)" status

# Convert timestamped migrations to versioned (to be run in CI);
# merge timestamped files to prevent conflict
.PHONY: version_migrations
version_migrations:
	cd db/migrations; $(GOOSE) fix

# Import a psql schema to the database
.PHONY: import
import:
	test -n "$(NAME)" # $$NAME
	psql $(NAME) < db/schema.sql

# Build plugin
.PHONY: plugin
plugin:
	go build -buildmode=plugin -o $(OUTPUT_LOCATION) $(TARGET_LOCATION)

# Docker actions
# Build any docker image in dockerfiles
.PHONY: dockerbuild
dockerbuild:
	test -n "$(IMAGE)" # $$IMAGE
	docker build -t $(IMAGE) -f dockerfiles/$(IMAGE)/Dockerfile .

.PHONY: execute
execute: HOST ?= host.docker.internal
execute: DATABASE_PASSWORD ?= postgres
execute:
	test -n "$(NAME)" # $$(NAME) - Database Name
	docker run \
		-it \
		-p "5432:5432" \
		-e "DATABASE_NAME=$(NAME)" \
		-e "DATABASE_HOSTNAME=$(HOST)" \
		-e "DATABASE_PORT=$(PORT)" \
		-e "DATABASE_USER=$(USER)" \
		-e "DATABASE_PASSWORD=$(DATABASE_PASSWORD)" \
		execute:latest
