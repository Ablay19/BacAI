import { Routes, Route, useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";
import { useLanguageStore } from "./store/language";
import { useThemeStore } from "./store/theme";
import Layout from "./components/Layout";
import HomePage from "./pages/HomePage";
import SolverPage from "./pages/SolverPage";
import SubjectsPage from "./pages/SubjectsPage";
import ProgressPage from "./pages/ProgressPage";
import SettingsPage from "./pages/SettingsPage";
import { validateEnv } from "./config/api";
import "./App.css";

function App() {
  const navigate = useNavigate();
  const { language, detectLanguage } = useLanguageStore();
  const { theme } = useThemeStore();
  const [envError, setEnvError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Validate environment on mount
  useEffect(() => {
    const validation = validateEnv();
    if (!validation.valid) {
      setEnvError(validation.errors.join(", "));
      console.error("Environment validation failed:", validation.errors);
    }
    setIsLoading(false);
  }, []);

  // Auto-detect browser language
  useEffect(() => {
    detectLanguage();
  }, [detectLanguage]);

  // Apply theme and language to body
  useEffect(() => {
    document.body.className = language === "ar" ? "arabic" : "";
    document.documentElement.dir = language === "ar" ? "rtl" : "ltr";
    document.documentElement.lang =
      language === "ar" ? "ar" : language === "fr" ? "fr" : "en";
  }, [language, theme]);

  // Show loading state
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="spinner w-12 h-12 mx-auto mb-4"></div>
          <p className="text-gray-600">Loading BACAI...</p>
        </div>
      </div>
    );
  }

  // Show error state
  if (envError) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-50 p-4">
        <div className="card max-w-md w-full text-center">
          <div className="w-12 h-12 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <svg
              className="w-6 h-6 text-red-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
          </div>
          <h2 className="text-xl font-bold text-gray-900 mb-2">
            Configuration Error
          </h2>
          <p className="text-gray-600 mb-4">{envError}</p>
          <p className="text-sm text-gray-500">
            Please check your environment variables and refresh the page.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div
      className={`min-h-screen bg-gray-50 ${
        language === "ar" ? "arabic-text" : ""
      }`}
    >
      <Layout>
        <Routes>
          <Route path="/" element={<HomePage onNavigate={navigate} />} />
          <Route path="/solver" element={<SolverPage />} />
          <Route
            path="/subjects"
            element={
              <SubjectsPage
                onSubjectSelect={(subject) =>
                  navigate(`/solver?subject=${subject.id}`)
                }
              />
            }
          />
          <Route path="/progress" element={<ProgressPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Routes>
      </Layout>
    </div>
  );
}

export default App;
