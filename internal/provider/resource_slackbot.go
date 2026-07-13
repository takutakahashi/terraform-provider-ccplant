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
	_ resource.Resource                = (*slackbotResource)(nil)
	_ resource.ResourceWithConfigure   = (*slackbotResource)(nil)
	_ resource.ResourceWithImportState = (*slackbotResource)(nil)
)

type slackbotResource struct{ baseResource }

type slackbotModel struct {
	baseModel
	Name                   types.String `tfsdk:"name"`
	Scope                  types.String `tfsdk:"scope"`
	TeamID                 types.String `tfsdk:"team_id"`
	Teams                  types.List   `tfsdk:"teams"`
	Status                 types.String `tfsdk:"status"`
	BotTokenSecretName     types.String `tfsdk:"bot_token_secret_name"`
	BotTokenSecretKey      types.String `tfsdk:"bot_token_secret_key"`
	AllowedEventTypes      types.List   `tfsdk:"allowed_event_types"`
	AllowedChannelNames    types.List   `tfsdk:"allowed_channel_names"`
	AllowedUserIDs         types.List   `tfsdk:"allowed_user_ids"`
	SessionConfig          types.Object `tfsdk:"session_config"`
	MaxSessions            types.Int64  `tfsdk:"max_sessions"`
	NotifyOnSessionCreated types.Bool   `tfsdk:"notify_on_session_created"`
	AllowBotMessages       types.Bool   `tfsdk:"allow_bot_messages"`
	BotToken               types.String `tfsdk:"bot_token"`
	AppToken               types.String `tfsdk:"app_token"`
	UserID                 types.String `tfsdk:"user_id"`
	CreatedAt              types.String `tfsdk:"created_at"`
	UpdatedAt              types.String `tfsdk:"updated_at"`
}

type slackbotSessionConfigModel struct {
	InitialMessageTemplate types.String `tfsdk:"initial_message_template"`
	ReuseMessageTemplate   types.String `tfsdk:"reuse_message_template"`
	Tags                   types.Map    `tfsdk:"tags"`
	Environment            types.Map    `tfsdk:"environment"`
	Params                 types.Object `tfsdk:"params"`
	MemoryKey              types.Map    `tfsdk:"memory_key"`
}

type slackbotParamsModel struct {
	AgentType    types.String `tfsdk:"agent_type"`
	Oneshot      types.Bool   `tfsdk:"oneshot"`
	AuthProxy    types.Bool   `tfsdk:"auth_proxy"`
	RepoFullName types.String `tfsdk:"repo_full_name"`
}

type slackbotPayload struct {
	ID                     string                 `json:"id,omitempty"`
	Name                   string                 `json:"name,omitempty"`
	UserID                 string                 `json:"user_id,omitempty"`
	Scope                  string                 `json:"scope,omitempty"`
	TeamID                 string                 `json:"team_id,omitempty"`
	Teams                  []string               `json:"teams,omitempty"`
	Status                 string                 `json:"status,omitempty"`
	BotTokenSecretName     string                 `json:"bot_token_secret_name,omitempty"`
	BotTokenSecretKey      string                 `json:"bot_token_secret_key,omitempty"`
	AllowedEventTypes      []string               `json:"allowed_event_types,omitempty"`
	AllowedChannelNames    []string               `json:"allowed_channel_names,omitempty"`
	AllowedUserIDs         []string               `json:"allowed_user_ids,omitempty"`
	SessionConfig          *slackbotSessionConfig `json:"session_config,omitempty"`
	MaxSessions            int                    `json:"max_sessions,omitempty"`
	NotifyOnSessionCreated *bool                  `json:"notify_on_session_created,omitempty"`
	AllowBotMessages       *bool                  `json:"allow_bot_messages,omitempty"`
	BotToken               string                 `json:"bot_token,omitempty"`
	AppToken               string                 `json:"app_token,omitempty"`
	CreatedAt              string                 `json:"created_at,omitempty"`
	UpdatedAt              string                 `json:"updated_at,omitempty"`
}

