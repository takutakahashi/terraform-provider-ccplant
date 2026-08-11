# ccplant_settings

Manages typed user- or team-scoped agentapi-proxy settings.

## Example

```hcl
resource "ccplant_settings" "team" {
  scope = "team"
  name  = "ccplant/platform"

  bedrock = {
    enabled = true
    model   = "anthropic.claude-sonnet-4-20250514-v1:0"
  }

  mcp_servers = {
    github = {
      type    = "stdio"
      command = "npx"
      args    = ["-y", "@modelcontextprotocol/server-github"]
    }
  }
}
```

`scope` must be `user` or `team`. `name` is the user ID for user scope and the
`org/team-slug` identifier for team scope.

OAuth tokens, Bedrock credentials, MCP environment variables and headers, and
custom environment variables are marked sensitive. Terraform still stores
sensitive values in state, so use an encrypted, access-controlled backend.

## Import

Import IDs use `scope:name`:

```shell
terraform import ccplant_settings.user user:alice
terraform import ccplant_settings.team team:ccplant/platform
```
