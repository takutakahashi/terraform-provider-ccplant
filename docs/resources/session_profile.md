# ccplant_session_profile

Manages a session profile through `POST /session-profiles`,
`GET /session-profiles/{id}`, `PUT /session-profiles/{id}`, and
`DELETE /session-profiles/{id}`.

## Example Usage

```hcl
resource "ccplant_session_profile" "example" {
  name        = "default-claude-profile"
  description = "Default Claude session profile."
  scope       = "user"
  is_default  = true

  selector_tags = {
    profile = "default"
  }

  config = {
    environment = {
      LOG_LEVEL = "info"
    }
    tags = {
      managed_by = "terraform"
    }
    initial_message_template = "Start working on {{ .message }}"
    reuse_message_template   = "Continue working on {{ .message }}"
    params = {
      message    = "Start from the profile template."
      agent_type = "claude"
      oneshot    = false
      auth_proxy = true
    }
    reuse_session       = true
    memory_key          = { project = "ccplant" }
    sandbox_policy_id   = "sandbox-policy-id"
    session_ttl         = "24h"
    unsynced_file_paths = [".env"]
  }
}
```

### Team Scope

```hcl
resource "ccplant_session_profile" "team_example" {
  name        = "terraform-example-team-profile"
  description = "Team-scoped session profile managed by terraform-provider-ccplant."
  scope       = "team"
  team_id     = var.team_id

  selector_tags = {
    managed_by = "terraform"
    profile    = "team-example"
  }

  config = {
    tags = {
      managed_by = "terraform"
      scope      = "team"
    }
    params = {
      message    = "Start from the Terraform-managed team profile."
      agent_type = "claude"
      oneshot    = false
      auth_proxy = true
    }
  }
}
```

`scope = "team"` makes the session profile owned by `team_id`.

## Terraform Schema

### Required

- `name` (String) Profile name.
- `config` (Object) Session profile config. May be `{}`.

### Optional

- `description` (String) Profile description.
- `scope` (String) `user` or `team`. Create-only; changing it replaces the resource.
- `team_id` (String) Required when `scope` is `team`. Create-only; changing it replaces the resource.
- `is_default` (Boolean) Whether this is the default profile.
- `selector_tags` (Map of String) Tags used to select this profile.

### Optional `config` Fields

- `environment` (Map of String) Environment variables.
- `tags` (Map of String) Tags for created sessions.
- `initial_message_template` (String) Template for new sessions.
- `reuse_message_template` (String) Template for reused sessions.
- `params` (Object) Session parameters.
- `reuse_session` (Boolean) Reuse matching sessions.
- `memory_key` (Map of String) Memory lookup tags.
- `sandbox_policy_id` (String) Sandbox policy ID.
- `session_ttl` (String) Session TTL such as `24h`.
- `unsynced_file_paths` (List of String) Paths excluded from sync.

### Optional `config.params` Fields

- `message` (String) Initial message.
- `agent_type` (String) Agent type.
- `oneshot` (Boolean) Whether the session runs in one-shot mode.
- `auth_proxy` (Boolean) Whether auth proxy behavior is enabled.
- `repo_full_name` (String) GitHub repository full name.

### Computed

- `id` (String) Session profile ID.
- `user_id` (String) Owner user ID.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
- `response_json` (String) Latest normalized API response JSON.

## Import

Import by session profile ID:

```bash
terraform import ccplant_session_profile.example <session-profile-id>
```
