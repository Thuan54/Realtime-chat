import { execSync } from 'child_process'
import path from 'path'

export default async () => {
  console.log('Tearing down test environment...')
  const infraPath = path.resolve(__dirname, '../../infra/docker-compose.yml')
  
  // Remove containers and anonymous volumes to ensure clean state for next run
  execSync(`docker compose -f ${infraPath} down -v`, { stdio: 'inherit' })
  console.log('Teardown complete')
}