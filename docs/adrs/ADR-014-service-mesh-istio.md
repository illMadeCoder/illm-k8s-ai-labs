# ADR-014: Hub Service Mesh with Istio

## Status

Accepted

## Context

The hub cluster needed service mesh capabilities for:
- **mTLS** - Encrypted service-to-service communication
- **Ingress** - Path-based routing to internal services
- **Observability** - Service graph visualization, distributed tracing integration
- **Traffic management** - Future canary deployments, fault injection

Previously using Traefik for ingress only (no mesh features). Needed to evaluate mesh options and implement.

### Requirements

| Requirement | Priority | Notes |
|-------------|----------|-------|
| mTLS between services | Must | Zero-trust networking |
| Path-based ingress routing | Must | /grafana, /mimir, /loki, /tempo, /argocd, /openbao, /kiali |
| Talos Linux compatibility | Must | No iptables in containers |
| Kiali visualization | High | Service graph, traffic flow |
| Tempo tracing integration | High | Distributed tracing |
| Resource efficient | High | Single-node N100 homelab |

### Environment Constraints

- **Talos Linux**: Immutable OS, no iptables kernel modules available in containers
- **Single node**: Limited resources (~16GB RAM shared across all workloads)
- **Tailscale ingress**: External traffic via Tailscale LoadBalancer, not traditional ingress

## Options Considered

### Option 1: Istio (Sidecar Model)

**Architecture:**
- Envoy sidecar proxy per pod
- istiod control plane
- CNI plugin for Talos compatibility

**Pros:**
- Most feature-complete mesh
- Excellent Kiali integration
- Strong Tempo/tracing support
- Large community, extensive docs

**Cons:**
- Higher resource overhead (sidecar per pod)
- Complexity

### Option 2: Linkerd

**Architecture:**
- Lightweight Rust-based sidecar (linkerd2-proxy)
- Simpler control plane

**Pros:**
- Lower resource usage than Istio
- Simpler operations
- Faster startup

**Cons:**
- Fewer features than Istio
- Less ecosystem integration

### Option 3: Cilium Service Mesh

**Architecture:**
- eBPF-based, no sidecars
- Integrated with Cilium CNI

**Pros:**
- No sidecar overhead
- Native eBPF performance

**Cons:**
- Requires Cilium as CNI (we use Flannel on Talos)
- L7 features still maturing

## Decision

**Use Istio with CNI plugin** for service mesh.

### Why Istio

| Factor | Reasoning |
|--------|-----------|
| **Feature completeness** | Full traffic management, security, observability |
| **Kiali** | Best-in-class service graph visualization |
| **Tracing** | Native Tempo integration via Zipkin protocol |
| **Learning value** | Industry standard, highly resume-relevant |
| **CNI mode** | Solves Talos iptables limitation |

### Why CNI Plugin (Critical for Talos)

Standard Istio uses an init container to configure iptables rules. This fails on Talos because:
- Talos is immutable - no iptables kernel modules in container context
- Init containers can't modify network namespace via iptables

**Solution:** Istio CNI plugin runs as DaemonSet on nodes, configures iptables rules at node level before pods start.

Key configuration:
```yaml
# istio-istiod values
pilot:
  cni:
    enabled: true  # Critical - tells sidecar injector to skip iptables init

global:
  istio_cni:
    enabled: true
    chained: true  # Chain with existing CNI (Flannel)
```

## Implementation Details

### Component Architecture

```
Tailscale (hub.tailbdf608.ts.net)
         |
         v
Istio Ingress Gateway (ClusterIP + Tailscale LB)
         |
    +----+----+----+----+----+----+
    v    v    v    v    v    v    v
/grafana /mimir /loki /tempo /argocd /openbao /kiali
    |
    v
Sidecars (envoy) <--- mTLS ---> between services
```

### ArgoCD Apps (Sync Waves)

