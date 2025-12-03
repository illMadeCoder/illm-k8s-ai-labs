# Multi-Cloud Demo Experiment - Azure Environment
# Deploys AKS cluster with Crossplane Workload Identity enabled

terraform {
  required_version = ">= 1.0.0"

  # Uncomment to use remote state via Spacelift
  # backend "remote" {
  #   organization = "illm-k8s-lab"
  # }
}

variable "experiment_name" {
  description = "Name of this experiment"
  type        = string
  default     = "multi-cloud-demo"
}

variable "location" {
  description = "Azure region"
  type        = string
  default     = "eastus"
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.29"
}

# Cluster configuration
variable "clusters" {
  description = "Map of cluster configs"
  type = map(object({
    vm_size    = string
    node_count = number
    min_nodes  = number
    max_nodes  = number
  }))
  default = {
    target = {
      vm_size    = "Standard_D4s_v3"
      node_count = 3
      min_nodes  = 2
      max_nodes  = 10
    }
  }
}

# Create AKS clusters
module "clusters" {
  source   = "../../../../terraform-modules/aks"
  for_each = var.clusters

  cluster_name        = "${var.experiment_name}-${each.key}"
  resource_group_name = "${var.experiment_name}-${each.key}-rg"
  location            = var.location
  kubernetes_version  = var.kubernetes_version

  node_count          = each.value.node_count
  vm_size             = each.value.vm_size
  enable_auto_scaling = true
  min_nodes           = each.value.min_nodes
  max_nodes           = each.value.max_nodes

  enable_monitoring = true
  create_acr        = false

  tags = {
    environment = "demo"
    experiment  = var.experiment_name
    cluster     = each.key
    managed_by  = "spacelift"
    cloud       = "azure"
  }
}

# Write kubeconfigs to files
resource "local_file" "kubeconfigs" {
  for_each        = var.clusters
  content         = module.clusters[each.key].kube_config
  filename        = "${path.module}/kubeconfig-${each.key}"
  file_permission = "0600"
}

# Outputs
output "cluster_names" {
  description = "List of cluster names"
  value       = [for name, _ in var.clusters : name]
}

output "cluster_endpoints" {
  description = "API server endpoints"
  value = {
    for name, cluster in module.clusters :
    name => cluster.cluster_fqdn
  }
}

output "kubeconfig_paths" {
  description = "Paths to kubeconfig files"
  value = {
    for name, _ in var.clusters :
    name => "${path.module}/kubeconfig-${name}"
  }
}
