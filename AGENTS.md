# AGENTS.md - BACAI AI Agent System

## 🤖 Agent Architecture

This document outlines the AI agent system powering BACAI (Mauritanian AI Educational System). Our multi-agent approach provides specialized, contextually-aware educational support for Mauritanian students.

---

## 🧠 Core Agents

### 1. **Problem Solver Agent** 🧮

**Specialization**: Mathematical and scientific problem-solving

- **Models**: Qwen3-8B, Qwen3-235B-A22B
- **Capabilities**:
  - Step-by-step mathematical solutions
  - Scientific concept explanations
  - BEPC/Baccalaureate preparation
  - Equation solving and proofs
- **Languages**: Arabic, French, English
- **Curriculum**: Mauritanian secondary and university levels

### 2. **Language Specialist Agent** 📚

**Specialization**: Arabic, French, and English language tutoring

- **Models**: AraGPT2, Qwen3-8B
- **Capabilities**:
  - Grammar and syntax correction
  - Composition guidance
  - Literary analysis
  - Mauritanian dialect support
- **Focus**: Classical Arabic, Mauritanian literature, French composition

### 3. **Islamic Studies Agent** 🕌

**Specialization**: Islamic education and cultural context

- **Models**: Specialized Arabic models
- **Capabilities**:
  - Quran interpretation
  - Hadith explanation
  - Islamic jurisprudence (Fiqh)
  - Mauritanian Islamic scholarship
- **Cultural Context**: Mahdara system integration

### 4. **Sciences Agent** 🔬

**Specialization**: Physics, Chemistry, Biology

- **Models**: Qwen3-235B-A22B (reasoning focus)
- **Capabilities**:
  - Laboratory exercise explanations
  - Scientific method guidance
  - Concept visualization
  - Mauritanian environmental context

### 5. **Conversational Tutor Agent** 🎓

**Specialization**: Adaptive tutoring and dialogue

- **Models**: Qwen3-8B, custom fine-tuned models
- **Capabilities**:
  - Personalized learning paths
  - Progress tracking
  - Motivational support
  - Cultural adaptation
- **Modes**: Friendly, formal, encouraging

---

## 🔄 Agent Collaboration System

### **Agent Router** 🚦

```python
class AgentRouter:
    def route_request(self, user_input: Request):
        # Language detection
        language = self.detect_language(user_input.text)

        # Subject classification
        subject = self.classify_subject(user_input.text, language)

        # Difficulty assessment
        level = self.assess_level(user_input.text, subject)

        # Agent selection
        if subject == "mathematics":
            return self.problem_solver_agent
        elif subject == "arabic":
            return self.language_specialist_agent
        elif subject == "islamic_studies":
            return self.islamic_studies_agent
        elif subject in ["physics", "chemistry", "biology"]:
            return self.sciences_agent
        else:
            return self.conversational_tutor_agent
```

### **Agent Fusion Engine** ⚡

```python
class AgentFusion:
    def fuse_responses(self, primary_response, context_responses):
        # Combine insights from multiple agents
        fused = {
            "solution": primary_response.solution,
            "cultural_context": context_responses.cultural,
            "language_enhancement": context_responses.linguistic,
            "additional_resources": context_responses.resources
        }
        return fused
```

---

## 🎯 Specialized Capabilities

### **Cultural Context Integration** 🇲🇷

- **Mauritanian Curriculum Alignment**: BEPC, Baccalaureate formats
- **Local Examples**: Sahara desert physics, Islamic mathematics
- **Cultural References**: Traditional Mahdara teaching methods
- **Regional Dialects**: Hassaniya Arabic support

### **Multilingual Intelligence** 🌍

- **Code-Switching**: Seamless language transitions
- **Cross-Lingual Support**: Arabic-French-English connections
- **Cultural Nuances**: Context-appropriate explanations
- **Language Preservation**: Promoting Arabic and French excellence

### **Adaptive Learning** 📈

- **Difficulty Calibration**: Automatic level adjustment
- **Progress Analytics**: Learning pattern recognition
- **Weakness Identification**: Targeted improvement areas
- **Study Path Optimization**: Personalized curriculum

---

## 🏗️ Technical Implementation

### **Model Orchestration**

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Request Router │───▶│  Agent Manager  │───▶│  Model Cluster  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                │
                    ┌───────────┼───────────┐
                    │           │           │
        ┌───────────▼───┐ ┌───▼────┐ ┌──▼────────┐
        │   Math Agent  │ │Arabic   │ │Islamic  │
        │               │ │Agent    │ │Agent    │
        └───────────────┘ └─────────┘ └──────────┘
