# ADR-016: Hub Metrics Backend — VictoriaMetrics Single

## Status

Accepted (2026-02-08)

**Supersedes:** ADR-009's hub environment decision (Prometheus + Thanos). ADR-009 remains valid for its general TSDB comparison and for tutorials/Kind environments.

## Context

ADR-009 chose Prometheus + Thanos for the hub. In practice, we deployed Mimir distributed as the hub metrics backend because experiment targets remote-write metrics over Tailscale and Mimir's multi-tenant `X-Scope-OrgID` header seemed like a clean isolation model.

Mimir distributed runs 7-8 pods on the hub:

| Component | Mem req | Mem lim |
|-----------|---------|---------|
| ingester | 256Mi | 512Mi |
| distributor | 128Mi | 256Mi |
| querier | 128Mi | 256Mi |
| query_frontend | 128Mi | 256Mi |
| compactor | 128Mi | 256Mi |
| store_gateway | 128Mi | 256Mi |
| gateway | 64Mi | 256Mi |
| **Total** | **960Mi** | **2048Mi** |

On a single-node N100 (16GB) running the full platform stack, this is significant. Multi-tenancy proved unnecessary — we only have one operator deploying experiments, and filtering by experiment name is sufficient.

## Options Considered

### Option 1: Keep Mimir Distributed

**Pros:** Multi-tenant, Grafana-native, already deployed.
**Cons:** 7-8 pods, 960Mi requests / 2GiB limits for a single-user lab. Requires 3 S3 buckets. Multi-tenancy adds header plumbing throughout the stack (Alloy, operator, Grafana datasources) for no real benefit.

### Option 2: Prometheus + Remote Write Receiver

**Pros:** Industry standard, `--web.enable-remote-write-receiver` flag makes it a remote-write target. Huge ecosystem.
**Cons:** No built-in long-term retention without Thanos sidecar (adds complexity back). No built-in downsampling. Single-node scaling limits. Would need a separate Thanos stack for the features we get for free with VM.

### Option 3: VictoriaMetrics Single

**Pros:** Single binary/pod, 256Mi memory, built-in retention, PromQL-compatible, accepts Prometheus remote-write natively at `/api/v1/write`. No additional components needed.
**Cons:** Smaller community. MetricsQL has subtle differences (rarely hit in practice). No multi-tenancy in single-node mode.

### Option 4: Thanos Receive

**Pros:** Accepts remote-write, multi-tenant via headers, stores to S3. Extends existing Prometheus ecosystem.
**Cons:** Still needs multiple components (receive, query, store-gateway, compactor). More complex than VM Single for same result. 3-4 pods minimum.

## Tradeoff Summary

| Concern | Prometheus | VM Single | Mimir | Thanos Receive |
|---------|-----------|-----------|-------|----------------|
| **Pod count** | 1 (+Thanos) | 1 | 7-8 | 3-4 |
| **Memory requests** | ~500Mi (+Thanos) | ~256Mi | ~960Mi | ~500Mi |
| **Remote-write native** | Flag needed | Yes | Yes | Yes |
| **Multi-tenancy** | No | No | Yes | Yes |
| **Long-term retention** | Thanos needed | Built-in | Built-in | Built-in |
| **PromQL compat** | Native | 99%+ | Native | Native |
| **Compression** | Baseline | 5-10x better | ~Baseline | ~Baseline |
| **Operational complexity** | Low (+high w/ Thanos) | Low | High | Medium |
| **Community size** | Largest | Growing | Medium | Large |

## Decision

**Use VictoriaMetrics Single** as the hub metrics backend.

Drop multi-tenancy (`X-Scope-OrgID`) entirely. Use an `experiment` external label on remote-write for per-experiment filtering. PromQL queries use `{experiment="<name>"}`.

### Why

| Factor | Reasoning |
|--------|-----------|
| **Resource efficiency** | 73% less memory requests (960Mi → 256Mi), 86% fewer pods |
| **Simplicity** | Single binary, zero S3 dependency, no distributed coordination |
| **Sufficient for use case** | Lab has one operator, one user — multi-tenancy is overhead |
| **PromQL compatible** | Existing dashboards and operator queries work unchanged |
| **Built-in retention** | `retentionPeriod: 7d` without Thanos/compactor machinery |

### What We Lose

| Lost Capability | Mitigation |
|----------------|-----------|
| Multi-tenancy | `experiment` label filter — same query isolation, simpler plumbing |
| Mimir-specific Grafana integration | VM serves standard Prometheus API; Grafana works natively |
| Horizontal scaling | Not needed — single-node hub, <10k active series |
| S3-backed storage | Local PV (10Gi). Acceptable for 7d retention on a lab |

## Implementation

See commit `feat: Swap Mimir → VictoriaMetrics Single for metrics backhaul`.

Key changes:
- VictoriaMetrics Single deployed via ArgoCD (helm chart `victoria-metrics-single` v0.14.2)
- Tailscale-exposed at `vm-hub` hostname (port 8428)
- `metrics-egress` component replaces `mimir-egress` on target clusters
- Alloy metrics-agent: remote-write to `/api/v1/write`, `external_labels = { experiment = env("EXPERIMENT_NAME") }`
- Operator: `METRICS_URL` env var, queries `/api/v1/query_range` with `{experiment=%q}` filter
- Grafana: single `VictoriaMetrics` datasource replaces `Mimir` + `Mimir (Gateway Lab)`
- Mimir app, values, service, and 3 S3 buckets removed

## Consequences

### Positive

- ~704Mi memory requests freed on hub (~4.3% of 16GB node)
- 6-7 fewer pods reduces scheduling pressure and API server load
- Simpler observability plumbing — no tenant headers anywhere in the stack
- 3 fewer S3 buckets to manage

### Negative

- If a future experiment needs true tenant isolation, we'd need to add VM cluster mode or re-introduce Mimir
- MetricsQL differences could surface with advanced PromQL (not yet encountered)
- Local PV means metrics are lost if the node's disk fails (acceptable for lab)

### Mimir Retained for Tutorials

`components/observability/mimir/component.yaml` is kept for the tsdb-comparison experiment. Mimir can still be deployed to target clusters for side-by-side evaluation — it's just no longer the hub's metrics backend.

## References

- [ADR-009: TSDB Selection](ADR-009-tsdb-selection.md) — Original comparison
- [ADR-011: Observability Architecture](ADR-011-observability-architecture.md) — Stack overview
- [VictoriaMetrics Single Docs](https://docs.victoriametrics.com/single-server-victoriametrics/)
- [VictoriaMetrics Helm Chart](https://github.com/VictoriaMetrics/helm-charts/tree/master/charts/victoria-metrics-single)
