terraform {
  required_providers {
    ccplant = {
      source  = "takutakahashi/ccplant"
      version = "~> 0.1"
    }
  }
}

provider "ccplant" {
  endpoint = var.agentapi_proxy_endpoint
  api_key  = var.agentapi_proxy_api_key
}

variable "agentapi_proxy_endpoint" {
  type = string
}

variable "agentapi_proxy_api_key" {
  type      = string
  sensitive = true
}
