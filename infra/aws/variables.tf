variable "aws_region" {
  type    = string
  default = "ap-northeast-1"
}

variable "app_name" {
  type    = string
  default = "hitori"
}

variable "domain" {
  type    = string
  default = "hitori.topi-log.com"
}

variable "database_url" {
  type      = string
  sensitive = true
}

variable "jwt_secret" {
  type      = string
  sensitive = true
}

variable "task_cpu" {
  type    = string
  default = "256"
}

variable "task_memory" {
  type    = string
  default = "512"
}

variable "desired_count" {
  type    = number
  default = 1
}

variable "github_repository" {
  type    = string
  default = "topi0247/hitori"
}
