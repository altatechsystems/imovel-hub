# Visão Geral da Infraestrutura

Este documento fornece uma visão geral completa da infraestrutura de desenvolvimento e produção do Ecossistema Imobiliário.

## Diagrama de Arquitetura

```
┌─────────────────────────────────────────────────────────────────────┐
│                            GITHUB                                    │
│                                                                      │
│  ┌────────────┐     ┌────────────┐     ┌────────────┐             │
│  │  develop   │────▶│ Pull       │────▶│    main    │             │
│  │  (DEV)     │     │ Request    │     │  (PROD)    │             │
│  └────────────┘     └────────────┘     └────────────┘             │
│       │                                       │                     │
│       │ (auto)                               │ (manual approval)   │
│       ▼                                       ▼                     │
│  ┌────────────────────────────────────────────────────┐           │
│  │          GitHub Actions Workflows                  │           │
│  │  • backend-deploy.yml                             │           │
│  │  • frontend-public-deploy.yml                     │           │
│  │  • frontend-admin-deploy.yml                      │           │
│  └────────────────────────────────────────────────────┘           │
└─────────────────────────────────────────────────────────────────────┘
              │                                   │
              │                                   │
    ┌─────────▼────────┐               ┌─────────▼────────┐
    │   AMBIENTE DEV    │               │  AMBIENTE PROD    │
    │                   │               │                   │
    │ ┌───────────────┐ │               │ ┌───────────────┐ │
    │ │ Cloud Run     │ │               │ │ Cloud Run     │ │
    │ │ Backend (Go)  │ │               │ │ Backend (Go)  │ │
    │ │ Port: 8080    │ │               │ │ Port: 8080    │ │
    │ └───────────────┘ │               │ └───────────────┘ │
    │         │         │               │         │         │
    │         ▼         │               │         ▼         │
    │ ┌───────────────┐ │               │ ┌───────────────┐ │
    │ │  Firestore    │ │               │ │  Firestore    │ │
    │ │  (Database)   │ │               │ │  (Database)   │ │
    │ └───────────────┘ │               │ └───────────────┘ │
    │         │         │               │         │         │
    │         ▼         │               │         ▼         │
    │ ┌───────────────┐ │               │ ┌───────────────┐ │
    │ │ Cloud Storage │ │               │ │ Cloud Storage │ │
    │ │    (Images)   │ │               │ │    (Images)   │ │
    │ └───────────────┘ │               │ └───────────────┘ │
    │         │         │               │         │         │
    │         ▼         │               │         ▼         │
    │ ┌───────────────┐ │               │ ┌───────────────┐ │
    │ │   Firebase    │ │               │ │   Firebase    │ │
    │ │     Auth      │ │               │ │     Auth      │ │
    │ └───────────────┘ │               │ └───────────────┘ │
    └───────────────────┘               └───────────────────┘
              │                                   │
              │                                   │
    ┌─────────▼────────┐               ┌─────────▼────────┐
    │      VERCEL      │               │      VERCEL      │
    │                  │               │                  │
    │ ┌──────────────┐ │               │ ┌──────────────┐ │
    │ │Frontend      │ │               │ │Frontend      │ │
    │ │Public (Next) │ │               │ │Public (Next) │ │
    │ │dev.example   │ │               │ │www.example   │ │
    │ └──────────────┘ │               │ └──────────────┘ │
    │                  │               │                  │
    │ ┌──────────────┐ │               │ ┌──────────────┐ │
    │ │Frontend      │ │               │ │Frontend      │ │
    │ │Admin (Next)  │ │               │ │Admin (Next)  │ │
    │ │admin-dev.ex  │ │               │ │app.example   │ │
    │ └──────────────┘ │               │ └──────────────┘ │
    └───────────────────┘               └───────────────────┘
              │                                   │
              │                                   │
              ▼                                   ▼
    ┌───────────────────┐               ┌───────────────────┐
    │   Usuários DEV    │               │  Usuários PROD    │
    └───────────────────┘               └───────────────────┘
```

## Componentes por Ambiente

### Ambiente DEV (Desenvolvimento)

| Componente | Serviço | URL | Observações |
|------------|---------|-----|-------------|
| **Backend** | Cloud Run | `backend-api-xxxxx.run.app` | Deploy automático |
| **Database** | Firestore | `ecosistema-imob-dev` | Dados de teste |
| **Storage** | GCS | `ecosistema-imob-dev-storage` | Imagens de teste |
| **Auth** | Firebase | `ecosistema-imob-dev` | Usuários de teste |
| **Frontend Public** | Vercel | `*.vercel.app` ou `dev.example.com` | Preview deployments |
| **Frontend Admin** | Vercel | `*.vercel.app` ou `admin-dev.example.com` | Preview deployments |

