# Checklist de Configuração CI/CD

Use este checklist para acompanhar o progresso da configuração dos ambientes dev e prod com CI/CD.

## Fase 1: Preparação Inicial (15 min)

### Pré-requisitos
- [ ] Conta no Google Cloud Platform (GCP) criada
- [ ] Projeto GCP dev existe: `ecosistema-imob-dev`
- [ ] Conta no GitHub com acesso admin ao repositório
- [ ] Conta no Vercel criada
- [ ] gcloud CLI instalado e autenticado (`gcloud auth login`)
- [ ] Node.js 20+ instalado (`node --version`)
- [ ] Go 1.21+ instalado (`go version`)

### Verificar Instalações
```bash
# Verificar ferramentas
gcloud --version
node --version
npm --version
go version
git --version
```

---

## Fase 2: Google Cloud Platform (30 min)

### Projeto PROD
- [ ] Criar projeto GCP produção: `ecosistema-imob-prod`
- [ ] Vincular conta de faturamento ao projeto prod
- [ ] Definir projeto como padrão no gcloud

### APIs do GCP
- [ ] Habilitar Cloud Run API (dev e prod)
- [ ] Habilitar Firestore API (dev e prod)
- [ ] Habilitar Cloud Storage API (dev e prod)
- [ ] Habilitar Cloud Build API (dev e prod)
- [ ] Habilitar Artifact Registry API (dev e prod)
- [ ] Habilitar IAM API (dev e prod)

**Comando rápido:**
```bash
# Executar script automatizado
./scripts/setup-gcp-environments.sh
```

### Firestore
- [ ] Criar banco Firestore DEV (se não existir)
- [ ] Criar banco Firestore PROD
- [ ] Verificar location: `southamerica-east1`

### Cloud Storage
- [ ] Criar bucket DEV: `ecosistema-imob-dev-storage`
- [ ] Criar bucket PROD: `ecosistema-imob-prod-storage`
- [ ] Configurar acesso público nos buckets (objectViewer)
- [ ] Verificar region: `southamerica-east1`

### Service Accounts
- [ ] Criar SA para GitHub Actions DEV: `github-actions-dev`
- [ ] Criar SA para GitHub Actions PROD: `github-actions-prod`
- [ ] Atribuir role `roles/run.admin` (dev e prod)
- [ ] Atribuir role `roles/iam.serviceAccountUser` (dev e prod)
- [ ] Atribuir role `roles/storage.admin` (dev e prod)

### Chaves das Service Accounts
- [ ] Gerar chave JSON para SA DEV
- [ ] Gerar chave JSON para SA PROD
- [ ] Salvar chaves em local seguro (NÃO commitar!)
- [ ] Verificar que `.gcp-keys/` está no `.gitignore`

---

## Fase 3: Firebase Authentication (15 min)

