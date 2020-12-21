terraform {
  required_providers {
    random = {
      source = "hashicorp/random"
      version = "~> 3.0.0"
    }

    local = {
      source = "hashicorp/local"
      version = "~> 2.0.0"
    }

    time = {
      source = "hashicorp/time"
      version = "~> 0.6.0"
    }
  }
}
