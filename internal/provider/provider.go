package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/takutakahashi/terraform-provider-ccplant/internal/client"
)

var _ provider.Provider = (*ccplantProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ccplantProvider{version: version}
	}
}

type ccplantProvider struct {
	version string
}

type providerModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *ccplantProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ccplant"
	resp.Version = p.version
}

func (p *ccplantProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Terraform provider for managing agentapi-proxy resources.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "agentapi-proxy base URL. May also be set with CCPLANT_ENDPOINT.",
			},
			"api_key": schema.StringAttribute{
				Optional:            true,
				Sensitive:           true,
				MarkdownDescription: "agentapi-proxy API key. May also be set with CCPLANT_API_KEY.",
			},
		},
	}
}

func (p *ccplantProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("CCPLANT_ENDPOINT")
	if !config.Endpoint.IsNull() && !config.Endpoint.IsUnknown() {
		endpoint = config.Endpoint.ValueString()
	}

	apiKey := os.Getenv("CCPLANT_API_KEY")
	if !config.APIKey.IsNull() && !config.APIKey.IsUnknown() {
		apiKey = config.APIKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing agentapi-proxy endpoint",
			"Set the endpoint provider attribute or CCPLANT_ENDPOINT.",
		)
		return
	}

	apiClient, err := client.New(endpoint, apiKey)
	if err != nil {
		resp.Diagnostics.AddError("Invalid provider configuration", err.Error())
		return
	}

	resp.DataSourceData = apiClient
	resp.ResourceData = apiClient
}

func (p *ccplantProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		newSettingsResource,
		newJSONResource("webhook", "/webhooks"),
		newJSONResource("schedule", "/schedules"),
		newJSONResource("slackbot", "/slackbots"),
		newJSONResource("memory", "/memories"),
		newJSONResource("session_profile", "/session-profiles"),
		newJSONResource("sandbox_policy", "/sandbox-policies"),
	}
}

func (p *ccplantProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		newJSONDataSource("user_info", "/user/info"),
		newJSONDataSource("sessions", "/search"),
		newJSONDataSource("settings_managers", "/settings/managers"),
	}
}
