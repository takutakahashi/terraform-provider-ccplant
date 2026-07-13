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
	_ resource.Resource                = (*webhookResource)(nil)
	_ resource.ResourceWithConfigure   = (*webhookResource)(nil)
	_ resource.ResourceWithImportState = (*webhookResource)(nil)
)

type webhookResource struct{ baseResource }

type webhookModel struct {
	baseModel
	Name            types.String `tfsdk:"name"`
	Scope           types.String `tfsdk:"scope"`
	TeamID          types.String `tfsdk:"team_id"`
	Status          types.String `tfsdk:"status"`
	Type            types.String `tfsdk:"type"`
	Secret          types.String `tfsdk:"secret"`
	SignatureHeader types.String `tfsdk:"signature_header"`
	SignatureType   types.String `tfsdk:"signature_type"`
	SignaturePrefix types.String `tfsdk:"signature_prefix"`
	GitHub          types.Object `tfsdk:"github"`
	Triggers        types.List   `tfsdk:"triggers"`
	SessionConfig   types.Object `tfsdk:"session_config"`
	MaxSessions     types.Int64  `tfsdk:"max_sessions"`
	UserID          types.String `tfsdk:"user_id"`
	WebhookURL      types.String `tfsdk:"webhook_url"`
	DeliveryCount   types.Int64  `tfsdk:"delivery_count"`
	CreatedAt       types.String `tfsdk:"created_at"`
	UpdatedAt       types.String `tfsdk:"updated_at"`
}

type webhookGitHubModel struct {
	EnterpriseURL       types.String `tfsdk:"enterprise_url"`
	AllowedEvents       types.List   `tfsdk:"allowed_events"`
	AllowedRepositories types.List   `tfsdk:"allowed_repositories"`
}

type webhookTriggerModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Priority      types.Int64  `tfsdk:"priority"`
	Enabled       types.Bool   `tfsdk:"enabled"`
	Conditions    types.Object `tfsdk:"conditions"`
	SessionConfig types.Object `tfsdk:"session_config"`
	StopOnMatch   types.Bool   `tfsdk:"stop_on_match"`
}

type webhookTriggerConditionsModel struct {
	GitHub     types.Object `tfsdk:"github"`
	GoTemplate types.String `tfsdk:"go_template"`
}

type webhookGitHubConditionsModel struct {
	Events       types.List `tfsdk:"events"`
	Actions      types.List `tfsdk:"actions"`
	Branches     types.List `tfsdk:"branches"`
	Repositories types.List `tfsdk:"repositories"`
	Labels       types.List `tfsdk:"labels"`
	Paths        types.List `tfsdk:"paths"`
	BaseBranches types.List `tfsdk:"base_branches"`
	Draft        types.Bool `tfsdk:"draft"`
	Sender       types.List `tfsdk:"sender"`
}

type webhookSessionConfigModel struct {
	Environment            types.Map    `tfsdk:"environment"`
	Tags                   types.Map    `tfsdk:"tags"`
	InitialMessageTemplate types.String `tfsdk:"initial_message_template"`
	ReuseMessageTemplate   types.String `tfsdk:"reuse_message_template"`
	Params                 types.Object `tfsdk:"params"`
	ReuseSession           types.Bool   `tfsdk:"reuse_session"`
	MountPayload           types.Bool   `tfsdk:"mount_payload"`
	SessionProfileID       types.String `tfsdk:"session_profile_id"`
}

type webhookPayload struct {
	ID              string                       `json:"id,omitempty"`
	Name            string                       `json:"name,omitempty"`
	UserID          string                       `json:"user_id,omitempty"`
	Scope           string                       `json:"scope,omitempty"`
	TeamID          string                       `json:"team_id,omitempty"`
	Status          string                       `json:"status,omitempty"`
	Type            string                       `json:"type,omitempty"`
	Secret          string                       `json:"secret,omitempty"`
	SignatureHeader string                       `json:"signature_header,omitempty"`
	SignatureType   string                       `json:"signature_type,omitempty"`
	SignaturePrefix string                       `json:"signature_prefix,omitempty"`
	WebhookURL      string                       `json:"webhook_url,omitempty"`
	GitHub          *webhookGitHubPayload        `json:"github,omitempty"`
	Triggers        []webhookTriggerPayload      `json:"triggers,omitempty"`
	SessionConfig   *webhookSessionConfigPayload `json:"session_config,omitempty"`
	MaxSessions     int                          `json:"max_sessions,omitempty"`
	DeliveryCount   int64                        `json:"delivery_count,omitempty"`
	CreatedAt       string                       `json:"created_at,omitempty"`
	UpdatedAt       string                       `json:"updated_at,omitempty"`
}

