# Minikube Module Variables

variable "cluster_name" {
  description = "Name of the minikube cluster/profile"
  type        = string
  default     = "illm-lab"
}

variable "cpus" {
  description = "Number of CPUs to allocate"
  type        = number
  default     = 4
}

variable "memory" {
  description = "Memory in MB to allocate"
  type        = number
  default     = 8192
}

variable "driver" {
  description = "Minikube driver (docker, hyperkit, virtualbox, etc.)"
  type        = string
  default     = "docker"
}

variable "kubernetes_version" {
  description = "Kubernetes version to deploy"
  type        = string
  default     = "v1.29.0"
}
