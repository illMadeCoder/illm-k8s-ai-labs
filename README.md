# illm-k8s-ai-lab

A GitOps-driven Kubernetes learning environment with reproducible, scenario-based experiments.

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                            Hub Cluster                              │
│  ┌──────────┐  ┌──────────┐  ┌────────────┐  ┌─────────────────┐   │
│  │  ArgoCD  │  │ OpenBao  │  │ MetalLB    │  │ k8s_gateway     │   │
│  │ (GitOps) │  │ (Secrets)│  │ (LoadBal)  │  │ (DNS)           │   │
│  └────┬─────┘  └──────────┘  └────────────┘  └─────────────────┘   │
│       │                                                             │
│       ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────┐   │
│  │                    Experiment Namespace                      │   │
│  │  ┌─────────┐  ┌─────────────┐  ┌───────┐  ┌──────────────┐  │   │
│  │  │ Demo    │  │ Prometheus  │  │ k6    │  │ Argo         │  │   │
│  │  │ App     │  │ + Grafana   │  │ Tests │  │ Workflows    │  │   │
│  │  └─────────┘  └─────────────┘  └───────┘  └──────────────┘  │   │
│  └─────────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────────┘
```

## Quick Start

```bash
# Prerequisites: Docker, kubectl, Task, Helm

# Create Kind cluster and bootstrap ArgoCD
task kind:bootstrap

# Run an experiment
task kind:conduct -- prometheus-tutorial

# Teardown experiment
task kind:teardown -- prometheus-tutorial

# Delete cluster
task kind:delete
```

## Project Structure

```
├── hub/                    # Hub cluster configuration
│   ├── app-of-apps/        # ArgoCD application-of-apps pattern
│   │   └── kind/           # Kind-specific apps and values
│   └── bootstrap/          # Initial cluster setup
│
├── experiments/
│   ├── scenarios/          # Runnable experiments
│   │   ├── prometheus-tutorial/
│   │   ├── http-baseline/
│   │   └── ...
│   └── components/         # Reusable ArgoCD components
│       ├── observability/  # Prometheus, Grafana, Loki
│       ├── testing/        # k6, Argo Workflows
│       └── ...
│
├── platforms/
│   └── kind/               # Kind cluster tasks and config
│
└── docs/
    ├── adrs/               # Architecture Decision Records
    └── roadmap/            # Learning phases and topics
```

## Available Experiments

| Experiment | Description |
|------------|-------------|
| `prometheus-tutorial` | Interactive observability walkthrough with Prometheus, Grafana, and custom metrics |
| `http-baseline` | Load testing baseline with k6 |
| `hello-app` | Minimal deployment example |

## Task Commands

| Command | Description |
|---------|-------------|
| `task kind:bootstrap` | Create Kind cluster with ArgoCD + core infra |
| `task kind:conduct -- <name>` | Deploy and run experiment |
| `task kind:teardown -- <name>` | Clean up experiment |
| `task kind:delete` | Delete Kind cluster |
| `task kind:argocd` | Open ArgoCD UI |

## License

MIT
