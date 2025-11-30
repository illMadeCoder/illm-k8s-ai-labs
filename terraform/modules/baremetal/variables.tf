# Bare Metal Module Variables

variable "control_plane_endpoint" {
  description = "Kubernetes API endpoint (load balancer or VIP)"
  type        = string
}

variable "control_plane_nodes" {
  description = "List of control plane node IPs"
  type        = list(string)
  default     = []
}

variable "worker_nodes" {
  description = "List of worker node IPs"
  type        = list(string)
  default     = []
}

variable "ssh_user" {
  description = "SSH user for provisioning"
  type        = string
  default     = "ubuntu"
}

variable "ssh_private_key_path" {
  description = "Path to SSH private key"
  type        = string
  sensitive   = true
}

variable "pod_cidr" {
  description = "Pod network CIDR"
  type        = string
  default     = "10.244.0.0/16"
}

variable "service_cidr" {
  description = "Service network CIDR"
  type        = string
  default     = "10.96.0.0/12"
}

variable "kubernetes_version" {
  description = "Kubernetes version to install"
  type        = string
  default     = "1.29.0"
}
