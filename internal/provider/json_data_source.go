package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/takutakahashi/terraform-provider-ccplant/internal/client"
)

var (
	_ datasource.DataSource              = (*jsonDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*jsonDataSource)(nil)
)

type jsonDataSource struct {
	typeName string
	path     string
	client   *client.Client
}

type jsonDataSourceModel struct {
	ResponseJSON types.String `tfsdk:"response_json"`
}

func newJSONDataSource(typeName, path string) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &jsonDataSource{typeName: typeName, path: path}
	}
}

func (d *jsonDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + d.typeName
}

func (d *jsonDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: fmt.Sprintf("Reads agentapi-proxy %s data.", strings.ReplaceAll(d.typeName, "_", " ")),
		Attributes: map[string]schema.Attribute{
			"response_json": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "JSON response returned by agentapi-proxy.",
			},
		},
	}
}

func (d *jsonDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *client.Client.")
		return
	}

	d.client = apiClient
}

func (d *jsonDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	body, err := d.client.GetJSON(ctx, d.path)
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}

	state := jsonDataSourceModel{
		ResponseJSON: types.StringValue(normalizeJSON(body)),
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
