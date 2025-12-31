# ADR-004: TLS Certificate Persistence via OpenBao

## Status

Accepted (Let's Encrypt temporarily disabled - see Current Status)

## Context

Let's Encrypt has strict rate limits: 5 certificates per week for the same domain set. In a development environment where clusters are frequently recreated (Kind clusters for testing), each recreation triggers a new certificate request, quickly exhausting rate limits.

Additionally, zero-trust internal communication requires certificates that include both external domains (for browser access) and internal service names (for pod-to-pod TLS).

**Problems to solve:**
1. Avoid Let's Encrypt rate limits on cluster recreation
2. Enable zero-trust internal TLS with proper certificate validation
3. Automate certificate lifecycle without manual intervention

## Decision

**Persist TLS certificates in OpenBao and restore them on cluster recreation using External Secrets Operator (ESO).**

### Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                    CERTIFICATE ISSUANCE / RENEWAL                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│  cert-manager ──────▶ argocd-server-tls-letsencrypt                  │
│  (Let's Encrypt)                 │                                    │
│                                  │ PushSecret (automatic)             │
│                                  ▼                                    │
│                              OpenBao ◀─── ~/.illmlab/openbao-data    │
│                         secret/tls/argocd    (persistent storage)    │
│                                  │                                    │
│                                  │ ExternalSecret (automatic)         │
│                                  ▼                                    │
│                         argocd-server-tls ──▶ ArgoCD                 │
│                                                                       │
├─────────────────────────────────────────────────────────────────────┤
│                      CLUSTER RECREATION                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                       │
│                              OpenBao                                 │
│                    (data persists on host)                           │
│                                  │                                    │
│              ┌───────────────────┼───────────────────┐                │
│              │ ExternalSecret    │ ExternalSecret    │                │
│              ▼                   ▼                   │                │
│  argocd-server-tls-letsencrypt  argocd-server-tls    │                │
│  (cert-manager sees valid       (ArgoCD uses this)   │                │
│   cert → no new request!)                            │                │
│                                                                       │
└─────────────────────────────────────────────────────────────────────┘
```

### Components

| Component | Purpose | File |
|-----------|---------|------|
| **Certificate** | Requests cert from Let's Encrypt | `cert-manager-config/argocd-certificate.yaml` |
| **PushSecret** | Pushes cert to OpenBao on issuance/renewal | `external-secrets-config/argocd-tls-push-secret.yaml` |
| **ExternalSecret (letsencrypt)** | Restores cert-manager's secret from OpenBao | `external-secrets-config/argocd-tls-letsencrypt-external-secret.yaml` |
| **ExternalSecret (argocd)** | Syncs production secret from OpenBao | `external-secrets-config/argocd-tls-external-secret.yaml` |

### Key Implementation Details

1. **Separate secrets**: cert-manager writes to `argocd-server-tls-letsencrypt`, ArgoCD uses `argocd-server-tls`. This decoupling allows ESO to manage the production secret.

2. **cert-manager annotations**: The ExternalSecret that restores `argocd-server-tls-letsencrypt` includes cert-manager annotations so cert-manager recognizes it as a valid, already-issued certificate:
   ```yaml
   template:
     metadata:
       annotations:
         cert-manager.io/issuer-name: letsencrypt-prod
         cert-manager.io/issuer-kind: ClusterIssuer
   ```

3. **Persistent OpenBao storage**: OpenBao data is mounted from `~/.illmlab/openbao-data`, surviving cluster deletion.

4. **Automatic sync**: No manual commands needed. PushSecret and ExternalSecret handle the bidirectional sync automatically.

### Internal TLS (Zero-Trust)

For internal service-to-service communication, OpenBao PKI can issue certificates with all required SANs:

```yaml
dnsNames:
  - argocd.illmlab.xyz              # external domain
  - argocd-server                    # short name
  - argocd-server.argocd             # namespace
  - argocd-server.argocd.svc         # svc
  - argocd-server.argocd.svc.cluster.local  # FQDN
```

This allows internal services (like webhook-relay) to connect via HTTPS using internal service names with proper TLS verification.

## Current Status (2025-12-30)

**Let's Encrypt is temporarily disabled due to rate limit exhaustion.** The rate limit resets approximately January 1, 2025.

### What Changed

The full Let's Encrypt + PushSecret architecture proved complex and ran into issues:

1. **Rate limit exhausted**: Repeated testing during development consumed all 5 certificates
2. **PushSecret circular dependency**: ESO cannot push secrets it owns back to OpenBao (the secret must not be managed by ExternalSecret)
3. **Certificate status mismatch**: Restored certificates showed `False` status in cert-manager

### Files Removed

- `cert-manager-config/argocd-certificate.yaml` - Let's Encrypt Certificate request
- `external-secrets-config/argocd-tls-push-secret.yaml` - PushSecret to OpenBao
- `external-secrets-config/argocd-tls-letsencrypt-external-secret.yaml` - ExternalSecret for cert-manager

### Current Architecture (Simplified)

```
OpenBao PKI
    │
    ├── argocd-tls-external-secret.yaml
    │   (pulls from secret/data/tls/argocd)
    │
    └── argocd-server-tls
        (Kubernetes TLS secret)
```

The existing certificate in OpenBao (persisted from before rate limit exhaustion) is served via ExternalSecret. No new Let's Encrypt requests are made.

### Browser Trust

Currently using OpenBao PKI certificate which shows a browser warning (not publicly trusted). After rate limit reset, the Let's Encrypt architecture can be re-enabled for browser-trusted certificates.

### ArgoCD ignoreDifferences

The ArgoCD application for ArgoCD itself includes `ignoreDifferences` to prevent drift detection on the ESO-managed TLS secret:

```yaml
ignoreDifferences:
  - group: ""
    kind: Secret
    name: argocd-server-tls
    jsonPointers:
      - /data
      - /metadata/labels
      - /metadata/ownerReferences
```

## Consequences

### Positive

- **No rate limits**: Cluster recreation restores existing cert, no new Let's Encrypt request
- **Zero-trust internal TLS**: Certificates include internal service names as SANs
- **Fully automated**: PushSecret/ExternalSecret handle sync without manual intervention
- **Separation of concerns**: cert-manager handles issuance, ESO handles distribution

### Negative

- **Complexity**: Multiple ESO resources to manage the flow
- **Bootstrap chicken-egg**: First cluster creation needs Let's Encrypt (OpenBao empty), subsequent recreations use cached cert
- **cert-manager status**: Certificate resource shows `False` after restore (cosmetic - cert is valid, renewal works)

### Trade-offs

| Approach | Rate Limit Safe | Automation | Complexity |
|----------|-----------------|------------|------------|
| cert-manager only | No | Full | Low |
| Manual backup to OpenBao | Yes | Manual | Medium |
| **PushSecret + ExternalSecret** | Yes | Full | Higher |

## Alternatives Considered

1. **Use OpenBao PKI for everything**: Simpler, but browsers show certificate warnings (not publicly trusted CA).

2. **Manual backup command**: Added `task kind:cert-backup` but requires human intervention.

3. **Use Let's Encrypt staging**: No rate limits but not browser-trusted.

## Files

### Current (Simplified Architecture)

```
hub/app-of-apps/kind/manifests/
├── cert-manager-config/
│   ├── cluster-issuer.yaml          # Let's Encrypt ClusterIssuers (ready for re-enabling)
│   └── openbao-issuer.yaml          # OpenBao PKI ClusterIssuer
└── external-secrets-config/
    ├── argocd-tls-external-secret.yaml  # OpenBao → argocd-server-tls
    └── cluster-secret-store.yaml        # OpenBao connection
```

### Original Architecture (For Re-enabling Later)

When Let's Encrypt rate limits reset, recreate these files to restore full automation:

```
hub/app-of-apps/kind/manifests/
├── cert-manager-config/
│   └── argocd-certificate.yaml      # Let's Encrypt Certificate (DELETED)
└── external-secrets-config/
    ├── argocd-tls-letsencrypt-external-secret.yaml  # OpenBao → cert-manager secret (DELETED)
    └── argocd-tls-push-secret.yaml  # cert-manager → OpenBao (DELETED - but has circular dependency issue)
```

**Note**: The PushSecret approach has a fundamental issue - ESO cannot push secrets it owns. A different pattern may be needed, such as:
- Using a CronJob to periodically backup certs to OpenBao
- Using cert-manager's built-in secret syncing
- Using a dedicated operator for certificate backup

## References

- [cert-manager Certificate resources](https://cert-manager.io/docs/usage/certificate/)
- [ESO PushSecret](https://external-secrets.io/latest/api/pushsecret/)
- [Let's Encrypt Rate Limits](https://letsencrypt.org/docs/rate-limits/)
- [ADR-002: Secrets Management](./ADR-002-secrets-management.md)

## Decision Date

2025-12-30