**Características DEV:**
- Deploy automático ao fazer push para `develop`
- Sem necessidade de aprovação manual
- GIN_MODE=debug (logs detalhados)
- Dados podem ser resetados
- Ideal para testes rápidos

### Ambiente PROD (Produção)

| Componente | Serviço | URL | Observações |
|------------|---------|-----|-------------|
| **Backend** | Cloud Run | `api.example.com` | Deploy com aprovação |
| **Database** | Firestore | `ecosistema-imob-prod` | Dados reais |
| **Storage** | GCS | `ecosistema-imob-prod-storage` | Imagens reais |
| **Auth** | Firebase | `ecosistema-imob-prod` | Usuários reais |
| **Frontend Public** | Vercel | `www.example.com` | Production |
| **Frontend Admin** | Vercel | `app.example.com` | Production |

**Características PROD:**
- Deploy requer aprovação manual (GitHub Environment)
- GIN_MODE=release (otimizado)
- Backup automático (configurar)
- Monitoramento ativo
- Alta disponibilidade

## Fluxo de Deploy

### 1. Desenvolvimento (Feature → DEV)

```bash
# Desenvolvedor cria feature branch
git checkout develop
git checkout -b feature/nova-funcionalidade

# Faz alterações e commit
git add .
git commit -m "feat: adicionar nova funcionalidade"
git push origin feature/nova-funcionalidade

# Cria Pull Request: feature/* → develop
# GitHub Actions executa:
#  ✓ Testes automatizados
#  ✓ Lint
#  ✓ Build

# Após merge em develop:
#  → Deploy automático para DEV
#  → Testes de integração (opcional)
```

### 2. Homologação → Produção (DEV → PROD)

```bash
# Após validação em DEV, criar PR para PROD
git checkout develop
git pull origin develop

# Criar Pull Request: develop → main
# GitHub Actions executa:
#  ✓ Testes automatizados
#  ✓ Security scan
#  ✓ Build

# Aguarda aprovação manual (reviewer)
# Após aprovação e merge:
#  → Deploy automático para PROD (após aprovação no environment)
```

## Configuração de Secrets

### GitHub Secrets (Repository Level)

```
# GCP
GCP_PROJECT_ID_DEV=ecosistema-imob-dev
GCP_SA_KEY_DEV=<JSON da service account DEV>
GCP_PROJECT_ID_PROD=ecosistema-imob-prod
GCP_SA_KEY_PROD=<JSON da service account PROD>

# Firebase
FIREBASE_PROJECT_ID_DEV=ecosistema-imob-dev
FIREBASE_ADMIN_SDK_DEV=<JSON do Firebase Admin SDK DEV>
FIREBASE_API_KEY_DEV=<API Key do Firebase Console DEV>
FIREBASE_AUTH_DOMAIN_DEV=ecosistema-imob-dev.firebaseapp.com

FIREBASE_PROJECT_ID_PROD=ecosistema-imob-prod
FIREBASE_ADMIN_SDK_PROD=<JSON do Firebase Admin SDK PROD>
FIREBASE_API_KEY_PROD=<API Key do Firebase Console PROD>
FIREBASE_AUTH_DOMAIN_PROD=ecosistema-imob-prod.firebaseapp.com

# Storage
GCS_BUCKET_NAME_DEV=ecosistema-imob-dev-storage
GCS_BUCKET_NAME_PROD=ecosistema-imob-prod-storage

# Vercel
VERCEL_TOKEN=<Token do Vercel>
NEXT_PUBLIC_API_URL_DEV=<URL do Cloud Run DEV>
NEXT_PUBLIC_API_URL_PROD=<URL do Cloud Run PROD>
```

### Vercel Environment Variables

**Frontend Public (DEV - Preview):**
```
NEXT_PUBLIC_API_URL=https://backend-api-dev.run.app
NEXT_PUBLIC_ENVIRONMENT=development
```

**Frontend Public (PROD - Production):**
```
NEXT_PUBLIC_API_URL=https://api.example.com
NEXT_PUBLIC_ENVIRONMENT=production
```

**Frontend Admin (DEV - Preview):**
```
NEXT_PUBLIC_API_URL=https://backend-api-dev.run.app
NEXT_PUBLIC_FIREBASE_API_KEY=<API Key DEV>
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=ecosistema-imob-dev.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=ecosistema-imob-dev
NEXT_PUBLIC_ENVIRONMENT=development
```

**Frontend Admin (PROD - Production):**
```
NEXT_PUBLIC_API_URL=https://api.example.com
NEXT_PUBLIC_FIREBASE_API_KEY=<API Key PROD>
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=ecosistema-imob-prod.firebaseapp.com
NEXT_PUBLIC_FIREBASE_PROJECT_ID=ecosistema-imob-prod
NEXT_PUBLIC_ENVIRONMENT=production
```

## Custos Mensais Estimados

### Ambiente DEV