type webhookGitHubPayload struct {
	EnterpriseURL       string   `json:"enterprise_url,omitempty"`
	AllowedEvents       []string `json:"allowed_events,omitempty"`
	AllowedRepositories []string `json:"allowed_repositories,omitempty"`
}

type webhookTriggerPayload struct {
	ID            string                          `json:"id,omitempty"`
	Name          string                          `json:"name,omitempty"`
	Priority      int                             `json:"priority,omitempty"`
	Enabled       *bool                           `json:"enabled,omitempty"`
	Conditions    webhookTriggerConditionsPayload `json:"conditions"`
	SessionConfig *webhookSessionConfigPayload    `json:"session_config,omitempty"`
	StopOnMatch   *bool                           `json:"stop_on_match,omitempty"`
}

type webhookTriggerConditionsPayload struct {
	GitHub     *webhookGitHubConditionsPayload `json:"github,omitempty"`
	GoTemplate string                          `json:"go_template,omitempty"`
}

type webhookGitHubConditionsPayload struct {
	Events       []string `json:"events,omitempty"`
	Actions      []string `json:"actions,omitempty"`
	Branches     []string `json:"branches,omitempty"`
	Repositories []string `json:"repositories,omitempty"`
	Labels       []string `json:"labels,omitempty"`
	Paths        []string `json:"paths,omitempty"`
	BaseBranches []string `json:"base_branches,omitempty"`
	Draft        *bool    `json:"draft,omitempty"`
	Sender       []string `json:"sender,omitempty"`
}

type webhookSessionConfigPayload struct {
	Environment            map[string]string     `json:"environment,omitempty"`
	Tags                   map[string]string     `json:"tags,omitempty"`
	InitialMessageTemplate string                `json:"initial_message_template,omitempty"`
	ReuseMessageTemplate   string                `json:"reuse_message_template,omitempty"`
	Params                 *sessionParamsPayload `json:"params,omitempty"`
	ReuseSession           bool                  `json:"reuse_session,omitempty"`
	MountPayload           bool                  `json:"mount_payload,omitempty"`
	SessionProfileID       string                `json:"session_profile_id,omitempty"`
}

func newWebhookResource() resource.Resource {
	return &webhookResource{baseResource: baseResource{typeName: "webhook", basePath: "/webhooks"}}
}

func (r *webhookResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.metadata(ctx, req, resp)
}

func (r *webhookResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.configure(ctx, req, resp)
}

func (r *webhookResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an agentapi-proxy webhook.",
		Attributes: map[string]schema.Attribute{
			"id":               idAttribute(),
			"response_json":    responseJSONAttribute(),
			"name":             schema.StringAttribute{Required: true, MarkdownDescription: "Webhook name."},
			"scope":            replaceableOptionalComputedString("Ownership scope: `user` or `team`. This field is only sent on create; changing it requires replacement."),
			"team_id":          replaceableOptionalComputedString("Team identifier. Required by agentapi-proxy when scope is `team`. This field is only sent on create; changing it requires replacement."),
			"status":           schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Webhook status."},
			"type":             schema.StringAttribute{Required: true, MarkdownDescription: "Webhook type: `github` or `custom`. Changing it requires replacement.", PlanModifiers: replaceOnChangeString()},
			"secret":           schema.StringAttribute{Optional: true, Computed: true, Sensitive: true, PlanModifiers: useStateString(), MarkdownDescription: "Webhook secret."},
			"signature_header": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Header containing the webhook signature."},
			"signature_type":   schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Signature verification type such as `hmac` or `static`."},
			"signature_prefix": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Prefix stripped before signature verification."},
			"github":           webhookGitHubAttribute(),
			"triggers": schema.ListNestedAttribute{
				Required:            true,
				MarkdownDescription: "Webhook trigger rules. The API requires at least one trigger.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: webhookTriggerAttributes(),
				},
			},
			"session_config": webhookSessionConfigAttribute(),
			"max_sessions":   schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: useStateInt64(), MarkdownDescription: "Maximum concurrent sessions created by this webhook."},
			"user_id":        computedStringAttribute("User ID that owns the webhook."),
			"webhook_url":    computedStringAttribute("Webhook delivery URL returned by agentapi-proxy."),
			"delivery_count": schema.Int64Attribute{Computed: true, MarkdownDescription: "Total delivery count."},
			"created_at":     computedStringAttribute("Creation timestamp returned by agentapi-proxy."),
			"updated_at":     computedStringAttribute("Last update timestamp returned by agentapi-proxy."),
		},
	}
}

