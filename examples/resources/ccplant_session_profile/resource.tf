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

resource "ccplant_session_profile" "team_example" {
  name        = "terraform-example-team-profile"
  description = "Team-scoped session profile managed by terraform-provider-ccplant."
  scope       = "team"
  team_id     = var.team_id

  selector_tags = {
    managed_by = "terraform"
    profile    = "team-example"
  }

  config = {
    tags = {
      managed_by = "terraform"
      scope      = "team"
    }
    params = {
      message    = "Start from the Terraform-managed team profile."
      agent_type = "claude"
      oneshot    = false
      auth_proxy = true
    }
  }
}

variable "team_id" {
  type        = string
  description = "agentapi-proxy team ID, for example org/team-slug."
}
