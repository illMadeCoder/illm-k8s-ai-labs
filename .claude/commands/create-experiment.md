---
description: Scaffold a new experiment directory with validated experiment.yaml
allowed-tools: Bash, Read, Write, Glob, Grep, Edit, AskUserQuestion
argument-hint: [experiment-name]
---

# Create Experiment

Scaffold a new experiment with a validated `experiment.yaml`. Guide the user through an interactive flow, then generate the file and validate it.

## Step 1: Parse Arguments

Extract the experiment name from `$ARGUMENTS`. If empty, ask the user with `AskUserQuestion`.

The experiment name becomes the directory under `experiments/`. It must be lowercase, hyphenated, and descriptive (e.g., `metrics-comparison`, `tempo-tutorial`).

Check that `experiments/{name}/` does not already exist. If it does, tell the user and stop.

## Step 2: Gather Experiment Details

Use `AskUserQuestion` to collect inputs. You may combine questions into a single call (up to 4 questions per call).

### Question Set 1: Type & Domain

Ask these two questions together:

**Experiment type** (determines structure):
- `comparison` — Compares two or more tools/approaches. Gets `publish: true`, metrics, hypothesis, analyzerConfig. Workflow completion mode: `workflow`.
- `tutorial` — Interactive learning. Gets `tutorial:` section with `exposeKubeconfig: true` and service refs. No publish, no metrics. Workflow completion mode: `manual`.
- `demo` — Simple demonstration. Minimal config, no publish. Workflow completion mode: `workflow`.
- `baseline` — Establishes baseline measurements. Gets metrics but simpler than comparison. Workflow completion mode: `workflow`.

**Domain** (determines tags):
- `observability` — metrics, logging, tracing, slos, cost
- `networking` — gateways, ingress, service-mesh
- `storage` — object-storage, database
- `cicd` — pipelines, ci, cd, supply-chain

### Question Set 2: Target Configuration

Ask about target setup:

**Number of targets:**
- `Single target (app only)` — Most tutorials and simple experiments
- `Two targets (app + loadgen)` — Comparisons with load testing (loadgen depends on app)
- `Custom` — Let user specify target names

**Machine type for main target:**
- `e2-standard-2` — Lightweight (loadgen, simple demos)
- `e2-standard-4` — Default (most experiments)
- `e2-standard-8` — Heavy compute (large observability stacks)

All targets default to `preemptible: true` and observability enabled with `transport: tailscale`.

### Question Set 3: Components

Present the component catalog grouped by category. Ask the user which components to include for each target. You can suggest relevant components based on the domain:

**Component Catalog (28 components):**

