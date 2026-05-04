output "cloudfront_domain" {
  description = "CloudFrontディストリビューションのドメイン名"
  value       = aws_cloudfront_distribution.frontend.domain_name
}

output "cloudfront_distribution_id" {
  description = "CloudFrontディストリビューションID（デプロイ時のキャッシュ無効化に使用）"
  value       = aws_cloudfront_distribution.frontend.id
}

output "alb_dns_name" {
  description = "ALBのDNS名"
  value       = aws_lb.api.dns_name
}

output "ecr_repository_url" {
  description = "ECRリポジトリURL（Dockerイメージのpush先）"
  value       = aws_ecr_repository.app.repository_url
}

output "ecs_cluster_name" {
  description = "ECSクラスター名"
  value       = aws_ecs_cluster.main.name
}

output "ecs_service_name" {
  description = "ECSサービス名"
  value       = aws_ecs_service.app.name
}

output "github_actions_role_arn" {
  description = "GitHub Actions用IAMロールARN（ワークフローのrole-to-assumeに設定）"
  value       = aws_iam_role.github_actions.arn
}

output "acm_frontend_validation_records" {
  description = "CloudFront用ACM証明書のDNS検証レコード（Cloudflareに登録が必要）"
  value = {
    for dvo in aws_acm_certificate.frontend.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      record = dvo.resource_record_value
    }
  }
}

output "acm_api_validation_records" {
  description = "ALB用ACM証明書のDNS検証レコード（Cloudflareに登録が必要）"
  value = {
    for dvo in aws_acm_certificate.api.domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      type   = dvo.resource_record_type
      record = dvo.resource_record_value
    }
  }
}
