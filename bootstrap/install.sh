#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

usage() {
  echo "Usage: $0 [--help]"
  echo ""
  echo "Bootstrap the Kubernetes lab environment with ArgoCD."
  echo ""
  echo "This script installs ArgoCD with app-of-apps bootstrapping enabled."
  echo "The app-of-apps Application is created as part of the Helm install via values,"
  echo "automatically deploying all workloads managed in manifests/applications/."
  echo ""
  echo "Components automatically deployed:"
  echo "  - cert-manager (TLS certificate management)"
  echo "  - Gateway API with Envoy Gateway (ingress controller)"
  echo "  - Prometheus & Grafana (observability)"
  echo "  - Backstage (developer portal)"
  echo "  - HashiCorp Vault (secrets management)"
  echo "  - Demo application"
  echo "  - HTTPRoutes (external access configuration)"
  echo ""
  echo "For more information, see:"
  echo "  - bootstrap/argocd/values.yaml (ArgoCD config + app-of-apps bootstrap)"
  echo "  - manifests/applications/ (all Application manifests)"
  echo "  - manifests/workloads/ (service configurations)"
}

# Parse arguments
case "${1:-}" in
  -h|--help)
    usage
    exit 0
    ;;
  "")
    # No arguments - proceed with install
    ;;
  *)
    echo "Unknown option: $1"
    usage
    exit 1
    ;;
esac

echo "========================================"
echo "Installing ArgoCD with app-of-apps..."
echo "========================================"
echo ""

helm repo add argo https://argoproj.github.io/argo-helm 2>/dev/null || true
helm repo update

helm install argocd argo/argo-cd \
  --namespace argocd \
  --create-namespace \
  --values "$SCRIPT_DIR/argocd/values.yaml" \
  --wait

echo ""
echo "========================================"
echo "Bootstrap complete!"
echo "========================================"
echo ""
echo "ArgoCD is installed and the app-of-apps Application is being deployed."
echo "All workloads will sync automatically."
echo ""
echo "Credentials:"
ARGOCD_PASS=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d)
echo "  admin / $ARGOCD_PASS"
echo ""
echo "Access ArgoCD:"
echo "  kubectl port-forward svc/argocd-server -n argocd 8080:443"
echo "  https://localhost:8080"
echo ""
echo "Monitor applications:"
echo "  kubectl get applications -n argocd"
echo "  kubectl get applications -n argocd -o wide"
echo ""