**apps/** (7):
- `hello-app` — Simple hello world app for load testing
- `nginx` — NGINX web server
- `otel-demo` — Multi-service OTel demo (user-service, order-service, payment-service)
- `cardinality-generator` — High-cardinality Prometheus metrics for cost analysis
- `station-monitor` — Station monitoring app for Prometheus tutorial
- `log-generator` — Structured logs for logging pipeline tutorials
- `metrics-app` — Demo app exposing Prometheus metrics

**core/** (4):
- `nginx-ingress` — NGINX Ingress Controller
- `envoy-gateway` — Envoy Gateway (Gateway API)
- `traefik` — Traefik proxy with Gateway API support
- `tailscale-operator` — Tailscale K8s operator for mesh networking

**observability/** (15):
- `kube-prometheus-stack` — Prometheus + Grafana + Alertmanager
- `victoria-metrics` — VictoriaMetrics single-node TSDB
- `mimir` — Grafana Mimir TSDB (monolithic mode)
- `loki` — Loki log aggregation
- `promtail` — Promtail log collector (ships to Loki)
- `fluent-bit` — Fluent Bit log processor/forwarder
- `elasticsearch` — Elasticsearch for log storage
- `kibana` — Kibana dashboard for Elasticsearch
- `tempo` — Grafana Tempo distributed tracing
- `jaeger` — Jaeger distributed tracing (all-in-one)
- `otel-collector` — OpenTelemetry Collector
- `pyrra` — Pyrra SLO management
- `seaweedfs` — SeaweedFS S3-compatible object storage
- `metrics-agent` — Grafana Alloy agent (scrapes + remote-writes to hub VictoriaMetrics)
- `metrics-egress` — Tailscale egress service to hub VictoriaMetrics

**storage/** (1):
- `minio` — MinIO S3-compatible object storage

**testing/** (1):
- `k6-gateway-loadtest` — k6 load test for gateway comparison (runs on loadgen cluster)

Suggest relevant components based on domain and type. For example:
- Observability comparison → suggest the stacks being compared + `kube-prometheus-stack` + `metrics-agent` + `metrics-egress`
- Networking comparison → suggest gateway controllers + `hello-app` + `k6-gateway-loadtest`
- Tutorial → suggest the learning target + `kube-prometheus-stack` for dashboards

### Question Set 4: Publish & Analysis (comparison/baseline only)

If the experiment type is `comparison` or `baseline`, ask:

**Publish to benchmark site?**
- `Yes` — Sets `publish: true`, includes `analyzerConfig`
- `No` — Results stored in S3 only

If publishing, the standard analyzer sections are included automatically. Ask if they want to customize (most users won't):
- Core: `abstract`, `targetAnalysis`, `performanceAnalysis`, `metricInsights`
- FinOps: `finopsAnalysis` (include by default), `secopsAnalysis` (omit by default unless security-relevant)
- Deep dive: `body`, `capabilitiesMatrix` (comparisons only), `feedback`
- Diagram: `architectureDiagram`

## Step 3: Generate the Experiment YAML

### GKE Name Length Validation

Before writing, validate GKE cluster name lengths:

```
GKE name = "illm-" (5) + experimentName + "-" (1) + targetName + "-" (1) + xrSuffix (5) = 12 + len(experimentName) + len(targetName)
experimentName = generateNamePrefix + k8sSuffix (5 chars)
Total = 17 + len(generateNamePrefix) + len(targetName), must be <= 40
```

The `generateName` prefix should be an abbreviation of the experiment name followed by a hyphen:
- `gateway-comparison` → `gw-comp-`
- `logging-comparison` → `logging-comparison-` (if it fits)
- `prometheus-tutorial` → `prometheus-tutorial-`
- Long names → abbreviate: `observability-cost-tutorial` → `obs-cost-tut-`

**For each target, check:** `17 + len(prefix) + len(targetName) <= 40`

If any target exceeds 40 chars, shorten the `generateName` prefix and warn the user.

### YAML Structure

Generate `experiments/{name}/experiment.yaml` following these patterns:

**All experiments:**
```yaml
apiVersion: experiments.illm.io/v1alpha1
kind: Experiment
metadata:
  generateName: {prefix}-
  namespace: experiments
spec:
  description: "{description}"
  tags: ["{type}", "{domain}", ...additional-tags]
  # publish: true  — only for comparisons/baselines being published
  targets:
    - name: app
      cluster:
        type: gke
        machineType: {machineType}
        preemptible: true
      observability:
        enabled: true
        transport: tailscale
      components:
        - app: {component1}
        - app: {component2}
          params:
            key: "value"
  workflow:
    template: {name}-validation
    completion:
      mode: workflow  # or manual for tutorials
```

**Comparisons add** (when published):
```yaml
  publish: true

  analyzerConfig:
    sections:
      - abstract
      - targetAnalysis
      - performanceAnalysis
      - metricInsights
      - finopsAnalysis
      # - secopsAnalysis      # Uncomment for security-relevant experiments
      - body
      - capabilitiesMatrix    # Only for comparisons
      - feedback
      - architectureDiagram

  hypothesis:
    claim: "{user-provided or placeholder}"
    questions:
      - "{question1}"
    focus:
      - "{focus1}"

  metrics:
    # TODO: Add PromQL queries for the specific stacks being compared
    # See experiments/logging-comparison/experiment.yaml for examples
    # Metric names must match ^[a-z][a-z0-9_]*$
    # Types: instant (bar charts) or range (time-series line charts)
    # Variables: $EXPERIMENT, $NAMESPACE, $DURATION
    []
```

**Tutorials add:**
```yaml
  tutorial:
    exposeKubeconfig: true
    services:
      - name: grafana
        target: app
        service: kube-prometheus-stack-grafana
        namespace: {experiment-name}
```

**Two-target (app + loadgen) adds:**
```yaml
    - name: loadgen
      depends: [app]
      cluster:
        type: gke
        machineType: e2-standard-2
        preemptible: true
      components:
        - app: k6-gateway-loadtest
```

### Tag Conventions

Tags are used for site categorization. Always include:
1. The experiment type: `comparison`, `tutorial`, `demo`, or `baseline`
2. The domain: `observability`, `networking`, `storage`, `cicd`
3. Specific technology tags (lowercase): `prometheus`, `loki`, `gateway`, `envoy`, `tracing`, etc.

### Hypothesis Guidance (comparisons)

Ask the user for a 1-2 sentence claim about expected outcome. If they don't have one, generate a reasonable placeholder based on the components being compared. Include 2-3 questions the experiment should answer and 2-3 focus areas.

### Metrics Guidance (comparisons/baselines)

For published experiments, add a `metrics: []` placeholder with TODO comments explaining the pattern. Metrics are highly experiment-specific (PromQL queries referencing exact pod labels) and are best added after the first run reveals the actual label set. Point the user to `experiments/logging-comparison/experiment.yaml` or `experiments/gateway-comparison/experiment.yaml` as examples.

## Step 4: Write the File

Use the `Write` tool to create `experiments/{name}/experiment.yaml`.

## Step 5: Validate

After writing, run the `experiment-validator` agent on the new experiment to catch any issues:

```
Run experiment-validator agent with the experiment name
```

Report the validation results. If there are failures, fix them and re-validate.

## Step 6: Summary

Tell the user:
1. What was created: `experiments/{name}/experiment.yaml`
2. The `generateName` prefix and GKE name length status
3. Any TODOs they need to fill in (metrics queries, hypothesis details, workflow template)
4. How to run: `kubectl create -f experiments/{name}/experiment.yaml`
5. Remind them the workflow template (`{name}-validation`) needs to exist as an Argo WorkflowTemplate before the experiment can run
