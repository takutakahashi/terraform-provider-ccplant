# terraform-provider-ccplant

Terraform provider for managing long-lived resources exposed by
agentapi-proxy.

## Installation

After the provider is published to the Terraform Registry:

```hcl
terraform {
  required_providers {
    ccplant = {
      source  = "takutakahashi/ccplant"
      version = "~> 0.1"
    }
  }
}
```

## Provider Configuration

```hcl
provider "ccplant" {
  endpoint = var.agentapi_proxy_endpoint
  api_key  = var.agentapi_proxy_api_key
}
resource "ccplant_settings" "me" {
  scope = "user"
  name  = "alice"

  auth_mode            = "oauth"
  enabled_plugins      = ["commit@claude-plugins-official"]
  notification_channels = ["web"]
}

resource "ccplant_settings" "platform" {
  scope = "team"
  name  = "ccplant/platform"

  bedrock = {
    enabled = true
    model   = "anthropic.claude-sonnet-4-20250514-v1:0"
  }

  env_vars = {
    ENVIRONMENT = "production"
  }
}

data "ccplant_sessions" "all" {}
```

`endpoint` may also be set with `CCPLANT_ENDPOINT`.
`api_key` may also be set with `CCPLANT_API_KEY`.

See [provider documentation](docs/index.md) for the complete schema.

## Resources

Resources support create, read, update, delete, and import by agentapi-proxy
resource ID. All resources expose typed Terraform attributes based on the
agentapi-proxy API request and response schemas.

- [`ccplant_settings`](docs/resources/settings.md)
- [`ccplant_memory`](docs/resources/memory.md)
- [`ccplant_session_profile`](docs/resources/session_profile.md)
- [`ccplant_sandbox_policy`](docs/resources/sandbox_policy.md)
- [`ccplant_schedule`](docs/resources/schedule.md)
- [`ccplant_webhook`](docs/resources/webhook.md)
- [`ccplant_slackbot`](docs/resources/slackbot.md)

Managed fields are refreshed from API responses so drift is visible in plans.
Each resource also exposes `response_json` for debugging the latest normalized
API response.

Webhook secrets and Slack tokens are marked sensitive, but Terraform still
stores sensitive values in state. Prefer server-side Secret references for Slack
tokens when possible.

The settings resource uses `scope` (`user` or `team`) and `name` as its stable
identifier. For user scope, `name` is the user ID. For team scope, it is the
`org/team-slug` ID. Secret-bearing attributes such as OAuth tokens, Bedrock
credentials, MCP environment variables and custom environment variables are
marked sensitive. Import IDs use `scope:name`:

```shell
terraform import ccplant_settings.me user:alice
terraform import ccplant_settings.platform team:ccplant/platform
```

## Data Sources

- [`ccplant_user_info`](docs/data-sources/user_info.md)
- [`ccplant_sessions`](docs/data-sources/sessions.md)
- [`ccplant_settings_managers`](docs/data-sources/settings_managers.md)

## Import

All resources support import by agentapi-proxy resource ID:

```bash
terraform import ccplant_memory.example <memory-id>
terraform import ccplant_session_profile.example <session-profile-id>
terraform import ccplant_sandbox_policy.example <sandbox-policy-id>
terraform import ccplant_schedule.example <schedule-id>
terraform import ccplant_webhook.example <webhook-id>
terraform import ccplant_slackbot.example <slackbot-id>
```

After import, `terraform plan` shows any differences between configuration and
the API-populated state.

## Examples

Examples are available under:

- [`examples/provider`](examples/provider/provider.tf)
- [`examples/resources`](examples/resources)
- [`examples/data-sources`](examples/data-sources)
