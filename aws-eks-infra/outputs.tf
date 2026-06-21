output "cluster_name" {
  value = module.eks.cluster_name
}

output "aws_region" {
  value = var.aws_region
}

output "go_api_ecr_repository_url" {
  value = aws_ecr_repository.go_api.repository_url
}
