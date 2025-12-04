# illm-k8s-lab TODO

## Overview

A learning-focused experiment roadmap for mastering Kubernetes ecosystem tools. Each experiment is tutorial-style with hands-on exercises. Benchmarks and comparisons come later after fundamentals are solid.

**Target:** ~35-40 experiments across 11 phases
**Environment:** Kind (multi-cluster), Spacelift, Crossplane, Azure/AWS
**Focus:** Portfolio-ready, demonstrable experiments

---

## Phase 1: Platform Bootstrap

*Get the multi-cloud GitOps foundation running. This unlocks everything else.*

### 1.1 Spacelift + Kind Local Setup

**Goal:** Establish local Kind cluster with Spacelift managing remote state

**Learning objectives:**
- Understand Spacelift stacks, contexts, and policies
- Configure Kind for multi-cluster experiments
- Set up GitOps workflow for infrastructure changes

**Tasks:**
- [ ] Create Spacelift account and connect GitHub repo
- [ ] Create `spacelift-root` stack (administrative=true)
- [ ] Configure cloud credential contexts (Azure, AWS)
- [ ] Set up Kind cluster with sufficient resources
- [ ] Verify Spacelift can plan/apply to local state
- [ ] Document Spacelift workflow patterns

---

### 1.2 Crossplane Fundamentals

**Goal:** Master Crossplane for cloud resource provisioning

**Learning objectives:**
- Understand Crossplane architecture (providers, XRDs, compositions)
- Create and use Composite Resource Definitions
- Build reusable compositions for common patterns

**Tasks:**
- [ ] Deploy Crossplane to Kind cluster
- [ ] Install AWS and Azure providers
- [ ] Create first XRD: SimpleDatabase (abstracts RDS/Azure SQL)
- [ ] Create first XRD: SimpleBucket (abstracts S3/Azure Blob)
- [ ] Create first XRD: SimpleQueue (abstracts SQS/Service Bus)
- [ ] Test claims provision real cloud resources
- [ ] Document XRD authoring patterns

---

## Phase 2: Security Foundations

*Security first - TLS, certificates, and secrets management are prerequisites for everything else.*

### 2.1 cert-manager & TLS Automation

**Goal:** Automate TLS certificate lifecycle in Kubernetes

**Learning objectives:**
- Understand PKI fundamentals and certificate lifecycle
- Configure cert-manager issuers (self-signed, ACME, private CA)
- Automate certificate renewal and monitoring

**Tasks:**
- [ ] Create `experiments/cert-manager-tutorial/`
- [ ] Deploy cert-manager via ArgoCD
- [ ] Configure Issuers:
  - [ ] SelfSigned (for development)
  - [ ] Let's Encrypt staging (ACME HTTP-01)
  - [ ] Let's Encrypt production
  - [ ] Private CA (for internal mTLS)
- [ ] Create Certificate resources:
  - [ ] Ingress/Gateway TLS termination
  - [ ] Wildcard certificate
  - [ ] Short-lived certificate (test renewal)
- [ ] Implement DNS-01 challenge (Azure DNS or Route53 via Crossplane)
- [ ] Set up certificate expiry alerting
- [ ] Test failure scenarios (issuer down, challenge failure)
- [ ] Document certificate patterns for different use cases

---

### 2.2 HashiCorp Vault Secrets Management

**Goal:** Centralized secrets management with dynamic credentials

**Learning objectives:**
- Understand Vault architecture (secrets engines, auth methods, policies)
- Inject secrets into Kubernetes pods
- Implement dynamic database credentials

**Tasks:**
- [ ] Create `experiments/vault-tutorial/`
- [ ] Deploy Vault (dev mode first, then HA)
- [ ] Configure auth methods:
  - [ ] Kubernetes auth (ServiceAccount JWT)
  - [ ] AppRole (for CI/CD)
- [ ] Set up secrets engines:
  - [ ] KV v2 (static secrets)
  - [ ] Database (dynamic PostgreSQL creds)
  - [ ] PKI (dynamic certificates)
- [ ] Implement secret injection:
  - [ ] Vault Agent Sidecar
  - [ ] Vault CSI Provider
  - [ ] External Secrets Operator
- [ ] Create policies for namespace isolation
- [ ] Test secret rotation workflow
- [ ] Implement audit logging
- [ ] Document secrets management patterns

---

### 2.3 Network Policies & Pod Security

**Goal:** Implement defense-in-depth with network segmentation and pod security

**Learning objectives:**
- Write effective NetworkPolicy resources
- Understand Pod Security Standards (PSS)
- Implement least-privilege pod configurations

**Tasks:**
- [ ] Create `experiments/network-security-tutorial/`
- [ ] Deploy Calico or Cilium CNI (for NetworkPolicy support)
- [ ] Implement NetworkPolicy patterns:
  - [ ] Default deny all ingress/egress
  - [ ] Allow specific service-to-service communication
  - [ ] Allow egress to specific external CIDRs
  - [ ] Namespace isolation
