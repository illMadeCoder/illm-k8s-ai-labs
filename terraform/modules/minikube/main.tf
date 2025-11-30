# Minikube Cluster Module
# Manages a local minikube cluster for development/testing

terraform {
  required_version = ">= 1.0.0"
}

# Minikube cluster management via local-exec
# Note: Requires minikube CLI installed locally

resource "null_resource" "minikube_cluster" {
  triggers = {
    cluster_name = var.cluster_name
    cpus         = var.cpus
    memory       = var.memory
    driver       = var.driver
  }

  provisioner "local-exec" {
    command = <<-EOT
      minikube start \
        --profile=${var.cluster_name} \
        --cpus=${var.cpus} \
        --memory=${var.memory} \
        --driver=${var.driver} \
        --kubernetes-version=${var.kubernetes_version} \
        --addons=metrics-server,ingress
    EOT
  }

  provisioner "local-exec" {
    when    = destroy
    command = "minikube delete --profile=${self.triggers.cluster_name}"
  }
}

# Generate kubeconfig
resource "local_file" "kubeconfig" {
  depends_on = [null_resource.minikube_cluster]

  filename = "${path.module}/kubeconfig"
  content  = <<-EOT
    # Generated kubeconfig for minikube cluster: ${var.cluster_name}
    # Use: export KUBECONFIG=${path.module}/kubeconfig
  EOT

  provisioner "local-exec" {
    command = "minikube --profile=${var.cluster_name} kubectl config view --flatten > ${path.module}/kubeconfig"
  }
}
