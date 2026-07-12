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
