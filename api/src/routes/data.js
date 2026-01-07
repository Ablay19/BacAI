import { Hono } from 'hono';
import { zValidator } from '@hono/zod-validator';
import { z } from 'zod';
const dataRouter = new Hono();
// Available subjects endpoint
dataRouter.get('/subjects', (c) => {
    const subjects = [
        { id: 'mathematics', name: 'Mathematics', name_ar: 'الرياضيات', name_fr: 'Mathématiques' },
        { id: 'arabic', name: 'Arabic', name_ar: 'اللغة العربية', name_fr: 'Arabe' },
        { id: 'french', name: 'French', name_ar: 'اللغة الفرنسية', name_fr: 'Français' },
        { id: 'english', name: 'English', name_ar: 'اللغة الإنجليزية', name_fr: 'Anglais' },
        { id: 'sciences', name: 'Sciences', name_ar: 'العلوم', name_fr: 'Sciences' },
        { id: 'islamic_studies', name: 'Islamic Studies', name_ar: 'الدراسات الإسلامية', name_fr: 'Études islamiques' }
    ];
    return c.json({
        success: true,
        data: { subjects }
    });
});
// Education levels endpoint
dataRouter.get('/levels', (c) => {
    const levels = [
        {
            id: 'secondary_basic',
            name: 'Secondary Basic',
            name_ar: 'التعليم الثانوي الأساسي',
            name_fr: 'Secondaire fondamental',
            certificate: 'BEPC',
            duration: '3 years'
        },
        {
            id: 'secondary_lycee',
            name: 'Secondary Lycée',
            name_ar: 'التعليم الثانوي الثانوي',
            name_fr: 'Secondaire lycée',
            certificate: 'Baccalaureate',
            duration: '3 years'
        },
        {
            id: 'university',
            name: 'University',
            name_ar: 'التعليم الجامعي',
            name_fr: 'Université',
            certificate: 'Bachelor/Master/PhD',
            duration: '3-8 years'
        }
    ];
    return c.json({
        success: true,
        data: { levels }
    });
});
// Curriculum mapping endpoint
dataRouter.get('/curriculum', (c) => {
    const curriculum = {
        mauritania: {
            system: "Baccalaureate system with BEPC and Baccalaureate certificates",
            languages: ["Arabic", "French", "English"],
            structure: {
                "6-15 years": "Compulsory basic education",
                "16-18 years": "Secondary lycée with specializations",
                "18+ years": "Higher education"
            },
            specializations: [
                "Mathematics",
                "Sciences (Physics, Chemistry, Biology)",
                "Literature and Arts",
                "Quran and Arabic Studies",
                "Technical Education"
            ]
        }
    };
    return c.json({
        success: true,
        data: curriculum
    });
});
// Upload data schema (for future bulk upload feature)
const uploadSchema = z.object({
    exercises: z.array(z.object({
        subject: z.string(),
        level: z.string(),
        language: z.string(),
        exercise: z.string(),
        solution: z.string().optional(),
        difficulty: z.enum(['easy', 'medium', 'hard']).optional(),
        tags: z.array(z.string()).optional()
    }))
});
dataRouter.post('/upload', zValidator('json', uploadSchema), async (c) => {
    try {
        const validatedData = c.req.valid('json');
        // TODO: Implement bulk upload logic
        // - Validate against curriculum
        // - Store in database
        // - Update training data
        return c.json({
            success: true,
            message: `Received ${validatedData.exercises.length} exercises for processing`,
            data: {
                processed: 0,
                pending: validatedData.exercises.length
            }
        });
    }
    catch (error) {
        return c.json({
            success: false,
            error: 'Failed to upload data',
            code: 'UPLOAD_ERROR'
        }, 500);
    }
});
export { dataRouter };
