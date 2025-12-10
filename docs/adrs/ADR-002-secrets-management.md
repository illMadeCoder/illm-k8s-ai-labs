# ADR-002: Secrets Management with ESO + Vault

## Status

Accepted

## Context

Need GitOps-friendly secrets management that solves the bootstrap problem: cloud credentials are needed to deploy Vault, but Vault is where we'd store credentials.

## Decision

**Use External Secrets Operator (ESO)** with a two-phase approach:

1. **Bootstrap**: Cloud secret managers (Azure KV, AWS SM) → ESO → K8s Secrets
2. **Production**: Vault → ESO → K8s Secrets (cloud managers retained for ESO auth)

## Comparison

| Factor | ESO | Vault Agent | Vault CSI | Sealed Secrets | SOPS |
|--------|-----|-------------|-----------|----------------|------|
| **GitOps-friendly** | Yes (CRs) | No (annotations) | Partial | Yes | Yes |
| **Central management** | Yes | Yes | Yes | No | No |
| **Multiple backends** | Vault, Azure, AWS, GCP | Vault only | Vault only | N/A | KMS only |
| **Rotation** | Yes (refreshInterval) | Yes | Limited | No | No |
| **Per-pod config needed** | No | Yes (sidecar) | Yes (CSI volume) | No | No |
| **Dynamic credentials** | Via Vault backend | Yes | Yes | No | No |

## Why ESO + Two-Phase

- **Solves bootstrap**: Cloud secret managers don't need K8s to exist; ESO uses Workload Identity/IRSA
- **GitOps-native**: ExternalSecret CRs in Git, actual secrets never stored
- **Backend flexibility**: Switch from cloud → Vault without app changes
- **No sidecars**: Works with any workload, no per-pod configuration

## Why Not Others

**Sealed Secrets**: No central management, no rotation, secrets still version-controlled in Git

**SOPS + Age/KMS**: Complex key management, no dynamic credentials, requires CI/CD tooling

**Vault Agent Only**: Per-pod annotations required, not GitOps-friendly, sidecar complexity

**Vault CSI Only**: Requires CSI driver, volume management per pod, limited rotation

## Consequences

**Positive**: GitOps-compatible, automatic rotation, centralized audit (Vault), no static credentials

**Negative**: ESO adds operational overhead, two systems during transition

## References

- [External Secrets Operator](https://external-secrets.io/)
- [ESO Vault Provider](https://external-secrets.io/latest/provider/hashicorp-vault/)

## Decision Date

2025-12-10
