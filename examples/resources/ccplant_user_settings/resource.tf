resource "ccplant_user_settings" "example" {
  user_id = "alice"

  body_json = jsonencode({
    auth_mode = "oauth"
    env_vars = {
      MANAGED_BY = "terraform"
    }
  })
}
