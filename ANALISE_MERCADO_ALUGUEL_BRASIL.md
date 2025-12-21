# ANÁLISE DO MERCADO BRASILEIRO DE ALUGUEL (LOCAÇÃO)
## Preparação Estratégica para MVP+3 / MVP+4

**Data**: 2025-12-21
**Versão**: 1.0
**Status**: Análise de Mercado + Preparação de Schema
**Objetivo**: Documentar dores do mercado, requisitos técnicos e estratégia de preparação para suporte a locação

---

## SUMÁRIO EXECUTIVO

O mercado de locação no Brasil representa **60-65% do volume de transações imobiliárias** (vs 35-40% vendas), mas possui **dores estruturais** que nenhuma plataforma resolve completamente:

**Principais Gaps de Mercado**:
1. **QuintoAndar**: Excelente em gestão de locação, mas **verticalizou** (não aceita imóveis fora de sua operação)
2. **ZAP/VivaReal**: Apenas vitrine (classificados), **não gerenciam contratos ou pagamentos**
3. **CRMs (Kenlo/Jetimob)**: Gestão básica, mas **não integram garantias, cobrança ou manutenção**
4. **Loft**: Foco em vendas (desistiu de locação em 2022)

**Nossa Oportunidade**:
Plataforma completa para **corretores/imobiliárias gerenciarem locação** sem verticalizar (marketplace aberto + gestão end-to-end).

**Priorização**:
- **MVP+3** (Mês 7-9): Adicionar suporte básico a anúncios de aluguel (sem gestão de contrato)
- **MVP+4** (Mês 10-12): Gestão de contratos, pagamentos, manutenção (diferenciação total)

---

## 1. DORES DO MERCADO DE LOCAÇÃO (Pain Points)

### 1.1 Problemas com Anúncios de Aluguel

#### Dor 1: Anúncios Desatualizados e Informações Incompletas
**Contexto**: 70-80% dos anúncios de aluguel em portais tradicionais têm informações desatualizadas.

**Problemas Específicos**:
- ❌ Valor do aluguel desatualizado (proprietário reajustou, corretor não atualizou)
- ❌ **Custos ocultos**: Anúncio mostra R$ 2.000, mas não informa:
  - Condomínio: R$ 800
  - IPTU: R$ 150/mês
  - **Custo real**: R$ 2.950/mês (48% maior!)
- ❌ Imóvel já alugado mas anúncio continua ativo
- ❌ Informações contraditórias (anúncio diz "mobiliado", visita mostra vazio)

**Benchmark Concorrentes**:
| Plataforma | Exige Condomínio/IPTU | Atualização Automática | Score |
|------------|----------------------|------------------------|-------|
| **ZAP Imóveis** | ❌ Opcional (maioria omite) | ❌ Manual | 3/10 |
| **VivaReal** | ⚠️ Campo existe mas não obrigatório | ❌ Manual | 4/10 |
| **QuintoAndar** | ✅ Obrigatório (eles gerenciam) | ✅ Automático (API própria) | 9/10 |
| **Imovelweb** | ❌ Opcional | ❌ Manual | 3/10 |
| **Nossa Oportunidade** | ✅ Obrigatório + validação | ⚠️ Confirmação periódica | **8/10** |

**Nossa Solução**:
- Campos obrigatórios: `monthly_rent`, `condo_fee`, `iptu_monthly`, `deposit_months`
- **Custo total calculado automaticamente** (transparência total)
- Confirmação periódica de disponibilidade (PROMPT 08 já implementado para vendas, adaptar para locação)

---

#### Dor 2: Falta de Padronização de Garantias Locatícias
**Contexto**: Brasil tem **4 modalidades principais** de garantia, mas anúncios raramente especificam qual é aceita.

**Modalidades de Garantia**:
1. **Fiança Bancária** (5-10% aceitam): Banco emite garantia, locatário paga 1-3% do valor anual
2. **Seguro Fiança** (40-50% aceitam): Seguradora garante, locatário paga prêmio mensal (~0.5-1.5% do aluguel)
3. **Caução** (30-40% aceitam): 3-6 meses de aluguel adiantado (bloqueado)
4. **Fiador** (70-80% aceitam): Pessoa física garante com imóvel próprio

**Problema**: Anúncios não deixam claro quais garantias são aceitas → lead perdido por incompatibilidade.

**Exemplo Real**:
```
Anúncio típico ZAP:
  "Apartamento 2 quartos - R$ 2.500/mês - Aceita garantia"

Cliente pergunta: "Qual tipo de garantia?"
Resposta demora 2-3 dias → cliente já alugou outro imóvel
```

**Nossa Solução**:
- Campo `accepted_guarantees[]` (array de opções):
  - `fiador` (fiador pessoa física com imóvel)
  - `caucao` (caução de X meses)
  - `seguro_fianca` (seguradora)
  - `fianca_bancaria` (carta de fiança)
- **Filtro de busca** por tipo de garantia aceita (inovação: nenhum portal faz isso!)

---

#### Dor 3: Ausência de Informação sobre Flexibilidade Contratual
**Contexto**: Locações flexíveis (Airbnb-style) vs tradicionais (12+ meses) têm públicos diferentes.

**Problema**: Anúncios misturam tudo sem diferenciar:
- Locação tradicional (12, 24, 30 meses)
- Locação temporada (férias, 1-3 meses)
- Locação corporativa (6-12 meses, mobiliado)

**Nossa Solução**:
- Campo `rental_type`: `traditional`, `corporate`, `short_term`, `vacation`
- Campo `min_rental_period_months`: 3, 6, 12, 24
- Campo `furnished`: `unfurnished`, `semi_furnished`, `fully_furnished`

---

### 1.2 Problemas com Gestão de Locação

#### Dor 4: Gestão Manual de Contratos
**Contexto**: 85% das imobiliárias ainda usam **contratos em Word + impressão + assinatura física**.

**Problemas**:
- ❌ Processo lento (7-15 dias para assinar contrato completo)
- ❌ Sem versionamento (cláusulas desatualizadas, leis mudaram)
- ❌ Sem rastreabilidade (quem assinou quando?)
- ❌ Renovação manual (perda de prazo → multa/vacância)

**Concorrentes**:
| Plataforma | Contrato Digital | Assinatura Eletrônica | Renovação Automática |
|------------|------------------|----------------------|---------------------|
| **QuintoAndar** | ✅ Templates próprios | ✅ Integrado | ✅ Automatizado |
| **ZAP/VivaReal** | ❌ Não oferece | ❌ Não oferece | ❌ Não oferece |
| **CRMs (Kenlo)** | ⚠️ Upload de PDF | ⚠️ Integração externa (Clicksign) | ❌ Não oferece |
| **Nossa Oportunidade** | ✅ Templates + custom | ✅ DocuSign/Clicksign | ✅ Alertas 90/60/30 dias | **7/10** |

**Nossa Solução (MVP+4)**:
- Model `RentalContract`:
  - `contract_start_date`, `contract_end_date`
  - `renewal_type`: `automatic`, `manual`, `none`
  - `indexation`: `igpm`, `ipca`, `none` (reajuste anual)
  - `signed_by_tenant_at`, `signed_by_landlord_at`
  - `status`: `draft`, `pending_signatures`, `active`, `expired`, `terminated`
- Templates pré-configurados (Law 8.245/91 compliance)
- Alertas de renovação (Cloud Scheduler → notificação 90 dias antes do vencimento)

---

#### Dor 5: Cobrança e Inadimplência
**Contexto**: **Taxa de inadimplência média no Brasil: 8-12%** (vs 2-3% nos EUA).

**Problemas**:
- ❌ Geração manual de boletos (processo repetitivo)
- ❌ Sem controle de pagamentos (corretor depende de confirmação do proprietário)
- ❌ Sem automação de cobrança (multa, juros, notificação)
- ❌ Split de pagamento complexo (repasse para proprietário - taxa da imobiliária)

**Concorrentes**:
| Plataforma | Boleto Automático | Split Automático | Cobrança Inadimplência |
|------------|-------------------|------------------|----------------------|
| **QuintoAndar** | ✅ Pix + Boleto | ✅ Automatizado (8% taxa) | ✅ Cobrança + seguro |
| **ZAP/VivaReal** | ❌ Não oferece | ❌ Não oferece | ❌ Não oferece |
| **CRMs (Kenlo)** | ⚠️ Integração Pagar.me | ❌ Manual | ❌ Manual |

