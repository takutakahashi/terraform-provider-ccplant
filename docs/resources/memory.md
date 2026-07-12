# ccplant_memory

Manages a memory entry through `POST /memories`, `GET /memories/{id}`,
`PUT /memories/{id}`, and `DELETE /memories/{id}`.

## Example Usage

```hcl
resource "ccplant_memory" "example" {
  title   = "Terraform managed memory"
  content = "This memory entry is managed by terraform-provider-ccplant."
  scope   = "user"

  tags = {
    managed_by = "terraform"
  }
}
```

## Schema

### Required

- `title` (String) Memory title.
- `content` (String) Memory content.
- `scope` (String) Ownership scope: `user` or `team`. This field is only sent
  on create; changing it requires replacement.

### Optional

- `team_id` (String) Team identifier. Required by agentapi-proxy when
  `scope` is `team`. This field is only sent on create; changing it requires
  replacement.
- `tags` (Map of String) Tags attached to the memory. Updates replace the
  complete tag map.

### Computed

- `id` (String) Memory ID.
- `owner_id` (String) Owner user ID.
- `created_at` (String) Creation timestamp.
- `updated_at` (String) Last update timestamp.
- `response_json` (String) Latest normalized API response.

## Import

Import by memory ID:

```bash
terraform import ccplant_memory.example <memory-id>
```
