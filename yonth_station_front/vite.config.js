import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://127.0.0.1:8888',
        changeOrigin: true,
        ws: true,
        configure: (proxy) => {
          proxy.on('proxyRes', (proxyRes, req) => {
            // 流式响应不缓冲，立即转发
            if (req.url.includes('/chat/ask/stream')) {
              proxyRes.headers['cache-control'] = 'no-cache'
              proxyRes.headers['x-no-buffering'] = 'true'
            }
          })
        }
      }
    }
  }
})
