# ccplant_sandbox_policy

Manages a sandbox policy through `POST /sandbox-policies`,
`GET /sandbox-policies/{id}`, `PUT /sandbox-policies/{id}`, and
`DELETE /sandbox-policies/{id}`.

## Example Usage

```hcl
resource "ccplant_sandbox_policy" "example" {
  body_json = jsonencode({
    name            = "restricted-network"
    description     = "Allow only approved domains."
    scope           = "user"
    allowed_domains = ["example.com"]
    denied_domains  = ["blocked.example.com"]
    count_mode      = true
  })
}
```

## Terraform Schema

### Required

- `body_json` (String) JSON request body.

### Computed

- `id` (String) Sandbox policy ID.
- `response_json` (String) Latest API response JSON.

## `body_json` Parameters

### Required

- `name` (String) Policy name.
- `scope` (String) `user` or `team`.

### Optional

- `description` (String) Policy description.
- `allowed_domains` (List of String) Allowed domains.
- `denied_domains` (List of String) Denied domains.
- `count_mode` (Boolean) Count policy hits instead of enforcing.
- `team_id` (String) Required when `scope` is `team`.

## Import

Import by sandbox policy ID:

```bash
terraform import ccplant_sandbox_policy.example <sandbox-policy-id>
```

