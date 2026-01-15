## Phase 4: Traffic Management

*Control how traffic flows before learning deployment strategies that depend on it.*

### 4.1 Gateway Tutorial: Ingress → Gateway API Evolution

**Goal:** Understand L7 traffic management from legacy Ingress through modern Gateway API

**Learning objectives:**
- Master Kubernetes Ingress and its limitations
- Understand Gateway API resources (Gateway, HTTPRoute, GRPCRoute)
- Experience the migration path from Ingress to Gateway API
- Implement advanced routing and traffic manipulation patterns

**Scenario:** `experiments/scenarios/gateway-tutorial/`

**Part 1: Ingress Basics**
- [ ] Deploy nginx-ingress controller
- [ ] Configure basic Ingress resources:
  - [ ] Path-based routing
  - [ ] Host-based routing (virtual hosts)
  - [ ] TLS termination with cert-manager
- [ ] Understand Ingress annotations pattern

**Part 2: Hitting Ingress Limitations**
- [ ] Attempt rate limiting (annotation hell begins)
- [ ] Try header-based routing (limited support)
- [ ] Add authentication (nginx-specific annotations)
- [ ] Traffic splitting for canary (awkward with Ingress)
- [ ] Document the pain points

**Part 3: Migrate to Gateway API**
- [ ] Deploy Envoy Gateway (CNCF reference implementation)
- [ ] Create Gateway resource
- [ ] Migrate Ingress rules to HTTPRoute
- [ ] Compare configuration complexity
- [ ] Side-by-side: same routes, both approaches

**Part 4: Gateway API Deep Dive**
- [ ] HTTPRoute patterns:
  - [ ] Path/host/header/query/method matching
  - [ ] Weight-based traffic splitting
  - [ ] Request mirroring
  - [ ] URL rewriting
  - [ ] Header modification
  - [ ] Redirects
- [ ] Advanced features:
  - [ ] Timeouts and retries
  - [ ] Rate limiting (BackendTrafficPolicy)
  - [ ] CORS configuration
- [ ] TLS configuration:
  - [ ] TLS termination
  - [ ] TLS passthrough
  - [ ] mTLS with client certificates
- [ ] Multi-gateway patterns:
  - [ ] Internal vs external gateways
  - [ ] Namespace isolation (ReferenceGrant)
- [ ] GRPCRoute for gRPC services

**Deliverables:**
- [ ] Working tutorial with all four parts
- [ ] **ADR:** Gateway API implementation choice (Envoy Gateway)
- [ ] Comparison table: Ingress vs Gateway API

---

### 4.2 Gateway Comparison: In-Cluster Implementations

**Goal:** Compare in-cluster gateway implementations for informed selection

**Learning objectives:**
- Understand trade-offs between gateway implementations
- Compare configuration patterns and complexity
- Evaluate feature availability and resource consumption

**Scenario:** `experiments/scenarios/gateway-comparison/`

**Implementations to compare:**
- [ ] nginx-ingress (most widely deployed)
- [ ] Traefik (popular, good Gateway API support)
- [ ] Envoy Gateway (CNCF reference, pure Gateway API)

**Same demo app, same routes on each:**
- [ ] Path-based routing to multiple services
- [ ] Host-based virtual hosting
- [ ] TLS termination
- [ ] Rate limiting
- [ ] Header manipulation

**Comparison metrics:**
| Metric | nginx | Traefik | Envoy Gateway |
|--------|-------|---------|---------------|
| Config complexity | | | |
| Gateway API support | | | |
| Resource usage (CPU/mem) | | | |
| Feature completeness | | | |
| Observability integration | | | |
| Community/docs quality | | | |

**Deliverables:**
- [ ] Side-by-side deployment of all three
- [ ] Comparison matrix with findings
- [ ] Recommendation criteria document

---

### 4.3 Cloud Gateway Comparison: Managed vs In-Cluster

**Goal:** Compare cloud-native application gateways with in-cluster solutions

**Learning objectives:**
- Understand cloud provider L7 gateway offerings
- Evaluate cost, performance, and operational trade-offs
- Make informed decisions for production architectures

**Scenario:** `experiments/scenarios/cloud-gateway-comparison/`

**Implementations to compare:**
- [ ] AWS ALB Ingress Controller (provisions AWS ALBs)
- [ ] Azure AGIC (provisions Azure Application Gateways)
- [ ] Envoy Gateway in-cluster (baseline comparison)

**Infrastructure:**
- [ ] AWS EKS cluster via Crossplane
- [ ] Azure AKS cluster via Crossplane
- [ ] Talos cluster (in-cluster baseline)

**Comparison Metrics:**

| Category | Metrics |
|----------|---------|
| **Latency** | p50, p95, p99 request latency |
| **Throughput** | Max RPS, RPS at 10/100/1000 concurrency |
| **Provisioning** | Time to create gateway, time to add route |
| **Cost** | Hourly gateway cost, per-request cost, data transfer |
| **Config propagation** | Time from kubectl apply → traffic flowing |
| **Reliability** | Multi-AZ behavior, failover time, SLA |
| **Features** | Rate limiting, auth (JWT/mTLS), WAF, WebSocket/gRPC |
| **Observability** | Metrics export, access logs, tracing integration |
| **Blast radius** | Impact of misconfig, rollback ease |
| **Vendor lock-in** | Portability of configuration |

**Load testing:**
- [ ] Use k6 for consistent load generation
- [ ] Test at multiple concurrency levels
- [ ] Measure during route changes

**Deliverables:**
- [ ] Crossplane compositions for AWS ALB IC and Azure AGIC
- [ ] Automated metrics collection
- [ ] Scoring matrix with all metrics
- [ ] **ADR:** When to use cloud-native vs in-cluster gateways
- [ ] Cost calculator for different traffic levels

---

## Dependencies

```
gateway-tutorial
       ↓
gateway-comparison
       ↓
cloud-gateway-comparison
```

## Backlog

After completing Phase 4:
- [ ] Play through gateway-tutorial
- [ ] Play through gateway-comparison
- [ ] Play through cloud-gateway-comparison
