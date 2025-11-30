# Proxmox Kubernetes Cluster Module
# Provisions VMs on Proxmox VE and bootstraps Kubernetes
#
# TODO: Implement using bpg/proxmox provider
# https://registry.terraform.io/providers/bpg/proxmox/latest

terraform {
  required_version = ">= 1.0.0"

  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = ">= 0.38.0"
    }
  }
}

# Placeholder - implement VM creation and k8s bootstrap
resource "null_resource" "placeholder" {
  triggers = {
    message = "Proxmox module not yet implemented"
  }
}
