// Client for calling the model service (Railway)
export async function callModelService(payload) {
    const modelServiceUrl = process.env.MODEL_SERVICE_URL || 'http://localhost:8000';
    try {
        const response = await fetch(`${modelServiceUrl}/api/process`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${process.env.MODEL_SERVICE_TOKEN || 'default-token'}`,
                'User-Agent': 'BACAI-API/1.0'
            },
            body: JSON.stringify(payload),
            signal: AbortSignal.timeout(30000) // 30 second timeout
        });
        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(`Model service error: ${response.status} - ${errorText}`);
        }
        const result = await response.json();
        return result;
    }
    catch (error) {
        console.error('Model service call failed:', error);
        return {
            success: false,
            error: error.message || 'Failed to connect to model service'
        };
    }
}
// Fallback response for when model service is unavailable
export function getFallbackResponse(requestType, language = 'en') {
    const fallbackMessages = {
        solve: {
            en: 'I apologize, but I am currently unable to solve your problem. Please try again later.',
            ar: 'أعتذر، ولكن لا أستطيع حاليًا حل مشكلتك. يرجى المحاولة مرة أخرى لاحقًا.',
            fr: 'Je suis désolé, mais je ne peux pas résoudre votre problème actuellement. Veuillez réessayer plus tard.'
        },
        explain: {
            en: 'I apologize, but I am currently unable to provide explanations. Please try again later.',
            ar: 'أعتذر، ولكن لا أستطيع حاليًا تقديم شروحات. يرجى المحاولة مرة أخرى لاحقًا.',
            fr: 'Je suis désolé, mais je ne peux pas fournir d\'explications actuellement. Veuillez réessayer plus tard.'
        },
        converse: {
            en: 'I apologize, but I am currently unavailable for conversation. Please try again later.',
            ar: 'أعتذر، ولكن لا أتوفر حاليًا للمحادثة. يرجى المحاولة مرة أخرى لاحقًا.',
            fr: 'Je suis désolé, mais je ne suis pas disponible pour la conversation actuellement. Veuillez réessayer plus tard.'
        }
    };
    const message = fallbackMessages[requestType]?.[language] || fallbackMessages[requestType]?.['en'] || 'Service temporarily unavailable.';
    return {
        success: false,
        error: message
    };
}
