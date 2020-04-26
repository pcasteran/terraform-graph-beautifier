resource "null_resource" "noop" {
}

module "file_1" {
  source = "./random_file"
  file_prefix = "foo"
}
