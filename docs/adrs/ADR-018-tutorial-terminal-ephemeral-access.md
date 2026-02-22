# ADR-018: Tutorial Terminal & Ephemeral Web Access

## Status

Proposed (2026-02-21)

## Context

The lab's identity series (id-101 through id-104) introduces interactive tutorials — step-by-step walkthroughs where users deploy authentication infrastructure and experience the evolution from separate passwords to SSO. Unlike benchmark experiments that produce metrics for the static site, tutorials require live interaction:

1. **Terminal access** — Users run `kubectl` commands, inspect LDAP directories, decode JWT tokens, and observe protocol flows from a shell session.
2. **Ephemeral web UIs** — Tutorial components expose web interfaces (OpenLDAP admin console, Keycloak admin, three tutorial web apps) that users interact with directly in their browser.

Currently, the only way to interact with a running experiment is via `kubectl` on a machine with cluster access. The benchmark site (ADR-017) is purely static — it renders results after experiments complete but cannot provide live access to running infrastructure.

The site already has a pattern for proxying real-time connections through Cloudflare Workers: the VoiceChat component (`voice-proxy.illmadecoder.workers.dev`) proxies WebRTC sessions between the browser and OpenAI's Realtime API. This ADR extends that pattern to terminal access and ephemeral web pages.

### Requirements

| Requirement | Priority | Notes |
|-------------|----------|-------|
| Browser-based terminal (no local kubectl) | Must | Users should experience tutorials without installing anything |
| Ephemeral web UIs accessible from public internet | Must | Tutorial apps (LDAP admin, Keycloak, web apps) must be reachable |
| Zero permanent public exposure | Must | Nothing exposed when experiments aren't running |
| No changes to hub cluster network (no port forwarding) | Must | Hub is on a home LAN behind NAT |
| Minimal infrastructure cost | High | No always-on VMs or load balancers |
| Secure — no unauthenticated access to cluster | High | Terminal = cluster access, must be gated |
| Works with static site (no SSR) | High | Site is Astro static on GitHub Pages |
| Service URLs discoverable by the site | Medium | Site should render clickable links when services are live |

### Constraints

- **WSL ~8GB RAM** — No Docker, no local builds. All infrastructure runs on the Talos hub cluster or GKE target clusters.
- **Home LAN** — Hub cluster is behind NAT at `192.168.1.178`. No public IP, no port forwarding. Outbound-only connectivity.
- **Static site** — GitHub Pages serves pre-built HTML/JS. No server-side rendering, no dynamic URL injection at request time.
- **Ephemeral experiments** — Target clusters are provisioned on-demand and torn down after completion. URLs are transient.

## Options Considered

### Axis 1: Terminal Access

| Option | How it works | Security | Complexity |
|--------|-------------|----------|------------|
| **ttyd + cloudflared sidecar** | ttyd pod provides web terminal (HTTP/WebSocket). cloudflared sidecar creates a Cloudflare Tunnel, exposing it as `https://<random>.cfargotunnel.com` or a named route. Browser connects directly. | Tunnel auth (Cloudflare Access or token) | Low — two containers, zero ingress config |
| **ttyd + Cloudflare Worker proxy** | Same as above, but a Cloudflare Worker (`terminal-proxy.illmadecoder.workers.dev`) sits in front, adding auth and stable URL. Browser connects to Worker, Worker proxies WebSocket to tunnel. | Worker validates token before proxying | Medium — Worker + tunnel |
| **code-server (VS Code in browser)** | Full VS Code instance with integrated terminal. More than just a terminal — file editing, extensions. | Same tunnel pattern | High — heavier resource footprint (1-2 GB RAM), overkill for tutorial commands |
| **Teleport / Boundary** | Enterprise-grade zero-trust access. Session recording, RBAC, audit. | Excellent | Very high — requires Teleport cluster, overkill for a lab |
| **kubectl proxy + WebSocket adapter** | `kubectl exec` piped through a WebSocket adapter. No ttyd needed. | Kubernetes RBAC | Medium — custom adapter needed |

### Axis 2: Ephemeral Web UI Access

