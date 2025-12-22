# Sistema de Cadastro de Imobiliárias - Implementação Unificada

## 📋 Visão Geral

Implementação completa de um sistema unificado de cadastro de imobiliárias que funciona em **múltiplos pontos de entrada** (Frontend Public e Frontend Admin), compartilhando o mesmo endpoint backend e componente React reutilizável.

## 🏗️ Arquitetura

```
┌─────────────────────────────────────────────────────────────┐
│                    BACKEND GO (API)                          │
│  POST /api/v1/auth/signup                                    │
│  - Cria Tenant no Firestore                                  │
│  - Cria Broker como Admin                                    │
│  - Cria usuário no Firebase Auth                             │
│  - Define custom claims (tenant_id, role)                    │
└─────────────────────────────────────────────────────────────┘
                            ▲
                            │
                ┌───────────┴───────────┐
                │                       │
┌───────────────▼─────────┐  ┌─────────▼──────────────┐
│   FRONTEND PUBLIC        │  │   FRONTEND ADMIN       │
│   (Porta 3000)           │  │   (Porta 3002)         │
│                          │  │                        │
│  Landing Page:           │  │  Páginas:              │
│  - Homepage com CTA      │  │  - /login              │
│  - Seção "Para           │  │  - /signup             │
│    Imobiliárias"         │  │                        │
│  - Link no header        │  │  Ambos redirecionam    │
│                          │  │  para o dashboard      │
│  Signup:                 │  │  após cadastro         │
│  - /cadastro-imobiliaria │  │                        │
│                          │  │                        │
└──────────────────────────┘  └────────────────────────┘
         │                              │
         └──────────┬───────────────────┘
                    │
         ┌──────────▼──────────┐
         │ COMPONENTE SHARED   │
         │ <SignupForm />      │
         │                     │
         │ Usado em ambos os   │
         │ frontends com props │
         │ diferentes          │
         └─────────────────────┘
```

## 📁 Estrutura de Arquivos

### Backend (Go)
```
backend/
└── internal/
    └── handlers/
        └── auth_handler.go
            └── POST /api/v1/auth/signup
                - Validação de dados
                - Criação de tenant
                - Criação de broker (primeiro = admin)
                - Geração de slug único
                - Custom claims no Firebase
```

### Frontend Admin
```
frontend-admin/
├── app/
│   ├── login/
│   │   └── page.tsx           # Login com link para signup
│   └── signup/
│       └── page.tsx            # Usa <SignupForm variant="standalone" />
└── components/
    └── auth/
        └── signup-form.tsx     # Componente compartilhado
```

### Frontend Public
```
frontend-public/
├── app/
│   ├── page.tsx                      # Homepage com CTA
│   └── cadastro-imobiliaria/
│       └── page.tsx                  # Usa <SignupForm variant="standalone" />
└── components/
    └── auth/
        └── signup-form.tsx           # Componente compartilhado (cópia)
```

## 🎯 Pontos de Entrada

### 1. Frontend Public - Landing Page (Melhor UX)

**URL**: `http://localhost:3000`

**Fluxo**:
1. Usuário acessa a homepage
2. Vê seção "Você é uma Imobiliária?" com benefícios
3. Clica em "Cadastre sua Imobiliária"
4. É redirecionado para `/cadastro-imobiliaria`
5. Preenche formulário
6. Após sucesso, é redirecionado para `http://localhost:3002/dashboard` (Admin)

**Elementos na Homepage**:
- Header: Link "Para Imobiliárias"
- Header: Botão "Login Admin"
- Seção CTA: Card com benefícios + botões de ação

### 2. Frontend Admin - Signup Direto

**URL**: `http://localhost:3002/signup`

**Fluxo**:
1. Usuário acessa diretamente a página de signup
2. Preenche formulário
3. Após sucesso, é redirecionado para `/dashboard`

**Acesso via**:
- Link na página de login: "Ainda não tem uma conta? Cadastre sua imobiliária"

## 🔧 Componente Reutilizável: `<SignupForm />`

### Props

```typescript
interface SignupFormProps {
  onSuccess?: () => void;           // Callback após sucesso
  redirectTo?: string;              // URL de redirecionamento
  variant?: 'standalone' | 'embedded'; // Estilo do componente
}
```

### Variantes

**`standalone`** (padrão):
- Página completa com fundo gradiente
- Card centralizado com logo
- Inclui footer com links de termos

**`embedded`**:
- Sem estilização externa
- Pode ser integrado em modais/páginas existentes

### Campos do Formulário

