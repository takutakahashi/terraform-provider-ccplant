resource "ccplant_sandbox_policy" "example" {
  name        = "terraform-example-policy"
  description = "Sandbox policy managed by terraform-provider-ccplant."
  scope       = "user"

  allowed_domains = ["example.com"]
  denied_domains  = ["blocked.example.com"]
  count_mode      = true
}

resource "ccplant_sandbox_policy" "team_example" {
  name        = "terraform-example-team-policy"
  description = "Team-scoped sandbox policy managed by terraform-provider-ccplant."
  scope       = "team"
  team_id     = var.team_id

  allowed_domains = ["example.com"]
  denied_domains  = ["blocked.example.com"]
  count_mode      = true
}

variable "team_id" {
  type        = string
  description = "agentapi-proxy team ID, for example org/team-slug."
}
