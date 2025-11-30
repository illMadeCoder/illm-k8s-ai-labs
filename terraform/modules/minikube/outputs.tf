# Minikube Module Outputs

output "cluster_name" {
  description = "Name of the minikube cluster"
  value       = var.cluster_name
}

output "kubeconfig_path" {
  description = "Path to generated kubeconfig file"
  value       = local_file.kubeconfig.filename
}

output "cluster_endpoint" {
  description = "Kubernetes API endpoint"
  value       = "https://$(minikube ip --profile=${var.cluster_name}):8443"
}
