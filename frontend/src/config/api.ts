// API Configuration
// Uses environment variables with fallbacks for development

const getEnv = (key: string, defaultValue: string): string => {
  const value = import.meta.env[key] || defaultValue;
  if (!value) {
    console.warn(
      `Environment variable ${key} is not set, using default: ${defaultValue}`,
    );
  }
  return value;
};

export const API_CONFIG = {
  // Base API URL
  baseUrl: getEnv(
    "VITE_API_URL",
    "https://bacai-api-development.abdoullahelvogani.workers.dev",
  ),

  // API Endpoints
  endpoints: {
    solve: "/api/solve",
    explain: "/api/explain",
    converse: "/api/converse",
    health: "/api/health",
    data: "/api/data",
  },

  // Request Configuration
  timeout: 30000, // 30 seconds
  retries: 3,

  // Feature Flags
  features: {
    darkMode: getEnv("VITE_ENABLE_DARK_MODE", "true") === "true",
    analytics: getEnv("VITE_SHOW_ANALYTICS", "false") === "true",
    rtlSupport: getEnv("VITE_ENABLE_RTL", "true") === "true",
    progressTracking: getEnv("VITE_SHOW_PROGRESS", "true") === "true",
  },

  // Supported Languages
  supportedLanguages: ["ar", "fr", "en"],

  // Supported Subjects
  supportedSubjects: [
    "mathematics",
    "arabic",
    "french",
    "english",
    "sciences",
    "islamic_studies",
  ],

  // Education Levels
  educationLevels: ["secondary_basic", "secondary_lycee", "university"],

  // App Info
  app: {
    name: "BACAI",
    version: getEnv("VITE_APP_VERSION", "1.0.0"),
  },
};

// Helper function to build API URLs
export const buildApiUrl = (endpoint: string): string => {
  const endpointPath =
    API_CONFIG.endpoints[endpoint as keyof typeof API_CONFIG.endpoints];
  return `${API_CONFIG.baseUrl}${endpointPath || endpoint}`;
};

// Validate environment
export const validateEnv = (): { valid: boolean; errors: string[] } => {
  const errors: string[] = [];

  if (!API_CONFIG.baseUrl) {
    errors.push("VITE_API_URL is not configured");
  }

  return {
    valid: errors.length === 0,
    errors,
  };
};

export default API_CONFIG;
