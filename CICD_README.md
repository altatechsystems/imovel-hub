# CI/CD - Guia Completo

Bem-vindo ao guia completo de CI/CD do Ecossistema Imobiliário! Este documento é seu ponto de partida para configurar e gerenciar os ambientes de desenvolvimento e produção.

## 📚 Documentação Disponível

A documentação está organizada em diferentes níveis de detalhe. Escolha o que melhor atende sua necessidade:

### 🚀 Para Começar Rapidamente

| Documento | Descrição | Tempo | Quando Usar |
|-----------|-----------|-------|-------------|
| **[CONFIGURACAO_INTERFACE_WEB.md](CONFIGURACAO_INTERFACE_WEB.md)** ⭐ | Configuração via interface web (SEM CLI) | 4-5 horas | Prefere usar interfaces gráficas |
| **[QUICK_START_CICD.md](QUICK_START_CICD.md)** | Guia rápido de configuração (com CLI) | 45-60 min | Primeira configuração, setup rápido |
| **[CHECKLIST_CONFIGURACAO_CICD.md](CHECKLIST_CONFIGURACAO_CICD.md)** | Checklist passo a passo | 3-4 horas | Acompanhar progresso da configuração |

### 📖 Para Entender em Profundidade

| Documento | Descrição | Quando Usar |
|-----------|-----------|-------------|
| **[CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md)** | Documentação completa e detalhada | Referência completa, troubleshooting avançado |
| **[INFRASTRUCTURE_OVERVIEW.md](INFRASTRUCTURE_OVERVIEW.md)** | Visão geral da arquitetura | Entender arquitetura, diagramas, custos |

### 🔧 Para Uso Diário

| Documento | Descrição | Quando Usar |
|-----------|-----------|-------------|
| **[COMANDOS_UTEIS_CICD.md](COMANDOS_UTEIS_CICD.md)** | Referência rápida de comandos | Operação diária, troubleshooting |
| **[.github/workflows/README.md](.github/workflows/README.md)** | Documentação dos workflows | Entender/modificar workflows |

---

## 🎯 Começe por Aqui

### Primeira Vez Configurando?

**Opção A: Via Interface Web (Recomendado para iniciantes) 🖱️**
1. ✅ Leia [CONFIGURACAO_INTERFACE_WEB.md](CONFIGURACAO_INTERFACE_WEB.md)
2. ✅ Siga o guia passo a passo (tudo via web)
3. ✅ Faça o primeiro deploy

**Opção B: Via Linha de Comando (Mais rápido) ⚡**
1. ✅ Leia [QUICK_START_CICD.md](QUICK_START_CICD.md) (10 min)
2. ✅ Execute os scripts de setup:
   ```bash
   ./scripts/setup-gcp-environments.sh
   ./scripts/setup-vercel.sh
   ```
3. ✅ Siga o [CHECKLIST_CONFIGURACAO_CICD.md](CHECKLIST_CONFIGURACAO_CICD.md)
4. ✅ Faça o primeiro deploy

### Já Configurado?

- 📱 Use [COMANDOS_UTEIS_CICD.md](COMANDOS_UTEIS_CICD.md) para operação diária
- 🔍 Consulte [INFRASTRUCTURE_OVERVIEW.md](INFRASTRUCTURE_OVERVIEW.md) para visão geral
- 🐛 Problemas? Veja seção de Troubleshooting em qualquer documento

---

## 🏗️ Arquitetura em Resumo

```
GitHub (develop/main)
         ↓
   GitHub Actions
         ↓
    ┌────┴────┐
    ↓         ↓
  DEV       PROD
    ↓         ↓
Cloud Run   Cloud Run
Firestore   Firestore
   GCS        GCS
Vercel      Vercel
```

**2 Ambientes Isolados:**
- **DEV**: Deploy automático ao fazer push para `develop`
- **PROD**: Deploy após aprovação manual ao fazer merge em `main`

---

## 🚢 Fluxo de Deploy Simplificado

### Para Desenvolvimento (DEV)

```bash
# 1. Criar feature branch
git checkout develop
git checkout -b feature/minha-feature

# 2. Fazer alterações
# ... código ...

# 3. Commit e push
git add .
git commit -m "feat: nova funcionalidade"
git push origin feature/minha-feature

# 4. Criar PR no GitHub: feature/minha-feature → develop
# 5. Após merge → Deploy automático para DEV ✨
```

### Para Produção (PROD)

