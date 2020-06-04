provider "google" {
  # Credentials to use to provision the resources.
  # No credentials specified : use the ones inferred from the execution environment.

  # Version of the provider to use.
  version = "~> v3.24.0"

  # Default values to use for the resources created with this provider.
  project = var.project_id
}

provider "google-beta" {
  # Credentials to use to provision the resources.
  # No credentials specified : use the ones inferred from the execution environment.

  # Version of the provider to use.
  version = "~> v3.24.0"

  # Default values to use for the resources created with this provider.
  project = var.project_id
}
