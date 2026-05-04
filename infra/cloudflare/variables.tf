variable "cloudflare_api_token" {
  type      = string
  sensitive = true
}

variable "zone_id" {
  type    = string
  default = "6bcdf9e86929897a7682c0b2952c642f"
}

variable "zone_name" {
  type    = string
  default = "topi-log.com"
}

variable "subdomain" {
  type    = string
  default = "hitori"
}