func (r *webhookResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan webhookModel
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

func (r *webhookResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state webhookModel
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

func (r *webhookResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan webhookModel
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

func (r *webhookResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state webhookModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	err := r.client.Delete(ctx, r.itemPath(state.ID.ValueString()))
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Delete failed", err.Error())
	}
}

func (r *webhookResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importByID(ctx, req, resp)
}

func (r *webhookResource) request(ctx context.Context, model webhookModel, diags *diag.Diagnostics, includeCreateOnly bool) []byte {
	payload := webhookPayload{
		Name:            model.Name.ValueString(),
		Status:          optionalString(model.Status),
		Secret:          optionalString(model.Secret),
		SignatureHeader: optionalString(model.SignatureHeader),
		SignatureType:   optionalString(model.SignatureType),
		SignaturePrefix: optionalString(model.SignaturePrefix),
		GitHub:          webhookGitHubPayloadFromTerraform(ctx, model.GitHub, diags),
		Triggers:        webhookTriggersPayloadFromTerraform(ctx, model.Triggers, diags),
		SessionConfig:   webhookSessionConfigPayloadFromTerraform(ctx, model.SessionConfig, diags),
		MaxSessions:     optionalInt(model.MaxSessions),
	}
	if includeCreateOnly {
		payload.Scope = optionalString(model.Scope)
		payload.TeamID = optionalString(model.TeamID)
		payload.Type = model.Type.ValueString()
	}
	return marshalRequest(payload, diags)
}

func (r *webhookResource) applyResponse(ctx context.Context, body []byte, model *webhookModel, diags *diag.Diagnostics) {
	var payload webhookPayload
	unmarshalResponse(body, &payload, diags)
	if diags.HasError() {
		return
	}
	model.ID = types.StringValue(payload.ID)
	model.Name = types.StringValue(payload.Name)
	model.UserID = types.StringValue(payload.UserID)
	model.Scope = stringOrNull(payload.Scope)
	model.TeamID = stringOrNull(payload.TeamID)
	model.Status = stringOrNull(payload.Status)
	model.Type = types.StringValue(payload.Type)
	model.Secret = stringOrEmptyState(payload.Secret, model.Secret)
	model.SignatureHeader = stringOrEmptyState(payload.SignatureHeader, model.SignatureHeader)
	model.SignatureType = stringOrEmptyState(payload.SignatureType, model.SignatureType)
	model.SignaturePrefix = stringOrEmptyState(payload.SignaturePrefix, model.SignaturePrefix)
	model.GitHub = webhookGitHubToTerraform(ctx, payload.GitHub, model.GitHub, diags)
	model.Triggers = webhookTriggersToTerraform(ctx, payload.Triggers, model.Triggers, diags)
	model.SessionConfig = webhookSessionConfigToTerraform(ctx, payload.SessionConfig, model.SessionConfig, diags)
	model.MaxSessions = types.Int64Value(int64(payload.MaxSessions))
	model.WebhookURL = stringOrNull(payload.WebhookURL)
	model.DeliveryCount = types.Int64Value(payload.DeliveryCount)
	model.CreatedAt = types.StringValue(payload.CreatedAt)
	model.UpdatedAt = types.StringValue(payload.UpdatedAt)
	model.ResponseJSON = types.StringValue(normalizeJSON(body))
}

func webhookGitHubPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) *webhookGitHubPayload {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var gh webhookGitHubModel
	diags.Append(value.As(ctx, &gh, objectAsOptions)...)
	if diags.HasError() {
		return nil
	}
	return &webhookGitHubPayload{
		EnterpriseURL:       optionalString(gh.EnterpriseURL),
		AllowedEvents:       stringListFromTerraform(ctx, gh.AllowedEvents, diags),
		AllowedRepositories: stringListFromTerraform(ctx, gh.AllowedRepositories, diags),
	}
}

func webhookTriggersPayloadFromTerraform(ctx context.Context, value types.List, diags *diag.Diagnostics) []webhookTriggerPayload {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var triggers []webhookTriggerModel
	diags.Append(value.ElementsAs(ctx, &triggers, false)...)
	if diags.HasError() {
		return nil
	}
	payloads := make([]webhookTriggerPayload, 0, len(triggers))
	for _, trigger := range triggers {
		payloads = append(payloads, webhookTriggerPayload{
			ID:            optionalString(trigger.ID),
			Name:          trigger.Name.ValueString(),
			Priority:      optionalInt(trigger.Priority),
			Enabled:       optionalBoolPointer(trigger.Enabled),
			Conditions:    webhookTriggerConditionsPayloadFromTerraform(ctx, trigger.Conditions, diags),
			SessionConfig: webhookSessionConfigPayloadFromTerraform(ctx, trigger.SessionConfig, diags),
			StopOnMatch:   optionalBoolPointer(trigger.StopOnMatch),
		})
	}
	return payloads
}

func webhookTriggerConditionsPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) webhookTriggerConditionsPayload {
	if value.IsNull() || value.IsUnknown() {
		return webhookTriggerConditionsPayload{}
	}
	var conditions webhookTriggerConditionsModel
	diags.Append(value.As(ctx, &conditions, objectAsOptions)...)
	if diags.HasError() {
		return webhookTriggerConditionsPayload{}
	}
	return webhookTriggerConditionsPayload{
		GitHub:     webhookGitHubConditionsPayloadFromTerraform(ctx, conditions.GitHub, diags),
		GoTemplate: optionalString(conditions.GoTemplate),
	}
}

func webhookGitHubConditionsPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) *webhookGitHubConditionsPayload {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var gh webhookGitHubConditionsModel
	diags.Append(value.As(ctx, &gh, objectAsOptions)...)
	if diags.HasError() {
		return nil
	}
	return &webhookGitHubConditionsPayload{
		Events:       stringListFromTerraform(ctx, gh.Events, diags),
		Actions:      stringListFromTerraform(ctx, gh.Actions, diags),
		Branches:     stringListFromTerraform(ctx, gh.Branches, diags),
		Repositories: stringListFromTerraform(ctx, gh.Repositories, diags),
		Labels:       stringListFromTerraform(ctx, gh.Labels, diags),
		Paths:        stringListFromTerraform(ctx, gh.Paths, diags),
		BaseBranches: stringListFromTerraform(ctx, gh.BaseBranches, diags),
		Draft:        optionalBoolPointer(gh.Draft),
		Sender:       stringListFromTerraform(ctx, gh.Sender, diags),
	}
}

func webhookSessionConfigPayloadFromTerraform(ctx context.Context, value types.Object, diags *diag.Diagnostics) *webhookSessionConfigPayload {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}
	var cfg webhookSessionConfigModel
	diags.Append(value.As(ctx, &cfg, objectAsOptions)...)
	if diags.HasError() {
		return nil
	}
	return &webhookSessionConfigPayload{
		Environment:            stringMapFromTerraform(ctx, cfg.Environment, diags),
		Tags:                   stringMapFromTerraform(ctx, cfg.Tags, diags),
		InitialMessageTemplate: optionalString(cfg.InitialMessageTemplate),
		ReuseMessageTemplate:   optionalString(cfg.ReuseMessageTemplate),
		Params:                 genericSessionParamsPayloadFromTerraform(ctx, cfg.Params, diags),
		ReuseSession:           !cfg.ReuseSession.IsNull() && !cfg.ReuseSession.IsUnknown() && cfg.ReuseSession.ValueBool(),
		MountPayload:           !cfg.MountPayload.IsNull() && !cfg.MountPayload.IsUnknown() && cfg.MountPayload.ValueBool(),
		SessionProfileID:       optionalString(cfg.SessionProfileID),
	}
}

