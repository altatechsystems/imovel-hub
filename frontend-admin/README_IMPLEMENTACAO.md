# Frontend Admin - Implementação (Fase Inicial)

## Status da Implementação: 🔶 40% Concluído

Este documento descreve a implementação **inicial** do **Frontend Admin (Prompt 04b)** do ecosistema imobiliário.

## 📋 Resumo Executivo

O Frontend Admin foi configurado e a estrutura base foi implementada com sucesso. As funcionalidades principais de autenticação, layout e dashboard estão operacionais.

### ✅ Implementado

- [x] Setup do projeto Next.js 14 + TypeScript
- [x] Configuração de dependências (React Query, Zod, Firebase, Recharts)
- [x] Sistema de tipos TypeScript alinhado com backend
- [x] API Client configurado com auth automática
- [x] Firebase Client configurado
- [x] Autenticação protegida (AuthGuard)
- [x] Página de Login
- [x] Layout admin com Sidebar e Header
- [x] Dashboard com métricas básicas
- [x] Build do projeto sem erros ✅

### ⏳ Pendente (60%)

- [ ] CRUD completo de imóveis
- [ ] Upload de fotos (drag & drop)
- [ ] Gerenciamento de leads
- [ ] Sistema de importação UI
- [ ] Gerenciamento de proprietários
- [ ] Gerenciamento de corretores
- [ ] Relatórios e gráficos
- [ ] Configurações de tenant

## 🏗️ Arquitetura

### Estrutura de Diretórios

```
frontend-admin/
├── app/
│   ├── layout.tsx                 # Root layout com Providers
│   ├── page.tsx                   # Home (redirect para dashboard/login)
│   ├── login/
│   │   └── page.tsx              # Página de login
│   └── dashboard/
│       ├── layout.tsx            # Layout com sidebar
│       └── page.tsx              # Dashboard principal
├── components/
│   ├── auth-guard.tsx            # HOC para proteger rotas
│   ├── admin-sidebar.tsx         # Sidebar de navegação
│   └── admin-header.tsx          # Header com search e user
├── lib/
│   ├── api.ts                    # Cliente da API (admin endpoints)
│   ├── firebase.ts               # Configuração do Firebase
│   ├── providers.tsx             # React Query Provider
│   └── utils.ts                  # Funções utilitárias
├── hooks/
│   └── use-auth.ts               # Hook de autenticação
├── types/
│   ├── property.ts               # Tipos de propriedades
│   └── lead.ts                   # Tipos de leads
└── .env.local                    # Variáveis de ambiente
```

### Stack Tecnológica

- **Framework**: Next.js 16.1.0 (App Router)
- **Linguagem**: TypeScript 5
- **Estilização**: Tailwind CSS 4
- **State Management**: React Query (@tanstack/react-query)
- **Validação**: Zod + React Hook Form
- **HTTP Client**: Axios
- **Autenticação**: Firebase Auth
- **Backend Database**: Firestore (named database: imob-dev)
- **Gráficos**: Recharts 2.15.0

## 🎨 Páginas Implementadas

### 1. Login Page (`/login`)

**Arquivo**: [app/login/page.tsx](app/login/page.tsx)

**Funcionalidades**:
- Login com email e senha
- Integração com Firebase Auth
- Validação de formulário
- Error handling
- Loading states
- Redirect automático para dashboard após login

### 2. Dashboard (`/dashboard`)

**Arquivo**: [app/dashboard/page.tsx](app/dashboard/page.tsx)

**Funcionalidades**:
- Métricas principais em cards:
  - Total de imóveis
  - Imóveis disponíveis
  - Leads total
  - Leads novos
  - Proprietários
  - Negócios fechados
- Seção de imóveis recentes (placeholder)
- Seção de leads recentes (placeholder)
- Ações rápidas (Novo imóvel, Importar XML, etc.)
- Loading states

## 🔌 API Client

**Arquivo**: [lib/api.ts](lib/api.ts)

### Configuração

```typescript
baseURL: process.env.NEXT_PUBLIC_ADMIN_API_URL  // http://localhost:8080/api/admin
```

### Autenticação Automática

O client adiciona automaticamente o token do Firebase em todas as requisições:

```typescript
// Request interceptor - add auth token
this.client.interceptors.request.use(
  async (config) => {
    if (auth?.currentUser) {
      const token = await auth.currentUser.getIdToken();
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  }
);
```

### Endpoints Disponíveis

#### Properties
```typescript
getProperties(filters?, pagination?) → PropertyListResponse
getProperty(id) → Property
createProperty(data) → Property
updateProperty(id, data) → Property
deleteProperty(id) → void
updatePropertyStatus(id, status) → Property
updatePropertyVisibility(id, visibility) → Property
```

#### Leads
```typescript
getLeads(filters?, pagination?) → LeadListResponse
getLead(id) → Lead
updateLeadStatus(id, status) → Lead
assignLeadToBroker(leadId, brokerId) → Lead
```

