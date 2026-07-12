# ccplant Provider

The ccplant provider manages long-lived resources exposed by agentapi-proxy.

## Example Usage

```hcl
terraform {
  required_providers {
    ccplant = {
      source  = "takutakahashi/ccplant"
      version = "~> 0.1"
    }
  }
}

provider "ccplant" {
  endpoint = "https://agentapi-proxy.example.com"
  api_key  = var.agentapi_proxy_api_key
}
```

## Schema

### Optional

- `endpoint` (String) agentapi-proxy base URL. May also be set with
  `CCPLANT_ENDPOINT`. One of `endpoint` or `CCPLANT_ENDPOINT` is required.
- `api_key` (String, Sensitive) agentapi-proxy API key. May also be set with
  `CCPLANT_API_KEY`. Required when the target agentapi-proxy requires
  authentication.

## State and Secrets

The initial provider version uses JSON request bodies directly. Resource
`body_json` values are stored as Terraform configuration and state. Avoid
putting long-lived secrets in `body_json` when a server-side reference, such as
a Kubernetes Secret name, is available.

## Resources

- [`ccplant_memory`](resources/memory.md)
- [`ccplant_session_profile`](resources/session_profile.md)
- [`ccplant_sandbox_policy`](resources/sandbox_policy.md)
- [`ccplant_schedule`](resources/schedule.md)
- [`ccplant_webhook`](resources/webhook.md)
- [`ccplant_slackbot`](resources/slackbot.md)

## Data Sources

- [`ccplant_user_info`](data-sources/user_info.md)
- [`ccplant_sessions`](data-sources/sessions.md)
- [`ccplant_settings_managers`](data-sources/settings_managers.md)
