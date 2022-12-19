#!/usr/bin/env bats

load "node_modules/bats-support/load"
load "node_modules/bats-assert/load"
load "node_modules/bats-file/load"

setup_file() {
  test_folder=$(pwd)
  cd ../samples/config1
  terraform init
  terraform graph > ${test_folder}/config1_raw.gv
  cd ${test_folder}
}

setup() {
  # Remove all artifacts.
  rm -rf *actual*
}

get_executable_cmd() {
  if [ "${USE_DOCKER_IMAGE}" == "true" ]; then
    # Run the Docker image.
    # TODO
    echo "Not implemented"
    return 1
  else
    # Run the binary.
    echo "./terraform-graph-beautifier"
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

@test "graphviz output" {
  $(get_executable_cmd) --output-type=graphviz < config1_raw.gv > config1_actual.gv
  run diff -w config1_expected.gv config1_actual.gv
  assert_success
}

@test "exclusion" {
  run $(get_executable_cmd) --output-type=graphviz < config1_raw.gv
  assert_success
  assert_output --partial "module.root.output.file_name"

  run $(get_executable_cmd) --exclude=module.root.output.file_name --output-type=graphviz < config1_raw.gv
  assert_success
  refute_output --partial "module.root.output.file_name"
}

@test "Cytoscape.js JSON" {
  run $(get_executable_cmd) --output-type=cyto-json < config1_raw.gv
  assert_success
}

@test "Cytoscape.js JSON output" {
  $(get_executable_cmd) --output-type=cyto-json < config1_raw.gv > config1_actual.json
  run diff -w \
    <(jq --sort-keys '(. | (.. | arrays) |= sort)' config1_expected.json) \
    <(jq --sort-keys '(. | (.. | arrays) |= sort)' config1_actual.json)
  assert_success
}

@test "Cytoscape.js HTML" {
  run $(get_executable_cmd) --output-type=cyto-html < config1_raw.gv
  assert_success
}