type slackbotSessionConfig struct {
	InitialMessageTemplate string                 `json:"initial_message_template,omitempty"`
	ReuseMessageTemplate   string                 `json:"reuse_message_template,omitempty"`
	Tags                   map[string]string      `json:"tags,omitempty"`
	Environment            map[string]string      `json:"environment,omitempty"`
	Params                 *slackbotParamsPayload `json:"params,omitempty"`
	MemoryKey              map[string]string      `json:"memory_key,omitempty"`
}

type slackbotParamsPayload struct {
	AgentType    string `json:"agent_type,omitempty"`
	Oneshot      *bool  `json:"oneshot,omitempty"`
	AuthProxy    *bool  `json:"auth_proxy,omitempty"`
	RepoFullName string `json:"repo_full_name,omitempty"`
}

func newSlackbotResource() resource.Resource {
	return &slackbotResource{baseResource: baseResource{typeName: "slackbot", basePath: "/slackbots"}}
}

func (r *slackbotResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.metadata(ctx, req, resp)
}

func (r *slackbotResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.configure(ctx, req, resp)
}

func (r *slackbotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an agentapi-proxy SlackBot.",
		Attributes: map[string]schema.Attribute{
			"id":            idAttribute(),
			"response_json": responseJSONAttribute(),
			"name":          schema.StringAttribute{Required: true, MarkdownDescription: "SlackBot name."},
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
			"teams": schema.ListAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Team IDs whose settings are merged into sessions created by this bot.",
			},
			"status":                    schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "SlackBot status."},
			"bot_token_secret_name":     schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Kubernetes Secret name containing Slack tokens."},
			"bot_token_secret_key":      schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Secret key containing the Slack bot token."},
			"allowed_event_types":       stringListAttribute("Allowed Slack event types."),
			"allowed_channel_names":     stringListAttribute("Allowed Slack channel names."),
			"allowed_user_ids":          stringListAttribute("Allowed Slack user IDs."),
			"session_config":            slackbotSessionConfigAttribute(),
			"max_sessions":              schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: useStateInt64(), MarkdownDescription: "Maximum concurrent sessions created by this bot."},
			"notify_on_session_created": schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: useStateBool(), MarkdownDescription: "Whether the bot posts a Slack message when creating a session."},
			"allow_bot_messages":        schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: useStateBool(), MarkdownDescription: "Whether the bot processes messages posted by other bots."},
			"bot_token":                 schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Write-only Slack bot token."},
			"app_token":                 schema.StringAttribute{Optional: true, Sensitive: true, MarkdownDescription: "Write-only Slack app token."},
			"user_id":                   computedStringAttribute("User ID that owns the SlackBot."),
			"created_at":                computedStringAttribute("Creation timestamp returned by agentapi-proxy."),
			"updated_at":                computedStringAttribute("Last update timestamp returned by agentapi-proxy."),
		},
	}
}

