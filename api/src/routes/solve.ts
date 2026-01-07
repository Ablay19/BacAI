import { Hono } from "hono";
import { zValidator } from "@hono/zod-validator";
import { z } from "zod";
import { callModelService } from "../utils/model-client";
import { detectLanguage, extractSubject } from "../utils/language-utils";

const solveSchema = z.object({
  exercise: z.string().min(1, "Exercise text is required"),
  subject: z
    .enum([
      "mathematics",
      "arabic",
      "french",
      "english",
      "sciences",
      "islamic_studies",
    ])
    .optional(),
  level: z
    .enum(["secondary_basic", "secondary_lycee", "university"])
    .optional(),
  language: z.enum(["ar", "fr", "en"]).optional(),
  mode: z.enum(["step-by-step", "answer-only"]).default("step-by-step"),
  cultural_context: z.boolean().default(true),
});

const solveRouter = new Hono();

solveRouter.post("/", zValidator("json", solveSchema), async (c) => {
  try {
    const startTime = Date.now();
    const validatedData = c.req.valid("json");

    // Auto-detect language if not provided
    const detectedLanguage =
      validatedData.language || detectLanguage(validatedData.exercise);

    // Auto-detect subject if not provided
    const detectedSubject =
      validatedData.subject ||
      extractSubject(validatedData.exercise, detectedLanguage);

    // Prepare payload for model service
    const payload = {
      exercise: validatedData.exercise,
      subject: detectedSubject,
      level: validatedData.level || "secondary_lycee",
      language: detectedLanguage,
      mode: validatedData.mode,
      cultural_context: validatedData.cultural_context,
      request_type: "solve" as const,
    };

    // Call model service
    const modelResponse = await callModelService(payload);

    if (!modelResponse.success) {
      return c.json(
        {
          success: false,
          error: modelResponse.error || "Model service unavailable",
          code: "MODEL_ERROR",
        },
        503,
      );
    }

    const processingTime = Date.now() - startTime;

    // Format response
    const response = {
      success: true,
      data: {
        solution: modelResponse.data.solution,
        subject: detectedSubject,
        level: validatedData.level || "secondary_lycee",
        language: detectedLanguage,
        mode: validatedData.mode,
        metadata: {
          processing_time: processingTime,
          model_used: modelResponse.metadata?.model_used || "unknown",
          language_detected: detectedLanguage,
          subject_detected: detectedSubject,
          timestamp: new Date().toISOString(),
        },
      },
    };

    return c.json(response);
  } catch (error) {
    console.error("Solve endpoint error:", error);

    return c.json(
      {
        success: false,
        error: "Failed to solve exercise",
        code: "SOLVE_ERROR",
        details: error.message,
      },
      500,
    );
  }
});

export { solveRouter };
