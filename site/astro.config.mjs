import { defineConfig } from 'astro/config';
import tailwind from '@astrojs/tailwind';

export default defineConfig({
  site: 'https://illmadecoder.github.io',
  base: '/k8s-ai-cloud-testbed/',
  trailingSlash: 'always',
  integrations: [tailwind()],
});
