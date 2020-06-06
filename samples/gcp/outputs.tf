output "sa_foo_email" {
  description = "Foo service account email"
  value       = module.service_accounts_v2.emails["foo"]
}

output "sa_bar_email" {
  description = "Bar service account email"
  value       = module.service_accounts_v2.emails["bar"]
}

output "sa_baz_email" {
  description = "Baz service account email"
  value       = module.service_accounts_v3.emails["baz"]
}

output "bucket_raw_data" {
  description = "Raw data bucket"
  value       = module.gcs_buckets.names["raw"]
}

output "bucket_processed_data" {
  description = "Processed data bucket"
  value       = module.gcs_buckets.names["processed"]
}
