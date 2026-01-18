# Advanced Metrics & eBPF Tooling Strategy

**Branch:** `claude/review-project-roadmap-psMLb`
**Status:** Planning / Not Yet Committed to Roadmap
**Date:** 2026-01-17

---

## Problem Statement

**Current metrics focus:** CPU and RAM utilization

**Missing dimensions:**
- **I/O performance** - Disk latency, IOPS, throughput
- **Network I/O** - TCP connections, packet loss, retransmits
- **System calls** - Overhead from kernel interactions
- **File system** - VFS operations, page cache efficiency

**Why this matters:**
- A database with 20% CPU can still be I/O bound
- Mesh sidecars can saturate network buffers without high CPU
- Container image pulls can bottleneck on IOPS, not CPU

---

## The Missing Metrics

### 1. Block I/O Metrics

**What to measure:**
- **Latency:** Time from I/O request to completion (p50, p95, p99)
- **IOPS:** Read/write operations per second
- **Throughput:** MB/s read/write
- **Queue depth:** How many I/O operations are waiting
- **Device saturation:** % time device is busy

**Why it matters:**
```
Scenario: PostgreSQL performance degradation
- CPU: 15% (looks fine)
- RAM: 40% (looks fine)
- Disk latency: p99 = 500ms (PROBLEM!)
- IOPS: 95% of limit (saturated)

Root cause: Database write amplification saturating disk
```

**Tools:**
- **biosnoop** (eBPF) - Trace block I/O with latency
- **biotop** (eBPF) - Top-like view of block I/O by process
- **iostat** - Classic I/O stats
- **blktrace** - Detailed block layer tracing

### 2. Network I/O Metrics

**What to measure:**
- **TCP connections:** Active connections, connection rate
- **Retransmits:** % of packets retransmitted (indicates congestion)
- **Packet loss:** Dropped packets (network issues)
- **Socket buffer saturation:** Send/receive buffers full
- **Bandwidth:** Network throughput vs capacity

**Why it matters:**
```
Scenario: Service mesh latency spikes
- CPU: 30% (looks fine)
- RAM: 50% (looks fine)
- TCP retransmits: 5% (PROBLEM!)
- Socket buffer: 90% full (congestion)

Root cause: Insufficient socket buffers for high-throughput mesh
```

**Tools:**
- **tcptop** (eBPF) - TCP connections by throughput
- **tcpretrans** (eBPF) - TCP retransmit tracer
- **tcplife** (eBPF) - TCP connection lifespan
- **ss** (socket statistics) - Connection states
- **netstat** - Network statistics

### 3. File System Metrics

**What to measure:**
- **VFS operations:** read/write/open/close rates
- **Page cache hit rate:** % of reads from memory vs disk
- **Inode usage:** File descriptor exhaustion
- **Filesystem latency:** Time for VFS operations

**Why it matters:**
```
Scenario: Observability stack slow queries
- CPU: 25% (looks fine)
- RAM: 60% (looks fine)
- Page cache hit rate: 40% (PROBLEM! Should be >90%)
- VFS reads: 10k/sec hitting disk

Root cause: Insufficient RAM for index caching
```

**Tools:**
- **vfsstat** (eBPF) - VFS operation statistics
- **cachestat** (eBPF) - Page cache hit/miss stats
- **filetop** (eBPF) - File reads/writes by process

### 4. System Call Metrics

**What to measure:**
- **Syscall latency:** Time spent in kernel
- **Syscall frequency:** Number of syscalls per second
- **Context switches:** Process switching overhead

**Why it matters:**
```
Scenario: High-throughput API slow despite low CPU
- CPU: 35% (looks fine)
- Syscalls: 100k/sec (PROBLEM!)
- Context switches: 50k/sec (high overhead)

Root cause: Inefficient I/O patterns (lots of small reads)
```

**Tools:**
- **syscount** (eBPF) - Count syscalls by type
- **funclatency** (eBPF) - Latency histogram for syscalls

---

## eBPF: The Key Technology

**Why eBPF?**
- **Zero overhead when not tracing** - No permanent performance impact
- **Production-safe** - Kernel verifier prevents crashes
- **Real-time visibility** - See what's happening NOW, not after reboot
- **No instrumentation** - No code changes required

**BCC Tools Suite:**

| Tool | What It Measures | Use Case |
|------|------------------|----------|
| **biosnoop** | Block I/O latency per operation | Find slow disk operations |
| **biotop** | Top processes by block I/O | Identify I/O-heavy workloads |
| **tcptop** | TCP throughput by connection | Find network bandwidth hogs |
| **tcpretrans** | TCP retransmits | Diagnose network congestion |
| **tcplife** | TCP connection duration | Understand connection churn |
| **vfsstat** | VFS operation rate | Filesystem operation rate |
| **cachestat** | Page cache hit rate | Memory efficiency |
| **execsnoop** | Process execution | Container startup analysis |
| **opensnoop** | File opens | Track file access patterns |

