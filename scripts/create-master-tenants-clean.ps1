# Script para criar os tenants master da plataforma
# Uso: powershell -ExecutionPolicy Bypass -File scripts/create-master-tenants.ps1

$ErrorActionPreference = "Stop"

# Configurao da API
$API_URL = "http://localhost:8080/api/v1"

Write-Host "=====================================" -ForegroundColor Cyan
Write-Host "Criando Tenants Master da Plataforma" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

# ========================================
# 1. TENANT MASTER - ALTATECH SYSTEMS
# ========================================
Write-Host "1. Criando Tenant Master: ALTATECH Systems..." -ForegroundColor Yellow

$masterTenantData = @{
    email = "daniel.garcia@altatechsystems.com"
    password = "`$%&AltatechSystems`$%&"
    name = "Daniel Garcia"
    phone = "+5511999999999"  # Voc pode ajustar para seu telefone real
    tenant_name = "ALTATECH Systems"
} | ConvertTo-Json

try {
    $masterResponse = Invoke-RestMethod -Uri "$API_URL/auth/signup" `
        -Method POST `
        -Body $masterTenantData `
        -ContentType "application/json" `
        -ErrorAction Stop

    Write-Host " Tenant Master criado com sucesso!" -ForegroundColor Green
    Write-Host "   Tenant ID: $($masterResponse.tenant_id)" -ForegroundColor Gray
    Write-Host "   Broker ID: $($masterResponse.broker_id)" -ForegroundColor Gray
    Write-Host "   Email: daniel.garcia@altatechsystems.com" -ForegroundColor Gray
    Write-Host ""

    $masterTenantId = $masterResponse.tenant_id
    $masterBrokerId = $masterResponse.broker_id
    $masterToken = $masterResponse.firebase_token

} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 409) {
        Write-Host "  Tenant Master j existe (email j cadastrado)" -ForegroundColor Yellow
        Write-Host "   Tentando fazer login..." -ForegroundColor Gray

        $loginData = @{
            email = "daniel.garcia@altatechsystems.com"
            password = "`$%&AltatechSystems`$%&"
        } | ConvertTo-Json

        try {
            $loginResponse = Invoke-RestMethod -Uri "$API_URL/auth/login" `
                -Method POST `
                -Body $loginData `
                -ContentType "application/json"

            Write-Host " Login realizado com sucesso!" -ForegroundColor Green
            Write-Host "   Tenant ID: $($loginResponse.tenant_id)" -ForegroundColor Gray
            Write-Host ""

            $masterTenantId = $loginResponse.tenant_id
            $masterToken = $loginResponse.firebase_token

        } catch {
            Write-Host " Erro ao fazer login do Tenant Master: $_" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host " Erro ao criar Tenant Master: $_" -ForegroundColor Red
        exit 1
    }
}

# ========================================
# 2. TENANT ZERO - ALTATECH IMVEIS
# ========================================
Write-Host "2. Criando Tenant Zero: ALTATECH Imveis..." -ForegroundColor Yellow

$tenantZeroData = @{
    email = "administracao@altatechimoveis.com"
    password = "`$%&AltatechImoveis`$%&"
    name = "Administrao"
    phone = "+5511988888888"  # Voc pode ajustar para o telefone da imobiliria
    tenant_name = "ALTATECH Imveis"
} | ConvertTo-Json

try {
    $tenantZeroResponse = Invoke-RestMethod -Uri "$API_URL/auth/signup" `
        -Method POST `
        -Body $tenantZeroData `
        -ContentType "application/json" `
        -ErrorAction Stop

    Write-Host " Tenant Zero criado com sucesso!" -ForegroundColor Green
    Write-Host "   Tenant ID: $($tenantZeroResponse.tenant_id)" -ForegroundColor Gray
    Write-Host "   Broker ID: $($tenantZeroResponse.broker_id)" -ForegroundColor Gray
    Write-Host "   Email: administracao@altatechimoveis.com" -ForegroundColor Gray
    Write-Host ""

    $tenantZeroId = $tenantZeroResponse.tenant_id
    $tenantZeroBrokerId = $tenantZeroResponse.broker_id
    $tenantZeroToken = $tenantZeroResponse.firebase_token

} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 409) {
        Write-Host "  Tenant Zero j existe (email j cadastrado)" -ForegroundColor Yellow
        Write-Host "   Tentando fazer login..." -ForegroundColor Gray

        $loginData = @{
            email = "administracao@altatechimoveis.com"
            password = "`$%&AltatechImoveis`$%&"
        } | ConvertTo-Json

        try {
            $loginResponse = Invoke-RestMethod -Uri "$API_URL/auth/login" `
                -Method POST `
                -Body $loginData `
                -ContentType "application/json"

            Write-Host " Login realizado com sucesso!" -ForegroundColor Green
            Write-Host "   Tenant ID: $($loginResponse.tenant_id)" -ForegroundColor Gray
            Write-Host ""

            $tenantZeroId = $loginResponse.tenant_id
            $tenantZeroToken = $loginResponse.firefox_token

        } catch {
            Write-Host " Erro ao fazer login do Tenant Zero: $_" -ForegroundColor Red
            exit 1
        }
    } else {
        Write-Host " Erro ao criar Tenant Zero: $_" -ForegroundColor Red
        exit 1
    }
}

