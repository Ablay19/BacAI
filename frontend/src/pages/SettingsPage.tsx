import React from "react";
import { useLanguageStore } from "../store/language";
import { useThemeStore } from "../store/theme";
import { UserPreferences, Language } from "../types";

interface SettingsPageProps {}

export default function SettingsPage({}: SettingsPageProps) {
  const { language, setLanguage } = useLanguageStore();
  const { theme, setTheme } = useThemeStore();

  const getTranslation = (key: string) => {
    const translations: Record<string, Record<string, string>> = {
      title: {
        en: "Settings",
        fr: "Paramètres",
        ar: "الإعدادات",
      },
      language: {
        en: "Language",
        fr: "Langue",
        ar: "اللغة",
      },
      theme: {
        en: "Theme",
        fr: "Thème",
        ar: "المظهر",
      },
      notifications: {
        en: "Notifications",
        fr: "Notifications",
        ar: "الإشعارات",
      },
      stepByStep: {
        en: "Step-by-Step Solutions",
        fr: "Solutions Étape par Étape",
        ar: "الحلول خطوة بخطوة",
      },
      culturalContext: {
        en: "Cultural Context",
        fr: "Contexte Culturel",
        ar: "السياق الثقافي",
      },
      save: {
        en: "Save Settings",
        fr: "Enregistrer",
        ar: "حفظ الإعدادات",
      },
    };
    return translations[key]?.[language] || key;
  };

  // Mock user preferences - in real app this would come from API/store
  const [preferences, setPreferences] = React.useState<UserPreferences>({
    language,
    theme,
    notifications: true,
    auto_detect_language: true,
    preferred_model: "qwen-8b",
    step_by_step: true,
    cultural_context: true,
  });

  const handleLanguageChange = (newLanguage: Language) => {
    setLanguage(newLanguage);
    setPreferences((prev) => ({ ...prev, language: newLanguage }));
  };

  const handleThemeChange = (newTheme: "light" | "dark") => {
    setTheme(newTheme);
    setPreferences((prev) => ({ ...prev, theme: newTheme }));
  };

  const handleToggle = (key: keyof UserPreferences) => {
    setPreferences((prev) => ({
      ...prev,
      [key]: typeof prev[key] === "boolean" ? !prev[key] : prev[key],
    }));
  };

  const handleSave = () => {
    // In real app, this would save to backend
    console.log("Saving preferences:", preferences);
    alert(
      getTranslation("save") +
        " - " +
        (language === "ar"
          ? "تم الحفظ بنجاح"
          : language === "fr"
            ? "Enregistré avec succès"
            : "Saved successfully"),
    );
  };

  return (
    <div className="max-w-2xl mx-auto p-6">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">
        {getTranslation("title")}
      </h1>

      <div className="bg-white rounded-lg shadow-md p-6 space-y-6">
        {/* Language Setting */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            {getTranslation("language")}
          </label>
          <select
            value={preferences.language}
            onChange={(e) => handleLanguageChange(e.target.value as Language)}
            className="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="en">English</option>
            <option value="fr">Français</option>
            <option value="ar">العربية</option>
          </select>
        </div>

        {/* Theme Setting */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            {getTranslation("theme")}
          </label>
          <div className="flex space-x-4">
            <button
              onClick={() => handleThemeChange("light")}
              className={`px-4 py-2 rounded-md border ${
                preferences.theme === "light"
                  ? "bg-blue-600 text-white border-blue-600"
                  : "bg-white text-gray-700 border-gray-300"
              }`}
            >
              Light
            </button>
            <button
              onClick={() => handleThemeChange("dark")}
              className={`px-4 py-2 rounded-md border ${
                preferences.theme === "dark"
                  ? "bg-blue-600 text-white border-blue-600"
                  : "bg-white text-gray-700 border-gray-300"
              }`}
            >
              Dark
            </button>
          </div>
        </div>

        {/* Toggle Settings */}
        <div className="space-y-4">
          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-gray-700">
              {getTranslation("notifications")}
            </span>
            <button
              onClick={() => handleToggle("notifications")}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                preferences.notifications ? "bg-blue-600" : "bg-gray-200"
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  preferences.notifications ? "translate-x-6" : "translate-x-1"
                }`}
              />
            </button>
          </div>

          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-gray-700">
              {getTranslation("stepByStep")}
            </span>
            <button
              onClick={() => handleToggle("step_by_step")}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                preferences.step_by_step ? "bg-blue-600" : "bg-gray-200"
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  preferences.step_by_step ? "translate-x-6" : "translate-x-1"
                }`}
              />
            </button>
          </div>

          <div className="flex items-center justify-between">
            <span className="text-sm font-medium text-gray-700">
              {getTranslation("culturalContext")}
            </span>
            <button
              onClick={() => handleToggle("cultural_context")}
              className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                preferences.cultural_context ? "bg-blue-600" : "bg-gray-200"
              }`}
            >
              <span
                className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                  preferences.cultural_context
                    ? "translate-x-6"
                    : "translate-x-1"
                }`}
              />
            </button>
          </div>
        </div>

        {/* Save Button */}
        <div className="pt-4">
          <button
            onClick={handleSave}
            className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 transition-colors duration-200"
          >
            {getTranslation("save")}
          </button>
        </div>
      </div>
    </div>
  );
}
