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

.PHONY: dep
dep: ## Download the dependencies
	go mod download

.PHONY: tidy
tidy: ## Add the missing and remove the unused dependencies
	go mod tidy

.PHONY: fmt
fmt: ## Format the project sources
	go fmt ${project_sources}

.PHONY: build
build: dep ## Build the project
	CGO_ENABLED=0 go build -v -ldflags="-w -s -X 'main.version=${version}'" -o dist/ .

.PHONY: install
install: ## Compile and install the project binary
	go install .

.PHONY: clean
clean: ## Clean the temporary files
	go clean

.PHONY: doc_generate
doc_generate: install ## Generates the static documentation files
	$(eval doc_dir := "$(shell pwd)/doc")
	cd "samples/config1" && \
		terraform init && \
		terraform graph | terraform-graph-beautifier \
			--output-type=cyto-html \
			> ${doc_dir}/config1.html && \
		\
		terraform graph | terraform-graph-beautifier \
			--output-type=cyto-json \
			| jq . \
			> ${doc_dir}/config1.json && \
		\
		terraform graph | terraform-graph-beautifier \
			--output-type=graphviz \
			> ${doc_dir}/config1.gv

.PHONY: test_build_image
test_build_image: ## Build the Docker image used for the tests
	docker buildx build --tag terraform-graph-beautifier-test test/

.PHONY: test
test: build ## Run the tests
	docker run --rm -it \
	  --user $(shell id -u):$(shell id -g) \
	  -v $(shell pwd)/samples:/workspace/samples\
	  -v $(shell pwd)/dist/terraform-graph-beautifier:/workspace/test/terraform-graph-beautifier:ro \
	  -v $(shell pwd)/test/test.bats:/workspace/test/test.bats:ro \
	  -v $(shell pwd)/test/config1_expected.gv:/workspace/test/config1_expected.gv:ro \
	  -v $(shell pwd)/test/config1_expected.json:/workspace/test/config1_expected.json:ro \
	  -v $(shell pwd)/test/test_template.gohtml:/workspace/test/test_template.gohtml:ro \
	  terraform-graph-beautifier-test \
	  npx bats .

define update_terraform_lock_fn
    $(eval $@config_dir = $(1))
    cd ${$@config_dir} && \
		terraform init && \
		terraform providers lock -platform=linux_amd64 && \
		terraform providers lock -platform=linux_arm64 && \
		terraform providers lock -platform=darwin_amd64 && \
		terraform providers lock -platform=darwin_arm64 && \
		terraform providers lock -platform=windows_amd64
endef

.PHONY: update_terraform_lock
update_terraform_lock: ## Update the Terraform dependency lock file for all the supported platforms.
	@$(call update_terraform_lock_fn, "samples/config1")
	@$(call update_terraform_lock_fn, "samples/gcp")
