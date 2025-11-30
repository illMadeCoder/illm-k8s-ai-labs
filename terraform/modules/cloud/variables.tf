# Cloud Module Variables

variable "provider" {
  description = "Cloud provider (aws, gcp, azure)"
  type        = string
  default     = "aws"

  validation {
    condition     = contains(["aws", "gcp", "azure"], var.provider)
    error_message = "Provider must be one of: aws, gcp, azure"
  }
}

variable "region" {
  description = "Cloud region to deploy in"
  type        = string
  default     = "us-west-2"
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "illm-lab"
}

variable "kubernetes_version" {
  description = "Kubernetes version"
  type        = string
  default     = "1.29"
}

variable "node_count" {
  description = "Number of worker nodes"
  type        = number
  default     = 3
}

variable "node_type" {
  description = "Instance type for nodes"
  type        = string
  default     = "t3.large" # AWS default, will vary by provider
}

variable "min_nodes" {
  description = "Minimum nodes for autoscaling"
  type        = number
  default     = 1
}

variable "max_nodes" {
  description = "Maximum nodes for autoscaling"
  type        = number
  default     = 10
}

# Provider-specific credentials (use environment variables or secrets manager)
variable "aws_region" {
  description = "AWS region (if using AWS)"
  type        = string
  default     = ""
}

variable "gcp_project" {
  description = "GCP project ID (if using GCP)"
  type        = string
  default     = ""
}

variable "azure_subscription_id" {
  description = "Azure subscription ID (if using Azure)"
  type        = string
  default     = ""
}
