import { Hono } from 'hono';
const healthRouter = new Hono();
// Basic health check
healthRouter.get('/', (c) => {
    return c.json({
        success: true,
        status: 'healthy',
        timestamp: new Date().toISOString(),
        service: 'bacai-api',
        version: '1.0.0',
        environment: c.env.ENVIRONMENT || 'development'
    });
});
// Detailed health check with dependencies
healthRouter.get('/detailed', async (c) => {
    const health = {
        success: true,
        status: 'healthy',
        timestamp: new Date().toISOString(),
        service: 'bacai-api',
        version: '1.0.0',
        environment: c.env.ENVIRONMENT || 'development',
        dependencies: {
            model_service: 'unknown',
            database: 'unknown',
            cache: 'unknown'
        },
        metrics: {
            uptime: 0, // Will be calculated
            memory_usage: 0,
            cpu_usage: 0
        }
    };
    try {
        // Check model service connectivity
        const modelServiceUrl = c.env.MODEL_SERVICE_URL || 'http://localhost:8000';
        const response = await fetch(`${modelServiceUrl}/health`, {
            method: 'GET',
            signal: AbortSignal.timeout(5000) // 5 second timeout
        });
        health.dependencies.model_service = response.ok ? 'healthy' : 'unhealthy';
    }
    catch (error) {
        health.dependencies.model_service = 'unreachable';
        health.status = 'degraded';
    }
    return c.json(health);
});
// Readiness check (for K8s/deployment)
healthRouter.get('/ready', (c) => {
    // The service is ready when all critical dependencies are healthy
    const isReady = true; // This would check actual dependencies
    return c.json({
        success: true,
        ready: isReady,
        timestamp: new Date().toISOString()
    }, isReady ? 200 : 503);
});
// Liveness check (for K8s/deployment)
healthRouter.get('/live', (c) => {
    // Basic liveness check - if this fails, the service should be restarted
    return c.json({
        success: true,
        alive: true,
        timestamp: new Date().toISOString()
    });
});
export { healthRouter };