```

### **Agent Memory System**

- **Session Memory**: Conversation continuity
- **Knowledge Base**: Mauritanian curriculum data
- **Cultural Context**: Local examples and references
- **Learning Analytics**: Student progress tracking

### **Quality Assurance**

- **Cross-Agent Validation**: Multiple perspective verification
- **Cultural Appropriateness**: Mauritanian context checking
- **Educational Soundness**: Pedagogical validation
- **Language Quality**: Grammar and cultural accuracy

---

## 🚀 Agent Evolution Roadmap

### **Phase 1**: Core Implementation ✅

- [x] Problem Solver Agent
- [x] Language Specialist Agent
- [x] Islamic Studies Agent
- [x] Sciences Agent
- [x] Conversational Tutor Agent

### **Phase 2**: Advanced Features 🔄

- [ ] Visual Learning Agent (diagrams, charts)
- [ ] Assessment Agent (quizzes, evaluations)
- [ ] Study Group Agent (collaborative learning)
- [ ] Parent Dashboard Agent (progress reports)

### **Phase 3**: AI Enhancement 🔮

- [ ] Emotional Intelligence Agent
- [ ] Learning Style Adaptation Agent
- [ ] Career Guidance Agent
- [ ] University Preparation Agent

---

## 📊 Performance Metrics

### **Agent Efficiency**

```python
class AgentMetrics:
    def __init__(self):
        self.response_time = {}
        self.accuracy_scores = {}
        self.user_satisfaction = {}
        self.cultural_appropriateness = {}

    def track_performance(self, agent_name, request, response, feedback):
        self.response_time[agent_name] = response.time
        self.accuracy_scores[agent_name] = feedback.correctness
        self.user_satisfaction[agent_name] = feedback.satisfaction
        self.cultural_appropriateness[agent_name] = feedback.cultural_fit
```

### **Quality Indicators**

- **Response Accuracy**: 95% target across all agents
- **Cultural Relevance**: Mauritanian context integration
- **Language Proficiency**: Native-level Arabic/French support
- **Educational Impact**: Learning outcome improvements

---

## 🔧 Agent Configuration

### **Environment Setup**

```yaml
agents:
  problem_solver:
    model: "Qwen/Qwen3-8B-Instruct"
    temperature: 0.3 # Precision-focused
    max_tokens: 2048
    specialties: ["mathematics", "physics", "chemistry"]

  language_specialist:
    model: "aubmindlab/aragpt2-base"
    temperature: 0.5 # Creative language use
    max_tokens: 1024
    specialties: ["arabic", "french", "english"]

  islamic_studies:
    model: "Qwen/Qwen3-8B-Instruct"
    temperature: 0.2 # Conservative knowledge
    max_tokens: 2048
    specialties: ["quran", "hadith", "fiqh"]

  sciences:
    model: "Qwen/Qwen3-235B-A22B-Instruct"
    temperature: 0.4 # Balanced reasoning
    max_tokens: 3072
    specialties: ["physics", "chemistry", "biology"]

  conversational_tutor:
    model: "Qwen/Qwen3-8B-Instruct"
    temperature: 0.7 # Engaging conversation
    max_tokens: 1024
    specialties: ["dialogue", "motivation", "guidance"]
```

---

## 🌟 Unique Features

### **Mauritanian Specialization**

- **Local Context**: Sahara examples, traditional knowledge
- **Cultural Alignment**: Islamic values integration
- **Educational Standards**: BEPC/Baccalaureate preparation
- **Language Preservation**: Arabic excellence promotion

### **Agent Symbiosis**

- **Cross-Referencing**: Agents share insights
- **Context Enrichment**: Multiple perspective integration
- **Quality Assurance**: Peer validation between agents
- **Continuous Learning**: Shared knowledge improvement

### **Accessibility Focus**

- **Free Tier Usage**: No cost barriers for Mauritanian students
- **Mobile Optimization**: Smartphone-compatible interface
- **Offline Capability**: Essential content caching
- **Low Bandwidth**: Efficient for limited internet

---

## 🎓 Educational Impact

### **Student Benefits**

- **Personalized Learning**: Custom-paced education
- **24/7 Availability**: Always-on tutoring support
- **Cultural Relevance**: Localized educational content
- **Language Excellence**: Multilingual mastery development

### **Teacher Support**

- **Teaching Assistant**: AI-powered classroom aid
- **Resource Generation**: Exercise and test creation
- **Student Analytics**: Performance tracking tools
- **Curriculum Alignment**: Mauritanian standards compliance

### **System Advantages**

- **Scalable**: Supports unlimited students
- **Consistent**: Quality education delivery
- **Adaptable**: Curriculum evolution ready
- **Cost-Effective**: Free educational access

---

**🚀 BACAI Agents: Empowering Mauritanian Education Through AI Excellence**
