# Decisão Arquitetural: Frontends Separados por Contexto

**Data**: 2025-12-21
**Versão**: 1.0
**Status**: ✅ **APROVADO**
**Contexto**: Arquitetura de frontends para suporte a vendas + locação

---

## 📋 Resumo Executivo

**Decisão**: Implementar **frontends separados** por bounded context (Domain-Driven Design):
- `frontend-public` (portal público SEO)
- `frontend-admin-sales` (dashboard de vendas)
- `frontend-admin-rentals` (dashboard de locação - MVP+4)

**Alternativa Rejeitada**: Dashboard admin único (monolito frontend)

**Rationale**: Separação de contextos, personas distintas, deploy independente, UX otimizada, escalabilidade.

---

## 🤔 Problema

Com a adição da vertical de **locação/aluguel** (MVP+3 a MVP+5), precisamos decidir a arquitetura de frontend:

**Opção 1**: Dashboard admin único
```
/frontend-admin
  ├── /app/imoveis      (vendas)
  ├── /app/leads        (vendas)
  ├── /app/contratos    (locação) ← ADICIONA AQUI
  ├── /app/pagamentos   (locação)
  └── /app/manutencoes  (locação)
```

**Opção 2**: Frontends separados por contexto
```
/frontend-admin-sales    (dashboard vendas)
/frontend-admin-rentals  (dashboard locação) ← SEPARADO
```

---

## ⚖️ Opções Avaliadas

### Opção 1: Dashboard Admin Único (Monolito Frontend) ❌

**Descrição**: Adicionar rotas de locação ao dashboard admin existente.

**Vantagens**:
- ✅ Menos projetos para gerenciar (1 frontend admin)
- ✅ Design system compartilhado automaticamente
- ✅ Menos custo de infra inicial (1 deploy Vercel)

**Desvantagens**:
- ❌ **Personas diferentes**: Corretor de vendas vs Administrador de locação
- ❌ **Navegação confusa**: Menu sobrecarregado (vendas + locação + lançamentos)
- ❌ **Permissões complexas**: Imobiliária pode ter corretores apenas de vendas OU apenas de locação
- ❌ **Performance**: Bundle JavaScript pesado (~400kb+ com todas as features)
- ❌ **Deploy acoplado**: Bug na gestão de locação afeta dashboard de vendas
- ❌ **Escalabilidade limitada**: Difícil adicionar novos contextos (ex: construtoras)
- ❌ **Manutenção**: Código acoplado, difícil refatorar uma área sem afetar outra

**Conclusão**: ❌ **REJEITADO** - Problemas de UX, performance e escalabilidade.

---

### Opção 2: Frontends Separados por Contexto ✅

**Descrição**: Criar projetos Next.js separados para cada bounded context.

**Estrutura**:
```
/frontend-public          (Portal SEO - compradores/locatários)
/frontend-admin-sales     (Dashboard Vendas - corretores/imobiliárias)
/frontend-admin-rentals   (Dashboard Locação - gestores) ← MVP+4
```

**Vantagens**:
- ✅ **Separação de contextos** (Domain-Driven Design)
- ✅ **Personas distintas**: UX otimizada para cada perfil
  - Vendas: Foco em leads, imóveis, co-corretagem
  - Locação: Foco em contratos, inadimplência, manutenções
- ✅ **Performance**: Bundles menores
  - Sales: ~200kb (apenas features de vendas)
  - Rentals: ~250kb (calendário, pagamentos, SLA)
- ✅ **Deploy independente**: Bug em locação NÃO afeta vendas (zero downtime)
- ✅ **Desenvolvimento paralelo**: Equipes trabalham sem conflitos de merge
- ✅ **Permissões granulares**: Usuário pode ter acesso a um OU ambos dashboards
- ✅ **Escalabilidade**: Fácil adicionar novos contextos (ex: `frontend-admin-developers`)
- ✅ **Monitoramento**: Erros e métricas isoladas por contexto (Sentry/DataDog)
- ✅ **CI/CD otimizado**: Build/deploy apenas do que mudou

