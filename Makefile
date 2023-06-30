.PHONY: start
start:
	go run src/cmd/main.go

.PHONY: gen_docs
gen_docs:
	swag init -g ./src/cmd/main.go -o ./docs --parseDependency --parseInternal --quiet