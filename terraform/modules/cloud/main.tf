# Cloud Kubernetes Cluster Module
# Provisions managed Kubernetes on AWS/GCP/Azure
#
# TODO: Implement provider-specific modules
# - AWS EKS
# - GCP GKE
# - Azure AKS

terraform {
  required_version = ">= 1.0.0"
}

# Placeholder - implement cloud provider kubernetes
resource "null_resource" "placeholder" {
  triggers = {
    message = "Cloud module not yet implemented - use submodules for specific providers"
  }
}

# See subdirectories for provider-specific implementations:
# - ./aws/   - Amazon EKS
# - ./gcp/   - Google GKE
# - ./azure/ - Azure AKS
