resource "ccplant_memory" "example" {
  title   = "Terraform managed memory"
  content = "This memory entry is managed by terraform-provider-ccplant."
  scope   = "user"

  tags = {
    managed_by = "terraform"
  }
}

resource "ccplant_memory" "team_example" {
  title   = "Terraform managed team memory"
  content = "This team-scoped memory entry is managed by terraform-provider-ccplant."
  scope   = "team"
  team_id = var.team_id

  tags = {
    managed_by = "terraform"
    scope      = "team"
  }
}

variable "team_id" {
  type        = string
  description = "agentapi-proxy team ID, for example org/team-slug."
}