**bpftrace:**
- One-liner custom tracing scripts
- Ad-hoc investigation during incidents
- Example: `bpftrace -e 'kprobe:tcp_retransmit_skb { @[comm] = count(); }'`

---

## Integration into Roadmap Phases

### Phase 3: Observability (Current - Add eBPF)

**New Sub-Phase: 3.6 eBPF & System Metrics**

**Goal:** Add I/O and network metrics to observability stack

**Components to deploy:**
- **Pixie** (CNCF) - Auto-instrumented observability via eBPF
  - No code changes required
  - Captures network traffic, system calls, app-level traces
  - UI for exploring eBPF data

- **Parca** (CNCF) - Continuous profiling via eBPF
  - CPU flamegraphs showing where time is spent
  - Memory profiling
  - Correlate with other metrics

- **Tetragon** (Cilium) - Runtime security observability
  - Process execution
  - Network connections
  - File access
  - Syscall filtering

**What to measure:**
- Deploy database workload
- Capture metrics with traditional tools (Prometheus) vs eBPF tools
- Compare visibility:
  - Prometheus: CPU 20%, RAM 40%
  - eBPF: + Disk p99 latency 500ms, IOPS 95% saturated

**Experiments:**
1. **I/O bottleneck detection**
   - Create `io-bound-app` that writes frequently
   - Measure with `biosnoop` to see latency spikes
   - Show how CPU/RAM metrics miss the problem

2. **Network congestion analysis**
   - Deploy service mesh with high throughput
   - Use `tcptop` and `tcpretrans` to find retransmits
   - Correlate with latency spikes

3. **Page cache efficiency**
   - Database query workload
   - Use `cachestat` to measure hit rate
   - Show impact of RAM on query performance

**FinOps Angle:**
- "We thought we needed more CPU, but eBPF showed we needed faster disks"
- "Saved $500/month by right-sizing based on I/O, not CPU"

**ADR-XXX:** Document why eBPF is essential for production observability

---

### Phase 5: Data & Persistence (Add I/O Benchmarking)

**Enhanced Database Benchmark (5.X)**

**Current plan:** pgbench with TPS and latency
**Enhanced:** Add I/O-aware benchmarking

**Metrics to add:**
- **Disk latency breakdown:**
  - Application latency (total time)
  - Database processing time
  - Disk I/O time (measured with biosnoop)

- **I/O patterns:**
  - Sequential vs random I/O ratio
  - Read vs write ratio
  - Average I/O size

- **Page cache efficiency:**
  - Cache hit rate
  - Dirty pages
  - Writeback frequency

**Experiment: "Why Your Database is Slow"**
```
Scenario 1: CPU-bound (sorting, aggregation)
├─ CPU: 80%
├─ Disk latency: <1ms
└─ Solution: More CPU cores

Scenario 2: I/O-bound (table scans, indexes)
├─ CPU: 20%
├─ Disk latency: 100ms
└─ Solution: Faster disks or more RAM for cache

Scenario 3: Network-bound (replication lag)
├─ CPU: 15%
├─ Disk latency: 2ms
├─ Network retransmits: 10%
└─ Solution: Better network or reduce replication traffic
```

**Cost Analysis:**
- Compare: 2x CPU @ $100/month vs faster disk @ $50/month
- Use eBPF to prove disk is the bottleneck, not CPU
- **FinOps win:** Right-size based on actual bottleneck

---

### Phase 7: Service Mesh (Add Network I/O Analysis)

**Enhanced Mesh Overhead Benchmark (7.X)**

**Current plan:** Measure sidecar CPU/memory overhead
**Enhanced:** Add network I/O overhead

**Metrics to add:**
- **TCP connection overhead:**
  - Connection establishment time (SYN → SYN-ACK → ACK)
  - Connections/second rate
  - Connection pool saturation

- **Network retransmits:**
  - Baseline (no mesh): 0.1% retransmit rate
  - With Istio: 2% retransmit rate (sidecar adds latency → timeout → retransmit)

- **Socket buffer usage:**
  - Send buffer saturation
  - Receive buffer saturation
  - Impact on throughput

