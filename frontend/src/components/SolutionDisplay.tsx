import { useState, useEffect } from "react";
import { useLanguageStore } from "../store/language";
import { SolutionStep } from "../types";

interface SolutionDisplayProps {
  steps: SolutionStep[];
  finalAnswer: string;
  confidence: number;
  language: string;
  processingTime?: number;
}

export default function SolutionDisplay({
  steps,
  finalAnswer,
  confidence,
  language,
  processingTime,
}: SolutionDisplayProps) {
  const { language: currentLanguage } = useLanguageStore();
  const [showStepDetails, setShowStepDetails] = useState<number | null>(null);

  const getDirection = () => {
    return language === "ar" || currentLanguage === "ar" ? "rtl" : "ltr";
  };

  const getTextAlign = () => {
    return language === "ar" || currentLanguage === "ar"
      ? "text-right"
      : "text-left";
  };

  const getConfidenceColor = (confidence: number) => {
    if (confidence >= 0.9) return "text-green-600";
    if (confidence >= 0.7) return "text-yellow-600";
    return "text-red-600";
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    // Show toast notification (you could use react-hot-toast here)
  };

  const getStepNumber = (step: SolutionStep) => {
    return language === "ar"
      ? `الخطوة ${step.step_number}`
      : language === "fr"
        ? `Étape ${step.step_number}`
        : `Step ${step.step_number}`;
  };

  return (
    <div className={`space-y-6 ${getDirection()}`}>
      {/* Header */}
      <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div className="flex items-center justify-between mb-4">
          <h3
            className={`text-lg font-semibold text-gray-900 ${getTextAlign()}`}
          >
            {language === "ar"
              ? "الحل خطوة بخطوة"
              : language === "fr"
                ? "Solution étape par étape"
                : "Step-by-Step Solution"}
          </h3>

          <div className="flex items-center space-x-4">
            <div className="text-sm">
              <span className="text-gray-500">
                {language === "ar"
                  ? "الثقة:"
                  : language === "fr"
                    ? "Confiance:"
                    : "Confidence:"}
              </span>
              <span className={`font-medium ${getConfidenceColor(confidence)}`}>
                {(confidence * 100).toFixed(1)}%
              </span>
            </div>

            {processingTime && (
              <div className="text-sm text-gray-500">{processingTime}s</div>
            )}
          </div>
        </div>
      </div>

      {/* Solution Steps */}
      <div className="space-y-4">
        {steps.map((step, index) => (
          <div
            key={index}
            className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden"
          >
            <div
              className="p-4 bg-gray-50 border-b border-gray-200 cursor-pointer hover:bg-gray-100 transition-colors"
              onClick={() =>
                setShowStepDetails(showStepDetails === index ? null : index)
              }
            >
              <div className="flex items-center justify-between">
                <div className="flex items-center space-x-3">
                  <div className="w-8 h-8 bg-primary-600 text-white rounded-full flex items-center justify-center text-sm font-medium">
                    {step.step_number}
                  </div>
                  <h4 className={`font-medium text-gray-900 ${getTextAlign()}`}>
                    {getStepNumber(step)}
                  </h4>
                </div>

                <div className="flex items-center space-x-2">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      copyToClipboard(step.output);
                    }}
                    className="p-2 text-gray-400 hover:text-gray-600 transition-colors"
                    title={
                      language === "ar"
                        ? "نسخ"
                        : language === "fr"
                          ? "Copier"
                          : "Copy"
                    }
                  >
                    📋
                  </button>

                  <div className="text-gray-400">
                    {showStepDetails === index ? "▼" : "▶"}
                  </div>
                </div>
              </div>
            </div>

            {/* Step Content */}
            <div
              className={`
              transition-all duration-300 ease-in-out
              ${showStepDetails === index ? "max-h-96 opacity-100" : "max-h-0 opacity-0 overflow-hidden"}
            `}
            >
              <div className="p-6 space-y-4">
                {step.explanation && (
                  <div className="p-4 bg-blue-50 rounded-lg border border-blue-200">
                    <h5
                      className={`text-sm font-medium text-blue-800 mb-2 ${getTextAlign()}`}
                    >
                      {language === "ar"
                        ? "الشرح:"
                        : language === "fr"
                          ? "Explication:"
                          : "Explanation:"}
                    </h5>
                    <p className={`text-blue-700 ${getTextAlign()}`}>
                      {step.explanation}
                    </p>
                  </div>
                )}

                <div className="p-4 bg-gray-50 rounded-lg">
                  <h5
                    className={`text-sm font-medium text-gray-700 mb-2 ${getTextAlign()}`}
                  >
                    {language === "ar"
                      ? "الناتج:"
                      : language === "fr"
                        ? "Résultat:"
                        : "Output:"}
                  </h5>
                  <div
                    className={`font-mono text-sm bg-white p-3 rounded border border-gray-300 ${getTextAlign()}`}
                  >
                    {step.output}
                  </div>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Final Answer */}
      <div className="bg-gradient-to-r from-primary-500 to-primary-600 rounded-lg shadow-lg text-white p-6">
        <div className="flex items-center justify-between">
          <div>
            <h3 className="text-lg font-semibold mb-2">
              {language === "ar"
                ? "الإجابة النهائية"
                : language === "fr"
                  ? "Réponse finale"
                  : "Final Answer"}
            </h3>
            <div className={`text-2xl font-bold ${getTextAlign()}`}>
              {finalAnswer}
            </div>
          </div>

          <button
            onClick={() => copyToClipboard(finalAnswer)}
            className="p-3 bg-white bg-opacity-20 hover:bg-opacity-30 rounded-lg transition-colors"
            title={
              language === "ar"
                ? "نسخ الإجابة"
                : language === "fr"
                  ? "Copier la réponse"
                  : "Copy answer"
            }
          >
            📋
          </button>
        </div>
      </div>

      {/* Actions */}
      <div className="flex flex-wrap gap-3 justify-center">
        <button
          onClick={() => window.history.back()}
          className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
        >
          {language === "ar"
            ? "مسألة جديدة"
            : language === "fr"
              ? "Nouveau problème"
              : "New Problem"}
        </button>

        <button
          onClick={() => copyToClipboard(finalAnswer)}
          className="px-4 py-2 bg-primary-600 text-white rounded-lg hover:bg-primary-700 transition-colors"
        >
          {language === "ar"
            ? "نسخ الحل"
            : language === "fr"
              ? "Copier la solution"
              : "Copy Solution"}
        </button>

        <button
          onClick={() => {
            // Implement share functionality
            if (navigator.share) {
              navigator.share({
                title:
                  language === "ar"
                    ? "حل المسألة"
                    : language === "fr"
                      ? "Solution du problème"
                      : "Problem Solution",
                text: finalAnswer,
              });
            }
          }}
          className="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        >
          {language === "ar"
            ? "مشاركة"
            : language === "fr"
              ? "Partager"
              : "Share"}
        </button>
      </div>
    </div>
  );
}
