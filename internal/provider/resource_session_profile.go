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
	_ resource.Resource                = (*sessionProfileResource)(nil)
	_ resource.ResourceWithConfigure   = (*sessionProfileResource)(nil)
	_ resource.ResourceWithImportState = (*sessionProfileResource)(nil)
)

type sessionProfileResource struct{ baseResource }

type sessionProfileModel struct {
	baseModel
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Scope        types.String `tfsdk:"scope"`
	TeamID       types.String `tfsdk:"team_id"`
	IsDefault    types.Bool   `tfsdk:"is_default"`
	SelectorTags types.Map    `tfsdk:"selector_tags"`
	Config       types.Object `tfsdk:"config"`
	UserID       types.String `tfsdk:"user_id"`
	CreatedAt    types.String `tfsdk:"created_at"`
	UpdatedAt    types.String `tfsdk:"updated_at"`
}

type sessionProfileConfigModel struct {
	Environment            types.Map    `tfsdk:"environment"`
	Tags                   types.Map    `tfsdk:"tags"`
	InitialMessageTemplate types.String `tfsdk:"initial_message_template"`
	ReuseMessageTemplate   types.String `tfsdk:"reuse_message_template"`
	Params                 types.Object `tfsdk:"params"`
	ReuseSession           types.Bool   `tfsdk:"reuse_session"`
	MemoryKey              types.Map    `tfsdk:"memory_key"`
	SandboxPolicyID        types.String `tfsdk:"sandbox_policy_id"`
	SessionTTL             types.String `tfsdk:"session_ttl"`
	UnsyncedFilePaths      types.List   `tfsdk:"unsynced_file_paths"`
}

type sessionParamsModel struct {
	Message      types.String `tfsdk:"message"`
	AgentType    types.String `tfsdk:"agent_type"`
	Oneshot      types.Bool   `tfsdk:"oneshot"`
	AuthProxy    types.Bool   `tfsdk:"auth_proxy"`
	RepoFullName types.String `tfsdk:"repo_full_name"`
}

type sessionProfilePayload struct {
	ID           string                      `json:"id,omitempty"`
	Name         string                      `json:"name,omitempty"`
	Description  string                      `json:"description,omitempty"`
	UserID       string                      `json:"user_id,omitempty"`
	Scope        string                      `json:"scope,omitempty"`
	TeamID       string                      `json:"team_id,omitempty"`
	IsDefault    *bool                       `json:"is_default,omitempty"`
	SelectorTags map[string]string           `json:"selector_tags,omitempty"`
	Config       sessionProfileConfigPayload `json:"config"`
	CreatedAt    string                      `json:"created_at,omitempty"`
	UpdatedAt    string                      `json:"updated_at,omitempty"`
}

type sessionProfileUpdatePayload struct {
	Name         *string                      `json:"name,omitempty"`
	Description  *string                      `json:"description,omitempty"`
	IsDefault    *bool                        `json:"is_default,omitempty"`
	SelectorTags map[string]string            `json:"selector_tags,omitempty"`
	Config       *sessionProfileConfigPayload `json:"config,omitempty"`
}

type sessionProfileConfigPayload struct {
	Environment            map[string]string     `json:"environment,omitempty"`
	Tags                   map[string]string     `json:"tags,omitempty"`
	InitialMessageTemplate string                `json:"initial_message_template,omitempty"`
	ReuseMessageTemplate   string                `json:"reuse_message_template,omitempty"`
	Params                 *sessionParamsPayload `json:"params,omitempty"`
	ReuseSession           bool                  `json:"reuse_session,omitempty"`
	MemoryKey              map[string]string     `json:"memory_key,omitempty"`
	SandboxPolicyID        string                `json:"sandbox_policy_id,omitempty"`
	SessionTTL             string                `json:"session_ttl,omitempty"`
	UnsyncedFilePaths      []string              `json:"unsynced_file_paths,omitempty"`
}

