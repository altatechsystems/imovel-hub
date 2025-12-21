# Claude Code Skills - Ecosistema Imob

Este projeto possui **2 skills especializadas** para garantir implementação de código de alta qualidade seguindo as melhores práticas de cada stack.

## 📋 Skills Disponíveis

### 1. `golang-dev` - Backend Go

**Arquivo**: `golang-dev.skill` (691 linhas)

**Quando é ativada**: Automaticamente quando trabalhando com código Go, APIs REST, lógica de negócio backend.

**Especialidades**:
- ✅ Go idiomático seguindo Effective Go
- ✅ Naming conventions (mixedCaps, package names)
- ✅ Error handling (`fmt.Errorf` com `%w`)
- ✅ Concurrency (goroutines, channels, context)
- ✅ Standard library preference (io, net/http, encoding/json, sync)
- ✅ Interface design (small, focused)
- ✅ Memory management (make, defer, slices)
- ✅ Table-driven tests
- ✅ Project structure (cmd/, internal/, pkg/)

**Tecnologias cobertas**:
- Gin framework (HTTP handlers)
- Firestore (context, structs)
- Google Cloud Vision API (async, Pub/Sub)
- ffmpeg/ffprobe (os/exec)
- Cloud Functions (HTTP, context)
- Multi-tenancy (dependency injection)
- Google Cloud Storage (io.Reader/Writer)

**Exemplo de código gerado**:
```go
func (v *VisionAnalyzer) AnalyzePhoto(ctx context.Context, photoURL string) (*PhotoAnalysis, error) {
    client, err := vision.NewImageAnnotatorClient(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to create vision client: %w", err)
    }
    defer client.Close()

    // ...
}
```

---

### 2. `nextjs-dev` - Frontend Next.js + React

**Arquivo**: `nextjs-dev.skill` (criado agora - 725 linhas)

**Quando é ativada**: Automaticamente quando trabalhando com componentes React, páginas Next.js, TypeScript frontend.

**Especialidades**:
- ✅ Next.js 14+ App Router
- ✅ Server Components (default) vs Client Components
- ✅ TypeScript strict mode
- ✅ shadcn/ui component library
- ✅ React Query (client-side caching)
- ✅ Zod validation
- ✅ @dnd-kit (drag & drop)
- ✅ react-dropzone (file upload)
- ✅ Tailwind CSS
- ✅ SEO optimization (metadata, OpenGraph)
- ✅ Accessibility (ARIA, semantic HTML)

**Tecnologias cobertas**:
- **PROMPT 04 (Public)**: Property listings, search, LGPD forms
- **PROMPT 04b (Admin)**: Photo/video upload, drag-and-drop reordering, quality indicators

**Exemplo de código gerado**:
```typescript
// app/imoveis/[id]/page.tsx
export default async function PropertyPage({ params }: PageProps) {
  const property = await fetchProperty(params.id)

  return (
    <div>
      <h1>{property.title}</h1>
      <PhotoGallery photos={property.photos} />
    </div>
  )
}

export async function generateMetadata({ params }: PageProps): Promise<Metadata> {
  const property = await fetchProperty(params.id)

  return {
    title: `${property.title} - Ecosistema Imob`,
    description: property.description,
  }
}
```

---

## 🔧 Comandos Pré-Aprovados

O arquivo [`.claude/settings.local.json`](../settings.local.json) possui **266 comandos pré-aprovados**, incluindo:

### Go Commands
```bash
go test ./...
go build ./cmd/myapp
go run ./cmd/myapp
go get package/path
go mod tidy
```

### Node.js Commands
```bash
npm install
npm run dev
npm run build
npx create-next-app@latest
```

### AWS Commands
```bash
aws dynamodb describe-table
aws dynamodb scan
aws secretsmanager list-secrets
aws ecr get-login-password
```

Isso significa que durante a implementação, **não será necessário pedir permissão** para executar estes comandos.

---

## ✅ Verificação de Compatibilidade

### Backend (PROMPTs 01-03)

| Tecnologia | Skill | Cobertura |
|-----------|-------|----------|
| **Gin framework** | `golang-dev` | ✅ HTTP handlers |
| **Firestore** | `golang-dev` | ✅ Context, structs |
| **Vision API** | `golang-dev` | ✅ Async patterns |
| **ffmpeg** | `golang-dev` | ✅ os/exec |
| **Pub/Sub** | `golang-dev` | ✅ Goroutines, channels |
| **Multi-tenancy** | `golang-dev` | ✅ Dependency injection |

### Frontend (PROMPTs 04, 04b)

| Tecnologia | Skill | Cobertura |
|-----------|-------|----------|
| **Next.js 14+ App Router** | `nextjs-dev` | ✅ File-based routing |
| **shadcn/ui** | `nextjs-dev` | ✅ Component patterns |
| **React Query** | `nextjs-dev` | ✅ Client caching |
| **@dnd-kit** | `nextjs-dev` | ✅ Drag & drop |
| **TypeScript** | `nextjs-dev` | ✅ Strict mode, Zod |
| **Tailwind CSS** | `nextjs-dev` | ✅ Utility classes |
| **LGPD compliance** | `nextjs-dev` | ✅ Form validation, privacy |

---

## 📊 Estatísticas

- **Total de linhas de diretrizes**: 1416 linhas (691 Go + 725 Next.js)
- **Comandos pré-aprovados**: 266
- **Cobertura de tecnologias**: 100% (backend + frontend)
- **Padrões documentados**: Architecture patterns, testing, performance, accessibility

---

## 🚀 Como Usar

### Durante Implementação Backend
As skills Go são ativadas **automaticamente** quando:
1. Modificando arquivos `.go`
2. Criando packages Go
3. Escrevendo testes `*_test.go`
4. Trabalhando com Firestore, Vision API, Cloud Functions

### Durante Implementação Frontend
As skills Next.js são ativadas **automaticamente** quando:
1. Modificando arquivos `.tsx` ou `.ts` no frontend
2. Criando páginas (`page.tsx`)
3. Desenvolvendo componentes React
4. Implementando formulários, validação, drag-and-drop

### Verificar Skills Ativas
Para confirmar que as skills estão carregadas:
```bash
# Claude Code carrega skills automaticamente da pasta .claude/skills/
# Basta iniciar a implementação!
```

---

## 📝 Notas de Implementação

### Regras Críticas do Go Skill
1. **SEMPRE** usar `fmt.Errorf` com `%w` para wrapping de erros
2. **NUNCA** usar underscores em nomes de variáveis/funções (usar mixedCaps)
3. **PREFERIR** Go standard library ao invés de pacotes externos
4. **SEMPRE** documentar funções/tipos exportados
5. **USAR** table-driven tests para casos múltiplos

### Regras Críticas do Next.js Skill
1. **Server Components by default** - Só adicionar `'use client'` quando necessário
2. **SEMPRE** usar TypeScript strict mode (sem `any`)
3. **VALIDAR** com Zod em boundaries (forms, API)
4. **OTIMIZAR** imagens com `next/image`
5. **ACESSIBILIDADE** - Semantic HTML, ARIA labels, keyboard navigation

---

## 🎯 Conclusão

Este projeto está **100% preparado** para implementação com:
- ✅ **2 skills especializadas** cobrindo todo o stack
- ✅ **266 comandos pré-aprovados** para workflow fluido
- ✅ **1416 linhas de best practices** documentadas
- ✅ **Compatibilidade total** com todas as tecnologias nos PROMPTs

**Claude Code respeitará todas estas diretrizes durante a implementação!**
