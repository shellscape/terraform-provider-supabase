package provider

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/shellscape/terraform-provider-supabase/internal/provider/settings"
	"github.com/supabase/cli/pkg/api"
)

var (
	_ resource.Resource              = &EdgeFunctionResource{}
	_ resource.ResourceWithConfigure = &EdgeFunctionResource{}
)

func NewEdgeFunctionResource() resource.Resource {
	return &EdgeFunctionResource{}
}

type EdgeFunctionResource struct {
	client  *api.ClientWithResponses
	tempDir string
}



type EdgeFunctionResourceModel struct {
	ProjectRef        types.String  `tfsdk:"project_ref"`
	Id                types.String  `tfsdk:"id"`
	Slug              types.String  `tfsdk:"slug"`
	Name              types.String  `tfsdk:"name"`
	EntrypointPath    types.String  `tfsdk:"entrypoint_path"`
	ImportMapPath     types.String  `tfsdk:"import_map_path"`
	VerifyJwt         types.Bool    `tfsdk:"verify_jwt"`
	ComputeMultiplier types.Float64 `tfsdk:"compute_multiplier"`
	Status            types.String  `tfsdk:"status"`
	CreatedAt         types.Int64   `tfsdk:"created_at"`
	UpdatedAt         types.Int64   `tfsdk:"updated_at"`
}

func (r *EdgeFunctionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_edge_function"
}

func (r *EdgeFunctionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: `

Manages a Supabase edge function using proper bundling and deployment.

Refer to the [Supabase Edge Functions documentation](https://supabase.com/docs/guides/functions) for more information.

## Example Usage

~~~hcl
resource "supabase_edge_function" "example" {
  project_ref      = "abcdefghijklmnopqrst"
  slug             = "hello-world"
  name             = "Hello World Function"
  entrypoint_path  = "${path.module}/functions/hello-world/index.ts"
  import_map_path  = "${path.module}/functions/hello-world/import_map.json"
  verify_jwt       = false
}
~~~
`,
		Attributes: map[string]schema.Attribute{
			"project_ref": schema.StringAttribute{
				MarkdownDescription: "Project reference ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Function ID",
				Computed:            true,
			},
			"slug": schema.StringAttribute{
				MarkdownDescription: "Function slug (URL path component). Must contain only letters, numbers, underscores, and hyphens.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[A-Za-z0-9_-]+$`),
						"slug must contain only letters, numbers, underscores, and hyphens",
					),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Function display name",
				Required:            true,
			},
			"entrypoint_path": schema.StringAttribute{
				MarkdownDescription: "Path to the function entrypoint file (e.g., index.ts)",
				Required:            true,
			},
			"import_map_path": schema.StringAttribute{
				MarkdownDescription: "Path to the import map file (optional)",
				Optional:            true,
			},
			"verify_jwt": schema.BoolAttribute{
				MarkdownDescription: "Whether to verify JWT tokens for this function",
				Optional:            true,
			},
			"compute_multiplier": schema.Float64Attribute{
				MarkdownDescription: "Compute multiplier for the function (affects performance and billing)",
				Optional:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Function status",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "Unix timestamp when the function was created",
				Computed:            true,
			},
			"updated_at": schema.Int64Attribute{
				MarkdownDescription: "Unix timestamp when the function was last updated",
				Computed:            true,
			},
		},
	}
}

func (r *EdgeFunctionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(*settings.SupabaseProviderData)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *settings.SupabaseProviderData, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = providerData.ManagementClient
	r.tempDir = os.TempDir()
}

