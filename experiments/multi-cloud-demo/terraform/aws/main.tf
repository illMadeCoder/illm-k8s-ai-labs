# Multi-Cloud Demo Experiment - AWS Environment
# Deploys EKS cluster with Crossplane IRSA enabled

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

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
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
    instance_type = string
    node_count    = number
    min_nodes     = number
    max_nodes     = number
  }))
  default = {
    target = {
      instance_type = "t3.large"
      node_count    = 3
      min_nodes     = 2
      max_nodes     = 10
    }
  }
}

# Create EKS clusters
module "clusters" {
  source   = "../../../../terraform-modules/eks"
  for_each = var.clusters

  cluster_name       = "${var.experiment_name}-${each.key}"
  region             = var.region
  kubernetes_version = var.kubernetes_version

  node_count    = each.value.node_count
  instance_type = each.value.instance_type
  min_nodes     = each.value.min_nodes
  max_nodes     = each.value.max_nodes

  enable_monitoring      = true
  enable_crossplane_irsa = true
  create_ecr             = false

  tags = {
    environment = "demo"
    experiment  = var.experiment_name
    cluster     = each.key
    managed_by  = "spacelift"
    cloud       = "aws"
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
    name => cluster.cluster_endpoint
  }
}

output "kubeconfig_paths" {
  description = "Paths to kubeconfig files"
  value = {
    for name, _ in var.clusters :
    name => "${path.module}/kubeconfig-${name}"
  }
}

output "crossplane_role_arns" {
  description = "IRSA role ARNs for Crossplane (if enabled)"
  value = {
    for name, cluster in module.clusters :
    name => cluster.crossplane_role_arn
  }
}
