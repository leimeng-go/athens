import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  build: {
    outDir: '../pkg/admin/static',
    emptyOutDir: true
  },
  server: {
    proxy: {
      '/admin/api': {
        target: 'http://localhost:3000',
        changeOrigin: true
      }
    }
  }
})