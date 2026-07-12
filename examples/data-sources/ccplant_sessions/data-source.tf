data "ccplant_sessions" "all" {}

output "sessions_json" {
  value = data.ccplant_sessions.all.response_json
}
