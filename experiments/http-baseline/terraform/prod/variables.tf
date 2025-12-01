# HTTP Baseline Prod - Variables

variable "experiment_name" {
  description = "Name of the experiment"
  type        = string
  default     = "http-baseline"
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
