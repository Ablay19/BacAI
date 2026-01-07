import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { Theme, UserPreferences } from '../types'

type Theme = 'light' | 'dark'

interface ThemeStore {
  theme: Theme
  isDark: boolean
  
  // Actions
  setTheme: (theme: Theme) => void
  toggleTheme: () => void
  resetTheme: () => void
}

export const useThemeStore = create<ThemeStore>()(
  persist(
    (set, get) => ({
      theme: 'light',
      isDark: false,

      setTheme: (theme: Theme) => {
        set({
          theme,
          isDark: theme === 'dark'
        })
        
        // Apply theme to document
        if (theme === 'dark') {
          document.documentElement.classList.add('dark')
          document.body.classList.add('bg-gray-900', 'text-gray-100')
          document.body.classList.remove('bg-gray-50', 'text-gray-900')
        } else {
          document.documentElement.classList.remove('dark')
          document.body.classList.remove('bg-gray-900', 'text-gray-100')
          document.body.classList.add('bg-gray-50', 'text-gray-900')
        }
      },

      toggleTheme: () => {
        const currentTheme = get().theme
        const newTheme = currentTheme === 'light' ? 'dark' : 'light'
        get().setTheme(newTheme)
      },

      resetTheme: () => {
        get().setTheme('light')
      }
    }),
    {
      name: 'bacai-theme-storage',
      partialize: (state) => ({ 
        theme: state.theme 
      })
    }
  )
)

interface UserPreferencesStore extends UserPreferences {
  // Additional preferences
  autoSaveProgress: boolean
  showHints: boolean
  solutionDetail: 'simple' | 'detailed' | 'comprehensive'
  
  // Actions
  updatePreferences: (preferences: Partial<UserPreferences>) => void
  resetPreferences: () => void
}

const defaultPreferences: UserPreferences = {
  language: 'en',
  theme: 'light',
  notifications: true,
  auto_detect_language: true,
  preferred_model: 'qwen2-8b',
  step_by_step: true,
  cultural_context: true
}

export const useUserPreferencesStore = create<UserPreferencesStore>()(
  persist(
    (set, get) => ({
      ...defaultPreferences,
      autoSaveProgress: true,
      showHints: true,
      solutionDetail: 'detailed',

      updatePreferences: (newPreferences: Partial<UserPreferences>) => {
        const current = get()
        set({
          ...current,
          ...newPreferences
        })
      },

      resetPreferences: () => {
        set({
          ...defaultPreferences,
          autoSaveProgress: true,
          showHints: true,
          solutionDetail: 'detailed'
        })
      }
    }),
    {
      name: 'bacai-preferences-storage'
    }
  )
)