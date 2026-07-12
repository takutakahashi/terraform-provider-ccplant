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
```

`endpoint` may also be set with `CCPLANT_ENDPOINT`.
`api_key` may also be set with `CCPLANT_API_KEY`.

See [provider documentation](docs/index.md) for the complete schema.

## Resources

Resources support create, read, update, delete, and import by agentapi-proxy
resource ID. `ccplant_memory`, `ccplant_sandbox_policy`, and
`ccplant_schedule` expose typed Terraform attributes. `ccplant_webhook`,
`ccplant_slackbot`, and `ccplant_session_profile` currently use `body_json`
request bodies until they are migrated.

- [`ccplant_memory`](docs/resources/memory.md)
- [`ccplant_session_profile`](docs/resources/session_profile.md)
- [`ccplant_sandbox_policy`](docs/resources/sandbox_policy.md)
- [`ccplant_schedule`](docs/resources/schedule.md)
- [`ccplant_webhook`](docs/resources/webhook.md)
- [`ccplant_slackbot`](docs/resources/slackbot.md)

For typed resources, managed fields are refreshed from API responses so drift is
visible in plans. JSON resources expose `body_json` and `response_json`; drift
is visible in `response_json`, but Terraform compares desired changes through
`body_json`.

Sensitive values placed in JSON resource `body_json`, such as webhook secrets or
Slack tokens, live in Terraform state. Prefer server-side Secret references for
Slack tokens when possible.

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
