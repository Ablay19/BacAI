import { useLanguageStore } from "../store/language";
import { UserProgress } from "../types";

interface ProgressPageProps {}

export default function ProgressPage({}: ProgressPageProps) {
  const { language } = useLanguageStore();

  const getTranslation = (key: string) => {
    const translations: Record<string, Record<string, string>> = {
      title: {
        en: "Your Progress",
        fr: "Votre Progrès",
        ar: "تقدمك",
      },
      totalExercises: {
        en: "Total Exercises",
        fr: "Total des Exercices",
        ar: "إجمالي التمارين",
      },
      completed: {
        en: "Completed",
        fr: "Terminés",
        ar: "مكتمل",
      },
      strongAreas: {
        en: "Strong Areas",
        fr: "Domaines Forts",
        ar: "المجالات القوية",
      },
      weakAreas: {
        en: "Areas to Improve",
        fr: "Domaines à Améliorer",
        ar: "المجالات التي تحتاج تحسين",
      },
    };
    return translations[key]?.[language] || key;
  };

  // Mock progress data - in real app this would come from API
  const mockProgress: UserProgress = {
    total_exercises: 45,
    completed_exercises: 32,
    subjects_progress: {
      mathematics: 85,
      arabic: 72,
      french: 68,
      sciences: 90,
    },
    weak_areas: ["algebra", "grammar"],
    strong_areas: ["geometry", "vocabulary"],
    last_activity: new Date().toISOString(),
  };

  const completionRate = Math.round(
    (mockProgress.completed_exercises / mockProgress.total_exercises) * 100,
  );

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">
        {getTranslation("title")}
      </h1>

      {/* Overview Card */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="text-center">
            <div className="text-3xl font-bold text-blue-600">
              {mockProgress.total_exercises}
            </div>
            <div className="text-gray-600">
              {getTranslation("totalExercises")}
            </div>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-green-600">
              {mockProgress.completed_exercises}
            </div>
            <div className="text-gray-600">{getTranslation("completed")}</div>
          </div>
          <div className="text-center">
            <div className="text-3xl font-bold text-purple-600">
              {completionRate}%
            </div>
            <div className="text-gray-600">Complete</div>
          </div>
        </div>
      </div>

      {/* Subject Progress */}
      <div className="bg-white rounded-lg shadow-md p-6 mb-6">
        <h2 className="text-xl font-semibold mb-4">Subject Progress</h2>
        <div className="space-y-4">
          {Object.entries(mockProgress.subjects_progress).map(
            ([subject, progress]) => (
              <div key={subject}>
                <div className="flex justify-between mb-1">
                  <span className="capitalize">{subject}</span>
                  <span>{progress}%</span>
                </div>
                <div className="w-full bg-gray-200 rounded-full h-2">
                  <div
                    className="bg-blue-600 h-2 rounded-full"
                    style={{ width: `${progress}%` }}
                  />
                </div>
              </div>
            ),
          )}
        </div>
      </div>

      {/* Areas */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold mb-3 text-green-700">
            {getTranslation("strongAreas")}
          </h3>
          <ul className="space-y-2">
            {mockProgress.strong_areas.map((area, index) => (
              <li key={index} className="flex items-center">
                <span className="w-2 h-2 bg-green-500 rounded-full mr-3"></span>
                <span className="capitalize">{area}</span>
              </li>
            ))}
          </ul>
        </div>

        <div className="bg-white rounded-lg shadow-md p-6">
          <h3 className="text-lg font-semibold mb-3 text-orange-700">
            {getTranslation("weakAreas")}
          </h3>
          <ul className="space-y-2">
            {mockProgress.weak_areas.map((area, index) => (
              <li key={index} className="flex items-center">
                <span className="w-2 h-2 bg-orange-500 rounded-full mr-3"></span>
                <span className="capitalize">{area}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>
    </div>
  );
}
