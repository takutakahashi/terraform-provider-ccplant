# ccplant_sandbox_policy

Manages a sandbox policy through `POST /sandbox-policies`,
`GET /sandbox-policies/{id}`, `PUT /sandbox-policies/{id}`, and
`DELETE /sandbox-policies/{id}`.

## Example Usage

```hcl
resource "ccplant_sandbox_policy" "example" {
  name        = "terraform-example-policy"
  description = "Sandbox policy managed by terraform-provider-ccplant."
  scope       = "user"

  allowed_domains = ["example.com"]
  denied_domains  = ["blocked.example.com"]
  count_mode      = true
}
```

### Team Scope

```hcl
resource "ccplant_sandbox_policy" "team_example" {
  name        = "terraform-example-team-policy"
  description = "Team-scoped sandbox policy managed by terraform-provider-ccplant."
  scope       = "team"
  team_id     = var.team_id

  allowed_domains = ["example.com"]
  denied_domains  = ["blocked.example.com"]
  count_mode      = true
}
```

`scope = "team"` makes the sandbox policy owned by `team_id`.

## Schema

### Required

- `name` (String) Policy name.
- `scope` (String) Ownership scope: `user` or `team`. This field is only sent
  on create; changing it requires replacement.

### Optional

- `description` (String) Policy description.
- `allowed_domains` (List of String) Domains allowed by the sandbox policy.
- `denied_domains` (List of String) Domains denied by the sandbox policy.
- `count_mode` (Boolean) Whether the policy counts matches instead of enforcing.
- `team_id` (String) Team identifier. Required by agentapi-proxy when
  `scope` is `team`. This field is only sent on create; changing it requires
  replacement.

### Computed

- `id` (String) Sandbox policy ID.
- `owner_id` (String) Owner user ID.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
- `response_json` (String, Sensitive) Latest normalized API response.

## Import

Import by sandbox policy ID:

```bash
terraform import ccplant_sandbox_policy.example <sandbox-policy-id>
```
