import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  build: {
    emptyOutDir: false,
    rollupOptions: {
      input: {
        "": "index.html",
        "kleo/index.html": "kleo.html",
      }
    }
  },
  // root: "/home/rosta/skrat-org/web/",
  plugins: [svelte()],
})