func webhookGitHubToTerraform(ctx context.Context, payload *webhookGitHubPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	if payload == nil {
		if isEmptyObject(current) {
			return current
		}
		return types.ObjectNull(webhookGitHubAttrTypes())
	}
	var currentGH webhookGitHubModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentGH, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(webhookGitHubAttrTypes())
		}
	}
	gh := webhookGitHubModel{
		EnterpriseURL:       stringOrEmptyState(payload.EnterpriseURL, currentGH.EnterpriseURL),
		AllowedEvents:       listToTerraformPreserveEmpty(ctx, payload.AllowedEvents, currentGH.AllowedEvents, diags),
		AllowedRepositories: listToTerraformPreserveEmpty(ctx, payload.AllowedRepositories, currentGH.AllowedRepositories, diags),
	}
	value, d := types.ObjectValueFrom(ctx, webhookGitHubAttrTypes(), gh)
	diags.Append(d...)
	return value
}

func webhookTriggersToTerraform(ctx context.Context, payloads []webhookTriggerPayload, current types.List, diags *diag.Diagnostics) types.List {
	if payloads == nil {
		if isEmptyList(current) {
			return current
		}
		return types.ListNull(types.ObjectType{AttrTypes: webhookTriggerAttrTypes()})
	}
	var currentTriggers []webhookTriggerModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.ElementsAs(ctx, &currentTriggers, false)...)
		if diags.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: webhookTriggerAttrTypes()})
		}
	}
	triggers := make([]webhookTriggerModel, 0, len(payloads))
	for i, payload := range payloads {
		var currentTrigger webhookTriggerModel
		if i < len(currentTriggers) {
			currentTrigger = currentTriggers[i]
		}
		triggers = append(triggers, webhookTriggerModel{
			ID:            stringOrEmptyState(payload.ID, currentTrigger.ID),
			Name:          types.StringValue(payload.Name),
			Priority:      types.Int64Value(int64(payload.Priority)),
			Enabled:       boolPointerOrFalseState(payload.Enabled, currentTrigger.Enabled),
			Conditions:    webhookTriggerConditionsToTerraform(ctx, payload.Conditions, currentTrigger.Conditions, diags),
			SessionConfig: webhookSessionConfigToTerraform(ctx, payload.SessionConfig, currentTrigger.SessionConfig, diags),
			StopOnMatch:   boolPointerOrFalseState(payload.StopOnMatch, currentTrigger.StopOnMatch),
		})
	}
	value, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: webhookTriggerAttrTypes()}, triggers)
	diags.Append(d...)
	return value
}

func webhookTriggerConditionsToTerraform(ctx context.Context, payload webhookTriggerConditionsPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	var currentConditions webhookTriggerConditionsModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentConditions, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(webhookTriggerConditionsAttrTypes())
		}
	}
	conditions := webhookTriggerConditionsModel{
		GitHub:     webhookGitHubConditionsToTerraform(ctx, payload.GitHub, currentConditions.GitHub, diags),
		GoTemplate: stringOrEmptyState(payload.GoTemplate, currentConditions.GoTemplate),
	}
	value, d := types.ObjectValueFrom(ctx, webhookTriggerConditionsAttrTypes(), conditions)
	diags.Append(d...)
	return value
}

func webhookGitHubConditionsToTerraform(ctx context.Context, payload *webhookGitHubConditionsPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	if payload == nil {
		if isEmptyObject(current) {
			return current
		}
		return types.ObjectNull(webhookGitHubConditionsAttrTypes())
	}
	var currentGH webhookGitHubConditionsModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentGH, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(webhookGitHubConditionsAttrTypes())
		}
	}
	gh := webhookGitHubConditionsModel{
		Events:       listToTerraformPreserveEmpty(ctx, payload.Events, currentGH.Events, diags),
		Actions:      listToTerraformPreserveEmpty(ctx, payload.Actions, currentGH.Actions, diags),
		Branches:     listToTerraformPreserveEmpty(ctx, payload.Branches, currentGH.Branches, diags),
		Repositories: listToTerraformPreserveEmpty(ctx, payload.Repositories, currentGH.Repositories, diags),
		Labels:       listToTerraformPreserveEmpty(ctx, payload.Labels, currentGH.Labels, diags),
		Paths:        listToTerraformPreserveEmpty(ctx, payload.Paths, currentGH.Paths, diags),
		BaseBranches: listToTerraformPreserveEmpty(ctx, payload.BaseBranches, currentGH.BaseBranches, diags),
		Draft:        boolPointerOrNull(payload.Draft),
		Sender:       listToTerraformPreserveEmpty(ctx, payload.Sender, currentGH.Sender, diags),
	}
	value, d := types.ObjectValueFrom(ctx, webhookGitHubConditionsAttrTypes(), gh)
	diags.Append(d...)
	return value
}

