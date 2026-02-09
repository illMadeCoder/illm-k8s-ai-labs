# ADR-017: Benchmark Results Site

## Status

Proposed (2026-02-08)

## Context

The lab runs 15+ experiment scenarios (TSDB comparisons, logging comparisons, gateway benchmarks, etc.) that produce `ExperimentSummary` JSON stored in SeaweedFS S3. The operator's `collectAndStoreResults()` uploads `summary.json` and `metrics-snapshot.json` to `s3://experiment-results/<name>/` on completion (see ADR-015).

Currently the only way to view results is:

```bash
kubectl run -n seaweedfs s3check --rm -it --restart=Never \
  --image=curlimages/curl:8.5.0 -- \
  curl -s http://seaweedfs-s3.seaweedfs.svc.cluster.local:8333/experiment-results/<name>/summary.json
```

There is no public-facing way to display, analyze, or cross-reference benchmark results. The data exists but is locked behind `kubectl` and a Tailscale-only network.

### Requirements

| Requirement | Priority | Notes |
|-------------|----------|-------|
| Public access (no auth) | Must | Portfolio piece, shareable links |
| General-purpose metrics | Must | CPU, memory, latency, throughput, cost — not tied to specific experiment types |
| Auto-publish on experiment completion | Must | Operator-driven, no manual steps |
| Cross-experiment comparison | High | Overlay metrics from different runs |
| Interactive charts | High | Filtering, selection, tooltips |
| Zero hosting cost | High | Lab budget is $0 for static hosting |
| Low operational complexity | Medium | No new runtime infrastructure on hub |

## Options Considered

### Axis 1: Hosting / Site Framework

| Option | How it works | Cost | Complexity |
|--------|-------------|------|------------|
| **GitHub Pages + Astro** | Static site in repo, GitHub Actions builds, Pages serves. Astro's island architecture ships minimal JS. | Free | Medium |
| **GitHub Pages + github-action-benchmark** | Purpose-built GH Action that ingests JSON benchmarks and generates Chart.js pages. | Free | Low |
| **Bencher.dev (SaaS)** | Continuous benchmarking platform. Free for public projects. API-driven data ingestion, hosted dashboard. | Free (public) | Low |
| **Grafana Cloud public dashboards** | Push metrics to Grafana Cloud free tier, share dashboards as public snapshots. | Free tier (10k series) | Low |
| **Self-hosted on hub (Tailscale + Cloudflare Tunnel)** | Host site on hub cluster, expose via Cloudflare Tunnel for public access. | Free | Medium-High |

### Axis 2: Visualization

| Approach | Interactivity | Cross-experiment | Custom metrics |
|----------|--------------|-----------------|----------------|
| **Vega-Lite** (declarative JSON specs) | Filtering, selection, tooltips | Yes (concat datasets) | Yes (any JSON shape) |
| **Observable Plot** (JS library) | Good interactivity | Yes (JS data joins) | Yes |
| **Chart.js** (via github-action-benchmark) | Basic (zoom, hover) | Limited (single page) | Rigid format |
| **Grafana dashboards** | Excellent (native) | Via variables/templating | Must be time-series |
| **D3.js** (low-level) | Unlimited | Yes | Yes |

### Axis 3: Data Pipeline (operator → site)

| Approach | How data flows | Latency | Complexity |
|----------|---------------|---------|------------|
| **Extend experiment-operator** | On completion: push summary JSON to GitHub data branch via API → triggers Actions build | ~2-5 min | Low (add one function) |
| **Separate publisher operator** | Watches Experiment CRs for Complete phase → reads S3 → pushes to GitHub | ~2-5 min | Medium (new operator) |
| **Argo Workflow post-step** | Add a final workflow step that pushes results to the site repo | ~1 min | Low (WorkflowTemplate change) |
| **Periodic GitHub Action** | Cron action pulls from S3 (via Tailscale) and rebuilds site | ~15 min | Low but requires S3 access from GH |

## Decision

**GitHub Pages + Astro + Vega-Lite**, with data pushed by extending the existing experiment-operator.

### Why This Combination

| Factor | Reasoning |
|--------|-----------|
| **Free public hosting** | GitHub Pages — zero cost, global CDN |
| **General-purpose** | Vega-Lite specs handle any JSON shape — not locked to specific metric types |
| **Operator integration** | Extend `collectAndStoreResults()` to also push JSON to a GitHub data branch — minimal new code |
| **Interactive** | Vega-Lite provides filtering, selection, tooltips, cross-experiment comparison |
| **Static** | No runtime infrastructure (no DB, no server process on hub) |
| **Portfolio value** | Public site doubles as a portfolio piece showing real benchmark methodology |

### Why Not the Others

| Option | Reason to reject |
|--------|-----------------|
| **github-action-benchmark** | Too rigid — expects a specific benchmark format (`name`, `unit`, `value`). Weak cross-experiment comparison. Cannot display arbitrary `ExperimentSummary` fields like `costEstimate` or `mimirMetrics`. |
| **Bencher.dev** | Excellent for CI regression tracking but less flexible for arbitrary system design metrics. Locked to their dashboard UX. |
| **Grafana Cloud** | Must structure everything as time-series. Public dashboard snapshots are ephemeral (auto-expire). Free tier limited to 10k series. |
| **Self-hosted** | Unnecessary complexity when GitHub Pages is free and always-on. Adds a process to monitor on the hub. |
| **D3.js** | Maximum flexibility but requires writing chart logic from scratch. Vega-Lite gets 90% of the value declaratively. |
| **Separate publisher operator** | Extra Go binary to build, deploy, and maintain. The experiment-operator already has S3 access and completion hooks. |
| **Periodic cron** | 15-min latency, requires Tailscale access from GitHub Actions runners (network complexity). |