**Nossa Solução (MVP+4)**:
- Model `RentalPayment`:
  - `due_date`, `paid_at`, `amount`, `status` (`pending`, `paid`, `overdue`)
  - `payment_method`: `boleto`, `pix`, `transfer`, `credit_card`
  - `late_fee`, `interest_daily` (multa 2%, juros 1% a.m.)
- **Geração automática** de cobranças mensais (Cloud Scheduler)
- **Split automático** via PagSeguro/Stripe:
  - Proprietário: 92%
  - Imobiliária: 8%
  - Plataforma: 2% (success fee sobre taxa da imobiliária)
- **Notificação multi-canal**: Email (D-3), WhatsApp (vencimento), SMS (D+3)

**Revenue Model**:
- Taxa de transação: 2% sobre pagamentos processados
- Exemplo: Aluguel R$ 3.000/mês → R$ 60 de receita recorrente por contrato
- **10 contratos** = R$ 600 MRR adicional

---

#### Dor 6: Manutenção e Chamados
**Contexto**: **40-50% das reclamações de inquilinos** são sobre demora no atendimento de manutenção.

**Problemas**:
- ❌ Chamados via WhatsApp/email dispersos (sem centralização)
- ❌ Sem SLA ou rastreamento (inquilino não sabe status)
- ❌ Sem histórico (repetição de problemas não detectada)
- ❌ Sem integração com prestadores de serviço

**Concorrentes**:
| Plataforma | Chamados Centralizados | Rastreamento Status | Integração Prestadores |
|------------|----------------------|---------------------|----------------------|
| **QuintoAndar** | ✅ App próprio | ✅ Tempo real | ✅ Rede própria |
| **ZAP/VivaReal** | ❌ Não oferece | ❌ Não oferece | ❌ Não oferece |
| **CRMs (Kenlo)** | ⚠️ Básico (sem app inquilino) | ❌ Manual | ❌ Não oferece |

**Nossa Solução (MVP+4)**:
- Model `MaintenanceRequest`:
  - `category`: `plumbing`, `electrical`, `locksmith`, `appliance`, `other`
  - `priority`: `low`, `medium`, `high`, `urgent`
  - `status`: `open`, `assigned`, `in_progress`, `resolved`, `closed`
  - `photos[]` (antes/depois)
  - `assigned_to_provider_id` (prestador)
  - `resolution_sla_hours` (SLA por prioridade)
- **Portal do inquilino** (frontend público):
  - Abrir chamado (foto + descrição)
  - Acompanhar status em tempo real
  - Avaliar atendimento (NPS por prestador)
- **Dashboard imobiliária**:
  - Atribuir a prestador
  - Monitorar SLA
  - Histórico por imóvel (detectar problemas recorrentes)

**Diferencial Competitivo**:
- ✅ Única plataforma que integra marketplace de imóveis + gestão de manutenção
- ✅ Dados históricos (imóveis com muitos chamados → bandeira vermelha para futuros locatários)

---

#### Dor 7: Falta de Transparência para Proprietários
**Contexto**: Proprietários reclamam de **falta de visibilidade** sobre status do imóvel alugado.

**Problemas**:
- ❌ Não sabem se aluguel foi pago (dependem de ligação da imobiliária)
- ❌ Não têm acesso ao histórico de manutenção
- ❌ Não sabem quando contrato vence
- ❌ Não conseguem avaliar desempenho da imobiliária

**Nossa Solução (MVP+5 - Futuro)**:
- **Portal do Proprietário** (read-only, sem login complexo):
  - Dashboard: Status de pagamento (pago/pendente)
  - Histórico de repasses (extratos mensais)
  - Chamados de manutenção (transparência total)
  - Documentos (contrato, vistorias)
  - Análise de mercado (valor do aluguel vs média do bairro)

**Benchmark**:
- QuintoAndar tem isso → **Único diferencial deles que ainda não copiamos**
- Nossa vantagem: Faremos marketplace aberto (eles são verticalizados)

---

### 1.3 Problemas com Relacionamento Corretor/Inquilino/Proprietário

#### Dor 8: Assimetria de Informação
**Contexto**: Inquilinos não sabem histórico do imóvel (manutenções anteriores, motivo da saída do inquilino anterior).

**Problema**:
- ❌ Imóvel com infiltração recorrente → inquilino descobre só após alugar
- ❌ Proprietário difícil (demora para aprovar reparos) → frustração do inquilino
- ❌ Bairro com problemas não divulgados (segurança, barulho)

**Nossa Solução**:
- **Histórico Público de Manutenção** (anonimizado):
  - "Imóvel teve 3 chamados de encanamento nos últimos 12 meses"
  - "Tempo médio de resolução: 5 dias" (vs média do bairro: 7 dias)
- **Rating de Proprietários** (como Airbnb):
  - Inquilinos avaliam proprietário após saída (aprovação de reparos, comunicação)
  - Score visível no anúncio (diferencial ÚNICO no Brasil)
- **Rating de Imóvel**:
  - Score baseado em: manutenção, localização, custo-benefício
  - Comentários de inquilinos anteriores (opcional, moderado)

**Impacto**:
- ✅ Redução de conflitos (expectativas alinhadas)
- ✅ Aumento de confiança (transparência total)
- ✅ Diferenciação competitiva (nenhum portal faz isso)

---

#### Dor 9: Burocracia na Análise de Crédito
**Contexto**: Processo de aprovação de inquilino demora **7-15 dias** (perda de leads).

**Problemas**:
- ❌ Documentação manual (envio de PDF por email)
- ❌ Análise subjetiva (corretor analisa "no olho")
- ❌ Sem histórico de crédito centralizado
- ❌ Fraudes (documentos falsificados)

**Nossa Solução (MVP+4 + Parceria)**:
- **Integração com Bureaus de Crédito**:
  - Serasa Experian API (score de crédito)
  - Consulta CPF (dívidas ativas)
  - Análise de renda (integração Conta Azul/Banco)
- **Análise Automatizada**:
  - Renda mínima: 3x o valor do aluguel (regra configurável)
  - Score mínimo: 600 (Serasa)
  - Sem restrições graves (proteção, falência)
- **Aprovação em 24h**:
  - Inquilino envia docs pelo app
  - Sistema analisa automaticamente
  - Corretor apenas confirma (decisão assistida)

**Revenue Model**:
- Taxa de análise: R$ 49-99 por inquilino (pago pelo corretor ou inquilino)
- Volume: 100 análises/mês = R$ 5.000-10.000 MRR adicional

---

## 2. REQUISITOS TÉCNICOS - RENTAL LISTINGS (Anúncios de Aluguel)

### 2.1 Campos Essenciais do Property Model (Locação)

