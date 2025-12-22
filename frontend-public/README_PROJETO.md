# Frontend Público - Ecosistema Imob

Portal público para busca de imóveis e geração de leads.

## 🎯 Objetivo

Interface pública onde visitantes podem:
- Buscar imóveis disponíveis
- Filtrar por tipo, localização, preço, características
- Ver detalhes completos do imóvel
- Agendar visitas
- Entrar em contato via WhatsApp/Formulário

## 🏗️ Tecnologias

- **Next.js 15** - Framework React com SSR/SSG
- **TypeScript** - Tipagem estática
- **Tailwind CSS** - Estilização
- **Firebase Auth** - Autenticação (opcional para favoritos)
- **React Query** - Data fetching e cache
- **Zustand** - State management
- **React Hook Form** - Formulários
- **Zod** - Validação de schemas

## 📁 Estrutura do Projeto

```
frontend-public/
├── app/                    # App Router (Next.js 15)
│   ├── (public)/          # Layout público
│   │   ├── page.tsx       # Home - Busca de imóveis
│   │   ├── imoveis/       # Listagem e detalhes
│   │   ├── sobre/         # Sobre a imobiliária
│   │   └── contato/       # Formulário de contato
│   ├── layout.tsx         # Layout raiz
│   └── globals.css        # Estilos globais
├── components/            # Componentes React
│   ├── ui/               # Componentes de UI base
│   ├── property/         # Componentes de imóveis
│   ├── search/           # Componentes de busca
│   └── forms/            # Formulários
├── lib/                   # Bibliotecas e utilitários
│   ├── api.ts            # Cliente API
│   ├── firebase.ts       # Config Firebase
│   └── utils.ts          # Funções utilitárias
├── hooks/                 # Custom React hooks
├── types/                 # TypeScript types
└── public/               # Assets estáticos
```

## 🚀 Instalação

```bash
# Instalar dependências
npm install

# Configurar variáveis de ambiente
cp .env.example .env.local
# Editar .env.local com suas configurações

# Rodar em desenvolvimento
npm run dev

# Build para produção
npm run build
npm start
```

## 🔧 Configuração

### Variáveis de Ambiente

Criar arquivo `.env.local`:

```env
# API Backend
NEXT_PUBLIC_API_URL=http://localhost:8080/api

# Firebase
NEXT_PUBLIC_FIREBASE_API_KEY=your-api-key
NEXT_PUBLIC_FIREBASE_PROJECT_ID=your-project-id

# Tenant
NEXT_PUBLIC_TENANT_ID=your-tenant-id
NEXT_PUBLIC_TENANT_SLUG=your-slug
```

### Firebase

1. Criar projeto no [Firebase Console](https://console.firebase.google.com)
2. Habilitar Authentication (Email/Password, Google)
3. Copiar configurações do Web App
4. Adicionar ao `.env.local`

## 📄 Páginas Principais

### Home (`/`)
- Hero section com busca rápida
- Destaques de imóveis
- Filtros principais (tipo, cidade, faixa de preço)
- CTA para agendamento

### Listagem (`/imoveis`)
- Grid de imóveis
- Filtros avançados (sidebar)
- Ordenação (preço, data, relevância)
- Paginação
- Mapa de localização

### Detalhes (`/imoveis/[slug]`)
- Galeria de fotos
- Informações completas
- Mapa de localização
- Calculadora de financiamento
- Formulário de contato/agendamento
- Imóveis similares

### Sobre (`/sobre`)
- História da imobiliária
- Equipe de corretores
- Diferenciais
- Depoimentos de clientes

### Contato (`/contato`)
- Formulário de contato
- Informações de contato
- Mapa com localização
- Horário de atendimento

## 🎨 Design System

### Cores Principais
- Primary: Azul (#0066CC)
- Secondary: Laranja (#FF6B35)
- Success: Verde (#10B981)
- Error: Vermelho (#EF4444)

### Componentes Base
- Button (primary, secondary, outline, ghost)
- Input (text, email, tel, number, select)
- Card (property card, info card)
- Modal (contact, gallery, schedule)
- Badge (status, featured, new)

## 🔌 Integração com Backend

### Endpoints Utilizados

```typescript
// Listar imóveis públicos
GET /api/{tenant_id}/properties?status=available&visibility=public

// Detalhes do imóvel
GET /api/{tenant_id}/properties/{id}

// Imagens do imóvel
GET /api/{tenant_id}/properties/{property_id}/images

// Criar lead
POST /api/{tenant_id}/leads
```

## 📱 Responsividade

- **Mobile First**: Design otimizado para mobile
- **Breakpoints**:
  - sm: 640px (mobile landscape)
  - md: 768px (tablet)
  - lg: 1024px (desktop)
  - xl: 1280px (large desktop)

## 🔍 SEO

- Meta tags dinâmicas por página
- Open Graph para redes sociais
- Schema.org structured data
- Sitemap.xml automático
- robots.txt configurado

## 📊 Analytics

- Google Analytics 4
- Facebook Pixel
- Eventos customizados:
  - property_view
  - property_favorite
  - contact_form_submit
  - whatsapp_click

## 🚀 Deploy

### Vercel (Recomendado)

```bash
# Instalar Vercel CLI
npm i -g vercel

# Deploy
vercel

# Deploy em produção
vercel --prod
```

### Configurações no Vercel
1. Importar repositório do GitHub
2. Configurar variáveis de ambiente
3. Deploy automático em cada push

## 🧪 Testes

```bash
# Rodar testes
npm test

# Coverage
npm run test:coverage

# E2E com Playwright
npm run test:e2e
```

## 📝 Próximos Passos

- [ ] Implementar página de favoritos (requer auth)
- [ ] Adicionar comparador de imóveis
- [ ] Tour virtual 360°
- [ ] Chat em tempo real
- [ ] Notificações push
- [ ] PWA (Progressive Web App)

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Propriedade de Altatech Systems - Todos os direitos reservados.
