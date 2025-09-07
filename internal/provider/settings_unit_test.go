// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Test the new struct-based model
func TestSettingsResourceModel(t *testing.T) {
	model := SettingsResourceModel{
		ProjectRef: types.StringValue("test-project"),
		Database: &DatabaseConfig{
			StatementTimeout: types.StringValue("10s"),
			MaxConnections:   types.Int64Value(100),
		},
		Auth: &AuthConfig{
			SiteUrl:        types.StringValue("http://localhost:3000"),
			DisableSignup:  types.BoolValue(false),
			SmtpHost:      types.StringValue("smtp.example.com"),
			SmtpPort:      types.Int64Value(587),
			SmtpUser:      types.StringValue("user@example.com"),
		},
		Api: &ApiConfig{
			DbSchema:          types.StringValue("public"),
			DbExtraSearchPath: types.StringValue("public,extensions"),
			MaxRows:          types.Int64Value(1000),
		},
		Network: &NetworkConfig{
			DbAllowedCidrs: []types.String{
				types.StringValue("0.0.0.0/0"),
			},
		},
	}

	if model.ProjectRef.ValueString() != "test-project" {
		t.Errorf("Expected project_ref 'test-project', got %s", model.ProjectRef.ValueString())
	}
	
	if model.Database.StatementTimeout.ValueString() != "10s" {
		t.Errorf("Expected statement_timeout '10s', got %s", model.Database.StatementTimeout.ValueString())
	}
	
	if model.Database.MaxConnections.ValueInt64() != 100 {
		t.Errorf("Expected max_connections 100, got %d", model.Database.MaxConnections.ValueInt64())
	}
	
	if model.Auth.SiteUrl.ValueString() != "http://localhost:3000" {
		t.Errorf("Expected site_url 'http://localhost:3000', got %s", model.Auth.SiteUrl.ValueString())
	}
	
	if model.Auth.DisableSignup.ValueBool() != false {
		t.Errorf("Expected disable_signup false, got %t", model.Auth.DisableSignup.ValueBool())
	}
	
	if model.Api.MaxRows.ValueInt64() != 1000 {
		t.Errorf("Expected max_rows 1000, got %d", model.Api.MaxRows.ValueInt64())
	}
	
	if len(model.Network.DbAllowedCidrs) != 1 {
		t.Errorf("Expected 1 CIDR, got %d", len(model.Network.DbAllowedCidrs))
	}
	
	if model.Network.DbAllowedCidrs[0].ValueString() != "0.0.0.0/0" {
		t.Errorf("Expected CIDR '0.0.0.0/0', got %s", model.Network.DbAllowedCidrs[0].ValueString())
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"587", 587},
		{"25", 25},
		{"0", 0},
		{"invalid", 0}, // should default to 0 on error
		{"", 0},        // should default to 0 on error
	}

	for _, test := range tests {
		result := parseInt(test.input)
		if result != test.expected {
			t.Errorf("parseInt(%s): expected %d, got %d", test.input, test.expected, result)
		}
	}
}

func TestExternalProviders(t *testing.T) {
	model := SettingsResourceModel{
		ProjectRef: types.StringValue("test-project"),
		Auth: &AuthConfig{
			SiteUrl:        types.StringValue("http://localhost:3000"),
			DisableSignup:  types.BoolValue(false),
			ExternalGithub: &ExternalProviderConfig{
				Enabled:  types.BoolValue(true),
				ClientId: types.StringValue("github_client_123"),
				Secret:   types.StringValue("github_secret_456"),
			},
			ExternalGoogle: &ExternalProviderConfig{
				Enabled:              types.BoolValue(true),
				ClientId:             types.StringValue("google_client_789"),
				Secret:               types.StringValue("google_secret_000"),
				AdditionalClientIds:  types.StringValue("additional_123,additional_456"),
			},
			ExternalKeycloak: &ExternalProviderConfig{
				Enabled:  types.BoolValue(true),
				ClientId: types.StringValue("keycloak_client"),
				Secret:   types.StringValue("keycloak_secret"),
				Url:      types.StringValue("https://keycloak.example.com"),
			},
		},
	}

	// Test that the external providers are properly configured
	if model.Auth.ExternalGithub == nil {
		t.Error("Expected GitHub provider to be configured")
	}
	
	if !model.Auth.ExternalGithub.Enabled.ValueBool() {
		t.Error("Expected GitHub provider to be enabled")
	}
	
	if model.Auth.ExternalGithub.ClientId.ValueString() != "github_client_123" {
		t.Errorf("Expected GitHub client ID 'github_client_123', got %s", model.Auth.ExternalGithub.ClientId.ValueString())
	}
	
	if model.Auth.ExternalGoogle.AdditionalClientIds.ValueString() != "additional_123,additional_456" {
		t.Errorf("Expected Google additional client IDs, got %s", model.Auth.ExternalGoogle.AdditionalClientIds.ValueString())
	}
	
	if model.Auth.ExternalKeycloak.Url.ValueString() != "https://keycloak.example.com" {
		t.Errorf("Expected Keycloak URL, got %s", model.Auth.ExternalKeycloak.Url.ValueString())
	}
}