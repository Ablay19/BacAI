import { useState } from "react";
import { useLanguageStore } from "../store/language";
import { Solution } from "../types";

interface ProblemSolverProps {
  onSolve: (
    exercise: string,
    subject: string,
    level: string,
    language: string,
    mode: string,
    culturalContext: boolean,
  ) => Promise<Solution>;
  isLoading: boolean;
}

export default function ProblemSolver({
  onSolve,
  isLoading,
}: ProblemSolverProps) {
  const { language } = useLanguageStore();
  const [exercise, setExercise] = useState("");
  const [subject, setSubject] = useState("mathematics");
  const [level, setLevel] = useState("secondary_lycee");
  const [solutionMode, setSolutionMode] = useState("step-by-step");
  const [culturalContext, setCulturalContext] = useState(true);
  const [solution, setSolution] = useState<Solution | null>(null);
  const [error, setError] = useState<string | null>(null);

  const subjects = [
    {
      value: "mathematics",
      label: "Mathematics",
      labelAr: "الرياضيات",
      labelFr: "Mathématiques",
    },
    {
      value: "sciences",
      label: "Sciences",
      labelAr: "العلوم",
      labelFr: "Sciences",
    },
    {
      value: "arabic",
      label: "Arabic",
      labelAr: "اللغة العربية",
      labelFr: "Arabe",
    },
    {
      value: "french",
      label: "French",
      labelAr: "اللغة الفرنسية",
      labelFr: "Français",
    },
    {
      value: "english",
      label: "English",
      labelAr: "اللغة الإنجليزية",
      labelFr: "Anglais",
    },
    {
      value: "islamic_studies",
      label: "Islamic Studies",
      labelAr: "الدراسات الإسلامية",
      labelFr: "Études islamiques",
    },
  ];

  const levels = [
    {
      value: "secondary_basic",
      label: "Secondary Basic",
      labelAr: "التعليم الثانوي الأساسي",
      labelFr: "Secondaire fondamental",
    },
    {
      value: "secondary_lycee",
      label: "Secondary Lycée",
      labelAr: "التعليم الثانوي الثانوي",
      labelFr: "Secondaire lycée",
    },
    {
      value: "university",
      label: "University",
      labelAr: "التعليم الجامعي",
      labelFr: "Université",
    },
  ];

  const solutionModes = [
    {
      value: "step-by-step",
      label: "Step by Step",
      labelAr: "خطوة بخطوة",
      labelFr: "Étape par étape",
    },
    {
      value: "answer-only",
      label: "Answer Only",
      labelAr: "الإجابة فقط",
      labelFr: "Réponse seulement",
    },
  ];

  const handleSolve = async () => {
    if (!exercise.trim()) {
      setError(
        language === "ar"
          ? "الرجاء إدخال مسألة"
          : language === "fr"
            ? "Veuillez entrer un exercice"
            : "Please enter an exercise",
      );
      return;
    }

    setError(null);
    setSolution(null);

    try {
      const result = await onSolve(
        exercise,
        subject,
        level,
        language,
        solutionMode,
        culturalContext,
      );
      setSolution(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : "An error occurred");
    }
  };

  const getLabel = (item: any) => {
    switch (language) {
      case "ar":
        return item.labelAr || item.label;
      case "fr":
        return item.labelFr || item.label;
      default:
        return item.label;
    }
  };

  const getTextDirection = () => (language === "ar" ? "rtl" : "ltr");
  const getTextAlign = () => (language === "ar" ? "text-right" : "text-left");

  return (
    <div className="max-w-4xl mx-auto p-6">
      <div className="bg-white rounded-lg shadow-lg border border-gray-200">
        {/* Header */}
        <div className="p-6 border-b border-gray-200">
          <h1 className="text-2xl font-bold text-gray-900 mb-2">
            {language === "ar"
              ? "حل المسائل"
              : language === "fr"
                ? "Résolveur de Problèmes"
                : "Problem Solver"}
          </h1>
          <p className="text-gray-600">
            {language === "ar"
              ? "أدخل مسألتك واحصل على حل خطوة بخطوة مع الشرح"
              : language === "fr"
                ? "Entrez votre exercice et obtenez une solution étape par étape avec explications"
                : "Enter your problem and get a step-by-step solution with explanations"}
          </p>
        </div>

        {/* Input Form */}
        <div className="p-6 space-y-6">
          {/* Exercise Input */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {language === "ar"
                ? "المسألة:"
                : language === "fr"
                  ? "Exercice:"
                  : "Exercise:"}
            </label>
            <textarea
              value={exercise}
              onChange={(e) => setExercise(e.target.value)}
              className={`w-full h-32 p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500 resize-none ${getTextAlign()}`}
              placeholder={
                language === "ar"
                  ? "أدخل مسألتك هنا..."
                  : language === "fr"
                    ? "Entrez votre exercice ici..."
                    : "Enter your exercise here..."
              }
              dir={getTextDirection()}
            />
          </div>

          {/* Subject Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {language === "ar"
                ? "المادة:"
                : language === "fr"
                  ? "Matière:"
                  : "Subject:"}
            </label>
            <select
              value={subject}
              onChange={(e) => setSubject(e.target.value)}
              className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            >
              {subjects.map((sub) => (
                <option key={sub.value} value={sub.value}>
                  {getLabel(sub)}
                </option>
              ))}
            </select>
          </div>

          {/* Level Selection */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {language === "ar"
                ? "المستوى:"
                : language === "fr"
                  ? "Niveau:"
                  : "Level:"}
            </label>
            <select
              value={level}
              onChange={(e) => setLevel(e.target.value)}
              className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary-500 focus:border-primary-500"
            >
              {levels.map((lvl) => (
                <option key={lvl.value} value={lvl.value}>
                  {getLabel(lvl)}
                </option>
              ))}
            </select>
          </div>

          {/* Solution Mode */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              {language === "ar"
                ? "طريقة الحل:"
                : language === "fr"
                  ? "Mode de solution:"
                  : "Solution Mode:"}
            </label>
            <div className="flex space-x-4">
              {solutionModes.map((mode) => (
                <label key={mode.value} className="flex items-center">
                  <input
                    type="radio"
                    value={mode.value}
                    checked={solutionMode === mode.value}
                    onChange={(e) => setSolutionMode(e.target.value)}
                    className="mr-2"
                  />
                  <span>{getLabel(mode)}</span>
                </label>
              ))}
            </div>
          </div>

          {/* Cultural Context Toggle */}
          <div className="flex items-center">
            <input
              type="checkbox"
              id="cultural-context"
              checked={culturalContext}
              onChange={(e) => setCulturalContext(e.target.checked)}
              className="mr-2"
            />
            <label htmlFor="cultural-context" className="text-sm text-gray-700">
              {language === "ar"
                ? "تضمين السياق الثقافي الموريتاني"
                : language === "fr"
                  ? "Inclure le contexte culturel mauritanien"
                  : "Include Mauritanian cultural context"}
            </label>
          </div>

          {/* Error Display */}
          {error && (
            <div className="bg-red-50 border border-red-200 text-red-700 p-4 rounded-lg">
              {error}
            </div>
          )}

          {/* Solve Button */}
          <div className="flex justify-center">
            <button
              onClick={handleSolve}
              disabled={isLoading || !exercise.trim()}
              className="px-8 py-3 bg-primary-600 text-white font-medium rounded-lg hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {isLoading ? (
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  <span>
                    {language === "ar"
                      ? "جاري الحل..."
                      : language === "fr"
                        ? "Résolution en cours..."
                        : "Solving..."}
                  </span>
                </div>
              ) : (
                <>
                  {language === "ar"
                    ? "حل المسألة"
                    : language === "fr"
                      ? "Résoudre"
                      : "Solve Problem"}
                </>
              )}
            </button>
          </div>
        </div>
      </div>

      {/* Solution Display */}
      {solution && (
        <div className="mt-6">
          {/* Import SolutionDisplay component here */}
          <div className="bg-white rounded-lg shadow-lg border border-gray-200 p-6">
            <h2 className="text-xl font-bold text-gray-900 mb-4">
              {language === "ar"
                ? "الحل"
                : language === "fr"
                  ? "Solution"
                  : "Solution"}
            </h2>

            {/* Steps */}
            <div className="space-y-4 mb-6">
              {solution.steps?.map((step, index) => (
                <div key={index} className="border-l-4 border-primary-500 pl-4">
                  <div className="font-medium text-gray-900 mb-1">
                    {language === "ar"
                      ? `الخطوة ${step.step_number}`
                      : language === "fr"
                        ? `Étape ${step.step_number}`
                        : `Step ${step.step_number}`}
                  </div>
                  <div className="text-gray-600 mb-2">{step.explanation}</div>
                  <div className="bg-gray-50 p-3 rounded font-mono text-sm">
                    {step.output}
                  </div>
                </div>
              ))}
            </div>

            {/* Final Answer */}
            <div className="bg-gradient-to-r from-primary-500 to-primary-600 text-white p-4 rounded-lg">
              <div className="font-bold mb-2">
                {language === "ar"
                  ? "الإجابة النهائية:"
                  : language === "fr"
                    ? "Réponse finale:"
                    : "Final Answer:"}
              </div>
              <div className="text-lg">{solution.final_answer}</div>
              <div className="text-sm opacity-75 mt-2">
                {language === "ar"
                  ? `مستوى الثقة: ${(solution.confidence * 100).toFixed(1)}%`
                  : language === "fr"
                    ? `Confiance: ${(solution.confidence * 100).toFixed(1)}%`
                    : `Confidence: ${(solution.confidence * 100).toFixed(1)}%`}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
