resource "helm_release" "kube_prometheus_stack" {
  name             = "kube-prometheus-stack"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  namespace        = "monitoring"
  create_namespace = true
  timeout          = 900
  wait             = true

  values = [
    file("${path.module}/kube-prometheus-stack-values/values.yaml")
  ]

  depends_on = [
    module.eks,
    helm_release.aws_load_balancer_controller,
  ]
}
