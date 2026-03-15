import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

export default {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      pages: '../cmd/server/dist',
      assets: '../cmd/server/dist',
      fallback: 'index.html',
      precompress: false,
    }),
  },
};
