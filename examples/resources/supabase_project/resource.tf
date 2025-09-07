resource "supabase_project" "test" {
  organization_id   = "continued-brown-smelt"
  name              = "foo"
  db_pass           = "bar"
  region            = "us-east-1"
  instance_size     = "micro"

  lifecycle {
    ignore_changes = [
      db_pass,
      instance_size,
    ]
  }
}