### Pipeline Architecture

```
Experiment CR (phase: Complete)
        │
        ▼
experiment-operator
  collectAndStoreResults()
        │
        ├──► SeaweedFS S3       (summary.json, metrics-snapshot.json)
        │    [existing path]
        │
        └──► GitHub API          (push JSON to illm-benchmarks repo, data/ branch)
             [new: ~50 lines]
                │
                ▼
        GitHub Actions trigger   (on push to data/**)
                │
                ▼
        Astro build              (reads data/*.json, generates Vega-Lite pages)
                │
                ▼
        GitHub Pages             (public, CDN-served)
```

### Operator Change

The `collectAndStoreResults()` function in `internal/controller/experiment_controller.go` currently uploads to S3 and sets `status.resultsURL`. The change adds one call after the S3 upload:

```go
// After S3 upload succeeds:
if r.GitHubPublisher != nil {
    if err := r.GitHubPublisher.Publish(ctx, summary); err != nil {
        log.Error(err, "Failed to publish to GitHub — results are in S3, will retry")
        // Non-fatal: S3 is the source of truth, GitHub is best-effort
    }
}
```

The publisher uses the GitHub Contents API (`PUT /repos/{owner}/{repo}/contents/data/{name}.json`) with a Personal Access Token stored as a Kubernetes Secret. No Git clone required.

### Site Structure

```
illm-benchmarks/                    # Separate public repo
├── src/
│   ├── pages/
│   │   ├── index.astro             # Landing: experiment list with status cards
│   │   └── [experiment].astro      # Per-experiment detail page
│   ├── components/
│   │   ├── VegaChart.astro         # Reusable Vega-Lite chart island
│   │   ├── ComparisonTable.astro   # Side-by-side metric comparison
│   │   └── CostBreakdown.astro     # Cost estimate visualization
│   └── layouts/
│       └── Base.astro
├── data/                           # Auto-populated by operator via GitHub API
│   ├── hello-app-x7k2q.json
│   ├── tsdb-comparison-a1b2.json
│   └── ...
├── astro.config.mjs
└── .github/workflows/
    └── deploy.yaml                 # On push to data/** → build + deploy to Pages
```

### Vega-Lite Example

A single experiment detail page renders specs like:

```json
{
  "$schema": "https://vega.github.io/schema/vega-lite/v5.json",
  "data": {"url": "/data/tsdb-comparison-a1b2.json"},
  "mark": "bar",
  "encoding": {
    "x": {"field": "targets[].name", "type": "nominal"},
    "y": {"field": "targets[].metrics.memoryPeakMB", "type": "quantitative"},
    "color": {"field": "targets[].components[0]", "type": "nominal"}
  }
}
```

The spec is generated at build time from the JSON shape — any metric field in `ExperimentSummary` becomes a chartable dimension without code changes.

## Implementation Phases

### Phase 1: Operator Publisher (~50 LOC)

- Add `GitHubPublisher` interface to experiment-operator
- Implement GitHub Contents API client
- Store PAT as K8s Secret (`github-publisher-token`)
- Push `summary.json` to `illm-benchmarks` repo `data/` directory on experiment completion
- Non-fatal: log error and continue if publish fails (S3 remains source of truth)

### Phase 2: Static Site Scaffold

- Create `illm-benchmarks` repo with Astro + Vega-Lite
- GitHub Actions workflow: on push to `data/**` → `astro build` → deploy to Pages
- Landing page: reads `data/*.json`, renders experiment cards (name, phase, duration, date)
- Detail page: per-experiment Vega-Lite charts for all numeric fields

### Phase 3: Cross-Experiment Comparison

- Comparison view: select 2+ experiments, overlay metrics
- Filter by experiment type (TSDB, logging, gateway)
- Tag-based grouping from `ExperimentSummary` metadata

## Consequences

### Positive

- Public benchmark portfolio at zero hosting cost
- Operator-driven pipeline — experiments auto-publish on completion, no manual steps
- Vega-Lite specs are declarative JSON — add new chart types without code changes
- Cross-experiment comparison (overlay metrics from different runs)
- `ExperimentSummary` JSON is the only contract — site works with any experiment type
- Static site has no runtime to monitor or restart

### Negative

- GitHub API dependency in operator (needs PAT or GitHub App token as K8s Secret)
- Astro build step adds ~30s latency to publishing (total: ~2-5 min end-to-end)
- Vega-Lite has a learning curve vs raw Chart.js for custom interactions
- Site design/layout requires frontend work upfront
- PAT rotation is a manual toil item (mitigate with GitHub App in future)

### Future

- Add Bencher.dev integration for regression detection once experiment cadence increases
- Embed Grafana panel links for live experiment monitoring (VictoriaMetrics dashboards)
- Add experiment metadata tags to `ExperimentSummary` for richer filtering
- Migrate from PAT to GitHub App for token management

## References

- [ADR-009: TSDB Selection](ADR-009-tsdb-selection.md) — Metrics comparison experiments that will populate the site
- [ADR-011: Observability Architecture](ADR-011-observability-architecture.md) — Stack producing the metrics
- [ADR-015: Experiment Operator](ADR-015-experiment-operator.md) — Operator lifecycle and `collectAndStoreResults()`
- [ADR-016: Hub Metrics Backend](ADR-016-hub-metrics-backend.md) — VictoriaMetrics as metrics source
- [Astro Documentation](https://docs.astro.build/)
- [Vega-Lite Documentation](https://vega.github.io/vega-lite/)
- [GitHub Pages Documentation](https://docs.github.com/en/pages)
- [GitHub Contents API](https://docs.github.com/en/rest/repos/contents)
- Operator source: `operators/experiment-operator/internal/controller/experiment_controller.go`
