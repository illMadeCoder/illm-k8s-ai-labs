# Experiments

Experiments are self-contained scenarios that deploy infrastructure and applications for testing specific hypotheses.

## Structure

```
experiments/
└── <experiment-name>/
    ├── argocd/              # ArgoCD applications (one per cluster)
    │   ├── target.yaml      # Apps for "target" cluster
    │   └── loadgen.yaml     # Apps for "loadgen" cluster
    ├── terraform/           # Infrastructure per environment
    │   └── prod/
    │       ├── main.tf
    │       ├── variables.tf
    │       └── outputs.tf
    └── k6/                  # Load test scripts
        └── baseline.js
```

## Conventions

### ArgoCD Files → Cluster Matching

ArgoCD application files are named after the cluster they deploy to:

| File | Deploys To |
|------|------------|
| `argocd/target.yaml` | Cluster named "target" |
| `argocd/loadgen.yaml` | Cluster named "loadgen" |
| `argocd/worker.yaml` | Cluster named "worker" |

This convention enables automatic deployment - the task iterates through clusters and finds the matching ArgoCD file.

### Terraform Cluster Definitions

Clusters are defined in `terraform/prod/main.tf`:

```hcl
variable "clusters" {
  default = {
    target = {
      vm_size    = "Standard_D4s_v3"
      node_count = 3
      min_nodes  = 2
      max_nodes  = 10
    }
    loadgen = {
      vm_size    = "Standard_D2s_v3"
      node_count = 2
      min_nodes  = 1
      max_nodes  = 5
    }
  }
}
```

Cluster names here **must match** ArgoCD filenames.

## Deployment

### Local (Minikube)

Deploys **all** ArgoCD files to a single cluster:

```bash
# Prerequisites
task deploy:core                              # Deploy core infra (k6, observability)

# Deploy experiment
task exp:deploy:minikube NAME=http-baseline   # Deploys all argocd/*.yaml

# Run load test
task exp:run USERS=10 DURATION=60s

# Teardown
task exp:undeploy:minikube NAME=http-baseline
```

### Production (Azure AKS)

Creates multiple clusters and deploys matching ArgoCD apps:

```bash
# Deploy (creates infra + deploys apps)
task exp:deploy:prod NAME=http-baseline

# This runs:
#   1. terraform apply → creates target + loadgen clusters
#   2. For each cluster:
#      - Writes kubeconfig-{cluster} file
#      - Deploys argocd/{cluster}.yaml using that kubeconfig

# View kubeconfigs
task exp:kubeconfig:prod NAME=http-baseline

# Connect to specific cluster
KUBECONFIG=experiments/http-baseline/terraform/prod/kubeconfig-target kubectl get pods

# Teardown (destroys ALL clusters)
task exp:undeploy:prod NAME=http-baseline
```

## Creating a New Experiment

1. **Create directory structure:**
   ```bash
   mkdir -p experiments/my-experiment/{argocd,terraform/prod,k6}
   ```

2. **Define clusters** in `terraform/prod/main.tf`:
   ```hcl
   variable "clusters" {
     default = {
       server = { ... }  # Your cluster names
       client = { ... }
     }
   }
   ```

3. **Create matching ArgoCD apps:**
   ```bash
   # argocd/server.yaml - apps for server cluster
   # argocd/client.yaml - apps for client cluster
   ```

4. **Add load test scripts** in `k6/`

5. **Test locally first:**
   ```bash
   task exp:deploy:minikube NAME=my-experiment
   task exp:run
   ```

## Available Tasks

| Task | Description |
|------|-------------|
| `task exp:list` | List available experiments |
| `task exp:deploy:minikube NAME=x` | Deploy to minikube (all argocd files) |
| `task exp:undeploy:minikube NAME=x` | Remove from minikube |
| `task exp:deploy:prod NAME=x` | Create AKS clusters + deploy apps |
| `task exp:plan:prod NAME=x` | Preview terraform changes |
| `task exp:undeploy:prod NAME=x` | Destroy AKS clusters |
| `task exp:kubeconfig:prod NAME=x` | List kubeconfig files |
| `task exp:run` | Run k6 load test |
| `task exp:status` | Show k6 pod status |
| `task exp:clean` | Clean up k6 pods |

## Core Infrastructure

Experiments depend on core infrastructure deployed separately:

```bash
task deploy:core
```

This deploys shared services to the cluster:
- **k6** - Load testing namespace and scripts
- **Observability** - Prometheus, Grafana, Loki
- **Gateway API** - Ingress/routing
- **Cert Manager** - TLS certificates

Core infra is **not** managed by experiments - it persists across experiment deployments.
