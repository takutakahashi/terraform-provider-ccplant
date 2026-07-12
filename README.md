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

The initial provider version uses JSON request bodies directly. Each resource
has the same Terraform-level attributes:

- `body_json` (required): JSON body sent to agentapi-proxy on create and update.
- `id` (computed): resource ID assigned by agentapi-proxy.
- `response_json` (computed): latest JSON response from agentapi-proxy.

Because `body_json` is a normal Terraform string in the initial provider
version, values such as webhook secrets or Slack tokens can be stored in
Terraform state if you put them there. Prefer references to server-side
Kubernetes Secrets when possible.

Supported resources:

- [`ccplant_memory`](docs/resources/memory.md)
- [`ccplant_session_profile`](docs/resources/session_profile.md)
- [`ccplant_sandbox_policy`](docs/resources/sandbox_policy.md)
- [`ccplant_schedule`](docs/resources/schedule.md)
- [`ccplant_webhook`](docs/resources/webhook.md)
- [`ccplant_slackbot`](docs/resources/slackbot.md)

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

After import, keep the matching `body_json` in configuration and run
`terraform plan`. Terraform may need one update after import to store
`body_json` in state.

## Examples

Examples are available under:

- [`examples/provider`](examples/provider/provider.tf)
- [`examples/resources`](examples/resources)
- [`examples/data-sources`](examples/data-sources)
