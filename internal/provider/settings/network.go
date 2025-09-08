package settings

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/supabase/cli/pkg/api"
)

// NetworkConfig represents network restrictions configuration
type NetworkConfig struct {
	DbAllowedCidrs   []types.String `tfsdk:"db_allowed_cidrs"`
	DbAllowedCidrsV6 []types.String `tfsdk:"db_allowed_cidrs_v6"`
}

func GetNetworkSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"db_allowed_cidrs": schema.ListAttribute{
			MarkdownDescription: "List of allowed IPv4 CIDR blocks for database access",
			Optional:            true,
			ElementType:         types.StringType,
		},
		"db_allowed_cidrs_v6": schema.ListAttribute{
			MarkdownDescription: "List of allowed IPv6 CIDR blocks for database access",
			Optional:            true,
			ElementType:         types.StringType,
		},
	}
}

// ReadNetworkConfig reads network configuration from the API
func ReadNetworkConfig(ctx context.Context, client *api.ClientWithResponses, state *SettingsResourceModel) diag.Diagnostics {
	httpResp, err := client.V1GetNetworkRestrictionsWithResponse(ctx, state.Id.ValueString())
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read network settings: %s", err))}
	}

	switch httpResp.StatusCode() {
	case http.StatusNotFound, http.StatusNotAcceptable:
		return nil
	}

	if httpResp.JSON200 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to read network settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	if state.Network == nil {
		state.Network = &NetworkConfig{}
	}

	// Initialize slices based on API response
	if v4 := httpResp.JSON200.Config.DbAllowedCidrs; v4 != nil {
		state.Network.DbAllowedCidrs = []types.String{}
		for _, cidr := range *v4 {
			state.Network.DbAllowedCidrs = append(state.Network.DbAllowedCidrs, types.StringValue(cidr))
		}
	} else {
		state.Network.DbAllowedCidrs = nil
	}

	if v6 := httpResp.JSON200.Config.DbAllowedCidrsV6; v6 != nil {
		state.Network.DbAllowedCidrsV6 = []types.String{}
		for _, cidr := range *v6 {
			state.Network.DbAllowedCidrsV6 = append(state.Network.DbAllowedCidrsV6, types.StringValue(cidr))
		}
	} else {
		state.Network.DbAllowedCidrsV6 = nil
	}

	return nil
}

// UpdateNetworkConfig updates network configuration via the API
func UpdateNetworkConfig(ctx context.Context, client *api.ClientWithResponses, plan *SettingsResourceModel) diag.Diagnostics {
	body := api.NetworkRestrictionsRequest{
		DbAllowedCidrs:   &[]string{},
		DbAllowedCidrsV6: &[]string{},
	}

	for _, cidr := range plan.Network.DbAllowedCidrs {
		cidrStr := cidr.ValueString()
		ip, _, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Invalid CIDR: %s", cidrStr))}
		}
		if ip.IsPrivate() {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Private IP not allowed: %s", cidrStr))}
		}
		if ip.To4() != nil {
			*body.DbAllowedCidrs = append(*body.DbAllowedCidrs, cidrStr)
		}
	}

	for _, cidr := range plan.Network.DbAllowedCidrsV6 {
		cidrStr := cidr.ValueString()
		ip, _, err := net.ParseCIDR(cidrStr)
		if err != nil {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Invalid CIDR: %s", cidrStr))}
		}
		if ip.IsPrivate() {
			return diag.Diagnostics{diag.NewErrorDiagnostic("Validation Error", fmt.Sprintf("Private IP not allowed: %s", cidrStr))}
		}
		if ip.To4() == nil {
			*body.DbAllowedCidrsV6 = append(*body.DbAllowedCidrsV6, cidrStr)
		}
	}

	httpResp, err := client.V1UpdateNetworkRestrictionsWithResponse(ctx, plan.ProjectRef.ValueString(), body)
	if err != nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update network settings: %s", err))}
	}

	if httpResp.JSON201 == nil {
		return diag.Diagnostics{diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Unable to update network settings, got status %d: %s", httpResp.StatusCode(), httpResp.Body))}
	}

	return nil
}
