.PHONY: help
help: ## makeコマンド一覧を表示
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s033[0m %s\n", $$1, $$2}'

.PHONY: release
release: ## release build
	@go build -o bin -ldflags="-s -w" -trimpath

.PHONY: build
build: ## build
	@go build -o bin

.PHONY: clean
clean: ## clean
	@rm -fr bin/*