### Firebase Console
- [ ] Acessar [Firebase Console](https://console.firebase.google.com)
- [ ] Verificar projeto DEV: `ecosistema-imob-dev`
- [ ] Adicionar projeto PROD: `ecosistema-imob-prod`

### Authentication
- [ ] Ativar Email/Password no projeto DEV
- [ ] Ativar Email/Password no projeto PROD
- [ ] Configurar domínios autorizados DEV
- [ ] Configurar domínios autorizados PROD

### Credenciais Firebase
- [ ] Baixar Firebase Admin SDK DEV (JSON)
- [ ] Baixar Firebase Admin SDK PROD (JSON)
- [ ] Copiar API Key do projeto DEV (Project Settings)
- [ ] Copiar API Key do projeto PROD (Project Settings)
- [ ] Salvar credenciais em local seguro

---

## Fase 4: Repositório GitHub (20 min)

### Estrutura de Branches
- [ ] Verificar branch `main` existe
- [ ] Criar branch `develop` (se não existir)
- [ ] Fazer push dos workflows para `develop`

### GitHub Secrets
Acessar: `Settings > Secrets and variables > Actions`

#### GCP Secrets
- [ ] `GCP_PROJECT_ID_DEV` = `ecosistema-imob-dev`
- [ ] `GCP_SA_KEY_DEV` = conteúdo do JSON da SA DEV
- [ ] `GCP_PROJECT_ID_PROD` = `ecosistema-imob-prod`
- [ ] `GCP_SA_KEY_PROD` = conteúdo do JSON da SA PROD

#### Firebase Secrets DEV
- [ ] `FIREBASE_PROJECT_ID_DEV` = `ecosistema-imob-dev`
- [ ] `FIREBASE_ADMIN_SDK_DEV` = conteúdo do JSON Admin SDK DEV
- [ ] `FIREBASE_API_KEY_DEV` = API Key do Firebase DEV
- [ ] `FIREBASE_AUTH_DOMAIN_DEV` = `ecosistema-imob-dev.firebaseapp.com`

#### Firebase Secrets PROD
- [ ] `FIREBASE_PROJECT_ID_PROD` = `ecosistema-imob-prod`
- [ ] `FIREBASE_ADMIN_SDK_PROD` = conteúdo do JSON Admin SDK PROD
- [ ] `FIREBASE_API_KEY_PROD` = API Key do Firebase PROD
- [ ] `FIREBASE_AUTH_DOMAIN_PROD` = `ecosistema-imob-prod.firebaseapp.com`

#### Storage Secrets
- [ ] `GCS_BUCKET_NAME_DEV` = `ecosistema-imob-dev-storage`
- [ ] `GCS_BUCKET_NAME_PROD` = `ecosistema-imob-prod-storage`

#### Vercel Secrets (adicionar depois)
- [ ] `VERCEL_TOKEN` = (obter após configurar Vercel)
- [ ] `NEXT_PUBLIC_API_URL_DEV` = (obter após primeiro deploy)
- [ ] `NEXT_PUBLIC_API_URL_PROD` = (obter após primeiro deploy)

### Branch Protection
Acessar: `Settings > Branches > Add rule`

#### Proteger `main`
- [ ] Branch name pattern: `main`
- [ ] Require a pull request before merging
- [ ] Require approvals: 1
- [ ] Require status checks to pass before merging
- [ ] Require conversation resolution before merging
- [ ] Do not allow bypassing the above settings

### GitHub Environment
Acessar: `Settings > Environments > New environment`

#### Criar ambiente `production`
- [ ] Nome: `production`
- [ ] Required reviewers: adicionar seu usuário
- [ ] Deployment branches: Only protected branches
- [ ] Environment URL: `https://api.example.com` (opcional)

---

## Fase 5: Vercel (15 min)

### Instalação e Login
- [ ] Instalar Vercel CLI: `npm install -g vercel`
- [ ] Fazer login: `vercel login`
- [ ] Verificar autenticação: `vercel whoami`

### Frontend Public
- [ ] Acessar diretório: `cd frontend-public`
- [ ] Linkar projeto: `vercel link`
- [ ] Nome sugerido: `ecosistema-imob-public`
- [ ] Verificar arquivo `.vercel/project.json` criado

### Frontend Admin
- [ ] Acessar diretório: `cd frontend-admin`
- [ ] Linkar projeto: `vercel link`
- [ ] Nome sugerido: `ecosistema-imob-admin`
- [ ] Verificar arquivo `.vercel/project.json` criado

### Obter Token
- [ ] Acessar [Vercel Tokens](https://vercel.com/account/tokens)
- [ ] Criar novo token
- [ ] Copiar token gerado
- [ ] Adicionar ao GitHub Secret: `VERCEL_TOKEN`

**Comando rápido:**
```bash
# Executar script automatizado
./scripts/setup-vercel.sh
```

---

## Fase 6: Primeiro Deploy (30 min)

### Backend DEV
- [ ] Criar branch: `git checkout -b test/first-deploy`
- [ ] Fazer alteração mínima no backend
- [ ] Commit: `git commit -m "test: primeiro deploy backend"`
- [ ] Push: `git push origin test/first-deploy`
- [ ] Criar PR: `test/first-deploy` → `develop`
- [ ] Aguardar CI passar
- [ ] Fazer merge
- [ ] Aguardar deploy completar no Actions
- [ ] Copiar URL do Cloud Run: `https://backend-api-xxxxx.run.app`
- [ ] Atualizar GitHub Secret: `NEXT_PUBLIC_API_URL_DEV`

### Frontend Public DEV
- [ ] Fazer alteração mínima no frontend-public
- [ ] Commit e push para `develop`
- [ ] Aguardar deploy completar
- [ ] Acessar URL do Vercel e testar

### Frontend Admin DEV
- [ ] Fazer alteração mínima no frontend-admin
- [ ] Commit e push para `develop`
- [ ] Aguardar deploy completar
- [ ] Acessar URL do Vercel e testar
- [ ] Testar login com Firebase

### Validação DEV
- [ ] Backend responde: `curl https://backend-api-xxxxx.run.app/health`
- [ ] Frontend Public carrega corretamente
- [ ] Frontend Admin carrega corretamente
- [ ] Frontend Admin conecta ao backend
- [ ] Autenticação Firebase funciona

---

## Fase 7: Deploy para Produção (20 min)

### Backend PROD
- [ ] Criar PR: `develop` → `main`
- [ ] Adicionar descrição detalhada no PR
- [ ] Aguardar aprovação
- [ ] Fazer merge
- [ ] Aguardar aprovação no environment `production`
- [ ] Aprovar deploy
- [ ] Aguardar deploy completar
- [ ] Copiar URL do Cloud Run: `https://backend-api-xxxxx.run.app`
- [ ] Atualizar GitHub Secret: `NEXT_PUBLIC_API_URL_PROD`

### Frontend PROD
- [ ] Criar PR: `develop` → `main` (se já não foi criado)
- [ ] Aguardar aprovação e merge
- [ ] Aguardar deploy completar
- [ ] Acessar URLs de produção e testar

### Validação PROD
- [ ] Backend PROD responde: `curl https://api.example.com/health`
- [ ] Frontend Public PROD carrega corretamente
- [ ] Frontend Admin PROD carrega corretamente
- [ ] Autenticação Firebase PROD funciona

---

## Fase 8: Configurações Finais (30 min)

### Domínios Personalizados

#### Cloud Run (Backend)
- [ ] Mapear domínio DEV: `api-dev.example.com`
- [ ] Mapear domínio PROD: `api.example.com`
- [ ] Adicionar registros DNS (CNAME)
- [ ] Aguardar propagação DNS
- [ ] Verificar HTTPS automático

#### Vercel (Frontend)
- [ ] Adicionar domínio Frontend Public DEV: `dev.example.com`
- [ ] Adicionar domínio Frontend Public PROD: `www.example.com`
- [ ] Adicionar domínio Frontend Admin DEV: `admin-dev.example.com`
- [ ] Adicionar domínio Frontend Admin PROD: `app.example.com`
- [ ] Adicionar registros DNS (CNAME)
- [ ] Aguardar propagação DNS
- [ ] Verificar SSL certificados

### Vercel Environment Variables

#### Frontend Public DEV
- [ ] `NEXT_PUBLIC_API_URL` (Preview)
- [ ] `NEXT_PUBLIC_ENVIRONMENT=development` (Preview)

#### Frontend Public PROD
- [ ] `NEXT_PUBLIC_API_URL` (Production)
- [ ] `NEXT_PUBLIC_ENVIRONMENT=production` (Production)

#### Frontend Admin DEV
- [ ] `NEXT_PUBLIC_API_URL` (Preview)
- [ ] `NEXT_PUBLIC_FIREBASE_API_KEY` (Preview)
- [ ] `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` (Preview)
- [ ] `NEXT_PUBLIC_FIREBASE_PROJECT_ID` (Preview)
- [ ] `NEXT_PUBLIC_ENVIRONMENT=development` (Preview)

#### Frontend Admin PROD
- [ ] `NEXT_PUBLIC_API_URL` (Production)
- [ ] `NEXT_PUBLIC_FIREBASE_API_KEY` (Production)
- [ ] `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` (Production)
- [ ] `NEXT_PUBLIC_FIREBASE_PROJECT_ID` (Production)
- [ ] `NEXT_PUBLIC_ENVIRONMENT=production` (Production)

---

## Fase 9: Monitoramento e Alertas (15 min)

### Cloud Monitoring
- [ ] Acessar Cloud Console > Monitoring
- [ ] Criar alerta: Latência > 1s
- [ ] Criar alerta: Taxa de erro > 5%
- [ ] Criar alerta: CPU > 80%
- [ ] Criar alerta: Memória > 80%
- [ ] Configurar notificações por email

### Logs
- [ ] Testar logs do backend DEV
- [ ] Testar logs do backend PROD
- [ ] Configurar retenção de logs (30 dias)

### Vercel Analytics
- [ ] Habilitar Analytics no Frontend Public
- [ ] Habilitar Analytics no Frontend Admin
- [ ] Configurar Core Web Vitals monitoring

---

## Fase 10: Documentação e Testes (10 min)

### Documentação
- [ ] Atualizar README.md com URLs dos ambientes
- [ ] Documentar processo de deploy
- [ ] Documentar processo de rollback
- [ ] Adicionar diagramas de arquitetura (opcional)

### Testes Finais
- [ ] Testar fluxo completo em DEV
- [ ] Testar fluxo completo em PROD
- [ ] Testar rollback do backend
- [ ] Testar rollback do frontend
- [ ] Verificar logs em ambos ambientes

### Backup
- [ ] Configurar export automático do Firestore DEV
- [ ] Configurar export automático do Firestore PROD
- [ ] Definir política de retenção (30 dias)
- [ ] Testar restauração de backup

---

## Checklist de Validação Final

### Funcionalidades
- [ ] Backend DEV responde corretamente
- [ ] Backend PROD responde corretamente
- [ ] Frontend Public DEV funciona
- [ ] Frontend Public PROD funciona
- [ ] Frontend Admin DEV funciona
- [ ] Frontend Admin PROD funciona
- [ ] Autenticação funciona em ambos ambientes
- [ ] Upload de imagens funciona
- [ ] APIs retornam dados corretamente

### CI/CD
- [ ] Deploy automático para DEV funciona
- [ ] Deploy para PROD requer aprovação
- [ ] Testes executam antes do deploy
- [ ] Rollback funciona
- [ ] Logs estão acessíveis
- [ ] Notificações de falha funcionam

### Segurança
- [ ] HTTPS em todos os endpoints
- [ ] Secrets não estão no código
- [ ] Service accounts têm permissões mínimas
- [ ] CORS configurado corretamente
- [ ] Firebase rules configuradas
- [ ] Rate limiting implementado (se aplicável)

### Performance
- [ ] Latência < 500ms (P95)
- [ ] Taxa de erro < 1%
- [ ] Core Web Vitals verdes
- [ ] Imagens otimizadas
- [ ] Bundle size otimizado

---

## Problemas Comuns e Soluções

### ❌ Deploy falha: "Permission denied"
**Solução:**
```bash
# Verificar permissões
gcloud projects get-iam-policy PROJECT_ID
# Adicionar roles necessárias
```

### ❌ Backend não inicia: "Firebase error"
**Solução:**
- Verificar formato do JSON do Firebase Admin SDK
- Verificar se secret está correto no GitHub
- Verificar se Firebase API está habilitada

### ❌ Frontend não conecta ao backend
**Solução:**
- Verificar CORS no backend
- Verificar URL da API nas env vars
- Verificar se Cloud Run permite tráfego não autenticado

### ❌ Vercel build falha
**Solução:**
- Verificar se `package-lock.json` está commitado
- Verificar node version
- Verificar env vars no Vercel Dashboard

---

## Próximos Passos Recomendados

### Curto Prazo (1-2 semanas)
- [ ] Implementar testes E2E no pipeline
- [ ] Configurar deploy canary para PROD
- [ ] Adicionar análise de qualidade de código
- [ ] Configurar alertas de Slack/Discord

### Médio Prazo (1 mês)
- [ ] Implementar feature flags
- [ ] Adicionar APM (Application Performance Monitoring)
- [ ] Configurar CDN para assets
- [ ] Implementar cache Redis

### Longo Prazo (3 meses)
- [ ] Multi-region deployment
- [ ] Disaster recovery testing
- [ ] Security audit completo
- [ ] Performance optimization avançada

---

## Tempo Total Estimado

| Fase | Tempo Estimado |
|------|----------------|
| Fase 1: Preparação | 15 min |
| Fase 2: GCP | 30 min |
| Fase 3: Firebase | 15 min |
| Fase 4: GitHub | 20 min |
| Fase 5: Vercel | 15 min |
| Fase 6: Primeiro Deploy | 30 min |
| Fase 7: Deploy PROD | 20 min |
| Fase 8: Config Finais | 30 min |
| Fase 9: Monitoramento | 15 min |
| Fase 10: Docs e Testes | 10 min |
| **TOTAL** | **3h 20min** |

---

## Referências

- [CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md) - Documentação completa
- [QUICK_START_CICD.md](QUICK_START_CICD.md) - Guia rápido
- [INFRASTRUCTURE_OVERVIEW.md](INFRASTRUCTURE_OVERVIEW.md) - Visão geral da infraestrutura
- [.github/workflows/README.md](.github/workflows/README.md) - Docs dos workflows

---

**Última atualização**: 2026-01-07
**Status**: ✅ Pronto para uso
