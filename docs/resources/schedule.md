# ccplant_schedule

Manages a schedule through `POST /schedules`, `GET /schedules/{id}`,
`PUT /schedules/{id}`, and `DELETE /schedules/{id}`.

## Example Usage

```hcl
resource "ccplant_schedule" "example" {
  body_json = jsonencode({
    name      = "weekday-review"
    scope     = "user"
    cron_expr = "0 9 * * 1-5"
    timezone  = "UTC"
    session_config = {
      tags = {
        managed_by = "terraform"
      }
      params = {
        message    = "Run the weekday review."
        agent_type = "claude"
        oneshot    = true
      }
    }
  })
}
```

## Terraform Schema

### Required

- `body_json` (String) JSON request body.

### Computed

- `id` (String) Schedule ID.
- `response_json` (String) Latest API response JSON.

## `body_json` Parameters

### Required

- `name` (String) Schedule name.
- One of:
  - `scheduled_at` (String) RFC3339 timestamp for a one-time schedule.
  - `cron_expr` (String) Five-field cron expression for recurring schedules.
- `session_config` (Object) Session creation config. May be `{}`.

### Optional

- `scope` (String) `user` or `team`.
- `team_id` (String) Required when `scope` is `team`.
- `timezone` (String) IANA timezone. Defaults server-side when omitted.

### Optional `session_config` Fields

- `environment` (Map of String) Environment variables.
- `tags` (Map of String) Session tags.
- `params` (Object) Session parameters. Common fields include `message`,
  `agent_type`, `oneshot`, and `auth_proxy`.
- `memory_key` (Map of String) Memory lookup tags.
- `reuse_session` (Boolean) Reuse an existing active session.
- `reuse_message` (String) Message sent to reused sessions.
- `session_profile_id` (String) Session profile ID.

## Import

Import by schedule ID:

```bash
terraform import ccplant_schedule.example <schedule-id>
```