| Option | How it works | URL stability | Complexity |
|--------|-------------|---------------|------------|
| **cloudflared per-service** | Each tutorial component gets a cloudflared sidecar. Each creates its own tunnel. URLs are random unless named routes are configured. | Random unless configured | Low per-service, but N tunnels for N services |
| **Single cloudflared ingress** | One cloudflared pod with an ingress config mapping paths to multiple services. All services share one tunnel with path-based routing. | Stable paths under one domain | Medium — ingress config, path rewriting |
| **Cloudflare Tunnel named routes** | cloudflared registers named routes (`ldap-admin.illmadecoder.com`, `app-portal.illmadecoder.com`) via Cloudflare dashboard or API. Stable, human-readable URLs. | Stable, DNS-based | Medium — requires Cloudflare DNS zone |
| **Tailscale Funnel** | Tailscale's built-in tunnel feature. Each service gets a `https://<node>.ts.net/<path>` URL. | Stable per-node | Low — but ties to Tailscale account, not Cloudflare ecosystem |
| **Ngrok** | Per-service ngrok tunnels. Random URLs on free tier, stable on paid. | Random (free) / Stable (paid) | Low — but adds another vendor |

### Axis 3: Service Discovery (how the site knows URLs)

| Option | How it works | Latency | Complexity |
|--------|-------------|---------|------------|
| **Static configuration** | Tutorial JSON includes hardcoded URLs (named Cloudflare routes). Site always renders them; services are only reachable when experiment is running. | Zero | Minimal — but shows broken links when experiment is down |
| **Operator publishes URLs** | Operator writes service URLs to experiment status. A webhook or polling mechanism updates the site. | Minutes | High — requires dynamic site or API |
| **Status API endpoint** | A lightweight API (Cloudflare Worker or hub-hosted) that returns current experiment status including live service URLs. Site fetches on page load. | Seconds | Medium — new API surface |
| **Convention-based** | URLs follow a predictable pattern (`https://{service}-{experiment}.illmadecoder.com`). Site constructs them from experiment metadata. No discovery needed. | Zero | Low — requires Cloudflare DNS setup |

### Axis 4: Authentication

| Option | How it works | UX | Security |
|--------|-------------|-----|----------|
| **Cloudflare Access (Zero Trust)** | Users authenticate via Cloudflare's identity provider (GitHub OAuth, email OTP, etc.) before reaching any tunneled service. | Login prompt before first access | Strong — enterprise-grade, per-request auth |
| **Shared secret (URL token)** | Tutorial page includes a time-limited token in service URLs. Worker or cloudflared validates it. | Seamless (token in URL) | Moderate — token leakage risk via URL sharing |
| **No auth (accept the risk)** | Tunnels are unprotected. Security relies on obscurity (random tunnel URLs) and short experiment lifetimes. | Seamless | Weak — acceptable only if terminal has restricted RBAC |
| **mTLS client certificates** | cloudflared requires client cert. Site provides cert download or embeds it. | Complex — cert installation | Strong — but terrible UX for tutorials |

## Decision

**ttyd + single cloudflared ingress with named routes + convention-based URLs + Cloudflare Access for terminal, no auth for web UIs.**

### Architecture

```
Browser (benchmark site)
  │
  ├── Terminal: wss://terminal.illmadecoder.com
  │     │
  │     ▼
  │   Cloudflare Access (GitHub OAuth gate)
  │     │
  │     ▼
  │   cloudflared tunnel (hub cluster)
  │     │
  │     ▼
  │   ttyd pod (port 7681, bash + kubectl + kubeconfig)
  │
  ├── LDAP Admin: https://ldap-admin.illmadecoder.com
  │     │
  │     ▼
  │   cloudflared tunnel (same tunnel, path-based route)
  │     │
  │     ▼
  │   openldap-admin Service (port 443)
  │
  ├── App Portal: https://app-portal.illmadecoder.com
  │     ...same tunnel pattern...
  │
  └── Keycloak: https://keycloak.illmadecoder.com
        ...same tunnel pattern...
```

### Why This Combination

