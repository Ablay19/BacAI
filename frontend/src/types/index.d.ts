export interface Exercise {
    id: string;
    subject: Subject;
    level: EducationLevel;
    language: Language;
    exercise: string;
    difficulty: 'easy' | 'medium' | 'hard';
    tags?: string[];
}
export interface Solution {
    steps: SolutionStep[];
    final_answer: string;
    confidence: number;
    model_used: string;
    processing_time: number;
    language_detected: Language;
}
export interface SolutionStep {
    explanation: string;
    output: string;
    step_number: number;
}
export interface Subject {
    id: string;
    name: string;
    name_ar: string;
    name_fr: string;
    icon: string;
    color: string;
    topics: Topic[];
}
export interface Topic {
    id: string;
    name: string;
    name_ar: string;
    name_fr: string;
    description: string;
}
export type Language = 'ar' | 'fr' | 'en';
export type EducationLevel = 'secondary_basic' | 'secondary_lycee' | 'university';
export interface UserProgress {
    total_exercises: number;
    completed_exercises: number;
    subjects_progress: Record<string, number>;
    weak_areas: string[];
    strong_areas: string[];
    last_activity: string;
}
export interface ApiResponse<T> {
    data: T;
    success: boolean;
    message?: string;
    metadata?: {
        processing_time: number;
        model_used: string;
        timestamp: string;
    };
}
export interface ErrorResponse {
    success: false;
    error: string;
    code?: string;
    details?: any;
}
export interface Config {
    primary_model: string;
    fallback_model: string;
    max_tokens: number;
    temperature: number;
    supported_languages: Language[];
    subjects: Subject[];
}
export interface UserPreferences {
    language: Language;
    theme: 'light' | 'dark';
    notifications: boolean;
    auto_detect_language: boolean;
    preferred_model: string;
    step_by_step: boolean;
    cultural_context: boolean;
}
//# sourceMappingURL=index.d.ts.map