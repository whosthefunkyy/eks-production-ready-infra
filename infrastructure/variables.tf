variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "cluster_name" {
  type    = string
  default = "terraform-eks-v3"
}

variable "vpc_name" {
  type    = string
  default = "my-vpc"
}
variable "vpc_cidr" {
  type    = string
  default = "10.0.0.0/16"
}
variable "cluster_version" {
  type    = string
  default = "1.31"
}

variable "github_actions_role_arn" {
  type    = string
  default = "arn:aws:iam::262778473495:role/GitHubActionsEKSRole"
}

variable "ecr_repository_name" {
  type    = string
  default = "go-api-image"
}

variable "local_admin_principal_arn" {
  description = "Optional IAM user or role ARN that should get admin access to the EKS cluster for local kubectl usage."
  type        = string
  default     = ""
}

variable "db_identifier" {
  type    = string
  default = "go-api-postgres"
}

variable "db_name" {
  type    = string
  default = "goapidb"
}

variable "db_username" {
  type    = string
  default = "goapi"
}

variable "db_instance_class" {
  type    = string
  default = "db.t3.micro"
}
