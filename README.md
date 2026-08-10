# terraform-provider-ccplant

Terraform provider for managing long-lived resources exposed by
agentapi-proxy.

## Status

This repository currently contains an initial provider implementation. The
first resources use JSON request bodies directly so the provider can track the
agentapi-proxy API while the typed Terraform schemas are stabilized.

## Example

```hcl
terraform {
  required_providers {
    ccplant = {
      source = "takutakahashi/ccplant"
    }
  }
}

provider "ccplant" {
  endpoint = var.agentapi_proxy_endpoint
  api_key  = var.agentapi_proxy_api_key
}

resource "ccplant_memory" "team_notes" {
  body_json = jsonencode({
    title   = "Platform notes"
    content = "Shared operating notes"
    scope   = "team"
    team_id = "org/platform"
    tags = {
      managed_by = "terraform"
    }
  })
}

resource "ccplant_user_settings" "me" {
  user_id = "alice"
  body_json = jsonencode({
    auth_mode       = "oauth"
    enabled_plugins = ["commit@claude-plugins-official"]
  })
}

resource "ccplant_team_settings" "platform" {
  team_id = "ccplant/platform"
  body_json = jsonencode({
    env_vars = {
      ENVIRONMENT = "production"
    }
  })
}

data "ccplant_sessions" "all" {}
```

Provider configuration may also be supplied with `CCPLANT_ENDPOINT` and
`CCPLANT_API_KEY`.

## Resources

- `ccplant_user_settings`
- `ccplant_team_settings`
- `ccplant_webhook`
- `ccplant_schedule`
- `ccplant_slackbot`
- `ccplant_memory`
- `ccplant_session_profile`
- `ccplant_sandbox_policy`

The JSON-backed collection resources accept:

- `body_json`: JSON request body for create and update.
- `id`: resource ID returned by agentapi-proxy.
- `response_json`: latest JSON response from agentapi-proxy.

Settings resources use `user_id` or `team_id` as their stable identifier.
Their `body_json` attribute is marked sensitive because settings can contain
credentials. They can be imported with that identifier, for example:

```shell
terraform import ccplant_user_settings.me alice
terraform import ccplant_team_settings.platform ccplant/platform
```

## Data Sources

- `ccplant_user_info`
- `ccplant_sessions`
- `ccplant_settings_managers`
