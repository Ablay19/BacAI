/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  readonly VITE_APP_NAME: string;
  readonly VITE_APP_VERSION: string;
  readonly VITE_ENABLE_DARK_MODE: string;
  readonly VITE_ENABLE_RTL: string;
  readonly VITE_SHOW_ANALYTICS: string;
  readonly VITE_SHOW_PROGRESS: string;
  readonly VITE_GA_ID: string;
  readonly VITE_ANALYTICS_URL: string;
  readonly VITE_MODEL_NAME: string;
  readonly VITE_MAX_TOKENS: string;
  readonly VITE_TEMPERATURE: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
