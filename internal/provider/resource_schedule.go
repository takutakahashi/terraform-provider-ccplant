package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/takutakahashi/terraform-provider-ccplant/internal/client"
)

var (
	_ resource.Resource                = (*scheduleResource)(nil)
	_ resource.ResourceWithConfigure   = (*scheduleResource)(nil)
	_ resource.ResourceWithImportState = (*scheduleResource)(nil)
)

type scheduleResource struct{ baseResource }

type scheduleModel struct {
	baseModel
	Name            types.String `tfsdk:"name"`
	Scope           types.String `tfsdk:"scope"`
	TeamID          types.String `tfsdk:"team_id"`
	Status          types.String `tfsdk:"status"`
	ScheduledAt     types.String `tfsdk:"scheduled_at"`
	CronExpr        types.String `tfsdk:"cron_expr"`
	Timezone        types.String `tfsdk:"timezone"`
	SessionConfig   types.Object `tfsdk:"session_config"`
	UserID          types.String `tfsdk:"user_id"`
	NextExecutionAt types.String `tfsdk:"next_execution_at"`
	ExecutionCount  types.Int64  `tfsdk:"execution_count"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

type scheduleSessionConfigModel struct {
	Environment      types.Map    `tfsdk:"environment"`
	Tags             types.Map    `tfsdk:"tags"`
	Params           types.Object `tfsdk:"params"`
	MemoryKey        types.Map    `tfsdk:"memory_key"`
	ReuseSession     types.Bool   `tfsdk:"reuse_session"`
	ReuseMessage     types.String `tfsdk:"reuse_message"`
	SessionProfileID types.String `tfsdk:"session_profile_id"`
}

type scheduleParamsModel struct {
	Message   types.String `tfsdk:"message"`
	AgentType types.String `tfsdk:"agent_type"`
	Oneshot   types.Bool   `tfsdk:"oneshot"`
	AuthProxy types.Bool   `tfsdk:"auth_proxy"`
}

type schedulePayload struct {
	ID              string                       `json:"id,omitempty"`
	Name            string                       `json:"name,omitempty"`
	UserID          string                       `json:"user_id,omitempty"`
	Scope           string                       `json:"scope,omitempty"`
	TeamID          string                       `json:"team_id,omitempty"`
	Status          string                       `json:"status,omitempty"`
	ScheduledAt     string                       `json:"scheduled_at,omitempty"`
	CronExpr        string                       `json:"cron_expr,omitempty"`
	Timezone        string                       `json:"timezone,omitempty"`
	SessionConfig   scheduleSessionConfigPayload `json:"session_config"`
	NextExecutionAt string                       `json:"next_execution_at,omitempty"`
	ExecutionCount  int64                        `json:"execution_count,omitempty"`
	CreatedAt       string                       `json:"created_at,omitempty"`
	UpdatedAt       string                       `json:"updated_at,omitempty"`
}

type scheduleSessionConfigPayload struct {
	Environment      map[string]string      `json:"environment,omitempty"`
	Tags             map[string]string      `json:"tags,omitempty"`
	Params           *scheduleParamsPayload `json:"params,omitempty"`
	MemoryKey        map[string]string      `json:"memory_key,omitempty"`
	ReuseSession     bool                   `json:"reuse_session,omitempty"`
	ReuseMessage     string                 `json:"reuse_message,omitempty"`
	SessionProfileID string                 `json:"session_profile_id,omitempty"`
}

type scheduleParamsPayload struct {
	Message   string `json:"message,omitempty"`
	AgentType string `json:"agent_type,omitempty"`
	Oneshot   *bool  `json:"oneshot,omitempty"`
	AuthProxy *bool  `json:"auth_proxy,omitempty"`
}

func newScheduleResource() resource.Resource {
	return &scheduleResource{baseResource: baseResource{typeName: "schedule", basePath: "/schedules"}}
}

func (r *scheduleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.metadata(ctx, req, resp)
}

func (r *scheduleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.configure(ctx, req, resp)
}

func (r *scheduleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an agentapi-proxy schedule.",
		Attributes: map[string]schema.Attribute{
			"id":            idAttribute(),
			"response_json": responseJSONAttribute(),
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Schedule name.",
			},
			"scope": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Ownership scope: `user` or `team`. This field is only sent on create; changing it requires replacement.",
				PlanModifiers:       replaceOnChangeString(),
			},
			"team_id": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Team identifier. Required by agentapi-proxy when scope is `team`. This field is only sent on create; changing it requires replacement.",
				PlanModifiers:       replaceOnChangeString(),
			},
			"status": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Schedule status: `active`, `paused`, or `completed`.",
			},
			"scheduled_at": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "One-time or first execution time as an RFC3339 timestamp.",
			},
			"cron_expr": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Recurring cron expression. Either `scheduled_at` or `cron_expr` is required by the API.",
			},
			"timezone": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "IANA timezone for schedule evaluation.",
			},
			"session_config": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Session creation configuration.",
				Attributes: map[string]schema.Attribute{
					"environment": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Environment variables for the session.",
					},
					"tags": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Tags applied to sessions created by this schedule.",
					},
					"params": schema.SingleNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Session parameters.",
						Attributes: map[string]schema.Attribute{
							"message": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Initial message sent to the agent session.",
							},
							"agent_type": schema.StringAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Agent type passed to agentapi-proxy.",
							},
							"oneshot": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Whether the session should run in one-shot mode.",
							},
							"auth_proxy": schema.BoolAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Whether auth proxy behavior is enabled for the session.",
							},
						},
					},
					"memory_key": schema.MapAttribute{
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Memory lookup tags injected into the session.",
					},
					"reuse_session": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Reuse an active matching session instead of creating a new one.",
					},
					"reuse_message": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Message sent to a reused session.",
					},
					"session_profile_id": schema.StringAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "Session profile ID used as the base session configuration.",
					},
				},
			},
			"user_id":           computedStringAttribute("User ID that owns the schedule."),
			"next_execution_at": computedStringAttribute("Next scheduled execution time."),
			"execution_count":   schema.Int64Attribute{Computed: true, MarkdownDescription: "Total execution count."},
			"created_at":        computedStringAttribute("Creation timestamp returned by agentapi-proxy."),
			"updated_at":        computedStringAttribute("Last update timestamp returned by agentapi-proxy."),
		},
	}
}

func (r *scheduleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan scheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := r.request(ctx, plan, &resp.Diagnostics, true)
	if resp.Diagnostics.HasError() {
		return
	}
	responseBody, err := r.client.CreateJSON(ctx, r.basePath, body)
	if err != nil {
		resp.Diagnostics.AddError("Create failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scheduleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state scheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	id, err := ensureID(state.ID)
	if err != nil {
		resp.State.RemoveResource(ctx)
		return
	}
	responseBody, err := r.client.GetJSON(ctx, r.itemPath(id))
	if client.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Read failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &state, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *scheduleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan scheduleModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	body := r.request(ctx, plan, &resp.Diagnostics, false)
	if resp.Diagnostics.HasError() {
		return
	}
	responseBody, err := r.client.UpdateJSON(ctx, r.itemPath(plan.ID.ValueString()), body)
	if err != nil {
		resp.Diagnostics.AddError("Update failed", err.Error())
		return
	}
	r.applyResponse(ctx, responseBody, &plan, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *scheduleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state scheduleModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Delete(ctx, r.itemPath(state.ID.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *scheduleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, req, resp)
}

func (r *scheduleResource) request(ctx context.Context, model scheduleModel, diags *diag.Diagnostics, includeCreateOnly bool) []byte {
	if model.ScheduledAt.IsNull() && model.CronExpr.IsNull() {
		diags.AddError("Missing schedule time", "Set either scheduled_at or cron_expr.")
		return nil
	}
	sessionConfig := sessionConfigPayloadFromTerraform(ctx, model.SessionConfig, diags)
	if diags.HasError() {
		return nil
	}
	payload := schedulePayload{
		Name:          model.Name.ValueString(),
		Status:        optionalString(model.Status),
		ScheduledAt:   optionalString(model.ScheduledAt),
		CronExpr:      optionalString(model.CronExpr),
		Timezone:      optionalString(model.Timezone),
		SessionConfig: sessionConfig,
	}
	if includeCreateOnly {
		payload.Scope = optionalString(model.Scope)
		payload.TeamID = optionalString(model.TeamID)
	}
	return marshalRequest(payload, diags)
}

func (r *scheduleResource) applyResponse(ctx context.Context, body []byte, model *scheduleModel, diags *diag.Diagnostics) {
	var payload schedulePayload
	unmarshalResponse(body, &payload, diags)
	if diags.HasError() {
		return
	}
	model.ID = types.StringValue(payload.ID)
	model.Name = types.StringValue(payload.Name)
	model.UserID = types.StringValue(payload.UserID)
	model.Scope = stringOrNull(payload.Scope)
	model.TeamID = stringOrNull(payload.TeamID)
	model.Status = types.StringValue(payload.Status)
	model.ScheduledAt = stringOrNull(payload.ScheduledAt)
	model.CronExpr = stringOrNull(payload.CronExpr)
	model.Timezone = stringOrNull(payload.Timezone)
	model.SessionConfig = sessionConfigToTerraform(ctx, payload.SessionConfig, model.SessionConfig, diags)
	model.NextExecutionAt = stringOrNull(payload.NextExecutionAt)
	model.ExecutionCount = types.Int64Value(payload.ExecutionCount)
	model.CreatedAt = types.StringValue(payload.CreatedAt)
	model.UpdatedAt = types.StringValue(payload.UpdatedAt)
	model.ResponseJSON = types.StringValue(normalizeJSON(body))
}

func sessionConfigPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) scheduleSessionConfigPayload {
	var cfg scheduleSessionConfigModel
	diags.Append(value.As(ctx, &cfg, objectAsOptions)...)
	if diags.HasError() {
		return scheduleSessionConfigPayload{}
	}
	return scheduleSessionConfigPayload{
		Environment:      stringMapFromTerraform(ctx, cfg.Environment, diags),
		Tags:             stringMapFromTerraform(ctx, cfg.Tags, diags),
		Params:           paramsPayloadFromTerraform(ctx, cfg.Params, diags),
		MemoryKey:        stringMapFromTerraform(ctx, cfg.MemoryKey, diags),
		ReuseSession:     !cfg.ReuseSession.IsNull() && !cfg.ReuseSession.IsUnknown() && cfg.ReuseSession.ValueBool(),
		ReuseMessage:     optionalString(cfg.ReuseMessage),
		SessionProfileID: optionalString(cfg.SessionProfileID),
	}
}

func paramsPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) *scheduleParamsPayload {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var params scheduleParamsModel
	diags.Append(value.As(ctx, &params, objectAsOptions)...)
	if diags.HasError() {
		return nil
	}
	return &scheduleParamsPayload{
		Message:   optionalString(params.Message),
		AgentType: optionalString(params.AgentType),
		Oneshot:   optionalBoolPointer(params.Oneshot),
		AuthProxy: optionalBoolPointer(params.AuthProxy),
	}
}

func sessionConfigToTerraform(ctx context.Context, payload scheduleSessionConfigPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	var currentConfig scheduleSessionConfigModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentConfig, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(sessionConfigAttrTypes())
		}
	}
	cfg := scheduleSessionConfigModel{
		Environment:      mapToTerraformPreserveEmpty(ctx, payload.Environment, currentConfig.Environment, diags),
		Tags:             mapToTerraformPreserveEmpty(ctx, payload.Tags, currentConfig.Tags, diags),
		Params:           paramsToTerraform(ctx, payload.Params, currentConfig.Params, diags),
		MemoryKey:        mapToTerraformPreserveEmpty(ctx, payload.MemoryKey, currentConfig.MemoryKey, diags),
		ReuseSession:     types.BoolValue(payload.ReuseSession),
		ReuseMessage:     stringOrNull(payload.ReuseMessage),
		SessionProfileID: stringOrNull(payload.SessionProfileID),
	}
	value, d := types.ObjectValueFrom(ctx, sessionConfigAttrTypes(), cfg)
	diags.Append(d...)
	return value
}

func paramsToTerraform(ctx context.Context, payload *scheduleParamsPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	if payload == nil {
		if isEmptyObject(current) {
			return current
		}
		return types.ObjectNull(paramsAttrTypes())
	}
	var currentParams scheduleParamsModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentParams, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(paramsAttrTypes())
		}
	}
	params := scheduleParamsModel{
		Message:   stringOrEmptyState(payload.Message, currentParams.Message),
		AgentType: stringOrEmptyState(payload.AgentType, currentParams.AgentType),
		Oneshot:   boolPointerOrFalseState(payload.Oneshot, currentParams.Oneshot),
		AuthProxy: boolPointerOrFalseState(payload.AuthProxy, currentParams.AuthProxy),
	}
	value, d := types.ObjectValueFrom(ctx, paramsAttrTypes(), params)
	diags.Append(d...)
	return value
}

func sessionConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"environment":        types.MapType{ElemType: types.StringType},
		"tags":               types.MapType{ElemType: types.StringType},
		"params":             types.ObjectType{AttrTypes: paramsAttrTypes()},
		"memory_key":         types.MapType{ElemType: types.StringType},
		"reuse_session":      types.BoolType,
		"reuse_message":      types.StringType,
		"session_profile_id": types.StringType,
	}
}

func paramsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"message":    types.StringType,
		"agent_type": types.StringType,
		"oneshot":    types.BoolType,
		"auth_proxy": types.BoolType,
	}
}

func boolPointerOrNull(value *bool) types.Bool {
	if value == nil {
		return types.BoolNull()
	}
	return types.BoolValue(*value)
}

func boolPointerOrFalseState(value *bool, current types.Bool) types.Bool {
	if value != nil {
		return types.BoolValue(*value)
	}
	if !current.IsNull() && !current.IsUnknown() && !current.ValueBool() {
		return current
	}
	return types.BoolNull()
}

func stringOrEmptyState(value string, current types.String) types.String {
	if value != "" {
		return types.StringValue(value)
	}
	if !current.IsNull() && !current.IsUnknown() {
		return current
	}
	return types.StringNull()
}

func stringOrNull(value string) types.String {
	if value == "" {
		return types.StringNull()
	}
	return types.StringValue(value)
}

func (p *scheduleSessionConfigPayload) UnmarshalJSON(data []byte) error {
	type alias scheduleSessionConfigPayload
	var raw struct {
		alias
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = scheduleSessionConfigPayload(raw.alias)
	if len(raw.Params) > 0 && string(raw.Params) != "null" {
		var params scheduleParamsPayload
		if err := json.Unmarshal(raw.Params, &params); err != nil {
			return err
		}
		p.Params = &params
	}
	return nil
}