**Desvantagens**:
- ⚠️ Custo de infra maior: +R$ 100/mês por frontend (R$ 300/mês total)
- ⚠️ Compartilhamento de código manual (sem monorepo inicial)
- ⚠️ Mais projetos para gerenciar (3 frontends)

**Mitigações**:
- Custo adicional (R$ 100/mês) é **insignificante** vs economia de desenvolvimento (10-20h/mês = R$ 1.5k-3k/mês)
- Compartilhamento de código: copiar componentes manualmente no MVP, migrar para monorepo (Turborepo) no futuro
- Gestão de projetos: CI/CD automatizado (Vercel) reduz overhead

**Conclusão**: ✅ **APROVADO** - Benefícios de UX, performance e escalabilidade superam desvantagens.

---

### Opção 3: Micro-Frontends (Module Federation) ⚠️

**Descrição**: Usar Webpack Module Federation para compartilhar código entre frontends dinamicamente.

**Estrutura**:
```
/frontend-shell         (Container/Orchestrator)
  ├── @remote/sales     (Micro-frontend vendas)
  ├── @remote/rentals   (Micro-frontend locação)
  └── @remote/public    (Micro-frontend público)
```

**Vantagens**:
- ✅ Compartilhamento de código dinâmico (design system, utils)
- ✅ Deploy independente de cada micro-frontend
- ✅ Versionamento granular

**Desvantagens**:
- ❌ **Complexidade altíssima**: Webpack Module Federation, orquestração de dependências
- ❌ **Overhead de infraestrutura**: Múltiplos builds, CDN para cada micro-frontend
- ❌ **Debugging difícil**: Erros podem ocorrer entre frontends, stack traces complexos
- ❌ **Over-engineering**: Adequado para 50+ desenvolvedores, NOT para MVP/startup
- ❌ **TypeScript complexo**: Shared types entre micro-frontends, race conditions

**Conclusão**: ⚠️ **DESCARTADO** - Over-engineering para o estágio atual do projeto. Reavaliar quando equipe > 20 desenvolvedores.

---

## ✅ Decisão Final: Opção 2 (Frontends Separados)

### Estrutura de Projetos

```
ecosistema-imob/
├── backend/                    # Go/Gin (API única para todos)
│
├── frontend-public/            # Next.js (Portal SEO)
│   └── Domínio: www.example.com
│
├── frontend-admin-sales/       # Next.js (Dashboard Vendas)
│   └── Domínio: admin-vendas.example.com
│
├── frontend-admin-rentals/     # Next.js (Dashboard Locação) - MVP+4
│   └── Domínio: admin-locacao.example.com
│
└── shared/                     # Código compartilhado (opcional)
    ├── ui/                     # Design system (shadcn/ui)
    ├── lib/                    # API client, auth utilities
    └── types/                  # TypeScript types
```

### Autenticação Unificada

**Firebase Auth Compartilhado**:
- Usuário faz login UMA vez (Firebase Auth)
- Token JWT válido para TODOS os frontends admin
- Cookie httpOnly compartilhado entre subdomínios (`*.example.com`)

**Navegação Entre Dashboards**:
```typescript
// AppSwitcher component (compartilhado)
const apps = [
  {
    name: 'Vendas',
    url: 'https://admin-vendas.example.com',
    enabled: userPermissions.can_manage_properties,
  },
  {
    name: 'Locação',
    url: 'https://admin-locacao.example.com',
    enabled: userPermissions.can_manage_contracts,
  },
]
```

**Experiência do Usuário**:
1. Usuário faz login em `admin-vendas.example.com`
2. Clica em "Apps" (Grid icon) → vê "Locação" disponível
3. Clica em "Locação" → redireciona para `admin-locacao.example.com`
4. **Não precisa fazer login novamente** (token compartilhado)

### Permissões Granulares

```go
// backend/internal/models/broker.go
type BrokerRole string
const (
    BrokerRoleAdmin         BrokerRole = "admin"          // Acesso total
    BrokerRoleSalesAgent    BrokerRole = "sales_agent"    // Só vendas
    BrokerRoleRentalManager BrokerRole = "rental_manager" // Só locação
    BrokerRoleBoth          BrokerRole = "both"           // Vendas + Locação
)

type BrokerPermissions struct {
    // Vendas
    CanManageProperties  bool
    CanManageLeads       bool

    // Locação (MVP+4)
    CanManageContracts   bool
    CanManagePayments    bool
    CanManageMaintenance bool
}
```

