output "unzipped_data_bucket" {
  value = google_storage_bucket.unzipped_data_bucket.name
}

output "ungzip_function_name" {
  value = google_cloudfunctions_function.ungzip_function.name
}

