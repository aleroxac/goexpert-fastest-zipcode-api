conf ?= cmd/server/.env
include $(conf)
export $(shell sed 's/=.*//' $(conf))



## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@rm -f coverage.* cmd/server/.env



## ---------- SETUP
.PHONY: install
install: ## install requirements
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go mod tidy



## ---------- FMT && VET && DOCS
.PHONY: fmt
fmt: ## format the code
	@go fmt ./...

.PHONY: vet
vet: ## run static analysis
	@go vet ./...

.PHONY: docs
docs: ## generate/update swagger docs
	@swag init -g cmd/server/main.go || true
	@swag fmt



## ----- MAIN
run: ## run the code
	@go run main.go

scenarios: ## make some api requests
	@echo -e "\n-------------------- 422 --------------------"; curl -s "http://localhost:8080/cep/1234567"
	@echo -e "\n-------------------- 402 --------------------"; curl -s "http://localhost:8080/cep/12345678"
	@echo -e "\n-------------------- 200 --------------------"; curl -s "http://localhost:8080/cep/13330250"