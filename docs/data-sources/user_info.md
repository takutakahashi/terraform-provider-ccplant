# ccplant_user_info

Reads authenticated user information from `GET /user/info`.

## Example Usage

```hcl
data "ccplant_user_info" "current" {}

output "user_info_json" {
  value = data.ccplant_user_info.current.response_json
}
```

## Schema

### Computed

- `response_json` (String) JSON response returned by agentapi-proxy.

## Parameters

This data source has no configurable parameters.

