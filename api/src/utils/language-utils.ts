// Language detection and subject extraction utilities

export function detectLanguage(text: string): 'ar' | 'fr' | 'en' {
  // Simple language detection based on character patterns
  const arabicPattern = /[\u0600-\u06FF]/
  const frenchPattern = /[àâäçéèêëïîôöùûüÿ]/gi
  
  // Count characters for each language
  const arabicChars = (text.match(arabicPattern) || []).length
  const frenchChars = (text.match(frenchPattern) || []).length
  const totalChars = text.replace(/\s/g, '').length
  
  if (arabicChars / totalChars > 0.3) {
    return 'ar'
  } else if (frenchChars / totalChars > 0.1) {
    return 'fr'
  } else {
    return 'en' // Default to English
  }
}

export function extractSubject(text: string, language: string): string {
  const subjectKeywords = {
    ar: {
      mathematics: ['الرياضيات', 'معادلة', 'حل', 'رقم', 'مساحة', 'حجم', 'زاوية'],
      sciences: ['الفيزياء', 'الكيمياء', 'الأحياء', 'تجربة', 'مادة', 'طاقة'],
      arabic: ['قواعد', 'نحو', 'صرف', 'بلاغة', 'أدب', 'شعر'],
      french: ['français', 'grammaire', 'conjugaison', 'vocabulaire'],
      english: ['english', 'grammar', 'vocabulary', 'verb'],
      islamic_studies: ['القرآن', 'الحديث', 'الشريعة', 'الفقه', 'السنة']
    },
    fr: {
      mathematics: ['mathématiques', 'équation', 'résoudre', 'nombre', 'surface', 'volume', 'angle'],
      sciences: ['physique', 'chimie', 'biologie', 'expérience', 'matière', 'énergie'],
      arabic: ['arabe', 'grammaire', 'conjugaison', 'vocabulaire'],
      french: ['français', 'grammaire', 'conjugaison', 'littérature'],
      english: ['anglais', 'grammar', 'vocabulary', 'verb'],
      islamic_studies: ['coran', 'islam', 'charia', 'fiqh', 'sounna']
    },
    en: {
      mathematics: ['mathematics', 'equation', 'solve', 'number', 'area', 'volume', 'angle', 'calculate'],
      sciences: ['physics', 'chemistry', 'biology', 'experiment', 'matter', 'energy'],
      arabic: ['arabic', 'grammar', 'vocabulary', 'language'],
      french: ['french', 'grammar', 'vocabulary', 'conjugation'],
      english: ['english', 'grammar', 'vocabulary', 'verb'],
      islamic_studies: ['quran', 'islam', 'sharia', 'fiqh', 'hadith', 'sunnah']
    }
  }

  const keywords = subjectKeywords[language] || subjectKeywords.en
  const textLower = text.toLowerCase()

  // Count keyword matches for each subject
  const scores: Record<string, number> = {}
  
  for (const [subject, words] of Object.entries(keywords)) {
    let score = 0
    for (const word of words) {
      const regex = new RegExp(`\\b${word}\\b`, 'gi')
      const matches = textLower.match(regex)
      if (matches) {
        score += matches.length
      }
    }
    scores[subject] = score
  }

  // Find subject with highest score
  let bestSubject = 'general'
  let maxScore = 0
  
  for (const [subject, score] of Object.entries(scores)) {
    if (score > maxScore) {
      maxScore = score
      bestSubject = subject
    }
  }

  // If no strong matches, default to general
  if (maxScore === 0) {
    return 'general'
  }

  return bestSubject
}

export function getSubjectDisplayName(subject: string, language: string): string {
  const names = {
    mathematics: {
      en: 'Mathematics',
      ar: 'الرياضيات',
      fr: 'Mathématiques'
    },
    sciences: {
      en: 'Sciences',
      ar: 'العلوم',
      fr: 'Sciences'
    },
    arabic: {
      en: 'Arabic',
      ar: 'اللغة العربية',
      fr: 'Arabe'
    },
    french: {
      en: 'French',
      ar: 'اللغة الفرنسية',
      fr: 'Français'
    },
    english: {
      en: 'English',
      ar: 'اللغة الإنجليزية',
      fr: 'Anglais'
    },
    islamic_studies: {
      en: 'Islamic Studies',
      ar: 'الدراسات الإسلامية',
      fr: 'Études islamiques'
    },
    general: {
      en: 'General',
      ar: 'عام',
      fr: 'Général'
    }
  }

  return names[subject]?.[language] || names[subject]?.['en'] || subject
}

export function validateCurriculumAlignment(
  subject: string, 
  level: string, 
  language: string
): boolean {
  // Mauritanian curriculum validation
  const validCombinations = {
    secondary_basic: {
      subjects: ['mathematics', 'sciences', 'arabic', 'french', 'english', 'islamic_studies'],
      languages: ['ar', 'fr', 'en']
    },
    secondary_lycee: {
      subjects: ['mathematics', 'sciences', 'arabic', 'french', 'english', 'islamic_studies'],
      languages: ['ar', 'fr', 'en']
    },
    university: {
      subjects: ['mathematics', 'sciences', 'arabic', 'french', 'english', 'islamic_studies'],
      languages: ['ar', 'fr', 'en']
    }
  }

  const levelConfig = validCombinations[level as keyof typeof validCombinations]
  if (!levelConfig) {
    return false
  }

  return levelConfig.subjects.includes(subject) && levelConfig.languages.includes(language)
}