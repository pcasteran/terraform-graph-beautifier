resource "random_string" "content" {
  keepers = {
    file_prefix = var.file_prefix
  }
  length = 20
}

resource "time_static" "creation_time" {
  triggers = {
    file_prefix = var.file_prefix
  }
}

resource "local_file" "file" {
  content = random_string.content.result
  filename = "generated/${time_static.creation_time.triggers.file_prefix}-${time_static.creation_time.rfc3339}.txt"
}