**Controle de Acesso**:
- Usuário com `BrokerRoleSalesAgent`: Acessa APENAS `frontend-admin-sales`
- Usuário com `BrokerRoleRentalManager`: Acessa APENAS `frontend-admin-rentals`
- Usuário com `BrokerRoleBoth` ou `Admin`: Acessa AMBOS

---

## 💰 Custo de Infraestrutura

### Vercel (Hosting)

| Frontend | Custo/Mês | Justificativa |
|----------|-----------|---------------|
| `frontend-public` | R$ 100 | High traffic (SEO, público geral) |
| `frontend-admin-sales` | R$ 100 | Team collaboration, staging |
| `frontend-admin-rentals` | R$ 100 | Team collaboration, staging (MVP+4) |
| **Total** | **R$ 300/mês** | (~$60 USD/mês) |

### ROI do Custo Adicional

**Opção 1** (Dashboard Único): R$ 100/mês
**Opção 2** (Frontends Separados): R$ 300/mês
**Diferença**: +R$ 200/mês

**Economia de Desenvolvimento**:
- Menos bugs de deploy acoplado: -5h/mês
- Desenvolvimento paralelo sem conflitos: -5h/mês
- UX otimizada (menos refactoring): -5h/mês
- **Total**: -15h/mês × R$ 150/h = **-R$ 2.250/mês economizados**

**ROI**: +R$ 2.250/mês economizados - R$ 200/mês custo = **+R$ 2.050/mês líquido**

**Conclusão**: Custo adicional de R$ 200/mês é **amplamente compensado** por economia de desenvolvimento.

---

## 🚀 Roadmap de Implementação

### MVP (Agora) - Semanas 1-12

**Criar**:
- ✅ `backend/` (Go/Gin)
- ✅ `frontend-public/` (Next.js)
- ✅ `frontend-admin-sales/` (Next.js)

**NÃO criar ainda**:
- ❌ `frontend-admin-rentals/` (APENAS MVP+4)

### MVP+3 (Mês 7-9) - Anúncios de Aluguel

**Atualizar**:
- ✅ `frontend-public/` → adicionar rotas `/busca/aluguel`, `/imoveis/aluguel/[slug]`
- ✅ `backend/` → adicionar endpoints `/properties?transaction_type=rent`

**NÃO criar ainda**:
- ❌ `frontend-admin-rentals/` (gestão de contratos é MVP+4)

### MVP+4 (Mês 10-12) - Gestão de Contratos

**Criar**:
- ✅ `frontend-admin-rentals/` (NOVO projeto Next.js)
- ✅ Setup: `npx create-next-app@latest frontend-admin-rentals`
- ✅ Copiar design system de `frontend-admin-sales/components/ui/*`
- ✅ Configurar Firebase Auth (compartilhado)
- ✅ Deploy em `admin-locacao.example.com`

**Backend**:
- ✅ Adicionar handlers: `contract_handler.go`, `payment_handler.go`
- ✅ Endpoints: `/contracts`, `/payments`, `/maintenance`

---

## 📦 Compartilhamento de Código

### MVP (Manual) - Copiar/Colar

**Componentes Compartilhados**:
```bash
# Copiar manualmente
cp -r frontend-admin-sales/components/ui/* frontend-admin-rentals/components/ui/
cp frontend-admin-sales/lib/api.ts frontend-admin-rentals/lib/api.ts
cp frontend-admin-sales/lib/firebase.ts frontend-admin-rentals/lib/firebase.ts
```

**Vantagens**:
- ✅ Simples, sem setup complexo
- ✅ Cada frontend pode customizar componentes conforme necessário

**Desvantagens**:
- ⚠️ Duplicação de código
- ⚠️ Bugfix precisa ser aplicado manualmente em ambos

---

### Futuro (Monorepo Turborepo) - Quando Equipe > 5 Devs

