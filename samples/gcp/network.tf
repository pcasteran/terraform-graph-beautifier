module "vpc" {
  source  = "terraform-google-modules/network/google"
  version = "~> 2.6.0"
  project_id   = var.project_id
  network_name = var.network
  routing_mode = "GLOBAL"
  subnets = [
    {
      subnet_name           = var.subnet
      subnet_ip             = var.subnet_cidr
      subnet_region         = var.region
      subnet_private_access = "true"
      subnet_flow_logs      = "true"
      description           = "The one and only subnet"
    }
  ]
}
