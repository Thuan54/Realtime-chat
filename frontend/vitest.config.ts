import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'
import path from 'path'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  test: {
    environment: 'jsdom',
    globals: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html'],
      thresholds: {
        lines: 80,
        branches: 80,
        functions: 80,
        statements: 80,
      },
      include: ['src/**/*.{ts,tsx}'],
      exclude: ['src/**/*.d.ts', 'src/main.tsx', 'src/vite-env.d.ts'],
    },
    setupFiles: ['./src/test/setup.ts'],
    mockReset: true,
    restoreMocks: true,
    reporters: ['default', 'junit'],
    outputFile: {
      junit: './test-results/frontend-junit.xml',
    },
    // Deterministic WS/REST mock adapters for unit isolation
    // WS/REST mocks should be imported explicitly in test files via @/test/mocks/*
  },
})