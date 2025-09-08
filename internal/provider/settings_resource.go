package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/shellscape/terraform-provider-supabase/internal/provider/settings"
)

func NewSettingsResource() resource.Resource {
	return settings.NewSettingsResource()
}