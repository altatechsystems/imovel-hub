#!/bin/bash

# Script de configuração automática do Vercel
# Uso: ./setup-vercel.sh

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Configuração do Vercel${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""

# Verificar se vercel está instalado
if ! command -v vercel &> /dev/null; then
    echo -e "${YELLOW}Vercel CLI não encontrado. Instalando...${NC}"
    npm install -g vercel
    echo -e "${GREEN}✓ Vercel CLI instalado${NC}"
fi

echo ""
echo -e "${YELLOW}[1/4] Fazendo login no Vercel...${NC}"
vercel login

echo ""
echo -e "${YELLOW}[2/4] Configurando Frontend Public...${NC}"
cd frontend-public

if [ -f ".vercel/project.json" ]; then
    echo -e "${GREEN}✓ Frontend Public já está linkado${NC}"
else
    echo -e "${YELLOW}Configure o projeto:${NC}"
    echo "  → Escolha: Create new project"
    echo "  → Nome sugerido: ecosistema-imob-public"
    echo ""
    vercel link
    echo -e "${GREEN}✓ Frontend Public configurado${NC}"
fi

echo ""
echo "Configurando variáveis de ambiente para DEV (Preview)..."
vercel env add NEXT_PUBLIC_API_URL preview || true
vercel env add NEXT_PUBLIC_ENVIRONMENT preview <<< "development" || true

echo ""
echo "Configurando variáveis de ambiente para PROD (Production)..."
vercel env add NEXT_PUBLIC_API_URL production || true
vercel env add NEXT_PUBLIC_ENVIRONMENT production <<< "production" || true

cd ..

echo ""
echo -e "${YELLOW}[3/4] Configurando Frontend Admin...${NC}"
cd frontend-admin

if [ -f ".vercel/project.json" ]; then
    echo -e "${GREEN}✓ Frontend Admin já está linkado${NC}"
else
    echo -e "${YELLOW}Configure o projeto:${NC}"
    echo "  → Escolha: Create new project"
    echo "  → Nome sugerido: ecosistema-imob-admin"
    echo ""
    vercel link
    echo -e "${GREEN}✓ Frontend Admin configurado${NC}"
fi

echo ""
echo "Configurando variáveis de ambiente para DEV (Preview)..."
vercel env add NEXT_PUBLIC_API_URL preview || true
vercel env add NEXT_PUBLIC_FIREBASE_API_KEY preview || true
vercel env add NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN preview || true
vercel env add NEXT_PUBLIC_FIREBASE_PROJECT_ID preview || true
vercel env add NEXT_PUBLIC_ENVIRONMENT preview <<< "development" || true

echo ""
echo "Configurando variáveis de ambiente para PROD (Production)..."
vercel env add NEXT_PUBLIC_API_URL production || true
vercel env add NEXT_PUBLIC_FIREBASE_API_KEY production || true
vercel env add NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN production || true
vercel env add NEXT_PUBLIC_FIREBASE_PROJECT_ID production || true
vercel env add NEXT_PUBLIC_ENVIRONMENT production <<< "production" || true

cd ..

echo ""
echo -e "${YELLOW}[4/4] Obtendo informações dos projetos...${NC}"

# Obter project IDs
cd frontend-public
PROJECT_ID_PUBLIC=$(vercel inspect --token=$(vercel whoami -t) 2>/dev/null | grep "ID:" | awk '{print $2}' || echo "N/A")
cd ..

cd frontend-admin
PROJECT_ID_ADMIN=$(vercel inspect --token=$(vercel whoami -t) 2>/dev/null | grep "ID:" | awk '{print $2}' || echo "N/A")
cd ..

# Obter token
echo ""
echo -e "${YELLOW}Para obter seu token do Vercel:${NC}"
echo "1. Acesse: https://vercel.com/account/tokens"
echo "2. Crie um novo token"
echo "3. Copie o token gerado"

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Configuração do Vercel Concluída!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo -e "${YELLOW}Informações para GitHub Secrets:${NC}"
echo ""
echo "VERCEL_TOKEN=<obter em https://vercel.com/account/tokens>"
echo "VERCEL_ORG_ID=<obter no Vercel Dashboard > Settings > General>"
echo "VERCEL_PROJECT_ID_PUBLIC=$PROJECT_ID_PUBLIC"
echo "VERCEL_PROJECT_ID_ADMIN=$PROJECT_ID_ADMIN"
echo ""
echo -e "${YELLOW}Próximos passos:${NC}"
echo ""
echo "1. Adicione os secrets acima no GitHub:"
echo "   → https://github.com/SEU-USUARIO/ecosistema-imob/settings/secrets/actions"
echo ""
echo "2. Atualize as variáveis de ambiente no Vercel Dashboard:"
echo "   → Frontend Public: https://vercel.com/dashboard"
echo "   → Frontend Admin: https://vercel.com/dashboard"
echo ""
echo "3. Faça o primeiro deploy:"
echo "   → git push origin develop"
echo ""
