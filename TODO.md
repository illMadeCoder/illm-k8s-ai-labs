# illm-k8s-lab TODO

## Phase 1: Platform Setup

### Spacelift + Crossplane Multi-Cloud Setup

**Prerequisites**: None - this is the foundation for everything else

#### Spacelift Setup
- [ ] Create Spacelift account at [spacelift.io](https://spacelift.io) (free tier)
- [ ] Connect GitHub repository (illMadeCoder/illm-k8s-lab)
- [ ] Create `spacelift-root` stack manually (administrative=true, path=spacelift-stacks/base)
- [ ] Configure `azure-credentials` context with environment variables:
  - `ARM_CLIENT_ID`
  - `ARM_CLIENT_SECRET` (write-only)
  - `ARM_SUBSCRIPTION_ID`
  - `ARM_TENANT_ID`
- [ ] Configure `aws-credentials` context with environment variables:
  - `AWS_ACCESS_KEY_ID`
  - `AWS_SECRET_ACCESS_KEY` (write-only)
  - `AWS_DEFAULT_REGION`
- [ ] Trigger initial run on `spacelift-root` to create child stacks
- [ ] Verify foundation stacks appear (azure-foundation, aws-foundation)

#### Crossplane Local Testing
- [ ] Deploy Crossplane to minikube: `task deploy:crossplane`
- [ ] Verify providers are healthy: `task crossplane:status`
- [ ] Test XRDs are installed: `kubectl get xrds`
- [ ] (Optional) Configure LocalStack for AWS resource mocking

#### Validation
- [ ] Deploy multi-cloud-demo to Azure via Spacelift
- [ ] Verify Crossplane claims provision real cloud resources
- [ ] Run load test against deployed application
- [ ] Destroy infrastructure and verify cleanup

---

## Phase 2: Observability Foundations

*Learn to see what's happening before trying to improve it. These skills are essential for all future experiments.*

### Prometheus & Grafana Deep Dive

**Prerequisites**: Basic Kubernetes knowledge
**Why first**: You need metrics to measure everything else

**Goal**: Master metrics collection, PromQL queries, alerting, and dashboard creation

**Learning objectives**:
- Understand Prometheus architecture (scraping, storage, federation)
- Write effective PromQL queries
- Create actionable Grafana dashboards
- Configure alerting rules and notification channels

**Tasks**:
- [ ] Create `experiments/prometheus-tutorial/` directory structure
- [ ] Deploy sample application with custom metrics:
  - [ ] Counter (requests_total)
  - [ ] Gauge (connections_active)
  - [ ] Histogram (request_duration_seconds)
  - [ ] Summary (response_size_bytes)
- [ ] Instrument application with Prometheus client library (Go or Python)
- [ ] Create ServiceMonitor for automatic scrape discovery
- [ ] Write PromQL queries tutorial:
  - [ ] Rate and irate for counters
  - [ ] Aggregations (sum, avg, max by labels)
  - [ ] Histogram quantiles
  - [ ] Absent() for alert on missing metrics
  - [ ] Predict_linear() for capacity planning
- [ ] Build Grafana dashboards:
  - [ ] RED metrics dashboard (Rate, Errors, Duration)
  - [ ] USE metrics dashboard (Utilization, Saturation, Errors)
  - [ ] Application-specific dashboard with variables
- [ ] Configure alerting:
  - [ ] PrometheusRule CRD for alert definitions
  - [ ] Alertmanager routing and silences
  - [ ] Alert fatigue prevention (grouping, inhibition)
- [ ] Document query optimization tips

---

### Loki & Log Aggregation

**Prerequisites**: Prometheus & Grafana Deep Dive
**Why here**: Logs complement metrics - together they tell the full story

**Goal**: Master centralized logging with Loki, Promtail, and Grafana

**Learning objectives**:
- Understand Loki's label-based indexing (vs full-text)
- Write effective LogQL queries
- Correlate logs with metrics
- Set up log-based alerts

**Tasks**:
- [ ] Create `experiments/loki-tutorial/` directory structure
- [ ] Deploy application with structured logging (JSON)
- [ ] Configure Promtail:
  - [ ] Pipeline stages (regex, json, labels)
  - [ ] Drop/keep log lines
  - [ ] Multi-line log handling
- [ ] Write LogQL queries tutorial:
  - [ ] Label matchers and filters
  - [ ] Log line parsing (pattern, regexp, json)
  - [ ] Metric queries from logs (rate, count_over_time)
  - [ ] Unwrap for numeric extraction
- [ ] Build log dashboards in Grafana:
  - [ ] Log panel with live tail
  - [ ] Log volume visualization
  - [ ] Correlation with metrics (split view)
- [ ] Configure log-based alerts:
  - [ ] Error rate from logs
  - [ ] Pattern matching alerts
- [ ] Compare Loki vs ELK:
  - [ ] Resource usage
  - [ ] Query capabilities
  - [ ] Operational complexity
- [ ] Document retention and storage optimization

---

### OpenTelemetry & Distributed Tracing

**Prerequisites**: Prometheus & Grafana, Loki
**Why here**: Tracing connects metrics and logs across services

**Goal**: Implement end-to-end observability with OpenTelemetry (traces, metrics, logs)

**Learning objectives**:
- Understand OpenTelemetry architecture (SDK, Collector, backends)
- Instrument applications for tracing
- Correlate traces, metrics, and logs
- Analyze distributed system behavior

**Tasks**:
- [ ] Create `experiments/opentelemetry-tutorial/` directory structure
- [ ] Deploy multi-service application (3+ services calling each other)
- [ ] Instrument with OpenTelemetry SDK:
  - [ ] Auto-instrumentation (Go, Python, Node, Java)
  - [ ] Manual span creation
  - [ ] Span attributes and events
  - [ ] Baggage propagation
- [ ] Configure OpenTelemetry Collector:
  - [ ] Receivers (OTLP, Jaeger, Zipkin)
  - [ ] Processors (batch, memory_limiter, attributes)
  - [ ] Exporters (Jaeger, Tempo, Prometheus)
- [ ] Deploy tracing backend:
  - [ ] Jaeger for trace visualization
  - [ ] Or Grafana Tempo (integrates with Grafana)
- [ ] Create trace-aware dashboards:
  - [ ] Service dependency graph
  - [ ] Latency breakdown by span
  - [ ] Error trace exploration
- [ ] Implement exemplars (link metrics to traces)
- [ ] Correlate logs with trace IDs
- [ ] Document sampling strategies (head vs tail sampling)

---

## Phase 3: Application Fundamentals

*Build foundational understanding of web servers and runtimes. You'll use observability skills from Phase 2 to analyze results.*

### Web Server Threading Models: nginx vs Apache

**Prerequisites**: Prometheus & Grafana (for measurement)
**Why here**: Foundational understanding of how web servers handle load

**Goal**: Compare event-driven (nginx) vs process/thread-per-connection (Apache) under load

**Metrics to capture**:
- Requests per second at various concurrency levels
- Memory usage under load
- Latency percentiles (p50, p95, p99)
- Connection handling behavior

**Tasks**:
- [ ] Create `experiments/threading-models/` directory structure
- [ ] Build nginx container with static content
- [ ] Build Apache (prefork MPM) container with same static content
- [ ] Build Apache (worker MPM) container for comparison
- [ ] Create k6 scripts for:
  - Gradual ramp-up (10 → 1000 concurrent users)
  - Sustained high concurrency
  - Spike testing (sudden traffic bursts)
- [ ] Create Argo Workflow to run A/B comparison
- [ ] Configure Prometheus/Grafana dashboards for visualization
- [ ] Document findings in experiment README

---

### Web Server Runtime Comparison: Rust vs Go vs .NET vs Node vs Bun

**Prerequisites**: Threading Models experiment
**Why here**: Builds on web server understanding with modern runtimes

**Goal**: Compare modern web server runtimes for a simple JSON API workload

**Runtimes to test**:
- Rust (Actix-web or Axum)
- Go (net/http or Gin)
- .NET (Kestrel / ASP.NET Core minimal API)
- Node.js (Express or Fastify)
- Bun (native HTTP server)

**Metrics to capture**:
- Requests per second (throughput)
- Latency distribution
- Memory footprint (idle and under load)
- CPU utilization
- Cold start time
- Container image size

**Tasks**:
- [ ] Create `experiments/runtime-comparison/` directory structure
- [ ] Create `images/` entries for each runtime:
  - [ ] `images/api-rust/` - Actix-web or Axum
  - [ ] `images/api-go/` - Go net/http
  - [ ] `images/api-dotnet/` - ASP.NET Core minimal API
  - [ ] `images/api-node/` - Fastify
  - [ ] `images/api-bun/` - Bun native
- [ ] Implement identical API for each:
  - `GET /health` - health check
  - `GET /json` - return static JSON
  - `POST /echo` - echo request body
  - `GET /compute` - CPU-bound work (e.g., fibonacci)
- [ ] Create k6 scripts for:
  - JSON serialization throughput
  - Echo latency
  - CPU-bound workload comparison
- [ ] Create multi-container deployment (all runtimes in same cluster)
- [ ] Build Grafana dashboard comparing all runtimes side-by-side
- [ ] Document methodology and findings

---

### Real-Time Communication: WebSockets vs gRPC

**Prerequisites**: Runtime Comparison experiment
**Why here**: Communication patterns between services - essential for microservices

**Goal**: Compare WebSockets and gRPC for real-time bidirectional communication

**Scenarios to test**:
- High-frequency message streaming (market data simulation)
- Request-response patterns
- Connection management under client churn
- Reconnection behavior

**Metrics to capture**:
- Messages per second
- End-to-end latency
- Connection establishment time
- Memory per connection
- Behavior under network issues (packet loss, latency)

**Tasks**:
- [ ] Create `experiments/realtime-protocols/` directory structure
- [ ] Build WebSocket server (Go or Node)
- [ ] Build gRPC server with streaming (Go)
- [ ] Build gRPC-Web proxy for browser comparison (optional)
- [ ] Create load test clients:
  - [ ] WebSocket client (k6 has WebSocket support)
  - [ ] gRPC client (ghz or custom)
- [ ] Implement test scenarios:
  - [ ] Server push (1 server → N clients)
  - [ ] Bidirectional streaming
  - [ ] Many short-lived connections
- [ ] Configure network policies to simulate latency/packet loss
- [ ] Create comparison dashboard
- [ ] Document protocol trade-offs and findings

---

## Phase 4: Traffic Management

*Learn how to route traffic before learning deployment strategies that depend on it.*

### Gateway API Deep Dive

**Prerequisites**: Basic application deployment knowledge
**Why here**: Traffic routing is fundamental to deployment strategies

**Goal**: Master Kubernetes Gateway API for ingress and traffic management

**Learning objectives**:
- Understand Gateway API resources (Gateway, HTTPRoute, etc.)
- Configure advanced routing
- Implement traffic policies
- Compare with Ingress

**Tasks**:
- [ ] Create `experiments/gateway-api-tutorial/` directory structure
- [ ] Deploy Gateway and HTTPRoutes:
  - [ ] Basic path routing
  - [ ] Host-based routing
  - [ ] Header matching
  - [ ] Query parameter routing
- [ ] Implement traffic management:
  - [ ] Weight-based splitting
  - [ ] Request mirroring
  - [ ] URL rewriting
  - [ ] Header modification
- [ ] Configure advanced features:
  - [ ] Rate limiting
  - [ ] Timeouts and retries
  - [ ] CORS policies
  - [ ] Request authentication
- [ ] Deploy multiple gateways:
  - [ ] Internal vs external
  - [ ] Namespace isolation
  - [ ] Gateway sharing (ReferenceGrant)
- [ ] Implement TLS:
  - [ ] TLS termination
  - [ ] TLS passthrough
  - [ ] mTLS with client certs
- [ ] Document Gateway API vs Ingress migration

---

### Gateway API vs Ingress Performance

**Prerequisites**: Gateway API Deep Dive
**Why here**: Compare ingress solutions now that you understand Gateway API

**Goal**: Compare Kubernetes ingress solutions

**Implementations to test**:
- Gateway API (Envoy Gateway - already deployed)
- Nginx Ingress Controller
- Traefik
- Contour

**Scenarios to test**:
- HTTP routing performance
- TLS termination overhead
- Path-based routing at scale (100+ routes)
- Header manipulation
- Rate limiting

**Tasks**:
- [ ] Create `experiments/ingress-comparison/` directory structure
- [ ] Deploy ingress controllers (one at a time or different namespaces)
- [ ] Create identical route configurations
- [ ] Run throughput and latency tests
- [ ] Compare configuration complexity
- [ ] Test advanced features:
  - [ ] Traffic splitting
  - [ ] Request mirroring
  - [ ] Custom error pages
- [ ] Measure resource consumption

---

## Phase 5: Deployment Strategies

*Progressive complexity: rolling → blue-green → canary → shadow. Each builds on the previous.*

### Rolling Update Optimization

**Prerequisites**: Gateway API, Prometheus
**Why first in deployment**: This is the default Kubernetes deployment method

**Goal**: Optimize Kubernetes default rolling update for minimal disruption

**Parameters to tune**:
- maxSurge and maxUnavailable combinations
- minReadySeconds
- progressDeadlineSeconds
- terminationGracePeriodSeconds
- Readiness probe timing

**Scenarios to test**:
- Fast rollout (aggressive surge)
- Safe rollout (conservative, 1 at a time)
- Large deployments (50+ replicas)
- Slow-starting applications
- Resource-constrained clusters

**Tasks**:
- [ ] Create `experiments/rolling-update/` directory structure
- [ ] Test parameter combinations:
  - [ ] maxSurge=25%, maxUnavailable=25% (default)
  - [ ] maxSurge=100%, maxUnavailable=0 (blue-green-ish)
  - [ ] maxSurge=1, maxUnavailable=0 (safest)
- [ ] Measure impact of readiness probe configuration
- [ ] Test preStop hooks for graceful shutdown
- [ ] Compare with Argo Rollouts rolling strategy
- [ ] Document recommended configurations per use case

---

### Blue-Green Deployment Under Load

**Prerequisites**: Rolling Update Optimization
**Why here**: Simplest advanced deployment strategy (instant cutover)

**Goal**: Measure zero-downtime deployment using blue-green strategy with instant traffic cutover

**Scenarios to test**:
- Instant switchover during steady load
- Switchover during traffic spike
- Rollback timing and impact
- Database connection draining
- Session affinity handling

**Metrics to capture**:
- Request failures during switchover
- Latency spike duration
- Time to complete cutover
- Resource overhead (2x infrastructure)
- Rollback time

**Tasks**:
- [ ] Create `experiments/blue-green-deployment/` directory structure
- [ ] Implement blue-green with:
  - [ ] Kubernetes Services (label selector swap)
  - [ ] Gateway API traffic switching
  - [ ] Argo Rollouts BlueGreen strategy
- [ ] Build test application with health endpoints
- [ ] Create k6 script for continuous load during deployment
- [ ] Test switchover scenarios:
  - [ ] Clean switchover (new version healthy)
  - [ ] Failed health check (should not switch)
  - [ ] Rollback after bad deployment
- [ ] Measure connection draining behavior
- [ ] Test with stateful workloads (database connections)
- [ ] Document infrastructure cost implications

---

### Canary Deployment & Progressive Rollout

**Prerequisites**: Blue-Green Deployment
**Why here**: Adds gradual traffic shifting to blue-green concepts

**Goal**: Measure gradual traffic shifting with automated rollback based on metrics

**Scenarios to test**:
- Weight-based traffic splitting (1% → 10% → 50% → 100%)
- Automatic promotion based on success rate
- Automatic rollback on error threshold
- Header-based routing (internal testing)
- Geographic canary (specific regions first)

**Metrics to capture**:
- Error rate comparison (canary vs stable)
- Latency comparison
- Time to full rollout
- Rollback detection time
- False positive/negative rates for promotion

**Tasks**:
- [ ] Create `experiments/canary-deployment/` directory structure
- [ ] Implement canary with:
  - [ ] Argo Rollouts Canary strategy
  - [ ] Flagger with Gateway API
  - [ ] Istio traffic splitting (if service mesh experiment done)
- [ ] Configure analysis templates:
  - [ ] Success rate (Prometheus query)
  - [ ] Latency p99 threshold
  - [ ] Custom business metrics
- [ ] Create "bad" deployment versions:
  - [ ] Increased latency version
  - [ ] Error-prone version
  - [ ] Memory leak version
- [ ] Test automated rollback triggers
- [ ] Compare manual vs automated promotion
- [ ] Document metric selection best practices

---

### A/B Testing & Feature Flags

**Prerequisites**: Canary Deployment
**Why here**: Variation of canary with user segmentation

**Goal**: Route traffic based on user segments for feature experimentation

**Scenarios to test**:
- Header-based routing (internal users)
- Cookie-based routing (user segments)
- Percentage-based random assignment
- Sticky sessions (consistent user experience)
- Multi-variant testing (A/B/C/n)

**Tasks**:
- [ ] Create `experiments/ab-testing/` directory structure
- [ ] Implement A/B routing with:
  - [ ] Gateway API header matching
  - [ ] Feature flag service (Flagsmith, LaunchDarkly OSS, or custom)
  - [ ] Argo Rollouts experiments
- [ ] Build analytics pipeline for variant comparison
- [ ] Test statistical significance calculation
- [ ] Document experiment design best practices

---

### Shadow Deployment (Dark Launching / Traffic Mirroring)

**Prerequisites**: Canary Deployment
**Why here**: Most advanced deployment pattern - testing without user impact

**Goal**: Test new versions with production traffic without affecting users

**Scenarios to test**:
- Mirror 100% of traffic to shadow
- Mirror subset of traffic (sampling)
- Compare response correctness (diff testing)
- Performance comparison under real load
- Database write handling (shadow should not write)

**Metrics to capture**:
- Shadow vs production response differences
- Shadow version latency
- Resource overhead of mirroring
- Data consistency (read-only shadow)

**Tasks**:
- [ ] Create `experiments/shadow-deployment/` directory structure
- [ ] Implement traffic mirroring with:
  - [ ] Gateway API HTTPRoute mirroring
  - [ ] Istio traffic mirroring
  - [ ] Envoy mirror policy
- [ ] Build shadow-aware application:
  - [ ] Read-only mode detection
  - [ ] Response comparison endpoint
- [ ] Create diff analysis pipeline:
  - [ ] Log shadow responses
  - [ ] Compare with production responses
  - [ ] Alert on differences
- [ ] Test scenarios:
  - [ ] API response format changes
  - [ ] Performance regression detection
  - [ ] New feature validation
- [ ] Handle stateful operations (prevent shadow writes)
- [ ] Document shadow testing patterns

---

### GitOps Deployment Patterns with ArgoCD

**Prerequisites**: All deployment strategies above
**Why here**: Wraps up deployment section - how to manage all strategies via GitOps

**Goal**: Compare ArgoCD sync strategies and their impact during deployments

**Strategies to test**:
- Auto-sync vs manual sync
- Sync waves and hooks
- Progressive sync (ApplicationSet progressive rollout)
- Multi-cluster sync ordering
- Selective sync (resource hooks)

**Tasks**:
- [ ] Create `experiments/gitops-patterns/` directory structure
- [ ] Test sync wave configurations
- [ ] Implement pre-sync and post-sync hooks:
  - [ ] Database migration job
  - [ ] Smoke test job
  - [ ] Notification webhook
- [ ] Test ApplicationSet rolling update across clusters
- [ ] Measure sync performance at scale (100+ resources)
- [ ] Document GitOps deployment patterns

---

### Database Schema Migrations During Deployment

**Prerequisites**: Blue-Green or Canary Deployment
**Why here**: Adds stateful complexity to deployment strategies

**Goal**: Test zero-downtime database migrations with different strategies

**Migration patterns to test**:
- Expand-contract (add column → backfill → remove old)
- Online schema change (pt-online-schema-change, gh-ost)
- Blue-green database (full copy)
- Versioned APIs with schema compatibility

**Scenarios to test**:
- Add nullable column during load
- Add non-nullable column with default
- Rename column (breaking change)
- Add index on large table
- Change column type
- Split table (normalization)

**Metrics to capture**:
- Lock duration and impact
- Replication lag during migration
- Application error rate during migration
- Migration duration vs table size
- Rollback complexity

**Tasks**:
- [ ] Create `experiments/database-migrations/` directory structure
- [ ] Deploy PostgreSQL with test data (10M+ rows)
- [ ] Implement migration strategies:
  - [ ] Plain ALTER TABLE (baseline - causes locks)
  - [ ] Expand-contract pattern
  - [ ] pg_repack for table rewrites
  - [ ] Flyway/Liquibase migration tooling
- [ ] Test application compatibility:
  - [ ] Old app + new schema
  - [ ] New app + old schema
  - [ ] Mixed deployment during migration
- [ ] Create k6 load during migrations
- [ ] Measure lock wait times
- [ ] Test rollback procedures
- [ ] Document migration runbook template

---

## Phase 6: Autoscaling & Resource Management

*Learn to scale applications efficiently before optimizing data and messaging layers.*

### Horizontal Pod Autoscaler Tuning

**Prerequisites**: Prometheus (for metrics-based scaling)
**Why first in scaling**: HPA is the foundation for Kubernetes autoscaling

**Goal**: Find optimal HPA configurations for different workload patterns

**Scenarios to test**:
- CPU-bound workloads
- Memory-bound workloads
- Custom metrics (requests per second, queue depth)
- Sudden traffic spikes
- Gradual ramp-up

**Metrics to capture**:
- Time to scale up/down
- Resource utilization during scaling
- Request latency during scale events
- Over-provisioning waste
- Under-provisioning failures

**Tasks**:
- [ ] Create `experiments/hpa-tuning/` directory structure
- [ ] Build test application with configurable resource usage
- [ ] Test HPA configurations:
  - [ ] CPU threshold variations (50%, 70%, 80%)
  - [ ] Scale up/down stabilization windows
  - [ ] Min/max replica bounds
  - [ ] Multiple metrics (CPU + custom)
- [ ] Integrate KEDA for event-driven scaling comparison
- [ ] Test with Vertical Pod Autoscaler (VPA)
- [ ] Create decision tree for HPA configuration

---

### KEDA Event-Driven Autoscaling

**Prerequisites**: HPA Tuning
**Why here**: Builds on HPA with external event sources

**Goal**: Master scaling based on external metrics and events

**Learning objectives**:
- Understand KEDA architecture (ScaledObject, ScaledJob, triggers)
- Configure various event sources
- Tune scaling parameters
- Compare with HPA

**Tasks**:
- [ ] Create `experiments/keda-tutorial/` directory structure
- [ ] Deploy KEDA scalers tutorial:
  - [ ] Prometheus scaler (scale on custom metrics)
  - [ ] RabbitMQ scaler (scale on queue depth)
  - [ ] Kafka scaler (scale on consumer lag)
  - [ ] Cron scaler (scheduled scaling)
  - [ ] Azure Service Bus / AWS SQS (via Crossplane queues)
- [ ] Build worker application that processes queue messages
- [ ] Configure ScaledObject:
  - [ ] Triggers and thresholds
  - [ ] Cooldown periods
  - [ ] Min/max replicas
  - [ ] Fallback behavior
- [ ] Test ScaledJob for batch processing:
  - [ ] Job per message
  - [ ] Parallel job execution
  - [ ] Completion handling
- [ ] Compare KEDA vs HPA:
  - [ ] Scale-to-zero capability
  - [ ] External metric support
  - [ ] Scaling responsiveness
- [ ] Document scaler selection guide

---

### Spot/Preemptible VMs with KEDA

**Prerequisites**: KEDA Event-Driven Autoscaling
**Why here**: Advanced scaling with cost optimization

**Goal**: Optimize cloud costs using spot instances with graceful handling

**Learning objectives**:
- Understand spot/preemptible VM characteristics
- Handle interruption gracefully
- Combine with KEDA for cost-effective scaling
- Implement fallback strategies

**Tasks**:
- [ ] Create `experiments/spot-instances/` directory structure
- [ ] Configure spot node pools:
  - [ ] Azure: Spot VMs with AKS
  - [ ] AWS: Spot Instances with EKS
- [ ] Implement interruption handling:
  - [ ] Node termination handler (AWS/Azure)
  - [ ] Pod disruption budgets
  - [ ] Graceful shutdown with preStop hooks
- [ ] Deploy workload across spot and on-demand:
  - [ ] Tolerations and node affinity
  - [ ] Priority classes for critical workloads
  - [ ] Pod topology spread constraints
- [ ] Configure KEDA with spot awareness:
  - [ ] Scale on queue depth to spot nodes
  - [ ] Fallback to on-demand when spots unavailable
  - [ ] Scheduled scaling for predictable workloads
- [ ] Measure cost savings:
  - [ ] Track spot vs on-demand hours
  - [ ] Calculate effective discount
  - [ ] Document interruption frequency
- [ ] Test failure scenarios:
  - [ ] Spot interruption during processing
  - [ ] All spots reclaimed simultaneously
  - [ ] Recovery time measurement

---

## Phase 7: Data & Storage

*Now that you can deploy and scale, learn to manage stateful workloads.*

### Database Performance: PostgreSQL vs MySQL vs MariaDB

**Prerequisites**: Prometheus, basic SQL knowledge
**Why first in data**: Databases are the foundation of stateful applications

**Goal**: Compare relational databases for OLTP workloads

**Workloads to test**:
- Simple key-value lookups
- Complex joins across multiple tables
- Write-heavy (INSERT/UPDATE)
- Mixed read/write (80/20)

**Metrics to capture**:
- Queries per second
- Latency percentiles
- Connection pool efficiency
- Memory and CPU usage
- WAL/binlog write performance

**Tasks**:
- [ ] Create `experiments/database-comparison/` directory structure
- [ ] Use Crossplane Database claims for managed instances (Azure/AWS)
- [ ] Alternatively, deploy containerized versions for local testing:
  - [ ] PostgreSQL 16
  - [ ] MySQL 8.0
  - [ ] MariaDB 11
- [ ] Create schema with realistic tables (users, orders, products)
- [ ] Generate test data (1M+ rows)
- [ ] Build k6 scripts using xk6-sql extension
- [ ] Create pgbench/sysbench comparison baseline
- [ ] Test with connection pooling (PgBouncer, ProxySQL)
- [ ] Document query optimizer differences

---

### Caching Layer: Redis vs Memcached vs Dragonfly

**Prerequisites**: Database Performance experiment
**Why here**: Caching complements databases for read-heavy workloads

**Goal**: Compare in-memory caching solutions for high-throughput scenarios

**Scenarios to test**:
- Simple GET/SET operations
- Pipelined batch operations
- Pub/Sub messaging throughput
- Large value handling (1KB, 10KB, 100KB)
- Cluster mode vs single instance

**Metrics to capture**:
- Operations per second
- Latency (p50, p95, p99)
- Memory efficiency
- CPU utilization
- Eviction behavior under memory pressure

**Tasks**:
- [ ] Create `experiments/caching-comparison/` directory structure
- [ ] Deploy cache instances:
  - [ ] Redis 7 (standalone and cluster)
  - [ ] Memcached 1.6
  - [ ] Dragonfly (Redis-compatible, multi-threaded)
- [ ] Use Crossplane Cache claims for managed versions
- [ ] Build benchmark client using redis-benchmark and custom k6 scripts
- [ ] Test serialization formats (JSON, MessagePack, Protobuf)
- [ ] Measure memory fragmentation over time
- [ ] Compare persistence options (Redis AOF/RDB vs Dragonfly snapshots)
- [ ] Document use case recommendations

---

### Object Storage: S3 vs MinIO vs Azure Blob

**Prerequisites**: Basic storage understanding
**Why here**: Object storage for unstructured data complements databases

**Goal**: Compare object storage throughput and latency

**Scenarios to test**:
- Small file uploads (1KB - 100KB)
- Large file uploads (100MB - 1GB) with multipart
- Concurrent downloads
- Listing operations on large buckets
- Presigned URL generation

**Metrics to capture**:
- Upload/download throughput (MB/s)
- Time to first byte
- Operations per second for small objects
- Multipart upload efficiency
- S3 API compatibility differences

**Tasks**:
- [ ] Create `experiments/object-storage/` directory structure
- [ ] Use Crossplane ObjectStorage claims for cloud buckets
- [ ] Deploy MinIO for local/self-hosted comparison
- [ ] Build test client with AWS SDK
- [ ] Generate test datasets (small files, large files, mixed)
- [ ] Test with various concurrency levels
- [ ] Measure cost per operation (cloud only)
- [ ] Compare lifecycle policy behavior

---

## Phase 8: Messaging & Event Streaming

*Asynchronous communication patterns - builds on database knowledge for event-driven architectures.*

### Message Queues: Kafka vs RabbitMQ vs NATS

**Prerequisites**: Database concepts, basic messaging understanding
**Why first in messaging**: Core message broker comparison

**Goal**: Compare message brokers for event streaming and task queues

**Scenarios to test**:
- High-throughput event streaming (1M+ msgs/sec)
- Fan-out (1 producer, N consumers)
- Request-reply patterns
- Message ordering guarantees
- Consumer group rebalancing
- Exactly-once semantics

**Metrics to capture**:
- Messages per second (produce and consume)
- End-to-end latency
- Broker CPU and memory usage
- Disk I/O (for persistent messages)
- Recovery time after broker failure

**Tasks**:
- [ ] Create `experiments/message-queues/` directory structure
- [ ] Deploy message brokers:
  - [ ] Kafka (via Strimzi operator - already in argocd-apps)
  - [ ] RabbitMQ (via RabbitMQ operator - already in argocd-apps)
  - [ ] NATS with JetStream
- [ ] Build producer/consumer applications in Go
- [ ] Create load test scenarios:
  - [ ] Sustained throughput
  - [ ] Burst traffic
  - [ ] Slow consumer simulation
- [ ] Test durability (kill broker during writes)
- [ ] Compare operational complexity
- [ ] Document when to use each

---

### Cloud Message Services: SQS vs Azure Service Bus

**Prerequisites**: Message Queues experiment
**Why here**: Cloud-managed messaging via Crossplane

**Goal**: Compare managed cloud messaging via Crossplane

**Scenarios to test**:
- Standard queue throughput
- FIFO queue ordering
- Dead letter queue behavior
- Long polling efficiency
- Batch operations

**Tasks**:
- [ ] Create `experiments/cloud-messaging/` directory structure
- [ ] Use Crossplane Queue claims (demonstrates cloud abstraction!)
- [ ] Deploy same producer/consumer to both clouds
- [ ] Compare:
  - [ ] Message visibility handling
  - [ ] Retry policies
  - [ ] DLQ processing
  - [ ] Pricing per message
- [ ] Document cloud-specific gotchas

---

## Phase 9: Security & Infrastructure

*Secure your applications and automate certificate management.*

### cert-manager & TLS Automation

**Prerequisites**: Gateway API, basic PKI understanding
**Why first in security**: TLS is fundamental to secure communication

**Goal**: Automate TLS certificate management in Kubernetes

**Learning objectives**:
- Understand PKI and certificate lifecycle
- Configure cert-manager issuers
- Automate certificate renewal
- Secure ingress with TLS

**Tasks**:
- [ ] Create `experiments/cert-manager-tutorial/` directory structure
- [ ] Configure issuers:
  - [ ] Self-signed (development)
  - [ ] Let's Encrypt (ACME - staging and prod)
  - [ ] Private CA (internal services)
  - [ ] Vault PKI backend
- [ ] Create certificates:
  - [ ] Ingress/Gateway TLS
  - [ ] Service-to-service mTLS
  - [ ] Wildcard certificates
- [ ] Configure automatic renewal:
  - [ ] Renewal thresholds
  - [ ] Certificate monitoring
  - [ ] Alerting on expiry
- [ ] Implement ACME challenges:
  - [ ] HTTP-01 challenge
  - [ ] DNS-01 challenge (Azure DNS, Route53)
- [ ] Test failure scenarios:
  - [ ] Issuer unavailable
  - [ ] Challenge failures
  - [ ] Certificate expiry handling
- [ ] Document certificate lifecycle management

---

### HashiCorp Vault Secrets Management

**Prerequisites**: cert-manager, basic security concepts
**Why here**: Builds on TLS with comprehensive secrets management

**Goal**: Master secrets management with Vault in Kubernetes

**Learning objectives**:
- Understand Vault architecture and auth methods
- Inject secrets into pods
- Implement dynamic secrets
- Configure secret rotation

**Tasks**:
- [ ] Create `experiments/vault-tutorial/` directory structure
- [ ] Configure Vault auth methods:
  - [ ] Kubernetes auth (ServiceAccount)
  - [ ] AppRole for applications
  - [ ] OIDC for users
- [ ] Implement secret injection patterns:
  - [ ] Vault Agent Sidecar
  - [ ] Vault CSI Provider
  - [ ] External Secrets Operator
- [ ] Configure dynamic secrets:
  - [ ] Database credentials (PostgreSQL)
  - [ ] AWS IAM credentials
  - [ ] PKI certificates
- [ ] Build secret rotation workflow:
  - [ ] Automatic credential rotation
  - [ ] Application restart strategies
  - [ ] Zero-downtime rotation
- [ ] Implement policies:
  - [ ] Path-based access control
  - [ ] Namespace isolation
  - [ ] Audit logging
- [ ] Test disaster recovery:
  - [ ] Seal/unseal procedures
  - [ ] Backup and restore
  - [ ] HA failover
- [ ] Document secrets management patterns

---

## Phase 10: Advanced Platform

*Advanced topics that tie together everything learned so far.*

### Argo Workflows Orchestration

**Prerequisites**: ArgoCD, GitOps patterns
**Why first in advanced**: Complex workflow orchestration builds on deployment knowledge

**Goal**: Master workflow orchestration for CI/CD, data pipelines, and batch processing

**Learning objectives**:
- Understand Argo Workflows concepts (DAG, steps, artifacts)
- Build complex multi-step workflows
- Handle failures and retries
- Integrate with other systems

**Tasks**:
- [ ] Create `experiments/argo-workflows-tutorial/` directory structure
- [ ] Build workflow patterns:
  - [ ] Sequential steps
  - [ ] Parallel execution
  - [ ] DAG dependencies
  - [ ] Conditional execution
  - [ ] Loops and recursion
- [ ] Implement artifact passing:
  - [ ] S3/MinIO artifact storage
  - [ ] Artifact between steps
  - [ ] Output parameters
- [ ] Configure workflow templates:
  - [ ] WorkflowTemplate for reuse
  - [ ] ClusterWorkflowTemplate
  - [ ] Template composition
- [ ] Build practical workflows:
  - [ ] CI pipeline (build → test → deploy)
  - [ ] Data processing pipeline
  - [ ] ML training workflow
- [ ] Handle failures:
  - [ ] Retry strategies
  - [ ] Timeout configuration
  - [ ] Exit handlers
  - [ ] Workflow-level error handling
- [ ] Integrate with events:
  - [ ] Argo Events triggers
  - [ ] Webhook triggers
  - [ ] Scheduled workflows (CronWorkflow)
- [ ] Document workflow design patterns

---

### Service Mesh Comparison: Istio vs Linkerd vs Cilium

**Prerequisites**: Gateway API, OpenTelemetry, deployment strategies
**Why here**: Service mesh ties together networking, observability, and traffic management

**Goal**: Compare service mesh implementations for observability and traffic management

**Features to test**:
- mTLS performance overhead
- Traffic splitting (canary deployments)
- Circuit breaking behavior
- Retry policies
- Observability (traces, metrics)
- Resource consumption (sidecars vs eBPF)

**Metrics to capture**:
- Latency overhead (with/without mesh)
- Memory per pod (sidecar cost)
- CPU overhead
- Control plane resource usage
- Time to propagate config changes

**Tasks**:
- [ ] Create `experiments/service-mesh/` directory structure
- [ ] Deploy baseline application without mesh
- [ ] Install and configure:
  - [ ] Istio (sidecar-based)
  - [ ] Linkerd (lightweight sidecar)
  - [ ] Cilium service mesh (eBPF, sidecar-free)
- [ ] Run identical load tests on each
- [ ] Test traffic management features:
  - [ ] A/B testing
  - [ ] Canary releases
  - [ ] Fault injection
- [ ] Compare observability integrations with existing Prometheus/Grafana
- [ ] Document operational complexity

---

## Phase 11: Chaos Engineering & Resilience

*The final test - validate everything works under failure conditions.*

### Pod Failure Recovery

**Prerequisites**: All deployment strategies, observability
**Why first in chaos**: Simplest chaos scenario - pod-level failures

**Goal**: Measure application resilience to pod failures

**Scenarios to test**:
- Single pod termination
- Multiple pod termination (50% of replicas)
- All pods termination
- OOMKill simulation
- Liveness probe failures

**Metrics to capture**:
- Time to detect failure
- Time to reschedule
- Time to ready
- Request error rate during failure
- Recovery time to full capacity

**Tasks**:
- [ ] Create `experiments/chaos-pod-failure/` directory structure
- [ ] Deploy Chaos Mesh (already in argocd-apps)
- [ ] Create PodChaos experiments:
  - [ ] Pod kill
  - [ ] Pod failure (container crash)
  - [ ] Container kill
- [ ] Instrument application with detailed metrics
- [ ] Run chaos during load test
- [ ] Test with different:
  - [ ] Replica counts
  - [ ] Pod disruption budgets
  - [ ] Readiness probe configurations
- [ ] Document SLO impact

---

### Node Drain Impact

**Prerequisites**: Pod Failure Recovery
**Why here**: Infrastructure-level failures - more complex than pod failures

**Goal**: Measure impact of node maintenance on application availability

**Scenarios to test**:
- Graceful node drain
- Sudden node failure
- Multiple node failure
- Zone failure simulation

**Tasks**:
- [ ] Create `experiments/chaos-node-drain/` directory structure
- [ ] Deploy multi-replica application across nodes
- [ ] Configure pod anti-affinity rules
- [ ] Test drain scenarios:
  - [ ] `kubectl drain` with grace period
  - [ ] Sudden node shutdown (VM stop)
  - [ ] Cordon + eviction
- [ ] Measure:
  - [ ] Request failures during drain
  - [ ] Time to rebalance workloads
  - [ ] PVC reattachment time (if applicable)
- [ ] Test with pod disruption budgets
- [ ] Document node pool sizing recommendations

---

### Network Partition & Latency Injection

**Prerequisites**: Node Drain, Service Mesh (recommended)
**Why last**: Most complex chaos - network-level failures

**Goal**: Test application behavior during network issues

**Scenarios to test**:
- Latency injection (50ms, 100ms, 500ms)
- Packet loss (1%, 5%, 20%)
- Bandwidth throttling
- DNS failures
- Partition between services

**Tasks**:
- [ ] Create `experiments/chaos-network/` directory structure
- [ ] Use Chaos Mesh NetworkChaos:
  - [ ] Delay
  - [ ] Loss
  - [ ] Duplicate
  - [ ] Corrupt
  - [ ] Partition
- [ ] Test microservice communication patterns
- [ ] Measure:
  - [ ] Timeout behavior
  - [ ] Retry storms
  - [ ] Circuit breaker activation
  - [ ] Fallback behavior
- [ ] Test database connection resilience
- [ ] Document timeout tuning recommendations

---

## Learning Path Summary

| Phase | Focus | Key Skills Gained |
|-------|-------|-------------------|
| 1 | Platform Setup | Spacelift, Crossplane, multi-cloud |
| 2 | Observability | Prometheus, Loki, OpenTelemetry, metrics/logs/traces |
| 3 | App Fundamentals | Web servers, runtimes, protocols |
| 4 | Traffic Management | Gateway API, routing, load balancing |
| 5 | Deployment Strategies | Rolling, blue-green, canary, shadow, GitOps |
| 6 | Autoscaling | HPA, KEDA, spot instances, cost optimization |
| 7 | Data & Storage | Databases, caching, object storage |
| 8 | Messaging | Kafka, RabbitMQ, NATS, cloud queues |
| 9 | Security | TLS, cert-manager, Vault, secrets |
| 10 | Advanced Platform | Argo Workflows, service mesh |
| 11 | Chaos Engineering | Pod/node/network failures, resilience |

---

## Notes

- All experiments should follow the `experiments/_template/` structure
- Use Crossplane claims for cloud resources when applicable
- Prefer Spacelift for cloud deployments, Taskfile for local minikube
- Document cost estimates for cloud experiments
- Always include cleanup instructions
- Each phase builds on previous phases - follow the order for best learning progression
