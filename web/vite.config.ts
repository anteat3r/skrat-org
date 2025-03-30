import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

import * as child from "child_process";
const commitHash = child.execSync("git rev-parse --short HEAD").toString();

export default defineConfig({
  build: {
    emptyOutDir: false,
    rollupOptions: {
      input: {
        "index.html": "index.html",
        "kleo/index.html": "kleo/index.html",
      }
    }
  },
  define: {
    __COMMIT_HASH__: JSON.stringify(commitHash),
  },
  // root: "/home/rosta/skrat-org/web/",
  plugins: [svelte()],
})
