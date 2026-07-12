# ccplant_settings_managers

Reads available external session managers from `GET /settings/managers`.

## Example Usage

```hcl
data "ccplant_settings_managers" "all" {}

output "settings_managers_json" {
  value = data.ccplant_settings_managers.all.response_json
}
```

## Schema

### Computed

- `response_json` (String) JSON response returned by agentapi-proxy.

## Parameters

This data source has no configurable parameters.

