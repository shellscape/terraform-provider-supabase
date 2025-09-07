data "supabase_storage_buckets" "all" {
  project_ref = "abcdefghijklmnopqrst"
}

# You can access individual buckets from the list
output "storage_buckets" {
  value = data.supabase_storage_buckets.all.buckets
}

# Example of accessing a specific bucket by name
output "images_bucket" {
  value = [for bucket in data.supabase_storage_buckets.all.buckets : bucket if bucket.name == "images"]
}