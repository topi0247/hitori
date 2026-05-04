# Terraform リソース対応表

Terraformで定義した各リソースが、AWSコンソール・Cloudflareダッシュボードのどこにあるかをまとめたドキュメント。

---

## AWS

### VPC（`vpc.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_vpc` | VPC → **Your VPCs** |
| `aws_internet_gateway` | VPC → **Internet Gateways** |
| `aws_subnet` | VPC → **Subnets** |
| `aws_route_table` | VPC → **Route Tables** |
| `aws_route_table_association` | Route Tables → 対象テーブル → **Subnet associations** タブ |

> **補足**: サブネットが「パブリック」かどうかは Route Tables の Routes タブで `0.0.0.0/0 → igw-xxx` の経路があるかで判断できる。

---

### ACM 証明書（`acm.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_acm_certificate.frontend` | **バージニア北部（us-east-1）** → Certificate Manager → **Certificates** |
| `aws_acm_certificate.api` | **東京（ap-northeast-1）** → Certificate Manager → **Certificates** |
| `aws_acm_certificate_validation` | 同上 → 対象証明書 → Status が **Issued** になれば検証完了 |

> **注意**: CloudFront用の証明書はリージョンを **us-east-1** に切り替えないとコンソールに表示されない。

---

### S3（`s3.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_s3_bucket` | S3 → **Buckets** → `hitori-frontend-<アカウントID>` |
| `aws_s3_bucket_public_access_block` | バケット → **Permissions** タブ → Block public access |

---

### CloudFront（`cloudfront.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_cloudfront_origin_access_control` | CloudFront → **Origin access** |
| `aws_cloudfront_distribution` | CloudFront → **Distributions** |
| `aws_s3_bucket_policy` | S3 → バケット → **Permissions** タブ → Bucket policy |

> **補足**: DistributionのStatusが **Enabled** かつ Last modified が更新されていればデプロイ完了。

---

### ECR（`ecr.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_ecr_repository` | ECR → **Repositories** → `hitori` |
| `aws_ecr_lifecycle_policy` | リポジトリ → **Lifecycle Policy** タブ |

> **補足**: DockerイメージをpushするとここにImageが一覧表示される。

---

### セキュリティグループ（`security_groups.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_security_group.alb` | EC2 → **Security Groups** → `hitori-alb-sg` |
| `aws_security_group.ecs` | EC2 → **Security Groups** → `hitori-ecs-sg` |

> **補足**: Inbound rules で ALB は 80/443、ECS は 8080（ALBからのみ）を許可しているか確認できる。

---

### ALB（`alb.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_lb` | EC2 → **Load Balancers** → `hitori-alb` |
| `aws_lb_target_group` | EC2 → **Target Groups** → `hitori-tg` |
| `aws_lb_listener.http` | Load Balancers → **Listeners** タブ → HTTP:80 |
| `aws_lb_listener.https` | Load Balancers → **Listeners** タブ → HTTPS:443 |

> **補足**: Target Groups の **Targets** タブでECSタスクが登録されHealthyになっているか確認できる。

---

### IAM（`iam.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_iam_role.ecs_task_execution` | IAM → **Roles** → `hitori-ecs-task-execution` |
| `aws_iam_role.ecs_task` | IAM → **Roles** → `hitori-ecs-task` |
| `aws_iam_role_policy_attachment` | ロール → **Permissions** タブ → Attached policies |
| `aws_iam_role_policy` | ロール → **Permissions** タブ → Inline policies |

> **補足**: `ecs_task_execution` はECSがECRからイメージをpullしたりCloudWatchにログを送るための権限。`ecs_task` はコンテナが実行中にAWSサービスを呼ぶための権限。

---

### Secrets Manager（`secrets.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_secretsmanager_secret.database_url` | Secrets Manager → **Secrets** → `hitori/database_url` |
| `aws_secretsmanager_secret.jwt_secret` | Secrets Manager → **Secrets** → `hitori/jwt_secret` |

> **補足**: シークレットの値は **Retrieve secret value** ボタンで確認できる。

---

### ECS（`ecs.tf`）

| Terraformリソース | AWSコンソールの場所 |
|---|---|
| `aws_cloudwatch_log_group` | CloudWatch → **Log groups** → `/ecs/hitori` |
| `aws_ecs_cluster` | ECS → **Clusters** → `hitori` |
| `aws_ecs_task_definition` | ECS → **Task Definitions** → `hitori` |
| `aws_ecs_service` | ECS → Clusters → `hitori` → **Services** タブ → `hitori` |

> **補足**: Serviceの **Tasks** タブでタスクの起動状況を確認できる。ログは CloudWatch → Log groups → `/ecs/hitori` から確認。

---

## Cloudflare

対象ゾーン: `topi-log.com`

**確認場所**: Cloudflareダッシュボード → `topi-log.com` → **DNS** → **Records**

| Terraformリソース | レコード | 用途 |
|---|---|---|
| `cloudflare_record.frontend` | `hitori.topi-log.com` CNAME → CloudFrontドメイン | フロントエンド配信 |
| `cloudflare_record.api` | `api.hitori.topi-log.com` CNAME → ALB DNSドメイン | APIサーバー |
| `cloudflare_record.acm_frontend_validation` | `_xxx.hitori.topi-log.com` CNAME | CloudFront用ACM証明書のDNS検証 |
| `cloudflare_record.acm_api_validation` | `_xxx.api.hitori.topi-log.com` CNAME | ALB用ACM証明書のDNS検証 |

> **補足**: ACM検証レコードは証明書がIssuedになった後も削除しなくてよい（再発行時にも使われる）。プロキシ（オレンジの雲）はOFFにしてあるため、Cloudflareを経由せず直接AWS側に向く。
