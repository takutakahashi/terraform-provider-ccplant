# ccplant_webhook

Manages a webhook through `POST /webhooks`, `GET /webhooks/{id}`,
`PUT /webhooks/{id}`, and `DELETE /webhooks/{id}`.

## Example Usage

```hcl
resource "ccplant_webhook" "example" {
  body_json = jsonencode({
    name             = "github-pr-review"
    scope            = "user"
    type             = "github"
    secret           = var.webhook_secret
    signature_header = "X-Hub-Signature-256"
    signature_type   = "hmac"
    signature_prefix = "sha256="
    github = {
      allowed_events       = ["pull_request"]
      allowed_repositories = ["org/repo"]
    }
    triggers = [{
      name          = "opened-pr"
      priority      = 10
      enabled       = true
      stop_on_match = true
      conditions = {
        github = {
          events  = ["pull_request"]
          actions = ["opened", "reopened", "synchronize"]
        }
      }
      session_config = {
        initial_message_template = "Review PR #{{ .pull_request.number }}"
      }
    }]
    session_config = {
      tags = {
        managed_by = "terraform"
      }
      params = {
        agent_type = "claude"
        oneshot    = false
      }
    }
    max_sessions = 3
  })
}
```

## Terraform Schema

### Required

- `body_json` (String) JSON request body.

### Computed

- `id` (String) Webhook ID.
- `response_json` (String) Latest API response JSON.

## `body_json` Parameters

### Required

- `name` (String) Webhook name.
- `type` (String) `github` or `custom`.
- `triggers` (List of Object) At least one trigger.

### Optional

- `scope` (String) `user` or `team`.
- `team_id` (String) Required when `scope` is `team`.
- `secret` (String) Webhook secret. Sensitive in practice, but this initial
  provider stores `body_json` as a normal Terraform string.
- `signature_header` (String) Header containing the signature.
- `signature_type` (String) `hmac` or `static`.
- `signature_prefix` (String) Prefix stripped before verification.
- `github` (Object) GitHub webhook config.
- `session_config` (Object) Default session config.
- `max_sessions` (Number) Maximum concurrent sessions.

### Optional `github` Fields

- `enterprise_url` (String) GitHub Enterprise URL.
- `allowed_events` (List of String) Allowed GitHub event names.
- `allowed_repositories` (List of String) Allowed repositories.

### Required `triggers` Fields

- `name` (String) Trigger name.

### Optional `triggers` Fields

- `id` (String) Trigger ID. Generated when omitted.
- `priority` (Number) Trigger priority.
- `enabled` (Boolean) Whether the trigger is enabled.
- `conditions` (Object) Trigger conditions.
- `session_config` (Object) Trigger-specific session config.
- `stop_on_match` (Boolean) Stop evaluating after this trigger matches.

### Optional `conditions` Fields

- `go_template` (String) Go template expression.
- `github` (Object) GitHub condition fields:
  - `events`
  - `actions`
  - `branches`
  - `repositories`
  - `labels`
  - `paths`
  - `base_branches`
  - `draft`
  - `sender`

### Optional `session_config` Fields

- `environment` (Map of String) Environment variables.
- `tags` (Map of String) Session tags.
- `initial_message_template` (String) Template for new sessions.
- `reuse_message_template` (String) Template for reused sessions.
- `params` (Object) Session parameters. Common fields include `message`,
  `agent_type`, `oneshot`, and `auth_proxy`.
- `reuse_session` (Boolean) Reuse matching sessions.
- `mount_payload` (Boolean) Mount webhook payload in the session.
- `session_profile_id` (String) Session profile ID.

## Import

Import by webhook ID:

```bash
terraform import ccplant_webhook.example <webhook-id>
```

