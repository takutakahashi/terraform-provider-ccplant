resource "ccplant_team_settings" "example" {
  team_id = "ccplant/platform"

  body_json = jsonencode({
    enabled_plugins = ["commit@claude-plugins-official"]
    env_vars = {
      MANAGED_BY = "terraform"
    }
  })
}