func webhookSessionConfigToTerraform(ctx context.Context, payload *webhookSessionConfigPayload, current types.Object, diags *diag.Diagnostics) types.Object {
	if payload == nil {
		if isEmptyObject(current) {
			return current
		}
		return types.ObjectNull(webhookSessionConfigAttrTypes())
	}
	var currentConfig webhookSessionConfigModel
	if !current.IsNull() && !current.IsUnknown() {
		diags.Append(current.As(ctx, &currentConfig, objectAsOptions)...)
		if diags.HasError() {
			return types.ObjectNull(webhookSessionConfigAttrTypes())
		}
	}
	cfg := webhookSessionConfigModel{
		Environment:            mapToTerraformPreserveEmpty(ctx, payload.Environment, currentConfig.Environment, diags),
		Tags:                   mapToTerraformPreserveEmpty(ctx, payload.Tags, currentConfig.Tags, diags),
		InitialMessageTemplate: stringOrEmptyState(payload.InitialMessageTemplate, currentConfig.InitialMessageTemplate),
		ReuseMessageTemplate:   stringOrEmptyState(payload.ReuseMessageTemplate, currentConfig.ReuseMessageTemplate),
		Params:                 genericSessionParamsToTerraform(ctx, payload.Params, currentConfig.Params, diags),
		ReuseSession:           types.BoolValue(payload.ReuseSession),
		MountPayload:           types.BoolValue(payload.MountPayload),
		SessionProfileID:       stringOrEmptyState(payload.SessionProfileID, currentConfig.SessionProfileID),
	}
	value, d := types.ObjectValueFrom(ctx, webhookSessionConfigAttrTypes(), cfg)
	diags.Append(d...)
	return value
}

func webhookGitHubAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers:       useStateObject(),
		MarkdownDescription: "GitHub-specific webhook configuration.",
		Attributes: map[string]schema.Attribute{
			"enterprise_url":       schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "GitHub Enterprise URL."},
			"allowed_events":       stringListAttribute("Allowed GitHub event names."),
			"allowed_repositories": stringListAttribute("Allowed GitHub repositories."),
		},
	}
}

func webhookTriggerAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id":             schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Trigger ID. Generated by agentapi-proxy when omitted."},
		"name":           schema.StringAttribute{Required: true, MarkdownDescription: "Trigger name."},
		"priority":       schema.Int64Attribute{Optional: true, Computed: true, PlanModifiers: useStateInt64(), MarkdownDescription: "Trigger priority."},
		"enabled":        schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: useStateBool(), MarkdownDescription: "Whether this trigger is enabled."},
		"conditions":     webhookTriggerConditionsAttribute(),
		"session_config": webhookSessionConfigAttribute(),
		"stop_on_match":  schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: useStateBool(), MarkdownDescription: "Stop evaluating triggers after this trigger matches."},
	}
}

func webhookTriggerConditionsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers:       useStateObject(),
		MarkdownDescription: "Trigger matching conditions.",
		Attributes: map[string]schema.Attribute{
			"github":      webhookGitHubConditionsAttribute(),
			"go_template": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Go template expression used as a condition."},
		},
	}
}

func webhookGitHubConditionsAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers:       useStateObject(),
		MarkdownDescription: "GitHub event condition fields.",
		Attributes: map[string]schema.Attribute{
			"events":        stringListAttribute("Allowed GitHub event names."),
			"actions":       stringListAttribute("Allowed GitHub event actions."),
			"branches":      stringListAttribute("Allowed branches."),
			"repositories":  stringListAttribute("Allowed repositories."),
			"labels":        stringListAttribute("Allowed labels."),
			"paths":         stringListAttribute("Allowed changed paths."),
			"base_branches": stringListAttribute("Allowed base branches."),
			"draft":         schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: useStateBool(), MarkdownDescription: "Whether pull request draft state must match."},
			"sender":        stringListAttribute("Allowed GitHub sender logins."),
		},
	}
}

func webhookSessionConfigAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Optional:            true,
		Computed:            true,
		PlanModifiers:       useStateObject(),
		MarkdownDescription: "Default session configuration.",
		Attributes: map[string]schema.Attribute{
			"environment":              sensitiveStringMapAttribute("Environment variables for created sessions."),
			"tags":                     stringMapAttribute("Tags applied to created sessions."),
			"initial_message_template": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Template for initial session messages."},
			"reuse_message_template":   schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Template for reused-session messages."},
			"params": schema.SingleNestedAttribute{
				Optional:            true,
				Computed:            true,
				PlanModifiers:       useStateObject(),
				MarkdownDescription: "Session parameters.",
				Attributes:          genericSessionParamsAttributes(true),
			},
			"reuse_session":      schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: useStateBool(), MarkdownDescription: "Reuse matching sessions instead of creating new sessions."},
			"mount_payload":      schema.BoolAttribute{Optional: true, Computed: true, PlanModifiers: useStateBool(), MarkdownDescription: "Mount webhook payload into the session."},
			"session_profile_id": schema.StringAttribute{Optional: true, Computed: true, PlanModifiers: useStateString(), MarkdownDescription: "Session profile ID used as the base session configuration."},
		},
	}
}

func replaceableOptionalComputedString(description string) schema.StringAttribute {
	return schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: description,
		PlanModifiers:       replaceOnChangeString(),
	}
}

func webhookGitHubAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enterprise_url":       types.StringType,
		"allowed_events":       types.ListType{ElemType: types.StringType},
		"allowed_repositories": types.ListType{ElemType: types.StringType},
	}
}

func webhookTriggerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"name":           types.StringType,
		"priority":       types.Int64Type,
		"enabled":        types.BoolType,
		"conditions":     types.ObjectType{AttrTypes: webhookTriggerConditionsAttrTypes()},
		"session_config": types.ObjectType{AttrTypes: webhookSessionConfigAttrTypes()},
		"stop_on_match":  types.BoolType,
	}
}

func webhookTriggerConditionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"github":      types.ObjectType{AttrTypes: webhookGitHubConditionsAttrTypes()},
		"go_template": types.StringType,
	}
}

func webhookGitHubConditionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"events":        types.ListType{ElemType: types.StringType},
		"actions":       types.ListType{ElemType: types.StringType},
		"branches":      types.ListType{ElemType: types.StringType},
		"repositories":  types.ListType{ElemType: types.StringType},
		"labels":        types.ListType{ElemType: types.StringType},
		"paths":         types.ListType{ElemType: types.StringType},
		"base_branches": types.ListType{ElemType: types.StringType},
		"draft":         types.BoolType,
		"sender":        types.ListType{ElemType: types.StringType},
	}
}

func webhookSessionConfigAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"environment":              types.MapType{ElemType: types.StringType},
		"tags":                     types.MapType{ElemType: types.StringType},
		"initial_message_template": types.StringType,
		"reuse_message_template":   types.StringType,
		"params":                   types.ObjectType{AttrTypes: genericSessionParamsAttrTypes(true)},
		"reuse_session":            types.BoolType,
		"mount_payload":            types.BoolType,
		"session_profile_id":       types.StringType,
	}
}

func (p *webhookSessionConfigPayload) UnmarshalJSON(data []byte) error {
	type alias webhookSessionConfigPayload
	var raw struct {
		alias
		Params json.RawMessage `json:"params"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*p = webhookSessionConfigPayload(raw.alias)
	if len(raw.Params) > 0 && string(raw.Params) != "null" {
		var params sessionParamsPayload
		if err := json.Unmarshal(raw.Params, &params); err != nil {
			return err
		}
		p.Params = &params
	}
	return nil
}