```bash
# 1. Após validar em DEV, criar PR: develop → main
# 2. Aguardar aprovação de reviewer
# 3. Fazer merge
# 4. Aguardar aprovação no GitHub Environment
# 5. Deploy automático para PROD ✨
```

---

## 📦 O Que Foi Configurado

### ✅ Arquivos Criados

```
.github/workflows/
├── backend-deploy.yml          # Deploy do backend Go
├── frontend-public-deploy.yml  # Deploy do frontend público
├── frontend-admin-deploy.yml   # Deploy do frontend admin
└── README.md                   # Docs dos workflows

backend/
├── Dockerfile                  # Container do backend
└── .dockerignore              # Arquivos ignorados no build

scripts/
├── setup-gcp-environments.sh  # Setup automático do GCP
└── setup-vercel.sh            # Setup automático do Vercel

docs/ (CI/CD)
├── CONFIGURACAO_AMBIENTES_CICD.md   # Documentação completa
├── QUICK_START_CICD.md              # Guia rápido
├── CHECKLIST_CONFIGURACAO_CICD.md   # Checklist passo a passo
├── INFRASTRUCTURE_OVERVIEW.md       # Visão geral da infra
├── COMANDOS_UTEIS_CICD.md          # Comandos úteis
└── CICD_README.md                   # Este arquivo
```

### ✅ Infraestrutura Configurada

**Google Cloud Platform (2 projetos):**
- `ecosistema-imob-dev` (desenvolvimento)
- `ecosistema-imob-prod` (produção)

**Serviços por Projeto:**
- Cloud Run (backend)
- Firestore (database)
- Cloud Storage (imagens)
- Firebase Authentication

**Vercel (4 projetos):**
- Frontend Public DEV
- Frontend Public PROD
- Frontend Admin DEV
- Frontend Admin PROD

**GitHub:**
- Workflows configurados
- Secrets adicionados
- Branch protection
- Environment `production`

---

## 🔐 Segurança

### Secrets no GitHub

Todos os secrets sensíveis estão configurados em:
`Settings > Secrets and variables > Actions`

**NUNCA commite:**
- ❌ Chaves de service accounts (`.json`)
- ❌ Credenciais do Firebase
- ❌ Tokens do Vercel
- ❌ API keys

### Arquivos Protegidos pelo .gitignore

```
.gcp-keys/              # Chaves GCP
*-key.json              # Service account keys
firebase-adminsdk*.json # Credenciais Firebase
.vercel/                # Config Vercel local
.env*                   # Environment variables
```

---

## 💰 Custos Estimados

| Ambiente | Custo Mensal | Detalhes |
|----------|--------------|----------|
| **DEV** | $12-25 | Baixo tráfego, poucos dados |
| **PROD** | $110-370 | Tráfego real, dados reais |
| **TOTAL** | **$122-395** | Pode variar com uso |

**Dicas para Economizar:**
- Use DEV apenas quando necessário
- Delete recursos não utilizados
- Configure budgets no GCP
- Monitore uso regularmente

---

## 📊 Monitoramento

### URLs de Monitoramento