```typescript
{
  tenant_name: string;      // Nome da Imobiliária *
  name: string;             // Nome do Corretor/Admin *
  email: string;            // Email *
  phone: string;            // Telefone (+5511999999999) *
  password: string;         // Senha (min 6 chars) *
  confirmPassword: string;  // Confirmação de senha *
}
```

### Validações

✅ Nome da imobiliária obrigatório
✅ Nome do usuário obrigatório
✅ Email válido obrigatório
✅ Telefone no formato E.164 (internacional)
✅ Senha mínima de 6 caracteres
✅ Senhas devem coincidir

### Exemplo de Uso

**Frontend Admin**:
```tsx
<SignupForm
  variant="standalone"
  redirectTo="/dashboard"
/>
```

**Frontend Public**:
```tsx
<SignupForm
  variant="standalone"
  redirectTo="http://localhost:3002/dashboard"
/>
```

## 🔄 Fluxo Completo de Cadastro

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Usuário preenche formulário                              │
│    - Nome da imobiliária: "Imobiliária XYZ"                 │
│    - Nome: "João Silva"                                      │
│    - Email: "joao@xyz.com"                                   │
│    - Telefone: "+5511999999999"                              │
│    - Senha: "senha123"                                       │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Frontend envia POST /api/v1/auth/signup                  │
│    {                                                         │
│      "email": "joao@xyz.com",                                │
│      "password": "senha123",                                 │
│      "name": "João Silva",                                   │
│      "phone": "+5511999999999",                              │
│      "tenant_name": "Imobiliária XYZ"                        │
│    }                                                         │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. Backend processa                                          │
│    ✅ Criar usuário no Firebase Auth                         │
│    ✅ Criar Tenant no Firestore                              │
│       - ID: "tenant-abc123"                                  │
│       - Slug: "imobiliaria-xyz"                              │
│       - Status: "active"                                     │
│    ✅ Criar Broker na subcoleção                             │
│       - Role: "admin" (primeiro usuário)                     │
│       - user_id: Firebase UID                                │
│    ✅ Setar custom claims no Firebase                        │
│       - tenant_id: "tenant-abc123"                           │
│       - role: "admin"                                        │
│    ✅ Criar ActivityLog (tenant_created, broker_created)     │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. Backend retorna sucesso                                   │
│    {                                                         │
│      "tenant_id": "tenant-abc123",                           │
│      "broker_id": "broker-xyz",                              │
│      "firebase_token": "eyJhbGc...",                         │
│      "user": { ... }                                         │
│    }                                                         │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 5. Frontend faz login automático                            │
│    - signInWithEmailAndPassword(auth, email, password)       │
│    - Firebase Auth confirma e carrega custom claims          │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. Redireciona para Dashboard                               │
│    - Frontend Public → http://localhost:3002/dashboard       │
│    - Frontend Admin → /dashboard                             │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 Como Testar

### Pré-requisitos

1. **Backend rodando**: Porta 3000
2. **Frontend Public rodando**: Porta 3000
3. **Frontend Admin rodando**: Porta 3002
4. **Firebase configurado**: Credenciais nos `.env.local`

### Cenário 1: Signup via Landing Page (Frontend Public)

```bash
# 1. Acessar homepage
http://localhost:3000

# 2. Navegar até "Para Imobiliárias" ou clicar no CTA

# 3. Preencher formulário em /cadastro-imobiliaria
Nome da Imobiliária: Imobiliária Teste
Seu Nome: João da Silva
Email: joao@teste.com
Telefone: +5511987654321
Senha: teste123
Confirmar Senha: teste123

# 4. Clicar em "Criar Conta"

# 5. Verificar redirecionamento para http://localhost:3002/dashboard
```

### Cenário 2: Signup via Admin Direto

```bash
# 1. Acessar login do admin
http://localhost:3002/login

# 2. Clicar em "Cadastre sua imobiliária"

# 3. Preencher formulário em /signup
[mesmos dados acima com email diferente]

# 4. Verificar redirecionamento para /dashboard
```

### Cenário 3: Login após Cadastro

```bash
# 1. Fazer logout do dashboard

# 2. Acessar http://localhost:3002/login

# 3. Fazer login com credenciais criadas

# 4. Verificar acesso ao dashboard com dados do tenant
```

## ✅ Checklist de Funcionalidades

### Backend
- [x] Endpoint POST /api/v1/auth/signup implementado
- [x] Criação de Tenant no Firestore
- [x] Criação de Broker como admin
- [x] Geração de slug único
- [x] Custom claims no Firebase Auth
- [x] Validação de formato E.164 para telefone
- [x] Activity Log registrado