- [ ] Configure Pod Security:
  - [ ] Pod Security Admission (PSA) labels
  - [ ] Restricted security context
  - [ ] Read-only root filesystem
  - [ ] Non-root containers
- [ ] Test policy enforcement (verify blocked traffic)
- [ ] Document security baseline for all experiments

---

## Phase 3: Observability Stack

*You need to see what's happening before you can improve it. These skills are used in every subsequent experiment.*

### 3.1 Prometheus & Grafana Deep Dive

**Goal:** Master metrics collection, PromQL, alerting, and dashboards

**Learning objectives:**
- Understand Prometheus architecture (scraping, TSDB, federation)
- Write effective PromQL queries
- Build actionable Grafana dashboards
- Configure alerting pipelines

**Tasks:**
- [ ] Create `experiments/prometheus-tutorial/`
- [ ] Deploy kube-prometheus-stack via ArgoCD
- [ ] Build sample app with custom metrics:
  - [ ] Counter (http_requests_total)
  - [ ] Gauge (active_connections)
  - [ ] Histogram (request_duration_seconds)
  - [ ] Summary (response_size_bytes)
- [ ] Create ServiceMonitor for scrape discovery
- [ ] Write PromQL tutorial queries:
  - [ ] rate() and irate() for counters
  - [ ] Aggregations (sum, avg, max by labels)
  - [ ] histogram_quantile() for percentiles
  - [ ] absent() for missing metric alerts
  - [ ] predict_linear() for capacity planning
- [ ] Build Grafana dashboards:
  - [ ] RED metrics (Rate, Errors, Duration)
  - [ ] USE metrics (Utilization, Saturation, Errors)
  - [ ] Dashboard variables and templating
- [ ] Configure alerting:
  - [ ] PrometheusRule CRDs
  - [ ] Alertmanager routing and silences
  - [ ] Alert grouping and inhibition
- [ ] Document PromQL patterns and anti-patterns

---

### 3.2 Loki & Log Aggregation

**Goal:** Centralized logging with Loki and LogQL

**Learning objectives:**
- Understand Loki's label-based architecture (vs full-text indexing)
- Write effective LogQL queries
- Correlate logs with metrics in Grafana

**Tasks:**
- [ ] Create `experiments/loki-tutorial/`
- [ ] Deploy Loki stack (Loki + Promtail)
- [ ] Build app with structured JSON logging
- [ ] Configure Promtail pipelines:
  - [ ] Label extraction (namespace, pod, container)
  - [ ] JSON field parsing
  - [ ] Regex extraction
  - [ ] Drop/keep filtering
  - [ ] Multiline log handling
- [ ] Write LogQL tutorial:
  - [ ] Label matchers and line filters
  - [ ] Parser expressions (json, pattern, regexp)
  - [ ] Metric queries (rate, count_over_time)
  - [ ] Unwrap for numeric fields
- [ ] Build log dashboards:
  - [ ] Log panel with live tail
  - [ ] Log volume over time
  - [ ] Error log filtering
- [ ] Set up log-based alerts (error rate threshold)
- [ ] Correlate logs ↔ metrics in Grafana (split view)
- [ ] Document logging best practices

---

### 3.3 OpenTelemetry & Distributed Tracing

**Goal:** End-to-end observability with traces, connecting metrics and logs

**Learning objectives:**
- Understand OpenTelemetry architecture (SDK, Collector, backends)
- Instrument applications for distributed tracing
- Correlate traces ↔ metrics ↔ logs

**Tasks:**
- [ ] Create `experiments/opentelemetry-tutorial/`
- [ ] Deploy OpenTelemetry Collector
- [ ] Deploy Tempo or Jaeger as trace backend
- [ ] Build multi-service demo app (3+ services):
  - [ ] Service A → Service B → Service C
  - [ ] Each service instrumented with OTel SDK
- [ ] Implement tracing:
  - [ ] Auto-instrumentation (HTTP, gRPC, DB)
  - [ ] Manual span creation
  - [ ] Span attributes and events
  - [ ] Context propagation (W3C Trace Context)
- [ ] Configure Collector:
  - [ ] OTLP receiver
  - [ ] Batch processor
  - [ ] Exporters (Tempo/Jaeger, Prometheus)
- [ ] Connect the three pillars:
  - [ ] Exemplars (metrics → traces)
  - [ ] Trace ID in logs (logs → traces)
  - [ ] Service graph from traces
- [ ] Build trace-aware dashboards:
  - [ ] Service dependency map
  - [ ] Latency breakdown by span
  - [ ] Error trace exploration
- [ ] Document sampling strategies (head vs tail)

---

## Phase 4: Traffic Management

*Control how traffic flows before learning deployment strategies that depend on it.*

### 4.1 Gateway API Deep Dive

**Goal:** Master Kubernetes Gateway API for ingress and traffic routing

