# Creates bucket to store the archive
resource "google_storage_bucket" "function_bucket" {
  name                        = "bkt-${var.project}-gcf-source"
  location                    = var.region
  project                     = var.project
  uniform_bucket_level_access = true
}

data "archive_file" "source" {
  type        = "zip"
  source_dir  = "${path.module}/src"
  output_path = "tmp/function.zip"
  excludes    = [""]
}

resource "google_storage_bucket_object" "zip" {
  source       = data.archive_file.source.output_path
  content_type = "application/zip"
  name         = "src-${data.archive_file.source.output_md5}.zip"
  bucket       = google_storage_bucket.function_bucket.name
}

resource "google_storage_bucket" "unzipped_data_bucket" {
  name                        = var.unzipped_destination_bucket
  storage_class               = "NEARLINE"
  project                     = var.project
  location                    = var.region
  uniform_bucket_level_access = true
  force_destroy               = true
}

# Create the Cloud function 
resource "google_cloudfunctions_function" "ungzip_function" {
  name                  = var.function_name
  project               = var.project
  region                = var.region
  runtime               = "go121"
  available_memory_mb   = 512
  source_archive_bucket = google_storage_bucket.function_bucket.name
  source_archive_object = google_storage_bucket_object.zip.name
  entry_point           = "UncompressFile"
  timeout               = 60
  service_account_email = var.project_sa
  event_trigger {
    event_type = "google.storage.object.finalize"
    resource   = var.source_bucket
    failure_policy {
      retry = true
    }
  }

  environment_variables = {
    "PROJECT_ID"  = var.project
    "DESTINATION" = google_storage_bucket.unzipped_data_bucket.name
  }

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}
