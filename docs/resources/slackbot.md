# ccplant_slackbot

Manages a SlackBot through `POST /slackbots`, `GET /slackbots/{id}`,
`PUT /slackbots/{id}`, and `DELETE /slackbots/{id}`.

## Example Usage

```hcl
resource "ccplant_slackbot" "example" {
  name                      = "engineering-slackbot"
  scope                     = "user"
  allowed_event_types       = ["message", "app_mention"]
  allowed_channel_names     = ["eng-alerts"]
  allowed_user_ids          = ["U0123456789"]
  max_sessions              = 5
  notify_on_session_created = true
  allow_bot_messages        = false

  session_config = {
    initial_message_template = "Handle Slack event: {{ .event.text }}"
    reuse_message_template   = "Continue Slack event: {{ .event.text }}"
    tags = {
      managed_by = "terraform"
    }
    environment = {
      LOG_LEVEL = "info"
    }
    params = {
      agent_type     = "claude"
      oneshot        = false
      auth_proxy     = true
      repo_full_name = "org/repo"
    }
    memory_key = {
      channel = "eng-alerts"
    }
  }
}
```

### Team Scope

```hcl
resource "ccplant_slackbot" "team_example" {
  name    = "engineering-team-slackbot"
  scope   = "team"
  team_id = var.team_id

  teams                     = [var.team_id]
  allowed_event_types       = ["message", "app_mention"]
  allowed_channel_names     = ["eng-alerts"]
  max_sessions              = 5
  notify_on_session_created = true
  allow_bot_messages        = false

  session_config = {
    initial_message_template = "Handle team Slack event: {{ .event.text }}"
    tags = {
      managed_by = "terraform"
      scope      = "team"
    }
    params = {
      agent_type = "claude"
      oneshot    = false
      auth_proxy = true
    }
  }
}
```

`scope = "team"` makes the SlackBot owned by `team_id`. `teams` controls which
team settings are merged into sessions created by the bot.

## Terraform Schema

### Required

- `name` (String) SlackBot name.

### Optional

- `scope` (String) `user` or `team`. Create-only; changing it replaces the resource.
- `team_id` (String) Required when `scope` is `team`. Create-only; changing it replaces the resource.
- `teams` (List of String) Team IDs whose settings are merged into sessions.
- `status` (String) SlackBot status.
- `bot_token_secret_name` (String) Kubernetes Secret name containing Slack tokens.
- `bot_token_secret_key` (String) Secret key for the bot token.
- `allowed_event_types` (List of String) Allowed Slack event types.
- `allowed_channel_names` (List of String) Allowed Slack channel name patterns.
- `allowed_user_ids` (List of String) Allowed Slack user IDs.
- `session_config` (Object) Session config for sessions created by this bot.
- `max_sessions` (Number) Maximum concurrent sessions.
- `notify_on_session_created` (Boolean) Notify in Slack when a session is created.
- `allow_bot_messages` (Boolean) Process messages from bots.
- `bot_token` (String, Sensitive) Write-only Slack bot token.
- `app_token` (String, Sensitive) Write-only Slack app token.

### Optional `session_config` Fields

- `initial_message_template` (String) Template for new sessions.
- `reuse_message_template` (String) Template for reused sessions.
- `tags` (Map of String) Session tags.
- `environment` (Map of String) Environment variables.
- `params` (Object) Session parameters.
- `memory_key` (Map of String) Memory lookup tags.

### Optional `session_config.params` Fields

- `agent_type` (String) Agent type.
- `oneshot` (Boolean) Whether the session runs in one-shot mode.
- `auth_proxy` (Boolean) Whether auth proxy behavior is enabled.
- `repo_full_name` (String) GitHub repository full name.

### Computed

- `id` (String) SlackBot ID.
- `user_id` (String) Owner user ID.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
- `response_json` (String) Latest normalized API response JSON.

## Import

Import by SlackBot ID:

```bash
terraform import ccplant_slackbot.example <slackbot-id>
```
