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

data "ccplant_sessions" "all" {}
```

Provider configuration may also be supplied with `CCPLANT_ENDPOINT` and
`CCPLANT_API_KEY`.

## Resources

- `ccplant_webhook`
- `ccplant_schedule`
- `ccplant_slackbot`
- `ccplant_memory`
- `ccplant_session_profile`
- `ccplant_sandbox_policy`

Each resource currently accepts:

- `body_json`: JSON request body for create and update.
- `id`: resource ID returned by agentapi-proxy.
- `response_json`: latest JSON response from agentapi-proxy.

## Data Sources

- `ccplant_user_info`
- `ccplant_sessions`
- `ccplant_settings_managers`