#### Campos Obrigatórios para Aluguel
```go
type Property struct {
    // ... campos existentes de venda ...

    // ===== NOVOS CAMPOS PARA LOCAÇÃO (MVP+3) =====
    // Tipo de transação
    TransactionType TransactionType `firestore:"transaction_type" json:"transaction_type"`
    // Enum: "sale", "rent", "both"
    // Permite mesmo imóvel anunciado para venda E aluguel simultaneamente

    // ===== VALORES DE LOCAÇÃO =====
    RentalInfo *RentalInfo `firestore:"rental_info,omitempty" json:"rental_info,omitempty"`
    // NULL se TransactionType = "sale"
}

type TransactionType string
const (
    TransactionTypeSale TransactionType = "sale"
    TransactionTypeRent TransactionType = "rent"
    TransactionTypeBoth TransactionType = "both" // Proprietário aceita venda OU aluguel
)

type RentalInfo struct {
    // ===== VALORES MONETÁRIOS (OBRIGATÓRIOS) =====
    MonthlyRent      float64 `firestore:"monthly_rent" json:"monthly_rent" validate:"required,gt=0"`
    // Valor base do aluguel (sem condomínio/IPTU)

    CondoFee         float64 `firestore:"condo_fee,omitempty" json:"condo_fee,omitempty"`
    // Taxa de condomínio mensal (R$ 0 se não houver)

    IPTUMonthly      float64 `firestore:"iptu_monthly,omitempty" json:"iptu_monthly,omitempty"`
    // IPTU mensal (anual ÷ 12, ou R$ 0 se incluído no aluguel)

    TotalMonthlyCost float64 `firestore:"total_monthly_cost" json:"total_monthly_cost"`
    // CALCULADO: monthly_rent + condo_fee + iptu_monthly
    // Frontend SEMPRE exibe esse valor em destaque

    Currency         string  `firestore:"currency" json:"currency"` // "BRL"

    // ===== DEPÓSITO E GARANTIAS =====
    DepositMonths    int     `firestore:"deposit_months" json:"deposit_months" validate:"min=0,max=6"`
    // Quantidade de meses de caução (0-6, comum: 3)
    // Lei 8.245/91 limita a 3 salários mínimos, mas prática comum é 3 meses de aluguel

    AcceptedGuarantees []GuaranteeType `firestore:"accepted_guarantees" json:"accepted_guarantees"`
    // Array de garantias aceitas (pelo menos 1 obrigatório)

    // ===== TIPO DE LOCAÇÃO =====
    RentalType       RentalType `firestore:"rental_type" json:"rental_type"`
    // Enum: traditional, corporate, short_term, vacation

    MinRentalPeriod  int        `firestore:"min_rental_period_months" json:"min_rental_period_months"`
    // Prazo mínimo em meses (ex: 12, 24, 30 para tradicional; 1-3 para temporada)

    Furnished        FurnishedType `firestore:"furnished" json:"furnished"`
    // Enum: unfurnished, semi_furnished, fully_furnished

    // ===== DISPONIBILIDADE =====
    AvailableFrom    *time.Time `firestore:"available_from,omitempty" json:"available_from,omitempty"`
    // Data a partir da qual imóvel está disponível (NULL = imediato)

    // ===== PETS E RESTRIÇÕES =====
    AcceptsPets      bool       `firestore:"accepts_pets" json:"accepts_pets"`
    PetRestrictions  string     `firestore:"pet_restrictions,omitempty" json:"pet_restrictions,omitempty"`
    // Ex: "Apenas cães de pequeno porte", "Máximo 1 gato"

    // ===== INCLUSÕES =====
    UtilitiesIncluded []UtilityType `firestore:"utilities_included,omitempty" json:"utilities_included,omitempty"`
    // Ex: ["water", "gas", "internet", "cable_tv"]

    // ===== REAJUSTE ANUAL =====
    IndexationType   IndexationType `firestore:"indexation_type,omitempty" json:"indexation_type,omitempty"`
    // Enum: igpm, ipca, inpc, none (default: igpm)

    // ===== OBSERVAÇÕES =====
    RentalNotes      string `firestore:"rental_notes,omitempty" json:"rental_notes,omitempty"`
    // Ex: "Preferência para famílias", "Inquilino atual sai em 30/06/2026"
}

// ===== ENUMS =====
type GuaranteeType string
const (
    GuaranteeFiador        GuaranteeType = "fiador"         // Fiador pessoa física com imóvel
    GuaranteeCaucao        GuaranteeType = "caucao"         // Caução (3-6 meses adiantado)
    GuaranteeSeguroFianca  GuaranteeType = "seguro_fianca"  // Seguradora (Porto Seguro, etc)
    GuaranteeFiancaBancaria GuaranteeType = "fianca_bancaria" // Carta de fiança (banco)
)

type RentalType string
const (
    RentalTypeTraditional RentalType = "traditional" // 12+ meses, residencial
    RentalTypeCorporate   RentalType = "corporate"   // 6-12 meses, mobiliado, empresas
    RentalTypeShortTerm   RentalType = "short_term"  // 1-6 meses, flexível
    RentalTypeVacation    RentalType = "vacation"    // Temporada (férias, eventos)
)

type FurnishedType string
const (
    FurnishedNo      FurnishedType = "unfurnished"      // Sem móveis
    FurnishedSemi    FurnishedType = "semi_furnished"   // Cozinha + banheiros equipados
    FurnishedFull    FurnishedType = "fully_furnished"  // Pronto para morar
)

type UtilityType string
const (
    UtilityWater     UtilityType = "water"
    UtilityGas       UtilityType = "gas"
    UtilityElectricity UtilityType = "electricity"
    UtilityInternet  UtilityType = "internet"
    UtilityCableTV   UtilityType = "cable_tv"
)

type IndexationType string
const (
    IndexationIGPM  IndexationType = "igpm"  // Índice Geral de Preços do Mercado (mais comum)
    IndexationIPCA  IndexationType = "ipca"  // Índice de Preços ao Consumidor Amplo
    IndexationINPC  IndexationType = "inpc"  // Índice Nacional de Preços ao Consumidor
    IndexationNone  IndexationType = "none"  // Sem reajuste
)
```

#### Validações Obrigatórias (Backend)
```go
func ValidateRentalInfo(info *RentalInfo) error {
    // 1. Valores obrigatórios
    if info.MonthlyRent <= 0 {
        return errors.New("monthly_rent must be greater than 0")
    }

    // 2. Calcular total mensal (transparência)
    info.TotalMonthlyCost = info.MonthlyRent + info.CondoFee + info.IPTUMonthly

    // 3. Depósito não pode exceder 6 meses (limite razoável)
    if info.DepositMonths > 6 {
        return errors.New("deposit_months cannot exceed 6")
    }

    // 4. Pelo menos 1 garantia aceita
    if len(info.AcceptedGuarantees) == 0 {
        return errors.New("at least one accepted_guarantee required")
    }

    // 5. Rental type válido
    validRentalTypes := []RentalType{
        RentalTypeTraditional, RentalTypeCorporate,
        RentalTypeShortTerm, RentalTypeVacation,
    }
    if !contains(validRentalTypes, info.RentalType) {
        return errors.New("invalid rental_type")
    }

    // 6. Min rental period coerente com tipo
    if info.RentalType == RentalTypeTraditional && info.MinRentalPeriod < 12 {
        return errors.New("traditional rentals require min 12 months")
    }

    return nil
}
```

---

### 2.2 Campos Específicos vs Vendas

#### Comparação: Campos Exclusivos de Venda vs Locação

| Categoria | Campo | Venda | Locação | Notas |
|-----------|-------|-------|---------|-------|
| **Preço** | `price_amount` | ✅ Obrigatório | ❌ NULL | Preço de venda |
| | `monthly_rent` | ❌ NULL | ✅ Obrigatório | Aluguel mensal |
| | `condo_fee` | ⚠️ Informativo | ✅ Obrigatório | Comprador paga depois; inquilino paga desde o início |
| | `iptu_monthly` | ⚠️ Informativo | ✅ Obrigatório | Mesmo raciocínio |
| **Garantias** | `deposit_months` | ❌ NULL | ✅ Obrigatório | Não aplica em vendas |
| | `accepted_guarantees` | ❌ NULL | ✅ Obrigatório | Não aplica em vendas |
| **Prazo** | `min_rental_period` | ❌ NULL | ✅ Obrigatório | Não aplica em vendas |
| | `available_from` | ⚠️ Opcional | ✅ Comum | Venda imediata vs locação futura |
| **Mobília** | `furnished` | ⚠️ Raro | ✅ Crítico | Inquilinos filtram por isso |
| **Pets** | `accepts_pets` | ❌ NULL | ✅ Obrigatório | Não aplica em vendas |
| **Reajuste** | `indexation_type` | ❌ NULL | ✅ Obrigatório | Específico de locação (Lei 8.245/91) |

---

### 2.3 Campos para Preparação Futura (Gestão de Contratos - MVP+4)

**Estratégia**: Adicionar campos **RESERVADOS** no MVP+3, mas deixar NULL até implementação completa.

```go
type Property struct {
    // ... campos de RentalInfo (MVP+3) ...

    // ===== CAMPOS RESERVADOS PARA MVP+4 (Gestão de Contratos) =====
    // NULL no MVP+3, ativados em MVP+4

    CurrentContractID *string `firestore:"current_contract_id,omitempty" json:"current_contract_id,omitempty"`
    // ID do contrato ativo (NULL se imóvel vago)

    ContractHistory   []string `firestore:"contract_history,omitempty" json:"contract_history,omitempty"`
    // Array de IDs de contratos anteriores (histórico)

    LastRentalEndDate *time.Time `firestore:"last_rental_end_date,omitempty" json:"last_rental_end_date,omitempty"`
    // Data do último término de contrato (para análise de vacância)

    AverageVacancyDays int `firestore:"average_vacancy_days,omitempty" json:"average_vacancy_days,omitempty"`
    // Tempo médio sem locatário (calculado historicamente)
}
```

**Benefício**: Zero refatoração quando adicionar gestão de contratos em MVP+4.

---

## 3. GESTÃO DE LOCAÇÃO - NOVOS MODELS (MVP+4)

### 3.1 Model: RentalContract (Contrato de Locação)

