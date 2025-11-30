# Proxmox Module Variables

variable "proxmox_endpoint" {
  description = "Proxmox API endpoint URL"
  type        = string
}

variable "proxmox_api_token" {
  description = "Proxmox API token (user@realm!tokenid=secret)"
  type        = string
  sensitive   = true
}

variable "target_node" {
  description = "Proxmox node to deploy VMs on"
  type        = string
  default     = "pve"
}

variable "node_count" {
  description = "Number of Kubernetes nodes"
  type        = number
  default     = 3
}

variable "node_memory" {
  description = "Memory per node in MB"
  type        = number
  default     = 8192
}

variable "node_cpu" {
  description = "CPU cores per node"
  type        = number
  default     = 4
}

variable "storage_pool" {
  description = "Proxmox storage pool for VM disks"
  type        = string
  default     = "local-lvm"
}

variable "network_bridge" {
  description = "Network bridge for VMs"
  type        = string
  default     = "vmbr0"
}

variable "vm_template" {
  description = "VM template to clone (should have cloud-init)"
  type        = string
  default     = "ubuntu-cloud"
}
