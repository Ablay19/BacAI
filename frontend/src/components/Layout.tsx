import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useLanguageStore } from "../store/language";
import { useThemeStore } from "../store/theme";

interface NavigationItem {
  id: string;
  label: string;
  labelAr: string;
  labelFr: string;
  icon: string;
  path: string;
  badge?: string;
}

const navigationItems: NavigationItem[] = [
  {
    id: "home",
    label: "Home",
    labelAr: "الرئيسية",
    labelFr: "Accueil",
    icon: "🏠",
    path: "/",
  },
  {
    id: "solver",
    label: "Problem Solver",
    labelAr: "حل المسائل",
    labelFr: "Résolveur",
    icon: "🧮",
    path: "/solver",
    badge: "AI",
  },
  {
    id: "subjects",
    label: "Subjects",
    labelAr: "المواد",
    labelFr: "Matières",
    icon: "📚",
    path: "/subjects",
  },
  {
    id: "progress",
    label: "Progress",
    labelAr: "التقدم",
    labelFr: "Progrès",
    icon: "📊",
    path: "/progress",
  },
  {
    id: "settings",
    label: "Settings",
    labelAr: "الإعدادات",
    labelFr: "Paramètres",
    icon: "⚙️",
    path: "/settings",
  },
];

const languageOptions = [
  { code: "ar", name: "العربية", nativeName: "العربية" },
  { code: "fr", name: "French", nativeName: "Français" },
  { code: "en", name: "English", nativeName: "English" },
];

export default function Layout({ children }: { children: React.ReactNode }) {
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const navigate = useNavigate();
  const { language, setLanguage, supportedLanguages } = useLanguageStore();
  const { theme, toggleTheme } = useThemeStore();

  const handleNavigation = (path: string) => {
    navigate(path);
    setIsMobileMenuOpen(false);
  };

  const handleLanguageChange = (newLanguage: "ar" | "fr" | "en") => {
    setLanguage(newLanguage);
  };

  const getLabel = (item: NavigationItem) => {
    switch (language) {
      case "ar":
        return item.labelAr;
      case "fr":
        return item.labelFr;
      default:
        return item.label;
    }
  };

  const currentPath = window.location.pathname;
  const isRTL = language === "ar";

  return (
    <div className={`min-h-screen bg-gray-50 ${isRTL ? "rtl" : "ltr"}`}>
      {/* Mobile Menu Button */}
      <div className="lg:hidden fixed top-4 right-4 z-50">
        <button
          onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          className="p-2 bg-white rounded-lg shadow-lg hover:shadow-xl transition-shadow"
          aria-label="Toggle menu"
        >
          <span className="text-2xl">{isMobileMenuOpen ? "✕" : "☰"}</span>
        </button>
      </div>

      {/* Sidebar */}
      <aside
        className={`
        fixed lg:static inset-y-0 ${isRTL ? "right-0" : "left-0"} z-40
        w-64 bg-white shadow-lg transform transition-transform duration-300 ease-in-out
        ${isMobileMenuOpen ? "translate-x-0" : isRTL ? "translate-x-full" : "-translate-x-full"}
        lg:translate-x-0
      `}
      >
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="p-6 border-b border-gray-200">
            <div className="flex items-center space-x-3">
              <div className="text-3xl">🇲🇷</div>
              <div>
                <h1 className="text-xl font-bold text-gray-900">BACAI</h1>
                <p className="text-sm text-gray-500">
                  {language === "ar"
                    ? "الذكاء الاصطناعي التعليمي الموريتاني"
                    : language === "fr"
                      ? "IA Éducative Mauritanienne"
                      : "Mauritanian Educational AI"}
                </p>
              </div>
            </div>
          </div>

          {/* Navigation */}
          <nav className="flex-1 p-4">
            <ul className="space-y-2">
              {navigationItems.map((item) => (
                <li key={item.id}>
                  <button
                    onClick={() => handleNavigation(item.path)}
                    className={`
                      w-full flex items-center justify-between p-3 rounded-lg transition-colors
                      ${
                        currentPath === item.path
                          ? "bg-primary-100 text-primary-700 font-medium"
                          : "hover:bg-gray-100 text-gray-700"
                      }
                    `}
                  >
                    <div className="flex items-center space-x-3">
                      <span className="text-xl">{item.icon}</span>
                      <span className={isRTL ? "font-arabic" : ""}>
                        {getLabel(item)}
                      </span>
                    </div>
                    {item.badge && (
                      <span className="px-2 py-1 text-xs bg-primary-600 text-white rounded-full">
                        {item.badge}
                      </span>
                    )}
                  </button>
                </li>
              ))}
            </ul>
          </nav>

          {/* Footer */}
          <div className="p-4 border-t border-gray-200">
            {/* Language Selector */}
            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 mb-2">
                {language === "ar"
                  ? "اللغة"
                  : language === "fr"
                    ? "Langue"
                    : "Language"}
              </label>
              <div className="flex flex-wrap gap-2">
                {supportedLanguages.map((lang) => (
                  <button
                    key={lang.code}
                    onClick={() =>
                      handleLanguageChange(lang.code as "ar" | "fr" | "en")
                    }
                    className={`
                      px-3 py-1 text-sm rounded-full transition-colors
                      ${
                        language === lang.code
                          ? "bg-primary-600 text-white"
                          : "bg-gray-200 text-gray-700 hover:bg-gray-300"
                      }
                    `}
                  >
                    {lang.nativeName}
                  </button>
                ))}
              </div>
            </div>

            {/* Theme Toggle */}
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-700">
                {language === "ar"
                  ? "المظهر"
                  : language === "fr"
                    ? "Thème"
                    : "Theme"}
              </span>
              <button
                onClick={toggleTheme}
                className="p-2 rounded-lg bg-gray-100 hover:bg-gray-200 transition-colors"
                aria-label="Toggle theme"
              >
                {theme === "light" ? "🌙" : "☀️"}
              </button>
            </div>
          </div>
        </div>
      </aside>

      {/* Mobile Overlay */}
      {isMobileMenuOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-30 lg:hidden"
          onClick={() => setIsMobileMenuOpen(false)}
        />
      )}

      {/* Main Content */}
      <main className={`flex-1 lg:ml-64 ${isRTL ? "lg:mr-64" : ""}`}>
        {/* Top Bar */}
        <header className="bg-white shadow-sm border-b border-gray-200">
          <div className="px-4 sm:px-6 lg:px-8 py-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-lg font-semibold text-gray-900">
                  {navigationItems.find((item) => item.path === currentPath)
                    ?.label || "BACAI"}
                </h2>
              </div>

              {/* Status Indicators */}
              <div className="flex items-center space-x-4">
                <div className="flex items-center space-x-2 text-sm text-gray-600">
                  <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                  <span>
                    {language === "ar"
                      ? "متصل"
                      : language === "fr"
                        ? "En ligne"
                        : "Online"}
                  </span>
                </div>

                {/* Quick Actions */}
                <div className="flex items-center space-x-2">
                  <button
                    onClick={() => handleNavigation("/solver")}
                    className="px-3 py-1 text-sm bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
                  >
                    {language === "ar"
                      ? "حل مسألة"
                      : language === "fr"
                        ? "Résoudre"
                        : "Solve Problem"}
                  </button>
                </div>
              </div>
            </div>
          </div>
        </header>

        {/* Page Content */}
        <div className="flex-1">{children}</div>
      </main>
    </div>
  );
}
