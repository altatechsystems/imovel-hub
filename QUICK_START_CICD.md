# Quick Start - Configuração CI/CD

Guia rápido para começar a usar CI/CD. Para instruções completas, consulte [CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md).

## Pré-requisitos

- [ ] Conta no Google Cloud Platform (GCP)
- [ ] Projeto GCP dev criado: `ecosistema-imob-dev`
- [ ] Conta no GitHub com acesso ao repositório
- [ ] Conta no Vercel
- [ ] gcloud CLI instalado localmente
- [ ] Node.js 20+ instalado
- [ ] Go 1.21+ instalado

## Checklist de Configuração Rápida

### 1. GCP - Criar Projeto de Produção (5 min)

```bash
# Criar projeto
gcloud projects create ecosistema-imob-prod --name="Ecosistema Imob - Produção"

# Habilitar APIs
gcloud config set project ecosistema-imob-prod
gcloud services enable run.googleapis.com firestore.googleapis.com \
  storage.googleapis.com cloudbuild.googleapis.com \
  artifactregistry.googleapis.com
```

### 2. Firestore - Criar Banco PROD (2 min)

```bash
gcloud firestore databases create \
  --location=southamerica-east1 \
  --type=firestore-native \
  --project=ecosistema-imob-prod
```

### 3. Cloud Storage - Criar Buckets (2 min)

```bash
# DEV
gcloud config set project ecosistema-imob-dev
gsutil mb -l southamerica-east1 gs://ecosistema-imob-dev-storage
gsutil iam ch allUsers:objectViewer gs://ecosistema-imob-dev-storage

# PROD
gcloud config set project ecosistema-imob-prod
gsutil mb -l southamerica-east1 gs://ecosistema-imob-prod-storage
gsutil iam ch allUsers:objectViewer gs://ecosistema-imob-prod-storage
```

### 4. Service Accounts - Criar e Gerar Chaves (5 min)

```bash
# DEV
gcloud config set project ecosistema-imob-dev
gcloud iam service-accounts create github-actions-dev
gcloud projects add-iam-policy-binding ecosistema-imob-dev \
  --member="serviceAccount:github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com" \
  --role="roles/run.admin"
gcloud iam service-accounts keys create github-actions-dev-key.json \
  --iam-account=github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com

# PROD
gcloud config set project ecosistema-imob-prod
gcloud iam service-accounts create github-actions-prod
gcloud projects add-iam-policy-binding ecosistema-imob-prod \
  --member="serviceAccount:github-actions-prod@ecosistema-imob-prod.iam.gserviceaccount.com" \
  --role="roles/run.admin"
gcloud iam service-accounts keys create github-actions-prod-key.json \
  --iam-account=github-actions-prod@ecosistema-imob-prod.iam.gserviceaccount.com
```

### 5. Firebase - Configurar Autenticação (5 min)