| Factor | Reasoning |
|--------|-----------|
| **ttyd for terminal** | Purpose-built web terminal. Lightweight (single binary, ~20 MB image). WebSocket-native. No custom adapter code. Battle-tested. |
| **Single cloudflared ingress** | One tunnel process manages all service routes. Simpler than N sidecars. Ingress config is a YAML file — GitOps-friendly. |
| **Named routes** | Predictable URLs: `terminal.illmadecoder.com`, `ldap-admin.illmadecoder.com`. Site can hardcode them in tutorial JSON. No discovery API needed. |
| **Convention-based discovery** | Tutorial JSON includes service URLs at authoring time. URLs are deterministic (named routes). No runtime discovery. Works with static site. |
| **Cloudflare Access for terminal** | Terminal = kubectl access = cluster admin. Must be gated. GitHub OAuth is zero-config for a solo lab. Free tier supports 50 users. |
| **No auth for web UIs** | Tutorial apps (LDAP admin, web apps, Keycloak admin) contain synthetic test data only. Risk of exposure is negligible. Auth would add friction to the tutorial experience. The tunnel only exists while the experiment is running. |

### Why Not the Others

| Option | Reason to reject |
|--------|-----------------|
| **ttyd + Worker proxy** | Extra hop adds latency to terminal sessions. Worker WebSocket proxying is complex (must handle binary frames, connection lifecycle). Cloudflare Access provides the auth layer without custom code. |
| **code-server** | 1-2 GB RAM for VS Code is excessive when users only need a terminal. Tutorial commands are `kubectl`, `ldapsearch`, `curl` — no file editing. |
| **Teleport / Boundary** | Enterprise access management for a solo lab tutorial. Massive overkill. |
| **cloudflared per-service sidecars** | N tunnel processes for N services. More resource usage, more config, harder to manage. The single-ingress model is simpler. |
| **Tailscale Funnel** | Ties to Tailscale ecosystem. Lab already uses Cloudflare for DNS and Workers (VoiceChat proxy). Keeping one vendor for tunneling is simpler. |
| **Operator publishes URLs** | Requires dynamic content on a static site, or a polling mechanism. Named routes eliminate the need — URLs are known at authoring time. |
| **Status API** | Adds a new API surface to build and maintain. For V1, convention-based URLs with a "requires running experiment" fallback are sufficient. |
| **mTLS for auth** | Terrible UX for tutorials. Users would need to download and install client certificates. |

### Component Architecture

#### ttyd Pod

```yaml
# components/apps/tutorial-terminal/component.yaml
apiVersion: experiments.illm.io/v1alpha1
kind: Component
metadata:
  name: tutorial-terminal
spec:
  category: apps
  sources:
    - type: manifests
      path: components/apps/tutorial-terminal/manifests/
  params:
    - name: experiment-namespace
      description: Namespace to scope kubectl access
```

Deployment:
- **Container 1: ttyd** — `tsl0922/ttyd:latest` running `bash` with `kubectl` and `kubeconfig` mounted (read-only). Scoped to experiment namespace via RBAC (ServiceAccount with Role, not ClusterRole).
- **Container 2: cloudflared** — `cloudflare/cloudflared:latest` running `tunnel --config /etc/cloudflared/config.yml`. Config maps service routes to backend pods.
- **Service:** ClusterIP on port 7681 (ttyd). cloudflared handles external access.
- **RBAC:** ServiceAccount with a namespaced Role granting `get`, `list`, `watch` on pods, services, deployments, and `create` for `pods/exec`. No cluster-wide access.

#### cloudflared Tunnel Configuration

```yaml
# Tunnel config (mounted as ConfigMap)
tunnel: lab-tutorial
credentials-file: /etc/cloudflared/credentials.json

ingress:
  - hostname: terminal.illmadecoder.com
    service: http://tutorial-terminal:7681
    originRequest:
      noTLSVerify: true
  - hostname: ldap-admin.illmadecoder.com
    service: https://openldap-admin:443
  - hostname: app-portal.illmadecoder.com
    service: http://id-app-portal:8080
  - hostname: app-dashboard.illmadecoder.com
    service: http://id-app-dashboard:8080
  - hostname: app-wiki.illmadecoder.com
    service: http://id-app-wiki:8080
  - hostname: keycloak.illmadecoder.com
    service: https://keycloak:8443
  - service: http_status:404
```

#### Cloudflare Access Policy

```
Application: Lab Tutorial Terminal
Domain: terminal.illmadecoder.com
Policy: Allow — GitHub OAuth (restrict to repo owner's GitHub account)
Session duration: 24 hours
```

