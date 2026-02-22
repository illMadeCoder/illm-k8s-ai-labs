# ADR-019: AI Chat & Voice Assistant

## Status

Accepted (2026-02-15)

## Context

The benchmark site (ADR-017) publishes experiment results as static HTML pages on GitHub Pages. Visitors can read analysis summaries, inspect charts, and browse experiment metadata — but there is no way to ask questions about what they're looking at. A first-time visitor seeing a TSDB comparison has no easy path to understanding why VictoriaMetrics outperformed Prometheus in a specific metric, or what the experiment methodology was, or how the operator works under the hood.

The site is a portfolio piece. Making it conversational — letting visitors ask questions and get contextual answers about the experiment they're viewing — adds significant engagement value. Voice interaction further differentiates it from typical static documentation.

### Requirements

| Requirement | Priority | Notes |
|-------------|----------|-------|
| Text chat on every page | Must | Ask questions about the experiment being viewed |
| Voice input/output | High | Speak questions, hear answers — differentiator |
| Page-aware context | Must | AI knows which experiment/series the user is viewing |
| Source code browsing | High | AI can look up operator code, configs, component definitions |
| Zero backend infrastructure | Must | No servers to run — static site on GitHub Pages |
| Cost-controlled | Must | Public site, no auth — must cap API spend |
| Works without auth | Must | No login required for visitors |

### Constraints

- **Static site** — GitHub Pages, no SSR, no server-side API routes.
- **API key protection** — OpenAI API key cannot be in client-side JavaScript.
- **Budget** — OpenAI Realtime API charges per session. Uncontrolled public access would be expensive.

## Decision

**OpenAI Realtime API via WebRTC, proxied through a Cloudflare Worker for key protection and rate limiting.**

### Architecture

```
Browser (every page)
  │
  ├── Text mode: type question → WebRTC data channel → OpenAI Realtime API
  │
  └── Voice mode: microphone → WebRTC audio track → OpenAI Realtime API
                                                        │
                                                   ◄── audio + text response
Session bootstrap:
  Browser ──POST /session──► Cloudflare Worker (voice-proxy.illmadecoder.workers.dev)
                                   │
                                   ├── CORS check (allowed origins only)
                                   ├── Rate limit check (KV: 40 sessions/month)
                                   └── POST /v1/realtime/sessions ──► OpenAI API
                                         │
                                         ▼
                              Ephemeral client token (short-lived)
                                         │
                                         ▼
  Browser ◄── { client_secret } ◄──────────
  Browser ──WebRTC SDP offer──► api.openai.com/v1/realtime (direct, using ephemeral token)
  Browser ◄── SDP answer ◄──── OpenAI
  [WebRTC session established — audio + data channel, peer-to-peer with OpenAI]
```

The key insight: the Cloudflare Worker only handles session creation (one HTTP request). The actual WebRTC session runs directly between the browser and OpenAI — no proxy in the data path. This means zero latency overhead for the conversation itself.

### Components

#### Cloudflare Worker (`workers/voice-proxy/`)

**Endpoint:** `POST https://voice-proxy.illmadecoder.workers.dev/session`

Responsibilities:
1. **CORS enforcement** — Only allows requests from `illmadecoder.github.io` and `localhost` (dev).
2. **Rate limiting** — Monthly session budget (40/month) tracked in Cloudflare KV. Prevents runaway API costs from a public site.
3. **Key protection** — OpenAI API key stored as a Worker secret (`OPENAI_API_KEY`), never sent to the browser.
4. **Ephemeral token exchange** — Calls OpenAI's `/v1/realtime/sessions` to create a short-lived client token, returns it to the browser.

```
workers/voice-proxy/
├── src/index.ts          # ~95 lines: CORS, rate limit, session creation
├── wrangler.toml         # Worker config, KV namespace binding
├── package.json
└── tsconfig.json
```

Deployed via `npx wrangler deploy`. KV namespace stores monthly counters with 31-day TTL.

#### VoiceChat Component (`site/src/components/VoiceChat.astro`)

A single Astro component (~660 lines, all `<script is:inline>` + `<style>`) embedded in the Base layout. Renders on every page.

**UI:**
- Fixed-position panel at bottom of viewport, initially minimized
- Draggable resize handle (10–85vh)
- Message bubbles (AI + user) with monospace styling
- Text input with send button
- Voice toggle button (microphone) with voice selector (8 OpenAI voices)
- Animated gradient border when AI is responding
- Mini visualization dots (listening/speaking states)

**Behavior:**
- **Lazy connection** — No WebRTC session until the user sends their first message or clicks voice. Zero cost for page views without interaction.
- **Text mode** — Creates WebRTC connection with a silent audio track. Messages sent via data channel. Responses streamed as text deltas.
- **Voice mode** — Requests microphone access, sends real audio track. Server VAD (voice activity detection) handles turn-taking. Responses arrive as both audio and text transcript.
- **Mode switching** — Can switch between text and voice mid-conversation without reconnecting. Swaps the audio track and updates session modalities.
- **Auto-reconnect** — On WebRTC disconnection, automatically re-establishes the session in the same mode.

