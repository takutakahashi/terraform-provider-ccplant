resource "ccplant_slackbot" "example" {
  name                      = "terraform-example-slackbot"
  scope                     = "user"
  allowed_event_types       = ["message", "app_mention"]
  allowed_channel_names     = ["engineering"]
  max_sessions              = 5
  notify_on_session_created = true
  allow_bot_messages        = false

  session_config = {
    initial_message_template = "Handle Slack event: {{ .event.text }}"
    tags = {
      managed_by = "terraform"
    }
    params = {
      agent_type = "claude"
      oneshot    = false
      auth_proxy = true
    }
  }
}

resource "ccplant_slackbot" "team_example" {
  name    = "terraform-example-team-slackbot"
  scope   = "team"
  team_id = var.team_id

  teams                     = [var.team_id]
  allowed_event_types       = ["message", "app_mention"]
  allowed_channel_names     = ["engineering"]
  max_sessions              = 5
  notify_on_session_created = true
  allow_bot_messages        = false

  session_config = {
    initial_message_template = "Handle team Slack event: {{ .event.text }}"
    tags = {
      managed_by = "terraform"
      scope      = "team"
    }
    params = {
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
