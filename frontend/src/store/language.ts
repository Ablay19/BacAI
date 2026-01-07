import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Language } from "../types";

interface LanguageStore {
  language: Language;
  isRTL: boolean;
  supportedLanguages: { code: Language; name: string; nativeName: string }[];

  // Actions
  setLanguage: (language: Language) => void;
  detectLanguage: () => void;
  toggleDirection: () => void;
}

export const useLanguageStore = create<LanguageStore>()(
  persist(
    (set, get) => ({
      language: "en", // Default to English
      isRTL: false,
      supportedLanguages: [
        { code: "ar", name: "Arabic", nativeName: "العربية" },
        { code: "fr", name: "French", nativeName: "Français" },
        { code: "en", name: "English", nativeName: "English" },
      ],

      setLanguage: (language: Language) => {
        set({
          language,
          isRTL: language === "ar",
        });

        // Update document direction
        document.documentElement.dir = language === "ar" ? "rtl" : "ltr";
        document.documentElement.lang = language;
      },

      detectLanguage: () => {
        const browserLang = navigator.language.split("-")[0] as Language;
        const storedLang = localStorage.getItem("bacai-language-storage");

        let detectedLanguage: Language = "en";

        // Priority: stored preference > browser language > default
        if (storedLang) {
          try {
            const parsed = JSON.parse(storedLang);
            detectedLanguage = parsed.state?.language || "en";
          } catch {
            detectedLanguage = "en";
          }
        } else if (["ar", "fr", "en"].includes(browserLang)) {
          detectedLanguage = browserLang;
        }

        get().setLanguage(detectedLanguage);
      },

      toggleDirection: () => {
        const currentLang = get().language;
        const newLang = currentLang === "ar" ? "en" : "ar";
        get().setLanguage(newLang);
      },
    }),
    {
      name: "bacai-language-storage",
      partialize: (state) => ({
        language: state.language,
      }),
    },
  ),
);
