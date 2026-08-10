resource "ccplant_settings" "user" {
  scope = "user"
  name  = "alice"

  auth_mode       = "oauth"
  enabled_plugins = ["commit@claude-plugins-official"]

  env_vars = {
    MANAGED_BY = "terraform"
  }
}

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

  env_vars = {
    MANAGED_BY = "terraform"
  }
}
