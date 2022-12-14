#!/usr/bin/env bats

load "../node_modules/bats-support/load"
load "../node_modules/bats-assert/load"
load "../node_modules/bats-file/load"

load "../common.bash"

setup() {
    init_environment
}

teardown() {
    rm *.json *.gv
}

run_executable() {
    # Get the arguments to use.
    args=$@

    echo "USE_DOCKER_IMAGE: ${USE_DOCKER_IMAGE}"
    if [ "${USE_DOCKER_IMAGE}" == "true" ]; then
      # Run the Docker image.
      # TODO
      echo "Not implemented"
      return 1
    else
      # Run the binary.
      terraform-graph-beautifier "${args}"
    fi
}

@test "is executable" {
    run run_executable -v
    assert_success
    assert_output --partial "version:"
    assert_output --partial "commit:"
    assert_output --partial "built at:"
    assert_output --partial "built by:"
}
