data "ccplant_settings_managers" "all" {}

output "settings_managers_json" {
  value = data.ccplant_settings_managers.all.response_json
}
