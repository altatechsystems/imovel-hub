# Script para atualizar metadados dos Tenants Master
# Requer: Firebase CLI instalado e autenticado

$TENANT_MASTER_ID = "391b12f8-ebe4-426a-8c99-ec5a10b1f361"
$TENANT_ZERO_ID = "bd71c02b-5fa5-43df-8b46-a1df2206f1ef"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Atualizando Metadados dos Tenants Master" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""

# 1. Atualizar Tenant Master (ALTATECH Systems)
Write-Host "1. Atualizando TENANT MASTER - ALTATECH Systems..." -ForegroundColor Yellow
Write-Host "   Tenant ID: $TENANT_MASTER_ID" -ForegroundColor Gray

$tenantMasterUpdate = @{
    is_platform_admin = $true
    document = "36.077.869/0001-81"
    document_type = "cnpj"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/tenants/$TENANT_MASTER_ID" `
        -Method PUT `
        -Body $tenantMasterUpdate `
        -ContentType "application/json"

    Write-Host "   Metadados do Tenant Master atualizados!" -ForegroundColor Green
} catch {
    Write-Host "   Erro ao atualizar Tenant Master: $_" -ForegroundColor Red
}

Write-Host ""

# 2. Atualizar Tenant Zero (ALTATECH Imoveis)
Write-Host "2. Atualizando TENANT ZERO - ALTATECH Imoveis..." -ForegroundColor Yellow
Write-Host "   Tenant ID: $TENANT_ZERO_ID" -ForegroundColor Gray

$tenantZeroUpdate = @{
    is_default_tenant = $true
    document = "26.517.873/0001-60"
    document_type = "cnpj"
    creci = "05733-J/SP"
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/tenants/$TENANT_ZERO_ID" `
        -Method PUT `
        -Body $tenantZeroUpdate `
        -ContentType "application/json"

    Write-Host "   Metadados do Tenant Zero atualizados!" -ForegroundColor Green
} catch {
    Write-Host "   Erro ao atualizar Tenant Zero: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "Metadados atualizados com sucesso!" -ForegroundColor Green
Write-Host "=========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "PROXIMOS PASSOS:" -ForegroundColor Yellow
Write-Host "1. Executar script de vinculacao de imoveis ao Tenant Zero" -ForegroundColor White
Write-Host "2. Configurar custom claims para platform_admin no Firebase Console" -ForegroundColor White
Write-Host ""