#### Owners
```typescript
getOwners(pagination?) → OwnerListResponse
getOwner(id) → Owner
createOwner(data) → Owner
updateOwner(id, data) → Owner
```

#### Brokers
```typescript
getBrokers(pagination?) → BrokerListResponse
getBroker(id) → Broker
createBroker(data) → Broker
updateBroker(id, data) → Broker
```

#### Imports
```typescript
uploadImport(files) → ImportBatch
getImportBatch(batchId) → ImportBatch
getImportBatches(pagination?) → ImportBatchListResponse
```

#### Storage (Photos)
```typescript
uploadPropertyPhoto(propertyId, file) → Photo
deletePropertyPhoto(propertyId, imageId) → void
getPropertyPhotos(propertyId) → Photo[]
```

#### Dashboard Metrics
```typescript
getDashboardMetrics() → DashboardMetrics
getPropertyStats(period?) → PropertyStats
getLeadStats(period?) → LeadStats
```

## 🔐 Autenticação e Proteção de Rotas

### AuthGuard Component

**Arquivo**: [components/auth-guard.tsx](components/auth-guard.tsx)

Componente HOC que protege rotas autenticadas:
- Verifica se usuário está autenticado
- Redireciona para `/login` se não autenticado
- Mostra loading state durante verificação
- Usado no layout do dashboard

### Login Flow

1. Usuário acessa `/login`
2. Preenche email e senha
3. Firebase Auth valida credenciais
4. Token é armazenado automaticamente
5. Redirect para `/dashboard`
6. AuthGuard valida em cada rota protegida

## 📦 Componentes Principais

### AdminSidebar

**Arquivo**: [components/admin-sidebar.tsx](components/admin-sidebar.tsx)

**Menu Items**:
- Dashboard (/)
- Imóveis (/dashboard/imoveis)
- Leads (/dashboard/leads)
- Proprietários (/dashboard/proprietarios)
- Corretores (/dashboard/corretores)
- Importação (/dashboard/importacao)
- Relatórios (/dashboard/relatorios)
- Configurações (/dashboard/configuracoes)
- Sair (logout)

**Features**:
- Highlight automático do item ativo
- Ícones Lucide React
- Botão de logout integrado

### AdminHeader

**Arquivo**: [components/admin-header.tsx](components/admin-header.tsx)

**Features**:
- Barra de busca global
- Notificações (badge com count)
- Avatar e info do usuário
- Responsivo

## 🎯 Próximos Passos

### CRUD de Imóveis (Alta Prioridade)

1. **Listagem de Imóveis** (`/dashboard/imoveis`)
   - Tabela com paginação
   - Filtros (status, tipo, cidade)
   - Ordenação
   - Ações (editar, deletar, mudar status)

2. **Formulário de Criação/Edição**
   - Form wizard multi-step
   - Validação com Zod
   - Upload de fotos
   - Preview

3. **Detalhes do Imóvel**
   - Todas as informações
   - Histórico de atividades
   - Leads relacionados
   - Fotos

### Upload de Fotos

1. **Componente de Upload**
   - Drag & drop
   - Multi-file upload
   - Progress bar
   - Preview
   - Crop/resize (opcional)

2. **Galeria de Fotos**
   - Grid view
   - Reordenar (drag & drop)
   - Definir capa
   - Deletar

### Gerenciamento de Leads

1. **Listagem** (`/dashboard/leads`)
   - Tabela com status
   - Filtros (status, canal, data)
   - Atribuir corretor
   - Mudar status

2. **Detalhes do Lead**
   - Informações de contato
   - Propriedade relacionada
   - Histórico de interações
   - Notas

### Sistema de Importação UI

1. **Página de Importação** (`/dashboard/importacao`)
   - Upload de XML/XLS
   - Preview de dados
   - Validação
   - Progress tracking
   - Histórico de importações

2. **Detalhes do Batch**
   - Estatísticas
   - Erros/warnings
   - Imóveis importados
   - Logs

## 🚀 Como Executar

### Desenvolvimento

```bash
cd frontend-admin
npm install
npm run dev
```

Acesse: http://localhost:3001

**Login Test**:
- Email: (criar usuário no Firebase Console)
- Senha: (definir no Firebase Console)

### Build de Produção

```bash
npm run build
npm start
```

### Linting

```bash
npm run lint
```

## 🔗 Integração com Backend

### Endpoints Backend Necessários

O backend precisa implementar os seguintes endpoints admin:

