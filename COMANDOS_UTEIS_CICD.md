# Comandos Úteis - CI/CD

Referência rápida de comandos para gerenciar a infraestrutura e troubleshooting.

## Índice
- [GCP / Cloud Run](#gcp--cloud-run)
- [Firestore](#firestore)
- [Cloud Storage](#cloud-storage)
- [GitHub Actions](#github-actions)
- [Vercel](#vercel)
- [Docker Local](#docker-local)
- [Troubleshooting](#troubleshooting)

---

## GCP / Cloud Run

### Autenticação e Configuração

```bash
# Login
gcloud auth login

# Listar projetos
gcloud projects list

# Definir projeto padrão
gcloud config set project ecosistema-imob-dev

# Ver configuração atual
gcloud config list

# Listar regiões disponíveis
gcloud compute regions list
```

### Cloud Run - Deploy Manual

```bash
# Deploy do backend (DEV)
cd backend
gcloud builds submit --tag gcr.io/ecosistema-imob-dev/backend-api
gcloud run deploy backend-api \
  --image gcr.io/ecosistema-imob-dev/backend-api \
  --platform managed \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --project ecosistema-imob-dev

# Deploy do backend (PROD)
gcloud builds submit --tag gcr.io/ecosistema-imob-prod/backend-api \
  --project ecosistema-imob-prod
gcloud run deploy backend-api \
  --image gcr.io/ecosistema-imob-prod/backend-api \
  --platform managed \
  --region southamerica-east1 \
  --allow-unauthenticated \
  --project ecosistema-imob-prod
```

### Cloud Run - Gerenciamento

```bash
# Listar serviços
gcloud run services list --project ecosistema-imob-dev

# Descrever serviço
gcloud run services describe backend-api \
  --region southamerica-east1 \
  --project ecosistema-imob-dev

# Obter URL do serviço
gcloud run services describe backend-api \
  --region southamerica-east1 \
  --format 'value(status.url)' \
  --project ecosistema-imob-dev

# Listar revisões
gcloud run revisions list \
  --service backend-api \
  --region southamerica-east1 \
  --project ecosistema-imob-dev

# Deletar serviço
gcloud run services delete backend-api \
  --region southamerica-east1 \
  --project ecosistema-imob-dev
```

### Cloud Run - Tráfego e Rollback

```bash
# Ver distribuição de tráfego
gcloud run services describe backend-api \
  --region southamerica-east1 \
  --format 'value(status.traffic)' \
  --project ecosistema-imob-prod

# Rollback para revisão anterior (100% do tráfego)
gcloud run services update-traffic backend-api \
  --to-revisions REVISION-NAME=100 \
  --region southamerica-east1 \
  --project ecosistema-imob-prod

# Split de tráfego (canary deploy)
gcloud run services update-traffic backend-api \
  --to-revisions REVISION-NEW=10,REVISION-OLD=90 \
  --region southamerica-east1 \
  --project ecosistema-imob-prod
```

### Cloud Run - Logs

```bash
# Logs em tempo real (DEV)
gcloud logging tail "resource.type=cloud_run_revision AND resource.labels.service_name=backend-api" \
  --project ecosistema-imob-dev

# Logs das últimas 50 linhas
gcloud logging read "resource.type=cloud_run_revision AND resource.labels.service_name=backend-api" \
  --limit 50 \
  --project ecosistema-imob-dev

# Logs com filtro de severidade
gcloud logging read "resource.type=cloud_run_revision AND severity>=ERROR" \
  --limit 50 \
  --project ecosistema-imob-prod

# Logs de uma data específica
gcloud logging read "resource.type=cloud_run_revision AND timestamp>=\"2026-01-07T00:00:00Z\"" \
  --limit 50 \
  --project ecosistema-imob-prod
```

### Cloud Run - Métricas

```bash
# Ver métricas via console
# https://console.cloud.google.com/run?project=ecosistema-imob-prod

# Exportar métricas
gcloud monitoring time-series list \
  --filter='metric.type="run.googleapis.com/request_count"' \
  --project ecosistema-imob-prod
```

---

## Firestore

### Gerenciamento

```bash
# Criar banco de dados
gcloud firestore databases create \
  --location=southamerica-east1 \
  --type=firestore-native \
  --project ecosistema-imob-prod

# Descrever banco
gcloud firestore databases describe \
  --database=\(default\) \
  --project ecosistema-imob-dev

# Listar índices
gcloud firestore indexes composite list \
  --project ecosistema-imob-dev

# Criar índice
gcloud firestore indexes composite create \
  --collection-group=properties \
  --field-config field-path=tenant_id,order=ASCENDING \
  --field-config field-path=status,order=ASCENDING \
  --project ecosistema-imob-dev
```

### Backup e Restore

```bash
# Export (backup) - DEV
gcloud firestore export gs://ecosistema-imob-dev-backup \
  --project ecosistema-imob-dev

# Export (backup) - PROD
gcloud firestore export gs://ecosistema-imob-prod-backup \
  --project ecosistema-imob-prod

# Import (restore)
gcloud firestore import gs://ecosistema-imob-prod-backup/[TIMESTAMP] \
  --project ecosistema-imob-prod

# Listar backups
gsutil ls gs://ecosistema-imob-prod-backup/
```

### Operações de Dados

```bash
# Deletar documento (via gcloud não é direto, usar Firebase CLI)
firebase firestore:delete tenants/tenant123 \
  --project ecosistema-imob-dev \
  --recursive

# Deletar coleção inteira (com confirmação)
firebase firestore:delete --all-collections \
  --project ecosistema-imob-dev
```

---

## Cloud Storage

### Gerenciamento de Buckets

```bash
# Criar bucket
gsutil mb -c STANDARD -l southamerica-east1 gs://ecosistema-imob-dev-storage

# Listar buckets
gsutil ls

# Ver detalhes do bucket
gsutil ls -L -b gs://ecosistema-imob-dev-storage

# Deletar bucket
gsutil rm -r gs://ecosistema-imob-dev-storage
```

### Permissões

```bash
# Tornar bucket público para leitura
gsutil iam ch allUsers:objectViewer gs://ecosistema-imob-dev-storage

# Remover acesso público
gsutil iam ch -d allUsers:objectViewer gs://ecosistema-imob-dev-storage

# Ver IAM do bucket
gsutil iam get gs://ecosistema-imob-dev-storage
```

### Upload e Download

```bash
# Upload de arquivo
gsutil cp local-file.jpg gs://ecosistema-imob-dev-storage/images/

# Upload de diretório
gsutil -m cp -r local-directory/* gs://ecosistema-imob-dev-storage/images/

# Download de arquivo
gsutil cp gs://ecosistema-imob-dev-storage/images/file.jpg ./

# Sincronizar diretório
gsutil -m rsync -r local-directory gs://ecosistema-imob-dev-storage/images/
```

### Listar e Deletar

```bash
# Listar arquivos
gsutil ls gs://ecosistema-imob-dev-storage/

# Listar com detalhes
gsutil ls -l gs://ecosistema-imob-dev-storage/**

# Deletar arquivo
gsutil rm gs://ecosistema-imob-dev-storage/images/file.jpg

# Deletar com wildcard
gsutil -m rm gs://ecosistema-imob-dev-storage/images/*.jpg
```

---

## GitHub Actions

### Via GitHub CLI

```bash
# Instalar GitHub CLI
# Windows: winget install GitHub.cli
# Mac: brew install gh
# Linux: apt install gh

# Login
gh auth login

# Listar workflows
gh workflow list

# Ver status de workflow
gh run list --workflow=backend-deploy.yml

# Ver detalhes de uma run
gh run view RUN_ID

# Ver logs de uma run
gh run view RUN_ID --log

# Baixar logs
gh run download RUN_ID

# Cancelar run
gh run cancel RUN_ID

# Re-executar workflow
gh run rerun RUN_ID
```

### Executar Workflow Manualmente

```bash
# Via GitHub CLI
gh workflow run backend-deploy.yml

# Via interface web
# https://github.com/seu-usuario/ecosistema-imob/actions
# Selecione workflow > "Run workflow"
```

### Secrets

```bash
# Listar secrets
gh secret list

# Adicionar secret
gh secret set SECRET_NAME < secret-file.txt

# Deletar secret
gh secret delete SECRET_NAME
```

---

## Vercel

### Autenticação

```bash
# Instalar CLI
npm install -g vercel

# Login
vercel login

# Ver usuário logado
vercel whoami

# Logout
vercel logout
```

### Gerenciamento de Projetos

```bash
# Listar projetos
vercel list

# Linkar projeto ao diretório atual
vercel link

# Ver informações do projeto
vercel inspect

# Remover link
vercel unlink
```

### Deploy

```bash
# Deploy para preview (DEV)
vercel

# Deploy para produção
vercel --prod

# Deploy sem interação
vercel --yes

# Build local sem deploy
vercel build
```

### Logs e Monitoramento

```bash
# Ver deployments
vercel list

# Ver logs de deployment específico
vercel logs DEPLOYMENT-URL

# Ver logs em tempo real
vercel logs --follow

# Inspecionar deployment
vercel inspect DEPLOYMENT-URL
```

### Environment Variables

```bash
# Listar env vars
vercel env ls

# Adicionar env var
vercel env add VARIABLE_NAME

# Remover env var
vercel env rm VARIABLE_NAME

# Pull env vars para local
vercel env pull
```

### Domínios

```bash
# Listar domínios
vercel domains ls

# Adicionar domínio
vercel domains add example.com

# Remover domínio
vercel domains rm example.com

# Ver configuração de DNS
vercel domains inspect example.com
```

### Rollback

```bash
# Promover deployment anterior para produção
vercel rollback PREVIOUS-DEPLOYMENT-URL

# Via Vercel Dashboard
# https://vercel.com/[seu-time]/[projeto]/deployments
# Clique em "..." no deployment anterior > "Promote to Production"
```

---

## Docker Local

### Build e Teste Local

```bash
# Build da imagem do backend
cd backend
docker build -t ecosistema-backend:local .

# Executar container localmente
docker run -p 8080:8080 \
  -e ENVIRONMENT=development \
  -e FIREBASE_PROJECT_ID=ecosistema-imob-dev \
  -v $(pwd)/config:/home/appuser/config \
  ecosistema-backend:local

# Ver logs do container
docker logs CONTAINER_ID

# Parar container
docker stop CONTAINER_ID

# Remover container
docker rm CONTAINER_ID

# Remover imagem
docker rmi ecosistema-backend:local
```

### Docker Compose (desenvolvimento local)

```yaml
# Criar arquivo docker-compose.yml na raiz do backend
version: '3.8'
services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=development
      - FIREBASE_PROJECT_ID=ecosistema-imob-dev
      - GCS_BUCKET_NAME=ecosistema-imob-dev-storage
    volumes:
      - ./config:/home/appuser/config
```

```bash
# Iniciar
docker-compose up

# Iniciar em background
docker-compose up -d

# Parar
docker-compose down

# Ver logs
docker-compose logs -f
```

---

## Troubleshooting

### Verificar Status dos Serviços

```bash
# Backend DEV
curl https://backend-api-xxxxx.run.app/health

# Backend PROD
curl https://api.example.com/health

# Com detalhes
curl -v https://backend-api-xxxxx.run.app/health

# Com headers
curl -i https://backend-api-xxxxx.run.app/health
```

### Testar Autenticação

```bash
# Obter token Firebase (via Firebase CLI)
firebase login:ci

# Testar endpoint protegido
curl -H "Authorization: Bearer YOUR_TOKEN" \
  https://backend-api-xxxxx.run.app/api/v1/protected
```

### Verificar Conectividade

```bash
# Ping (não funciona com Cloud Run, usar curl)
curl -I https://backend-api-xxxxx.run.app

# Traceroute DNS
nslookup backend-api-xxxxx.run.app

# Verificar SSL
openssl s_client -connect backend-api-xxxxx.run.app:443
```

### Limpar Resources

```bash
# Deletar imagens antigas do Container Registry
gcloud container images list-tags gcr.io/ecosistema-imob-dev/backend-api

# Deletar imagens não tagueadas
gcloud container images list-tags gcr.io/ecosistema-imob-dev/backend-api \
  --format="get(digest)" \
  --filter="NOT tags:*" | \
  xargs -I {} gcloud container images delete gcr.io/ecosistema-imob-dev/backend-api@{} --quiet

# Limpar Cloud Build cache
gcloud builds list --limit=100 --project=ecosistema-imob-dev
```

### Verificar Permissões

```bash
# Ver IAM do projeto
gcloud projects get-iam-policy ecosistema-imob-dev

# Ver permissões de uma service account
gcloud projects get-iam-policy ecosistema-imob-dev \
  --flatten="bindings[].members" \
  --filter="bindings.members:serviceAccount:github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com"

# Testar permissões
gcloud projects test-iam-permissions ecosistema-imob-dev \
  --permissions=run.services.create,run.services.update
```

### Análise de Custos

```bash
# Ver billing account
gcloud billing accounts list

# Ver uso de recursos (via console)
# https://console.cloud.google.com/billing

# Export de billing para BigQuery (configurar via console)
# https://console.cloud.google.com/billing/export
```

### Performance Testing

```bash
# Teste de carga simples com Apache Bench
ab -n 1000 -c 10 https://backend-api-xxxxx.run.app/health

# Teste com wrk
wrk -t12 -c400 -d30s https://backend-api-xxxxx.run.app/health

# Teste de latência
time curl https://backend-api-xxxxx.run.app/api/v1/properties
```

---

## Scripts Úteis

### Script de Health Check

```bash
#!/bin/bash
# check-health.sh

URLS=(
  "https://backend-api-dev-xxxxx.run.app/health"
  "https://backend-api-prod-xxxxx.run.app/health"
  "https://dev.example.com"
  "https://www.example.com"
)

for url in "${URLS[@]}"; do
  echo "Checking $url..."
  status=$(curl -s -o /dev/null -w "%{http_code}" "$url")
  if [ $status -eq 200 ]; then
    echo "✓ $url is healthy"
  else
    echo "✗ $url returned $status"
  fi
done
```

### Script de Backup Automático

```bash
#!/bin/bash
# backup-firestore.sh

PROJECT_ID="ecosistema-imob-prod"
BUCKET="gs://ecosistema-imob-prod-backup"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "Starting Firestore backup for $PROJECT_ID..."

gcloud firestore export $BUCKET/$TIMESTAMP \
  --project=$PROJECT_ID

echo "Backup completed: $BUCKET/$TIMESTAMP"

# Limpar backups antigos (manter últimos 30 dias)
gsutil ls -l $BUCKET/ | awk '$1 ~ /^[0-9]/ {print $3}' | \
  sort -r | tail -n +31 | \
  xargs -I {} gsutil -m rm -r {}
```

---

## Referências Rápidas

### URLs Importantes

- **GCP Console**: https://console.cloud.google.com
- **Firebase Console**: https://console.firebase.google.com
- **Vercel Dashboard**: https://vercel.com/dashboard
- **GitHub Actions**: https://github.com/seu-usuario/ecosistema-imob/actions

### Documentação Oficial

- **Cloud Run**: https://cloud.google.com/run/docs
- **Firestore**: https://cloud.google.com/firestore/docs
- **Cloud Storage**: https://cloud.google.com/storage/docs
- **GitHub Actions**: https://docs.github.com/en/actions
- **Vercel**: https://vercel.com/docs

### Documentação do Projeto

- [CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md)
- [QUICK_START_CICD.md](QUICK_START_CICD.md)
- [CHECKLIST_CONFIGURACAO_CICD.md](CHECKLIST_CONFIGURACAO_CICD.md)
- [INFRASTRUCTURE_OVERVIEW.md](INFRASTRUCTURE_OVERVIEW.md)

---

**Dica**: Adicione estes comandos ao seu `.bash_aliases` ou `.zshrc` para acesso rápido!

```bash
# Adicione ao ~/.bash_aliases ou ~/.zshrc
alias gcp-dev='gcloud config set project ecosistema-imob-dev'
alias gcp-prod='gcloud config set project ecosistema-imob-prod'
alias cr-logs-dev='gcloud logging tail "resource.type=cloud_run_revision" --project ecosistema-imob-dev'
alias cr-logs-prod='gcloud logging tail "resource.type=cloud_run_revision" --project ecosistema-imob-prod'
alias v-deploy='vercel --prod'
alias gh-runs='gh run list --limit 5'
```
