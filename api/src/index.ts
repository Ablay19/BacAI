import { Hono } from 'hono'
import { cors } from 'hono/cors'
import { logger } from 'hono/logger'
import { zValidator } from '@hono/zod-validator'
import { z } from 'zod'

// Import route handlers
import { solveRouter } from './routes/solve'
import { explainRouter } from './routes/explain'
import { converseRouter } from './routes/converse'
import { dataRouter } from './routes/data'
import { healthRouter } from './routes/health'

const app = new Hono()

// Middleware
app.use('*', cors({
  origin: [
    'http://localhost:3000',
    'https://bacai.vercel.app',
    'https://bacai-*.vercel.app'
  ],
  allowMethods: ['GET', 'POST', 'PUT', 'DELETE', 'OPTIONS'],
  allowHeaders: ['Content-Type', 'Authorization'],
}))

app.use('*', logger())

// Rate limiting middleware
app.use('*', async (c, next) => {
  const ip = c.req.header('CF-Connecting-IP') || c.req.header('X-Forwarded-For') || 'unknown'
  const key = `rate_limit:${ip}`
  
  // Simple rate limiting using Cloudflare KV (if available)
  // For now, we'll just pass through
  
  await next()
})

// Root endpoint
app.get('/', (c) => {
  return c.json({
    name: 'BACAI API',
    version: '1.0.0',
    description: 'Mauritanian AI Educational System API',
    endpoints: {
      solve: '/api/solve',
      explain: '/api/explain',
      converse: '/api/converse',
      data: '/api/data',
      health: '/api/health'
    },
    supported_languages: ['ar', 'fr', 'en'],
    supported_subjects: [
      'mathematics', 'arabic', 'french', 'english', 
      'sciences', 'islamic_studies'
    ]
  })
})

// Route handlers
app.route('/api/solve', solveRouter)
app.route('/api/explain', explainRouter)
app.route('/api/converse', converseRouter)
app.route('/api/data', dataRouter)
app.route('/api/health', healthRouter)

// Error handling
app.onError((err, c) => {
  console.error('Error:', err)
  
  return c.json({
    success: false,
    error: 'Internal Server Error',
    message: err.message,
    code: 'INTERNAL_ERROR'
  }, 500)
})

// 404 handling
app.notFound((c) => {
  return c.json({
    success: false,
    error: 'Endpoint not found',
    code: 'NOT_FOUND'
  }, 404)
})

export default {
  fetch: app.fetch,
}