resource "ccplant_memory" "example" {
  body_json = jsonencode({
    title   = "Terraform managed memory"
    content = "This memory entry is managed by terraform-provider-ccplant."
    scope   = "user"
    tags = {
      managed_by = "terraform"
    }
  })
}
