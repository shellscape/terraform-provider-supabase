# Supabase Terraform Provider

[![liber manifesto][liberation-badge]][liberation-link]

The [Supabase Provider](https://registry.terraform.io/providers/supabase/supabase/latest/docs) allows Terraform to manage resources hosted on the [Supabase](https://supabase.com/) platform.

## Features

This provider enables you to:

- Version control your project settings in Git
- Set up CI/CD pipelines for automatically provisioning projects and branches  
- Deploy and manage Edge Functions
- Query storage bucket information

## Getting Started

### Requirements

- [Terraform](https://www.terraform.io/) >= 1.0
- [Supabase](https://supabase.com/) account with access token

### Resources and Documentation

- [Step-by-step tutorial](docs/tutorial.md)
- [CI/CD example](https://github.com/supabase/supabase-action-example/tree/main/supabase/remotes)
- [Contributing guide](CONTRIBUTING.md)

## Usage

### Basic Example

This example imports an existing Supabase project and synchronises its API settings:

```hcl
terraform {
  required_providers {
    supabase = {
      source  = "supabase/supabase"
      version = "~> 1.0"
    }
  }
}

provider "supabase" {
  access_token = file("${path.module}/access-token")
}

# Define a linked project variable as user input
variable "linked_project" {
  type = string
}

# Import the linked project resource
import {
  to = supabase_project.production
  id = var.linked_project
}

resource "supabase_project" "production" {
  organization_id   = "nknnyrtlhxudbsbuazsu"
  name              = "tf-project"
  database_password = "tf-example"
  region            = "ap-southeast-1"

  lifecycle {
    ignore_changes = [database_password]
  }
}

# Configure api settings for the linked project
resource "supabase_settings" "production" {
  project_ref = var.linked_project

  api = {
    db_schema            = "public,storage,graphql_public"
    db_extra_search_path = "public,extensions"
    max_rows             = 1000
  }

  auth = {
    site_url       = "https://example.com"
    disable_signup = false
  }
}

# Deploy an edge function
resource "supabase_edge_function" "hello" {
  project_ref = var.linked_project
  slug        = "hello"
  name        = "Hello World"
  body        = file("${path.module}/functions/hello.ts")
  verify_jwt  = false
}

# Query storage buckets
data "supabase_storage_buckets" "all" {
  project_ref = var.linked_project
}
```

## Contributing

We 💛 contributions! Please read our [Contributing Guide](CONTRIBUTING.md) to get started.

## License

[MIT](LICENSE)

[liberation-badge]: https://img.shields.io/badge/libera-manifesto-lightgrey.svg
[liberation-link]: https://liberamanifesto.com
