# Examples

This directory contains examples that are mostly used for documentation, but can also be run/tested manually via the Terraform CLI.

The document generation tool looks for files in the following locations by default. All other *.tf files besides the ones mentioned below are ignored by the documentation tool. This is useful for creating examples that can run and/or are testable even if some parts are not relevant for the documentation.

* **provider/provider.tf** example file for the provider index page
* **data-sources/`full data source name`/data-source.tf** example file for the named data source page  
* **resources/`full resource name`/resource.tf** example file for the named resource page

## Available Examples

### Resources
- **resources/supabase_project/** - Project creation and management
- **resources/supabase_settings/** - Project configuration (API, Auth, Database, Network, Storage, Pooler)
- **resources/supabase_branch/** - Branch management for database branching
- **resources/supabase_edge_function/** - Edge Function deployment and management

### Data Sources
- **data-sources/supabase_branch/** - Query branch information
- **data-sources/supabase_pooler/** - Query connection pooler details
- **data-sources/supabase_apikeys/** - Query project API keys
- **data-sources/supabase_storage_buckets/** - Query storage bucket information

## Running Examples

To run these examples:

1. Set your Supabase access token:
   ```bash
   export SUPABASE_ACCESS_TOKEN="your-access-token"
   ```

2. Navigate to an example directory:
   ```bash
   cd examples/resources/supabase_edge_function/
   ```

3. Initialize and apply Terraform:
   ```bash
   terraform init
   terraform plan
   terraform apply
   ```

## Edge Functions Example

The `supabase_edge_function` example includes sample TypeScript files:
- **hello-world.ts** - Basic function that accepts JSON and returns a greeting
- **protected-api.ts** - JWT-protected function that validates user authentication

These demonstrate common patterns for Edge Functions including request handling, user authentication, and Supabase client usage.