1. Acesse [Firebase Console](https://console.firebase.google.com)
2. Adicione o projeto: `ecosistema-imob-prod`
3. Ative Authentication > Email/Password
4. Baixe as credenciais:
   - Configurações > Contas de Serviço > Gerar nova chave privada
   - Salve como: `firebase-adminsdk-prod.json`

### 6. GitHub Secrets - Adicionar Credenciais (10 min)

Acesse: `https://github.com/SEU-USUARIO/ecosistema-imob/settings/secrets/actions`

**Adicione estes secrets:**

```
# GCP DEV
GCP_PROJECT_ID_DEV=ecosistema-imob-dev
GCP_SA_KEY_DEV=<conteúdo do github-actions-dev-key.json>

# GCP PROD
GCP_PROJECT_ID_PROD=ecosistema-imob-prod
GCP_SA_KEY_PROD=<conteúdo do github-actions-prod-key.json>

# Firebase DEV
FIREBASE_PROJECT_ID_DEV=ecosistema-imob-dev
FIREBASE_ADMIN_SDK_DEV=<conteúdo do firebase-adminsdk-dev.json>
FIREBASE_API_KEY_DEV=<obter no Firebase Console>
FIREBASE_AUTH_DOMAIN_DEV=ecosistema-imob-dev.firebaseapp.com

# Firebase PROD
FIREBASE_PROJECT_ID_PROD=ecosistema-imob-prod
FIREBASE_ADMIN_SDK_PROD=<conteúdo do firebase-adminsdk-prod.json>
FIREBASE_API_KEY_PROD=<obter no Firebase Console>
FIREBASE_AUTH_DOMAIN_PROD=ecosistema-imob-prod.firebaseapp.com

# Storage
GCS_BUCKET_NAME_DEV=ecosistema-imob-dev-storage
GCS_BUCKET_NAME_PROD=ecosistema-imob-prod-storage

# Vercel (obter em vercel.com/account/tokens)
VERCEL_TOKEN=<seu token do Vercel>

# URLs do Backend (serão atualizadas após primeiro deploy)
NEXT_PUBLIC_API_URL_DEV=https://backend-api-xxxxx-uc.a.run.app
NEXT_PUBLIC_API_URL_PROD=https://backend-api-xxxxx-uc.a.run.app
```

### 7. Vercel - Configurar Projetos (5 min)

```bash
# Instalar CLI
npm install -g vercel

# Login
vercel login

# Configurar Frontend Public
cd frontend-public
vercel link
# Escolha: Create new project
# Nome: ecosistema-imob-public

# Configurar Frontend Admin
cd ../frontend-admin
vercel link
# Escolha: Create new project
# Nome: ecosistema-imob-admin
```

### 8. GitHub - Proteger Branch Main (2 min)

1. Acesse: `Settings > Branches`
2. Add rule: `main`
3. Marque:
   - ✅ Require a pull request before merging
   - ✅ Require approvals (1)
   - ✅ Require status checks to pass

### 9. GitHub - Criar Ambiente Production (2 min)

1. Acesse: `Settings > Environments`
2. New environment: `production`
3. Marque:
   - ✅ Required reviewers (adicione seu usuário)
4. Deployment branches: `Only protected branches`

### 10. Testar Pipeline (5 min)

```bash
# Criar branch de teste
git checkout -b feature/test-cicd

# Fazer alteração mínima
echo "# CI/CD Test" >> README.md

# Commit e push
git add .
git commit -m "test: CI/CD pipeline setup"
git push origin feature/test-cicd

# Criar PR no GitHub: feature/test-cicd → develop
# Aguardar CI passar
# Fazer merge
# Deploy automático para DEV será executado
```

## Ordem de Primeiro Deploy

**IMPORTANTE**: Siga esta ordem para o primeiro deploy:

1. **Backend DEV** (primeiro)
   - Push para `develop` com alteração em `backend/`
   - Aguardar deploy completar
   - Copiar URL do Cloud Run
   - Atualizar secret: `NEXT_PUBLIC_API_URL_DEV`

2. **Frontend Public DEV**
   - Push para `develop` com alteração em `frontend-public/`
   - Deploy automático

3. **Frontend Admin DEV**
   - Push para `develop` com alteração em `frontend-admin/`
   - Deploy automático

4. **Testar DEV completo**
   - Acessar frontend public
   - Acessar frontend admin
   - Testar autenticação

5. **Backend PROD**
   - Criar PR: `develop` → `main`
   - Aguardar aprovação
   - Merge
   - Deploy para PROD

6. **Frontend PROD**
   - Mesmo processo
   - Atualizar `NEXT_PUBLIC_API_URL_PROD` se necessário

## URLs Esperadas Após Deploy

### DEV
- Backend: `https://backend-api-xxxxx-uc.a.run.app`
- Frontend Public: `https://ecosistema-imob-public-xxxxx.vercel.app`
- Frontend Admin: `https://ecosistema-imob-admin-xxxxx.vercel.app`

### PROD
- Backend: `https://backend-api-xxxxx-uc.a.run.app`
- Frontend Public: `https://ecosistema-imob-public.vercel.app`
- Frontend Admin: `https://ecosistema-imob-admin.vercel.app`

## Comandos Úteis

### Ver logs do backend
```bash
# DEV
gcloud logging read "resource.type=cloud_run_revision" \
  --limit 50 \
  --project ecosistema-imob-dev

# PROD
gcloud logging read "resource.type=cloud_run_revision" \
  --limit 50 \
  --project ecosistema-imob-prod
```

### Verificar status do Cloud Run
```bash
gcloud run services list --project ecosistema-imob-dev
gcloud run services list --project ecosistema-imob-prod
```

### Ver deployments no Vercel
```bash
vercel list
```

## Troubleshooting Rápido

### ❌ Deploy falha: "Permission denied"
→ Verificar se service account tem permissões corretas:
```bash
gcloud projects get-iam-policy ecosistema-imob-dev
```

### ❌ Backend não inicia: "Firebase error"
→ Verificar se secret `FIREBASE_ADMIN_SDK_DEV` está correto

### ❌ Frontend não conecta ao backend
→ Verificar se `NEXT_PUBLIC_API_URL_DEV` está atualizado com URL do Cloud Run

### ❌ GitHub Actions falha: "Secret not found"
→ Verificar se todos os secrets foram adicionados corretamente

### ❌ Vercel build falha
→ Verificar se `package-lock.json` está commitado
→ Verificar variáveis de ambiente no Vercel Dashboard

## Próximos Passos

Após configuração inicial:

1. [ ] Configurar domínios personalizados
2. [ ] Adicionar monitoramento (Cloud Monitoring)
3. [ ] Configurar alertas
4. [ ] Implementar testes E2E no CI
5. [ ] Configurar backup automático do Firestore

## Suporte

Para mais detalhes, consulte:
- [CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md) - Documentação completa
- [README.md](README.md) - Visão geral do projeto

---

**Tempo estimado de configuração**: 45-60 minutos