| Serviço | Custo Mensal | Observações |
|---------|--------------|-------------|
| Cloud Run | $5-10 | Baixo tráfego |
| Firestore | $5-10 | Poucos documentos |
| Cloud Storage | $2-5 | Poucas imagens |
| Firebase Auth | $0 | Free tier |
| Vercel | $0 | Hobby plan |
| **Total DEV** | **$12-25** | |

### Ambiente PROD

| Serviço | Custo Mensal | Observações |
|---------|--------------|-------------|
| Cloud Run | $20-50 | Tráfego médio |
| Firestore | $50-200 | Volume real |
| Cloud Storage | $10-30 | Imagens reais |
| Firebase Auth | $0-50 | Depende do volume |
| Vercel Pro | $20 | Production |
| Cloud CDN | $10-20 | Cache de assets |
| **Total PROD** | **$110-370** | |

**Total Geral: $122-395/mês**

## Monitoramento e Observabilidade

### Cloud Run (Backend)

```bash
# Logs em tempo real
gcloud logging tail "resource.type=cloud_run_revision" \
  --project=ecosistema-imob-prod

# Métricas
# Acesse: Cloud Console > Cloud Run > backend-api > Metrics
```

**Alertas Recomendados:**
- Latência P95 > 1s
- Taxa de erro > 5%
- CPU > 80%
- Memória > 80%

### Vercel (Frontend)

**Analytics:**
- Core Web Vitals
- Visitor analytics
- Performance monitoring

**Logs:**
- Runtime logs
- Build logs
- Edge network logs

### Firestore (Database)

**Métricas:**
- Document reads/writes
- Storage usage
- Query performance

**Backup:**
- Configurar export automático diário
- Retenção: 30 dias

## Segurança

### Backend (Cloud Run)

- ✅ HTTPS obrigatório
- ✅ Service account com least privilege
- ✅ Secrets via environment variables
- ✅ CORS configurado
- ✅ Rate limiting (implementar)
- ✅ JWT validation (Firebase)

### Frontend (Vercel)

- ✅ HTTPS obrigatório
- ✅ Environment variables separadas
- ✅ No secrets in client code
- ✅ CSP headers (configurar)
- ✅ CORS headers

### Database (Firestore)

- ✅ Security rules
- ✅ Multi-tenancy isolation
- ✅ Audit logging
- ✅ Backup automático

### Secrets Management

- ✅ GitHub Secrets encrypted
- ✅ Service account keys rotacionadas
- ✅ Firebase credentials protegidas
- ✅ No secrets em código

## Disaster Recovery

### Backup Strategy

**Firestore:**
- Export automático: Diário (2 AM)
- Retenção: 30 dias
- Location: GCS bucket separado

**Cloud Storage:**
- Versioning: Habilitado
- Lifecycle policy: 90 dias
- Backup: Mirror em outro bucket

### Rollback Procedures

**Backend (Cloud Run):**
```bash
# Listar revisões
gcloud run revisions list --service backend-api

# Rollback
gcloud run services update-traffic backend-api \
  --to-revisions PREVIOUS-REVISION=100
```

**Frontend (Vercel):**
```bash
# Listar deployments
vercel list

# Rollback (via Vercel Dashboard ou CLI)
vercel rollback DEPLOYMENT-URL
```

**Database (Firestore):**
```bash
# Restaurar de backup
gcloud firestore import gs://backup-bucket/[TIMESTAMP]
```

## Performance Optimization

### Backend

- [x] Dockerfile multi-stage build
- [ ] Redis caching (próxima fase)
- [ ] Connection pooling
- [ ] Query optimization

### Frontend

- [x] Next.js 14 App Router
- [x] SSR/SSG para páginas públicas
- [ ] Image optimization
- [ ] Bundle size analysis
- [ ] CDN para assets estáticos

### Database

- [x] Composite indexes
- [ ] Query caching
- [ ] Connection pooling
- [ ] Read replicas (se necessário)

## Próximos Passos

### Curto Prazo (1-2 semanas)
- [ ] Configurar domínios personalizados
- [ ] Implementar health checks
- [ ] Configurar alertas
- [ ] Documentar runbooks

### Médio Prazo (1 mês)
- [ ] Implementar testes E2E no CI
- [ ] Configurar CDN
- [ ] Adicionar APM (DataDog/New Relic)
- [ ] Implementar feature flags

### Longo Prazo (3 meses)
- [ ] Multi-region deployment
- [ ] Disaster recovery testing
- [ ] Performance optimization
- [ ] Security audit

## Referências

- [CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md) - Guia completo de configuração
- [QUICK_START_CICD.md](QUICK_START_CICD.md) - Guia rápido de início
- [.github/workflows/README.md](.github/workflows/README.md) - Documentação dos workflows
- [README.md](README.md) - Visão geral do projeto

---

**Versão**: 1.0
**Data**: 2026-01-07
**Última Atualização**: Configuração inicial de CI/CD