### Frontend Public
- [x] Homepage com seção CTA para imobiliárias
- [x] Link "Para Imobiliárias" no header
- [x] Botão "Login Admin" no header
- [x] Página `/cadastro-imobiliaria` funcional
- [x] Componente `<SignupForm />` implementado
- [x] Redirecionamento para admin após signup

### Frontend Admin
- [x] Página `/signup` funcional
- [x] Link de signup na página de login
- [x] Componente `<SignupForm />` implementado
- [x] Redirecionamento para dashboard após signup

### Componente Compartilhado
- [x] Props configuráveis (variant, redirectTo, onSuccess)
- [x] Validação de campos completa
- [x] Feedback de erros visual
- [x] Loading state durante processo
- [x] Toggle para mostrar/ocultar senha
- [x] Formatação de telefone internacional
- [x] Link para login (variante standalone)
- [x] Termos de uso e privacidade (variante standalone)

## 🎨 UX/UI

### Landing Page CTA
- Seção com fundo escuro (gray-900 to gray-800)
- Duas colunas: Benefícios + Steps
- Ícones de check para benefícios
- Card com passos numerados
- Botões primário e outline

### Formulário de Signup
- Fundo gradiente azul (standalone)
- Card branco centralizado
- Logo com ícone de prédio
- Campos com validação em tempo real
- Mensagens de erro destacadas
- Botão com loading spinner
- Links para login e termos

## 🔐 Segurança

### Validações Frontend
- Email formato válido
- Telefone formato E.164
- Senha mínima 6 caracteres
- Confirmação de senha match
- Campos obrigatórios

### Validações Backend
- Email único (não duplicado)
- Formato de dados
- Sanitização de inputs
- Slug único para tenant

### Firebase Auth
- Autenticação segura
- Custom claims protegidos
- Token JWT com expiração
- Middleware de autenticação

## 📊 Métricas e Monitoramento

### Activity Logs Criados
- `tenant_created`: Quando novo tenant é criado
- `broker_created`: Quando primeiro admin é criado

### Dados Rastreáveis
- Tenant ID
- Broker ID (admin)
- Timestamp de criação
- Email do admin
- Telefone do admin

## 🐛 Tratamento de Erros

### Erros Comuns

**409 Conflict**: Email já cadastrado
```
Mensagem: "Email já cadastrado. Faça login ou use outro email."
```

**400 Bad Request**: Dados inválidos
```
Mensagem: "Dados inválidos. Verifique os campos e tente novamente."
```

**Erro genérico**: Falha de conexão
```
Mensagem: "Erro ao criar conta. Tente novamente."
```

## 🔄 Próximos Passos

### Melhorias Futuras
- [ ] Verificação de email após cadastro
- [ ] Validação de telefone via SMS
- [ ] Upload de logo da imobiliária durante signup
- [ ] Wizard multi-step para onboarding
- [ ] Planos de assinatura (free, pro, enterprise)
- [ ] Convite para outros corretores pós-signup
- [ ] Tutorial guiado após primeiro login

### Integração com Outros Módulos
- [ ] Importação de imóveis (Prompt 02)
- [ ] Dashboard com métricas (Prompt 04b)
- [ ] Sistema de leads (Prompt 06)
- [ ] WhatsApp integration (Prompt 07)

## 📝 Notas de Implementação

### Diferenças entre Frontends

**Frontend Public**:
- `redirectTo`: `http://localhost:3002/dashboard` (URL completo para outro servidor)
- Link de login: `http://localhost:3002/login` (URL completo)

**Frontend Admin**:
- `redirectTo`: `/dashboard` (path relativo, mesmo servidor)
- Link de login: `/login` (path relativo)

### Firebase Config
Ambos os frontends usam as **mesmas credenciais Firebase** (mesmo projeto), pois compartilham:
- Mesma autenticação
- Mesmo Firestore
- Mesmos custom claims

### Ambiente de Desenvolvimento

```bash
# Frontend Public
PORT=3000 npm run dev

# Frontend Admin
PORT=3002 npm run dev

# Backend
PORT=3000 go run cmd/api/main.go
```

## 📞 Suporte

Para dúvidas ou problemas com a implementação:
1. Verificar logs do backend
2. Verificar console do navegador
3. Validar variáveis de ambiente
4. Conferir status do Firebase

---

**Última atualização**: 2025-12-22
**Status**: ✅ Implementação Completa
