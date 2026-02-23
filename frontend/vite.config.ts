import react from '@vitejs/plugin-react'
import path from "path"
import { defineConfig } from 'vite'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    host: '127.0.0.1',
    port: 51733,
    strictPort: true,
    proxy: {
      // 只代理后端 API 路径
      '/admin/api': {
        target: 'http://127.0.0.1:3000',
        changeOrigin: true,
      }
    },
  },
})
