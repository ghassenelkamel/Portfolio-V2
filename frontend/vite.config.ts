import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig(() => {
  const target = process.env.VITE_API_PROXY_TARGET || 'http://localhost:8080';

  return {
    plugins: [sveltekit()],
    server: {
      host: '0.0.0.0',
      port: 5173,
      proxy: {
        '/api': {
          target,
          changeOrigin: true,
          secure: false
        }
      }
    }
  };
});
