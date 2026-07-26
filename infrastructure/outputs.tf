output "cluster_name" {
  value = module.eks.cluster_name
}

output "aws_region" {
  value = var.aws_region
}

output "go_api_ecr_repository_url" {
  value = aws_ecr_repository.go_api.repository_url
}

output "go_api_db_endpoint" {
  value = aws_db_instance.go_api.endpoint
}

output "go_api_db_name" {
  value = aws_db_instance.go_api.db_name
}

output "go_api_db_master_user_secret_arn" {
  value     = aws_db_instance.go_api.master_user_secret[0].secret_arn
  sensitive = true
}

output "monitoring_namespace" {
  value = helm_release.kube_prometheus_stack.namespace
}

output "grafana_service_name" {
  value = "${helm_release.kube_prometheus_stack.name}-grafana"
}