**Learning objectives:**
- Understand Gateway API resources (Gateway, HTTPRoute, GRPCRoute)
- Implement advanced routing patterns
- Compare with legacy Ingress

**Tasks:**
- [ ] Create `experiments/gateway-api-tutorial/`
- [ ] Deploy Envoy Gateway (or Cilium Gateway)
- [ ] Configure Gateway resource
- [ ] Implement HTTPRoute patterns:
  - [ ] Path-based routing
  - [ ] Host-based routing (virtual hosts)
  - [ ] Header matching
  - [ ] Query parameter routing
  - [ ] Method matching (GET vs POST)
- [ ] Traffic manipulation:
  - [ ] Weight-based splitting (A/B)
  - [ ] Request mirroring
  - [ ] URL rewriting
  - [ ] Header modification (add/remove/set)
  - [ ] Redirects
- [ ] Advanced features:
  - [ ] Timeouts and retries
  - [ ] Rate limiting (via policy attachment)
  - [ ] CORS configuration
- [ ] TLS configuration:
  - [ ] TLS termination (with cert-manager certs)
  - [ ] TLS passthrough
  - [ ] mTLS with client certificates
- [ ] Multi-gateway setup:
  - [ ] Internal vs external gateways
  - [ ] Namespace isolation (ReferenceGrant)
- [ ] Document Gateway API vs Ingress migration

---

### 4.2 Ingress Controllers Comparison

**Goal:** Understand trade-offs between ingress implementations

**Learning objectives:**
- Compare nginx, Traefik, and Envoy-based controllers
- Understand feature/performance trade-offs
- Make informed controller selection

**Tasks:**
- [ ] Create `experiments/ingress-comparison/`
- [ ] Deploy and configure:
  - [ ] Nginx Ingress Controller
  - [ ] Traefik
  - [ ] Envoy Gateway (Gateway API)
- [ ] Implement equivalent routing on each
- [ ] Compare:
  - [ ] Configuration complexity
  - [ ] Feature availability
  - [ ] Resource consumption
  - [ ] Custom resource patterns
- [ ] Test advanced features:
  - [ ] Rate limiting implementation
  - [ ] Authentication integration
  - [ ] Custom error pages
- [ ] Document selection criteria

---

## Phase 5: Service Mesh

*Service mesh builds on traffic management with mTLS, observability, and advanced traffic control.*

### 5.1 Istio Deep Dive

**Goal:** Master Istio service mesh fundamentals

**Learning objectives:**
- Understand Istio architecture (control plane, data plane, sidecars)
- Configure traffic management policies
- Implement security with mTLS

**Tasks:**
- [ ] Create `experiments/istio-tutorial/`
- [ ] Install Istio (istioctl or Helm)
- [ ] Enable sidecar injection (namespace label)
- [ ] Deploy sample microservices app
- [ ] Traffic management:
  - [ ] VirtualService routing rules
  - [ ] DestinationRule load balancing
  - [ ] Traffic splitting (canary)
  - [ ] Fault injection (delays, aborts)
  - [ ] Circuit breaking
  - [ ] Retries and timeouts
- [ ] Security:
  - [ ] Automatic mTLS (PeerAuthentication)
  - [ ] Authorization policies (allow/deny)
  - [ ] JWT validation (RequestAuthentication)
- [ ] Observability:
  - [ ] Kiali service graph
  - [ ] Distributed tracing (Jaeger integration)
  - [ ] Metrics (Prometheus integration)
- [ ] Gateway:
  - [ ] Istio Gateway (vs Gateway API)
  - [ ] External traffic management
- [ ] Document Istio patterns and gotchas

---

### 5.2 Linkerd Tutorial

**Goal:** Learn lightweight service mesh alternative

**Learning objectives:**
- Understand Linkerd architecture (simpler than Istio)
- Compare operational complexity
- Evaluate for different use cases

**Tasks:**
- [ ] Create `experiments/linkerd-tutorial/`
- [ ] Install Linkerd (CLI + control plane)
- [ ] Inject proxies into workloads
- [ ] Deploy same sample app as Istio experiment
- [ ] Configure:
  - [ ] Automatic mTLS
  - [ ] Traffic splitting (TrafficSplit CRD)
  - [ ] Retries and timeouts (ServiceProfile)
  - [ ] Authorization policies
- [ ] Observability:
  - [ ] Linkerd dashboard
  - [ ] Tap for live traffic inspection
  - [ ] Metrics and golden signals
- [ ] Compare with Istio:
  - [ ] Resource consumption
  - [ ] Configuration complexity
  - [ ] Feature coverage
- [ ] Document when to choose Linkerd vs Istio

---

### 5.3 Cilium Service Mesh (eBPF)

**Goal:** Explore sidecar-free service mesh with eBPF

**Learning objectives:**
- Understand eBPF-based networking
- Compare sidecar vs sidecar-free architectures
- Evaluate Cilium for CNI + service mesh

