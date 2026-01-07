import { useLanguageStore } from "../store/language";
import { Subject } from "../types";

interface SubjectCardProps {
  subject: Subject;
  onSelect: (subject: Subject) => void;
  index: number;
}

export default function SubjectCard({
  subject,
  onSelect,
  index,
}: SubjectCardProps) {
  const { language } = useLanguageStore();

  const getDisplayName = () => {
    switch (language) {
      case "ar":
        return subject.name_ar;
      case "fr":
        return subject.name_fr;
      default:
        return subject.name;
    }
  };

  const getDescription = () => {
    switch (language) {
      case "ar":
        return "انقر للمزيد من المعلومات حول هذا الموضوع";
      case "fr":
        return "Cliquez pour en savoir plus sur cette matière";
      default:
        return "Click to learn more about this subject";
    }
  };

  const getButtonText = () => {
    switch (language) {
      case "ar":
        return "اختر المادة";
      case "fr":
        return "Choisir la matière";
      default:
        return "Select Subject";
    }
  };

  const animationDelay = `${index * 100}ms`;

  return (
    <div
      className="bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden hover:shadow-xl transition-all duration-300 transform hover:-translate-y-1 cursor-pointer"
      style={{ animationDelay }}
      onClick={() => onSelect(subject)}
    >
      {/* Header with gradient background */}
      <div
        className="p-6 text-white"
        style={{ backgroundColor: subject.color }}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <div className="text-4xl">{subject.icon}</div>
            <div>
              <h3 className="text-xl font-bold">{getDisplayName()}</h3>
              <p className="text-white text-opacity-90 text-sm">
                {getDescription()}
              </p>
            </div>
          </div>
          <div className="w-12 h-12 bg-white bg-opacity-20 rounded-full flex items-center justify-center">
            <span className="text-2xl">→</span>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="p-6">
        {/* Topics Overview */}
        <div className="mb-4">
          <h4 className="font-medium text-gray-900 mb-3">
            {language === "ar"
              ? "المواضيع المشمولة:"
              : language === "fr"
                ? "Sujets couverts:"
                : "Topics Covered:"}
          </h4>
          <div className="flex flex-wrap gap-2">
            {subject.topics.slice(0, 4).map((topic, idx) => (
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
            {subject.topics.length > 4 && (
              <span className="px-3 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-600">
                +{subject.topics.length - 4}{" "}
                {language === "ar"
                  ? "أخري"
                  : language === "fr"
                    ? "autres"
                    : "more"}
              </span>
            )}
          </div>
        </div>

        {/* Features */}
        <div className="grid grid-cols-2 gap-4 mb-6">
          <div className="flex items-center space-x-2">
            <span className="text-green-500">✓</span>
            <span className="text-sm text-gray-700">
              {language === "ar"
                ? "حل خطوة بخطوة"
                : language === "fr"
                  ? "Solution étape par étape"
                  : "Step-by-step solutions"}
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
          <div className="flex items-center space-x-2">
            <span className="text-green-500">✓</span>
            <span className="text-sm text-gray-700">
              {language === "ar"
                ? "منهجية متوافقة"
                : language === "fr"
                  ? "Alignement curriculaire"
                  : "Curriculum aligned"}
            </span>
          </div>
        </div>

        {/* Select Button */}
        <button
          className="w-full py-3 px-4 rounded-lg font-medium text-white transition-colors duration-200 hover:opacity-90"
          style={{ backgroundColor: subject.color }}
        >
          {getButtonText()}
        </button>
      </div>
    </div>
  );
}

interface SubjectsGridProps {
  subjects: Subject[];
  onSubjectSelect: (subject: Subject) => void;
}

export function SubjectsGrid({ subjects, onSubjectSelect }: SubjectsGridProps) {
  const { language } = useLanguageStore();

  const getGridDirection = () => (language === "ar" ? "rtl" : "ltr");

  return (
    <div className="space-y-8">
      {/* Header */}
      <div
        className={`text-center ${language === "ar" ? "text-right" : "text-left"}`}
      >
        <h2 className="text-3xl font-bold text-gray-900 mb-4">
          {language === "ar"
            ? "اختر الموضوع الذي تريد دراسته"
            : language === "fr"
              ? "Choisissez la matière que vous voulez étudier"
              : "Choose the subject you want to study"}
        </h2>
        <p className="text-lg text-gray-600 max-w-3xl mx-auto">
          {language === "ar"
            ? "غطي جميع المواد الرئيسية في المنهج الموريتاني مع دعم متعدد اللغات وسياق ثقافي مناسب"
            : language === "fr"
              ? "Couvrez toutes les matières principales du curriculum mauritanien avec un support multilingue et un contexte culturel approprié"
              : "Cover all main subjects in the Mauritanian curriculum with multilingual support and appropriate cultural context"}
        </p>
      </div>

      {/* Subjects Grid */}
      <div
        className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-7xl mx-auto"
        dir={getGridDirection()}
      >
        {subjects.map((subject, index) => (
          <SubjectCard
            key={subject.id}
            subject={subject}
            onSelect={onSubjectSelect}
            index={index}
          />
        ))}
      </div>

      {/* Additional Info */}
      <div className="bg-blue-50 rounded-xl p-6 mt-8 max-w-4xl mx-auto">
        <div className="flex items-center space-x-3 mb-3">
          <span className="text-2xl">💡</span>
          <h3 className="text-lg font-semibold text-blue-900">
            {language === "ar"
              ? "نصائح للاستخدام الأمثل"
              : language === "fr"
                ? "Conseils pour une meilleure utilisation"
                : "Tips for Best Use"}
          </h3>
        </div>
        <ul
          className={`space-y-2 text-blue-800 ${language === "ar" ? "text-right" : "text-left"}`}
        >
          <li className="flex items-start space-x-2">
            <span className="text-blue-600 mt-1">•</span>
            <span>
              {language === "ar"
                ? "اكتب المسائل باللغة التي تفضلها (العربية، الفرنسية، الإنجليزية)"
                : language === "fr"
                  ? "Écrivez les exercices dans la langue de votre choix (arabe, français, anglais)"
                  : "Write exercises in your preferred language (Arabic, French, English)"}
            </span>
          </li>
          <li className="flex items-start space-x-2">
            <span className="text-blue-600 mt-1">•</span>
            <span>
              {language === "ar"
                ? "اختر المستوى التعليمي المناسب للحصول على حلول مخصصة"
                : language === "fr"
                  ? "Choisissez le bon niveau éducatif pour obtenir des solutions personnalisées"
                  : "Select the appropriate education level for personalized solutions"}
            </span>
          </li>
          <li className="flex items-start space-x-2">
            <span className="text-blue-600 mt-1">•</span>
            <span>
              {language === "ar"
                ? "فّعل شرح السياق الثقافي للحصول على أمثلة مرتبطة بالبيئة الموريتانية"
                : language === "fr"
                  ? "Activez le contexte culturel pour obtenir des exemples liés à l'environnement mauritanien"
                  : "Enable cultural context for examples related to the Mauritanian environment"}
            </span>
          </li>
        </ul>
      </div>
    </div>
  );
}
