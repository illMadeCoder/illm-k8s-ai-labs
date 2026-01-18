# illm-k8s-ai-lab Roadmap (Final Structure)

**Updated:** 2026-01-17
**Structure:** 10 core phases + 18 appendices + AI-powered tech discovery

---

## Philosophy

**Component Isolation → System Composition → AI Evolution**

1. **Phases 1-9:** Deploy each component + Measure in isolation + FinOps cost analysis
2. **Phase 10:** Measure how components compose as a system + Cost per transaction end-to-end
3. **AI Discovery:** Web scraping to find emerging tech + Automated lab evolution

**Every phase answers:**
- ✅ How do I deploy this component?
- ✅ How do I measure its performance?
- ✅ What does it cost?
- ✅ How does it compare to alternatives?

---

## Core Learning Path (10 Phases)

| # | Phase | What You Build | What You Measure | FinOps Integration |
|---|-------|----------------|------------------|-------------------|
| **1** | Platform Bootstrap & GitOps | Hub, ArgoCD, Crossplane, OpenBao, Argo Workflows | Platform uptime, ArgoCD sync time | Platform running costs |
| **2** | CI/CD & Supply Chain | GitHub Actions, Cosign, SBOM, Kyverno, Image Updater | Build time, image size, scan duration | Build minutes, registry storage |
| **3** | Observability | Prometheus vs VictoriaMetrics, Loki vs ELK, Tempo vs Jaeger, SeaweedFS | Metrics cardinality, log volume, trace sampling | **Cost per metric, cost per GB logs, cost per trace** |
| **4** | Traffic Management | Gateway API, nginx vs Traefik vs Envoy comparison | Requests/sec, latency (p50/p95/p99), connection overhead | **Cost per request, ingress bandwidth cost** |
| **5** | Data & Persistence | PostgreSQL (CloudNativePG), Redis, backup/DR, **database benchmark** | Transactions/sec, query latency, IOPS | **Cost per transaction, cost per GB stored** |
| **6** | Security & Policy | TLS (cert-manager), secrets (ESO+OpenBao), RBAC, Kyverno, NetworkPolicy | Policy evaluation time, TLS handshake overhead | **Security tooling costs, compliance overhead** |
| **7** | Service Mesh | Istio vs Linkerd vs Cilium + **mesh overhead benchmark** | Sidecar latency overhead, control plane CPU/memory | **Mesh overhead cost (sidecar tax)** |
| **8** | Messaging & Events | Kafka vs RabbitMQ vs NATS + **messaging benchmark** | Messages/sec, end-to-end latency, fan-out performance | **Cost per million messages, retention storage cost** |
| **9** | Autoscaling & Resources | HPA, KEDA, VPA, cluster autoscaling | Scale-up time, resource efficiency, cost reduction | **Cost optimization via autoscaling** |
| **10** | **Performance & Cost Engineering** | **Runtime comparison + Full stack composition** | **System-level p99 latency, cost per transaction** | **Cost-efficiency as first-class metric** |

---

## Phase 10: The Grand Finale 🏆

**Goal:** Synthesize everything into data-driven system engineering

### 10.1 Runtime Performance Comparison
- Build identical API in: Go, Rust, .NET, Node.js, Bun
- Endpoints: /health, /json, /compute, /database
- Measure: RPS, latency distribution, memory, image size, cold start
- **Cost per million requests by runtime**

### 10.2 Full Stack Composition Benchmark
```
Client → Gateway → Service Mesh → App → Database
         ↓           ↓              ↓       ↓
      Measure    Measure        Measure Measure
```
- Deploy full stack: Runtime + Gateway (nginx/Envoy) + Mesh (Istio/Linkerd) + Database (PostgreSQL)
- Measure p99 latency through entire stack
- Isolate overhead: Baseline vs +Gateway vs +Mesh vs +Observability
- Answer: "What does each layer cost us in latency and $?"

### 10.3 System Trade-Off Analysis
- Performance vs Cost: "The mesh adds 5ms but costs $200/month - worth it?"
- Complexity vs Benefit: "3 layers of observability - which do we actually need?"
- Data-driven decision framework

