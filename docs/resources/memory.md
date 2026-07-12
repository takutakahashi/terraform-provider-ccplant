# ccplant_memory

Manages a memory entry through `POST /memories`, `GET /memories/{id}`,
`PUT /memories/{id}`, and `DELETE /memories/{id}`.

## Example Usage

```hcl
resource "ccplant_memory" "example" {
  body_json = jsonencode({
    title   = "Platform operating notes"
    content = "Shared notes for agent sessions."
    scope   = "team"
    team_id = "org/platform"
    tags = {
      managed_by = "terraform"
    }
  })
}
```

## Terraform Schema

### Required

- `body_json` (String) JSON request body.

### Computed

- `id` (String) Memory ID.
- `response_json` (String) Latest API response JSON.

## `body_json` Parameters

### Required

- `title` (String) Memory title.
- `content` (String) Memory content.
- `scope` (String) `user` or `team`.

### Optional

- `team_id` (String) Required when `scope` is `team`.
- `tags` (Map of String) Tags attached to the memory.

## Import

Import by memory ID:

```bash
terraform import ccplant_memory.example <memory-id>
```

