import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '^/api': {
        target: 'http://host.docker.internal:8039',
        ws: true,
        secure: false,
        changeOrigin: true,
      },
      '^/media': {
        target: 'http://host.docker.internal:8039',
        ws: true,
        secure: false,
        changeOrigin: true,
      },
    }
  }
})
