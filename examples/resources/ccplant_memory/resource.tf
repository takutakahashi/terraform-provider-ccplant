resource "ccplant_memory" "example" {
  title   = "Terraform managed memory"
  content = "This memory entry is managed by terraform-provider-ccplant."
  scope   = "user"

  tags = {
    managed_by = "terraform"
  }
}
