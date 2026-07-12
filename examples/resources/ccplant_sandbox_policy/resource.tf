resource "ccplant_sandbox_policy" "example" {
  name        = "terraform-example-policy"
  description = "Sandbox policy managed by terraform-provider-ccplant."
  scope       = "user"

  allowed_domains = ["example.com"]
  denied_domains  = ["blocked.example.com"]
  count_mode      = true
}