```go
// Backend routes (cmd/server/main.go)
protected := api.Group("/admin")
protected.Use(authMiddleware.AuthRequired())
{
  // Properties
  GET    /admin/properties
  GET    /admin/properties/:id
  POST   /admin/properties
  PUT    /admin/properties/:id
  DELETE /admin/properties/:id
  PATCH  /admin/properties/:id/status
  PATCH  /admin/properties/:id/visibility

  // Leads
  GET    /admin/leads
  GET    /admin/leads/:id
  PATCH  /admin/leads/:id/status
  POST   /admin/leads/:id/assign

  // Owners
  GET    /admin/owners
  GET    /admin/owners/:id
  POST   /admin/owners
  PUT    /admin/owners/:id

  // Brokers
  GET    /admin/brokers
  GET    /admin/brokers/:id
  POST   /admin/brokers
  PUT    /admin/brokers/:id

  // Imports
  POST   /admin/tenants/:id/import
  GET    /admin/import-batches
  GET    /admin/import-batches/:id

  // Storage
  POST   /admin/properties/:id/images
  GET    /admin/properties/:id/images
  DELETE /admin/properties/:id/images/:image_id

  // Dashboard
  GET    /admin/dashboard/metrics
  GET    /admin/dashboard/property-stats
  GET    /admin/dashboard/lead-stats
}
```

## ✅ Checklist de Implementação

### Setup & Infraestrutura
- [x] Projeto Next.js 14 configurado
- [x] TypeScript configurado
- [x] Tailwind CSS configurado
- [x] Firebase SDK instalado
- [x] React Query provider
- [x] API client com auth

### Autenticação
- [x] Página de login
- [x] Firebase Auth integration
- [x] AuthGuard component
- [x] Protected routes
- [x] Logout functionality

### Layout & Navegação
- [x] AdminSidebar com menu
- [x] AdminHeader com search
- [x] Dashboard layout
- [x] Highlight de rota ativa

### Dashboard
- [x] Métricas principais
- [x] Cards de estatísticas
- [x] Ações rápidas
- [ ] Gráficos (Recharts)
- [ ] Tabelas de atividades recentes

### CRUD Imóveis
- [ ] Listagem com filtros
- [ ] Formulário de criação
- [ ] Formulário de edição
- [ ] Detalhes do imóvel
- [ ] Deletar imóvel
- [ ] Mudar status/visibilidade

### Upload de Fotos
- [ ] Componente de upload
- [ ] Drag & drop
- [ ] Preview
- [ ] Progress bar
- [ ] Galeria de fotos
- [ ] Reordenar fotos
- [ ] Deletar fotos

### Gerenciamento de Leads
- [ ] Listagem de leads
- [ ] Filtros e busca
- [ ] Detalhes do lead
- [ ] Atribuir corretor
- [ ] Mudar status
- [ ] Adicionar notas

### Importação
- [ ] Upload de XML/XLS
- [ ] Preview de dados
- [ ] Validação
- [ ] Progress tracking
- [ ] Listagem de batches
- [ ] Detalhes do batch

### Proprietários & Corretores
- [ ] Listagem
- [ ] CRUD completo
- [ ] Detalhes e histórico

### Relatórios
- [ ] Gráficos de performance
- [ ] Filtros de período
- [ ] Export para PDF/Excel

## 📊 Progresso por Funcionalidade

| Funcionalidade | Status | Progresso |
|----------------|--------|-----------|
| Setup & Config | ✅ Completo | 100% |
| Autenticação | ✅ Completo | 100% |
| Layout Admin | ✅ Completo | 100% |
| Dashboard Base | ✅ Completo | 80% |
| CRUD Imóveis | ⏳ Pendente | 0% |
| Upload Fotos | ⏳ Pendente | 0% |
| Leads | ⏳ Pendente | 0% |
| Importação UI | ⏳ Pendente | 0% |
| Proprietários | ⏳ Pendente | 0% |
| Corretores | ⏳ Pendente | 0% |
| Relatórios | ⏳ Pendente | 0% |

**Progresso Geral**: 40% ✅

## 🎨 Design System

### Cores

- **Primary**: Blue 600 (`#2563eb`)
- **Secondary**: Gray 900 (sidebar)
- **Success**: Green 600
- **Error**: Red 600
- **Warning**: Orange 600
- **Info**: Blue 500

### Layout

- **Sidebar**: 256px (w-64)
- **Header**: Altura automática
- **Content**: Flex-1 com overflow-y
- **Padding**: 24px (p-6)

### Ícones

Todos os ícones são do [Lucide React](https://lucide.dev/):
- Home, Building2, Users, UserCog, MessageSquare, Upload, BarChart3, Settings, LogOut

## 🔧 Configuração de Ambiente

```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080/api
NEXT_PUBLIC_ADMIN_API_URL=http://localhost:8080/api/admin

# Firebase Config
NEXT_PUBLIC_FIREBASE_PROJECT_ID=ecosistema-imob-dev
# ... outras variáveis Firebase

# Tenant
NEXT_PUBLIC_TENANT_ID=default-tenant-id
```

## 📚 Referências

- [Next.js Documentation](https://nextjs.org/docs)
- [React Query](https://tanstack.com/query)
- [Firebase Auth](https://firebase.google.com/docs/auth)
- [Lucide Icons](https://lucide.dev)
- [Recharts](https://recharts.org)

---

**Última Atualização**: 2025-12-22
**Status**: 🔶 Frontend Admin 40% Implementado
**Próximo**: CRUD de Imóveis + Upload de Fotos