**Tasks:**
- [ ] Create `experiments/cilium-tutorial/`
- [ ] Install Cilium as CNI with service mesh features
- [ ] Deploy sample app (no sidecars needed)
- [ ] Configure:
  - [ ] L7 traffic policies (CiliumNetworkPolicy)
  - [ ] mTLS (Cilium encryption)
  - [ ] Load balancing
  - [ ] Ingress (Cilium Ingress or Gateway API)
- [ ] Observability:
  - [ ] Hubble for network visibility
  - [ ] Hubble UI
  - [ ] Prometheus metrics
- [ ] Compare with sidecar meshes:
  - [ ] Performance overhead
  - [ ] Resource consumption
  - [ ] Operational complexity
- [ ] Document eBPF advantages and limitations

---

## Phase 6: Messaging & Event Streaming

*Asynchronous communication patterns for event-driven architectures.*

### 6.1 Kafka with Strimzi

**Goal:** Deploy and operate Kafka on Kubernetes

**Learning objectives:**
- Understand Kafka architecture (brokers, topics, partitions, consumers)
- Use Strimzi operator for Kafka lifecycle
- Implement common messaging patterns

**Tasks:**
- [ ] Create `experiments/kafka-tutorial/`
- [ ] Deploy Strimzi operator via ArgoCD
- [ ] Create Kafka cluster (KafkaCluster CRD)
- [ ] Configure:
  - [ ] Topics (KafkaTopic CRD)
  - [ ] Users and ACLs (KafkaUser CRD)
  - [ ] Replication and partitions
- [ ] Build producer/consumer apps:
  - [ ] Simple pub/sub
  - [ ] Consumer groups
  - [ ] Exactly-once semantics
- [ ] Implement patterns:
  - [ ] Event sourcing
  - [ ] CQRS with Kafka
  - [ ] Dead letter queue
- [ ] Monitoring:
  - [ ] Kafka metrics in Prometheus
  - [ ] Consumer lag monitoring
  - [ ] Grafana dashboards
- [ ] Connect (optional):
  - [ ] Kafka Connect for integrations
  - [ ] Source/sink connectors
- [ ] Document Kafka operational patterns

---

### 6.2 RabbitMQ with Operator

**Goal:** Deploy and operate RabbitMQ for task queues

**Learning objectives:**
- Understand RabbitMQ architecture (exchanges, queues, bindings)
- Use RabbitMQ Cluster Operator
- Compare with Kafka use cases

**Tasks:**
- [ ] Create `experiments/rabbitmq-tutorial/`
- [ ] Deploy RabbitMQ Cluster Operator
- [ ] Create RabbitMQ cluster (RabbitmqCluster CRD)
- [ ] Configure:
  - [ ] Exchanges (direct, fanout, topic, headers)
  - [ ] Queues and bindings
  - [ ] Users and permissions
- [ ] Build producer/consumer apps:
  - [ ] Work queues (competing consumers)
  - [ ] Pub/sub (fanout)
  - [ ] Routing (topic exchange)
  - [ ] RPC pattern
- [ ] Implement reliability:
  - [ ] Publisher confirms
  - [ ] Consumer acknowledgments
  - [ ] Dead letter exchanges
  - [ ] Message TTL
- [ ] Monitoring:
  - [ ] RabbitMQ management UI
  - [ ] Prometheus metrics
  - [ ] Queue depth alerting
- [ ] Document RabbitMQ vs Kafka decision guide

---

### 6.3 NATS & JetStream

**Goal:** Learn lightweight, high-performance messaging

**Learning objectives:**
- Understand NATS core vs JetStream
- Implement request-reply patterns
- Compare with Kafka and RabbitMQ

**Tasks:**
- [ ] Create `experiments/nats-tutorial/`
- [ ] Deploy NATS with JetStream enabled
- [ ] Core NATS patterns:
  - [ ] Pub/sub (fire and forget)
  - [ ] Request/reply
  - [ ] Queue groups (load balancing)
- [ ] JetStream (persistence):
  - [ ] Streams and consumers
  - [ ] At-least-once delivery
  - [ ] Message replay
  - [ ] Key-value store
  - [ ] Object store
- [ ] Build demo apps showcasing each pattern
- [ ] Compare with Kafka/RabbitMQ:
  - [ ] Latency
  - [ ] Throughput
  - [ ] Operational complexity
  - [ ] Use case fit
- [ ] Document NATS patterns and when to use

---

### 6.4 Cloud Messaging with Crossplane

**Goal:** Abstract cloud message queues with Crossplane XRDs

**Learning objectives:**
- Use Crossplane for managed messaging services
- Create portable queue abstractions
- Compare managed vs self-hosted

**Tasks:**
- [ ] Create `experiments/cloud-messaging/`
- [ ] Create XRD: SimpleQueue
  - [ ] Abstracts AWS SQS and Azure Service Bus
  - [ ] Common interface for both clouds
- [ ] Deploy same producer/consumer app to both clouds
- [ ] Compare:
  - [ ] Message visibility handling
  - [ ] Dead letter queue behavior
  - [ ] FIFO vs standard queues
  - [ ] Pricing models
