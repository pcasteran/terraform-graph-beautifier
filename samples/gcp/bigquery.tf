locals {
  raw_data_labels = {
    env       = var.env
    source    = "xyz"
    data-type = "raw"
  }
}

module "bigquery" {
  source  = "terraform-google-modules/bigquery/google"
  version = "~> 4.2"

  dataset_id     = "raw"
  dataset_name   = "Raw data"
  description    = "This is the raw data received from xyz."
  project_id     = var.project_id
  location       = var.location
  dataset_labels = local.raw_data_labels

  // In real-life, it's ok to use Terraform to provision the datasets and the related IAM but
  // don't use it for the tables, views and other BigQuery resources.
  // Trust me, it is not a good idea to to manage your table schemas / view queries / ...
  // as part of your __infrastructure__ configuration ;)
  tables = [
    {
      table_id = "product",
      // No "description" field available in the module :(
      schema            = "schemas/product.json",
      clustering        = [],
      time_partitioning = null,
      expiration_time   = null,
      labels            = local.raw_data_labels
    },
    {
      table_id = "basket",
      schema   = "schemas/basket.json",
      time_partitioning = {
        type                     = "DAY",
        field                    = "creation_time",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      expiration_time = null,
      clustering      = ["store_id"],
      labels          = local.raw_data_labels
    },
    {
      table_id = "basket_line",
      schema   = "schemas/basket_line.json",
      time_partitioning = {
        type                     = "DAY",
        field                    = "creation_time",
        require_partition_filter = true,
        expiration_ms            = null,
      },
      expiration_time = null,
      clustering      = ["basket_id", "product_id"],
      labels          = local.raw_data_labels
    }
  ]
  views = [
    {
      view_id        = "today_baskets",
      use_legacy_sql = false,
      query          = file("views/today_baskets.sql"),
      labels         = local.raw_data_labels
    }
  ]
}
