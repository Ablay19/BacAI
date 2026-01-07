import { Routes, Route } from 'react-router-dom'
import { useEffect } from 'react'
import { useLanguageStore } from './store/language'
import { useThemeStore } from './store/theme'
import Layout from './components/Layout'
import HomePage from './pages/HomePage'
import SolverPage from './pages/SolverPage'
import SubjectsPage from './pages/SubjectsPage'
import ProgressPage from './pages/ProgressPage'
import SettingsPage from './pages/SettingsPage'
import './App.css'

function App() {
  const { language, detectLanguage } = useLanguageStore()
  const { theme } = useThemeStore()

  useEffect(() => {
    // Auto-detect browser language
    detectLanguage()
  }, [detectLanguage])

  useEffect(() => {
    // Apply theme and language to body
    document.body.className = language === 'ar' ? 'arabic' : ''
    document.documentElement.dir = language === 'ar' ? 'rtl' : 'ltr'
    document.documentElement.lang = language === 'ar' ? 'ar' : language === 'fr' ? 'fr' : 'en'
  }, [language, theme])

  return (
    <div className={`min-h-screen bg-gray-50 ${language === 'ar' ? 'arabic-text' : ''}`}>
      <Layout>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/solver" element={<SolverPage />} />
          <Route path="/subjects" element={<SubjectsPage />} />
          <Route path="/progress" element={<ProgressPage />} />
          <Route path="/settings" element={<SettingsPage />} />
        </Routes>
      </Layout>
    </div>
  )
}

export default App