- [ ] Test failover (queue in different region)
- [ ] Document cloud queue patterns

---

## Phase 7: Deployment Strategies

*Progressive complexity: rolling → blue-green → canary → GitOps patterns.*

### 7.1 Rolling Updates Optimization

**Goal:** Master Kubernetes native rolling deployments

**Learning objectives:**
- Understand rolling update parameters
- Optimize for zero-downtime deployments
- Handle graceful shutdown correctly

**Tasks:**
- [ ] Create `experiments/rolling-update-tutorial/`
- [ ] Build app with slow startup and graceful shutdown
- [ ] Test parameter combinations:
  - [ ] maxSurge/maxUnavailable variations
  - [ ] minReadySeconds impact
  - [ ] progressDeadlineSeconds
- [ ] Implement graceful shutdown:
  - [ ] preStop hooks
  - [ ] terminationGracePeriodSeconds
  - [ ] Connection draining
- [ ] Readiness probe tuning:
  - [ ] initialDelaySeconds
  - [ ] periodSeconds
  - [ ] failureThreshold
- [ ] Load test during rollout (measure errors)
- [ ] Document recommended configurations

---

### 7.2 Blue-Green Deployments

**Goal:** Implement instant cutover deployments

**Learning objectives:**
- Understand blue-green pattern
- Implement with different tools
- Handle rollback scenarios

**Tasks:**
- [ ] Create `experiments/blue-green-tutorial/`
- [ ] Implement blue-green with:
  - [ ] Kubernetes Services (label selector swap)
  - [ ] Gateway API traffic switching
  - [ ] Argo Rollouts BlueGreen strategy
- [ ] Test scenarios:
  - [ ] Successful deployment
  - [ ] Failed health check (no switch)
  - [ ] Rollback after deployment
- [ ] Measure:
  - [ ] Cutover time
  - [ ] Request failures during switch
  - [ ] Resource overhead (2x replicas)
- [ ] Handle stateful considerations:
  - [ ] Database compatibility
  - [ ] Session handling
- [ ] Document blue-green patterns

---

### 7.3 Canary Deployments with Argo Rollouts

**Goal:** Implement gradual traffic shifting with automated analysis

**Learning objectives:**
- Understand canary deployment pattern
- Configure Argo Rollouts
- Implement metric-based promotion/rollback

**Tasks:**
- [ ] Create `experiments/canary-tutorial/`
- [ ] Install Argo Rollouts
- [ ] Configure Rollout resource:
  - [ ] Traffic splitting steps (5% → 25% → 50% → 100%)
  - [ ] Pause durations
  - [ ] Manual gates
- [ ] Implement AnalysisTemplate:
  - [ ] Success rate query (Prometheus)
  - [ ] Latency threshold query
  - [ ] Custom business metrics
- [ ] Create "bad" versions to test:
  - [ ] High error rate version
  - [ ] High latency version
- [ ] Test automated rollback on failure
- [ ] Integrate with:
  - [ ] Gateway API (traffic splitting)
  - [ ] Istio (if mesh deployed)
- [ ] Document canary analysis patterns

---

### 7.4 GitOps Patterns with ArgoCD

**Goal:** Master ArgoCD for GitOps deployments

**Learning objectives:**
- Understand ArgoCD sync strategies
- Implement progressive delivery via Git
- Use ApplicationSets for multi-cluster

**Tasks:**
- [ ] Create `experiments/argocd-patterns/`
- [ ] Sync strategies:
  - [ ] Auto-sync vs manual
  - [ ] Self-heal behavior
  - [ ] Prune policies
- [ ] Sync waves and hooks:
  - [ ] Pre-sync hooks (DB migration job)
  - [ ] Sync waves (ordering)
  - [ ] Post-sync hooks (smoke tests)
  - [ ] SyncFail hooks (notifications)
- [ ] ApplicationSet patterns:
  - [ ] Git generator (directory/file)
  - [ ] Cluster generator (multi-cluster)
  - [ ] Matrix generator (combinations)
  - [ ] Progressive rollout across clusters
- [ ] App-of-apps pattern
- [ ] Document GitOps workflow patterns

---

## Phase 8: Autoscaling & Resource Management

*Scale applications efficiently based on various signals.*

### 8.1 Horizontal Pod Autoscaler Deep Dive

**Goal:** Master HPA configuration for different workloads

**Learning objectives:**
- Understand HPA algorithm and behavior
- Configure for various metric types
- Tune for responsiveness vs stability

**Tasks:**
- [ ] Create `experiments/hpa-tutorial/`
- [ ] Build test app with configurable CPU/memory load
- [ ] Configure HPA scenarios:
  - [ ] CPU-based scaling
  - [ ] Memory-based scaling
  - [ ] Custom metrics (Prometheus adapter)
  - [ ] External metrics
