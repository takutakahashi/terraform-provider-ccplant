resource "ccplant_session_profile" "example" {
  body_json = jsonencode({
    name        = "terraform-example-profile"
    description = "Session profile managed by terraform-provider-ccplant."
    scope       = "user"
    selector_tags = {
      managed_by = "terraform"
      profile    = "example"
    }
    config = {
      tags = {
        managed_by = "terraform"
      }
      params = {
        agent_type = "claude"
        oneshot    = false
      }
    }
  })
}
