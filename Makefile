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

# List the packages under the current folder.
project_sources := $(shell go list ./...)

# Get the package version.
version := $(shell git describe --always --long --dirty)

.PHONY: lint
lint: ## Lint the projet sources
	docker run --rm \
		--name super-linter \
		--env RUN_LOCAL=true \
		--env-file ".github/super-linter.env" \
		--volume "$(shell pwd)":/tmp/lint \
		github/super-linter:slim-v4

.PHONY: lint_shell
lint_shell: ## Open a shell in a linter container
	docker run --rm -it \
		--name super-linter \
		--env RUN_LOCAL=true \
		--env-file ".github/super-linter.env" \
		--volume "$(shell pwd)":/tmp/lint \
		--entrypoint /bin/bash \
		--workdir /tmp/lint \
		github/super-linter:slim-v4

.PHONY: setup
setup: ## Install the project requirements.
	go install github.com/markbates/pkger/cmd/pkger

.PHONY: dep
dep: ## Download the dependencies
	go mod download

.PHONY: tidy
tidy: ## Add the missing and remove the unused dependencies
	go mod tidy

.PHONY: fmt
fmt: ## Format the project sources
	go fmt ${project_sources}

.PHONY: generate
generate: dep ## Runs the code generation
	pkger

.PHONY: build
build: dep generate ## Build the project
	go build -v -ldflags="-w -s -X 'main.version=${version}'" .

.PHONY: install
install: ## Compile and install the project binary
	go install .

.PHONY: clean
clean: ## Clean the temporary files
	go clean

.PHONY: doc_generate
doc_generate: install ## Generates the static documentation files
	$(eval doc_dir := "$(shell pwd)/doc")
	$(eval config_dir := "samples/config1")
	cd ${config_dir} && terraform init

	cd ${config_dir} && terraform graph | terraform-graph-beautifier \
		--exclude="module.root.provider" \
		--output-type=cyto-html \
		> ${doc_dir}/config1.html

	cd ${config_dir} && terraform graph | terraform-graph-beautifier \
		--exclude="module.root.provider" \
		--output-type=cyto-json \
		| jq . \
		> ${doc_dir}/config1.json

	cd ${config_dir} && terraform graph | terraform-graph-beautifier \
		--exclude="module.root.provider" \
		--output-type=graphviz \
		> ${doc_dir}/config1.gv
