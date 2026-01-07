import { Hono } from 'hono'
import { zValidator } from '@hono/zod-validator'
import { z } from 'zod'
import { callModelService } from '../utils/model-client'
import { detectLanguage } from '../utils/language-utils'

const explainSchema = z.object({
  concept: z.string().min(1, 'Concept is required'),
  subject: z.enum(['mathematics', 'arabic', 'french', 'english', 'sciences', 'islamic_studies']),
  level: z.enum(['secondary_basic', 'secondary_lycee', 'university']).default('secondary_lycee'),
  language: z.enum(['ar', 'fr', 'en']).optional(),
  detail_level: z.enum(['basic', 'detailed', 'comprehensive']).default('detailed'),
  examples: z.boolean().default(true),
  cultural_context: z.boolean().default(true)
})

const explainRouter = new Hono()

explainRouter.post('/', zValidator('json', explainSchema), async (c) => {
  try {
    const startTime = Date.now()
    const validatedData = c.req.valid('json')
    
    // Auto-detect language if not provided
    const detectedLanguage = validatedData.language || detectLanguage(validatedData.concept)
    
    // Prepare payload for model service
    const payload = {
      concept: validatedData.concept,
      subject: validatedData.subject,
      level: validatedData.level,
      language: detectedLanguage,
      detail_level: validatedData.detail_level,
      examples: validatedData.examples,
      cultural_context: validatedData.cultural_context,
      request_type: 'explain'
    }

    // Call model service
    const modelResponse = await callModelService(payload)
    
    if (!modelResponse.success) {
      return c.json({
        success: false,
        error: modelResponse.error || 'Model service unavailable',
        code: 'MODEL_ERROR'
      }, 503)
    }

    const processingTime = Date.now() - startTime

    // Format response
    const response = {
      success: true,
      data: {
        explanation: modelResponse.data.explanation,
        concept: validatedData.concept,
        subject: validatedData.subject,
        level: validatedData.level,
        language: detectedLanguage,
        detail_level: validatedData.detail_level,
        metadata: {
          processing_time: processingTime,
          model_used: modelResponse.metadata?.model_used || 'unknown',
          language_detected: detectedLanguage,
          timestamp: new Date().toISOString()
        }
      }
    }

    return c.json(response)

  } catch (error) {
    console.error('Explain endpoint error:', error)
    
    return c.json({
      success: false,
      error: 'Failed to explain concept',
      code: 'EXPLAIN_ERROR',
      details: error.message
    }, 500)
  }
})

export { explainRouter }