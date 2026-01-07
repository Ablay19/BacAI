import { Hono } from 'hono'
import { zValidator } from '@hono/zod-validator'
import { z } from 'zod'
import { callModelService } from '../utils/model-client'

const converseSchema = z.object({
  message: z.string().min(1, 'Message is required'),
  conversation_id: z.string().optional(),
  subject: z.enum(['mathematics', 'arabic', 'french', 'english', 'sciences', 'islamic_studies', 'general']).default('general'),
  language: z.enum(['ar', 'fr', 'en']).optional(),
  tutor_mode: z.enum(['friendly', 'formal', 'encouraging']).default('friendly'),
  step_by_step: z.boolean().default(true)
})

const converseRouter = new Hono()

converseRouter.post('/', zValidator('json', converseSchema), async (c) => {
  try {
    const startTime = Date.now()
    const validatedData = c.req.valid('json')
    
    // Prepare payload for model service
    const payload = {
      message: validatedData.message,
      conversation_id: validatedData.conversation_id,
      subject: validatedData.subject,
      language: validatedData.language || 'en',
      tutor_mode: validatedData.tutor_mode,
      step_by_step: validatedData.step_by_step,
      request_type: 'converse'
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
        response: modelResponse.data.response,
        conversation_id: modelResponse.data.conversation_id,
        subject: validatedData.subject,
        language: validatedData.language || 'en',
        metadata: {
          processing_time: processingTime,
          model_used: modelResponse.metadata?.model_used || 'unknown',
          timestamp: new Date().toISOString()
        }
      }
    }

    return c.json(response)

  } catch (error) {
    console.error('Converse endpoint error:', error)
    
    return c.json({
      success: false,
      error: 'Failed to process conversation',
      code: 'CONVERSE_ERROR',
      details: error.message
    }, 500)
  }
})

export { converseRouter }