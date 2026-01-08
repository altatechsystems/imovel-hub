#!/bin/bash

# Script de configuração automática dos ambientes GCP (DEV e PROD)
# Uso: ./setup-gcp-environments.sh

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Variáveis
PROJECT_ID_DEV="ecosistema-imob-dev"
PROJECT_ID_PROD="ecosistema-imob-prod"
REGION="southamerica-east1"
BUCKET_DEV="${PROJECT_ID_DEV}-storage"
BUCKET_PROD="${PROJECT_ID_PROD}-storage"

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Configuração de Ambientes GCP${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Verificar se gcloud está instalado
if ! command -v gcloud &> /dev/null; then
    echo -e "${RED}Erro: gcloud CLI não está instalado${NC}"
    echo "Instale em: https://cloud.google.com/sdk/docs/install"
    exit 1
fi

echo -e "${YELLOW}[1/8] Verificando autenticação...${NC}"
gcloud auth list

echo ""
echo -e "${YELLOW}[2/8] Criando projeto PROD...${NC}"
if gcloud projects describe $PROJECT_ID_PROD &> /dev/null; then
    echo -e "${GREEN}✓ Projeto $PROJECT_ID_PROD já existe${NC}"
else
    gcloud projects create $PROJECT_ID_PROD --name="Ecosistema Imob - Produção"
    echo -e "${GREEN}✓ Projeto $PROJECT_ID_PROD criado${NC}"
fi

echo ""
echo -e "${YELLOW}[3/8] Habilitando APIs necessárias...${NC}"

# DEV
echo "  → Configurando DEV..."
gcloud config set project $PROJECT_ID_DEV
gcloud services enable \
  run.googleapis.com \
  firestore.googleapis.com \
  storage.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  iam.googleapis.com \
  --quiet

# PROD
echo "  → Configurando PROD..."
gcloud config set project $PROJECT_ID_PROD
gcloud services enable \
  run.googleapis.com \
  firestore.googleapis.com \
  storage.googleapis.com \
  cloudbuild.googleapis.com \
  artifactregistry.googleapis.com \
  cloudresourcemanager.googleapis.com \
  iam.googleapis.com \
  --quiet

echo -e "${GREEN}✓ APIs habilitadas${NC}"

echo ""
echo -e "${YELLOW}[4/8] Configurando Firestore...${NC}"

# DEV
gcloud config set project $PROJECT_ID_DEV
if gcloud firestore databases describe --database="\(default\)" &> /dev/null; then
    echo -e "${GREEN}✓ Firestore DEV já existe${NC}"
else
    gcloud firestore databases create \
      --location=$REGION \
      --type=firestore-native \
      --quiet
    echo -e "${GREEN}✓ Firestore DEV criado${NC}"
fi

# PROD
gcloud config set project $PROJECT_ID_PROD
if gcloud firestore databases describe --database="\(default\)" &> /dev/null; then
    echo -e "${GREEN}✓ Firestore PROD já existe${NC}"
else
    gcloud firestore databases create \
      --location=$REGION \
      --type=firestore-native \
      --quiet
    echo -e "${GREEN}✓ Firestore PROD criado${NC}"
fi

echo ""
echo -e "${YELLOW}[5/8] Criando buckets de storage...${NC}"

# DEV
gcloud config set project $PROJECT_ID_DEV
if gsutil ls -b gs://$BUCKET_DEV &> /dev/null; then
    echo -e "${GREEN}✓ Bucket DEV já existe${NC}"
else
    gsutil mb -c STANDARD -l $REGION gs://$BUCKET_DEV
    gsutil iam ch allUsers:objectViewer gs://$BUCKET_DEV
    echo -e "${GREEN}✓ Bucket DEV criado: gs://$BUCKET_DEV${NC}"
fi

# PROD
gcloud config set project $PROJECT_ID_PROD
if gsutil ls -b gs://$BUCKET_PROD &> /dev/null; then
    echo -e "${GREEN}✓ Bucket PROD já existe${NC}"
else
    gsutil mb -c STANDARD -l $REGION gs://$BUCKET_PROD
    gsutil iam ch allUsers:objectViewer gs://$BUCKET_PROD
    echo -e "${GREEN}✓ Bucket PROD criado: gs://$BUCKET_PROD${NC}"
fi

echo ""
echo -e "${YELLOW}[6/8] Criando Service Accounts para CI/CD...${NC}"

# DEV
gcloud config set project $PROJECT_ID_DEV
SA_EMAIL_DEV="github-actions-dev@${PROJECT_ID_DEV}.iam.gserviceaccount.com"

if gcloud iam service-accounts describe $SA_EMAIL_DEV &> /dev/null; then
    echo -e "${GREEN}✓ Service Account DEV já existe${NC}"
else
    gcloud iam service-accounts create github-actions-dev \
      --display-name="GitHub Actions - Development" \
      --quiet
    echo -e "${GREEN}✓ Service Account DEV criada${NC}"
fi

# PROD
gcloud config set project $PROJECT_ID_PROD
SA_EMAIL_PROD="github-actions-prod@${PROJECT_ID_PROD}.iam.gserviceaccount.com"

if gcloud iam service-accounts describe $SA_EMAIL_PROD &> /dev/null; then
    echo -e "${GREEN}✓ Service Account PROD já existe${NC}"
else
    gcloud iam service-accounts create github-actions-prod \
      --display-name="GitHub Actions - Production" \
      --quiet
    echo -e "${GREEN}✓ Service Account PROD criada${NC}"
fi

echo ""
echo -e "${YELLOW}[7/8] Atribuindo permissões...${NC}"

# DEV
gcloud config set project $PROJECT_ID_DEV
gcloud projects add-iam-policy-binding $PROJECT_ID_DEV \
  --member="serviceAccount:$SA_EMAIL_DEV" \
  --role="roles/run.admin" \
  --quiet

gcloud projects add-iam-policy-binding $PROJECT_ID_DEV \
  --member="serviceAccount:$SA_EMAIL_DEV" \
  --role="roles/iam.serviceAccountUser" \
  --quiet

gcloud projects add-iam-policy-binding $PROJECT_ID_DEV \
  --member="serviceAccount:$SA_EMAIL_DEV" \
  --role="roles/storage.admin" \
  --quiet

echo -e "${GREEN}✓ Permissões DEV configuradas${NC}"

# PROD
gcloud config set project $PROJECT_ID_PROD
gcloud projects add-iam-policy-binding $PROJECT_ID_PROD \
  --member="serviceAccount:$SA_EMAIL_PROD" \
  --role="roles/run.admin" \
  --quiet

gcloud projects add-iam-policy-binding $PROJECT_ID_PROD \
  --member="serviceAccount:$SA_EMAIL_PROD" \
  --role="roles/iam.serviceAccountUser" \
  --quiet

gcloud projects add-iam-policy-binding $PROJECT_ID_PROD \
  --member="serviceAccount:$SA_EMAIL_PROD" \
  --role="roles/storage.admin" \
  --quiet

echo -e "${GREEN}✓ Permissões PROD configuradas${NC}"

echo ""
echo -e "${YELLOW}[8/8] Gerando chaves das Service Accounts...${NC}"

# Criar diretório para as chaves
mkdir -p .gcp-keys
cd .gcp-keys

# DEV
gcloud config set project $PROJECT_ID_DEV
if [ -f "github-actions-dev-key.json" ]; then
    echo -e "${YELLOW}⚠ Chave DEV já existe, pulando...${NC}"
else
    gcloud iam service-accounts keys create github-actions-dev-key.json \
      --iam-account=$SA_EMAIL_DEV \
      --quiet
    echo -e "${GREEN}✓ Chave DEV gerada: .gcp-keys/github-actions-dev-key.json${NC}"
fi

# PROD
gcloud config set project $PROJECT_ID_PROD
if [ -f "github-actions-prod-key.json" ]; then
    echo -e "${YELLOW}⚠ Chave PROD já existe, pulando...${NC}"
else
    gcloud iam service-accounts keys create github-actions-prod-key.json \
      --iam-account=$SA_EMAIL_PROD \
      --quiet
    echo -e "${GREEN}✓ Chave PROD gerada: .gcp-keys/github-actions-prod-key.json${NC}"
fi

cd ..

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Configuração Concluída!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Próximos passos:${NC}"
echo ""
echo "1. Configurar Firebase Authentication:"
echo "   → Acesse: https://console.firebase.google.com"
echo "   → Adicione os projetos: $PROJECT_ID_DEV e $PROJECT_ID_PROD"
echo "   → Ative Authentication > Email/Password"
echo "   → Baixe as credenciais admin SDK"
echo ""
echo "2. Adicionar secrets no GitHub:"
echo "   → Acesse: https://github.com/SEU-USUARIO/ecosistema-imob/settings/secrets/actions"
echo "   → Adicione as chaves geradas em: .gcp-keys/"
echo ""
echo "3. Configurar Vercel:"
echo "   → Execute: npm install -g vercel"
echo "   → Execute: vercel login"
echo "   → Configure os projetos frontend-public e frontend-admin"
echo ""
echo "4. Consulte QUICK_START_CICD.md para instruções detalhadas"
echo ""
echo -e "${YELLOW}IMPORTANTE:${NC}"
echo -e "${RED}As chaves em .gcp-keys/ contêm credenciais sensíveis!${NC}"
echo -e "${RED}Adicione .gcp-keys/ ao .gitignore e não faça commit!${NC}"
echo ""
