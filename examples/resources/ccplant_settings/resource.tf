resource "ccplant_settings" "user" {
  scope = "user"
  name  = "alice"

  body_json = jsonencode({
    auth_mode = "oauth"
    env_vars = {
      MANAGED_BY = "terraform"
    }
  })
}

resource "ccplant_settings" "team" {
  scope = "team"
  name  = "ccplant/platform"

  body_json = jsonencode({
    enabled_plugins = ["commit@claude-plugins-official"]
    env_vars = {
      MANAGED_BY = "terraform"
    }
  })
}