type sessionParamsPayload struct {
	Message      string `json:"message,omitempty"`
	AgentType    string `json:"agent_type,omitempty"`
	Oneshot      *bool  `json:"oneshot,omitempty"`
	AuthProxy    *bool  `json:"auth_proxy,omitempty"`
	RepoFullName string `json:"repo_full_name,omitempty"`
}

func newSessionProfileResource() resource.Resource {
	return &sessionProfileResource{baseResource: baseResource{typeName: "session_profile", basePath: "/session-profiles"}}
}

func (r *sessionProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.metadata(ctx, req, resp)
}

func (r *sessionProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.configure(ctx, req, resp)
}

func (r *sessionProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an agentapi-proxy session profile.",
		Attributes: map[string]schema.Attribute{
			"id":            idAttribute(),
			"response_json": responseJSONAttribute(),
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Session profile name.",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Session profile description.",
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
			"is_default": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Whether this is the default profile.",
			},
			"selector_tags": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Tags used by agentapi-proxy to select this profile.",
			},
			"config": schema.SingleNestedAttribute{
				Required:            true,
				MarkdownDescription: "Session profile configuration.",
				Attributes:          profileConfigAttributes(),
			},
			"user_id":    computedStringAttribute("User ID that owns the session profile."),
			"created_at": computedStringAttribute("Creation timestamp returned by agentapi-proxy."),
			"updated_at": computedStringAttribute("Last update timestamp returned by agentapi-proxy."),
		},
	}
}

