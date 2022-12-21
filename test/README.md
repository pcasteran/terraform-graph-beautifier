# Tests

Suite of automated tests verifying that a binary or Docker image is functional. The tests are automatically launched by
a CI workflow following a build step.

The **[bats](https://www.npmjs.com/package/bats)** framework is used to define the tests, with the actions to execute
and the expecting behavior.

The image tag to be tested is specified using the `DOCKER_IMAGE_TAG` environment variable.
If the variable is not set, the tests are run on a binary named `terraform-graph-beautifier` expected to be present in
the `test` directory.

For an introduction on how to set up testing with **bats**
see [this](https://stefanzweifel.io/posts/2020/12/22/getting-started-with-bash-testing-with-bats) post.