**Setup**:
```bash
npm install turbo -g

# Estrutura
/apps
  /frontend-public
  /frontend-admin-sales
  /frontend-admin-rentals
/packages
  /ui              # Design system compartilhado
  /api-client      # Cliente HTTP compartilhado
  /auth            # Firebase Auth compartilhado
  /types           # TypeScript types

# Build
turbo run build --filter=frontend-admin-sales   # Build só vendas
turbo run build --filter=frontend-admin-rentals  # Build só locação
```

**Vantagens**:
- ✅ Zero duplicação de código
- ✅ Build cache inteligente (Turbo)
- ✅ Bugfix em 1 lugar, aplica em todos

**Desvantagens**:
- ⚠️ Setup inicial complexo (4-8h)
- ⚠️ Curva de aprendizado (Turborepo)

**Decisão**: Adiar para quando tiver 5+ desenvolvedores trabalhando simultaneamente.

---

## 📊 Comparativo de Opções

| Critério | Dashboard Único | Frontends Separados | Micro-Frontends |
|----------|----------------|---------------------|-----------------|
| **UX** | ⚠️ Navegação confusa | ✅ Otimizada por contexto | ✅ Otimizada |
| **Performance** | ❌ Bundle ~400kb | ✅ Bundles ~200kb | ✅ Bundles ~150kb |
| **Deploy** | ❌ Acoplado | ✅ Independente | ✅ Independente |
| **Desenvolvimento** | ⚠️ Conflitos de merge | ✅ Paralelo | ✅ Paralelo |
| **Permissões** | ⚠️ Complexas | ✅ Granulares | ✅ Granulares |
| **Escalabilidade** | ❌ Limitada | ✅ Alta | ✅ Altíssima |
| **Complexidade** | ✅ Baixa | ⚠️ Média | ❌ Altíssima |
| **Custo Infra** | ✅ R$ 100/mês | ⚠️ R$ 300/mês | ⚠️ R$ 400+/mês |
| **Manutenção** | ❌ Código acoplado | ✅ Código isolado | ⚠️ Overhead alto |
| **Adequação MVP** | ⚠️ OK para MVP | ✅ **IDEAL** | ❌ Over-engineering |

**Vencedor**: ✅ **Frontends Separados** - Melhor custo-benefício para MVP e escala futura.

---

## 🎯 Próximos Passos

### Semana 1 (Agora)

- [x] ✅ Documentar decisão arquitetural (este documento)
- [x] ✅ Atualizar `AI_DEV_DIRECTIVE.md` com nova estrutura de diretórios
- [ ] ⏳ Atualizar `prompts/04_frontend_mvp.txt` → renomear para `04_frontend_public_mvp.txt`
- [ ] ⏳ Criar `prompts/04b_frontend_admin_sales_mvp.txt` (novo prompt para dashboard vendas)

### MVP (Semanas 2-12)

- [ ] ⏳ Implementar `backend/` (Go/Gin)
- [ ] ⏳ Implementar `frontend-public/` (Portal SEO)
- [ ] ⏳ Implementar `frontend-admin-sales/` (Dashboard Vendas)
- [ ] ⏳ Configurar CI/CD (Vercel automático)

### MVP+4 (Mês 10-12)

- [ ] ⏳ Criar `frontend-admin-rentals/` (Dashboard Locação)
- [ ] ⏳ Implementar gestão de contratos (`/app/contratos`)
- [ ] ⏳ Implementar gestão de pagamentos (`/app/pagamentos`)
- [ ] ⏳ Setup Vercel deploy (`admin-locacao.example.com`)

---

## 📚 Referências

- [Domain-Driven Design (DDD)](https://martinfowler.com/bliki/DomainDrivenDesign.html)
- [Bounded Context Pattern](https://martinfowler.com/bliki/BoundedContext.html)
- [Micro-Frontends](https://micro-frontends.org/)
- [Turborepo Monorepo](https://turbo.build/repo/docs)
- [Vercel Multi-Project Deployment](https://vercel.com/docs/concepts/projects/overview)

---

**Versão**: 1.0
**Data**: 2025-12-21
**Aprovado por**: Equipe Altatech Systems + Claude Code
**Próxima Revisão**: MVP+3 (Mês 7, antes de implementar locação)
