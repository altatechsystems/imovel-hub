# Configuração de Ambientes Dev/Prod com CI/CD

## Visão Geral

Este documento descreve o processo completo para configurar os ambientes de desenvolvimento (dev) e produção (prod) para o projeto Ecossistema Imobiliário, incluindo CI/CD com GitHub Actions.

### Arquitetura dos Ambientes

```
┌─────────────────────────────────────────────────────────────────┐
│                         AMBIENTE DEV                             │
├─────────────────────────────────────────────────────────────────┤
│ Backend (Go)         → Cloud Run (dev)      → api-dev.example.com│
│ Frontend Public      → Vercel (dev)         → dev.example.com    │
│ Frontend Admin       → Vercel (dev)         → admin-dev.ex.com   │
│ Database             → Firestore (dev)      → ecosistema-imob-dev│
│ Storage              → GCS (dev)            → bucket-dev          │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                         AMBIENTE PROD                            │
├─────────────────────────────────────────────────────────────────┤
│ Backend (Go)         → Cloud Run (prod)     → api.example.com    │
│ Frontend Public      → Vercel (prod)        → www.example.com    │
│ Frontend Admin       → Vercel (prod)        → app.example.com    │
│ Database             → Firestore (prod)     → ecosistema-imob-prod│
│ Storage              → GCS (prod)           → bucket-prod         │
└─────────────────────────────────────────────────────────────────┘
```

---

## Fase 1: Preparação do Google Cloud Platform (GCP)

### 1.1 Criar Projeto de Produção no GCP

```bash
# Criar projeto de produção
gcloud projects create ecosistema-imob-prod --name="Ecosistema Imob - Produção"

# Definir projeto como padrão
gcloud config set project ecosistema-imob-prod

# Habilitar APIs necessárias
gcloud services enable \
  run.googleapis.com \
  firestore.googleapis.com \
  storage.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  iam.googleapis.com
```

### 1.2 Configurar Firestore (Dev já existe, criar Prod)

```bash
# Para PROD
gcloud config set project ecosistema-imob-prod
gcloud firestore databases create --location=southamerica-east1 --type=firestore-native

# Para DEV (verificar se já existe)
gcloud config set project ecosistema-imob-dev
gcloud firestore databases describe --database=\(default\)
```

### 1.3 Criar Buckets do Cloud Storage

```bash
# Bucket DEV
gcloud config set project ecosistema-imob-dev
gsutil mb -c STANDARD -l southamerica-east1 gs://ecosistema-imob-dev-storage
gsutil iam ch allUsers:objectViewer gs://ecosistema-imob-dev-storage

# Bucket PROD
gcloud config set project ecosistema-imob-prod
gsutil mb -c STANDARD -l southamerica-east1 gs://ecosistema-imob-prod-storage
gsutil iam ch allUsers:objectViewer gs://ecosistema-imob-prod-storage
```

### 1.4 Criar Service Accounts para CI/CD

```bash
# Service Account para DEV
gcloud config set project ecosistema-imob-dev
gcloud iam service-accounts create github-actions-dev \
  --display-name="GitHub Actions - Development"

# Service Account para PROD
gcloud config set project ecosistema-imob-prod
gcloud iam service-accounts create github-actions-prod \
  --display-name="GitHub Actions - Production"

# Atribuir permissões DEV
gcloud projects add-iam-policy-binding ecosistema-imob-dev \
  --member="serviceAccount:github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding ecosistema-imob-dev \
  --member="serviceAccount:github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding ecosistema-imob-dev \
  --member="serviceAccount:github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com" \
  --role="roles/storage.admin"

# Atribuir permissões PROD
gcloud projects add-iam-policy-binding ecosistema-imob-prod \
  --member="serviceAccount:github-actions-prod@ecosistema-imob-prod.iam.gserviceaccount.com" \
  --role="roles/run.admin"

gcloud projects add-iam-policy-binding ecosistema-imob-prod \
  --member="serviceAccount:github-actions-prod@ecosistema-imob-prod.iam.gserviceaccount.com" \
  --role="roles/iam.serviceAccountUser"

gcloud projects add-iam-policy-binding ecosistema-imob-prod \
  --member="serviceAccount:github-actions-prod@ecosistema-imob-prod.iam.gserviceaccount.com" \
  --role="roles/storage.admin"
```

### 1.5 Gerar Chaves das Service Accounts

