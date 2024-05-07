import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

console.log('Environment Variables:', import.meta.env)
// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: '0.0.0.0',
    port: 8000
  },
  build: {
    minify: 'terser', // 是否进行压缩,boolean | 'terser' | 'esbuild',默认使用terser
    manifest: false, // 是否产出manifest.json
    sourcemap: false, // 是否产出sourcemap.json
    outDir: 'dist' // 产出目录
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@use "@/style/element/index.scss" as *;`
      }
    }
  }
})
