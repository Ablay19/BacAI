import { useState } from "react";
import { useLanguageStore } from "../store/language";
import { Subject } from "../types";

interface SubjectsPageProps {
  onSubjectSelect: (subject: Subject) => void;
}

export default function SubjectsPage({ onSubjectSelect }: SubjectsPageProps) {
  const { language } = useLanguageStore();
  const [searchQuery, setSearchQuery] = useState("");

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

  const getTextDirection = () => (language === "ar" ? "rtl" : "ltr");
  const getTextAlign = () => (language === "ar" ? "text-right" : "text-left");

  return (
    <div
      className={`min-h-screen bg-gradient-to-br from-blue-50 via-white to-green-50 p-8 ${getTextDirection()}`}
    >
      <div className="max-w-7xl mx-auto">
        {/* Header */}
        <div className={`text-center mb-12 ${getTextAlign()}`}>
          <h1 className="text-4xl sm:text-5xl font-bold text-gray-900 mb-4">
            <span className="bg-gradient-to-r from-primary-600 to-green-600 bg-clip-text text-transparent">
              {language === "ar"
                ? "المواد الدراسية"
                : language === "fr"
                  ? "Matières Éducatives"
                  : "Educational Subjects"}
            </span>
          </h1>
          <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
            {language === "ar"
              ? "استكشف جميع المواد المتاحة في منصتنا الذكية للتعليم الموريتاني"
              : language === "fr"
                ? "Découvrez toutes les matières disponibles sur notre plateforme d'IA éducative mauritanienne"
                : "Discover all available subjects in our Mauritanian educational AI platform"}
          </p>
        </div>

        {/* Search and Filter */}
        <div className="mb-12">
          <div className="max-w-2xl mx-auto">
            <div className="relative">
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder={
                  language === "ar"
                    ? "ابحث عن مادة..."
                    : language === "fr"
                      ? "Rechercher une matière..."
                      : "Search for a subject..."
                }
                className={`w-full px-6 py-4 text-lg rounded-2xl border-2 border-gray-200 focus:border-primary-500 focus:outline-none shadow-lg ${getTextAlign()}`}
                dir={getTextDirection()}
              />
              <div className="absolute right-6 top-1/2 transform -translate-y-1/2 text-2xl text-gray-400">
                🔍
              </div>
            </div>
          </div>
        </div>

        {/* Curriculum Overview */}
        <div className="mb-12 bg-white rounded-2xl shadow-lg p-8">
          <h2
            className={`text-2xl font-bold text-gray-900 mb-6 ${getTextAlign()}`}
          >
            {language === "ar"
              ? "نظرة عامة على المنهج"
              : language === "fr"
                ? "Aperçu du curriculum"
                : "Curriculum Overview"}
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className={`text-center ${getTextAlign()}`}>
              <div className="text-4xl mb-3">🎓</div>
              <h3 className="font-semibold text-gray-900 mb-2">
                {language === "ar"
                  ? "التعليم الثانوي الأساسي"
                  : language === "fr"
                    ? "Secondaire fondamental"
                    : "Secondary Basic"}
              </h3>
              <p className="text-sm text-gray-600">
                {language === "ar"
                  ? "3 سنوات - شهادة BEPC"
                  : language === "fr"
                    ? "3 ans - Certificat BEPC"
                    : "3 years - BEPC Certificate"}
              </p>
            </div>
            <div className={`text-center ${getTextAlign()}`}>
              <div className="text-4xl mb-3">🏫</div>
              <h3 className="font-semibold text-gray-900 mb-2">
                {language === "ar"
                  ? "التعليم الثانوي الثانوي"
                  : language === "fr"
                    ? "Secondaire lycée"
                    : "Secondary Lycée"}
              </h3>
              <p className="text-sm text-gray-600">
                {language === "ar"
                  ? "3 سنوات - شهادة البكالوريا"
                  : language === "fr"
                    ? "3 ans - Baccalauréat"
                    : "3 years - Baccalaureate"}
              </p>
            </div>
            <div className={`text-center ${getTextAlign()}`}>
              <div className="text-4xl mb-3">🎓</div>
              <h3 className="font-semibold text-gray-900 mb-2">
                {language === "ar"
                  ? "التعليم الجامعي"
                  : language === "fr"
                    ? "Enseignement universitaire"
                    : "University Education"}
              </h3>
              <p className="text-sm text-gray-600">
                {language === "ar"
                  ? "بكالوريوس - ماجستير - دكتوراه"
                  : language === "fr"
                    ? "Licence - Master - Doctorat"
                    : "Bachelor - Master - PhD"}
              </p>
            </div>
          </div>
        </div>

        {/* Subject Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {subjects
            .filter((subject) => {
              if (!searchQuery.trim()) return true;
              const query = searchQuery.toLowerCase();
              return (
                subject.name.toLowerCase().includes(query) ||
                subject.name_ar.includes(query) ||
                subject.name_fr.toLowerCase().includes(query) ||
                subject.topics.some(
                  (topic) =>
                    topic.name.toLowerCase().includes(query) ||
                    topic.name_ar.includes(query) ||
                    topic.name_fr.toLowerCase().includes(query)
                )
              );
            })
            .map((subject, index) => (
            <div
              key={subject.id}
              onClick={() => onSubjectSelect(subject)}
              className="bg-white rounded-2xl shadow-lg border border-gray-200 overflow-hidden hover:shadow-2xl transform hover:-translate-y-2 transition-all duration-300 cursor-pointer"
              style={{ animationDelay: `${index * 100}ms` }}
            >
              {/* Header */}
              <div
                className="p-6 text-white"
                style={{ backgroundColor: subject.color }}
              >
                <div className="flex items-center justify-between mb-4">
                  <div className="text-5xl">{subject.icon}</div>
                  <div className="w-12 h-12 bg-white bg-opacity-20 rounded-full flex items-center justify-center">
                    <span className="text-2xl">→</span>
                  </div>
                </div>
                <h3 className="text-2xl font-bold mb-2">
                  {language === "ar"
                    ? subject.name_ar
                    : language === "fr"
                      ? subject.name_fr
                      : subject.name}
                </h3>
                <p className="text-white text-opacity-90 text-sm">
                  {subject.topics.length}{" "}
                  {language === "ar"
                    ? "مواضيع"
                    : language === "fr"
                      ? "sujets"
                      : "topics"}
                </p>
              </div>

              {/* Content */}
              <div className="p-6">
                {/* Topics Preview */}
                <div className="mb-4">
                  <h4 className="font-semibold text-gray-900 mb-3">
                    {language === "ar"
                      ? "المواضيع الرئيسية:"
                      : language === "fr"
                        ? "Sujets principaux:"
                        : "Main Topics:"}
                  </h4>
                  <div className="flex flex-wrap gap-2">
                    {subject.topics.slice(0, 3).map((topic, idx) => (
                      <span
                        key={idx}
                        className="px-3 py-1 text-xs font-medium rounded-full"
                        style={{
                          backgroundColor: `${subject.color}20`,
                          color: subject.color,
                        }}
                      >
                        {topic.name}
                      </span>
                    ))}
                    {subject.topics.length > 3 && (
                      <span className="px-3 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-600">
                        +{subject.topics.length - 3} more
                      </span>
                    )}
                  </div>
                </div>

                {/* Features */}
                <div className="space-y-2 mb-6">
                  <div className="flex items-center space-x-2">
                    <span className="text-green-500">✓</span>
                    <span className="text-sm text-gray-700">
                      {language === "ar"
                        ? "حلول متعددة الخطوات"
                        : language === "fr"
                          ? "Solutions multi-étapes"
                          : "Multi-step solutions"}
                    </span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span className="text-green-500">✓</span>
                    <span className="text-sm text-gray-700">
                      {language === "ar"
                        ? "دعم متعدد اللغات"
                        : language === "fr"
                          ? "Support multilingue"
                          : "Multilingual support"}
                    </span>
                  </div>
                  <div className="flex items-center space-x-2">
                    <span className="text-green-500">✓</span>
                    <span className="text-sm text-gray-700">
                      {language === "ar"
                        ? "سياق ثقافي موريتاني"
                        : language === "fr"
                          ? "Contexte culturel mauritanien"
                          : "Mauritanian cultural context"}
                    </span>
                  </div>
                </div>

                {/* Action Button */}
                <button
                  className="w-full py-3 px-4 rounded-xl font-medium text-white transition-colors duration-200 hover:opacity-90"
                  style={{ backgroundColor: subject.color }}
                >
                  {language === "ar"
                    ? "ابدأ الدراسة"
                    : language === "fr"
                      ? "Commencer l'étude"
                      : "Start Learning"}
                </button>
              </div>
            </div>
          ))}
        </div>

        {/* Learning Tips */}
        <div className="mt-16 bg-gradient-to-r from-blue-50 to-green-50 rounded-2xl p-8">
          <h2
            className={`text-2xl font-bold text-gray-900 mb-6 text-center ${getTextAlign()}`}
          >
            {language === "ar"
              ? "نصائح للتعليم الفعال"
              : language === "fr"
                ? "Conseils pour un apprentissage efficace"
                : "Tips for Effective Learning"}
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className={`text-center ${getTextAlign()}`}>
              <div className="text-4xl mb-3">🎯</div>
              <h3 className="font-semibold text-gray-900 mb-2">
                {language === "ar"
                  ? "حدد أهدافك"
                  : language === "fr"
                    ? "Définissez vos objectifs"
                    : "Set Clear Goals"}
              </h3>
              <p className="text-sm text-gray-700">
                {language === "ar"
                  ? "اختر المواد والمستويات التي تريد التركيز عليها"
                  : language === "fr"
                    ? "Choisissez les matières et niveaux sur lesquels vous concentrer"
                    : "Choose subjects and levels to focus on"}
              </p>
            </div>
            <div className={`text-center ${getTextAlign()}`}>
              <div className="text-4xl mb-3">⏰</div>
              <h3 className="font-semibold text-gray-900 mb-2">
                {language === "ar"
                  ? "مارسة منتظمة"
                  : language === "fr"
                    ? "Pratique régulière"
                    : "Regular Practice"}
              </h3>
              <p className="text-sm text-gray-700">
                {language === "ar"
                  ? "احل المسائل يومياً لتعزيز مهاراتك"
                  : language === "fr"
                    ? "Résolvez des exercices quotidiennement pour renforcer vos compétences"
                    : "Solve problems daily to strengthen skills"}
              </p>
            </div>
            <div className={`text-center ${getTextAlign()}`}>
              <div className="text-4xl mb-3">🤝</div>
              <h3 className="font-semibold text-gray-900 mb-2">
                {language === "ar"
                  ? "اطلب المساعدة"
                  : language === "fr"
                    ? "Demandez de l'aide"
                    : "Ask for Help"}
              </h3>
              <p className="text-sm text-gray-700">
                {language === "ar"
                  ? "لا تتردد في طلب شروح إضافية عند الحاجة"
                  : language === "fr"
                    ? "N'hésitez pas à demander des explications supplémentaires si nécessaire"
                    : "Don't hesitate to ask for additional explanations when needed"}
              </p>
            </div>
            <div className={`text-center ${getTextAlign()}`}>
              <div className="text-4xl mb-3">🌍</div>
              <h3 className="font-semibold text-gray-900 mb-2">
                {language === "ar"
                  ? "استخدم سياقك الثقافي"
                  : language === "fr"
                    ? "Utilisez votre contexte culturel"
                    : "Use Your Cultural Context"}
              </h3>
              <p className="text-sm text-gray-700">
                {language === "ar"
                  ? "استفد من الأمثلة والمراجع المتعلقة بالبيئة الموريتانية"
                  : language === "fr"
                    ? "Tirez parti des exemples et références liés au contexte mauritanien"
                    : "Benefit from examples and references related to Mauritanian context"}
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
