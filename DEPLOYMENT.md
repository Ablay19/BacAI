# BACAI Deployment Configuration

## Environment Variables

Required environment variables for deployment:

### Frontend (Vercel)
- `NEXT_PUBLIC_API_URL`: URL of the API service
- `NEXT_PUBLIC_APP_NAME`: Application name
- `NEXT_PUBLIC_VERSION`: Application version

### API (Cloudflare Workers)
- `MODEL_SERVICE_URL`: URL of the model service
- `ENVIRONMENT`: Environment (development/staging/production)
- `MODEL_SERVICE_TOKEN`: Authentication token for model service

### Model Service (Railway)
- `DATABASE_URL`: Database connection string
- `REDIS_URL`: Redis connection string
- `HUGGINGFACE_API_KEY`: Hugging Face API key
- `OPENAI_API_KEY`: OpenAI API key (optional)
- `ANTHROPIC_API_KEY`: Anthropic API key (optional)

## Deployment Commands

### Frontend to Vercel
```bash
cd frontend
npm run build
npx vercel --prod
```

### API to Cloudflare Workers
```bash
cd api
npm run build
npx wrangler deploy
```

### Model Service to Railway
```bash
cd model-service
railway up
```

## Free Tier Configuration

### Vercel
- 100GB bandwidth/month
- Serverless functions
- Automatic SSL
- Custom domains supported

### Cloudflare Workers
- 100,000 requests/day free
- 10ms CPU time per request
- Global edge network
- KV storage available

### Railway
- 500 hours/month free
- 1GB RAM limit
- PostgreSQL database available
- Redis add-on available

## Docker Deployment (Alternative)

```bash
# Build images
docker-compose build

# Run services
docker-compose up -d
```

## CI/CD Pipeline

GitHub Actions workflow for automatic deployment:

```yaml
name: Deploy to Production
on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      - name: Deploy to Vercel
        run: npx vercel --prod --token ${{ secrets.VERCEL_TOKEN }}
      - name: Deploy to Cloudflare
        run: npx wrangler deploy --token ${{ secrets.CLOUDFLARE_TOKEN }}
      - name: Deploy to Railway
        run: railway up --token ${{ secrets.RAILWAY_TOKEN }}
```