- [ ] Tune parameters:
  - [ ] Target utilization thresholds
  - [ ] Stabilization windows (scale up/down)
  - [ ] Scaling policies (pods vs percent)
- [ ] Test workload patterns:
  - [ ] Gradual ramp-up
  - [ ] Sudden spike
  - [ ] Oscillating load
- [ ] Measure:
  - [ ] Time to scale
  - [ ] Over/under provisioning
  - [ ] Request latency during scaling
- [ ] Document HPA tuning guide

---

### 8.2 KEDA Event-Driven Autoscaling

**Goal:** Scale based on external event sources

**Learning objectives:**
- Understand KEDA architecture
- Configure various scalers
- Implement scale-to-zero

**Tasks:**
- [ ] Create `experiments/keda-tutorial/`
- [ ] Install KEDA
- [ ] Implement scalers:
  - [ ] Prometheus scaler (custom metrics)
  - [ ] Kafka scaler (consumer lag)
  - [ ] RabbitMQ scaler (queue depth)
  - [ ] Cron scaler (scheduled scaling)
  - [ ] Azure Service Bus / AWS SQS (via Crossplane)
- [ ] Configure ScaledObject:
  - [ ] Triggers and thresholds
  - [ ] Cooldown periods
  - [ ] Min/max replicas
  - [ ] Scale-to-zero behavior
- [ ] Test ScaledJob for batch workloads
- [ ] Compare KEDA vs HPA:
  - [ ] Configuration complexity
  - [ ] Supported triggers
  - [ ] Scale-to-zero capability
- [ ] Document KEDA patterns

---

### 8.3 Vertical Pod Autoscaler

**Goal:** Right-size pod resource requests automatically

**Learning objectives:**
- Understand VPA modes and recommendations
- Combine VPA with HPA
- Implement resource optimization workflow

**Tasks:**
- [ ] Create `experiments/vpa-tutorial/`
- [ ] Install VPA
- [ ] Configure VPA modes:
  - [ ] Off (recommendations only)
  - [ ] Initial (set on pod creation)
  - [ ] Auto (update running pods)
- [ ] Test with various workloads:
  - [ ] CPU-bound application
  - [ ] Memory-bound application
  - [ ] Variable workload
- [ ] Analyze recommendations:
  - [ ] Lower bound, target, upper bound
  - [ ] Uncapped vs capped
- [ ] Combine with HPA (mutually exclusive metrics)
- [ ] Document resource optimization workflow

---

## Phase 9: Data & Storage

*Stateful workloads: databases, caching, and persistent storage.*

### 9.1 PostgreSQL with CloudNativePG

**Goal:** Operate PostgreSQL on Kubernetes with CloudNativePG

**Learning objectives:**
- Understand CloudNativePG operator
- Configure HA PostgreSQL clusters
- Implement backup and recovery

**Tasks:**
- [ ] Create `experiments/postgres-tutorial/`
- [ ] Deploy CloudNativePG operator
- [ ] Create PostgreSQL cluster:
  - [ ] Primary + replicas
  - [ ] Synchronous replication
  - [ ] Connection pooling (PgBouncer)
- [ ] Configure storage:
  - [ ] PVC sizing and storage class
  - [ ] WAL archiving to object storage
- [ ] Backup and recovery:
  - [ ] Scheduled backups (to S3/Azure via Crossplane)
  - [ ] Point-in-time recovery (PITR)
  - [ ] Restore to new cluster
- [ ] Monitoring:
  - [ ] pg_stat metrics in Prometheus
  - [ ] Grafana dashboards
  - [ ] Alerting on replication lag
- [ ] Failover testing:
  - [ ] Kill primary, verify promotion
  - [ ] Measure failover time
- [ ] Document PostgreSQL operational patterns

---

### 9.2 Redis with Spotahome Operator

**Goal:** Operate Redis on Kubernetes for caching

**Learning objectives:**
- Understand Redis sentinel vs cluster mode
- Configure persistence and HA
- Implement caching patterns

**Tasks:**
- [ ] Create `experiments/redis-tutorial/`
- [ ] Deploy Redis operator (Spotahome or similar)
- [ ] Create Redis deployments:
  - [ ] Standalone (development)
  - [ ] Sentinel (HA failover)
  - [ ] Cluster (horizontal scaling)
- [ ] Configure:
  - [ ] Persistence (RDB/AOF)
  - [ ] Memory limits and eviction
  - [ ] Password authentication
- [ ] Implement caching patterns:
  - [ ] Cache-aside
  - [ ] Write-through
  - [ ] Session storage
- [ ] Monitoring:
  - [ ] Redis metrics in Prometheus
  - [ ] Memory usage tracking
  - [ ] Hit/miss ratio
- [ ] Document Redis patterns for Kubernetes

---

### 9.3 Object Storage with MinIO

**Goal:** Deploy S3-compatible object storage

**Learning objectives:**
- Understand MinIO architecture
- Configure for different use cases
- Integrate with backup solutions

