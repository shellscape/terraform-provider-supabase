resource "supabase_settings" "production" {
  project_ref = "mayuaycdtijbctgqbycg"

  database = {
    statement_timeout = "10s"
  }

  network = {
    db_allowed_cidrs    = ["0.0.0.0/0"]
    db_allowed_cidrs_v6 = ["::/0"]
  }

  api = {
    db_schema            = "public,storage,graphql_public"
    db_extra_search_path = "public,extensions"
    max_rows             = 1000
  }

  auth = {
    site_url             = "http://localhost:3000"
    mailer_otp_exp       = 3600
    mfa_phone_otp_length = 6
    sms_otp_length       = 6

    # GitHub OAuth configuration example (structured)
    external_github = {
      enabled   = true
      client_id = "your_github_client_id"
      secret    = "your_github_client_secret"
    }
    
    # OR use direct properties (for CDKTF compatibility)
    # external_github_client_id = "your_github_client_id"
    # external_github_enabled = true
    # external_google_client_id = "your_google_client_id"
    # external_google_enabled = true
    # external_apple_client_id = "your_apple_client_id"
    # external_apple_enabled = true
    # external_facebook_client_id = "your_facebook_client_id"
    # external_facebook_enabled = true
    # external_azure_client_id = "your_azure_client_id"
    # external_azure_enabled = true
    # external_discord_client_id = "your_discord_client_id"
    # external_discord_enabled = true
  }

  storage = {
    file_size_limit = 52428800  # 50MB in bytes
    features = {
      image_transformation = {
        enabled = true  # Enable image transformation features
      }
      s3_protocol = {
        enabled = true  # Enable S3 protocol compatibility
      }
    }
  }

  pooler = {
    default_pool_size = 20
  }
}
