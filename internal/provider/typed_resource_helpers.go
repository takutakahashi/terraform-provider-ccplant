package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/takutakahashi/terraform-provider-ccplant/internal/client"
)

var objectAsOptions = basetypes.ObjectAsOptions{
	UnhandledNullAsEmpty:    true,
	UnhandledUnknownAsEmpty: true,
}

type baseResource struct {
	typeName string
	basePath string
	client   *client.Client
}

type baseModel struct {
	ID           types.String `tfsdk:"id"`
	ResponseJSON types.String `tfsdk:"response_json"`
}

func (r *baseResource) metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_" + r.typeName
}

func (r *baseResource) configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	apiClient, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected provider data", "Expected *client.Client.")
		return
	}
	r.client = apiClient
}

func (r *baseResource) itemPath(id string) string {
	return strings.TrimRight(r.basePath, "/") + "/" + id
}

func idAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Resource ID assigned by agentapi-proxy.",
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
}

func responseJSONAttribute() schema.StringAttribute {
	return schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: "Latest normalized JSON response returned by agentapi-proxy.",
	}
}

func computedStringAttribute(description string) schema.StringAttribute {
	return schema.StringAttribute{
		Computed:            true,
		MarkdownDescription: description,
	}
}

func replaceOnChangeString() []planmodifier.String {
	return []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	}
}

func optionalStringPointer(v types.String) *string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	s := v.ValueString()
	return &s
}

func optionalString(v types.String) string {
	if v.IsNull() || v.IsUnknown() {
		return ""
	}
	return v.ValueString()
}

func optionalBoolPointer(v types.Bool) *bool {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	b := v.ValueBool()
	return &b
}

func optionalIntPointer(v types.Int64) *int {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	i := int(v.ValueInt64())
	return &i
}

func optionalInt(v types.Int64) int {
	if v.IsNull() || v.IsUnknown() {
		return 0
	}
	return int(v.ValueInt64())
}

func stringMapFromTerraform(ctx context.Context, v types.Map, diags *diag.Diagnostics) map[string]string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	values := make(map[string]string, len(v.Elements()))
	diags.Append(v.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil
	}
	return values
}

func optionalStringMapPointer(ctx context.Context, v types.Map, diags *diag.Diagnostics) *map[string]string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	values := stringMapFromTerraform(ctx, v, diags)
	return &values
}

func stringListFromTerraform(ctx context.Context, v types.List, diags *diag.Diagnostics) []string {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	values := make([]string, 0, len(v.Elements()))
	diags.Append(v.ElementsAs(ctx, &values, false)...)
	if diags.HasError() {
		return nil
	}
	return values
}

func mapToTerraform(ctx context.Context, values map[string]string, diags *diag.Diagnostics) types.Map {
	if values == nil {
		return types.MapNull(types.StringType)
	}
	v, d := types.MapValueFrom(ctx, types.StringType, values)
	diags.Append(d...)
	return v
}

func mapToTerraformPreserveEmpty(ctx context.Context, values map[string]string, current types.Map, diags *diag.Diagnostics) types.Map {
	if values == nil && isEmptyMap(current) {
		values = map[string]string{}
	}
	return mapToTerraform(ctx, values, diags)
}

func listToTerraform(ctx context.Context, values []string, diags *diag.Diagnostics) types.List {
	if values == nil {
		return types.ListNull(types.StringType)
	}
	v, d := types.ListValueFrom(ctx, types.StringType, values)
	diags.Append(d...)
	return v
}

func listToTerraformPreserveEmpty(ctx context.Context, values []string, current types.List, diags *diag.Diagnostics) types.List {
	if values == nil && isEmptyList(current) {
		values = []string{}
	}
	return listToTerraform(ctx, values, diags)
}

func isEmptyMap(v types.Map) bool {
	return !v.IsNull() && !v.IsUnknown() && len(v.Elements()) == 0
}

func isEmptyList(v types.List) bool {
	return !v.IsNull() && !v.IsUnknown() && len(v.Elements()) == 0
}

func isEmptyObject(v types.Object) bool {
	if v.IsNull() || v.IsUnknown() {
		return false
	}
	for _, value := range v.Attributes() {
		if !value.IsNull() && !value.IsUnknown() {
			return false
		}
	}
	return true
}

func jsonStringForRequest(v types.String, attrPath path.Path, diags *diag.Diagnostics) (json.RawMessage, bool) {
	if v.IsNull() || v.IsUnknown() || v.ValueString() == "" {
		return nil, true
	}
	body := []byte(v.ValueString())
	if !json.Valid(body) {
		diags.AddAttributeError(attrPath, "Invalid JSON", "Value must contain valid JSON.")
		return nil, false
	}
	return json.RawMessage(body), true
}

func normalizedJSONFromRaw(raw json.RawMessage) types.String {
	if len(raw) == 0 || string(raw) == "null" {
		return types.StringNull()
	}
	return types.StringValue(normalizeJSON(raw))
}

func marshalRequest(v any, diags *diag.Diagnostics) []byte {
	body, err := json.Marshal(v)
	if err != nil {
		diags.AddError("Build request failed", err.Error())
		return nil
	}
	return body
}

func unmarshalResponse(body []byte, out any, diags *diag.Diagnostics) {
	if err := json.Unmarshal(body, out); err != nil {
		diags.AddError("Decode response failed", err.Error())
	}
}

func importByID(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func ensureID(id types.String) (string, error) {
	if id.IsNull() || id.IsUnknown() || id.ValueString() == "" {
		return "", fmt.Errorf("missing resource ID")
	}
	return id.ValueString(), nil
}

func normalizeJSON(body []byte) string {
	var value any
	if err := json.Unmarshal(body, &value); err != nil {
		return string(body)
	}

	normalized, err := json.Marshal(value)
	if err != nil {
		return string(body)
	}

	return string(normalized)
}
