resource "supabase_edge_function" "hello_world" {
  project_ref        = "abcdefghijklmnopqrst"
  slug               = "hello-world"
  name               = "Hello World Function"
  body               = file("${path.module}/hello-world.ts")
  verify_jwt         = false
  compute_multiplier = 1.0
}

# Example with JWT verification enabled
resource "supabase_edge_function" "protected_function" {
  project_ref        = "abcdefghijklmnopqrst"
  slug               = "protected-api"
  name               = "Protected API Function"
  body               = file("${path.module}/protected-api.ts")
  verify_jwt         = true
  compute_multiplier = 2.0
  entrypoint_path    = "index.ts"
  import_map         = true
  import_map_path    = "import_map.json"
}

# Output the function URLs
output "hello_world_function_id" {
  value = supabase_edge_function.hello_world.id
}

output "protected_function_status" {
  value = supabase_edge_function.protected_function.status
}