# HTTP Baseline Experiment - Production Environment
# Deploys multiple AKS clusters for load testing

terraform {
  required_version = ">= 1.0.0"

  # Uncomment to use remote state
  # backend "azurerm" {
  #   resource_group_name  = "tfstate"
  #   storage_account_name = "tfstate"
  #   container_name       = "tfstate"
  #   key                  = "http-baseline/prod.tfstate"
  # }
}

# Clusters to create - names match argocd/*.yaml files
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
    loadgen = {
      vm_size    = "Standard_D2s_v3"
      node_count = 2
      min_nodes  = 1
      max_nodes  = 5
    }
  }
}

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
    environment = "prod"
    experiment  = var.experiment_name
    cluster     = each.key
    managed_by  = "terraform"
  }
}

# Write kubeconfigs to files
resource "local_file" "kubeconfigs" {
  for_each        = var.clusters
  content         = module.clusters[each.key].kube_config
  filename        = "${path.module}/kubeconfig-${each.key}"
  file_permission = "0600"
}
