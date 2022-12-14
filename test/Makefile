#
# Project Makefile.
#

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@echo
	@echo "\033[1;94mProject Makefile\033[0m"
	@echo
	@echo "\033[1;93mAvailable targets:\033[0m"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2}'
	@echo


.PHONY: build_image
build_image: ## Build the Docker image used for the tests
	@echo
	@echo "\033[1;93m>>\033[0m \033[1;94mCreating the development image\033[0m"
	docker buildx build --tag terraform-graph-beautifier-test .

.PHONY: test
test: ## Run the tests
	docker run --rm -it \
	  --user $(shell id -u):$(shell id -g) \
	  -v $(shell pwd)/terraform-graph-beautifier:/workspace/terraform-graph-beautifier:ro \
	  -v $(shell pwd)/test.bats:/workspace/test.bats:ro \
	  terraform-graph-beautifier-test \
	  npx bats .
