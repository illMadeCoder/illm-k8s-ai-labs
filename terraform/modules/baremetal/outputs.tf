# Bare Metal Module Outputs

output "kubeconfig_path" {
  description = "Path to generated kubeconfig"
  value       = "" # TODO: implement
}

output "cluster_endpoint" {
  description = "Kubernetes API endpoint"
  value       = var.control_plane_endpoint
}

output "control_plane_ips" {
  description = "Control plane node IPs"
  value       = var.control_plane_nodes
}

output "worker_ips" {
  description = "Worker node IPs"
  value       = var.worker_nodes
}