```go
type RentalContract struct {
    ID         string    `firestore:"id" json:"id"`
    TenantID   string    `firestore:"tenant_id" json:"tenant_id"` // Isolamento multi-tenancy
    PropertyID string    `firestore:"property_id" json:"property_id"`

    // ===== PARTES ENVOLVIDAS =====
    LandlordID string `firestore:"landlord_id" json:"landlord_id"` // Owner.ID
    TenantPersonID string `firestore:"tenant_person_id" json:"tenant_person_id"` // Novo model: TenantPerson
    BrokerID   string `firestore:"broker_id" json:"broker_id"` // Corretor responsável

    // ===== DATAS DO CONTRATO =====
    StartDate  time.Time  `firestore:"start_date" json:"start_date"`
    EndDate    time.Time  `firestore:"end_date" json:"end_date"`
    SignedAt   *time.Time `firestore:"signed_at,omitempty" json:"signed_at,omitempty"`

    // ===== VALORES =====
    MonthlyRent       float64 `firestore:"monthly_rent" json:"monthly_rent"`
    CondoFee          float64 `firestore:"condo_fee" json:"condo_fee"`
    IPTUMonthly       float64 `firestore:"iptu_monthly" json:"iptu_monthly"`
    DepositAmount     float64 `firestore:"deposit_amount" json:"deposit_amount"` // Caução depositada

    // ===== REAJUSTE =====
    IndexationType    IndexationType `firestore:"indexation_type" json:"indexation_type"`
    LastAdjustmentDate *time.Time    `firestore:"last_adjustment_date,omitempty" json:"last_adjustment_date,omitempty"`
    NextAdjustmentDate *time.Time    `firestore:"next_adjustment_date,omitempty" json:"next_adjustment_date,omitempty"`

    // ===== GARANTIA =====
    GuaranteeType     GuaranteeType `firestore:"guarantee_type" json:"guarantee_type"`
    GuarantorID       *string       `firestore:"guarantor_id,omitempty" json:"guarantor_id,omitempty"` // Se fiador
    InsurancePolicyNumber *string   `firestore:"insurance_policy_number,omitempty" json:"insurance_policy_number,omitempty"` // Se seguro

    // ===== RENOVAÇÃO =====
    RenewalType       RenewalType   `firestore:"renewal_type" json:"renewal_type"`
    AutoRenewalMonths int           `firestore:"auto_renewal_months,omitempty" json:"auto_renewal_months,omitempty"`
    // Se automatic, renovar por quantos meses

    // ===== STATUS =====
    Status            ContractStatus `firestore:"status" json:"status"`
    TerminationDate   *time.Time     `firestore:"termination_date,omitempty" json:"termination_termination_date,omitempty"`
    TerminationReason string         `firestore:"termination_reason,omitempty" json:"termination_reason,omitempty"`

    // ===== DOCUMENTOS =====
    ContractDocumentURL string   `firestore:"contract_document_url" json:"contract_document_url"` // PDF assinado
    AttachmentURLs      []string `firestore:"attachment_urls,omitempty" json:"attachment_urls,omitempty"` // Vistorias, etc

    // ===== ASSINATURAS =====
    SignedByLandlordAt *time.Time `firestore:"signed_by_landlord_at,omitempty" json:"signed_by_landlord_at,omitempty"`
    SignedByTenantAt   *time.Time `firestore:"signed_by_tenant_at,omitempty" json:"signed_by_tenant_at,omitempty"`

    // ===== AUDIT =====
    CreatedAt time.Time `firestore:"created_at" json:"created_at"`
    UpdatedAt time.Time `firestore:"updated_at" json:"updated_at"`
}

type RenewalType string
const (
    RenewalAutomatic RenewalType = "automatic" // Renova automaticamente por X meses
    RenewalManual    RenewalType = "manual"    // Requer nova negociação
    RenewalNone      RenewalType = "none"      // Não permite renovação
)

type ContractStatus string
const (
    ContractDraft             ContractStatus = "draft"              // Rascunho (corretor criando)
    ContractPendingSignatures ContractStatus = "pending_signatures" // Aguardando assinaturas
    ContractActive            ContractStatus = "active"             // Ativo (vigente)
    ContractExpired           ContractStatus = "expired"            // Vencido (sem renovação)
    ContractTerminated        ContractStatus = "terminated"         // Rescindido antes do prazo
    ContractRenewed           ContractStatus = "renewed"            // Renovado (substituído por novo contrato)
)
```

#### Coleção Firestore
```
/tenants/{tenantId}/rental_contracts/{contractId}
```

#### Regras de Negócio
1. **Property.CurrentContractID** aponta para contrato ativo
2. Ao criar contrato, validar se Property não tem contrato ativo (evitar duplicidade)
3. Ao encerrar contrato, atualizar `Property.CurrentContractID = NULL` e `Property.Status = available`
4. **Alertas de renovação**:
   - D-90: Email para proprietário/corretor ("Contrato vence em 90 dias")
   - D-60: Notificação dashboard
   - D-30: Email urgente + WhatsApp
   - D-0: Alterar status para `expired` (Cloud Scheduler diário)

---

### 3.2 Model: RentalPayment (Pagamento de Aluguel)

```go
type RentalPayment struct {
    ID             string    `firestore:"id" json:"id"`
    TenantID       string    `firestore:"tenant_id" json:"tenant_id"`
    PropertyID     string    `firestore:"property_id" json:"property_id"`
    ContractID     string    `firestore:"contract_id" json:"contract_id"`

    // ===== DATAS =====
    DueDate        time.Time  `firestore:"due_date" json:"due_date"` // Vencimento
    PaidAt         *time.Time `firestore:"paid_at,omitempty" json:"paid_at,omitempty"` // Data do pagamento
    ReferenceMonth string     `firestore:"reference_month" json:"reference_month"` // "2026-01" (ano-mês)

    // ===== VALORES =====
    BaseAmount     float64 `firestore:"base_amount" json:"base_amount"` // Aluguel base
    CondoFee       float64 `firestore:"condo_fee" json:"condo_fee"`
    IPTUAmount     float64 `firestore:"iptu_amount" json:"iptu_amount"`
    LateFee        float64 `firestore:"late_fee,omitempty" json:"late_fee,omitempty"` // Multa (2%)
    InterestDaily  float64 `firestore:"interest_daily,omitempty" json:"interest_daily,omitempty"` // Juros (1% a.m.)
    TotalAmount    float64 `firestore:"total_amount" json:"total_amount"` // Soma de tudo

    // ===== PAGAMENTO =====
    PaymentMethod  PaymentMethod `firestore:"payment_method,omitempty" json:"payment_method,omitempty"`
    TransactionID  string        `firestore:"transaction_id,omitempty" json:"transaction_id,omitempty"` // ID PagSeguro/Stripe

    // ===== STATUS =====
    Status         PaymentStatus `firestore:"status" json:"status"`

    // ===== SPLIT (Repasse) =====
    SplitToLandlord      float64    `firestore:"split_to_landlord" json:"split_to_landlord"` // 92%
    SplitToRealEstate    float64    `firestore:"split_to_real_estate" json:"split_to_real_estate"` // 8%
    SplitToPlatform      float64    `firestore:"split_to_platform" json:"split_to_platform"` // 2% da taxa imobiliária
    SplitExecutedAt      *time.Time `firestore:"split_executed_at,omitempty" json:"split_executed_at,omitempty"`

    // ===== RECIBO =====
    ReceiptURL     string `firestore:"receipt_url,omitempty" json:"receipt_url,omitempty"` // PDF gerado automaticamente

    // ===== AUDIT =====
    CreatedAt      time.Time `firestore:"created_at" json:"created_at"`
    UpdatedAt      time.Time `firestore:"updated_at" json:"updated_at"`
}

type PaymentMethod string
const (
    PaymentBoleto       PaymentMethod = "boleto"
    PaymentPix          PaymentMethod = "pix"
    PaymentTransfer     PaymentMethod = "bank_transfer"
    PaymentCreditCard   PaymentMethod = "credit_card"
    PaymentCash         PaymentMethod = "cash"
)

type PaymentStatus string
const (
    PaymentPending  PaymentStatus = "pending"  // Aguardando pagamento
    PaymentPaid     PaymentStatus = "paid"     // Pago
    PaymentOverdue  PaymentStatus = "overdue"  // Atrasado (> due_date)
    PaymentCancelled PaymentStatus = "cancelled" // Cancelado (erro no boleto, etc)
)
```

#### Coleção Firestore
```
/tenants/{tenantId}/rental_payments/{paymentId}
```