Web UI hostnames (`ldap-admin`, `app-*`, `keycloak`) are NOT behind Cloudflare Access — they are directly accessible through the tunnel when it's running.

#### DNS Configuration

CNAME records in Cloudflare DNS (created once, persist across experiments):

```
terminal.illmadecoder.com        → <tunnel-id>.cfargotunnel.com
ldap-admin.illmadecoder.com      → <tunnel-id>.cfargotunnel.com
app-portal.illmadecoder.com      → <tunnel-id>.cfargotunnel.com
app-dashboard.illmadecoder.com   → <tunnel-id>.cfargotunnel.com
app-wiki.illmadecoder.com        → <tunnel-id>.cfargotunnel.com
keycloak.illmadecoder.com        → <tunnel-id>.cfargotunnel.com
```

All point to the same tunnel. cloudflared routes by hostname.

### Site Integration

#### WebTerminal Component (already built)

The `WebTerminal.astro` component (created in this branch) renders an expandable xterm.js terminal panel. It accepts a `terminalUrl` prop:

- **When URL is provided:** xterm.js connects via WebSocket, user interacts with ttyd.
- **When no URL:** Shows "Requires running experiment" message with instructions to start the experiment.

The terminal URL is static and known at build time: `wss://terminal.illmadecoder.com/ws`. Since ttyd uses WebSocket natively and Cloudflare Tunnel preserves WebSocket connections, no custom Worker proxy is needed.

#### Service Links

Tutorial JSON (`site/data/{name}.tutorial.json`) includes service definitions with URLs:

```json
{
  "services": [
    { "name": "terminal", "label": "Web Terminal", "icon": "terminal", "url": "wss://terminal.illmadecoder.com/ws" },
    { "name": "ldap-admin", "label": "LDAP Admin", "icon": "database", "url": "https://ldap-admin.illmadecoder.com" },
    { "name": "app-portal", "label": "App Portal", "icon": "globe", "url": "https://app-portal.illmadecoder.com" }
  ]
}
```

The `ServiceLinksPanel.astro` component renders these as clickable buttons. If the experiment isn't running, the links simply won't resolve (Cloudflare returns 502). The tutorial instructions explain this.

### Credential Management

#### Tunnel Credentials

cloudflared authenticates to Cloudflare using a tunnel token or credentials file. This is a one-time setup:

```bash
# On a machine with cloudflared installed (not the cluster):
cloudflared tunnel create lab-tutorial
# Produces a credentials JSON file

# Store in OpenBao:
kubectl exec -n openbao openbao-0 -- sh -c \
  "BAO_TOKEN='<root_token>' bao kv put secret/experiments/cloudflare-tunnel \
  credentials='$(cat ~/.cloudflared/<tunnel-id>.json)'"
```

ExternalSecret syncs to the experiment namespace as a K8s Secret. The cloudflared pod mounts it.

| OpenBao Path | K8s Secret | Namespace | Purpose |
|-------------|------------|-----------|---------|
| `secret/experiments/cloudflare-tunnel` | `cloudflare-tunnel-credentials` | `experiments` | Tunnel credentials JSON |

#### ttyd RBAC

ttyd runs as a ServiceAccount with minimal permissions:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: tutorial-terminal
  namespace: ${EXPERIMENT_NAMESPACE}
rules:
  - apiGroups: [""]
    resources: [pods, services, configmaps, secrets]
    verbs: [get, list, watch]
  - apiGroups: [""]
    resources: [pods/exec]
    verbs: [create]
  - apiGroups: [apps]
    resources: [deployments, statefulsets]
    verbs: [get, list, watch]
```

This gives enough access to run tutorial commands (`kubectl get pods`, `kubectl exec openldap-0 -- ldapsearch ...`) without cluster-admin privileges.

## Implementation Phases

### Phase 1: Tunnel Infrastructure (prerequisite)

1. Create Cloudflare Tunnel (`lab-tutorial`) via `cloudflared tunnel create`
2. Configure DNS CNAME records in Cloudflare dashboard
3. Set up Cloudflare Access application for `terminal.illmadecoder.com`
4. Store tunnel credentials in OpenBao, create ExternalSecret
5. Create `components/apps/tutorial-terminal/` with ttyd + cloudflared manifests

### Phase 2: Integration Testing

1. Deploy ttyd + cloudflared manually in the experiments namespace
2. Verify terminal access via `https://terminal.illmadecoder.com`
3. Verify Cloudflare Access gate prompts for GitHub OAuth
4. Test WebSocket connectivity from xterm.js (WebTerminal component)