func (r *slackbotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan slackbotModel
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

func (r *slackbotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state slackbotModel
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

func (r *slackbotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan slackbotModel
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

func (r *slackbotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state slackbotModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Delete(ctx, r.itemPath(state.ID.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *slackbotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, req, resp)
}

func (r *slackbotResource) request(ctx context.Context, model slackbotModel, diags *diag.Diagnostics, includeCreateOnly bool) []byte {
	payload := slackbotPayload{
		Name:                   model.Name.ValueString(),
		Status:                 optionalString(model.Status),
		Teams:                  stringListFromTerraform(ctx, model.Teams, diags),
		BotTokenSecretName:     optionalString(model.BotTokenSecretName),
		BotTokenSecretKey:      optionalString(model.BotTokenSecretKey),
		AllowedEventTypes:      stringListFromTerraform(ctx, model.AllowedEventTypes, diags),
		AllowedChannelNames:    stringListFromTerraform(ctx, model.AllowedChannelNames, diags),
		AllowedUserIDs:         stringListFromTerraform(ctx, model.AllowedUserIDs, diags),
		SessionConfig:          slackbotSessionConfigPayloadFromTerraform(ctx, model.SessionConfig, diags),
		MaxSessions:            optionalInt(model.MaxSessions),
		NotifyOnSessionCreated: optionalBoolPointer(model.NotifyOnSessionCreated),
		AllowBotMessages:       optionalBoolPointer(model.AllowBotMessages),
		BotToken:               optionalString(model.BotToken),
		AppToken:               optionalString(model.AppToken),
	}
	if includeCreateOnly {
		payload.Scope = optionalString(model.Scope)
		payload.TeamID = optionalString(model.TeamID)
	}
	return marshalRequest(payload, diags)
}

func (r *slackbotResource) applyResponse(ctx context.Context, body []byte, model *slackbotModel, diags *diag.Diagnostics) {
	var payload slackbotPayload
	unmarshalResponse(body, &payload, diags)
	if diags.HasError() {
		return
	}
	model.ID = types.StringValue(payload.ID)
	model.Name = types.StringValue(payload.Name)
	model.UserID = types.StringValue(payload.UserID)
	model.Scope = stringOrNull(payload.Scope)
	model.TeamID = stringOrNull(payload.TeamID)
	model.Teams = listToTerraformPreserveEmpty(ctx, payload.Teams, model.Teams, diags)
	model.Status = stringOrNull(payload.Status)
	model.BotTokenSecretName = stringOrEmptyState(payload.BotTokenSecretName, model.BotTokenSecretName)
	model.BotTokenSecretKey = stringOrEmptyState(payload.BotTokenSecretKey, model.BotTokenSecretKey)
	model.AllowedEventTypes = listToTerraformPreserveEmpty(ctx, payload.AllowedEventTypes, model.AllowedEventTypes, diags)
	model.AllowedChannelNames = listToTerraformPreserveEmpty(ctx, payload.AllowedChannelNames, model.AllowedChannelNames, diags)
	model.AllowedUserIDs = listToTerraformPreserveEmpty(ctx, payload.AllowedUserIDs, model.AllowedUserIDs, diags)
	model.SessionConfig = slackbotSessionConfigToTerraform(ctx, payload.SessionConfig, model.SessionConfig, diags)
	model.MaxSessions = types.Int64Value(int64(payload.MaxSessions))
	model.NotifyOnSessionCreated = boolPointerOrFalseState(payload.NotifyOnSessionCreated, model.NotifyOnSessionCreated)
	model.AllowBotMessages = boolPointerOrFalseState(payload.AllowBotMessages, model.AllowBotMessages)
	model.CreatedAt = types.StringValue(payload.CreatedAt)
	model.UpdatedAt = types.StringValue(payload.UpdatedAt)
	model.ResponseJSON = types.StringValue(normalizeJSON(body))
}

func slackbotSessionConfigPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) *slackbotSessionConfig {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var cfg slackbotSessionConfigModel
	diags.Append(value.As(ctx, &cfg, objectAsOptions)...)
	if diags.HasError() {
		return nil
	}
	return &slackbotSessionConfig{
		InitialMessageTemplate: optionalString(cfg.InitialMessageTemplate),
		ReuseMessageTemplate:   optionalString(cfg.ReuseMessageTemplate),
		Tags:                   stringMapFromTerraform(ctx, cfg.Tags, diags),
		Environment:            stringMapFromTerraform(ctx, cfg.Environment, diags),
		Params:                 slackbotParamsPayloadFromTerraform(ctx, cfg.Params, diags),
		MemoryKey:              stringMapFromTerraform(ctx, cfg.MemoryKey, diags),
	}
}

func slackbotParamsPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) *slackbotParamsPayload {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var params slackbotParamsModel
	diags.Append(value.As(ctx, &params, objectAsOptions)...)
	if diags.HasError() {
		return nil
	}
	return &slackbotParamsPayload{
		AgentType:    optionalString(params.AgentType),
		Oneshot:      optionalBoolPointer(params.Oneshot),
		AuthProxy:    optionalBoolPointer(params.AuthProxy),
		RepoFullName: optionalString(params.RepoFullName),
	}
}

func slackbotSessionConfigToTerraform(ctx context.Context, payload *slackbotSessionConfig, current types.Object, diags *diag.Diagnostics) types.Object {
	if payload == nil {
		if isEmptyObject(current) {
			return current
		}
		return types.ObjectNull(slackbotSessionConfigAttrTypes())
	}
	var currentConfig slackbotSessionConfigModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentConfig, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(slackbotSessionConfigAttrTypes())
		}
	}
	cfg := slackbotSessionConfigModel{
		InitialMessageTemplate: stringOrEmptyState(payload.InitialMessageTemplate, currentConfig.InitialMessageTemplate),
		ReuseMessageTemplate:   stringOrEmptyState(payload.ReuseMessageTemplate, currentConfig.ReuseMessageTemplate),
		Tags:                   mapToTerraformPreserveEmpty(ctx, payload.Tags, currentConfig.Tags, diags),
		Environment:            mapToTerraformPreserveEmpty(ctx, payload.Environment, currentConfig.Environment, diags),
		Params:                 slackbotParamsToTerraform(ctx, payload.Params, currentConfig.Params, diags),
		MemoryKey:              mapToTerraformPreserveEmpty(ctx, payload.MemoryKey, currentConfig.MemoryKey, diags),
	}
	value, d := types.ObjectValueFrom(ctx, slackbotSessionConfigAttrTypes(), cfg)
	diags.Append(d...)
	return value
}

