output "sa_foo_email" {
  description = "Foo service account email"
  value       = module.service_accounts.emails["foo"]
}

output "sa_bar_email" {
  description = "Bar service account email"
  value       = module.service_accounts.emails["bar"]
}

output "bucket_raw_data" {
  description = "Raw data bucket"
  value       = module.gcs_buckets.names["raw"]
}

output "bucket_processed_data" {
  description = "Processed data bucket"
  value       = module.gcs_buckets.names["processed"]
}
