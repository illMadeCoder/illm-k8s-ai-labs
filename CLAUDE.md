# CLAUDE.md

> Keep this file short and stable. Track progress in `TODO.md`, not here.
> Use `/compact` at ~70% context. Start sessions with "Continue from TODO.md Phase X".

## Project Overview

**illm-k8s-lab** is a learning-focused Kubernetes experiment platform for Cloud/Platform/Solutions Architect roles. Phased roadmap in `TODO.md`.

## Project Structure

```
workload-catalog/          # ArgoCD apps and Helm values
├── components/            # Platform components by category
└── stacks/                # Grouped deployments

experiments/{name}/        # Individual experiments
├── target/argocd/         # Target cluster apps
├── loadgen/               # Load generator configs
└── workflow/              # Argo Workflow definitions

platform/terraform/        # Infrastructure as Code
├── spacelift/             # Admin stack (manages other stacks)
├── azure/                 # Azure resources
└── aws/                   # AWS resources

docs/adrs/                 # Architecture Decision Records
```

## Key Decisions

- **IaC**: Spacelift + Terraform
- **Secrets**: ESO + Vault (ADR-002)
- **GitOps**: ArgoCD app-of-apps
- **CI/CD**: GitHub Actions primary

## Commands

```bash
task exp:run:{name}        # Run experiment
task exp:deploy:{name}     # Deploy experiment
task exp:undeploy:{name}   # Cleanup
git push                   # Triggers Spacelift autodeploy
```

## Conventions

- Terraform: Spacelift manages all stacks, contexts hold credentials
- ArgoCD: Multi-source pattern, sync waves, `ignoreDifferences` for CRDs
- Experiments: Use labels `experiment: {name}`, `cluster: target|loadgen`
