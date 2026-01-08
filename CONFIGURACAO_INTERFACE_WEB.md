# Configuração CI/CD via Interface Web

Guia completo para configurar os ambientes DEV e PROD usando apenas interfaces web, **SEM usar linha de comando**.

---

## 📋 Índice

1. [Google Cloud Platform (GCP)](#1-google-cloud-platform-gcp)
2. [Firebase](#2-firebase)
3. [GitHub](#3-github)
4. [Vercel](#4-vercel)
5. [Validação Final](#5-validação-final)

**Tempo estimado total:** 3-4 horas

---

## 1. Google Cloud Platform (GCP)

### 1.1 Criar Projeto de Produção (5 min)

1. **Acesse:** https://console.cloud.google.com
2. **Faça login** com sua conta Google
3. No topo da página, clique no **seletor de projetos** (ao lado do logo do Google Cloud)
4. Clique em **"NEW PROJECT"** (canto superior direito)
5. Preencha:
   - **Project name:** `Ecosistema Imob - Produção`
   - **Project ID:** `ecosistema-imob-prod` (anote este ID!)
   - **Location:** Deixe como está (Organization)
6. Clique em **"CREATE"**
7. Aguarde a criação (leva ~30 segundos)
8. **Selecione o novo projeto** no seletor de projetos

### 1.2 Habilitar APIs Necessárias (10 min)

Para cada API abaixo, siga este processo:

1. **Acesse:** https://console.cloud.google.com/apis/library
2. **Certifique-se** que o projeto correto está selecionado no topo
3. Use a **barra de busca** para encontrar a API
4. Clique na API
5. Clique em **"ENABLE"** (Ativar)
6. Aguarde ativação

**APIs para habilitar (DEV e PROD):**

**Para PROD (ecosistema-imob-prod):**
- ✅ Cloud Run API
- ✅ Cloud Firestore API
- ✅ Cloud Storage API
- ✅ Cloud Build API
- ✅ Artifact Registry API
- ✅ Identity and Access Management (IAM) API

**Para DEV (ecosistema-imob-dev):**
- Repita o processo acima, mas selecione o projeto `ecosistema-imob-dev`
- Habilite as mesmas APIs

**Atalho direto:** https://console.cloud.google.com/flows/enableapi?apiid=run.googleapis.com,firestore.googleapis.com,storage.googleapis.com,cloudbuild.googleapis.com,artifactregistry.googleapis.com,iam.googleapis.com

### 1.3 Criar Banco Firestore - PROD (5 min)

1. **Acesse:** https://console.cloud.google.com/firestore
2. Selecione projeto: **ecosistema-imob-prod**
3. Clique em **"SELECT NATIVE MODE"** (modo nativo)
4. Configurações:
   - **Database ID:** `(default)`
   - **Location type:** `Region`
   - **Location:** `southamerica-east1 (São Paulo)`
5. Clique em **"CREATE DATABASE"**
6. Aguarde criação (~2 minutos)

**Para DEV:** Verifique se já existe. Se não:
- Repita o processo acima com projeto `ecosistema-imob-dev`

### 1.4 Criar Buckets do Cloud Storage (10 min)

#### Bucket DEV

1. **Acesse:** https://console.cloud.google.com/storage/browser
2. Selecione projeto: **ecosistema-imob-dev**
3. Clique em **"CREATE BUCKET"** (+ CREATE no topo)
4. Preencha:
   - **Name:** `ecosistema-imob-dev-storage`
   - **Location type:** `Region`
   - **Location:** `southamerica-east1 (São Paulo)`
   - **Storage class:** `Standard`
   - **Access control:** `Fine-grained` (deixar marcado)
   - **Protection tools:** Desmarque "Enforce public access prevention"
5. Clique em **"CREATE"**
6. Após criar, clique no bucket criado
7. Vá na aba **"PERMISSIONS"**
8. Clique em **"GRANT ACCESS"**
9. Adicionar principal:
   - **New principals:** `allUsers`
   - **Role:** `Storage Object Viewer`
10. Clique em **"SAVE"**
11. Confirme o aviso de acesso público

#### Bucket PROD

1. Selecione projeto: **ecosistema-imob-prod**
2. Repita o processo acima com nome: `ecosistema-imob-prod-storage`

### 1.5 Criar Service Accounts para CI/CD (15 min)

#### Service Account DEV

1. **Acesse:** https://console.cloud.google.com/iam-admin/serviceaccounts
2. Selecione projeto: **ecosistema-imob-dev**
3. Clique em **"+ CREATE SERVICE ACCOUNT"** (topo)
4. Preencha:
   - **Service account name:** `github-actions-dev`
   - **Service account ID:** `github-actions-dev` (gerado automaticamente)
   - **Description:** `Service account for GitHub Actions - Development`
5. Clique em **"CREATE AND CONTINUE"**
6. **Grant this service account access to project:**
   - Clique em **"Select a role"**
   - Digite `Cloud Run Admin` e selecione **"Cloud Run Admin"**
   - Clique em **"+ ADD ANOTHER ROLE"**
   - Digite `Service Account User` e selecione **"Service Account User"**
   - Clique em **"+ ADD ANOTHER ROLE"**
   - Digite `Storage Admin` e selecione **"Storage Admin"**
7. Clique em **"CONTINUE"**
8. Clique em **"DONE"**

#### Service Account PROD

1. Selecione projeto: **ecosistema-imob-prod**
2. Repita o processo acima:
   - Nome: `github-actions-prod`
   - Mesmas 3 roles (Cloud Run Admin, Service Account User, Storage Admin)

### 1.6 Gerar Chaves das Service Accounts (10 min)

#### Chave DEV

1. **Acesse:** https://console.cloud.google.com/iam-admin/serviceaccounts
2. Selecione projeto: **ecosistema-imob-dev**
3. Encontre a service account: `github-actions-dev@ecosistema-imob-dev.iam.gserviceaccount.com`
4. Clique nos **3 pontos** (⋮) no final da linha
5. Clique em **"Manage keys"**
6. Clique em **"ADD KEY"** → **"Create new key"**
7. Selecione **"JSON"**
8. Clique em **"CREATE"**
9. O arquivo JSON será baixado automaticamente
10. **Renomeie** o arquivo para: `github-actions-dev-key.json`
11. **Guarde em local seguro** (você vai precisar depois)

#### Chave PROD

1. Selecione projeto: **ecosistema-imob-prod**
2. Repita o processo acima:
   - Service account: `github-actions-prod@ecosistema-imob-prod.iam.gserviceaccount.com`
   - Renomeie para: `github-actions-prod-key.json`

**⚠️ IMPORTANTE:**
- NÃO faça commit desses arquivos no Git
- Guarde-os em local seguro
- Você vai colar o conteúdo no GitHub Secrets

---

## 2. Firebase

### 2.1 Adicionar Projeto PROD ao Firebase (10 min)

1. **Acesse:** https://console.firebase.google.com
2. Clique em **"Add project"** (Adicionar projeto)
3. Selecione: **"Use an existing Google Cloud project"**
4. Escolha: **ecosistema-imob-prod**
5. Clique em **"Continue"**
6. **Confirme o plano:** Mantenha "Spark plan (free)"
7. Desmarque Google Analytics (opcional, pode deixar marcado)
8. Clique em **"Add Firebase"**
9. Aguarde configuração (~1 minuto)

**Para DEV:** Verifique se `ecosistema-imob-dev` já está no Firebase. Se não, repita o processo acima.

### 2.2 Ativar Firebase Authentication (5 min)

#### PROD

1. No Firebase Console, selecione projeto: **ecosistema-imob-prod**
2. No menu lateral esquerdo, procure pela seção **"Criação"** (ou "Build")
3. Clique em **"Authentication"**
   - Se não aparecer no menu, role para baixo ou expanda a seção "Criação"
4. Clique em **"Get started"** (botão grande no centro)
5. Você será levado para a tela de métodos de login
6. Clique em **"Email/Password"** (primeira opção na lista)
7. Na janela que abrir:
   - **Ative** o primeiro toggle: "Email/Password"
   - Deixe **desativado** o segundo toggle: "Email link (passwordless sign-in)"
8. Clique em **"Save"** (canto inferior direito da janela)
9. Pronto! Você verá "Email/Password" como "Enabled" na lista

#### DEV

1. Selecione projeto: **ecosistema-imob-dev**
2. Repita os passos 2-9 acima

### 2.3 Baixar Credenciais Firebase Admin SDK (10 min)

#### PROD

1. No Firebase Console, selecione projeto: **ecosistema-imob-prod**
2. Clique no **ícone de engrenagem** (⚙️) ao lado de "Project Overview"
3. Clique em **"Project settings"**
4. Vá na aba **"Service accounts"**
5. Clique em **"Generate new private key"**
6. Confirme clicando em **"Generate key"**
7. O arquivo JSON será baixado
8. **Renomeie** para: `firebase-adminsdk-prod.json`
9. **Guarde em local seguro**

#### DEV

1. Repita o processo para `ecosistema-imob-dev`
2. Renomeie para: `firebase-adminsdk-dev.json`

### 2.4 Obter API Keys do Firebase (5 min)

#### PROD

1. No Firebase Console, projeto: **ecosistema-imob-prod**
2. **Project Settings** (engrenagem)
3. Na aba **"General"**, role para baixo até **"Your apps"**
4. Se não há app web, clique em **"</>** (ícone web)
5. **App nickname:** `Admin Web App`
6. Clique em **"Register app"**
7. Copie as informações do **firebaseConfig**:
   ```javascript
   apiKey: "AIzaSy..." // ← Copie este valor
   authDomain: "ecosistema-imob-prod.firebaseapp.com"
   projectId: "ecosistema-imob-prod"
   ```
8. **Anote esses valores:**
   - `FIREBASE_API_KEY_PROD`: (o valor de apiKey)
   - `FIREBASE_AUTH_DOMAIN_PROD`: ecosistema-imob-prod.firebaseapp.com
   - `FIREBASE_PROJECT_ID_PROD`: ecosistema-imob-prod

#### DEV

1. Repita o processo para `ecosistema-imob-dev`
2. Anote os valores DEV

### 2.5 Configurar Domínios Autorizados (5 min)

#### PROD

1. Firebase Console → **ecosistema-imob-prod**
2. **Authentication** → **Settings** → aba **"Authorized domains"**
3. Clique em **"Add domain"**
4. Adicione (um por vez):
   - `api.example.com` (substitua pelo seu domínio real)
   - `www.example.com`
   - `app.example.com`
5. Clique em **"Add"** para cada

#### DEV

1. Repita para `ecosistema-imob-dev`
2. Adicione domínios DEV:
   - `api-dev.example.com`
   - `dev.example.com`
   - `admin-dev.example.com`

---

## 3. GitHub

### 3.1 Adicionar Secrets no GitHub (20 min)

1. **Acesse seu repositório:** https://github.com/seu-usuario/ecosistema-imob
2. Clique em **"Settings"** (aba superior)
3. No menu lateral esquerdo, clique em **"Secrets and variables"** → **"Actions"**
4. Para cada secret abaixo, clique em **"New repository secret"**

#### Secrets GCP

**Secret 1: GCP_PROJECT_ID_DEV**
- Name: `GCP_PROJECT_ID_DEV`
- Secret: `ecosistema-imob-dev`
- Clique em **"Add secret"**

**Secret 2: GCP_SA_KEY_DEV**
- Name: `GCP_SA_KEY_DEV`
- Secret: Abra o arquivo `github-actions-dev-key.json` no Notepad
- **Copie TODO o conteúdo** do arquivo (incluindo { e })
- Cole no campo Secret
- Clique em **"Add secret"**

**Secret 3: GCP_PROJECT_ID_PROD**
- Name: `GCP_PROJECT_ID_PROD`
- Secret: `ecosistema-imob-prod`

**Secret 4: GCP_SA_KEY_PROD**
- Name: `GCP_SA_KEY_PROD`
- Secret: Conteúdo completo do arquivo `github-actions-prod-key.json`

#### Secrets Firebase - DEV

**Secret 5: FIREBASE_PROJECT_ID_DEV**
- Name: `FIREBASE_PROJECT_ID_DEV`
- Secret: `ecosistema-imob-dev`

**Secret 6: FIREBASE_ADMIN_SDK_DEV**
- Name: `FIREBASE_ADMIN_SDK_DEV`
- Secret: Conteúdo completo do arquivo `firebase-adminsdk-dev.json`

**Secret 7: FIREBASE_API_KEY_DEV**
- Name: `FIREBASE_API_KEY_DEV`
- Secret: (o valor de apiKey que você anotou)

**Secret 8: FIREBASE_AUTH_DOMAIN_DEV**
- Name: `FIREBASE_AUTH_DOMAIN_DEV`
- Secret: `ecosistema-imob-dev.firebaseapp.com`

#### Secrets Firebase - PROD

**Secret 9: FIREBASE_PROJECT_ID_PROD**
- Name: `FIREBASE_PROJECT_ID_PROD`
- Secret: `ecosistema-imob-prod`

**Secret 10: FIREBASE_ADMIN_SDK_PROD**
- Name: `FIREBASE_ADMIN_SDK_PROD`
- Secret: Conteúdo completo do arquivo `firebase-adminsdk-prod.json`

**Secret 11: FIREBASE_API_KEY_PROD**
- Name: `FIREBASE_API_KEY_PROD`
- Secret: (o valor de apiKey PROD que você anotou)

**Secret 12: FIREBASE_AUTH_DOMAIN_PROD**
- Name: `FIREBASE_AUTH_DOMAIN_PROD`
- Secret: `ecosistema-imob-prod.firebaseapp.com`

#### Secrets Storage

**Secret 13: GCS_BUCKET_NAME_DEV**
- Name: `GCS_BUCKET_NAME_DEV`
- Secret: `ecosistema-imob-dev-storage`

**Secret 14: GCS_BUCKET_NAME_PROD**
- Name: `GCS_BUCKET_NAME_PROD`
- Secret: `ecosistema-imob-prod-storage`

#### Secrets Vercel (adicionar depois)

**Secret 15: VERCEL_TOKEN**
- Name: `VERCEL_TOKEN`
- Secret: (você vai obter na seção Vercel)

**Secret 16: NEXT_PUBLIC_API_URL_DEV**
- Name: `NEXT_PUBLIC_API_URL_DEV`
- Secret: (você vai obter após primeiro deploy do backend)

**Secret 17: NEXT_PUBLIC_API_URL_PROD**
- Name: `NEXT_PUBLIC_API_URL_PROD`
- Secret: (você vai obter após primeiro deploy do backend)

### 3.2 Criar Branch Develop (5 min)

1. No repositório GitHub, clique em **"main"** (dropdown de branches)
2. Digite no campo: `develop`
3. Clique em **"Create branch: develop from main"**
4. Pronto! Branch develop criada

### 3.3 Proteger Branch Main (5 min)

1. No repositório, clique em **"Settings"**
2. Menu lateral → **"Branches"**
3. Em "Branch protection rules", clique em **"Add rule"**
4. **Branch name pattern:** `main`
5. Marque as opções:
   - ✅ **Require a pull request before merging**
   - ✅ **Require approvals** (quantidade: 1)
   - ✅ **Require status checks to pass before merging**
   - ✅ **Require conversation resolution before merging**
   - ✅ **Do not allow bypassing the above settings**
6. Clique em **"Create"** (no final da página)

### 3.4 Criar Environment Production (5 min)

1. No repositório, clique em **"Settings"**
2. Menu lateral → **"Environments"**
3. Clique em **"New environment"**
4. **Name:** `production`
5. Clique em **"Configure environment"**
6. Configurações:
   - ✅ Marque **"Required reviewers"**
   - Adicione seu **usuário GitHub** como revisor
   - Em "Deployment branches", selecione **"Protected branches"**
7. Clique em **"Save protection rules"**

### 3.5 Fazer Commit dos Workflows (10 min)

#### Opção A: Via Interface Web do GitHub

1. No repositório, navegue para: `.github/workflows/`
2. Clique em **"Add file"** → **"Create new file"**
3. Nome do arquivo: `backend-deploy.yml`
4. Copie o conteúdo do arquivo [.github/workflows/backend-deploy.yml](.github/workflows/backend-deploy.yml)
5. Cole no editor
6. Role até o final, preencha:
   - **Commit message:** `ci: add backend deploy workflow`
   - Selecione: **"Create a new branch for this commit"**
   - Nome da branch: `feature/add-workflows`
7. Clique em **"Propose new file"**
8. Clique em **"Create pull request"**
9. Clique em **"Create pull request"** novamente
10. Clique em **"Merge pull request"** → **"Confirm merge"**

Repita o processo para:
- `frontend-public-deploy.yml`
- `frontend-admin-deploy.yml`

#### Opção B: Via Upload

1. Prepare os 3 arquivos `.yml` localmente
2. No GitHub, navegue para `.github/workflows/`
3. Clique em **"Add file"** → **"Upload files"**
4. Arraste os 3 arquivos `.yml`
5. Commit message: `ci: add CI/CD workflows`
6. Clique em **"Commit changes"**

---

## 4. Vercel

### 4.1 Criar Conta e Login (5 min)

1. **Acesse:** https://vercel.com
2. Clique em **"Sign Up"** (se não tem conta) ou **"Log In"**
3. Escolha: **"Continue with GitHub"**
4. Autorize o Vercel a acessar seu GitHub
5. Complete o cadastro

### 4.2 Criar Projeto Frontend Public (10 min)

1. No Vercel Dashboard, clique em **"Add New..."** → **"Project"**
2. Clique em **"Import"** ao lado do repositório `ecosistema-imob`
   - Se não aparecer, clique em **"Adjust GitHub App Permissions"**
3. Configure o projeto:
   - **Project Name:** `ecosistema-imob-public`
   - **Framework Preset:** `Next.js`
   - **Root Directory:** Clique em **"Edit"** → Selecione `frontend-public`
   - **Build Command:** `npm run build` (deixar padrão)
   - **Output Directory:** `.next` (deixar padrão)
   - **Install Command:** `npm install` (deixar padrão)
4. **Environment Variables:** (adicionar depois do primeiro deploy)
5. Clique em **"Deploy"**
6. Aguarde o deploy (~3-5 minutos)
7. **Anote a URL** gerada (ex: `ecosistema-imob-public.vercel.app`)

### 4.3 Criar Projeto Frontend Admin (10 min)

1. No Vercel Dashboard, clique em **"Add New..."** → **"Project"**
2. Clique em **"Import"** ao lado do repositório `ecosistema-imob`
3. Configure o projeto:
   - **Project Name:** `ecosistema-imob-admin`
   - **Framework Preset:** `Next.js`
   - **Root Directory:** Clique em **"Edit"** → Selecione `frontend-admin`
4. Clique em **"Deploy"**
5. Aguarde o deploy
6. **Anote a URL** gerada

### 4.4 Obter Token do Vercel (5 min)

1. No Vercel Dashboard, clique no seu **avatar** (canto superior direito)
2. Clique em **"Settings"**
3. Menu lateral → **"Tokens"**
4. Clique em **"Create"**
5. Preencha:
   - **Token Name:** `GitHub Actions CI/CD`
   - **Scope:** `Full Account`
   - **Expiration:** `No Expiration` (ou escolha período)
6. Clique em **"Create Token"**
7. **COPIE o token** (só será mostrado uma vez!)
8. **Guarde em local seguro**

Agora volte ao GitHub e adicione o secret:
- Name: `VERCEL_TOKEN`
- Secret: (cole o token copiado)

### 4.5 Configurar Environment Variables - Frontend Public (15 min)

#### Preview (DEV)

1. Acesse o projeto `ecosistema-imob-public` no Vercel
2. Clique em **"Settings"**
3. Menu lateral → **"Environment Variables"**
4. Adicione as seguintes variáveis (uma por vez):

**Variável 1:**
- **Key:** `NEXT_PUBLIC_API_URL`
- **Value:** `https://backend-api-xxxxx.run.app` (você vai atualizar depois)
- **Environments:** Marque apenas **Preview**
- Clique em **"Save"**

**Variável 2:**
- **Key:** `NEXT_PUBLIC_ENVIRONMENT`
- **Value:** `development`
- **Environments:** Marque apenas **Preview**
- Clique em **"Save"**

#### Production (PROD)

**Variável 3:**
- **Key:** `NEXT_PUBLIC_API_URL`
- **Value:** `https://api.example.com` (você vai atualizar depois)
- **Environments:** Marque apenas **Production**
- Clique em **"Save"**

**Variável 4:**
- **Key:** `NEXT_PUBLIC_ENVIRONMENT`
- **Value:** `production`
- **Environments:** Marque apenas **Production**
- Clique em **"Save"**

### 4.6 Configurar Environment Variables - Frontend Admin (20 min)

#### Preview (DEV)

1. Acesse o projeto `ecosistema-imob-admin` no Vercel
2. **Settings** → **Environment Variables**
3. Adicione:

**Variável 1:** `NEXT_PUBLIC_API_URL`
- Value: `https://backend-api-xxxxx.run.app`
- Environments: **Preview**

**Variável 2:** `NEXT_PUBLIC_FIREBASE_API_KEY`
- Value: (o apiKey DEV que você anotou)
- Environments: **Preview**

**Variável 3:** `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN`
- Value: `ecosistema-imob-dev.firebaseapp.com`
- Environments: **Preview**

**Variável 4:** `NEXT_PUBLIC_FIREBASE_PROJECT_ID`
- Value: `ecosistema-imob-dev`
- Environments: **Preview**

**Variável 5:** `NEXT_PUBLIC_ENVIRONMENT`
- Value: `development`
- Environments: **Preview**

#### Production (PROD)

Repita o processo acima, mas:
- Marque **Production**
- Use valores PROD (apiKey PROD, ecosistema-imob-prod, etc)
- `NEXT_PUBLIC_ENVIRONMENT`: `production`

---

## 5. Validação Final

### 5.1 Primeiro Deploy do Backend DEV (20 min)

#### Via GitHub Interface

1. Acesse seu repositório no GitHub
2. Navegue até: `backend/README.md` (ou qualquer arquivo do backend)
3. Clique no **ícone de lápis** (Edit this file)
4. Faça uma pequena alteração (adicione uma linha no README)
5. No final da página:
   - **Commit message:** `test: trigger backend deploy`
   - Selecione: **"Create a new branch"**
   - Nome: `test/backend-deploy`
6. Clique em **"Propose changes"**
7. Clique em **"Create pull request"**
8. Na página do PR:
   - **base:** `develop` (mude de main para develop)
   - **compare:** `test/backend-deploy`
9. Clique em **"Create pull request"**
10. Aguarde os checks do GitHub Actions
11. Quando ficar verde, clique em **"Merge pull request"**
12. Clique em **"Confirm merge"**

#### Acompanhar Deploy

1. Vá em **"Actions"** (aba superior do repositório)
2. Clique no workflow em execução: **"Backend Deploy"**
3. Acompanhe os logs em tempo real
4. Quando concluir (checkmark verde), clique em **"deploy-dev"**
5. Role até o final dos logs
6. **COPIE a URL** que aparece: `Backend DEV deployed to: https://backend-api-xxxxx.run.app`

#### Atualizar URL no GitHub e Vercel

1. **GitHub:** Atualize o secret `NEXT_PUBLIC_API_URL_DEV` com a URL copiada
2. **Vercel Frontend Public:** Atualize a env var `NEXT_PUBLIC_API_URL` (Preview)
3. **Vercel Frontend Admin:** Atualize a env var `NEXT_PUBLIC_API_URL` (Preview)

### 5.2 Testar Backend DEV (5 min)

1. **Abra uma nova aba** no navegador
2. Acesse: `https://backend-api-xxxxx.run.app/health` (substitua pela URL real)
3. Deve retornar algo como: `{"status":"ok"}`
4. Se retornar erro, verifique os logs no Cloud Run:
   - Acesse: https://console.cloud.google.com/run
   - Selecione projeto: `ecosistema-imob-dev`
   - Clique em `backend-api`
   - Vá em **"LOGS"**

### 5.3 Deploy do Frontend DEV (15 min)

1. No GitHub, edite qualquer arquivo do `frontend-public`
2. Crie PR para `develop`
3. Merge
4. Aguarde deploy automático
5. Acesse a URL do Vercel e teste

Repita para `frontend-admin`.

### 5.4 Deploy para PROD (30 min)

#### Criar PR de develop para main

1. No GitHub, clique em **"Pull requests"**
2. Clique em **"New pull request"**
3. Configure:
   - **base:** `main`
   - **compare:** `develop`
4. Clique em **"Create pull request"**
5. Adicione descrição: "Deploy inicial para produção"
6. Clique em **"Create pull request"**

#### Aprovar e Fazer Merge

1. (Se você for o revisor) Clique em **"Approve"** no PR
2. Aguarde checks passarem
3. Clique em **"Merge pull request"**
4. Clique em **"Confirm merge"**

#### Aprovar Deploy no Environment

1. Vá em **"Actions"**
2. Clique no workflow **"Backend Deploy"** em execução
3. O job `deploy-prod` vai aparecer com status **"Waiting"**
4. Clique em **"Review deployments"**
5. Marque **"production"**
6. Clique em **"Approve and deploy"**
7. Aguarde conclusão

#### Obter URL PROD

1. Quando concluir, copie a URL PROD
2. Atualize:
   - GitHub secret: `NEXT_PUBLIC_API_URL_PROD`
   - Vercel env vars (Production)

### 5.5 Teste Completo PROD (10 min)

1. **Backend PROD:** Acesse `/health`
2. **Frontend Public PROD:** Acesse a URL de produção
3. **Frontend Admin PROD:** Acesse e teste login

---

## 📋 Checklist de Validação

### GCP
- [ ] Projeto PROD criado
- [ ] APIs habilitadas (DEV e PROD)
- [ ] Firestore criado (DEV e PROD)
- [ ] Buckets criados e públicos
- [ ] Service Accounts criadas
- [ ] Chaves JSON baixadas

### Firebase
- [ ] Projeto PROD adicionado
- [ ] Authentication habilitada (DEV e PROD)
- [ ] Admin SDK keys baixadas
- [ ] API Keys anotadas
- [ ] Domínios autorizados configurados

### GitHub
- [ ] Todos os 17 secrets adicionados
- [ ] Branch develop criada
- [ ] Branch main protegida
- [ ] Environment production criado
- [ ] Workflows commitados

### Vercel
- [ ] Conta criada
- [ ] Projeto Frontend Public criado
- [ ] Projeto Frontend Admin criado
- [ ] Token gerado e adicionado ao GitHub
- [ ] Env vars configuradas (Preview e Production)

### Deploy
- [ ] Backend DEV deployado com sucesso
- [ ] Backend PROD deployado com sucesso
- [ ] Frontend Public deployado
- [ ] Frontend Admin deployado
- [ ] URLs atualizadas em todos os lugares

### Testes
- [ ] Backend DEV responde (`/health`)
- [ ] Backend PROD responde (`/health`)
- [ ] Frontend Public carrega
- [ ] Frontend Admin carrega
- [ ] Login funciona no Admin

---

## 🐛 Troubleshooting

### Deploy falha no GitHub Actions

**Sintoma:** Workflow fica vermelho

**Verificar:**
1. Clique no workflow com erro
2. Clique no job com erro
3. Expanda os steps para ver qual falhou
4. Leia a mensagem de erro

**Soluções comuns:**
- **"Secret not found":** Verifique se adicionou todos os secrets
- **"Permission denied":** Verifique roles da service account
- **"API not enabled":** Habilite a API no GCP Console

### Backend não inicia no Cloud Run

**Verificar logs:**
1. https://console.cloud.google.com/run
2. Selecione projeto
3. Clique em `backend-api`
4. Aba **"LOGS"**

**Soluções comuns:**
- Verificar se Firebase credentials estão corretas
- Verificar se variáveis de ambiente estão configuradas

### Frontend não carrega

**Verificar:**
1. Vercel Dashboard → Projeto → **"Deployments"**
2. Clique no deployment com erro
3. Veja os logs de build

**Soluções comuns:**
- Verificar se env vars estão configuradas
- Verificar se Root Directory está correto

---

## 📞 Suporte

Se ainda tiver problemas:

1. **Consulte os logs** do serviço específico
2. **Verifique o checklist** acima
3. **Leia a documentação completa:** [CONFIGURACAO_AMBIENTES_CICD.md](CONFIGURACAO_AMBIENTES_CICD.md)
4. **Entre em contato** com o time DevOps

---

## ⏱️ Tempo Total Estimado

| Fase | Tempo |
|------|-------|
| 1. GCP | 55 min |
| 2. Firebase | 35 min |
| 3. GitHub | 45 min |
| 4. Vercel | 65 min |
| 5. Validação | 80 min |
| **TOTAL** | **~4h 40min** |

---

**Última atualização:** 2026-01-07
**Versão:** 1.0
**Status:** ✅ Pronto para uso
