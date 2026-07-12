# ccplant_session_profile

Manages a session profile through `POST /session-profiles`,
`GET /session-profiles/{id}`, `PUT /session-profiles/{id}`, and
`DELETE /session-profiles/{id}`.

## Example Usage

```hcl
resource "ccplant_session_profile" "example" {
  body_json = jsonencode({
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
        agent_type = "claude"
        oneshot    = false
      }
      reuse_session      = true
      memory_key         = { project = "ccplant" }
      sandbox_policy_id  = ""
      session_ttl        = "24h"
      unsynced_file_paths = [
        ".env",
      ]
    }
  })
}
```

## Terraform Schema

### Required

- `body_json` (String) JSON request body.

### Computed

- `id` (String) Session profile ID.
- `response_json` (String) Latest API response JSON.

## `body_json` Parameters

### Required

- `name` (String) Profile name.
- `config` (Object) Session profile config. May be `{}`.

### Optional

- `description` (String) Profile description.
- `scope` (String) `user` or `team`. Defaults to server-side behavior when
  omitted.
- `team_id` (String) Required when `scope` is `team`.
- `is_default` (Boolean) Whether this is the default profile.
- `selector_tags` (Map of String) Tags used to select this profile.

### Optional `config` Fields

- `environment` (Map of String) Environment variables.
- `tags` (Map of String) Tags for created sessions.
- `initial_message_template` (String) Template for new sessions.
- `reuse_message_template` (String) Template for reused sessions.
- `params` (Object) Session parameters. Common fields include `message`,
  `agent_type`, `oneshot`, and `auth_proxy`.
- `reuse_session` (Boolean) Reuse matching sessions.
- `memory_key` (Map of String) Memory lookup tags.
- `sandbox_policy_id` (String) Sandbox policy ID.
- `session_ttl` (String) Session TTL such as `24h`.
- `unsynced_file_paths` (List of String) Paths excluded from sync.

## Import

Import by session profile ID:

```bash
terraform import ccplant_session_profile.example <session-profile-id>
```

