import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// The Go backend has no CORS middleware, so we proxy /api to it in dev.
export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    proxy: { '/api': { target: 'http://localhost:8080', changeOrigin: true } },
  },
})
