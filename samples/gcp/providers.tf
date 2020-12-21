provider "google" {
  # Credentials to use to provision the resources.
  # No credentials specified : use the ones inferred from the execution environment.

  # Default values to use for the resources created with this provider.
  project = var.project_id
}
