#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

usage() {
  echo "Usage: $0 [component...]"
  echo ""
  echo "Components:"
  echo "  argocd        Install ArgoCD"
  echo "  vault         Install HashiCorp Vault"
  echo "  gateway-api   Install Gateway API with Envoy Gateway"
  echo "  backstage     Install Backstage developer portal"
  echo "  observability Install Prometheus & Grafana"
  echo "  cert-manager  Install cert-manager for TLS"
  echo "  all           Install everything (default)"
  echo ""
  echo "Examples:"
  echo "  $0                  # Install all"
  echo "  $0 argocd             # Install only ArgoCD"
  echo "  $0 vault              # Install only Vault"
  echo "  $0 gateway-api        # Install only Gateway API"
  echo "  $0 backstage          # Install only Backstage"
  echo "  $0 observability      # Install only Prometheus & Grafana"
  echo "  $0 cert-manager       # Install only cert-manager"
  echo "  $0 argocd vault       # Install multiple components"
}

install_argocd() {
  echo "========================================"
  echo "Installing ArgoCD..."
  echo "========================================"

  helm repo add argo https://argoproj.github.io/argo-helm 2>/dev/null || true
  helm repo update

  helm install argocd argo/argo-cd \
    --namespace argocd \
    --create-namespace \
    --values "$SCRIPT_DIR/argocd/values.yaml" \
    --wait

  echo ""
  echo "ArgoCD installed!"
  echo "Password: $(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d)"
  echo ""
  echo "To access:"
  echo "  kubectl port-forward svc/argocd-server -n argocd 8080:443"
  echo "  https://localhost:8080 (admin)"
  echo ""
}

install_vault() {
  echo "========================================"
  echo "Installing Vault..."
  echo "========================================"
  "$SCRIPT_DIR/vault/install.sh"
}

install_gateway_api() {
  echo "========================================"
  echo "Installing Gateway API with Envoy Gateway..."
  echo "========================================"
  "$SCRIPT_DIR/gateway-api/install.sh"
}

install_backstage() {
  echo "========================================"
  echo "Installing Backstage..."
  echo "========================================"
  "$SCRIPT_DIR/backstage/install.sh"
}

install_observability() {
  echo "========================================"
  echo "Installing Observability Stack..."
  echo "========================================"
  "$SCRIPT_DIR/observability/install.sh"
}

install_cert_manager() {
  echo "========================================"
  echo "Installing cert-manager..."
  echo "========================================"
  "$SCRIPT_DIR/cert-manager/install.sh"
}

# Parse arguments
COMPONENTS=("$@")
if [ ${#COMPONENTS[@]} -eq 0 ] || [[ " ${COMPONENTS[*]} " =~ " all " ]]; then
  COMPONENTS=("argocd" "vault" "gateway-api" "backstage" "observability" "cert-manager")
fi

# Install requested components
for component in "${COMPONENTS[@]}"; do
  case $component in
    argocd)
      install_argocd
      ;;
    vault)
      install_vault
      ;;
    gateway-api)
      install_gateway_api
      ;;
    backstage)
      install_backstage
      ;;
    observability)
      install_observability
      ;;
    cert-manager)
      install_cert_manager
      ;;
    all)
      # Handled above
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Unknown component: $component"
      usage
      exit 1
      ;;
  esac
done

echo ""
echo "========================================"
echo "Bootstrap complete!"
echo "========================================"
