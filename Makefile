PROJECT_DIR = $(shell pwd)
ENV_TEST_PATH = ${PROJECT_DIR}/.env

.PHONY: start
start:
	go run src/cmd/main.go

.PHONY: start_test_server
start_test_server:
	#source $(ENV_TEST_PATH); env
	## go run src/cmd/main.go

.PHONY: gen_docs
gen_docs:
	swag init -g ./src/cmd/main.go -o ./docs --parseDependency --parseInternal --quiet

.PHONY: test
test:
	go test ./... -v

.PHONY: lint
lint:
	golangci-lint run ./... -c .\.golangci-lint.yml

.PHONY: lint_fix
lint_fix:
	golangci-lint run ./... -c .\.golangci-lint.yml --fix