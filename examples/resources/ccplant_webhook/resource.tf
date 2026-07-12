resource "ccplant_webhook" "example" {
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
      auth_proxy = true
    }
  }
}

resource "ccplant_webhook" "team_example" {
  name             = "terraform-example-team-webhook"
  scope            = "team"
  team_id          = var.team_id
  type             = "custom"
  secret           = var.webhook_secret
  signature_header = "X-Signature"
  signature_type   = "static"
  max_sessions     = 1

  triggers = [{
    name          = "team-example-trigger"
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
      scope      = "team"
    }
    params = {
      message    = "Handle team custom webhook payload."
      agent_type = "claude"
      oneshot    = true
      auth_proxy = true
    }
  }
}

variable "webhook_secret" {
  type      = string
  sensitive = true
}

variable "team_id" {
  type        = string
  description = "agentapi-proxy team ID, for example org/team-slug."
}
