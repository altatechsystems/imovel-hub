# GitHub Actions Workflows

Este diretório contém os workflows de CI/CD para o projeto Ecossistema Imobiliário.

## Workflows Disponíveis

### 1. Backend Deploy (`backend-deploy.yml`)

Deploy automático do backend Go para Google Cloud Run.

**Triggers:**
- Push para `develop` → Deploy automático para DEV
- Push para `main` → Deploy para PROD (requer aprovação manual)
- Alterações em `backend/**`
- Dispatch manual

**Passos:**
1. Checkout do código
2. Setup Go 1.21
3. Executar testes
4. Autenticar no GCP
5. Criar config Firebase
6. Build e deploy no Cloud Run

**Ambiente DEV:**
- URL: `https://backend-api-xxxxx.run.app`
- Region: southamerica-east1
- GIN_MODE: debug

**Ambiente PROD:**
- URL: `https://api.example.com`
- Region: southamerica-east1
- GIN_MODE: release
- Requer aprovação manual

### 2. Frontend Public Deploy (`frontend-public-deploy.yml`)

Deploy do frontend público (Next.js) para Vercel.

**Triggers:**
- Push para `develop` → Deploy para Preview
- Push para `main` → Deploy para Production
- Alterações em `frontend-public/**`
- Dispatch manual

**Passos:**
1. Checkout do código
2. Setup Node.js 20
3. Instalar dependências
4. Build com Vercel CLI
5. Deploy

**Variáveis de Ambiente:**
- DEV: `NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_ENVIRONMENT=development`
- PROD: `NEXT_PUBLIC_API_URL`, `NEXT_PUBLIC_ENVIRONMENT=production`

### 3. Frontend Admin Deploy (`frontend-admin-deploy.yml`)

Deploy do frontend administrativo (Next.js) para Vercel.

**Triggers:**
- Push para `develop` → Deploy para Preview
- Push para `main` → Deploy para Production
- Alterações em `frontend-admin/**`
- Dispatch manual

**Passos:**
1. Checkout do código
2. Setup Node.js 20
3. Instalar dependências
4. Build com Vercel CLI (com Firebase config)
5. Deploy

**Variáveis de Ambiente:**
- DEV: Firebase config + API URL
- PROD: Firebase config + API URL

## Secrets Necessários

### GCP Secrets
```
GCP_PROJECT_ID_DEV
GCP_SA_KEY_DEV
GCP_PROJECT_ID_PROD
GCP_SA_KEY_PROD
```

### Firebase Secrets
```
FIREBASE_PROJECT_ID_DEV
FIREBASE_ADMIN_SDK_DEV
FIREBASE_API_KEY_DEV
FIREBASE_AUTH_DOMAIN_DEV

FIREBASE_PROJECT_ID_PROD
FIREBASE_ADMIN_SDK_PROD
FIREBASE_API_KEY_PROD
FIREBASE_AUTH_DOMAIN_PROD
```

### Storage Secrets
```
GCS_BUCKET_NAME_DEV
GCS_BUCKET_NAME_PROD
```

### Vercel Secrets
```
VERCEL_TOKEN
NEXT_PUBLIC_API_URL_DEV
NEXT_PUBLIC_API_URL_PROD
```

## Estratégia de Branches

```
develop (DEV)
    ↓
  feature/* → PR → develop (auto-deploy DEV)
    ↓
  develop → PR → main (requer aprovação)
    ↓
  main (PROD - auto-deploy após aprovação)
```

## Workflow de Deploy

### Para DEV (Desenvolvimento)

1. Criar feature branch:
   ```bash
   git checkout develop
   git checkout -b feature/minha-feature
   ```

2. Fazer alterações e commit:
   ```bash
   git add .
   git commit -m "feat: adicionar nova funcionalidade"
   git push origin feature/minha-feature
   ```

3. Criar Pull Request para `develop`

4. Após aprovação e merge → Deploy automático para DEV

### Para PROD (Produção)

1. Criar Pull Request: `develop` → `main`

2. Aguardar:
   - Revisão de código
   - Testes passarem
   - Aprovação manual

3. Merge → Deploy automático para PROD (após aprovação no environment)

## Ambientes GitHub

### Production
- Usado para deploys em `main`
- Requer aprovação manual de reviewers
- URL do ambiente configurada

Para configurar:
1. GitHub repo → Settings → Environments
2. New environment: `production`
3. Required reviewers: Adicionar usuários
4. Deployment branches: Only protected branches

## Monitoramento

### Ver Status dos Workflows
- GitHub repo → Actions tab
- Filtrar por workflow
- Ver logs detalhados

### Notificações
- Configurar em: Settings → Notifications
- Receber alertas de falhas

## Troubleshooting

### Deploy Falha: "Permission denied"
**Solução:** Verificar permissões da service account no GCP
```bash
gcloud projects get-iam-policy PROJECT_ID
```

### Build Falha: "Secret not found"
**Solução:** Verificar se todos os secrets estão configurados
- Settings → Secrets and variables → Actions

### Backend Não Inicia: "Firebase error"
**Solução:** Verificar formato do secret `FIREBASE_ADMIN_SDK_*`

### Vercel Deploy Falha
**Solução:**
1. Verificar token: `vercel whoami`
2. Verificar projeto existe: `vercel list`
3. Verificar variáveis de ambiente no Vercel Dashboard

## Comandos Úteis

### Executar Workflow Manualmente
1. GitHub repo → Actions
2. Selecionar workflow
3. Run workflow → Escolher branch

### Ver Logs de Deploy
```bash
# GitHub CLI
gh run list
gh run view RUN_ID --log

# Cloud Run logs
gcloud logging read "resource.type=cloud_run_revision" --limit 50
```

### Cancelar Deploy
1. GitHub repo → Actions
2. Encontrar workflow em andamento
3. Cancel workflow

## Boas Práticas

1. **Sempre testar em DEV primeiro**
   - Nunca fazer push direto para `main`
   - Sempre criar PR de `develop` para `main`

2. **Revisar logs após deploy**
   - Verificar se aplicação iniciou corretamente
   - Testar endpoints principais

3. **Rollback em caso de problemas**
   - Reverter commit problemático
   - Fazer novo deploy

4. **Manter secrets atualizados**
   - Rotacionar credenciais periodicamente
   - Atualizar secrets após renovação

5. **Documentar mudanças**
   - Commits semânticos
   - PR descriptions claras

## Próximas Melhorias

- [ ] Adicionar testes E2E no pipeline
- [ ] Implementar deploy canary para PROD
- [ ] Adicionar análise de qualidade de código (SonarQube)
- [ ] Configurar deploy preview para PRs
- [ ] Adicionar notificações Slack/Discord

## Referências

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Google Cloud Run Deploy](https://cloud.google.com/run/docs/deploying)
- [Vercel Deploy](https://vercel.com/docs/deployments)
- [Documentação Completa](../CONFIGURACAO_AMBIENTES_CICD.md)
