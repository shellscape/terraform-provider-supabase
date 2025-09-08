package provider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/shellscape/terraform-provider-supabase/internal/provider/settings"
	"github.com/supabase/cli/pkg/api"
)

var _ resource.Resource = &DatabaseWebhookResource{}

func NewDatabaseWebhookResource() resource.Resource {
	return &DatabaseWebhookResource{}
}

type DatabaseWebhookResource struct {
	client *api.ClientWithResponses
}

type DatabaseWebhookResourceModel struct {
	ProjectRef types.String `tfsdk:"project_ref"`
	Id         types.String `tfsdk:"id"`
	Enabled    types.Bool   `tfsdk:"enabled"`
}

func (r *DatabaseWebhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_webhook"
}

func (r *DatabaseWebhookResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database Webhook resource for enabling database webhooks (Beta feature)",
		Attributes: map[string]schema.Attribute{
			"project_ref": schema.StringAttribute{
				MarkdownDescription: "Project reference",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "Resource identifier",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether database webhooks are enabled",
				Required:            true,
			},
		},
	}
}

func (r *DatabaseWebhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
}

func (r *DatabaseWebhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseWebhookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Enabled.ValueBool() {
		httpResp, err := r.client.V1EnableDatabaseWebhookWithResponse(ctx, data.ProjectRef.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to enable database webhook: %s", err))
			return
		}

		if httpResp.StatusCode() != http.StatusCreated {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to enable database webhook, got status %d: %s", httpResp.StatusCode(), httpResp.Body))
			return
		}
	}

	data.Id = types.StringValue(data.ProjectRef.ValueString() + "-webhook")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseWebhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseWebhookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: The API doesn't provide a way to check webhook status,
	// so we maintain the state as configured
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseWebhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseWebhookResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Enabled.ValueBool() {
		httpResp, err := r.client.V1EnableDatabaseWebhookWithResponse(ctx, data.ProjectRef.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to enable database webhook: %s", err))
			return
		}

		if httpResp.StatusCode() != http.StatusCreated && httpResp.StatusCode() != http.StatusOK {
			resp.Diagnostics.AddError("API Error", fmt.Sprintf("Unable to enable database webhook, got status %d: %s", httpResp.StatusCode(), httpResp.Body))
			return
		}
	} else {
		// Note: The API doesn't provide a disable endpoint
		resp.Diagnostics.AddWarning("Webhook Disable", "The Supabase API does not provide a way to disable database webhooks. The webhook will remain enabled in Supabase.")
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseWebhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DatabaseWebhookResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Note: The API doesn't provide a way to disable webhooks
	// The resource will be removed from Terraform state but the webhook remains enabled in Supabase
	resp.Diagnostics.AddWarning("Webhook Disable", "The Supabase API does not provide a way to disable database webhooks. The webhook will remain enabled in Supabase even though it's removed from Terraform state.")
}