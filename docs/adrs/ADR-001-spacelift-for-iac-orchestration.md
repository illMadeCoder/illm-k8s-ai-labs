# ADR-001: Spacelift for IaC Orchestration

## Status

Accepted

## Context

Need a platform to manage Terraform deployments across Azure and AWS with state management, credential handling, stack dependencies, and policy enforcement.

## Decision

**Use Spacelift** for Terraform state management and IaC orchestration.

## Comparison

| Factor | Spacelift | Terraform Cloud | Scalr | env0 | Self-Managed |
|--------|-----------|-----------------|-------|------|--------------|
| **Pricing** | Concurrency-based | Per-resource (RUM) | Per-run (50 free/mo) | Tier + RUM | Free |
| **Stack dependencies** | Native | Limited | Yes | Yes | Manual |
| **Policy-as-code** | OPA (unlimited free) | Sentinel (paid) | OPA | OPA | None |
| **OpenTofu support** | Yes | No | Yes (founding) | Yes | Yes |
| **Drift detection** | Yes | Paid tier | Yes | Yes | Manual |
| **Credential mgmt** | Contexts | Workspaces | Workspaces | Environments | Manual |

## Why Spacelift

- **Stack dependencies** - Foundation stacks → experiment clusters ordering
- **Unlimited OPA policies** on free tier - supports governance learning
- **OpenTofu support** - future-proofs against licensing changes
- **Predictable pricing** - concurrency-based, not per-resource

## Why Not Others

**Terraform Cloud**: Per-resource pricing gets expensive; no OpenTofu; Sentinel requires paid tier

**Scalr**: Close second - simpler pricing, OpenTofu founding member. Spacelift's stack dependencies and unlimited policies won out.

**env0**: Optimized for self-service/team use cases, less relevant for learning project

**Self-Managed**: No drift detection, no policy enforcement, misses enterprise pattern learning

## Consequences

**Positive**: Enterprise IaC patterns, policy foundation, multi-cloud credential isolation

**Negative**: Added complexity, vendor dependency (mitigated by standard Terraform)

## References

- [Spacelift Docs](https://docs.spacelift.io/)
- [Spacelift vs Terraform Cloud](https://spacelift.io/terraform-cloud-vs-spacelift)

## Decision Date

2025-12-10