### 10.4 Cost-Efficiency Dashboard
- Cost per transaction trending
- Cost breakdown by component
- Anomaly detection for cost spikes
- TCO comparison: Self-managed vs cloud-managed

**Portfolio Output:**
- Blog series: "I benchmarked 5 runtimes in Kubernetes"
- Interview material: "Here's how I reduced cost per transaction by 40%"
- GitHub showcase: Data-driven engineering

---

## AI-Powered Tech Discovery (Post Phase 10)

**Goal:** Keep the lab current with ecosystem evolution

### Components
1. **Web Scraping Jobs** (Argo Workflows)
   - Monitor CNCF landscape
   - Track GitHub trending (Kubernetes, Observability, etc.)
   - Parse tech blogs (The New Stack, KubeCon talks)
   - Reddit/HN for emerging patterns

2. **Technology Analysis**
   - Categorize: New component vs improvement vs noise
   - Assess: GitHub stars, contributor activity, production usage
   - Priority: P0 (disruptive), P1 (interesting), P2 (watch)

3. **Automated Suggestions**
   - "Cilium Tetragon is gaining traction - consider adding to Phase 6"
   - "Grafana Beyla (eBPF observability) - potential Phase 3 addition"
   - "Vector (log processor) now has 10k stars - compare vs Promtail?"

4. **Lab Evolution**
   - Generate experiment templates for new tech
   - Propose comparison scenarios
   - Suggest where to integrate in phases

**Implementation:**
- `experiments/ai-discovery/workflows/` - Argo Workflow CronJobs
- `experiments/ai-discovery/scrapers/` - Python scrapers with Beautiful Soup
- `experiments/ai-discovery/analysis/` - Claude integration for categorization
- `experiments/ai-discovery/suggestions/` - Markdown reports with recommendations

---

## Dependency Flow

```
Phase 1 (Platform)
   └─ ArgoCD, Crossplane, OpenBao, Workflows
         ↓
Phase 2 (CI/CD)
   └─ GitHub Actions, Cosign, SBOM, Kyverno
         ↓
Phase 3 (Observability) ←───────────────────────┐
   └─ Prometheus, Loki, Tempo, Grafana          │
         ↓                                       │
Phase 4 (Traffic Management)                     │
   └─ Gateway API, nginx/Traefik/Envoy          │
         ↓                                       │
Phase 5 (Data & Persistence)                     │
   └─ PostgreSQL, Redis, backups                │
         ↓                                       │
Phase 6 (Security & Policy) ─────────────────────┘
   └─ TLS, secrets, RBAC, NetworkPolicy
         ↓
Phase 7 (Service Mesh)
   └─ Istio, Linkerd, Cilium
         ↓
Phase 8 (Messaging & Events)
   └─ Kafka, RabbitMQ, NATS
         ↓
Phase 9 (Autoscaling & Resources)
   └─ HPA, KEDA, VPA, cluster autoscaling
         ↓
Phase 10 (Performance & Cost Engineering) ← THE CAPSTONE
   ├─ Runtime comparison (Go/Rust/.NET/Node/Bun)
   ├─ Full stack composition: Runtime→Gateway→Mesh→App→DB
   ├─ Cost per transaction end-to-end
   └─ System trade-off documentation
         ↓
AI-Powered Tech Discovery ← CONTINUOUS EVOLUTION
   └─ Web scraping → Analysis → Suggestions → Lab updates
```

---

## Advanced Topics (18 Appendices)

**Optional deep dives after core phases:**

### Cloud & Platform Engineering
- **A:** MLOps & AI Infrastructure
- **G:** Deployment Strategies (rolling, blue-green, canary, feature flags)
- **P:** Chaos Engineering
- **Q:** Advanced Workflow Patterns
- **R:** Internal Developer Platforms (Backstage)

### Security & Compliance
- **B:** Identity & Authentication
- **C:** PKI & Certificate Management
- **D:** Compliance & Security Operations
- **O:** SLSA Framework Deep Dive

### Architecture & Design
- **E:** Distributed Systems Fundamentals
- **F:** API Design & Contracts
- **H:** gRPC & HTTP/2 Patterns
- **K:** Event-Driven Architecture

### Performance & Operations
- **I:** Container & Runtime Internals
- **J:** Performance Engineering
- **L:** Database Internals
- **M:** SRE Practices & Incident Management
- **S:** Web Serving Internals

