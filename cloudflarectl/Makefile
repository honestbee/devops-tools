help: ## List targets & descriptions
	@cat Makefile* | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

binary-linux: ## linux build
	env GOOS=linux GOARCH=arm go build -o bin/linux/cloudflarectl

binary-osx: ## osx build
	go build  -o bin/osx/cloudflarectl

.PHONY: binary-linux, binary-osx