# ========================================
# 3. ATUALIZAR METADADOS VIA FIRESTORE
# ========================================
Write-Host "3. Atualizando metadados dos tenants..." -ForegroundColor Yellow
Write-Host ""

Write-Host "  ATENO: Metadados adicionais precisam ser configurados manualmente:" -ForegroundColor Yellow
Write-Host ""
Write-Host " TENANT MASTER (ALTATECH Systems):" -ForegroundColor Cyan
Write-Host "   Tenant ID: $masterTenantId" -ForegroundColor Gray
Write-Host "   Adicionar no Firestore:" -ForegroundColor Gray
Write-Host "   - is_platform_admin: true" -ForegroundColor White
Write-Host "   - cnpj: '36.077.869/0001-81'" -ForegroundColor White
Write-Host "   - settings.business_name: 'ALTATECH Systems'" -ForegroundColor White
Write-Host ""

Write-Host " TENANT ZERO (ALTATECH Imveis):" -ForegroundColor Cyan
Write-Host "   Tenant ID: $tenantZeroId" -ForegroundColor Gray
Write-Host "   Adicionar no Firestore:" -ForegroundColor Gray
Write-Host "   - is_default_tenant: true" -ForegroundColor White
Write-Host "   - cnpj: '26.517.873/0001-60'" -ForegroundColor White
Write-Host "   - creci: '5733-J'" -ForegroundColor White
Write-Host "   - settings.business_name: 'ALTATECH Imveis'" -ForegroundColor White
Write-Host ""

Write-Host " VINCULAR IMVEIS EXISTENTES:" -ForegroundColor Cyan
Write-Host "   Executar script para migrar imveis sem tenant_id para:" -ForegroundColor Gray
Write-Host "   Tenant ID: $tenantZeroId" -ForegroundColor White
Write-Host ""

# ========================================
# 4. RESUMO FINAL
# ========================================
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host " RESUMO DA CRIAO" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""

Write-Host " TENANT MASTER - ALTATECH Systems" -ForegroundColor Green
Write-Host "   Email: daniel.garcia@altatechsystems.com" -ForegroundColor Gray
Write-Host "   Senha: `$%&AltatechSystems`$%&" -ForegroundColor Gray
Write-Host "   Tenant ID: $masterTenantId" -ForegroundColor Gray
Write-Host "   Login Admin: http://localhost:3002/login" -ForegroundColor Blue
Write-Host ""

Write-Host " TENANT ZERO - ALTATECH Imveis" -ForegroundColor Green
Write-Host "   Email: administracao@altatechimoveis.com" -ForegroundColor Gray
Write-Host "   Senha: `$%&AltatechImoveis`$%&" -ForegroundColor Gray
Write-Host "   Tenant ID: $tenantZeroId" -ForegroundColor Gray
Write-Host "   Login Admin: http://localhost:3002/login" -ForegroundColor Blue
Write-Host ""

Write-Host "PROXIMOS PASSOS:" -ForegroundColor Yellow
Write-Host "1. Atualizar metadados no Firebase Console" -ForegroundColor White
Write-Host "2. Executar script de vinculacao de imoveis" -ForegroundColor White
Write-Host "3. Configurar custom claims para platform_admin" -ForegroundColor White
Write-Host ""
