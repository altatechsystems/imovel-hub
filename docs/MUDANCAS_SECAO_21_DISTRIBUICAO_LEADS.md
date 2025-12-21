# Mudanças nos Prompts - Seção 21: Distribuição de Leads

**Data**: 2025-12-20
**Motivo**: Resolução de GAPS críticos identificados nas regras de distribuição de leads
**Referência**: AI_DEV_DIRECTIVE.md Seção 20 (renumerada para 21 na conclusão)

---

## 📋 Resumo Executivo

A Seção 21 foi adicionada ao AI_DEV_DIRECTIVE para resolver **6 GAPS críticos** que bloqueariam a implementação:

1. ✅ Algoritmo de seleção do primary broker
2. ✅ Notificação multi-corretor
3. ✅ Permissões de visualização de leads
4. ✅ Tratamento de formulário de contato
5. ✅ Mudança de primary broker
6. ✅ Validação de phone obrigatório

---

## 🔄 Mudanças Aplicadas

### ✅ AI_DEV_DIRECTIVE.md
- **Adicionado**: Seção 20 completa (Regras de Distribuição de Leads e Co-Corretagem)
- **Subseções**:
  - 20.1: Papéis de Corretores (captador, vendedor, co-corretor)
  - 20.2: Visibilidade Escalonada (private, network, marketplace, public)
  - 20.3: Fluxo "Tenho um Cliente"
  - 20.4: Algoritmo de Seleção do Primary Broker
  - 20.5: Notificação Multi-Corretor
  - 20.6: Permissões de Visualização de Leads
  - 20.7: Campo Phone Obrigatório
  - 20.8: Mudança de Primary Broker
  - 20.9: Cadastro Direto pelo Proprietário (futuro)
  - 20.10: Resumo Executivo

### ✅ PROMPT 01 (Foundation MVP)
**Arquivo**: `prompts/01_foundation_mvp.txt`

**Mudanças no modelo Property**:
```go
// ADICIONADO: Visibilidade escalonada
Visibility         PropertyVisibility `firestore:"visibility" json:"visibility"`
// Valores: private, network, marketplace, public

// ADICIONADO: Comissão de co-corretagem
CoBrokerCommission float64 `firestore:"co_broker_commission" json:"co_broker_commission"`
// Exemplo: 40.0 = 40% para selling_broker

// ATUALIZADO: PropertyVisibility enum
const (
    PropertyVisibilityPrivate     PropertyVisibility = "private"
    PropertyVisibilityNetwork     PropertyVisibility = "network"
    PropertyVisibilityMarketplace PropertyVisibility = "marketplace"
    PropertyVisibilityPublic      PropertyVisibility = "public"
)
```

**Endpoints a ADICIONAR**:

1. **POST /api/v1/tenants/:tenantId/properties/:propertyId/brokers/interest**
   - Corretor manifesta interesse ("Tenho um cliente")
   - Cria PropertyBrokerRole com role: "selling_broker"
   - Notifica captador
   - Retorna: PropertyBrokerRole criado

2. **PATCH /api/v1/tenants/:tenantId/properties/:propertyId/primary-broker**
   - Altera is_primary de um corretor para outro
   - Apenas captador ou admin podem alterar
   - Transação atômica (apenas 1 primary por vez)
   - Retorna: novo primary_broker_id

3. **GET /api/v1/tenants/:tenantId/brokers/:brokerId/leads**
   - Lista leads de Properties onde corretor possui PropertyBrokerRole
   - Filtragem automática por permissão
   - Retorna: []Lead

4. **Atualizar POST /api/v1/properties/:propertyId/leads/whatsapp**
   - Implementar algoritmo GetPrimaryBroker() (Seção 20.4)
   - Notificar TODOS os PropertyBrokerRole ativos (Seção 20.5)
   - Retornar whatsapp_url do primary

5. **Atualizar POST /api/v1/properties/:propertyId/leads/form**
   - Enviar email IMEDIATO para primary
   - Notificar outros corretores via dashboard
   - Retornar success message

### ✅ PROMPT 09 (Autenticação e Multi-tenancy)
**Arquivo**: `prompts/09_autenticacao_multitenancy.txt`

