terraform {
  required_providers {
    random = {
      source  = "hashicorp/random"
      version = "~> 3.4.3"
    }

    local = {
      source  = "hashicorp/local"
      version = "~> 2.2.3"
    }

    time = {
      source  = "hashicorp/time"
      version = "~> 0.9.1"
    }
  }
}
