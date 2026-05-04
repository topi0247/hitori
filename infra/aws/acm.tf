# CloudFront用（us-east-1 必須）
resource "aws_acm_certificate" "frontend" {
  provider          = aws.us_east_1
  domain_name       = var.domain
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "${var.app_name}-frontend-cert"
  }
}

# ALB用（ap-northeast-1）
resource "aws_acm_certificate" "api" {
  domain_name       = "api.${var.domain}"
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }

  tags = {
    Name = "${var.app_name}-api-cert"
  }
}

# Cloudflare でDNS検証レコードを作成後、証明書がISSUEDになるまで待機
resource "aws_acm_certificate_validation" "frontend" {
  provider        = aws.us_east_1
  certificate_arn = aws_acm_certificate.frontend.arn
  validation_record_fqdns = [
    for dvo in aws_acm_certificate.frontend.domain_validation_options : dvo.resource_record_name
  ]
}

resource "aws_acm_certificate_validation" "api" {
  certificate_arn = aws_acm_certificate.api.arn
  validation_record_fqdns = [
    for dvo in aws_acm_certificate.api.domain_validation_options : dvo.resource_record_name
  ]
}
