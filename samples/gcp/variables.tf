variable "env" {
  description = "Environment (staging, prod)"
}

variable "project_id" {
  description = "GCP project id"
}

variable "region" {
  description = "GCP region"
}

variable "zone" {
  description = "GCP zone"
}

variable "location" {
  description = "GCP location for multi-region resources"
}

variable "network" {
  description = "VPC network name"
}

variable "subnet" {
  description = "VPC subnet name"
}
