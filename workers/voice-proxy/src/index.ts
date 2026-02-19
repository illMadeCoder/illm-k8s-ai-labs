export interface Env {
	OPENAI_API_KEY: string;
	RATE_LIMIT: KVNamespace;
}

const ALLOWED_ORIGINS = [
	'https://illmadecoder.github.io',
	'http://localhost:4321',
	'http://localhost:3000',
];

const MAX_PER_IP_PER_DAY = 2;
const MAX_PER_MONTH = 40;

async function hashIP(ip: string): Promise<string> {
	const data = new TextEncoder().encode(ip);
	const hash = await crypto.subtle.digest('SHA-256', data);
	return Array.from(new Uint8Array(hash))
		.slice(0, 8)
		.map((b) => b.toString(16).padStart(2, '0'))
		.join('');
}

function corsHeaders(origin: string): Record<string, string> {
	return {
		'Access-Control-Allow-Origin': origin,
		'Access-Control-Allow-Methods': 'POST, OPTIONS',
		'Access-Control-Allow-Headers': 'Content-Type',
	};
}

export default {
	async fetch(request: Request, env: Env): Promise<Response> {
		const origin = request.headers.get('Origin') || '';
		const isAllowed = ALLOWED_ORIGINS.includes(origin);
		const headers = isAllowed ? corsHeaders(origin) : {};

		if (request.method === 'OPTIONS') {
			return new Response(null, { status: 204, headers });
		}

		const url = new URL(request.url);
		if (request.method !== 'POST' || url.pathname !== '/session') {
			return new Response(JSON.stringify({ error: 'Not found' }), {
				status: 404,
				headers: { ...headers, 'Content-Type': 'application/json' },
			});
		}

		if (!isAllowed) {
			return new Response(JSON.stringify({ error: 'Forbidden' }), {
				status: 403,
				headers: { 'Content-Type': 'application/json' },
			});
		}

		// Rate limit: per IP per day
		const ip = request.headers.get('CF-Connecting-IP') || 'unknown';
		const ipHash = await hashIP(ip);
		const today = new Date().toISOString().slice(0, 10);
		const ipKey = `ip:${ipHash}:${today}`;

		const ipCount = parseInt((await env.RATE_LIMIT.get(ipKey)) || '0', 10);
		if (ipCount >= MAX_PER_IP_PER_DAY) {
			return new Response(
				JSON.stringify({ error: 'Daily limit reached. Try again tomorrow.' }),
				{ status: 429, headers: { ...headers, 'Content-Type': 'application/json' } }
			);
		}

		// Rate limit: monthly budget
		const month = new Date().toISOString().slice(0, 7);
		const budgetKey = `budget:${month}`;
		const budgetCount = parseInt((await env.RATE_LIMIT.get(budgetKey)) || '0', 10);
		if (budgetCount >= MAX_PER_MONTH) {
			return new Response(
				JSON.stringify({ error: 'Monthly budget exceeded. Voice chat is temporarily unavailable.' }),
				{ status: 429, headers: { ...headers, 'Content-Type': 'application/json' } }
			);
		}

		// Create ephemeral token via OpenAI Realtime API
		try {
			const response = await fetch('https://api.openai.com/v1/realtime/sessions', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${env.OPENAI_API_KEY}`,
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					model: 'gpt-4o-realtime-preview',
					voice: 'verse',
				}),
			});

			if (!response.ok) {
				const text = await response.text();
				return new Response(
					JSON.stringify({ error: 'Failed to create session', detail: text }),
					{ status: 502, headers: { ...headers, 'Content-Type': 'application/json' } }
				);
			}

			// Increment rate limits
			await env.RATE_LIMIT.put(ipKey, String(ipCount + 1), { expirationTtl: 86400 });
			await env.RATE_LIMIT.put(budgetKey, String(budgetCount + 1), { expirationTtl: 2678400 });

			const data = await response.json();
			return new Response(JSON.stringify(data), {
				status: 200,
				headers: { ...headers, 'Content-Type': 'application/json' },
			});
		} catch {
			return new Response(JSON.stringify({ error: 'Internal error' }), {
				status: 500,
				headers: { ...headers, 'Content-Type': 'application/json' },
			});
		}
	},
};
