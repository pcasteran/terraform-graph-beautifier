locals {
  service_accounts_prefix = "sa-${var.env}"
}

// Use the old style v2.0 version that uses `count`.
module "service_accounts_v2" {
  source     = "terraform-google-modules/service-accounts/google"
  version    = "~> 2.0.2"
  project_id = var.project_id
  prefix     = local.service_accounts_prefix
  names      = ["foo", "bar"]
  project_roles = [
    "${var.project_id}=>roles/bigquery.metadataViewer",
    "${var.project_id}=>roles/logging.viewer",
  ]
}

// Use the new style v3.0 version that uses `for_each`.
module "service_accounts_v3" {
  source       = "terraform-google-modules/service-accounts/google"
  version      = "~> 3.0.1"
  project_id   = var.project_id
  prefix       = local.service_accounts_prefix
  names        = ["baz"]
  description  = "Awesome service account"
  display_name = "baz"
  project_roles = [
    "${var.project_id}=>roles/pubsub.viewer",
  ]
}