```bash
# Chave DEV
gcloud iam service-accounts keys create github-actions-dev-key.json \
  --iam-account=github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com

# Chave PROD
gcloud config set project ecosistema-imob-prod
gcloud iam service-accounts keys create github-actions-prod-key.json \
  --iam-account=github-actions-prod@ecosistema-imob-prod.iam.gserviceaccount.com
```

**IMPORTANTE**: Guarde esses arquivos JSON em local seguro. Eles serão usados nos GitHub Secrets.

---

## Fase 2: Configuração do Firebase Authentication

### 2.1 Configurar Firebase para DEV (já existe)

Verificar se o projeto `ecosistema-imob-dev` já está configurado no Firebase Console.

### 2.2 Criar Projeto Firebase para PROD

1. Acesse [Firebase Console](https://console.firebase.google.com)
2. Clique em "Adicionar projeto"
3. Selecione o projeto GCP existente: `ecosistema-imob-prod`
4. Ative o Google Analytics (opcional)
5. Ative Authentication:
   - Vá em Authentication > Sign-in method
   - Ative "Email/Password"
6. Configure domínios autorizados:
   - Adicione: `api.example.com`, `www.example.com`, `app.example.com`

### 2.3 Obter Credenciais do Firebase

```bash
# DEV - Baixar service account (se ainda não tiver)
# Acesse: https://console.firebase.google.com
# Projeto: ecosistema-imob-dev
# Configurações > Contas de Serviço > Gerar nova chave privada
# Salve como: firebase-adminsdk-dev.json

# PROD - Baixar service account
# Projeto: ecosistema-imob-prod
# Mesmo processo acima
# Salve como: firebase-adminsdk-prod.json
```

---

## Fase 3: Configuração do GitHub Repository

### 3.1 Criar GitHub Secrets

Acesse: `https://github.com/seu-usuario/ecosistema-imob/settings/secrets/actions`

**Secrets para DEV:**
```
GCP_PROJECT_ID_DEV=ecosistema-imob-dev
GCP_SA_KEY_DEV=<conteúdo do arquivo github-actions-dev-key.json>
FIREBASE_PROJECT_ID_DEV=ecosistema-imob-dev
FIREBASE_ADMIN_SDK_DEV=<conteúdo do arquivo firebase-adminsdk-dev.json>
GCS_BUCKET_NAME_DEV=ecosistema-imob-dev-storage
```

**Secrets para PROD:**
```
GCP_PROJECT_ID_PROD=ecosistema-imob-prod
GCP_SA_KEY_PROD=<conteúdo do arquivo github-actions-prod-key.json>
FIREBASE_PROJECT_ID_PROD=ecosistema-imob-prod
FIREBASE_ADMIN_SDK_PROD=<conteúdo do arquivo firebase-adminsdk-prod.json>
GCS_BUCKET_NAME_PROD=ecosistema-imob-prod-storage
```

**Secrets para Vercel (obter em vercel.com/account/tokens):**
```
VERCEL_TOKEN=<token da sua conta Vercel>
VERCEL_ORG_ID=<ID da organização Vercel>
VERCEL_PROJECT_ID_PUBLIC=<ID do projeto frontend-public>
VERCEL_PROJECT_ID_ADMIN=<ID do projeto frontend-admin>
```

### 3.2 Criar Estrutura de Branches

```bash
# Criar branch de desenvolvimento (se não existir)
git checkout -b develop

# Criar branch de staging (opcional)
git checkout -b staging

# Branch main é produção
git checkout main
```

**Estratégia de Branches:**
- `develop` → Deploy automático para DEV
- `staging` → Deploy manual para ambiente de homologação (opcional)
- `main` → Deploy automático para PROD (com aprovação manual)

---

## Fase 4: Configurar CI/CD com GitHub Actions

### 4.1 Criar Workflow para Backend

Criar arquivo: `.github/workflows/backend-deploy.yml`

```yaml
name: Backend Deploy

on:
  push:
    branches:
      - develop  # Deploy automático para DEV
      - main     # Deploy para PROD (requer aprovação)
    paths:
      - 'backend/**'
      - '.github/workflows/backend-deploy.yml'

  workflow_dispatch:  # Permite deploy manual

env:
  GO_VERSION: '1.21'
  REGION: 'southamerica-east1'

jobs:
  # Job para ambiente DEV
  deploy-dev:
    name: Deploy Backend to DEV
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Run tests
        working-directory: ./backend
        run: |
          go test -v ./...

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY_DEV }}

      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Create Firebase config
        working-directory: ./backend
        run: |
          mkdir -p config
          echo '${{ secrets.FIREBASE_ADMIN_SDK_DEV }}' > config/firebase-adminsdk.json

      - name: Build and deploy to Cloud Run
        working-directory: ./backend
        run: |
          gcloud builds submit \
            --tag gcr.io/${{ secrets.GCP_PROJECT_ID_DEV }}/backend-api \
            --project=${{ secrets.GCP_PROJECT_ID_DEV }}

          gcloud run deploy backend-api \
            --image gcr.io/${{ secrets.GCP_PROJECT_ID_DEV }}/backend-api \
            --platform managed \
            --region ${{ env.REGION }} \
            --allow-unauthenticated \
            --set-env-vars "ENVIRONMENT=development,FIREBASE_PROJECT_ID=${{ secrets.FIREBASE_PROJECT_ID_DEV }},GCS_BUCKET_NAME=${{ secrets.GCS_BUCKET_NAME_DEV }}" \
            --project=${{ secrets.GCP_PROJECT_ID_DEV }}

      - name: Get service URL
        run: |
          URL=$(gcloud run services describe backend-api \
            --platform managed \
            --region ${{ env.REGION }} \
            --format 'value(status.url)' \
            --project=${{ secrets.GCP_PROJECT_ID_DEV }})
          echo "Backend DEV deployed to: $URL"

  # Job para ambiente PROD (requer aprovação manual)
  deploy-prod:
    name: Deploy Backend to PROD
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://api.example.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Run tests
        working-directory: ./backend
        run: |
          go test -v ./...

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY_PROD }}

      - name: Set up Cloud SDK
        uses: google-github-actions/setup-gcloud@v2

      - name: Create Firebase config
        working-directory: ./backend
        run: |
          mkdir -p config
          echo '${{ secrets.FIREBASE_ADMIN_SDK_PROD }}' > config/firebase-adminsdk.json

      - name: Build and deploy to Cloud Run
        working-directory: ./backend
        run: |
          gcloud builds submit \
            --tag gcr.io/${{ secrets.GCP_PROJECT_ID_PROD }}/backend-api \
            --project=${{ secrets.GCP_PROJECT_ID_PROD }}

          gcloud run deploy backend-api \
            --image gcr.io/${{ secrets.GCP_PROJECT_ID_PROD }}/backend-api \
            --platform managed \
            --region ${{ env.REGION }} \
            --allow-unauthenticated \
            --set-env-vars "ENVIRONMENT=production,FIREBASE_PROJECT_ID=${{ secrets.FIREBASE_PROJECT_ID_PROD }},GCS_BUCKET_NAME=${{ secrets.GCS_BUCKET_NAME_PROD }}" \
            --project=${{ secrets.GCP_PROJECT_ID_PROD }}

      - name: Get service URL
        run: |
          URL=$(gcloud run services describe backend-api \
            --platform managed \
            --region ${{ env.REGION }} \
            --format 'value(status.url)' \
            --project=${{ secrets.GCP_PROJECT_ID_PROD }})
          echo "Backend PROD deployed to: $URL"
```

### 4.2 Criar Dockerfile para Backend

Criar arquivo: `backend/Dockerfile`

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Copy config directory (firebase credentials will be added during build)
COPY --from=builder /app/config ./config

# Expose port
EXPOSE 8080

# Run
CMD ["./server"]
```

### 4.3 Criar Workflow para Frontend Public

Criar arquivo: `.github/workflows/frontend-public-deploy.yml`

```yaml
name: Frontend Public Deploy

on:
  push:
    branches:
      - develop
      - main
    paths:
      - 'frontend-public/**'
      - '.github/workflows/frontend-public-deploy.yml'

  workflow_dispatch:

jobs:
  deploy-dev:
    name: Deploy Frontend Public to DEV
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: frontend-public/package-lock.json

      - name: Install Vercel CLI
        run: npm install --global vercel@latest

      - name: Pull Vercel Environment Information
        working-directory: ./frontend-public
        run: vercel pull --yes --environment=preview --token=${{ secrets.VERCEL_TOKEN }}

      - name: Build Project Artifacts
        working-directory: ./frontend-public
        run: vercel build --token=${{ secrets.VERCEL_TOKEN }}
        env:
          NEXT_PUBLIC_API_URL: https://backend-api-dev-xxxx.run.app
          NEXT_PUBLIC_ENVIRONMENT: development

      - name: Deploy to Vercel
        working-directory: ./frontend-public
        run: vercel deploy --prebuilt --token=${{ secrets.VERCEL_TOKEN }}

  deploy-prod:
    name: Deploy Frontend Public to PROD
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://www.example.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: frontend-public/package-lock.json

      - name: Install Vercel CLI
        run: npm install --global vercel@latest

      - name: Pull Vercel Environment Information
        working-directory: ./frontend-public
        run: vercel pull --yes --environment=production --token=${{ secrets.VERCEL_TOKEN }}

      - name: Build Project Artifacts
        working-directory: ./frontend-public
        run: vercel build --prod --token=${{ secrets.VERCEL_TOKEN }}
        env:
          NEXT_PUBLIC_API_URL: https://api.example.com
          NEXT_PUBLIC_ENVIRONMENT: production

      - name: Deploy to Vercel
        working-directory: ./frontend-public
        run: vercel deploy --prebuilt --prod --token=${{ secrets.VERCEL_TOKEN }}
```

### 4.4 Criar Workflow para Frontend Admin

Criar arquivo: `.github/workflows/frontend-admin-deploy.yml`

```yaml
name: Frontend Admin Deploy

on:
  push:
    branches:
      - develop
      - main
    paths:
      - 'frontend-admin/**'
      - '.github/workflows/frontend-admin-deploy.yml'

  workflow_dispatch:

jobs:
  deploy-dev:
    name: Deploy Frontend Admin to DEV
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/develop'

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: frontend-admin/package-lock.json

      - name: Install Vercel CLI
        run: npm install --global vercel@latest

      - name: Pull Vercel Environment Information
        working-directory: ./frontend-admin
        run: vercel pull --yes --environment=preview --token=${{ secrets.VERCEL_TOKEN }}

      - name: Build Project Artifacts
        working-directory: ./frontend-admin
        run: vercel build --token=${{ secrets.VERCEL_TOKEN }}
        env:
          NEXT_PUBLIC_API_URL: https://backend-api-dev-xxxx.run.app
          NEXT_PUBLIC_FIREBASE_API_KEY: ${{ secrets.FIREBASE_API_KEY_DEV }}
          NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN: ecosistema-imob-dev.firebaseapp.com
          NEXT_PUBLIC_FIREBASE_PROJECT_ID: ecosistema-imob-dev
          NEXT_PUBLIC_ENVIRONMENT: development

      - name: Deploy to Vercel
        working-directory: ./frontend-admin
        run: vercel deploy --prebuilt --token=${{ secrets.VERCEL_TOKEN }}

  deploy-prod:
    name: Deploy Frontend Admin to PROD
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://app.example.com

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: frontend-admin/package-lock.json

      - name: Install Vercel CLI
        run: npm install --global vercel@latest

      - name: Pull Vercel Environment Information
        working-directory: ./frontend-admin
        run: vercel pull --yes --environment=production --token=${{ secrets.VERCEL_TOKEN }}

      - name: Build Project Artifacts
        working-directory: ./frontend-admin
        run: vercel build --prod --token=${{ secrets.VERCEL_TOKEN }}
        env:
          NEXT_PUBLIC_API_URL: https://api.example.com
          NEXT_PUBLIC_FIREBASE_API_KEY: ${{ secrets.FIREBASE_API_KEY_PROD }}
          NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN: ecosistema-imob-prod.firebaseapp.com
          NEXT_PUBLIC_FIREBASE_PROJECT_ID: ecosistema-imob-prod
          NEXT_PUBLIC_ENVIRONMENT: production

      - name: Deploy to Vercel
        working-directory: ./frontend-admin
        run: vercel deploy --prebuilt --prod --token=${{ secrets.VERCEL_TOKEN }}
```

---

## Fase 5: Configurar Ambientes no Vercel

### 5.1 Instalar Vercel CLI

```bash
npm install -g vercel
vercel login
```

### 5.2 Criar Projetos no Vercel

```bash
# Frontend Public
cd frontend-public
vercel link  # Seguir instruções para criar novo projeto

# Frontend Admin
cd ../frontend-admin
vercel link  # Seguir instruções para criar novo projeto
```

### 5.3 Configurar Variáveis de Ambiente no Vercel

**Frontend Public - DEV (Preview):**
```
NEXT_PUBLIC_API_URL=https://backend-api-dev-xxxxx.run.app
NEXT_PUBLIC_ENVIRONMENT=development
```

**Frontend Public - PROD (Production):**
```
NEXT_PUBLIC_API_URL=https://api.example.com
NEXT_PUBLIC_ENVIRONMENT=production
```

**Frontend Admin - DEV (Preview):**
```
NEXT_PUBLIC_API_URL=https://backend-api-dev-xxxxx.run.app
NEXT_PUBLIC_FIREBASE_API_KEY=<chave do Firebase DEV>
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=ecosistema-imob-dev.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=ecosistema-imob-dev
NEXT_PUBLIC_ENVIRONMENT=development
```

**Frontend Admin - PROD (Production):**
```
NEXT_PUBLIC_API_URL=https://api.example.com
NEXT_PUBLIC_FIREBASE_API_KEY=<chave do Firebase PROD>
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=ecosistema-imob-prod.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=ecosistema-imob-prod
NEXT_PUBLIC_ENVIRONMENT=production
```

---

## Fase 6: Configurar Domínios Personalizados

### 6.1 Backend (Cloud Run)

```bash
# DEV
gcloud beta run domain-mappings create \
  --service backend-api \
  --domain api-dev.example.com \
  --region southamerica-east1 \
  --project ecosistema-imob-dev

# PROD
gcloud beta run domain-mappings create \
  --service backend-api \
  --domain api.example.com \
  --region southamerica-east1 \
  --project ecosistema-imob-prod
```

### 6.2 Frontend (Vercel)

1. Acesse o projeto no Vercel Dashboard
2. Vá em Settings > Domains
3. Adicione os domínios:
   - **Frontend Public DEV**: `dev.example.com`
   - **Frontend Public PROD**: `www.example.com`
   - **Frontend Admin DEV**: `admin-dev.example.com`
   - **Frontend Admin PROD**: `app.example.com`

### 6.3 Configurar DNS

No seu provedor de DNS (ex: Cloudflare, Route53), adicione os registros CNAME:

```
# Backend
api-dev.example.com     CNAME  ghs.googlehosted.com
api.example.com         CNAME  ghs.googlehosted.com

# Frontend Public
dev.example.com         CNAME  cname.vercel-dns.com
www.example.com         CNAME  cname.vercel-dns.com

# Frontend Admin
admin-dev.example.com   CNAME  cname.vercel-dns.com
app.example.com         CNAME  cname.vercel-dns.com
```

---

## Fase 7: Configurar Proteção de Branch no GitHub

### 7.1 Proteger Branch Main (Produção)

1. Acesse: `https://github.com/seu-usuario/ecosistema-imob/settings/branches`
2. Clique em "Add rule"
3. Branch name pattern: `main`
4. Configurações recomendadas:
   - ✅ Require a pull request before merging
   - ✅ Require approvals (mínimo 1)
   - ✅ Require status checks to pass before merging
   - ✅ Require conversation resolution before merging
   - ✅ Do not allow bypassing the above settings

### 7.2 Configurar Ambientes com Aprovação Manual

1. Acesse: `https://github.com/seu-usuario/ecosistema-imob/settings/environments`
2. Clique em "New environment"
3. Nome: `production`
4. Configurações:
   - ✅ Required reviewers (adicionar usuários aprovadores)
   - ✅ Wait timer: 0 minutes (ou definir tempo de espera)
   - Deployment branches: Only protected branches

---

## Fase 8: Teste do Pipeline CI/CD

### 8.1 Teste em DEV

```bash
# Criar branch de feature
git checkout develop
git checkout -b feature/test-cicd

# Fazer uma pequena alteração
echo "# Test" >> backend/README.md

# Commit e push
git add .
git commit -m "test: CI/CD pipeline"
git push origin feature/test-cicd

# Criar Pull Request para develop
# O deploy deve ocorrer automaticamente após merge
```

### 8.2 Teste em PROD

```bash
# Após testar em DEV, criar PR de develop para main
git checkout develop
git pull origin develop

# Criar PR: develop → main
# Após aprovação e merge, o deploy para produção será executado
```

---

## Fase 9: Monitoramento e Logs

### 9.1 Logs do Backend (Cloud Run)

```bash
# DEV
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=backend-api" \
  --limit 50 \
  --project ecosistema-imob-dev

# PROD
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=backend-api" \
  --limit 50 \
  --project ecosistema-imob-prod
```

### 9.2 Logs do Frontend (Vercel)

Acesse:
- DEV: https://vercel.com/[seu-time]/frontend-public/logs
- PROD: https://vercel.com/[seu-time]/frontend-public/logs

### 9.3 Configurar Alertas

**Cloud Run (GCP):**
1. Acesse Cloud Console > Monitoring > Alerting
2. Criar alertas para:
   - Latência > 1s
   - Taxa de erro > 5%
   - CPU > 80%
   - Memória > 80%

**Vercel:**
1. Acesse projeto > Settings > Monitoring
2. Configurar notificações para:
   - Build failures
   - Deployment failures
   - Performance degradation

---

## Fase 10: Checklist de Validação

### Backend (Cloud Run)

- [ ] DEV: Deploy automático no push para `develop`
- [ ] PROD: Deploy requer aprovação no push para `main`
- [ ] Testes executam antes do deploy
- [ ] Variáveis de ambiente configuradas corretamente
- [ ] Credenciais Firebase carregadas
- [ ] Domínio personalizado configurado
- [ ] HTTPS ativo
- [ ] Logs acessíveis

### Frontend Public

- [ ] DEV: Deploy automático no push para `develop`
- [ ] PROD: Deploy requer aprovação no push para `main`
- [ ] Variáveis de ambiente configuradas
- [ ] Domínio personalizado configurado
- [ ] SSR/SSG funcionando
- [ ] API URL correta
- [ ] Build otimizado

### Frontend Admin

- [ ] DEV: Deploy automático no push para `develop`
- [ ] PROD: Deploy requer aprovação no push para `main`
- [ ] Variáveis de ambiente configuradas
- [ ] Firebase Auth configurado
- [ ] Domínio personalizado configurado
- [ ] Autenticação funcionando
- [ ] API URL correta

---

## Comandos Úteis

### Verificar Status do Cloud Run

```bash
# DEV
gcloud run services describe backend-api \
  --region southamerica-east1 \
  --project ecosistema-imob-dev

# PROD
gcloud run services describe backend-api \
  --region southamerica-east1 \
  --project ecosistema-imob-prod
```

### Rollback em caso de problemas

```bash
# Listar revisões
gcloud run revisions list \
  --service backend-api \
  --region southamerica-east1 \
  --project ecosistema-imob-prod

# Fazer rollback para revisão anterior
gcloud run services update-traffic backend-api \
  --to-revisions REVISION-NAME=100 \
  --region southamerica-east1 \
  --project ecosistema-imob-prod
```

### Limpar builds antigas

```bash
# Listar imagens antigas
gcloud container images list-tags gcr.io/ecosistema-imob-prod/backend-api

# Deletar imagens antigas (manter últimas 5)
gcloud container images list-tags gcr.io/ecosistema-imob-prod/backend-api \
  --format="get(digest)" --filter="NOT tags:*" | \
  tail -n +6 | \
  xargs -I {} gcloud container images delete gcr.io/ecosistema-imob-prod/backend-api@{}
```

---

## Custos Estimados (Mensal)

### Ambiente DEV
- Cloud Run: ~$5-10 (baixo tráfego)
- Firestore: ~$5-10 (leituras/escritas limitadas)
- Cloud Storage: ~$2-5
- Vercel: $0 (plano hobby)
- **Total DEV: ~$12-25/mês**

### Ambiente PROD
- Cloud Run: ~$20-50 (depende do tráfego)
- Firestore: ~$50-200 (depende do volume)
- Cloud Storage: ~$10-30
- Vercel Pro: ~$20/mês
- **Total PROD: ~$100-300/mês**

---

## Próximos Passos

1. **Configurar monitoramento avançado**: Sentry, DataDog, ou New Relic
2. **Implementar testes E2E**: Cypress ou Playwright
3. **Configurar backup automático**: Firestore Export
4. **CDN**: Cloud CDN para assets estáticos
5. **Performance**: Otimizar Core Web Vitals
6. **Segurança**: WAF (Web Application Firewall)

---

## Suporte e Troubleshooting

### Problema: Deploy falha no GitHub Actions

1. Verificar logs do workflow no GitHub
2. Verificar se secrets estão configurados corretamente
3. Verificar permissões da service account no GCP

### Problema: Backend não conecta ao Firestore

1. Verificar se `FIREBASE_PROJECT_ID` está correto
2. Verificar se credenciais Firebase estão presentes
3. Verificar se Firestore API está habilitada

### Problema: Frontend não consegue acessar API

1. Verificar CORS no backend
2. Verificar se `NEXT_PUBLIC_API_URL` está correto
3. Verificar se Cloud Run permite tráfego não autenticado

---

**Versão**: 1.0
**Data**: 2026-01-07
**Autor**: Altatech Systems