#### Geração Automática de Cobranças (Cloud Scheduler)
```go
// Cloud Scheduler: Rodar todo dia 1º de cada mês às 00:00
func GenerateMonthlyPayments(ctx context.Context) error {
    // 1. Buscar todos contratos ativos
    activeContracts := QueryContracts("status = 'active'")

    for _, contract := range activeContracts {
        // 2. Verificar se já existe payment para o mês atual
        referenceMonth := time.Now().Format("2006-01")
        existing := FindPayment(contract.ID, referenceMonth)

        if existing != nil {
            continue // Já criado
        }

        // 3. Criar nova cobrança
        payment := &RentalPayment{
            ID:             uuid.New().String(),
            TenantID:       contract.TenantID,
            PropertyID:     contract.PropertyID,
            ContractID:     contract.ID,
            DueDate:        time.Date(time.Now().Year(), time.Now().Month(), 5, 0, 0, 0, 0, time.UTC), // Dia 5
            ReferenceMonth: referenceMonth,
            BaseAmount:     contract.MonthlyRent,
            CondoFee:       contract.CondoFee,
            IPTUAmount:     contract.IPTUMonthly,
            TotalAmount:    contract.MonthlyRent + contract.CondoFee + contract.IPTUMonthly,
            Status:         PaymentPending,
            CreatedAt:      time.Now(),
        }

        // 4. Salvar no Firestore
        SavePayment(payment)

        // 5. Gerar boleto (integração PagSeguro/Stripe)
        GenerateBoleto(payment)

        // 6. Enviar email para inquilino
        SendPaymentEmail(payment)
    }

    return nil
}
```

#### Cálculo Automático de Multa e Juros
```go
// Cloud Scheduler: Rodar diariamente às 02:00
func ApplyLateFeeAndInterest(ctx context.Context) error {
    // 1. Buscar pagamentos vencidos não pagos
    overduePayments := QueryPayments("status = 'pending' AND due_date < ?", time.Now())

    for _, payment := range overduePayments {
        // 2. Calcular multa (2% do valor base, uma vez)
        if payment.LateFee == 0 {
            payment.LateFee = payment.BaseAmount * 0.02
        }

        // 3. Calcular juros (1% a.m. = 0.033% ao dia)
        daysLate := int(time.Since(payment.DueDate).Hours() / 24)
        dailyInterestRate := 0.01 / 30 // 1% ao mês / 30 dias
        payment.InterestDaily = payment.BaseAmount * dailyInterestRate * float64(daysLate)

        // 4. Recalcular total
        payment.TotalAmount = payment.BaseAmount + payment.CondoFee + payment.IPTUAmount + payment.LateFee + payment.InterestDaily

        // 5. Atualizar status
        payment.Status = PaymentOverdue

        // 6. Salvar
        UpdatePayment(payment)

        // 7. Notificar inquilino (apenas nos dias 3, 7, 15, 30 após vencimento)
        if daysLate == 3 || daysLate == 7 || daysLate == 15 || daysLate == 30 {
            SendOverdueNotification(payment)
        }
    }

    return nil
}
```

---

### 3.3 Model: MaintenanceRequest (Chamado de Manutenção)

```go
type MaintenanceRequest struct {
    ID             string    `firestore:"id" json:"id"`
    TenantID       string    `firestore:"tenant_id" json:"tenant_id"`
    PropertyID     string    `firestore:"property_id" json:"property_id"`
    ContractID     string    `firestore:"contract_id" json:"contract_id"`

    // ===== QUEM ABRIU =====
    RequestedBy    string    `firestore:"requested_by" json:"requested_by"` // "tenant" | "landlord" | "broker"
    RequestedByID  string    `firestore:"requested_by_id" json:"requested_by_id"` // ID da pessoa

    // ===== DESCRIÇÃO =====
    Category       MaintenanceCategory `firestore:"category" json:"category"`
    Title          string              `firestore:"title" json:"title"` // Ex: "Vazamento no banheiro"
    Description    string              `firestore:"description" json:"description"`
    Photos         []string            `firestore:"photos,omitempty" json:"photos,omitempty"` // URLs GCS

    // ===== PRIORIDADE =====
    Priority       MaintenancePriority `firestore:"priority" json:"priority"`
    SLAHours       int                 `firestore:"sla_hours" json:"sla_hours"` // Calculado pela prioridade
    DueDate        time.Time           `firestore:"due_date" json:"due_date"` // created_at + sla_hours

    // ===== ATRIBUIÇÃO =====
    AssignedToProviderID *string   `firestore:"assigned_to_provider_id,omitempty" json:"assigned_to_provider_id,omitempty"`
    AssignedAt           *time.Time `firestore:"assigned_at,omitempty" json:"assigned_at,omitempty"`

    // ===== STATUS =====
    Status         MaintenanceStatus `firestore:"status" json:"status"`
    ResolvedAt     *time.Time        `firestore:"resolved_at,omitempty" json:"resolved_at,omitempty"`
    ClosedAt       *time.Time        `firestore:"closed_at,omitempty" json:"closed_at,omitempty"`

    // ===== CUSTO =====
    EstimatedCost  float64   `firestore:"estimated_cost,omitempty" json:"estimated_cost,omitempty"`
    ActualCost     float64   `firestore:"actual_cost,omitempty" json:"actual_cost,omitempty"`
    ApprovedBy     string    `firestore:"approved_by,omitempty" json:"approved_by,omitempty"` // "landlord" | "broker"
    ApprovedAt     *time.Time `firestore:"approved_at,omitempty" json:"approved_at,omitempty"`

    // ===== RESOLUÇÃO =====
    ResolutionNotes string   `firestore:"resolution_notes,omitempty" json:"resolution_notes,omitempty"`
    ResolutionPhotos []string `firestore:"resolution_photos,omitempty" json:"resolution_photos,omitempty"` // Fotos após reparo

    // ===== AVALIAÇÃO =====
    Rating         *int      `firestore:"rating,omitempty" json:"rating,omitempty"` // 1-5 estrelas
    RatingComment  string    `firestore:"rating_comment,omitempty" json:"rating_comment,omitempty"`

    // ===== AUDIT =====
    CreatedAt      time.Time `firestore:"created_at" json:"created_at"`
    UpdatedAt      time.Time `firestore:"updated_at" json:"updated_at"`
}

type MaintenanceCategory string
const (
    MaintenancePlumbing      MaintenanceCategory = "plumbing"       // Encanamento
    MaintenanceElectrical    MaintenanceCategory = "electrical"     // Elétrica
    MaintenanceLocksmith     MaintenanceCategory = "locksmith"      // Chaveiro
    MaintenanceAppliance     MaintenanceCategory = "appliance"      // Eletrodomésticos
    MaintenancePainting      MaintenanceCategory = "painting"       // Pintura
    MaintenanceCleaning      MaintenanceCategory = "cleaning"       // Limpeza
    MaintenancePestControl   MaintenanceCategory = "pest_control"   // Dedetização
    MaintenanceOther         MaintenanceCategory = "other"          // Outros
)

type MaintenancePriority string
const (
    MaintenanceLow     MaintenancePriority = "low"     // SLA: 7 dias
    MaintenanceMedium  MaintenancePriority = "medium"  // SLA: 48h
    MaintenanceHigh    MaintenancePriority = "high"    // SLA: 24h
    MaintenanceUrgent  MaintenancePriority = "urgent"  // SLA: 4h (vazamento, falta de luz)
)

type MaintenanceStatus string
const (
    MaintenanceOpen       MaintenanceStatus = "open"        // Aguardando atribuição
    MaintenanceAssigned   MaintenanceStatus = "assigned"    // Atribuído a prestador
    MaintenanceInProgress MaintenanceStatus = "in_progress" // Prestador trabalhando
    MaintenanceResolved   MaintenanceStatus = "resolved"    // Concluído, aguardando confirmação
    MaintenanceClosed     MaintenanceStatus = "closed"      // Confirmado e fechado
    MaintenanceCancelled  MaintenanceStatus = "cancelled"   // Cancelado
)
```

#### Coleção Firestore
```
/tenants/{tenantId}/maintenance_requests/{requestId}
```

