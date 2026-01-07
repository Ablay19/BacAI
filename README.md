# Mauritanian AI Educational System

AI-powered tutoring system for Mauritanian education supporting Mathematics, Arabic language, Islamic Studies, Sciences, and more.

## Features

🧠 **Multilingual Support** - Arabic, French, English  
📚 **Comprehensive Subjects** - Mathematics, Sciences, Arabic, Islamic Studies  
🎓 **Education Levels** - Secondary (BEPC/Baccalaureate) & University  
🤖 **AI-Powered** - Step-by-step problem solving with explanations  
☁️ **Cloud Native** - Deployed on Vercel, Railway, Cloudflare  

## Architecture

```
Frontend (Vercel) ←→ API Router (Cloudflare) ←→ Model Service (Railway)
```

### Deployment Stack

- **Vercel**: Frontend web application
- **Cloudflare Workers**: API routing and caching
- **Railway**: Heavy AI model inference
- **Hugging Face**: Model hosting and fine-tuning

## Subjects Covered

### Mathematics
- Algebra, Geometry, Calculus
- BEPC & Baccalaureate preparation
- University level mathematics

### Languages
- **Arabic**: Classical Arabic, Mauritanian dialect, literature
- **French**: Grammar, literature, technical writing  
- **English**: Academic and technical English

### Sciences
- Physics, Chemistry, Biology
- Scientific methodology
- Laboratory exercises

### Islamic Studies
- Quran interpretation
- Hadith studies
- Islamic jurisprudence (Fiqh)
- Mauritanian Islamic scholarship

## Quick Start

### 1. Clone and Setup
```bash
git clone <repository>
cd bacai
npm install
```

### 2. Configure Environment
```bash
cp .env.example .env
# Add your API keys and configuration
```

### 3. Local Development
```bash
# Frontend
cd frontend
npm run dev

# API Service  
cd api
npm run dev

# Model Service
cd model-service
npm run dev
```

### 4. Deploy
```bash
# Deploy to free tier services
npm run deploy:vercel      # Frontend
npm run deploy:railway     # Model service  
npm run deploy:cloudflare  # API router
```

## API Endpoints

### Core Services
- `POST /api/solve` - Problem solving with step-by-step solutions
- `POST /api/explain` - Concept explanations
- `POST /api/converse` - Conversational tutoring
- `POST /api/translate` - Multi-language support

### Data Management
- `POST /api/data/upload` - Bulk content ingestion
- `GET /api/subjects` - Available subjects
- `GET /api/curriculum` - Curriculum mapping

## Project Structure

```
bacai/
├── frontend/              # Vercel web application
├── api/                   # Cloudflare Workers API
├── model-service/         # Railway AI model service
├── config/                # Configuration files
├── data/                  # Training and reference data
├── models/                # Fine-tuned models
├── notebooks/             # Development notebooks
├── tests/                 # Testing suite
└── deployment/            # Deployment configurations
```

## Technology Stack

### Frontend (Vercel)
- React 18 with TypeScript
- Tailwind CSS for styling
- Vite for development
- Zustand for state management

### API (Cloudflare Workers)
- Hono framework
- TypeScript
- Edge computing optimization
- Built-in rate limiting

### Model Service (Railway)
- Python 3.11
- FastAPI
- Transformers & PyTorch
- GPU acceleration (when available)

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

For Mauritanian educational content and curriculum questions, please open an issue in the repository.

---

Built with ❤️ for Mauritanian education