**GCP (Cloud Run):**
- [Console DEV](https://console.cloud.google.com/run?project=ecosistema-imob-dev)
- [Console PROD](https://console.cloud.google.com/run?project=ecosistema-imob-prod)

**Vercel:**
- [Dashboard](https://vercel.com/dashboard)

**GitHub Actions:**
- [Workflows](https://github.com/seu-usuario/ecosistema-imob/actions)

### Comandos de Monitoramento

```bash
# Ver logs do backend
gcloud logging tail "resource.type=cloud_run_revision" --project ecosistema-imob-prod

# Ver status do Cloud Run
gcloud run services describe backend-api --region southamerica-east1 --project ecosistema-imob-prod

# Ver deployments do Vercel
vercel list

# Ver workflows do GitHub
gh run list
```

---

## 🐛 Troubleshooting Rápido

### Problema: Deploy falha no GitHub Actions

1. Verificar logs do workflow no GitHub
2. Verificar se todos os secrets estão configurados
3. Verificar permissões das service accounts

**Solução rápida:**
```bash
# Verificar permissões
gcloud projects get-iam-policy PROJECT_ID
```

### Problema: Backend não inicia

1. Verificar logs do Cloud Run
2. Verificar variáveis de ambiente
3. Verificar credenciais do Firebase

**Solução rápida:**
```bash
# Ver logs
gcloud logging read "resource.type=cloud_run_revision AND severity>=ERROR" --limit 50
```

### Problema: Frontend não conecta ao backend

1. Verificar CORS no backend
2. Verificar URL da API nas env vars do Vercel
3. Verificar se Cloud Run permite tráfego não autenticado

**Solução rápida:**
```bash
# Testar endpoint
curl -v https://backend-api-xxxxx.run.app/health
```

### Mais Troubleshooting

Consulte a seção "Troubleshooting" em:
- [CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md)
- [COMANDOS_UTEIS_CICD.md](COMANDOS_UTEIS_CICD.md)

---

## 🎓 Melhores Práticas

### Deploy

✅ **Faça:**
- Sempre testar em DEV primeiro
- Fazer PRs descritivos
- Aguardar CI passar antes de merge
- Revisar logs após deploy

❌ **Evite:**
- Push direto para `main`
- Skip de aprovações
- Deploy em horário de pico
- Múltiplos deploys simultâneos

### Código

✅ **Faça:**
- Commits semânticos (feat, fix, docs, etc)
- Mensagens de commit descritivas
- Testes antes de push
- Code review

❌ **Evite:**
- Commitar secrets
- Commitar arquivos de build
- Commits muito grandes
- Código não testado

### Segurança

✅ **Faça:**
- Rotacionar secrets periodicamente
- Usar least privilege
- Revisar permissões regularmente
- Manter dependências atualizadas

❌ **Evite:**
- Compartilhar secrets
- Usar mesmas credenciais em dev/prod
- Ignorar alertas de segurança
- Acesso root desnecessário

---

## 📞 Suporte

### Documentação do Projeto

- [README.md](README.md) - Visão geral do projeto
- [AI_DEV_DIRECTIVE.md](AI_DEV_DIRECTIVE.md) - Diretrizes de desenvolvimento

### Documentação Externa

- [GCP Documentation](https://cloud.google.com/docs)
- [Firebase Documentation](https://firebase.google.com/docs)
- [Vercel Documentation](https://vercel.com/docs)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

### Contato

Para problemas específicos:
1. Consulte a documentação relevante
2. Verifique seção de troubleshooting
3. Consulte logs do serviço
4. Entre em contato com o time DevOps

---

## 🗺️ Roadmap

### ✅ Concluído (Fase 1)

- [x] Setup de ambientes DEV e PROD
- [x] CI/CD com GitHub Actions
- [x] Deploy automático para DEV
- [x] Deploy com aprovação para PROD
- [x] Documentação completa

### 🚧 Próximos Passos (Fase 2)

- [ ] Testes E2E no pipeline
- [ ] Análise de qualidade de código
- [ ] Deploy canary para PROD
- [ ] Configurar CDN

### 🔮 Futuro (Fase 3)

- [ ] Multi-region deployment
- [ ] APM (Application Performance Monitoring)
- [ ] Feature flags
- [ ] Disaster recovery testing

---

## ✅ Checklist Rápido

Já configurou tudo? Verifique:

**Infraestrutura:**
- [ ] GCP DEV configurado
- [ ] GCP PROD configurado
- [ ] Firestore DEV criado
- [ ] Firestore PROD criado
- [ ] Buckets GCS criados
- [ ] Service accounts criadas

**GitHub:**
- [ ] Secrets configurados
- [ ] Workflows commitados
- [ ] Branch protection ativo
- [ ] Environment `production` criado

**Vercel:**
- [ ] Projetos linkados
- [ ] Env vars configuradas
- [ ] Domínios configurados (opcional)

**Deploy:**
- [ ] Backend DEV deployado
- [ ] Backend PROD deployado
- [ ] Frontend Public deployado
- [ ] Frontend Admin deployado

**Validação:**
- [ ] Backend responde
- [ ] Frontend carrega
- [ ] Autenticação funciona
- [ ] CI/CD funciona

---

## 📝 Notas de Versão

### v1.0 (2026-01-07)
- Configuração inicial de CI/CD
- 3 workflows (backend, frontend-public, frontend-admin)
- 2 ambientes (DEV, PROD)
- Documentação completa
- Scripts de setup automatizado

---

**Status**: ✅ Pronto para uso

**Última atualização**: 2026-01-07

**Mantenedores**: Altatech Systems

---

## 🎉 Você está pronto!

Agora você tem tudo configurado para desenvolver com confiança. Consulte a documentação conforme necessário e bom desenvolvimento! 🚀
