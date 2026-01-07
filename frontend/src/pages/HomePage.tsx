import { useLanguageStore } from "../store/language";

interface HomePageProps {
  onNavigate: (path: string) => void;
}

interface Statistic {
  label: string;
  labelAr: string;
  labelFr: string;
  value: string;
  icon: string;
  color: string;
}

interface FeatureCard {
  title: string;
  titleAr: string;
  titleFr: string;
  description: string;
  descriptionAr: string;
  descriptionFr: string;
  icon: string;
  features: string[];
  featuresAr?: string[];
  featuresFr?: string[];
}

export default function HomePage({ onNavigate }: HomePageProps) {
  const { language, isRTL } = useLanguageStore();

  const statistics: Statistic[] = [
    {
      label: "Problems Solved",
      labelAr: "مسائل محلولة",
      labelFr: "Problèmes résolus",
      value: "10,000+",
      icon: "✓",
      color: "text-green-600",
    },
    {
      label: "Subjects Covered",
      labelAr: "مواد مغطاة",
      labelFr: "Matières couvertes",
      value: "6",
      icon: "📚",
      color: "text-blue-600",
    },
    {
      label: "Languages Supported",
      labelAr: "اللغات المدعومة",
      labelFr: "Langues supportées",
      value: "3",
      icon: "🌍",
      color: "text-purple-600",
    },
    {
      label: "Active Users",
      labelAr: "المستخدمون النشطون",
      labelFr: "Utilisateurs actifs",
      value: "500+",
      icon: "👥",
      color: "text-orange-600",
    },
  ];

  const features: FeatureCard[] = [
    {
      title: "AI Problem Solving",
      titleAr: "حل المسائل بالذكاء الاصطناعي",
      titleFr: "Résolution de problèmes par IA",
      description:
        "Step-by-step solutions with explanations for mathematics, sciences, and more.",
      descriptionAr: "حلول خطوة بخطوة مع شروح للرياضيات والعلوم وغيرها.",
      descriptionFr:
        "Solutions étape par étape avec explications pour les mathématiques, les sciences et plus.",
      icon: "🤖",
      features: [
        "Mathematics & Algebra",
        "Physics & Chemistry",
        "Arabic & French",
        "Cultural Context",
      ],
      featuresAr: [
        "الرياضيات والجبر",
        "الفيزياء والكيمياء",
        "العربية والفرنسية",
        "السياق الثقافي",
      ],
      featuresFr: [
        "Mathématiques et Algèbre",
        "Physique et Chimie",
        "Arabe et Français",
        "Contexte culturel",
      ],
    },
    {
      title: "Multilingual Support",
      titleAr: "الدعم متعدد اللغات",
      titleFr: "Support Multilingue",
      description:
        "Learn and solve problems in Arabic, French, or English with culturally relevant examples.",
      descriptionAr:
        "تعلم وحل المسائل بالعربية أو الفرنسية أو الإنجليزية مع أمثلة ثقافياً مناسبة.",
      descriptionFr:
        "Apprenez et résolvez des problèmes en arabe, français ou anglais avec des exemples culturellement pertinents.",
      icon: "🌐",
      features: [
        "Arabic (العربية)",
        "French (Français)",
        "English",
        "Code Switching",
      ],
      featuresAr: ["العربية", "الفرنسية", "الإنجليزية", "تبديل الكود"],
      featuresFr: ["Arabe", "Français", "Anglais", "Code switching"],
    },
    {
      title: "Mauritanian Curriculum",
      titleAr: "المنهج الموريتاني",
      titleFr: "Curriculum Mauritanien",
      description:
        "Aligned with BEPC and Baccalaureate standards with local examples and context.",
      descriptionAr: "متوافق مع معايير BEPC والبكالوريا مع أمثلة وسياق محليين.",
      descriptionFr:
        "Aligné avec les normes BEPC et Baccalauréat avec des exemples et contextes locaux.",
      icon: "🇲🇷",
      features: [
        "BEPC Preparation",
        "Baccalaureate Ready",
        "Local Examples",
        "Islamic Studies",
      ],
      featuresAr: [
        "التحضير لامتحان BEPC",
        "جاهزية البكالوريا",
        "أمثلة محلية",
        "الدراسات الإسلامية",
      ],
      featuresFr: [
        "Préparation BEPC",
        "Prêt Baccalauréat",
        "Exemples locaux",
        "Études islamiques",
      ],
    },
  ];

  const getLabel = (item: Statistic | FeatureCard) => {
    switch (language) {
      case "ar":
        return "labelAr" in item ? item.labelAr : item.titleAr;
      case "fr":
        return "labelFr" in item ? item.labelFr : item.titleFr;
      default:
        return "label" in item ? item.label : item.title;
    }
  };

  const getDescription = (item: FeatureCard) => {
    switch (language) {
      case "ar":
        return item.descriptionAr;
      case "fr":
        return item.descriptionFr;
      default:
        return item.description;
    }
  };

  const getFeatures = (item: FeatureCard) => {
    switch (language) {
      case "ar":
        return item.featuresAr || item.features;
      case "fr":
        return item.featuresFr || item.features;
      default:
        return item.features;
    }
  };

  const getTextDirection = () => (isRTL ? "rtl" : "ltr");
  const getTextAlign = () => (isRTL ? "text-right" : "text-left");

  return (
    <div
      className={`min-h-screen bg-gradient-to-br from-blue-50 via-white to-green-50 ${getTextDirection()}`}
    >
      {/* Hero Section */}
      <section className="relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 sm:py-20">
          <div className={`text-center ${getTextAlign()}`}>
            <div className="flex justify-center mb-6">
              <div className="text-6xl animate-pulse">🇲🇷</div>
            </div>

            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold text-gray-900 mb-6">
              <span className="bg-gradient-to-r from-primary-600 to-green-600 bg-clip-text text-transparent">
                BACAI
              </span>
              <br />
              <span className="text-2xl sm:text-3xl lg:text-4xl font-normal text-gray-700">
                {language === "ar"
                  ? "الذكاء الاصطناعي التعليمي الموريتاني"
                  : language === "fr"
                    ? "IA Éducative Mauritanienne"
                    : "Mauritanian Educational AI"}
              </span>
            </h1>

            <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
              {language === "ar"
                ? "منصة تعليمية متقدمة بالذكاء الاصطناعي للطلاب الموريتانيين، مع دعم متعدد اللغات وسياق ثقافي موريتاني أصيل"
                : language === "fr"
                  ? "Plateforme éducative avancée par IA pour les étudiants mauritaniens, avec support multilingue et contexte culturel mauritanien authentique"
                  : "Advanced AI-powered learning platform for Mauritanian students, with multilingual support and authentic Mauritanian cultural context"}
            </p>

            <div className="flex flex-col sm:flex-row gap-4 justify-center items-center">
              <button
                onClick={() => onNavigate("/solver")}
                className="px-8 py-4 bg-primary-600 text-white font-semibold rounded-lg hover:bg-primary-700 transform hover:scale-105 transition-all duration-200 shadow-lg"
              >
                <span className="mr-2">🧮</span>
                {language === "ar"
                  ? "ابدأ حل المسائل"
                  : language === "fr"
                    ? "Commencer à résoudre"
                    : "Start Solving Problems"}
              </button>

              <button
                onClick={() => onNavigate("/subjects")}
                className="px-8 py-4 bg-white text-primary-600 font-semibold rounded-lg hover:bg-gray-50 transform hover:scale-105 transition-all duration-200 shadow-lg border border-primary-200"
              >
                <span className="mr-2">📚</span>
                {language === "ar"
                  ? "استكشف المواد"
                  : language === "fr"
                    ? "Découvrir les matières"
                    : "Explore Subjects"}
              </button>
            </div>
          </div>
        </div>

        {/* Background decoration */}
        <div className="absolute inset-0 -z-10 opacity-10">
          <div className="absolute top-10 left-10 text-8xl">📐</div>
          <div className="absolute top-20 right-20 text-6xl">🔬</div>
          <div className="absolute bottom-20 left-1/4 text-7xl">📖</div>
          <div className="absolute bottom-10 right-10 text-6xl">🧪</div>
        </div>
      </section>

      {/* Statistics Section */}
      <section className="py-16 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2
              className={`text-3xl font-bold text-gray-900 mb-4 ${getTextAlign()}`}
            >
              {language === "ar"
                ? "إنجازاتنا بالأرقام"
                : language === "fr"
                  ? "Nos réalisations en chiffres"
                  : "Our Impact by Numbers"}
            </h2>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">
              {language === "ar"
                ? "نحن فخورون بالتأثير الذي نحدثه في تعليم الطلاب الموريتانيين"
                : language === "fr"
                  ? "Nous sommes fiers de l'impact que nous avons sur l'éducation des étudiants mauritaniens"
                  : "We're proud of the impact we're making on Mauritanian students' education"}
            </p>
          </div>

          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {statistics.map((stat, index) => (
              <div
                key={index}
                className="text-center p-6 rounded-xl bg-gray-50 hover:bg-gray-100 transition-colors duration-200"
              >
                <div className="text-4xl mb-3">{stat.icon}</div>
                <div
                  className={`text-3xl font-bold text-gray-900 mb-2 ${stat.color}`}
                >
                  {stat.value}
                </div>
                <div
                  className={`text-sm font-medium text-gray-700 ${getTextAlign()}`}
                >
                  {getLabel(stat)}
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-16 bg-gradient-to-br from-primary-50 to-green-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-12">
            <h2
              className={`text-3xl font-bold text-gray-900 mb-4 ${getTextAlign()}`}
            >
              {language === "ar"
                ? "المميزات الرئيسية"
                : language === "fr"
                  ? "Fonctionnalités principales"
                  : "Key Features"}
            </h2>
            <p className="text-lg text-gray-600 max-w-2xl mx-auto">
              {language === "ar"
                ? "أدوات وميزات متطورة مصممة خصيصاً للتعليم الموريتاني"
                : language === "fr"
                  ? "Outils et fonctionnalités avancés spécialement conçus pour l'éducation mauritanienne"
                  : "Advanced tools and features specifically designed for Mauritanian education"}
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <div
                key={index}
                className="bg-white rounded-xl shadow-lg p-8 hover:shadow-xl transform hover:-translate-y-1 transition-all duration-300"
              >
                <div className="text-5xl mb-6 text-center">{feature.icon}</div>

                <h3
                  className={`text-xl font-bold text-gray-900 mb-4 text-center ${getTextAlign()}`}
                >
                  {getLabel(feature)}
                </h3>

                <p
                  className={`text-gray-600 mb-6 text-center ${getTextAlign()}`}
                >
                  {getDescription(feature)}
                </p>

                <ul className="space-y-2">
                  {getFeatures(feature).map((item, idx) => (
                    <li key={idx} className="flex items-start">
                      <span className="text-green-500 mr-2 mt-1">✓</span>
                      <span className="text-sm text-gray-700">{item}</span>
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 bg-gradient-to-r from-primary-600 to-green-600">
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <h2
            className={`text-3xl sm:text-4xl font-bold text-white mb-6 ${getTextAlign()}`}
          >
            {language === "ar"
              ? "هل أنت مستعد للبدء؟"
              : language === "fr"
                ? "Prêt à commencer ?"
                : "Ready to Get Started?"}
          </h2>
          <p
            className={`text-xl text-white text-opacity-90 mb-8 ${getTextAlign()}`}
          >
            {language === "ar"
              ? "انضم إلى آلاف الطلاب الموريتانيين الذين يستخدمون الذكاء الاصطناعي لتحقيق تعليم أفضل"
              : language === "fr"
                ? "Rejoignez des milliers d'étudiants mauritaniens qui utilisent l'IA pour une meilleure éducation"
                : "Join thousands of Mauritanian students using AI for better education"}
          </p>

          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <button
              onClick={() => onNavigate("/solver")}
              className="px-8 py-4 bg-white text-primary-600 font-semibold rounded-lg hover:bg-gray-100 transform hover:scale-105 transition-all duration-200 shadow-lg"
            >
              {language === "ar"
                ? "حل مسألة الآن"
                : language === "fr"
                  ? "Résoudre un problème maintenant"
                  : "Solve a Problem Now"}
            </button>

            <button
              onClick={() => onNavigate("/subjects")}
              className="px-8 py-4 bg-white bg-opacity-20 text-white font-semibold rounded-lg hover:bg-opacity-30 transform hover:scale-105 transition-all duration-200 border border-white border-opacity-30"
            >
              {language === "ar"
                ? "استكشاف المزيد"
                : language === "fr"
                  ? "Explorer plus"
                  : "Learn More"}
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}
