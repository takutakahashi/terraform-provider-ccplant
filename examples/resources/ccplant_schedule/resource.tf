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