**Experiment: "Service Mesh Network Tax"**
```
Baseline (no mesh):
├─ Latency: p99 = 50ms
├─ Retransmits: 0.1%
└─ Throughput: 10k RPS

Istio sidecar:
├─ Latency: p99 = 55ms (+5ms from eBPF tracing)
├─ Retransmits: 2% (+1.9% - SIGNIFICANT)
├─ Throughput: 9.5k RPS (-5%)
└─ Socket buffers: 80% full (approaching saturation)

Root cause: Sidecar adds extra hop → latency → bufferbloat → retransmits
```

**eBPF Tools:**
- `tcptop` - Identify high-throughput connections
- `tcpretrans` - Trace retransmits back to specific services
- `tcplife` - Measure connection lifespan (short-lived = overhead)

**Cost Impact:**
- "Mesh retransmits cost us 5% throughput = need 5% more pods = $X/month"
- "Tuning socket buffers saved retransmits, recovered 5% throughput"

---

### Phase 10: Performance & Cost Engineering (Full Stack I/O Analysis)

**Enhanced Capstone: I/O in the Full Stack**

**Current plan:** Runtime → Gateway → Mesh → App → Database (latency only)
**Enhanced:** Add I/O breakdown at each layer

**Full Stack I/O Attribution:**
```
End-to-end p99 latency: 200ms

Breakdown with eBPF:
├─ Gateway (nginx): 5ms
│  ├─ CPU: 2ms
│  ├─ Network I/O: 3ms (routing decision, connection establishment)
│  └─ Disk I/O: 0ms (cached)
│
├─ Mesh (Istio sidecar): 10ms
│  ├─ CPU: 3ms
│  ├─ Network I/O: 7ms (mTLS handshake, routing to app)
│  └─ Disk I/O: 0ms
│
├─ App (Go runtime): 50ms
│  ├─ CPU: 40ms (business logic)
│  ├─ Network I/O: 5ms (database connection)
│  └─ Disk I/O: 5ms (local file cache lookup)
│
├─ Database (PostgreSQL): 130ms
│  ├─ CPU: 20ms (query planning, execution)
│  ├─ Network I/O: 10ms (replication to standby)
│  └─ Disk I/O: 100ms (INDEX SCAN - 95% of database latency!)
│
└─ Messaging (Kafka): 5ms
   ├─ CPU: 1ms
   ├─ Network I/O: 2ms
   └─ Disk I/O: 2ms (write to partition log)

INSIGHT: 100ms of 200ms (50%) is disk I/O in PostgreSQL
ACTION: Add read replicas or increase RAM for index caching
COST: Read replica $100/month vs more RAM $50/month → choose RAM
```

**Runtime Comparison with I/O:**

| Runtime | CPU Time | Network I/O | Disk I/O | Total Latency |
|---------|----------|-------------|----------|---------------|
| Go | 40ms | 5ms | 5ms | 50ms |
| Rust | 35ms | 5ms | 5ms | 45ms (faster CPU) |
| .NET | 50ms | 6ms | 4ms | 60ms |
| Node.js | 60ms | 7ms | 3ms | 70ms (async I/O helps) |
| Bun | 45ms | 5ms | 5ms | 55ms |

**Key Insight:** Node.js has highest CPU time but best async I/O handling

