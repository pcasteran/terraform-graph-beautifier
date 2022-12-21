module "gcs_buckets" {
  source        = "terraform-google-modules/cloud-storage/google"
  version       = "~> 3.4.0"
  project_id    = var.project_id
  location      = var.region
  storage_class = "REGIONAL"
  prefix        = "my-bucket-${var.env}"
  names         = ["raw", "processed"]
  force_destroy = {
    raw       = true
    processed = true
  }
  versioning = {
    raw       = false
    processed = true
  }
  labels = {
    env = var.env
    foo = "bar"
  }
  folders = {
    raw       = ["in", "archive"]
    processed = ["tmp", "process", "reject"]
  }

  set_viewer_roles = true
  bucket_viewers   = {
    raw       = module.service_accounts.iam_emails["foo"]
    processed = module.service_accounts.iam_emails["bar"]
  }

  set_creator_roles = true
  bucket_creators   = {
    raw       = module.service_accounts.iam_emails["foo"]
    processed = module.service_accounts.iam_emails["bar"]
  }
}