**Mudanças no modelo Broker**:
```go
// ATUALIZADO: Phone agora é OBRIGATÓRIO com validação E.164
Phone string `firestore:"phone" json:"phone" validate:"required,e164"`
// Exemplo: +5511999999999

// ADICIONAR validação no CreateBroker:
func CreateBroker(data BrokerInput) error {
    if data.Phone == "" {
        return errors.New("phone_required")
    }
    if !isValidE164(data.Phone) {
        return errors.New("invalid_phone_format")
    }
    // ...
}
```

**Impacto**: Todos os endpoints de signup/criação de broker devem validar phone.

---

## 📝 Prompts que PRECISAM de Atualização (Pendentes)

### ⚠️ PROMPT 07 (WhatsApp Flow)
**Arquivo**: `prompts/07_whatsapp_flow.txt`

**Adicionar**:

1. **Implementação do algoritmo GetPrimaryBroker()**
   ```go
   // Copiar da Seção 20.4 do AI_DEV_DIRECTIVE
   func GetPrimaryBroker(propertyID string) (*Broker, error) {
       // 1. Buscar is_primary = true
       // 2. Fallback: originating_broker
       // 3. Fallback: primeiro selling_broker
       // 4. Erro se nenhum tem phone
   }
   ```

2. **Tratamento de erro "no_phone_available"**
   - Frontend deve exibir apenas formulário de contato
   - Mensagem: "Imóvel indisponível para WhatsApp no momento"

3. **Notificação multi-corretor**
   - Primary → WhatsApp redirect (usuário final)
   - Outros → Email + Dashboard notification

### ⚠️ PROMPT 04b (Frontend Admin MVP)
**Arquivo**: `prompts/04b_frontend_admin_mvp.txt`

**Adicionar**:

1. **Busca Interna de Properties**
   - Endpoint: `GET /api/v1/tenants/:tenantId/properties/search`
   - Filtros: visibilidade (network, marketplace), tipo, cidade, preço
   - Resultado: Properties que corretor pode ver

2. **Componente PropertySearchCard** (no resultado da busca)
   ```tsx
   <PropertyCard>
     <PropertyInfo />
     <VisibilityBadge visibility={property.visibility} />
     <CommissionInfo>{property.co_broker_commission}%</CommissionInfo>

     {/* Botão principal */}
     <Button onClick={handleInterest}>
       Tenho um cliente para este imóvel
     </Button>
   </PropertyCard>
   ```

3. **Função handleInterest()**
   ```tsx
   const handleInterest = async (propertyId: string) => {
     await api.post(`/tenants/${tenantId}/properties/${propertyId}/brokers/interest`)
     toast.success('Interesse registrado! Captador foi notificado.')
   }
   ```

4. **LeadTable - Filtrar por PropertyBrokerRole**
   - Endpoint atual: `GET /tenants/:tenantId/leads` (ERRADO)
   - Novo endpoint: `GET /tenants/:tenantId/brokers/:brokerId/leads` (CORRETO)
   - Exibe apenas leads de Properties onde corretor participa

5. **PropertyForm - Campo Visibility**
   ```tsx
   <Select name="visibility">
     <Option value="private">Privado (apenas eu)</Option>
     <Option value="network">Rede (minha imobiliária)</Option>
     <Option value="marketplace">Marketplace (todos corretores)</Option>
     <Option value="public">Público (internet)</Option>
   </Select>

   {visibility === 'marketplace' && (
     <Input
       name="co_broker_commission"
       label="Comissão oferecida (%)"
       type="number"
       placeholder="40"
     />
   )}
   ```

### ⚠️ PROMPT 10 (Busca Pública)
**Arquivo**: `prompts/10_busca_publica.txt`

**Adicionar**:

1. **Filtro de visibilidade no backend**
   ```go
   // Busca PÚBLICA (frontend público)
   func SearchPublicProperties() {
       // WHERE visibility = 'public' AND status = 'available'
   }

   // Busca INTERNA (dashboard admin)
   func SearchInternalProperties(brokerID, tenantID string) {
       // WHERE visibility IN ('network', 'marketplace', 'public')
       // AND (
       //   (visibility = 'network' AND tenant_id = ?) OR
       //   (visibility = 'marketplace') OR
       //   (visibility = 'public')
       // )
   }
   ```

2. **Endpoint separado para busca interna**
   - Público: `GET /api/v1/properties/search` (apenas public)
   - Admin: `GET /api/v1/tenants/:tenantId/properties/search` (network + marketplace + public)

---

## 🎯 Cenários de Teste (Para PROMPT 03 e 05 - Auditorias)

