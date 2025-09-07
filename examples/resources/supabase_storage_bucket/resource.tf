resource "supabase_storage_bucket" "user_avatars" {
  project_ref        = "abcdefghijklmnopqrst"
  name               = "user-avatars"
  public             = true
  file_size_limit    = 5242880  # 5MB
  allowed_mime_types = ["image/jpeg", "image/png", "image/webp"]
}

resource "supabase_storage_bucket" "private_documents" {
  project_ref        = "abcdefghijklmnopqrst"
  name               = "user-documents"
  public             = false
  file_size_limit    = 104857600  # 100MB
  allowed_mime_types = [
    "application/pdf",
    "application/msword",
    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
    "text/plain"
  ]
}

resource "supabase_storage_bucket" "media_uploads" {
  project_ref        = "abcdefghijklmnopqrst"
  name               = "media-uploads"
  public             = false
  file_size_limit    = 52428800  # 50MB
  allowed_mime_types = ["image/*", "video/mp4", "video/quicktime"]
}

# Output bucket information
output "user_avatars_id" {
  value = supabase_storage_bucket.user_avatars.id
}

output "private_documents_created_at" {
  value = supabase_storage_bucket.private_documents.created_at
}