### Multi-Cloud
- **N:** Multi-Cloud & Disaster Recovery

---

## Current Status

**Phase 3: Observability - 60% Complete**

Validated:
- ✅ Prometheus + Grafana (metrics-app, RED dashboards)
- ✅ Victoria Metrics comparison
- ✅ SeaweedFS object storage

Backlog (needs validation):
- [ ] Loki tutorial + cost per GB logs
- [ ] Elasticsearch tutorial
- [ ] Logging comparison (Loki vs ELK)
- [ ] OpenTelemetry tutorial + cost per trace
- [ ] Tempo tutorial
- [ ] Jaeger tutorial
- [ ] Tracing comparison (Tempo vs Jaeger)
- [ ] Pyrra SLOs
- [ ] Observability cost management (cardinality, retention)

**Next:** Validate all 9 backlog experiments (2 weeks)

---

## Timeline to Portfolio-Ready

| Milestone | Duration | Cumulative |
|-----------|----------|------------|
| Phase 3 validation | 2 weeks | 2 weeks |
| Roadmap restructure | 1 week | 3 weeks |
| Phase 4 (Traffic Management) | 3-4 weeks | 6-7 weeks |
| Phase 5 (Data & Persistence) | 3-4 weeks | 9-11 weeks |
| Phase 6 (Security & Policy) | 4-5 weeks | 13-16 weeks |
| Phase 7 (Service Mesh) | 3-4 weeks | 16-20 weeks |
| Phase 8 (Messaging & Events) | 3-4 weeks | 19-24 weeks |
| Phase 9 (Autoscaling) | 2-3 weeks | 21-27 weeks |
| Phase 10 (Grand Finale) | 3-4 weeks | 24-31 weeks |
| AI Tech Discovery | 2-3 weeks | 26-34 weeks |

**Total:** ~5-8 months to complete (realistically 6 months)

---

## Success Metrics

### Current State
- ✅ 2 phases complete (Platform, CI/CD)
- 🚧 1 phase in progress (Observability ~60%)
- 📝 13 ADRs documented
- 🏗️ 8 experiments validated

### Target State (6 months)
- ✅ 10 core phases complete
- 📝 40+ ADRs documenting decisions
- 🏗️ 50-55 experiments validated
- 🎯 Phase 10 capstone: Runtime comparison + full stack benchmark
- 🤖 AI tech discovery running continuously
- 💼 Portfolio-ready: Blog posts, GitHub showcase, interview material

---

## What Changed (Consolidation)

**Before:** 16 phases, ~80-90 experiments, 10-12 months
**After:** 10 phases, ~50-55 experiments, 5-6 months

### Changes
- ✅ **Kept Phase 15** - Elevated to Phase 10 (the capstone)
- ✅ **FinOps integrated** - Every phase now includes cost measurements
- ✅ **Benchmarks preserved** - Database (Phase 5), Messaging (Phase 8), Mesh (Phase 7), Runtime (Phase 10)
- ✅ **Security consolidated** - Phase 7 + 8 → Phase 6
- ⬇️ **Moved to appendices:** Deployment strategies, Chaos, gRPC deep dive, Advanced workflows, Backstage, Web serving details

### Why This Works
1. **Component isolation** (Phases 3-9) teaches measurement expertise
2. **System composition** (Phase 10) teaches full-stack optimization
3. **FinOps first-class** demonstrates cost-conscious engineering
4. **AI discovery** demonstrates forward-thinking architecture

---

## Quick Start

```bash
# Prerequisites: Docker, kubectl, task, helm

task hub:bootstrap                      # Create hub cluster
task hub:conduct -- prometheus-tutorial # Run an experiment
task hub:down -- prometheus-tutorial    # Cleanup
task hub:destroy                        # Destroy cluster
```

---

## Documents

- [Strategic Review](strategic-review-2026-01.md) - Initial assessment
- [Consolidation Analysis](roadmap-consolidation-analysis.md) - Detailed phase analysis
- [Consolidation Summary](roadmap-consolidation-summary.md) - Visual before/after
- [GitOps Patterns](gitops-patterns.md)
- [ADRs](adrs/)
