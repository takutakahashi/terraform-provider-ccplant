# ccplant_schedule

Manages a schedule through `POST /schedules`, `GET /schedules/{id}`,
`PUT /schedules/{id}`, and `DELETE /schedules/{id}`.

## Example Usage

```hcl
resource "ccplant_schedule" "example" {
  name      = "terraform-example-schedule"
  scope     = "user"
  cron_expr = "0 9 * * *"
  timezone  = "UTC"

  session_config = {
    tags = {
      managed_by = "terraform"
    }
    params = {
      message    = "Run the scheduled task."
      agent_type = "claude"
      oneshot    = true
    }
  }
}
```

### Team Scope

```hcl
resource "ccplant_schedule" "team_example" {
  name      = "terraform-example-team-schedule"
  scope     = "team"
  team_id   = var.team_id
  cron_expr = "0 9 * * *"
  timezone  = "UTC"

  session_config = {
    tags = {
      managed_by = "terraform"
      scope      = "team"
    }
    params = {
      message    = "Run the scheduled team task."
      agent_type = "claude"
      oneshot    = true
    }
  }
}
```

`scope = "team"` makes the schedule owned by `team_id`.

## Schema

### Required

- `name` (String) Schedule name.
- `session_config` (Object) Session creation configuration.

The API requires at least one of:

- `scheduled_at` (String) RFC3339 timestamp for a one-time schedule.
- `cron_expr` (String) Five-field cron expression for recurring schedules.

### Optional

- `scope` (String) Ownership scope: `user` or `team`. This field is only sent
  on create; changing it requires replacement.
- `team_id` (String) Team identifier. Required by agentapi-proxy when
  `scope` is `team`. This field is only sent on create; changing it requires
  replacement.
- `status` (String) Schedule status: `active`, `paused`, or `completed`.
- `scheduled_at` (String) One-time or first execution time as RFC3339.
- `cron_expr` (String) Recurring cron expression.
- `timezone` (String) IANA timezone.

### Computed

- `id` (String) Schedule ID.
- `user_id` (String) Owner user ID.
- `next_execution_at` (String) Next execution time.
- `execution_count` (Number) Total execution count.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
- `response_json` (String, Sensitive) Latest normalized API response.

## `session_config` Fields

Optional fields:

- `environment` (Map of String, Sensitive) Environment variables.
- `tags` (Map of String) Session tags.
- `params` (Object) Session parameters.
- `memory_key` (Map of String) Memory lookup tags.
- `reuse_session` (Boolean) Reuse an existing active session.
- `reuse_message` (String) Message sent to reused sessions.
- `session_profile_id` (String) Session profile ID.

Optional `params` fields:

- `message` (String) Initial message sent to the agent session.
- `agent_type` (String) Agent type passed to agentapi-proxy.
- `oneshot` (Boolean) Whether the session should run in one-shot mode.
- `auth_proxy` (Boolean) Whether auth proxy behavior is enabled for the
  session.

## Import

Import by schedule ID:

```bash
terraform import ccplant_schedule.example <schedule-id>
```
