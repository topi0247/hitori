data "terraform_remote_state" "aws" {
  backend = "local"
  config = {
    path = "${path.module}/../aws/terraform.tfstate"
  }
}

# hitori.topi-log.com → CloudFront
resource "cloudflare_record" "frontend" {
  zone_id = var.zone_id
  name    = var.subdomain
  content = data.terraform_remote_state.aws.outputs.cloudfront_domain
  type    = "CNAME"
  ttl     = 1
  proxied = false
}

# api.hitori.topi-log.com → ALB
resource "cloudflare_record" "api" {
  zone_id = var.zone_id
  name    = "api.${var.subdomain}"
  content = data.terraform_remote_state.aws.outputs.alb_dns_name
  type    = "CNAME"
  ttl     = 1
  proxied = false
}

# CloudFront用ACM証明書のDNS検証レコード
resource "cloudflare_record" "acm_frontend_validation" {
  for_each = data.terraform_remote_state.aws.outputs.acm_frontend_validation_records

  zone_id = var.zone_id
  name    = replace(trimsuffix(each.value.name, "."), ".${var.zone_name}", "")
  content = trimsuffix(each.value.record, ".")
  type    = each.value.type
  ttl     = 60
  proxied = false
}

# ALB用ACM証明書のDNS検証レコード
resource "cloudflare_record" "acm_api_validation" {
  for_each = data.terraform_remote_state.aws.outputs.acm_api_validation_records

  zone_id = var.zone_id
  name    = replace(trimsuffix(each.value.name, "."), ".${var.zone_name}", "")
  content = trimsuffix(each.value.record, ".")
  type    = each.value.type
  ttl     = 60
  proxied = false
}
