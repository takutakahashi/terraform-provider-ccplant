resource "ccplant_session_profile" "example" {
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
      message    = "Start from the Terraform-managed profile."
      agent_type = "claude"
      oneshot    = false
      auth_proxy = true
    }
  }
}