| App | Sync Wave | Purpose |
|-----|-----------|---------|
| `istio-base` | 2 | CRDs and cluster resources |
| `istio-cni` | 2 | CNI plugin DaemonSet |
| `istio-istiod` | 3 | Control plane |
| `istio-ingress` | 4 | Ingress gateway |
| `istio-config` | 5 | Gateway, VirtualServices, policies |
| `kiali` | 5 | Service graph dashboard |

### Namespace Configuration

Namespaces must be labeled for sidecar injection:
```yaml
metadata:
  labels:
    istio-injection: enabled
    pod-security.kubernetes.io/enforce: privileged  # CNI requires NET_ADMIN
```

**In mesh (sidecars):**
- observability
- seaweedfs
- argocd
- openbao

**Excluded from mesh:**
- kyverno - Init containers need K8s API access before sidecar proxy starts
- istio-system - Control plane, no injection
- istio-ingress - Gateway, no injection

### mTLS Configuration

```yaml
# PeerAuthentication - mesh-wide STRICT
apiVersion: security.istio.io/v1
kind: PeerAuthentication
metadata:
  name: default
  namespace: istio-system
spec:
  mtls:
    mode: STRICT

# DestinationRule - per-service
apiVersion: networking.istio.io/v1
kind: DestinationRule
metadata:
  name: grafana
  namespace: istio-ingress
spec:
  host: grafana.observability.svc.cluster.local
  trafficPolicy:
    tls:
      mode: ISTIO_MUTUAL  # Use Istio certs for mTLS
```

**Exception:** Kiali uses `mode: DISABLE` because it runs in istio-system without a sidecar.

### Resource Configuration

Homelab-optimized (single N100 node):

| Component | CPU Request | Memory Request |
|-----------|-------------|----------------|
| istiod | 100m | 256Mi |
| ingress-gateway | 50m | 64Mi |
| per sidecar | 10m | 32Mi |
| kiali | 10m | 64Mi |

### Tracing Integration

Configured to send traces to Tempo:
```yaml
meshConfig:
  enableTracing: true
  defaultConfig:
    tracing:
      zipkin:
        address: tempo.observability.svc.cluster.local:9411
      sampling: 100.0  # 100% sampling for homelab
```

## Consequences

### Positive

- Full mTLS between all mesh services
- Kiali provides excellent service visualization
- Traces automatically captured and sent to Tempo
- Path-based routing via VirtualServices
- Foundation for advanced traffic management (canary, fault injection)

### Negative

- Resource overhead from sidecars (~10m CPU, 32Mi RAM each)
- Complexity - many CRDs and configuration options
- Debugging requires understanding Envoy proxy
- Some services excluded (Kyverno) due to init container requirements

### Operational Notes

**OpenBao auto-unseal:** OpenBao starts sealed after pod restart. Keys stored in `openbao-keys` secret for manual unseal. Consider auto-unseal in future.

**Kyverno exclusion:** Kyverno uses init containers that call K8s API. With Istio CNI, iptables rules redirect traffic before sidecar starts, breaking init container networking. Excluded via `istio-injection: disabled`.

## Lessons Learned

1. **`pilot.cni.enabled` is the critical setting** - Not just `global.istio_cni.enabled`. The sidecar injector checks `pilot.cni.enabled` to determine whether to add the iptables init container.

2. **PodSecurity must be privileged** - Istio CNI needs NET_ADMIN and NET_RAW capabilities which require privileged PSS.

3. **DestinationRules must match sidecar presence** - Services without sidecars (like Kiali in istio-system) must use `mode: DISABLE`, not `ISTIO_MUTUAL`.

4. **Envoy config is cached** - After changing DestinationRules, may need to restart ingress gateway for changes to take effect.

## References

- [Istio CNI Plugin](https://istio.io/latest/docs/setup/additional-setup/cni/)
- [Istio on Talos](https://www.talos.dev/v1.9/kubernetes-guides/network/istio/)
- [Kiali Documentation](https://kiali.io/docs/)
- [Istio + Tempo Integration](https://grafana.com/docs/tempo/latest/configuration/zipkin/)
