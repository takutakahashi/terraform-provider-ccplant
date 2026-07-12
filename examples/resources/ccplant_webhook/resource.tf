resource "ccplant_webhook" "example" {
  body_json = jsonencode({
    name             = "terraform-example-webhook"
    scope            = "user"
    type             = "custom"
    secret           = var.webhook_secret
    signature_header = "X-Signature"
    signature_type   = "static"
    max_sessions     = 1
    triggers = [{
      name          = "example-trigger"
      priority      = 1
      enabled       = true
      stop_on_match = true
      conditions = {
        go_template = "{{ true }}"
      }
    }]
    session_config = {
      tags = {
        managed_by = "terraform"
      }
      params = {
        message    = "Handle custom webhook payload."
        agent_type = "claude"
        oneshot    = true
      }
    }
  })
}

variable "webhook_secret" {
  type      = string
  sensitive = true
}
