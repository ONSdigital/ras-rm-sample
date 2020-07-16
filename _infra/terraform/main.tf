resource "google_pubsub_topic" "sample-jobs" {
  project = var.project
  name = "sample-jobs"

  message_storage_policy {
    allowed_persistence_regions = [
      "europe-west2",
    ]
  }
}

resource "google_pubsub_subscription" "sample-workers" {
  project = var.project
  name  = "sample-workers"
  topic = google_pubsub_topic.sample-jobs.name

  labels = {
    foo = "sample-service-workers"
  }
}

resource "google_service_account" "sample-workers" {
  project = var.project
  account_id   = var.service_account
  display_name = var.service_account_name
}

resource "google_service_account_key" "sample-workers-key" {
  service_account_id = google_service_account.sample-workers.name
}

output "google_service_account_key_json"  {
  value     = google_service_account_key.sample-workers-key.private_key
}

resource "google_project_iam_member" "sample-workers-iam-binding" {
  project           = var.project
  role               = "roles/pubsub.editor"
  member             = "serviceAccount:${google_service_account.sample-workers.email}"
}

resource "google_project_iam_binding" "sample-workers-iam-member" {
  project           = var.project
  role               = "roles/pubsub.editor"

  members = [
    "serviceAccount:${google_service_account.sample-workers.email}",
  ]
}


//resource "kubernetes_secret" "google-application-credentials" {
//  metadata {
//    name = "google-application-credentials"
//  }
//  data = {
//    "credentials.json" = base64decode(google_service_account_key.sample-workers-key.private_key)
//  }
//}