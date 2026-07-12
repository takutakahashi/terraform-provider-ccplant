# ccplant_webhook

Manages a webhook through `POST /webhooks`, `GET /webhooks/{id}`,
`PUT /webhooks/{id}`, and `DELETE /webhooks/{id}`.

## Example Usage

```hcl
resource "ccplant_webhook" "example" {
  name             = "github-pr-review"
  scope            = "user"
  type             = "github"
  secret           = var.webhook_secret
  signature_header = "X-Hub-Signature-256"
  signature_type   = "hmac"
  signature_prefix = "sha256="
  max_sessions     = 3

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
      message    = "Handle GitHub webhook payload."
      agent_type = "claude"
      oneshot    = false
      auth_proxy = true
    }
  }
}

variable "webhook_secret" {
  type      = string
  sensitive = true
}
```

### Team Scope

```hcl
resource "ccplant_webhook" "team_example" {
  name             = "terraform-example-team-webhook"
  scope            = "team"
  team_id          = var.team_id
  type             = "custom"
  secret           = var.webhook_secret
  signature_header = "X-Signature"
  signature_type   = "static"
  max_sessions     = 1

  triggers = [{
    name          = "team-example-trigger"
    priority      = 1
    enabled       = true
    stop_on_match = true
    conditions = {
      go_template = "{{ true }}"
    }
  }]

  session_config = {
    tags = {
      managed_by = "terraform"
      scope      = "team"
    }
    params = {
      message    = "Handle team custom webhook payload."
      agent_type = "claude"
      oneshot    = true
      auth_proxy = true
    }
  }
}
```

`scope = "team"` makes the webhook owned by `team_id`.

## Terraform Schema

### Required

- `name` (String) Webhook name.
- `type` (String) `github` or `custom`. Create-only; changing it replaces the resource.
- `triggers` (List of Object) At least one trigger.

### Optional

- `scope` (String) `user` or `team`. Create-only; changing it replaces the resource.
- `team_id` (String) Required when `scope` is `team`. Create-only; changing it replaces the resource.
- `status` (String) Webhook status.
- `secret` (String, Sensitive) Webhook secret.
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

### `triggers` Fields

- `name` (String, Required) Trigger name.
- `id` (String, Optional/Computed) Trigger ID. Generated when omitted.
- `priority` (Number) Trigger priority.
- `enabled` (Boolean) Whether the trigger is enabled.
- `conditions` (Object) Trigger conditions.
- `session_config` (Object) Trigger-specific session config.
- `stop_on_match` (Boolean) Stop evaluating after this trigger matches.

### Optional `conditions` Fields

- `go_template` (String) Go template expression.
- `github.events` (List of String)
- `github.actions` (List of String)
- `github.branches` (List of String)
- `github.repositories` (List of String)
- `github.labels` (List of String)
- `github.paths` (List of String)
- `github.base_branches` (List of String)
- `github.draft` (Boolean)
- `github.sender` (List of String)

### Optional `session_config` Fields

- `environment` (Map of String) Environment variables.
- `tags` (Map of String) Session tags.
- `initial_message_template` (String) Template for new sessions.
- `reuse_message_template` (String) Template for reused sessions.
- `params` (Object) Session parameters.
- `reuse_session` (Boolean) Reuse matching sessions.
- `mount_payload` (Boolean) Mount webhook payload in the session.
- `session_profile_id` (String) Session profile ID.

### Optional `session_config.params` Fields

- `message` (String) Initial message.
- `agent_type` (String) Agent type.
- `oneshot` (Boolean) Whether the session runs in one-shot mode.
- `auth_proxy` (Boolean) Whether auth proxy behavior is enabled.
- `repo_full_name` (String) GitHub repository full name.

### Computed

- `id` (String) Webhook ID.
- `user_id` (String) Owner user ID.
- `webhook_url` (String) Delivery URL.
- `delivery_count` (Number) Total delivery count.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
- `response_json` (String) Latest normalized API response JSON.

## Import

Import by webhook ID:

```bash
terraform import ccplant_webhook.example <webhook-id>
```
