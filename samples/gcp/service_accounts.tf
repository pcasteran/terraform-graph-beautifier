module "service_accounts" {
  source        = "terraform-google-modules/service-accounts/google"
  # Version constraint should be "~> 4.1.1" but there is a regression on version 4.1.1 (https://github.com/terraform-google-modules/terraform-google-service-accounts/issues/59)
  version       = "4.1.0"
  project_id    = var.project_id
  prefix        = "sa"
  names         = ["foo", "bar"]
  description   = "Awesome service account"
  display_name  = "baz"
  project_roles = [
    "${var.project_id}=>roles/pubsub.viewer",
  ]
}