### Teste 1: Visibilidade Escalonada
```
1. Captador cria Property com visibility: "private"
   → Apenas captador vê no admin
   → NÃO aparece no público
   → NÃO aparece para outros corretores

2. Captador muda para "network"
   → Todos do mesmo tenant veem no admin
   → Botão "Tenho cliente" aparece

3. Captador muda para "marketplace"
   → Todos os corretores (qualquer tenant) veem
   → Badge "Marketplace - 40% comissão"

4. Captador muda para "public"
   → Aparece no site público
   → Indexado pelo Google
```

### Teste 2: Fluxo "Tenho um Cliente"
```
1. Vendedor busca imóveis (visibility: marketplace)
2. Clica "Tenho um cliente" no Property X
3. Sistema cria PropertyBrokerRole (selling_broker)
4. Captador recebe email de notificação
5. Lead chega (WhatsApp):
   → Primary (captador) recebe redirect
   → Vendedor recebe email + dashboard
6. Vendedor vê lead na LeadTable
```

### Teste 3: Phone Obrigatório
```
1. Signup sem phone → ERRO "phone_required"
2. Signup com phone inválido → ERRO "invalid_phone_format"
3. Signup com +5511999999999 → SUCESSO
4. Lead criado, primary sem phone → fallback para próximo
5. Nenhum corretor tem phone → frontend exibe apenas formulário
```

### Teste 4: Mudança de Primary
```
1. Captador transfere primary para vendedor
2. Próximo lead vai para WhatsApp do vendedor
3. Notificação enviada ao novo primary
4. ActivityLog registra mudança
```

---

## 🚀 Ordem de Implementação Sugerida

1. **PROMPT 09 + PROMPT 01** (Foundation + Auth)
   - ✅ Phone obrigatório no Broker
   - ✅ Visibility e CoBrokerCommission no Property
   - ⚠️ Adicionar endpoints novos (interesse, primary, leads)

2. **PROMPT 07** (WhatsApp)
   - ⚠️ Implementar GetPrimaryBroker()
   - ⚠️ Notificação multi-corretor

3. **PROMPT 04b** (Frontend Admin)
   - ⚠️ Busca interna + botão "Tenho cliente"
   - ⚠️ LeadTable com filtro correto
   - ⚠️ PropertyForm com visibility

4. **PROMPT 10** (Busca)
   - ⚠️ Separar busca pública vs. interna
   - ⚠️ Filtros de visibilidade

5. **PROMPT 03 + 05** (Auditorias)
   - ⚠️ Adicionar testes dos cenários acima

---

## 📊 Impacto nas Decisões de Arquitetura

### Novo arquivo de decisão recomendado:

**`docs/decisions/007_visibilidade_escalonada.md`**
```markdown
# ADR 007: Visibilidade Escalonada de Properties

## Contexto
Necessidade de controlar quem pode ver imóveis para evitar duplicação
e permitir co-corretagem controlada.

## Decisão
Implementar 4 níveis de visibilidade:
- private: apenas captador
- network: imobiliária (tenant)
- marketplace: todos os corretores
- public: internet (SEO)

## Consequências
- Captador tem controle total
- Elimina duplicação (apenas 1 Property)
- Habilita marketplace cooperativo
- Preparado para cadastro pelo proprietário (futuro)
```

**`docs/decisions/008_distribuicao_leads.md`**
```markdown
# ADR 008: Distribuição de Leads Multi-Corretor

## Contexto
Múltiplos corretores podem estar envolvidos em um Property
(captador + vendedor + co-corretores).

## Decisão
- Lead pertence ao Property (não ao corretor)
- Primary recebe WhatsApp redirect (tempo real)
- Todos recebem notificação (email + dashboard)
- Algoritmo determinístico de fallback

## Consequências
- Transparência total
- Cooperação incentivada
- Nenhum lead perdido
- Auditoria completa
```

---

## ✅ Checklist de Conformidade

Antes de iniciar implementação, validar:

- [ ] AI_DEV_DIRECTIVE Seção 20 lida e compreendida
- [ ] Modelos Property e Broker atualizados
- [ ] Endpoints novos documentados
- [ ] Frontend Admin com busca interna planejado
- [ ] Algoritmo GetPrimaryBroker() copiado
- [ ] Notificações multi-corretor desenhadas
- [ ] Cenários de teste definidos
- [ ] ADRs 007 e 008 criados

---

**Documento gerado em**: 2025-12-20
**Por**: Claude Code (Análise + Consolidação)
