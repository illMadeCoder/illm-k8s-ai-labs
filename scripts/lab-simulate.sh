#!/bin/bash
# lab-simulate.sh — Apply all experiment CRs and monitor their progress
#
# Usage:
#   ./scripts/lab-simulate.sh              # Apply all experiments
#   ./scripts/lab-simulate.sh --status     # Just check status
#   ./scripts/lab-simulate.sh --dry-run    # Validate without applying

set -euo pipefail

EXPERIMENTS_DIR="$(cd "$(dirname "$0")/../experiments" && pwd)"
NAMESPACE="experiments"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

list_experiments() {
  for dir in "$EXPERIMENTS_DIR"/*/; do
    if [ -f "$dir/experiment.yaml" ]; then
      basename "$dir"
    fi
  done
}

apply_all() {
  echo "Applying all experiments to namespace $NAMESPACE..."
  for name in $(list_experiments); do
    echo -e "  ${GREEN}→${NC} $name"
    kubectl apply -f "$EXPERIMENTS_DIR/$name/experiment.yaml"
  done
  echo ""
  echo "$(list_experiments | wc -l) experiments applied."
  echo "Monitor with: $0 --status"
}

check_status() {
  echo "Experiment Status:"
  echo "─────────────────────────────────────────────────"
  printf "%-30s %-15s %-10s\n" "EXPERIMENT" "PHASE" "CLEANED"
  echo "─────────────────────────────────────────────────"

  for name in $(list_experiments); do
    phase=$(kubectl get experiment "$name" -n "$NAMESPACE" -o jsonpath='{.status.phase}' 2>/dev/null || echo "NotFound")
    cleaned=$(kubectl get experiment "$name" -n "$NAMESPACE" -o jsonpath='{.status.resourcesCleaned}' 2>/dev/null || echo "-")

    case "$phase" in
      Complete) color=$GREEN ;;
      Running|Ready) color=$YELLOW ;;
      Failed) color=$RED ;;
      *) color=$NC ;;
    esac

    printf "%-30s ${color}%-15s${NC} %-10s\n" "$name" "$phase" "$cleaned"
  done
}

dry_run() {
  echo "Dry-run validation:"
  echo ""
  for name in $(list_experiments); do
    file="$EXPERIMENTS_DIR/$name/experiment.yaml"
    if kubectl apply --dry-run=client -f "$file" > /dev/null 2>&1; then
      echo -e "  ${GREEN}✓${NC} $name"
    else
      echo -e "  ${RED}✗${NC} $name"
      kubectl apply --dry-run=client -f "$file" 2>&1 | sed 's/^/    /'
    fi
  done
}

case "${1:-apply}" in
  --status) check_status ;;
  --dry-run) dry_run ;;
  *) apply_all ;;
esac
