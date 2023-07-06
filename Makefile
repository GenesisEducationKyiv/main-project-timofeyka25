.PHONY: start
start:
	go run src/cmd/main.go

.PHONY: e2e_test
e2e_test:
	go run src/cmd/main.go --test=true & \
        PID=$$!; \
      	go test ./tests/e2e/... -v; \
        kill $$PID

.PHONY: gen_docs
gen_docs:
	swag init -g ./src/cmd/main.go -o ./docs --parseDependency --parseInternal --quiet

.PHONY: test
test:
	go test ./src/... -v

.PHONY: lint
lint:
	golangci-lint run ./... -c .\.golangci-lint.yml

.PHONY: lint_fix
lint_fix:
	golangci-lint run ./... -c .\.golangci-lint.yml --fix