**Context injection:**
- On data channel open, sends a `session.update` with system instructions containing:
  - Lab preamble (what the testbed is, tech stack, methodology)
  - Current page topic (experiment name, series name)
  - Page content (scraped from `#page-markdown` element or `pageContext` prop)
  - Brevity instruction ("Keep answers brief — a few sentences")
  - Tool definitions for GitHub repo browsing

**Tool use (GitHub repo browsing):**
- Two tools registered with the Realtime API session:
  - `get_file` — Reads a file from the GitHub repo via the public Contents API
  - `list_directory` — Lists directory contents
- When the AI calls a tool, the component fetches from `api.github.com`, sends the result back via data channel, and triggers a new response.
- This lets visitors ask "how does the operator handle metrics collection?" and get answers grounded in actual source code.

#### Base Layout Integration

`VoiceChat` is included in `Base.astro` (the root layout) with per-page props:

```astro
<VoiceChat
  pageTopic={pageTopic}
  openingMessage={openingMessage}
  pageContext={pageContext}
  placeholder={chatPlaceholder}
/>
```

Every page passes its topic and an opening message. Experiment detail pages include the hypothesis and abstract. The AI's first message is pre-rendered (no API call) to give immediate context.

### Cost Control

| Control | Mechanism |
|---------|-----------|
| Monthly budget cap | Worker KV counter, 40 sessions/month, returns 429 when exceeded |
| CORS origin check | Only accepts requests from the site's domain |
| Lazy connection | No session created until user interaction |
| No auth required | Deliberate — portfolio site, friction-free experience |
| Ephemeral tokens | Short-lived, cannot be reused after expiration |

At ~$0.06/min for Realtime API audio, 40 sessions averaging 3 minutes each = ~$7.20/month worst case. Text-only sessions are cheaper (no audio processing).

### Why This Approach

| Factor | Reasoning |
|--------|-----------|
| **OpenAI Realtime API** | Native voice support with server-side VAD. No speech-to-text + LLM + text-to-speech pipeline to build. Sub-second voice latency. |
| **WebRTC (not WebSocket)** | OpenAI's Realtime API uses WebRTC for lowest latency. Audio streams peer-to-peer. Data channel for text + events. |
| **Cloudflare Worker** | Serverless, globally distributed, free tier covers this usage. Keeps API key server-side. ~5ms cold start. |
| **KV rate limiting** | Simple, durable, no external DB. Monthly key with TTL = self-cleaning. |
| **Inline script** | `<script is:inline>` avoids Astro's module bundling. Single self-contained IIFE. No framework dependencies (no React, no state management). Vanilla JS keeps the bundle minimal. |
| **GitHub API for tools** | Public API, no auth needed for public repos. Lets the AI ground answers in actual source code. 60 req/hr unauthenticated rate limit is sufficient for conversational use. |

### Why Not Alternatives

| Alternative | Reason to reject |
|-------------|-----------------|
| **Claude API** | No native voice/WebRTC support at the time of implementation. Would require separate STT + TTS pipeline. |
| **Self-hosted LLM** | Hub cluster doesn't have GPU. API latency is already sub-second. |
| **WebSocket proxy** | Worker would need to relay all audio/text data, adding latency and bandwidth cost. WebRTC's direct connection is superior. |
| **API key in frontend** | Trivially extractable, no rate limiting possible. |
| **User authentication** | Adds friction to a portfolio site. Rate limiting at the session level is sufficient. |

## Consequences

### Positive

- Every page on the site has an AI assistant that understands the content being viewed
- Voice interaction with sub-second latency — differentiating feature for a portfolio site
- Zero infrastructure to maintain (Cloudflare Worker is serverless)
- Cost-controlled via monthly budget cap
- Source code browsing via tools grounds answers in reality
- No user auth required — zero-friction experience
- Lazy connection means no cost for visitors who don't interact

### Negative

- Dependency on OpenAI Realtime API (vendor lock-in for this feature)
- Monthly budget cap means the feature goes dark after 40 sessions
- GitHub API rate limit (60/hr unauthenticated) could be hit during heavy tool use
- ~660 lines of inline JavaScript in a single component (no framework, harder to test)
- Voice requires HTTPS (WebRTC security requirement) — works on GitHub Pages, not plain HTTP

### Future

- Increase monthly budget as experiment portfolio grows and traffic justifies it
- Add authenticated GitHub API access (PAT in Worker secret) for higher rate limits on tool calls
- Add conversation persistence (localStorage or KV) so users can resume across page navigations
- Consider Claude API if native voice/WebRTC support becomes available
- Add usage analytics (Worker logs) to understand session patterns and optimize budget

## References

- [ADR-017: Benchmark Results Site](ADR-017-benchmark-results-site.md) — Site architecture, static hosting
- [OpenAI Realtime API Documentation](https://platform.openai.com/docs/guides/realtime)
- [Cloudflare Workers Documentation](https://developers.cloudflare.com/workers/)
- Worker source: `workers/voice-proxy/src/index.ts`
- Component source: `site/src/components/VoiceChat.astro`
- Layout integration: `site/src/layouts/Base.astro`