**Tasks:**
- [ ] Create `experiments/minio-tutorial/`
- [ ] Deploy MinIO operator
- [ ] Create MinIO tenant:
  - [ ] Single node (development)
  - [ ] Multi-node distributed (HA)
- [ ] Configure:
  - [ ] Buckets and policies
  - [ ] Access keys and IAM
  - [ ] Lifecycle rules
  - [ ] Versioning
- [ ] Integrate with:
  - [ ] Loki (log storage)
  - [ ] Tempo (trace storage)
  - [ ] Velero (cluster backups)
  - [ ] CloudNativePG (WAL archive)
- [ ] Compare with Crossplane S3 buckets
- [ ] Document object storage patterns

---

## Phase 10: Argo Workflows & Automation

*Complex workflow orchestration for CI/CD, data pipelines, and experiments.*

### 10.1 Argo Workflows Deep Dive

**Goal:** Master workflow orchestration patterns

**Learning objectives:**
- Understand Argo Workflows concepts
- Build complex multi-step workflows
- Handle artifacts and parameters

**Tasks:**
- [ ] Create `experiments/argo-workflows-tutorial/`
- [ ] Workflow patterns:
  - [ ] Sequential steps
  - [ ] Parallel execution
  - [ ] DAG dependencies
  - [ ] Conditional execution (when)
  - [ ] Loops (withItems, withParam)
- [ ] Parameters and artifacts:
  - [ ] Input/output parameters
  - [ ] Artifact passing between steps
  - [ ] S3/MinIO artifact storage
- [ ] Templates:
  - [ ] Container templates
  - [ ] Script templates
  - [ ] WorkflowTemplate (reusable)
  - [ ] ClusterWorkflowTemplate
- [ ] Error handling:
  - [ ] Retry strategies
  - [ ] Timeout configuration
  - [ ] Exit handlers
  - [ ] ContinueOn failure
- [ ] Build practical workflows:
  - [ ] CI pipeline (build → test → deploy)
  - [ ] Data processing pipeline
  - [ ] Experiment runner (this lab!)
- [ ] Document workflow patterns

---

### 10.2 Argo Events

**Goal:** Event-driven workflow triggering

**Learning objectives:**
- Understand Argo Events architecture
- Configure event sources and sensors
- Integrate with Argo Workflows

**Tasks:**
- [ ] Create `experiments/argo-events-tutorial/`
- [ ] Deploy Argo Events
- [ ] Configure EventSources:
  - [ ] Webhook (HTTP triggers)
  - [ ] GitHub (push, PR events)
  - [ ] Kafka (message triggers)
  - [ ] Cron (scheduled triggers)
  - [ ] S3/MinIO (object events)
- [ ] Configure Sensors:
  - [ ] Event filtering
  - [ ] Parameter extraction
  - [ ] Trigger templates
- [ ] Integrate triggers:
  - [ ] Trigger Argo Workflow
  - [ ] Trigger Kubernetes resource
  - [ ] Trigger HTTP endpoint
- [ ] Build event-driven pipelines:
  - [ ] GitHub push → build workflow
  - [ ] S3 upload → processing workflow
  - [ ] Scheduled experiment runs
- [ ] Document event-driven patterns

---

## Phase 11: Advanced Topics & Benchmarks

*Deep dives and performance comparisons - now that fundamentals are solid.*

### 11.1 Database Performance Comparison

**Goal:** Compare relational databases for Kubernetes workloads

**Learning objectives:**
- Benchmark database performance objectively
- Understand trade-offs between options
- Make data-driven database selection

**Tasks:**
- [ ] Create `experiments/database-benchmark/`
- [ ] Deploy databases via Crossplane/operators:
  - [ ] PostgreSQL (CloudNativePG)
  - [ ] MySQL (via operator)
  - [ ] Cloud-managed (Azure SQL, RDS via Crossplane)
- [ ] Create benchmark schema and data
- [ ] Run benchmarks:
  - [ ] pgbench / sysbench
  - [ ] OLTP workloads (TPC-C style)
  - [ ] Read-heavy vs write-heavy
- [ ] Compare:
  - [ ] Throughput (TPS)
  - [ ] Latency percentiles
  - [ ] Resource consumption
  - [ ] Operational complexity
- [ ] Document findings and recommendations

---

### 11.2 Message Queue Performance Comparison

**Goal:** Compare messaging systems under load

**Learning objectives:**
- Benchmark throughput and latency
- Understand performance characteristics
- Inform technology selection

**Tasks:**
- [ ] Create `experiments/messaging-benchmark/`
- [ ] Deploy all three brokers (from Phase 6)
- [ ] Build benchmarking clients
- [ ] Test scenarios:
  - [ ] High throughput (max messages/sec)
  - [ ] Low latency (p99 measurement)
  - [ ] Fan-out (1 → N consumers)
  - [ ] Persistence impact
- [ ] Compare:
  - [ ] Messages per second
  - [ ] End-to-end latency
  - [ ] Resource consumption
  - [ ] Recovery time after failure
