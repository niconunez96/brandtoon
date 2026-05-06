import tailwindcss from '@tailwindcss/vite'
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vitest/config'

export default defineConfig(() => {
  const isTest = Boolean(process.env.VITEST)

  return {
    plugins: [react(), ...(isTest ? [] : [tailwindcss()])],
    server: {
      port: 5173,
    },
    test: {
      globals: true,
      environment: 'jsdom',
      setupFiles: './src/test/setup.ts',
      coverage: {
        provider: 'v8',
        reporter: ['text', 'html'],
      },
    },
  }
})