func slackbotParamsToTerraform(ctx context.Context, payload *slackbotParamsPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	if payload == nil {
		if isEmptyObject(current) {
			return current
		}
		return types.ObjectNull(slackbotParamsAttrTypes())
	}
	var currentParams slackbotParamsModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentParams, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(slackbotParamsAttrTypes())
		}
	}
	params := slackbotParamsModel{
		AgentType:    stringOrEmptyState(payload.AgentType, currentParams.AgentType),
		Oneshot:      boolPointerOrFalseState(payload.Oneshot, currentParams.Oneshot),
		AuthProxy:    boolPointerOrFalseState(payload.AuthProxy, currentParams.AuthProxy),
		RepoFullName: stringOrEmptyState(payload.RepoFullName, currentParams.RepoFullName),
	}
	value, d := types.ObjectValueFrom(ctx, slackbotParamsAttrTypes(), params)
	diags.Append(d...)
	return value
}

func slackbotSessionConfigAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers:       useStateObject(),
		MarkdownDescription: "Session configuration for sessions created by this SlackBot.",
		Attributes: map[string]schema.Attribute{
			"initial_message_template": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Template for initial session messages."},
			"reuse_message_template":   schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Template for reused-session messages."},
			"tags":                     stringMapAttribute("Tags applied to created sessions."),
			"environment":              sensitiveStringMapAttribute("Environment variables for created sessions."),
			"params": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       useStateObject(),
				MarkdownDescription: "SlackBot session parameters.",
				Attributes:          genericSessionParamsAttributes(false),
			},
			"memory_key": stringMapAttribute("Memory lookup tags injected into sessions."),
		},
	}
}

func slackbotSessionConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"initial_message_template": types.StringType,
		"reuse_message_template":   types.StringType,
		"tags":                     types.MapType{ElemType: types.StringType},
		"environment":              types.MapType{ElemType: types.StringType},
		"params":                   types.ObjectType{AttrTypes: slackbotParamsAttrTypes()},
		"memory_key":               types.MapType{ElemType: types.StringType},
	}
}

func slackbotParamsAttrTypes() map[string]attr.Type {
	return genericSessionParamsAttrTypes(false)
}

func stringMapAttribute(description string) schema.MapAttribute {
	return schema.MapAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		Computed:            true,
		PlanModifiers:       useStateMap(),
		MarkdownDescription: description,
	}
}

func sensitiveStringMapAttribute(description string) schema.MapAttribute {
	attribute := stringMapAttribute(description)
	attribute.Sensitive = true
	return attribute
}

func stringListAttribute(description string) schema.ListAttribute {
	return schema.ListAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		Computed:            true,
		PlanModifiers:       useStateList(),
		MarkdownDescription: description,
	}
}

func (p *slackbotSessionConfig) UnmarshalJSON(data []byte) error {
	type alias slackbotSessionConfig
	var raw struct {
		alias
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = slackbotSessionConfig(raw.alias)
	if len(raw.Params) > 0 && string(raw.Params) != "null" {
		var params slackbotParamsPayload
		if err := json.Unmarshal(raw.Params, &params); err != nil {
			return err
		}
		p.Params = &params
	}
	return nil
}