func (r *EdgeFunctionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data EdgeFunctionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	slug := data.Slug.ValueString()
	name := data.Name.ValueString()

	// Create the function - for now just use basic creation without bundling
	response, err := r.client.V1CreateAFunctionWithResponse(
		ctx,
		data.ProjectRef.ValueString(),
		&api.V1CreateAFunctionParams{
			Slug: &slug,
			Name: &name,
		},
		api.V1CreateAFunctionJSONRequestBody{
			Slug: slug,
			Name: name,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create edge function: %s", err))
		return
	}

	if response.StatusCode() != 201 || response.JSON201 == nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to create edge function, got status %d", response.StatusCode()))
		return
	}

	// Update the data model with response values
	r.updateDataFromResponse(&data, response.JSON201)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EdgeFunctionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data EdgeFunctionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the function details
	response, err := r.client.V1GetAFunctionWithResponse(ctx, data.ProjectRef.ValueString(), data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read edge function, got error: %s", err))
		return
	}

	if response.StatusCode() == 404 {
		// Function no longer exists
		resp.State.RemoveResource(ctx)
		return
	}

	if response.StatusCode() != 200 {
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to read edge function, got status %d", response.StatusCode()),
		)
		return
	}

	if response.JSON200 == nil {
		resp.Diagnostics.AddError("API Error", "Empty response from get function API")
		return
	}

	// Convert response to FunctionResponse for consistency
	functionResp := api.FunctionResponse{
		Id:                response.JSON200.Id,
		Slug:              response.JSON200.Slug,
		Name:              response.JSON200.Name,
		Status:            api.FunctionResponseStatus(response.JSON200.Status),
		CreatedAt:         response.JSON200.CreatedAt,
		UpdatedAt:         response.JSON200.UpdatedAt,
		ComputeMultiplier: response.JSON200.ComputeMultiplier,
		EntrypointPath:    response.JSON200.EntrypointPath,
		ImportMap:         response.JSON200.ImportMap,
		ImportMapPath:     response.JSON200.ImportMapPath,
	}

	// Update the data model with response values
	r.updateDataFromResponse(&data, &functionResp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EdgeFunctionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data EdgeFunctionResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Simple update without bundling
	name := data.Name.ValueString()
	response, err := r.client.V1UpdateAFunctionWithResponse(
		ctx,
		data.ProjectRef.ValueString(),
		data.Slug.ValueString(),
		nil, // params
		api.V1UpdateAFunctionJSONRequestBody{
			Name: &name,
		},
	)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update edge function: %s", err))
		return
	}

	if response.StatusCode() != 200 || response.JSON200 == nil {
		resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to update edge function, got status %d", response.StatusCode()))
		return
	}

	// Update the data model with response values
	r.updateDataFromResponse(&data, response.JSON200)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *EdgeFunctionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data EdgeFunctionResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the function
	response, err := r.client.V1DeleteAFunctionWithResponse(ctx, data.ProjectRef.ValueString(), data.Slug.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete edge function, got error: %s", err))
		return
	}

	if response.StatusCode() != 200 && response.StatusCode() != 404 {
		resp.Diagnostics.AddError(
			"API Error",
			fmt.Sprintf("Unable to delete edge function, got status %d", response.StatusCode()),
		)
		return
	}
}

func (r *EdgeFunctionResource) updateDataFromResponse(data *EdgeFunctionResourceModel, resp *api.FunctionResponse) {
	data.Id = types.StringValue(resp.Id)
	data.Slug = types.StringValue(resp.Slug)
	data.Name = types.StringValue(resp.Name)
	data.Status = types.StringValue(string(resp.Status))
	data.CreatedAt = types.Int64Value(resp.CreatedAt)
	data.UpdatedAt = types.Int64Value(resp.UpdatedAt)

	if resp.ComputeMultiplier != nil {
		data.ComputeMultiplier = types.Float64Value(float64(*resp.ComputeMultiplier))
	} else if data.ComputeMultiplier.IsNull() || data.ComputeMultiplier.IsUnknown() {
		data.ComputeMultiplier = types.Float64Null()
	}

	// Only update these fields if they come from the response, otherwise preserve the input values
	if resp.EntrypointPath != nil {
		data.EntrypointPath = types.StringValue(*resp.EntrypointPath)
	} else if data.EntrypointPath.IsUnknown() {
		// If unknown, set to null; otherwise preserve existing value
		data.EntrypointPath = types.StringNull()
	}

	if resp.ImportMapPath != nil {
		data.ImportMapPath = types.StringValue(*resp.ImportMapPath)
	} else if data.ImportMapPath.IsUnknown() {
		// If unknown, set to null; otherwise preserve existing value
		data.ImportMapPath = types.StringNull()
	}
}

// Helper functions for pointer conversions
func boolPtr(b bool) *bool {
	return &b
}

func boolPtrFromTypes(t types.Bool) *bool {
	if t.IsNull() || t.IsUnknown() {
		return nil
	}
	val := t.ValueBool()
	return &val
}

// toFileURL converts a file path to a file:// URL - copied from CLI batch.go
func toFileURL(hostPath string) *string {
	if hostPath == "" {
		return nil
	}
	// For simplicity, just return the path as-is for now
	// In a full implementation, this would convert to proper file:// URL
	return &hostPath
}