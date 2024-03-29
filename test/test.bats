#!/usr/bin/env bats

load "node_modules/bats-support/load"
load "node_modules/bats-assert/load"
load "node_modules/bats-file/load"

setup_file() {
  test_folder=$(pwd)
  cd ../samples/config1 || exit
  terraform init
  terraform graph > "${test_folder}"/config1_raw.gv
  cd "${test_folder}" || exit
}

setup() {
  # Remove all artifacts.
  rm -rf ./*actual*
}

get_executable_cmd() {
  if [ -z "${DOCKER_IMAGE_TAG}" ]; then
    # Run the binary.
    echo "./terraform-graph-beautifier"
  else
    # Run the Docker image.
    echo "docker run --rm -i \
      --name terraform-graph-beautifier \
      --workdir=/test \
      --volume $(pwd)/test_template.gohtml:/test/test_template.gohtml:ro \
      ${DOCKER_IMAGE_TAG}"
  fi
}

@test "is executable" {
  run $(get_executable_cmd) -v
  assert_success
  assert_output --partial "version:"
  assert_output --partial "commit:"
  assert_output --partial "built at:"
  assert_output --partial "built by:"
}

@test "graphviz" {
  run $(get_executable_cmd) --output-type=graphviz < config1_raw.gv
  assert_success
}

@test "graphviz - check output" {
  $(get_executable_cmd) --output-type=graphviz < config1_raw.gv > config1_actual.gv
  run diff -w config1_expected.gv config1_actual.gv
  assert_success
}

@test "Cytoscape.js JSON" {
  run $(get_executable_cmd) --output-type=cyto-json < config1_raw.gv
  assert_success
}

@test "Cytoscape.js JSON - check output" {
  $(get_executable_cmd) --output-type=cyto-json < config1_raw.gv > config1_actual.json
  run diff -w \
    <(jq --sort-keys '(. | (.. | arrays) |= sort)' config1_expected.json) \
    <(jq --sort-keys '(. | (.. | arrays) |= sort)' config1_actual.json)
  assert_success
}

@test "Cytoscape.js HTML - default template" {
  run $(get_executable_cmd) --output-type=cyto-html < config1_raw.gv
  assert_success
  assert_output --partial "<!-- Terraform graph beautifier default template -->"
  refute_output --partial "{{.PageTitle}}"
  refute_output --partial "{{.GraphElementsJSON}}"
}

@test "Cytoscape.js HTML - custom template" {
  run $(get_executable_cmd) --output-type=cyto-html --cyto-html-template=test_template.gohtml < config1_raw.gv
  assert_success
  assert_output --partial "<!-- Terraform graph beautifier test template -->"
  refute_output --partial "{{.PageTitle}}"
  refute_output --partial "{{.GraphElementsJSON}}"
}

@test "exclusion" {
  run $(get_executable_cmd) --output-type=graphviz < config1_raw.gv
  assert_success
  assert_output --partial "module.root.output.file_name"

  run $(get_executable_cmd) --exclude=module.root.output.file_name --output-type=graphviz < config1_raw.gv
  assert_success
  refute_output --partial "module.root.output.file_name"
}
