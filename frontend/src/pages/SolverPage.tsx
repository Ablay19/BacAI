import { useState } from "react";
import { SubjectsGrid } from "../components/SubjectCard";
import { useLanguageStore } from "../store/language";
import { Subject } from "../types";
import ProblemSolver from "../components/ProblemSolver";
import SolutionDisplay from "../components/SolutionDisplay";
import { Solution } from "../types";

export default function SolverPage() {
  const { language } = useLanguageStore();
  const [selectedSubject, setSelectedSubject] = useState<Subject | null>(null);
  const [solution] = useState<Solution | null>(null);
  const [isLoading, setIsLoading] = useState(false);

  const subjects: Subject[] = [
    {
      id: "mathematics",
      name: "Mathematics",
      name_ar: "الرياضيات",
      name_fr: "Mathématiques",
      icon: "🧮",
      color: "#3B82F6",
      topics: [
        {
          id: "algebra",
          name: "Algebra",
          name_ar: "الجبر",
          name_fr: "Algèbre",
          description: "Linear equations, quadratic equations, polynomials",
        },
        {
          id: "geometry",
          name: "Geometry",
          name_ar: "الهندسة",
          name_fr: "Géométrie",
          description: "Triangles, circles, area, volume",
        },
        {
          id: "calculus",
          name: "Calculus",
          name_ar: "حساب التفاضل والتكامل",
          name_fr: "Calcul",
          description: "Derivatives, integrals, limits",
        },
        {
          id: "statistics",
          name: "Statistics",
          name_ar: "الإحصاء",
          name_fr: "Statistiques",
          description: "Probability, data analysis, charts",
        },
      ],
    },
    {
      id: "sciences",
      name: "Sciences",
      name_ar: "العلوم",
      name_fr: "Sciences",
      icon: "🔬",
      color: "#10B981",
      topics: [
        {
          id: "physics",
          name: "Physics",
          name_ar: "الفيزياء",
          name_fr: "Physique",
          description: "Mechanics, electricity, magnetism",
        },
        {
          id: "chemistry",
          name: "Chemistry",
          name_ar: "الكيمياء",
          name_fr: "Chimie",
          description: "Atoms, molecules, reactions",
        },
        {
          id: "biology",
          name: "Biology",
          name_ar: "الأحياء",
          name_fr: "Biologie",
          description: "Cells, organisms, ecosystems",
        },
      ],
    },
    {
      id: "arabic",
      name: "Arabic",
      name_ar: "اللغة العربية",
      name_fr: "Arabe",
      icon: "📝",
      color: "#F59E0B",
      topics: [
        {
          id: "grammar",
          name: "Grammar",
          name_ar: "قواعد",
          name_fr: "Grammaire",
          description: "Syntax, morphology, sentence structure",
        },
        {
          id: "literature",
          name: "Literature",
          name_ar: "الأدب",
          name_fr: "Littérature",
          description: "Poetry, prose, literary analysis",
        },
        {
          id: "composition",
          name: "Composition",
          name_ar: "الإنشاء",
          name_fr: "Rédaction",
          description: "Essay writing, creative expression",
        },
      ],
    },
    {
      id: "french",
      name: "French",
      name_ar: "اللغة الفرنسية",
      name_fr: "Français",
      icon: "🥖",
      color: "#EF4444",
      topics: [
        {
          id: "french_grammar",
          name: "Grammar",
          name_ar: "قواعد الفرنسية",
          name_fr: "Grammaire française",
          description: "Verb conjugation, sentence structure",
        },
        {
          id: "french_literature",
          name: "Literature",
          name_ar: "الأدب الفرنسي",
          name_fr: "Littérature française",
          description: "French literary works and analysis",
        },
        {
          id: "french_composition",
          name: "Composition",
          name_ar: "الإنشاء الفرنسي",
          name_fr: "Rédaction française",
          description: "French essay and composition",
        },
      ],
    },
    {
      id: "english",
      name: "English",
      name_ar: "اللغة الإنجليزية",
      name_fr: "Anglais",
      icon: "🇬🇧",
      color: "#8B5CF6",
      topics: [
        {
          id: "english_grammar",
          name: "Grammar",
          name_ar: "قواعد الإنجليزية",
          name_fr: "Grammaire anglaise",
          description: "English grammar and structure",
        },
        {
          id: "vocabulary",
          name: "Vocabulary",
          name_ar: "المفردات",
          name_fr: "Vocabulaire",
          description: "Building English vocabulary",
        },
        {
          id: "comprehension",
          name: "Comprehension",
          name_ar: "الفهم",
          name_fr: "Compréhension",
          description: "Reading and listening comprehension",
        },
      ],
    },
    {
      id: "islamic_studies",
      name: "Islamic Studies",
      name_ar: "الدراسات الإسلامية",
      name_fr: "Études islamiques",
      icon: "🕌",
      color: "#059669",
      topics: [
        {
          id: "quran",
          name: "Quran",
          name_ar: "القرآن",
          name_fr: "Coran",
          description: "Quranic verses, recitation, memorization",
        },
        {
          id: "hadith",
          name: "Hadith",
          name_ar: "الحديث",
          name_fr: "Hadith",
          description: "Prophetic traditions and sayings",
        },
        {
          id: "fiqh",
          name: "Fiqh",
          name_ar: "الفقه",
          name_fr: "Fiqh",
          description: "Islamic jurisprudence and law",
        },
        {
          id: "islamic_history",
          name: "Islamic History",
          name_ar: "التاريخ الإسلامي",
          name_fr: "Histoire islamique",
          description: "Historical development of Islamic civilization",
        },
      ],
    },
  ];

  const handleSolve = async (
    exercise: string,
    subject: string,
    level: string,
    lang: string,
    mode: string,
    culturalContext: boolean,
  ): Promise<Solution> => {
    setIsLoading(true);

    try {
      // Call to API (mock implementation for now)
      const response = await fetch("/api/solve", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          exercise,
          subject,
          level,
          language: lang,
          mode,
          cultural_context: culturalContext,
        }),
      });

      if (!response.ok) {
        throw new Error("Failed to solve problem");
      }

      const result = await response.json();
      setIsLoading(false);
      return result.data.solution;
    } catch (error) {
      setIsLoading(false);

      // Mock solution for development
      return {
        steps: [
          {
            explanation:
              language === "ar"
                ? "تحليل المعادلة"
                : language === "fr"
                  ? "Analyse de l'équation"
                  : "Analyzing the equation",
            output: exercise,
            step_number: 1,
          },
          {
            explanation:
              language === "ar"
                ? "تطبيق طريقة الحل"
                : language === "fr"
                  ? "Application de la méthode"
                  : "Applying solution method",
            output:
              language === "ar"
                ? "الخطوة الأولى"
                : language === "fr"
                  ? "Première étape"
                  : "First step",
            step_number: 2,
          },
        ],
        final_answer:
          language === "ar"
            ? "الإجابة المحتملة"
            : language === "fr"
              ? "Réponse possible"
              : "Possible answer",
        confidence: 0.95,
        language_detected: language as "ar" | "fr" | "en",
        model_used: "qwen-8b",
        processing_time: 1500,
      };
    }
  };

  const handleSubjectSelect = (subject: Subject) => {
    setSelectedSubject(subject);
  };

  const getTextDirection = () => (language === "ar" ? "rtl" : "ltr");
  const getTextAlign = () => (language === "ar" ? "text-right" : "text-left");

  if (!selectedSubject) {
    return (
      <div className={`min-h-screen bg-gray-50 p-8 ${getTextDirection()}`}>
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className={`text-center mb-8 ${getTextAlign()}`}>
            <h1 className="text-4xl font-bold text-gray-900 mb-4">
              {language === "ar"
                ? "اختر الموضوع لحل المسائل"
                : language === "fr"
                  ? "Choisissez une matière pour résoudre les problèmes"
                  : "Choose a Subject to Solve Problems"}
            </h1>
            <p className="text-xl text-gray-600">
              {language === "ar"
                ? "اختر الموضوع الذي تريد العمل معه للحصول على مساعدة متخصصة"
                : language === "fr"
                  ? "Choisissez la matière avec laquelle vous voulez travailler pour obtenir une aide spécialisée"
                  : "Select the subject you want to work with for specialized assistance"}
            </p>
          </div>

          <SubjectsGrid
            subjects={subjects}
            onSubjectSelect={handleSubjectSelect}
          />
        </div>
      </div>
    );
  }

  return (
    <div className={`min-h-screen bg-gray-50 p-8 ${getTextDirection()}`}>
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className={`text-center mb-8 ${getTextAlign()}`}>
          <button
            onClick={() => setSelectedSubject(null)}
            className="mb-4 text-primary-600 hover:text-primary-700 font-medium"
          >
            ←{" "}
            {language === "ar"
              ? "العودة إلى المواد"
              : language === "fr"
                ? "Retour aux matières"
                : "Back to Subjects"}
          </button>

          <div className="flex items-center justify-center space-x-3 mb-4">
            <div className="text-4xl">{selectedSubject.icon}</div>
            <h1 className="text-3xl font-bold text-gray-900">
              {language === "ar"
                ? selectedSubject.name_ar
                : language === "fr"
                  ? selectedSubject.name_fr
                  : selectedSubject.name}
            </h1>
          </div>

          <p className="text-gray-600">
            {language === "ar"
              ? "أدخل مسألتك واحصل على حل خطوة بخطوة مع شروح مفصلة"
              : language === "fr"
                ? "Entrez votre exercice et obtenez une solution étape par étape avec des explications détaillées"
                : "Enter your problem and get step-by-step solutions with detailed explanations"}
          </p>
        </div>

        <ProblemSolver onSolve={handleSolve} isLoading={isLoading} />

        {solution && (
          <div className="mt-8">
            <SolutionDisplay
              steps={solution.steps}
              finalAnswer={solution.final_answer}
              confidence={solution.confidence}
              language={solution.language_detected}
              processingTime={2.3}
            />
          </div>
        )}
      </div>
    </div>
  );
}
