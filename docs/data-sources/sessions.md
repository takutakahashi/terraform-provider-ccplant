# ccplant_sessions

Reads current sessions from `GET /search`.

## Example Usage

```hcl
data "ccplant_sessions" "all" {}

output "sessions_json" {
  value = data.ccplant_sessions.all.response_json
}
```

## Schema

### Computed

- `response_json` (String) JSON response returned by agentapi-proxy.

## Parameters

This data source has no configurable parameters.