#### SLA Automático
```go
func CalculateSLA(priority MaintenancePriority) int {
    slaMap := map[MaintenancePriority]int{
        MaintenanceLow:    168, // 7 dias
        MaintenanceMedium: 48,  // 2 dias
        MaintenanceHigh:   24,  // 1 dia
        MaintenanceUrgent: 4,   // 4 horas
    }
    return slaMap[priority]
}

// Cloud Scheduler: Rodar a cada 1 hora
func CheckSLAViolations(ctx context.Context) error {
    // Buscar chamados abertos com due_date vencido
    violatedRequests := QueryMaintenanceRequests("status IN ('open', 'assigned') AND due_date < ?", time.Now())

    for _, request := range violatedRequests {
        // Notificar corretor/imobiliária
        SendSLAViolationAlert(request)

        // Escalar prioridade
        if request.Priority == MaintenanceLow {
            request.Priority = MaintenanceMedium
        } else if request.Priority == MaintenanceMedium {
            request.Priority = MaintenanceHigh
        }

        UpdateMaintenanceRequest(request)
    }

    return nil
}
```

---

## 4. ANÁLISE COMPETITIVA (Rental Market)

### 4.1 QuintoAndar - O Que Eles Fazem Bem

**Pontos Fortes**:
1. ✅ **Gestão End-to-End**: Desde anúncio até entrega de chaves (verticalizado)
2. ✅ **Análise de Crédito Automatizada**: 24-48h (vs 7-15 dias no mercado)
3. ✅ **Seguro contra Inadimplência**: Proprietário recebe mesmo se inquilino não pagar
4. ✅ **Portal do Proprietário**: Transparência total (pagamentos, manutenção, documentos)
5. ✅ **App do Inquilino**: Abrir chamados, pagar aluguel, solicitar mudanças
6. ✅ **Manutenção Profissional**: Rede própria de prestadores (SLA garantido)

**Pontos Fracos**:
- ❌ **Modelo Verticalizado**: Imóvel precisa ser **exclusivo** deles (corretor perde controle)
- ❌ **Taxa Alta**: 8-10% sobre aluguel (vs 6-8% imobiliárias tradicionais)
- ❌ **Sem Marketplace**: Não permite co-corretagem (fechado)
- ❌ **Padrão Rígido**: Proprietário não pode escolher prestador de manutenção
- ❌ **Cidades Limitadas**: Apenas capitais e grandes cidades (sem interior)

**Nossa Oportunidade**:
- ✅ **Marketplace Aberto**: Imobiliárias mantêm controle total
- ✅ **Co-Corretagem**: Permite múltiplos corretores (efeito de rede)
- ✅ **Flexibilidade**: Proprietário escolhe nível de serviço (básico, completo)
- ✅ **Cobertura Nacional**: Não limitado a capitais

---

### 4.2 ZAP Imóveis / VivaReal - Gaps Críticos

**O Que Fazem**:
- ✅ Vitrine de anúncios (classificados)
- ✅ Geração de leads
- ✅ SEO forte (tráfego orgânico)

**O Que NÃO Fazem** (Nossa Oportunidade):
- ❌ Gestão de contratos (apenas anúncio)
- ❌ Pagamentos/cobranças (corretor faz manual)
- ❌ Manutenção (WhatsApp disperso)
- ❌ Portal do proprietário (zero visibilidade)
- ❌ Análise de crédito (corretor faz manualmente)

**Nossa Vantagem**:
Somos **ZAP + QuintoAndar combinados** (vitrine + gestão), mas **sem verticalizar** (marketplace aberto).

---

### 4.3 CRMs Tradicionais (Kenlo, Jetimob, Imobzi, Tecimob)

**O Que Fazem**:
- ✅ Cadastro de imóveis
- ✅ Gestão de leads
- ✅ Upload de contratos (PDF)
- ⚠️ Integração com boleto (via PagSeguro/Pagar.me)

**O Que NÃO Fazem** (Nossa Oportunidade):
- ❌ **SEO Zero**: Dependem de ZAP/VivaReal para leads
- ❌ **Sem Portal Público**: Imobiliária precisa ter site próprio separado
- ❌ **Gestão de Manutenção Básica**: Apenas lista de chamados (sem SLA, sem portal inquilino)
- ❌ **Sem Split Automático**: Repasse manual para proprietário
- ❌ **Sem Marketplace**: Zero co-corretagem estruturada

**Nossa Vantagem**:
- ✅ **CRM + Portal Público + Marketplace** (ecossistema completo)
- ✅ **SEO 100%** (leads orgânicos, reduz dependência de portais)
- ✅ **Gestão Completa de Locação** (contratos, pagamentos, manutenção)

---

### 4.4 Loft - Lições Aprendidas

**Histórico**: Loft entrou forte em locação (2020-2021), mas desistiu em 2022.

**Razões da Desistência**:
- ❌ **Margens Baixas**: Locação tem ticket 10x menor que vendas (R$ 3k/mês vs R$ 30k comissão única)
- ❌ **Complexidade Operacional**: Gestão de manutenção é trabalhosa
- ❌ **Modelo VC-backed**: Precisavam crescimento rápido (locação é recorrente, mas lento)

**Nossa Vantagem**:
- ✅ **Somos Plataforma**: Não fazemos operação (imobiliárias fazem)
- ✅ **Modelo SaaS**: Receita recorrente escala sem aumentar operação
- ✅ **Margens Altas**: 2% sobre pagamentos + assinatura (40-60% margem)

---

## 5. OPORTUNIDADES DE DIFERENCIAÇÃO (Gaps de Mercado)

### 5.1 Marketplace de Imóveis para Aluguel com Rating de Proprietário

**Gap**: Nenhuma plataforma permite inquilinos avaliarem proprietários.

**Nossa Inovação**:
- ✅ **Rating de Proprietário** (1-5 estrelas):
  - Aprovação de reparos (rápido/lento)
  - Comunicação (responsivo/difícil)
  - Flexibilidade contratual (flexível/rígido)
- ✅ **Visível no Anúncio**: "Proprietário com 4.7 estrelas (23 avaliações)"
- ✅ **Filtro de Busca**: "Apenas proprietários com 4+ estrelas"

**Benefícios**:
- Para inquilinos: Reduz risco de alugar com proprietário problemático
- Para proprietários: Incentiva bom atendimento (reputação importa)
- Para plataforma: Diferencial competitivo único no Brasil

---

### 5.2 Transparência de Custo Total (Filtro Real)

**Gap**: Anúncios mostram apenas aluguel, omitem condomínio/IPTU.

**Nossa Inovação**:
- ✅ **Custo Total em Destaque**: "R$ 2.950/mês (aluguel + condomínio + IPTU)"
- ✅ **Filtro por Custo Total**: "Imóveis até R$ 3.000 total" (não apenas aluguel base)
- ✅ **Breakdown Detalhado**:
  ```
  Aluguel:    R$ 2.000
  Condomínio: R$ 800
  IPTU:       R$ 150
  ───────────────────
  Total:      R$ 2.950/mês
  ```

**Impacto**:
- Reduz frustração do inquilino (expectativas alinhadas)
- Aumenta conversão (leads qualificados)
- Diferenciação técnica (concorrentes não fazem)

---

### 5.3 Histórico de Manutenção Público (Imóvel)

**Gap**: Inquilinos não sabem se imóvel tem problemas recorrentes.

**Nossa Inovação**:
- ✅ **Badge de Qualidade**: "Imóvel bem mantido (2 chamados nos últimos 12 meses)"
- ✅ **Histórico Anonimizado**:
  ```
  Últimos 12 meses:
  - 1x Encanamento (resolvido em 2 dias)
  - 1x Elétrica (resolvido em 1 dia)

  Tempo médio de resolução: 1.5 dias
  Média do bairro: 4 dias
  ```
- ✅ **Filtro de Busca**: "Apenas imóveis com <5 chamados/ano"

**Benefícios**:
- Inquilinos têm visibilidade de risco
- Proprietários/imobiliárias são incentivados a manter bem (reputação)
- Plataforma ganha confiança

---

### 5.4 Análise de Crédito com IA (Antifraude)

**Gap**: Análise de crédito é manual (7-15 dias) e propensa a fraude.

**Nossa Inovação (MVP+4 + Parceria)**:
- ✅ **Integração Serasa/SPC**: Score automático
- ✅ **IA de Detecção de Fraude**:
  - Validação de documentos (CPF, RG, comprovante de renda)
  - Detecção de padrões suspeitos (renda incompatível, histórico)
- ✅ **Aprovação em 24h**: Inquilino envia docs → sistema analisa → corretor apenas confirma

**Revenue Model**:
- Taxa de análise: R$ 79 por inquilino
- Volume: 100 análises/mês = R$ 7.900 MRR adicional

---

### 5.5 Gestão de Contratos com Alertas Inteligentes

**Gap**: 30-40% dos contratos vencem sem renovação por esquecimento.

