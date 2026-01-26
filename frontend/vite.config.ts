import path from "path"
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    host: '127.0.0.1',  // 明确绑定 127.0.0.1 以支持 IDE 端口转发
    port: 51733,
    strictPort: true,  // 端口被占用时报错
  },
})