**eBPF-Powered Cost Optimization:**
1. **Find the bottleneck** (eBPF shows it's disk I/O, not CPU)
2. **Right-size resources** (add RAM for cache, not CPU cores)
3. **Measure savings** ($50/month RAM vs $200/month CPU)
4. **Validate with eBPF** (page cache hit rate improves 40% → 95%)

---

## Tooling Deployment Strategy

### Option 1: eBPF Tools as Sidecar (Lightweight)

**Deploy BCC tools in privileged pods:**
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ebpf-tools
spec:
  hostPID: true      # See all processes
  hostNetwork: true  # See all network traffic
  containers:
  - name: bcc-tools
    image: zlim/bcc   # BCC tools container
    securityContext:
      privileged: true  # Required for eBPF
    volumeMounts:
    - name: sys
      mountPath: /sys
    - name: debugfs
      mountPath: /sys/kernel/debug
```

**Usage:**
```bash
kubectl exec ebpf-tools -- biosnoop     # Trace block I/O
kubectl exec ebpf-tools -- tcptop       # TCP throughput
kubectl exec ebpf-tools -- cachestat    # Page cache stats
```

**Pros:**
- Minimal overhead
- Full control over tracing
- No permanent agents

**Cons:**
- Manual usage (not automated)
- No historical data
- No dashboards

---

### Option 2: Pixie (CNCF - Recommended for Observability Phase)

**Deploy via Helm:**
```bash
helm install pixie pixie-operator/pixie-operator-chart \
  --set deployKey=$PIXIE_DEPLOY_KEY
```

**What you get:**
- **Auto-instrumented** - No code changes, no sidecars
- **Network map** - Visual topology of service communication
- **Request tracing** - HTTP/gRPC/DNS/Kafka traces automatically
- **Resource profiling** - CPU flamegraphs, memory usage
- **Live debugging** - Real-time queries with PxL (Pixie Language)

**Example PxL query:**
```python
# Find slow database queries
px.display(px.DataFrame(
  table='mysql.query',
  start_time='-5m'
).groupby('query')
  .agg(latency_p99=('latency', px.quantiles, 0.99))
  .filter(latency_p99 > 100)  # >100ms
)
```

**FinOps Integration:**
- Query: "Show me the 10 slowest endpoints by total CPU time"
- Output: `/api/export` uses 40% of total CPU → optimize or rate-limit
- Savings: Reduce over-provisioning by 20% = $X/month

**ADR:** Pixie for Phase 3 observability (auto-instrumentation + eBPF visibility)

---

### Option 3: Parca (CNCF - Recommended for Performance Phase)

**Deploy via Helm:**
```bash
helm install parca parca/parca
```

**What you get:**
- **Continuous profiling** - Always-on CPU/memory profiling
- **Flamegraphs** - Visual representation of where time is spent
- **Differential profiling** - Compare before/after optimization
- **Historical analysis** - "What changed between 2pm and 3pm?"

**Use Cases:**
- **Phase 10 runtime comparison:**
  - Profile Go vs Rust vs .NET CPU usage
  - See exactly where each runtime spends time
  - Identify optimization opportunities

- **Database optimization:**
  - Flamegraph shows 60% time in `btree_search`
  - Add index to reduce search time
  - Re-profile: Now 10% time in `btree_search`

**FinOps Integration:**
- Before: CPU 70%, need to scale up
- Profiling: 50% of CPU in inefficient JSON parsing
- Fix: Use faster parser
- After: CPU 40%, cancel scale-up, save $200/month

---

### Option 4: Tetragon (Cilium - Recommended for Security Phase)

**Deploy via Helm:**
```bash
helm install tetragon cilium/tetragon
```

**What you get:**
- **Process execution** - Track all process spawns
- **Network connections** - TCP/UDP connections by pod
- **File access** - Which processes access which files
- **Syscall filtering** - Block dangerous syscalls

**Security Observability Use Cases:**
- **Detect unexpected behavior:**
  - Pod spawns `/bin/bash` (should only run app binary)
  - Pod makes outbound connection to unknown IP
  - Pod accesses `/etc/shadow`

**Performance Observability Use Cases:**
- **Process churn:** Identify pods constantly forking processes
- **File descriptor leaks:** Track file opens without closes
- **Network connection leaks:** Track connections without cleanup

**Phase 6 Integration:**
- Deploy Tetragon in Phase 6 (Security)
- Use for security AND performance observability
- Show how security tools can reveal performance issues

---

## Metrics Architecture

### Traditional Stack (CPU/RAM only)
```
┌─────────────────────────────────────┐
│         Applications                │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│    Prometheus (cAdvisor metrics)    │
│    - container_cpu_usage            │
│    - container_memory_usage         │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│          Grafana Dashboards         │
│    CPU: 20%  ✅                     │
│    RAM: 40%  ✅                     │
│    But... why is it slow?? 🤔      │
└─────────────────────────────────────┘
```

### Enhanced Stack (CPU/RAM/I/O/Network)
```
┌─────────────────────────────────────┐
│         Applications                │
└──┬────────┬─────────┬──────────┬───┘
   │        │         │          │
   │        │         │          │ (eBPF kernel hooks)
   │        │         │          │
   ▼        ▼         ▼          ▼
┌──────┐ ┌──────┐ ┌──────┐ ┌──────┐
│ CPU  │ │ RAM  │ │ I/O  │ │ Net  │
│ 20%  │ │ 40%  │ │ p99: │ │Retx: │
│      │ │      │ │500ms │ │ 5%   │
└──┬───┘ └──┬───┘ └──┬───┘ └──┬───┘
   │        │         │        │
   └────────┴─────────┴────────┘
                 │
                 ▼
┌─────────────────────────────────────┐
│     Unified Observability           │
│  (Prometheus + Pixie + Parca)       │
└────────────┬────────────────────────┘
             │
             ▼
┌─────────────────────────────────────┐
│       Grafana + Pixie UI            │
│  CPU: 20%  ✅                       │
│  RAM: 40%  ✅                       │
│  Disk: p99=500ms ❌ BOTTLENECK!     │
│  IOPS: 95% saturated                │
│  Solution: Faster disks or more RAM │
└─────────────────────────────────────┘
```

---

## Implementation Phases

### Phase 3.6: eBPF & System Metrics (NEW)

**Deliverables:**
1. Deploy Pixie for auto-instrumentation
2. Deploy Parca for continuous profiling
3. Create `ebpf-tutorial` experiment
4. Create `io-bottleneck-detection` experiment
5. Document eBPF tools in runbook

**Experiments:**
- **ebpf-tutorial:** Hands-on with biosnoop, tcptop, cachestat
- **io-bottleneck:** Create I/O-bound app, detect with eBPF
- **network-analysis:** Service mesh latency attribution with tcpretrans

**ADR-XXX:** eBPF for production observability

---

### Phase 5.X: I/O-Aware Database Benchmarking (ENHANCED)

**Deliverables:**
1. Add biosnoop to database benchmark
2. Measure page cache hit rate with cachestat
3. Document I/O patterns (sequential vs random)
4. Compare cloud disk types (standard vs premium vs ultra)

**Cost Analysis:**
- Measure: "Does premium disk justify 3x cost?"
- Benchmark: Standard (100 IOPS) vs Premium (1000 IOPS)
- Result: For our workload, premium saves 200ms → worth it for user-facing API

---

### Phase 7.X: Network I/O in Service Mesh (ENHANCED)

**Deliverables:**
1. Baseline network metrics (no mesh)
2. Istio with tcptop/tcpretrans analysis
3. Identify retransmit sources
4. Tune socket buffers based on data

**Cost Analysis:**
- Measure: "What's the real network cost of the mesh?"
- Result: 5% retransmit rate = 5% throughput loss = need 5% more pods
- Optimization: Tune buffers, reduce retransmits to 1%, save $X/month

---

### Phase 10.X: Full Stack I/O Attribution (ENHANCED)

**Deliverables:**
1. End-to-end trace with I/O breakdown
2. Per-layer I/O contribution analysis
3. Runtime comparison with I/O metrics
4. Cost optimization based on actual bottleneck

**Grand Finale Output:**
- Blog: "I used eBPF to find why our API was slow (spoiler: disk I/O)"
- Interview: "Saved 40% on infrastructure by right-sizing based on eBPF data"
- GitHub: eBPF-powered performance engineering showcase

---

## Tools Summary

| Tool | Phase | Purpose | Cost Impact |
|------|-------|---------|-------------|
| **Pixie** | 3 | Auto-instrumented observability | Find over-provisioned services |
| **Parca** | 10 | Continuous profiling | Optimize hot code paths |
| **Tetragon** | 6 | Security + perf observability | Detect process/network anomalies |
| **biosnoop** | 5 | Block I/O latency | Right-size disk vs RAM |
| **tcptop** | 7 | TCP throughput | Find network bottlenecks |
| **tcpretrans** | 7 | TCP retransmit analysis | Reduce mesh overhead |
| **cachestat** | 5 | Page cache efficiency | Optimize RAM allocation |

---

## Open Questions

1. **Should Pixie be in Phase 3 or its own phase?**
   - Option A: Phase 3.6 (eBPF & System Metrics)
   - Option B: Phase 3.7 (Auto-Instrumented Observability)

2. **How much eBPF in core vs appendix?**
   - Core: Basic tools (biosnoop, tcptop, cachestat)
   - Appendix: Advanced (bpftrace custom scripts, kernel internals)

3. **FinOps dashboard: Include I/O costs?**
   - Phase 10: Cost per transaction should include I/O breakdown
   - Example: "Disk I/O = $0.005 of $0.015 total cost per transaction"

4. **Should we have an "eBPF Deep Dive" appendix?**
   - Appendix T: eBPF Internals & Custom Tracing
   - For those who want to write bpftrace scripts

---

## Next Steps

1. **Review this document** - Discuss scope and priorities
2. **Decide on Pixie integration** - Phase 3.6 or separate?
3. **Plan eBPF experiments** - Which tools to include in tutorials
4. **Update Phase 3 plan** - Add eBPF section
5. **Update Phase 10 plan** - Add I/O attribution

**Status:** Planning - waiting for approval before committing to roadmap

---

## References

- [BCC Tools](https://github.com/iovisor/bcc) - eBPF tracing toolkit
- [Pixie](https://px.dev/) - CNCF auto-instrumented observability
- [Parca](https://www.parca.dev/) - CNCF continuous profiling
- [Tetragon](https://tetragon.io/) - Cilium runtime security observability
- [eBPF Documentation](https://ebpf.io/) - Official eBPF docs
- [Brendan Gregg's Blog](https://www.brendangregg.com/ebpf.html) - eBPF performance analysis
