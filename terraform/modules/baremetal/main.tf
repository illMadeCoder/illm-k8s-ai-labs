# Bare Metal Kubernetes Cluster Module
# Bootstraps Kubernetes on existing bare metal servers
#
# TODO: Implement using kubeadm or similar
# Options: kubeadm, k3s, rke2

terraform {
  required_version = ">= 1.0.0"
}

# Placeholder - implement kubeadm bootstrap
resource "null_resource" "placeholder" {
  triggers = {
    message = "Baremetal module not yet implemented"
  }
}