func (r *sessionProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan sessionProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cfg := profileConfigPayloadFromTerraform(ctx, plan.Config, &resp.Diagnostics)
	body := marshalRequest(sessionProfilePayload{
		Name:         plan.Name.ValueString(),
		Description:  optionalString(plan.Description),
		Scope:        optionalString(plan.Scope),
		TeamID:       optionalString(plan.TeamID),
		IsDefault:    optionalBoolPointer(plan.IsDefault),
		SelectorTags: stringMapFromTerraform(ctx, plan.SelectorTags, &resp.Diagnostics),
		Config:       cfg,
	}, &resp.Diagnostics)
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

func (r *sessionProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state sessionProfileModel
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

func (r *sessionProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan sessionProfileModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	cfg := profileConfigPayloadFromTerraform(ctx, plan.Config, &resp.Diagnostics)
	body := marshalRequest(sessionProfileUpdatePayload{
		Name:         optionalStringPointer(plan.Name),
		Description:  optionalStringPointer(plan.Description),
		IsDefault:    optionalBoolPointer(plan.IsDefault),
		SelectorTags: stringMapFromTerraform(ctx, plan.SelectorTags, &resp.Diagnostics),
		Config:       &cfg,
	}, &resp.Diagnostics)
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

func (r *sessionProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state sessionProfileModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Delete(ctx, r.itemPath(state.ID.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *sessionProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, req, resp)
}

func (r *sessionProfileResource) applyResponse(ctx context.Context, body []byte, model *sessionProfileModel, diags *diag.Diagnostics) {
	var payload sessionProfilePayload
	unmarshalResponse(body, &payload, diags)
	if diags.HasError() {
		return
	}
	model.ID = types.StringValue(payload.ID)
	model.Name = types.StringValue(payload.Name)
	model.Description = stringOrEmptyState(payload.Description, model.Description)
	model.UserID = types.StringValue(payload.UserID)
	model.Scope = stringOrNull(payload.Scope)
	model.TeamID = stringOrNull(payload.TeamID)
	model.IsDefault = boolPointerOrFalseState(payload.IsDefault, model.IsDefault)
	model.SelectorTags = mapToTerraformPreserveEmpty(ctx, payload.SelectorTags, model.SelectorTags, diags)
	model.Config = profileConfigToTerraform(ctx, payload.Config, model.Config, diags)
	model.CreatedAt = types.StringValue(payload.CreatedAt)
	model.UpdatedAt = types.StringValue(payload.UpdatedAt)
	model.ResponseJSON = types.StringValue(normalizeJSON(body))
}

func profileConfigPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) sessionProfileConfigPayload {
	var cfg sessionProfileConfigModel
	diags.Append(value.As(ctx, &cfg, objectAsOptions)...)
	if diags.HasError() {
		return sessionProfileConfigPayload{}
	}
	return sessionProfileConfigPayload{
		Environment:            stringMapFromTerraform(ctx, cfg.Environment, diags),
		Tags:                   stringMapFromTerraform(ctx, cfg.Tags, diags),
		InitialMessageTemplate: optionalString(cfg.InitialMessageTemplate),
		ReuseMessageTemplate:   optionalString(cfg.ReuseMessageTemplate),
		Params:                 genericSessionParamsPayloadFromTerraform(ctx, cfg.Params, diags),
		ReuseSession:           !cfg.ReuseSession.IsNull() && !cfg.ReuseSession.IsUnknown() && cfg.ReuseSession.ValueBool(),
		MemoryKey:              stringMapFromTerraform(ctx, cfg.MemoryKey, diags),
		SandboxPolicyID:        optionalString(cfg.SandboxPolicyID),
		SessionTTL:             optionalString(cfg.SessionTTL),
		UnsyncedFilePaths:      stringListFromTerraform(ctx, cfg.UnsyncedFilePaths, diags),
	}
}

func profileConfigToTerraform(ctx context.Context, payload sessionProfileConfigPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	var currentConfig sessionProfileConfigModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentConfig, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(profileConfigAttrTypes())
		}
	}
	cfg := sessionProfileConfigModel{
		Environment:            mapToTerraformPreserveEmpty(ctx, payload.Environment, currentConfig.Environment, diags),
		Tags:                   mapToTerraformPreserveEmpty(ctx, payload.Tags, currentConfig.Tags, diags),
		InitialMessageTemplate: stringOrEmptyState(payload.InitialMessageTemplate, currentConfig.InitialMessageTemplate),
		ReuseMessageTemplate:   stringOrEmptyState(payload.ReuseMessageTemplate, currentConfig.ReuseMessageTemplate),
		Params:                 genericSessionParamsToTerraform(ctx, payload.Params, currentConfig.Params, diags),
		ReuseSession:           types.BoolValue(payload.ReuseSession),
		MemoryKey:              mapToTerraformPreserveEmpty(ctx, payload.MemoryKey, currentConfig.MemoryKey, diags),
		SandboxPolicyID:        stringOrEmptyState(payload.SandboxPolicyID, currentConfig.SandboxPolicyID),
		SessionTTL:             stringOrEmptyState(payload.SessionTTL, currentConfig.SessionTTL),
		UnsyncedFilePaths:      listToTerraformPreserveEmpty(ctx, payload.UnsyncedFilePaths, currentConfig.UnsyncedFilePaths, diags),
	}
	value, d := types.ObjectValueFrom(ctx, profileConfigAttrTypes(), cfg)
	diags.Append(d...)
	return value
}

func profileConfigAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"environment": schema.MapAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Environment variables for sessions created from this profile.",
		},
		"tags": schema.MapAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Tags applied to sessions created from this profile.",
		},
		"initial_message_template": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Template for initial session messages.",
		},
		"reuse_message_template": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Template for messages sent to reused sessions.",
		},
		"params": schema.SingleNestedAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Session parameters.",
			Attributes:          genericSessionParamsAttributes(true),
		},
		"reuse_session": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Reuse matching sessions instead of creating new sessions.",
		},
		"memory_key": schema.MapAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Memory lookup tags injected into sessions.",
		},
		"sandbox_policy_id": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Sandbox policy ID applied to sessions.",
		},
		"session_ttl": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Session TTL as a Go duration string.",
		},
		"unsynced_file_paths": schema.ListAttribute{
			ElementType:         types.StringType,
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Managed file paths excluded from sync.",
		},
	}
}

func profileConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"environment":              types.MapType{ElemType: types.StringType},
		"tags":                     types.MapType{ElemType: types.StringType},
		"initial_message_template": types.StringType,
		"reuse_message_template":   types.StringType,
		"params":                   types.ObjectType{AttrTypes: genericSessionParamsAttrTypes(true)},
		"reuse_session":            types.BoolType,
		"memory_key":               types.MapType{ElemType: types.StringType},
		"sandbox_policy_id":        types.StringType,
		"session_ttl":              types.StringType,
		"unsynced_file_paths":      types.ListType{ElemType: types.StringType},
	}
}

func genericSessionParamsPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) *sessionParamsPayload {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var params sessionParamsModel
	diags.Append(value.As(ctx, &params, objectAsOptions)...)
	if diags.HasError() {
		return nil
	}
	return &sessionParamsPayload{
		Message:      optionalString(params.Message),
		AgentType:    optionalString(params.AgentType),
		Oneshot:      optionalBoolPointer(params.Oneshot),
		AuthProxy:    optionalBoolPointer(params.AuthProxy),
		RepoFullName: optionalString(params.RepoFullName),
	}
}

func genericSessionParamsToTerraform(ctx context.Context, payload *sessionParamsPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	if payload == nil {
		if isEmptyObject(current) {
			return current
		}
		return types.ObjectNull(genericSessionParamsAttrTypes(true))
	}
	var currentParams sessionParamsModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentParams, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(genericSessionParamsAttrTypes(true))
		}
	}
	params := sessionParamsModel{
		Message:      stringOrEmptyState(payload.Message, currentParams.Message),
		AgentType:    stringOrEmptyState(payload.AgentType, currentParams.AgentType),
		Oneshot:      boolPointerOrFalseState(payload.Oneshot, currentParams.Oneshot),
		AuthProxy:    boolPointerOrFalseState(payload.AuthProxy, currentParams.AuthProxy),
		RepoFullName: stringOrEmptyState(payload.RepoFullName, currentParams.RepoFullName),
	}
	value, d := types.ObjectValueFrom(ctx, genericSessionParamsAttrTypes(true), params)
	diags.Append(d...)
	return value
}

func genericSessionParamsAttributes(includeMessage bool) map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{
		"agent_type": schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Agent type passed to agentapi-proxy."},
		"oneshot":    schema.BoolAttribute{Optional: true, Computed: true, MarkdownDescription: "Whether the session should run in one-shot mode."},
		"auth_proxy": schema.BoolAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Whether auth proxy behavior is enabled for the session.",
		},
		"repo_full_name": schema.StringAttribute{
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "GitHub repository full name used when creating sessions.",
		},
	}
	if includeMessage {
		attrs["message"] = schema.StringAttribute{Optional: true, Computed: true, MarkdownDescription: "Initial message sent to the agent session."}
	}
	return attrs
}

func genericSessionParamsAttrTypes(includeMessage bool) map[string]attr.Type {
	attrs := map[string]attr.Type{
		"agent_type":     types.StringType,
		"oneshot":        types.BoolType,
		"auth_proxy":     types.BoolType,
		"repo_full_name": types.StringType,
	}
	if includeMessage {
		attrs["message"] = types.StringType
	}
	return attrs
}

func (p *sessionProfileConfigPayload) UnmarshalJSON(data []byte) error {
	type alias sessionProfileConfigPayload
	var raw struct {
		alias
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = sessionProfileConfigPayload(raw.alias)
	if len(raw.Params) > 0 && string(raw.Params) != "null" {
		var params sessionParamsPayload
		if err := json.Unmarshal(raw.Params, &params); err != nil {
			return err
		}
		p.Params = &params
	}
	return nil
}
