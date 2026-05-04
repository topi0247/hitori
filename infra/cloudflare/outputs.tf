output "frontend_record" {
  description = "フロントエンドのCNAMEレコード"
  value       = cloudflare_record.frontend.hostname
}

output "api_record" {
  description = "APIのCNAMEレコード"
  value       = cloudflare_record.api.hostname
}