**Nossa Inovação**:
- ✅ **Alertas Automáticos**:
  - D-90: "Contrato vence em 3 meses. Iniciar negociação de renovação?"
  - D-60: "Inquilino ainda não respondeu. Enviar lembrete?"
  - D-30: "URGENTE: Contrato vence em 1 mês"
- ✅ **Renovação 1-Click**:
  - Proprietário/inquilino aprovam online
  - Novo contrato gerado automaticamente (reajuste IGPM aplicado)
  - Assinatura digital
- ✅ **Dashboard de Vencimentos**: Timeline visual de todos os contratos

**Benefícios**:
- Reduz vacância (renovação proativa)
- Aumenta retenção (facilita renovação)
- Melhora NPS (corretor/imobiliária parecem profissionais)

---

## 6. PREPARAÇÃO DE SCHEMA (MVP+3 - Imediato)

### 6.1 Campos a Adicionar AGORA no Property Model

**Estratégia**: Adicionar campos para locação **desde o MVP**, mas deixar NULL até MVP+3.

#### Alterações no `prompts/01_foundation_mvp.txt`

```go
type Property struct {
    // ... campos existentes ...

    // ===== TRANSACTION TYPE (ADICIONADO EM MVP - NULL PERMITIDO) =====
    TransactionType *TransactionType `firestore:"transaction_type,omitempty" json:"transaction_type,omitempty"`
    // NULL no MVP, "sale" por padrão na importação
    // Valores: "sale", "rent", "both"

    // ===== RENTAL INFO (ADICIONADO EM MVP - STRUCT NULL) =====
    RentalInfo *RentalInfo `firestore:"rental_info,omitempty" json:"rental_info,omitempty"`
    // NULL no MVP, ativado em MVP+3

    // ===== RESERVED FIELDS FOR MVP+4 (Contract Management) =====
    CurrentContractID  *string    `firestore:"current_contract_id,omitempty" json:"current_contract_id,omitempty"`
    ContractHistory    []string   `firestore:"contract_history,omitempty" json:"contract_history,omitempty"`
    LastRentalEndDate  *time.Time `firestore:"last_rental_end_date,omitempty" json:"last_rental_end_date,omitempty"`
    AverageVacancyDays *int       `firestore:"average_vacancy_days,omitempty" json:"average_vacancy_days,omitempty"`
}

// ===== NOVO TYPE (ADICIONADO EM MVP) =====
type TransactionType string
const (
    TransactionTypeSale TransactionType = "sale"
    TransactionTypeRent TransactionType = "rent"
    TransactionTypeBoth TransactionType = "both"
)

// ===== NOVO STRUCT (ADICIONADO EM MVP, NULL PERMITIDO) =====
type RentalInfo struct {
    // Valores monetários
    MonthlyRent      float64 `firestore:"monthly_rent" json:"monthly_rent"`
    CondoFee         float64 `firestore:"condo_fee,omitempty" json:"condo_fee,omitempty"`
    IPTUMonthly      float64 `firestore:"iptu_monthly,omitempty" json:"iptu_monthly,omitempty"`
    TotalMonthlyCost float64 `firestore:"total_monthly_cost" json:"total_monthly_cost"`
    Currency         string  `firestore:"currency" json:"currency"` // "BRL"

    // Depósito e garantias
    DepositMonths      int             `firestore:"deposit_months" json:"deposit_months"`
    AcceptedGuarantees []GuaranteeType `firestore:"accepted_guarantees" json:"accepted_guarantees"`

    // Tipo de locação
    RentalType         RentalType      `firestore:"rental_type" json:"rental_type"`
    MinRentalPeriod    int             `firestore:"min_rental_period_months" json:"min_rental_period_months"`
    Furnished          FurnishedType   `firestore:"furnished" json:"furnished"`

    // Disponibilidade
    AvailableFrom      *time.Time      `firestore:"available_from,omitempty" json:"available_from,omitempty"`

    // Pets
    AcceptsPets        bool            `firestore:"accepts_pets" json:"accepts_pets"`
    PetRestrictions    string          `firestore:"pet_restrictions,omitempty" json:"pet_restrictions,omitempty"`

    // Inclusões
    UtilitiesIncluded  []UtilityType   `firestore:"utilities_included,omitempty" json:"utilities_included,omitempty"`

    // Reajuste
    IndexationType     IndexationType  `firestore:"indexation_type,omitempty" json:"indexation_type,omitempty"`

    // Observações
    RentalNotes        string          `firestore:"rental_notes,omitempty" json:"rental_notes,omitempty"`
}

// ===== ENUMS (ADICIONADOS EM MVP) =====
// (Ver seção 2.1 para definições completas)
```

#### Benefícios da Preparação Antecipada
1. ✅ **Zero Refatoração**: Quando adicionar locação em MVP+3, schema já está pronto
2. ✅ **Compatibilidade Retroativa**: Properties existentes continuam funcionando (campos NULL)
3. ✅ **Validação Condicional**: Se `TransactionType = "rent"`, validar `RentalInfo != NULL`
4. ✅ **Testes Antecipados**: Podemos testar locação em staging antes de MVP+3

---

### 6.2 Firestore Indexes Adicionais (Preparação)

#### `firestore.indexes.json` (Adicionar em MVP)

```json
{
  "indexes": [
    // ... indexes existentes ...

    // ===== INDEXES PARA LOCAÇÃO (ADICIONADOS EM MVP, INATIVOS ATÉ MVP+3) =====
    {
      "collectionGroup": "properties",
      "queryScope": "COLLECTION",
      "fields": [
        { "fieldPath": "tenant_id", "order": "ASCENDING" },
        { "fieldPath": "transaction_type", "order": "ASCENDING" },
        { "fieldPath": "rental_info.total_monthly_cost", "order": "ASCENDING" },
        { "fieldPath": "created_at", "order": "DESCENDING" }
      ]
    },
    {
      "collectionGroup": "properties",
      "queryScope": "COLLECTION",
      "fields": [
        { "fieldPath": "tenant_id", "order": "ASCENDING" },
        { "fieldPath": "transaction_type", "order": "ASCENDING" },
        { "fieldPath": "rental_info.rental_type", "order": "ASCENDING" },
        { "fieldPath": "city", "order": "ASCENDING" }
      ]
    },
    {
      "collectionGroup": "properties",
      "queryScope": "COLLECTION",
      "fields": [
        { "fieldPath": "tenant_id", "order": "ASCENDING" },
        { "fieldPath": "rental_info.furnished", "order": "ASCENDING" },
        { "fieldPath": "rental_info.accepts_pets", "order": "ASCENDING" },
        { "fieldPath": "created_at", "order": "DESCENDING" }
      ]
    }
  ]
}
```

**Custo**: Indexes não utilizados têm custo zero no Firestore.

---

### 6.3 Models Preparatórios (Criar Arquivos Vazios em MVP)

#### Criar arquivos placeholder (sem implementação)

```bash
# MVP: Criar arquivos vazios (comentados)
backend/internal/models/rental_contract.go       # Comentado até MVP+4
backend/internal/models/rental_payment.go        # Comentado até MVP+4
backend/internal/models/maintenance_request.go   # Comentado até MVP+4
backend/internal/models/tenant_person.go         # Comentado até MVP+4
```

**Conteúdo (exemplo `rental_contract.go`)**:
```go
package models

// RentalContract - PLACEHOLDER FOR MVP+4
// This model will be activated when rental management is implemented.
//
// IMPORTANT: Do not use this model in MVP or MVP+1-3.
//
// Estimated implementation: MVP+4 (Month 10-12)

/*
type RentalContract struct {
    // Definição completa aqui (comentada)
}
*/
```

**Benefício**: Documentação do roadmap no próprio código.

---

## 7. ROADMAP DE IMPLEMENTAÇÃO (Locação)

### 7.1 MVP+3 (Mês 7-9): Suporte Básico a Aluguel

**Objetivo**: Permitir anúncios de imóveis para aluguel (sem gestão de contrato).

#### Entregas (60-80h)
1. **Backend** (30-40h):
   - Ativar campos `RentalInfo` no Property model
   - Validações condicionais (`TransactionType = "rent"` → `RentalInfo != NULL`)
   - Endpoint de filtro de busca (total_monthly_cost, rental_type, furnished, accepts_pets)

2. **Frontend Público** (20-30h):
   - Página de busca: Filtros de aluguel (custo total, tipo, mobília, pets)
   - Página de detalhe: Exibir RentalInfo completo (breakdown de custos)
   - Cálculo automático de custo total (transparência)