### Phase 3: Experiment Integration

1. Add `tutorial-terminal` as a component in identity tutorial experiment YAMLs
2. Add service URLs to tutorial JSON files
3. Test full flow: create experiment → tunnel comes up → site links work → terminal accessible

### Phase 4: Hardening

1. Add health check endpoint to ttyd pod
2. Add tunnel status monitoring (cloudflared metrics)
3. Add session timeout / idle disconnect to ttyd
4. Consider rate limiting on Cloudflare Access

## Consequences

### Positive

- Browser-based terminal access with zero client-side installation
- Ephemeral exposure — services are only reachable while experiment is running
- Zero permanent public attack surface (tunnels are outbound-only)
- Convention-based URLs work with static site — no dynamic content needed
- Single tunnel process for all services — simple to operate
- Cloudflare Access for terminal provides enterprise-grade auth at zero cost
- Reuses existing Cloudflare infrastructure (DNS, Workers already in use for VoiceChat)
- GitOps-friendly — tunnel config is a YAML ConfigMap

### Negative

- Dependency on Cloudflare for tunneling (vendor lock-in for this feature)
- Named routes require a Cloudflare-managed DNS zone (`illmadecoder.com`)
- Terminal shows "connection refused" when experiment isn't running (mitigate: clear messaging in tutorial UI)
- Cloudflare Access adds a login step for terminal access (mitigate: 24h session, only applies to terminal, not web UIs)
- ttyd has no built-in session recording or audit trail (mitigate: sufficient for a lab; add if needed later)
- cloudflared pod is an extra resource (~50 MB RAM) in every tutorial experiment

### Security Considerations

| Surface | Risk | Mitigation |
|---------|------|------------|
| Terminal (kubectl) | Cluster access via browser | Cloudflare Access (GitHub OAuth), namespaced RBAC, no cluster-admin |
| Web UIs (LDAP admin, Keycloak) | Exposed admin consoles | Synthetic test data only, tunnel only exists during experiment, short-lived |
| Tunnel credentials | If leaked, attacker can route traffic through tunnel | Stored in OpenBao, synced via ESO, tunnel can be rotated |
| ttyd process | Command injection, shell escape | ttyd runs as non-root, RBAC limits kubectl scope, no host mounts |

### Future

- **Multi-experiment tunnels:** If multiple tutorials run simultaneously, each needs its own tunnel or the ingress config must be dynamic. For V1, only one tutorial runs at a time.
- **Session recording:** Add asciinema recording to ttyd sessions for tutorial replay / verification.
- **Status API:** If convention-based URLs prove insufficient (e.g., need to show "experiment running" indicator), add a lightweight status endpoint on a Cloudflare Worker.
- **Operator-managed tunnels:** The experiment operator could create/destroy tunnel routes as part of the experiment lifecycle, eliminating manual DNS setup for new tutorials.
- **Persistent terminal history:** Mount a PVC for bash history so users can resume where they left off.

## References

- [ADR-015: Experiment Operator](ADR-015-experiment-operator.md) — Experiment lifecycle and component deployment
- [ADR-017: Benchmark Results Site](ADR-017-benchmark-results-site.md) — Static site architecture, VoiceChat Worker pattern
- [Identity Series Roadmap](../roadmap/appendix-identity-auth.md) — Tutorial UX model (static mode vs live mode)
- [Cloudflare Tunnel Documentation](https://developers.cloudflare.com/cloudflare-one/connections/connect-networks/)
- [Cloudflare Access Documentation](https://developers.cloudflare.com/cloudflare-one/policies/access/)
- [ttyd — Share your terminal over the web](https://github.com/tsl0922/ttyd)
- [xterm.js — Terminal for the web](https://xtermjs.org/)
- VoiceChat component: `site/src/components/VoiceChat.astro` (existing Cloudflare Worker proxy pattern)
- WebTerminal component: `site/src/components/WebTerminal.astro` (xterm.js integration, built in this branch)
