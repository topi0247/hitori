# インフラ仕様

## 技術スタック

| カテゴリ | 技術 |
|---|---|
| コンテナ | Docker |
| CI/CD | GitHub Actions |
| ドメイン・DNS | Cloudflare |
| CDN・フロント配信 | CloudFront + S3 |
| APIサーバー | ECS Fargate + ECR |
| ロードバランサー | ALB |
| ネットワーク | VPC（パブリックサブネットのみ） |
| SSL証明書 | ACM |
| IaC | Terraform |

## 構成図

```
ユーザー
  ↓
CloudFront → S3（React）
  ↓
ALB
  ↓
ECS Fargate（Go）
  ↓ インターネット経由
Supabase（Auth + PostgreSQL）
```

## DNS・ドメイン

- 独自ドメインはCloudflareで取得・DNS管理
- サブドメインをCNAMEでCloudFront / ALBに向ける
- CloudflareのプロキシはOFF（グレーの雲マーク）

## Terraform構成

```
infra/
├── aws/          # S3, CloudFront, ECS, ALB, VPC, ACM
└── cloudflare/   # DNS レコード
```