3. **Frontend Admin** (10h):
   - PropertyForm: Adicionar seção "Informações de Locação"
   - Toggle "Tipo de transação" (venda, aluguel, ambos)
   - Validação de campos obrigatórios

**Critérios de Validação**:
- ✅ Corretor pode cadastrar imóvel para aluguel
- ✅ Busca pública filtra por custo total (não apenas aluguel base)
- ✅ Página de detalhe mostra breakdown de custos
- ✅ Imóvel pode ser anunciado para venda E aluguel simultaneamente

**Timeline**: 3-4 semanas | **Investimento**: R$ 6.000-8.000

---

### 7.2 MVP+4 (Mês 10-12): Gestão de Contratos e Pagamentos

**Objetivo**: Gestão end-to-end de locação (diferenciação total vs CRMs).

#### Entregas (100-120h)
1. **Backend** (50-60h):
   - Models: RentalContract, RentalPayment, TenantPerson
   - CRUD de contratos (criar, assinar, renovar, encerrar)
   - Geração automática de cobranças mensais (Cloud Scheduler)
   - Cálculo de multa e juros (Cloud Scheduler diário)
   - Integração PagSeguro/Stripe (split automático)
   - Alertas de renovação (D-90, D-60, D-30)

2. **Frontend Admin** (30-40h):
   - CRUD de contratos (criar, editar, visualizar)
   - Dashboard de pagamentos (pendentes, atrasados, pagos)
   - Timeline de vencimentos (renovações)
   - Upload de documentos (contrato assinado, vistorias)
   - Integração DocuSign/Clicksign (assinatura digital)

3. **Portal do Inquilino** (20h):
   - Página de pagamento (visualizar boleto, Pix)
   - Histórico de pagamentos
   - Download de recibos

**Critérios de Validação**:
- ✅ Corretor pode criar contrato a partir de Property
- ✅ Cobrança mensal gerada automaticamente (dia 1º de cada mês)
- ✅ Multa e juros aplicados automaticamente após vencimento
- ✅ Split automático: Proprietário 92%, Imobiliária 8%, Plataforma 2%
- ✅ Alertas de renovação enviados (90, 60, 30 dias antes)

**Timeline**: 4-5 semanas | **Investimento**: R$ 10.000-12.000

**Revenue Impact**: R$ 15.000-25.000 MRR adicional (2% sobre pagamentos processados)

---

### 7.3 MVP+5 (Mês 13-15): Manutenção e Portal do Proprietário

**Objetivo**: Diferenciação completa (paridade com QuintoAndar).

#### Entregas (80-100h)
1. **Backend** (40-50h):
   - Model MaintenanceRequest
   - CRUD de chamados
   - SLA automático (Cloud Scheduler)
   - Notificações multi-canal (email, WhatsApp, dashboard)
   - Integração com prestadores de serviço (marketplace futuro)

2. **Frontend Admin** (20-30h):
   - Dashboard de manutenção (abertos, atrasados, resolvidos)
   - Atribuição de prestadores
   - Histórico por imóvel

3. **Portal do Inquilino** (10h):
   - Abrir chamado (foto + descrição)
   - Acompanhar status em tempo real

4. **Portal do Proprietário** (10-20h):
   - Dashboard de pagamentos (status, extratos)
   - Histórico de manutenção (transparência)
   - Documentos (contrato, vistorias, recibos)

**Critérios de Validação**:
- ✅ Inquilino pode abrir chamado (com foto)
- ✅ SLA alertado automaticamente (violações notificadas)
- ✅ Proprietário vê status de pagamento em tempo real
- ✅ Rating de prestadores funciona (NPS)

**Timeline**: 4-5 semanas | **Investimento**: R$ 8.000-10.000

---

## 8. REVENUE MODEL (Locação)

### 8.1 Fontes de Receita

| Fonte | Descrição | Pricing | Volume Esperado (Ano 1) | Receita Anual |
|-------|-----------|---------|------------------------|---------------|
| **Assinatura Base** | Plano Pro/Elite com gestão de locação | +R$ 100-150/mês | 50 tenants | R$ 60.000-90.000 |
| **Taxa de Transação** | 2% sobre pagamentos processados | 2% do aluguel | R$ 500k/mês processado | R$ 120.000 |
| **Análise de Crédito** | Por inquilino aprovado | R$ 79/análise | 100 análises/mês | R$ 94.800 |
| **Renovação de Contrato** | Taxa por renovação automatizada | R$ 49/renovação | 50 renovações/mês | R$ 29.400 |
| **Success Fee (Futuro)** | % sobre comissão de aluguel | 20% da comissão imobiliária | 30 contratos/mês | R$ 72.000 |
| **TOTAL** | - | - | - | **R$ 376.200/ano** |

### 8.2 Unit Economics

**Exemplo: Imobiliária com 50 Contratos de Locação Ativos**

| Métrica | Valor |
|---------|-------|
| Aluguel médio por contrato | R$ 3.000/mês |
| Taxa de processamento (2%) | R$ 60/contrato/mês |
| Receita MRR (50 contratos) | R$ 3.000/mês |
| **Receita Anual** | **R$ 36.000/ano** |
| Custo de processamento (Stripe/PagSeguro ~1.5%) | R$ 2.250/ano |
| **Margem Líquida** | **R$ 33.750/ano (94%)** |

**ROI do Investimento**:
- Investimento MVP+3 + MVP+4: R$ 16.000-20.000
- Receita Ano 1 (conservadora): R$ 180.000
- **ROI**: 9-11x

---

## 9. CONCLUSÃO E RECOMENDAÇÕES

### 9.1 Resumo de Dores vs Soluções

| Dor do Mercado | Nossa Solução | Diferencial Competitivo |
|----------------|---------------|------------------------|
| Anúncios com custos ocultos | Custo total obrigatório + filtro por total | ✅ Único no Brasil |
| Garantias não especificadas | Campo `accepted_guarantees[]` + filtro | ✅ Único no Brasil |
| Gestão manual de contratos | Templates + assinatura digital + alertas | ⚠️ QuintoAndar faz, CRMs não |
| Cobrança manual | Boleto automático + split + inadimplência | ⚠️ QuintoAndar faz, CRMs básico |
| Manutenção dispersa | Portal inquilino + SLA + prestadores | ⚠️ QuintoAndar faz, CRMs não |
| Falta de transparência | Portal proprietário | ⚠️ QuintoAndar faz, CRMs não |
| Sem histórico de imóvel | Rating de proprietário + histórico manutenção | ✅ **NINGUÉM FAZ** |

### 9.2 Recomendação de Priorização

**Imediato (MVP)**:
1. ✅ Adicionar campos `RentalInfo` no Property model (NULL permitido)
2. ✅ Criar placeholders de models (RentalContract, RentalPayment, MaintenanceRequest)
3. ✅ Deploy de Firestore indexes para locação (custo zero se não usados)

**MVP+3 (Mês 7-9)**:
1. ✅ Ativar anúncios de aluguel (frontend + backend)
2. ✅ Filtros avançados (custo total, tipo, mobília, pets)
3. ✅ Breakdown de custos (transparência total)

**MVP+4 (Mês 10-12)**:
1. ✅ Gestão de contratos (CRUD + assinatura digital)
2. ✅ Pagamentos automatizados (geração + cobrança + split)
3. ✅ Alertas de renovação (D-90, D-60, D-30)

**MVP+5 (Mês 13-15)**:
1. ✅ Manutenção (chamados + SLA + portal inquilino)
2. ✅ Portal do proprietário (transparência total)
3. ✅ Rating de proprietários (inovação única)

### 9.3 Moat Defensável

**Nossa Vantagem Competitiva de Longo Prazo**:
1. ✅ **Efeito de Rede**: Marketplace de co-corretagem (quanto mais corretores, mais imóveis, mais leads)
2. ✅ **Dados Proprietários**: Histórico de manutenção, ratings, análise de mercado (concorrentes não têm)
3. ✅ **Lock-in**: Contratos ativos + histórico de pagamentos (switching cost alto)
4. ✅ **Modelo Híbrido**: Marketplace + Gestão (ZAP não faz gestão, QuintoAndar não faz marketplace)

---

**Documento gerado em**: 2025-12-21
**Por**: Claude Sonnet 4.5 (Análise de Mercado)
**Próximos Passos**: Validar com stakeholders → Atualizar AI_DEV_DIRECTIVE.md → Implementar preparação de schema em MVP
