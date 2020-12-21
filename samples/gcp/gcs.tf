module "gcs_buckets" {
  source        = "terraform-google-modules/cloud-storage/google"
  version       = "~> 1.7.2"
  project_id    = var.project_id
  location      = var.region
  storage_class = "REGIONAL"
  prefix        = var.project_id
  names         = ["raw", "processed"]
  force_destroy =  {
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
    processed = ["tmp", "processed", "reject", "archive", "foo/bar"]
  }

  set_viewer_roles = true
  bucket_viewers = {
    raw       = module.service_accounts_v2.iam_emails["foo"]
    processed = module.service_accounts_v2.iam_emails["bar"]
  }

  set_creator_roles = true
  bucket_creators = {
    raw       = module.service_accounts_v2.iam_emails["foo"]
    processed = module.service_accounts_v2.iam_emails["bar"]
  }
}