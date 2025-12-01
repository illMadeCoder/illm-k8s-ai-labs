# HTTP Baseline Prod - Outputs

# List of cluster names (for iteration)
output "cluster_names" {
  description = "List of cluster names"
  value       = keys(var.clusters)
}

# Cluster details as JSON (for Taskfile parsing)
output "clusters" {
  description = "Map of cluster name to details"
  value = {
    for name, _ in var.clusters : name => {
      name            = module.clusters[name].cluster_name
      resource_group  = module.clusters[name].resource_group_name
      fqdn            = module.clusters[name].cluster_fqdn
      kubeconfig_file = "${path.module}/kubeconfig-${name}"
      kubeconfig_cmd  = "az aks get-credentials --resource-group ${module.clusters[name].resource_group_name} --name ${module.clusters[name].cluster_name} --overwrite-existing"
    }
  }
}
