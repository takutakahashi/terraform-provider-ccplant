resource "ccplant_schedule" "example" {
  body_json = jsonencode({
    name     = "terraform-example"
    status   = "active"
    scope    = "user"
    cron     = "0 9 * * *"
    timezone = "UTC"
  })
}