- [ ] Document performance comparison

---

### 11.3 Service Mesh Performance Comparison

**Goal:** Measure service mesh overhead

**Learning objectives:**
- Quantify latency overhead
- Compare resource consumption
- Inform mesh selection

**Tasks:**
- [ ] Create `experiments/mesh-benchmark/`
- [ ] Deploy baseline app (no mesh)
- [ ] Deploy same app with:
  - [ ] Istio
  - [ ] Linkerd
  - [ ] Cilium
- [ ] Measure:
  - [ ] Latency overhead (p50, p95, p99)
  - [ ] CPU per pod (sidecar cost)
  - [ ] Memory per pod
  - [ ] Control plane resources
- [ ] Test at scale:
  - [ ] 10, 50, 100 services
  - [ ] High RPS scenarios
- [ ] Document mesh comparison

---

### 11.4 Runtime Performance Comparison

**Goal:** Compare web server runtimes for API workloads

**Learning objectives:**
- Benchmark different language runtimes
- Understand performance characteristics
- Portfolio piece for runtime expertise

**Tasks:**
- [ ] Create `experiments/runtime-benchmark/`
- [ ] Build identical API in:
  - [ ] Go (net/http)
  - [ ] Rust (Axum)
  - [ ] .NET (ASP.NET Core)
  - [ ] Node.js (Fastify)
  - [ ] Bun
- [ ] Implement endpoints:
  - [ ] GET /health
  - [ ] GET /json (serialize response)
  - [ ] POST /echo (deserialize + serialize)
  - [ ] GET /compute (CPU-bound work)
- [ ] Benchmark with k6:
  - [ ] Throughput (RPS)
  - [ ] Latency distribution
  - [ ] Memory footprint
  - [ ] Container image size
  - [ ] Cold start time
- [ ] Document runtime comparison

---

## Phase 12: Chaos Engineering (Nice to Have)

*Validate resilience - capstone experiments after everything else is solid.*

### 12.1 Pod Failure & Recovery

**Goal:** Measure application resilience to pod failures

**Tasks:**
- [ ] Create `experiments/chaos-pod-failure/`
- [ ] Deploy Chaos Mesh
- [ ] Test scenarios:
  - [ ] Single pod kill
  - [ ] Multiple pod kill (50%)
  - [ ] Container crash loop
- [ ] Measure recovery time and error rates
- [ ] Document resilience findings

---

### 12.2 Network Chaos

**Goal:** Test application behavior under network issues

**Tasks:**
- [ ] Create `experiments/chaos-network/`
- [ ] Test with Chaos Mesh NetworkChaos:
  - [ ] Latency injection (50ms, 200ms, 500ms)
  - [ ] Packet loss (1%, 5%, 20%)
  - [ ] Network partition
- [ ] Measure:
  - [ ] Timeout behavior
  - [ ] Circuit breaker activation
  - [ ] Retry storms
- [ ] Document network resilience patterns

---

### 12.3 Node Drain & Zone Failure

**Goal:** Test infrastructure-level failures

**Tasks:**
- [ ] Create `experiments/chaos-infrastructure/`
- [ ] Test scenarios:
  - [ ] Graceful node drain
  - [ ] Sudden node failure
  - [ ] Zone failure (multi-zone cluster)
- [ ] Measure:
  - [ ] Workload redistribution time
  - [ ] Request failures during event
  - [ ] PVC reattachment time
- [ ] Document infrastructure resilience

---

## Learning Path Summary

| Phase | Focus | Experiments | Key Skills |
|-------|-------|-------------|------------|
| 1 | Platform Bootstrap | 2 | Spacelift, Crossplane XRDs, Kind |
| 2 | Security | 3 | cert-manager, Vault, NetworkPolicy |
| 3 | Observability | 3 | Prometheus, Loki, OpenTelemetry |
| 4 | Traffic Management | 2 | Gateway API, Ingress controllers |
| 5 | Service Mesh | 3 | Istio, Linkerd, Cilium |
| 6 | Messaging | 4 | Kafka, RabbitMQ, NATS, Crossplane |
| 7 | Deployment Strategies | 4 | Rolling, Blue-Green, Canary, GitOps |
| 8 | Autoscaling | 3 | HPA, KEDA, VPA |
| 9 | Data & Storage | 3 | PostgreSQL, Redis, MinIO |
| 10 | Argo Workflows | 2 | Workflows, Events |
| 11 | Benchmarks | 4 | DB, Messaging, Mesh, Runtime comparisons |
| 12 | Chaos Engineering | 3 | Pod, Network, Infrastructure chaos |

**Total: 36 experiments**

---

## Notes

- All experiments follow `experiments/_template/` structure
- Use Crossplane claims for cloud resources where applicable
- Spacelift for cloud deployments, Taskfile for local Kind
- Each experiment should be portfolio-ready with clear README
- Security (TLS, secrets) established early and used throughout
- Tutorials first, benchmarks after fundamentals are solid
