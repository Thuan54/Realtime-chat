import { execSync } from 'child_process'
import path from 'path'

export default async () => {
  console.log('Starting containerized test environment...')

  const infraPath = path.resolve(__dirname, '../../infra/docker-compose.yml')
  
  // Start infrastructure dependencies
  execSync(`docker compose -f ${infraPath} up -d postgres redis`, { stdio: 'inherit' })
  
  // Apply migrations and seed test data
  execSync('make migrate-up seed', { cwd: path.resolve(__dirname, '../..'), stdio: 'inherit' })
  
  // Start application services
  execSync(`docker compose -f ${infraPath} up -d backend frontend nginx`, { stdio: 'inherit' })

  // Wait for backend health endpoint
  let retries = 30
  while (retries--) {
    try {
      execSync('curl -f http://localhost:80/api/v1/health', { stdio: 'ignore' })
      console.log('Environment ready for E2E tests')
      return
    } catch {
      await new Promise(resolve => setTimeout(resolve, 2000))
    }
  }
  throw new Error('Environment failed to become ready within timeout')
}