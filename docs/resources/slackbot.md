# ccplant_slackbot

Manages a SlackBot through `POST /slackbots`, `GET /slackbots/{id}`,
`PUT /slackbots/{id}`, and `DELETE /slackbots/{id}`.

## Example Usage

```hcl
resource "ccplant_slackbot" "example" {
  body_json = jsonencode({
    name                  = "engineering-slackbot"
    scope                 = "user"
    allowed_event_types   = ["message", "app_mention"]
    allowed_channel_names = ["eng-alerts"]
    allowed_user_ids      = ["U0123456789"]
    max_sessions          = 5
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
        agent_type = "claude"
        oneshot    = false
      }
      memory_key = {
        channel = "eng-alerts"
      }
    }
  })
}
```

## Terraform Schema

### Required

- `body_json` (String) JSON request body.

### Computed

- `id` (String) SlackBot ID.
- `response_json` (String) Latest API response JSON.

## `body_json` Parameters

### Required

- `name` (String) SlackBot name.

### Optional

- `scope` (String) `user` or `team`.
- `team_id` (String) Required when `scope` is `team`.
- `teams` (List of String) Team IDs whose settings are merged into sessions.
- `bot_token_secret_name` (String) Kubernetes Secret name containing Slack
  tokens.
- `bot_token_secret_key` (String) Secret key for the bot token.
- `app_token_secret_key` (String) Secret key for the app token.
- `allowed_event_types` (List of String) Allowed Slack event types.
- `allowed_channel_names` (List of String) Allowed Slack channel name patterns.
- `allowed_user_ids` (List of String) Allowed Slack user IDs.
- `session_config` (Object) Session config for sessions created by this bot.
- `max_sessions` (Number) Maximum concurrent sessions.
- `notify_on_session_created` (Boolean) Notify in Slack when a session is
  created.
- `allow_bot_messages` (Boolean) Process messages from bots.
- `bot_token` (String) Write-only Slack bot token. Sensitive in practice, but
  this initial provider stores `body_json` as a normal Terraform string.
- `app_token` (String) Write-only Slack app token. Sensitive in practice, but
  this initial provider stores `body_json` as a normal Terraform string.

### Optional `session_config` Fields

- `initial_message_template` (String) Template for new sessions.
- `reuse_message_template` (String) Template for reused sessions.
- `tags` (Map of String) Session tags.
- `environment` (Map of String) Environment variables.
- `params` (Object) Session parameters. Common fields include `agent_type`,
  `oneshot`, `auth_proxy`, and `repo_full_name`.
- `memory_key` (Map of String) Memory lookup tags.

## Import

Import by SlackBot ID:

```bash
terraform import ccplant_slackbot.example <slackbot-id>
```

