terraform {
  required_version = ">= 1.3.6"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "4.46.0"
    }
  }
}
