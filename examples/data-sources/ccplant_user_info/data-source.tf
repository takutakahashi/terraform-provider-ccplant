data "ccplant_user_info" "current" {}

output "user_info_json" {
  value = data.ccplant_user_info.current.response_json
}
