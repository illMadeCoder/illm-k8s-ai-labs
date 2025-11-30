# Cloud Module Outputs

output "kubeconfig_path" {
  description = "Path to generated kubeconfig"
  value       = "" # TODO: implement per provider
}

output "cluster_endpoint" {
  description = "Kubernetes API endpoint"
  value       = "" # TODO: implement per provider
}

output "cluster_ca_certificate" {
  description = "Cluster CA certificate (base64)"
  value       = "" # TODO: implement per provider
  sensitive   = true
}

output "provider" {
  description = "Cloud provider used"
  value       = var.provider
}

output "region" {
  description = "Region deployed in"
  value       = var.region
}
