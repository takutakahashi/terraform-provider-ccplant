resource "ccplant_schedule" "example" {
  name      = "terraform-example-schedule"
  scope     = "user"
  cron_expr = "0 9 * * *"
  timezone  = "UTC"

  session_config = {
    tags = {
      managed_by = "terraform"
    }
    params = {
      message    = "Run the scheduled task."
      agent_type = "claude"
      oneshot    = true
    }
  }
}

resource "ccplant_schedule" "team_example" {
  name      = "terraform-example-team-schedule"
  scope     = "team"
  team_id   = var.team_id
  cron_expr = "0 9 * * *"
  timezone  = "UTC"

  session_config = {
    tags = {
      managed_by = "terraform"
      scope      = "team"
    }
    params = {
      message    = "Run the scheduled team task."
      agent_type = "claude"
      oneshot    = true
    }
  }
}

variable "team_id" {
  type        = string
  description = "agentapi-proxy team ID, for example org/team-slug."
}
