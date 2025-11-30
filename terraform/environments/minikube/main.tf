# Minikube Environment Configuration
# Uses the minikube module with environment-specific settings

terraform {
  required_version = ">= 1.0.0"
}

module "cluster" {
  source = "../../modules/minikube"

  cluster_name       = var.cluster_name
  cpus               = var.cpus
  memory             = var.memory
  driver             = var.driver
  kubernetes_version = var.kubernetes_version
}

# Variables with defaults for this environment
variable "cluster_name" {
  default = "illm-lab"
}

variable "cpus" {
  default = 4
}

variable "memory" {
  default = 8192
}

variable "driver" {
  default = "docker"
}

variable "kubernetes_version" {
  default = "v1.29.0"
}

# Outputs
output "cluster_name" {
  value = module.cluster.cluster_name
}

output "kubeconfig_path" {
  value = module.cluster.kubeconfig_path
}
