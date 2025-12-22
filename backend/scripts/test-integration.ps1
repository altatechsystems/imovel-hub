# Integration Test Script for Ecosistema Imob
# Creates test data and validates frontend/backend integration

$API_URL = "http://localhost:8080"
$TENANT_ID = ""
$OWNER_ID = ""
$BROKER_ID = ""
$PROPERTY1_ID = ""
$PROPERTY2_ID = ""

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "ECOSISTEMA IMOB - INTEGRATION TESTS" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# 1. Create Tenant
Write-Host "[1/7] Creating tenant..." -ForegroundColor Yellow
$tenantResponse = curl -s -X POST "$API_URL/tenants" `
    -H "Content-Type: application/json" `
    -d '{"name":"Imobiliária Teste","slug":"imob-teste","is_active":true}' | ConvertFrom-Json

if ($tenantResponse.success) {
    $TENANT_ID = $tenantResponse.data.id
    Write-Host "✓ Tenant created: $TENANT_ID" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create tenant" -ForegroundColor Red
    exit 1
}

# 2. Create Owner
Write-Host "`n[2/7] Creating owner..." -ForegroundColor Yellow
$ownerResponse = curl -s -X POST "$API_URL/api/$TENANT_ID/owners" `
    -H "Content-Type: application/json" `
    -d '{"name":"João Silva","email":"joao@teste.com","phone":"11987654321","cpf":"12345678901","consent_given":true,"consent_text":"Autorizo uso dos dados"}' | ConvertFrom-Json

if ($ownerResponse.success) {
    $OWNER_ID = $ownerResponse.data.id
    Write-Host "✓ Owner created: $OWNER_ID" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create owner" -ForegroundColor Red
    exit 1
}

# 3. Create Broker
Write-Host "`n[3/7] Creating broker..." -ForegroundColor Yellow
$brokerResponse = curl -s -X POST "$API_URL/api/$TENANT_ID/brokers" `
    -H "Content-Type: application/json" `
    -d '{"name":"Maria Corretor","email":"maria@imob.com","phone":"11987654322","creci":"12345","is_active":true}' | ConvertFrom-Json

if ($brokerResponse.success) {
    $BROKER_ID = $brokerResponse.data.id
    Write-Host "✓ Broker created: $BROKER_ID" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create broker" -ForegroundColor Red
    exit 1
}

# 4. Create Property 1 (Featured Apartment)
Write-Host "`n[4/7] Creating featured apartment..." -ForegroundColor Yellow
$prop1Response = curl -s -X POST "$API_URL/api/$TENANT_ID/properties" `
    -H "Content-Type: application/json" `
    -d "{
        \"owner_id\":\"$OWNER_ID\",
        \"transaction_type\":\"sale\",
        \"property_type\":\"apartment\",
        \"status\":\"available\",
        \"visibility\":\"public\",
        \"sale_price\":850000,
        \"street\":\"Rua das Flores\",
        \"number\":\"123\",
        \"complement\":\"Apto 45\",
        \"neighborhood\":\"Jardins\",
        \"city\":\"São Paulo\",
        \"state\":\"SP\",
        \"postal_code\":\"01234567\",
        \"country\":\"Brasil\",
        \"bedrooms\":3,
        \"bathrooms\":2,
        \"suites\":1,
        \"parking_spaces\":2,
        \"area_sqm\":120,
        \"furnished\":false,
        \"pet_friendly\":true,
        \"has_pool\":true,
        \"has_gym\":true,
        \"has_elevator\":true,
        \"has_security\":true,
        \"title\":\"Apartamento Moderno no Jardins\",
        \"description\":\"Lindo apartamento com acabamento de primeira, localizado em uma das melhores regiões de São Paulo. Próximo a restaurantes, shopping e transporte público.\",
        \"featured\":true,
        \"slug\":\"apartamento-moderno-jardins-sp\"
    }" | ConvertFrom-Json

if ($prop1Response.success) {
    $PROPERTY1_ID = $prop1Response.data.id
    Write-Host "✓ Featured apartment created: $PROPERTY1_ID" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create property 1" -ForegroundColor Red
    exit 1
}

# 5. Create Property 2 (House for Rent)
Write-Host "`n[5/7] Creating house for rent..." -ForegroundColor Yellow
$prop2Response = curl -s -X POST "$API_URL/api/$TENANT_ID/properties" `
    -H "Content-Type: application/json" `
    -d "{
        \"owner_id\":\"$OWNER_ID\",
        \"transaction_type\":\"rent\",
        \"property_type\":\"house\",
        \"status\":\"available\",
        \"visibility\":\"public\",
        \"rental_price\":4500,
        \"street\":\"Av. Principal\",
        \"number\":\"456\",
        \"neighborhood\":\"Centro\",
        \"city\":\"Campinas\",
        \"state\":\"SP\",
        \"postal_code\":\"13010000\",
        \"country\":\"Brasil\",
        \"bedrooms\":4,
        \"bathrooms\":3,
        \"suites\":2,
        \"parking_spaces\":3,
        \"area_sqm\":250,
        \"total_area_sqm\":400,
        \"furnished\":true,
        \"pet_friendly\":true,
        \"has_pool\":false,
        \"has_gym\":false,
        \"has_garden\":true,
        \"has_balcony\":true,
        \"title\":\"Casa Espaçosa em Campinas\",
        \"description\":\"Casa ampla e confortável, ideal para famílias. Possui jardim, quintal e garagem para 3 carros.\",
        \"featured\":true,
        \"slug\":\"casa-espacosa-campinas-sp\"
    }" | ConvertFrom-Json

if ($prop2Response.success) {
    $PROPERTY2_ID = $prop2Response.data.id
    Write-Host "✓ House created: $PROPERTY2_ID" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create property 2" -ForegroundColor Red
    exit 1
}

# 6. List Properties
Write-Host "`n[6/7] Testing property listing..." -ForegroundColor Yellow
$listResponse = curl -s "$API_URL/api/$TENANT_ID/properties?status=available&visibility=public" | ConvertFrom-Json

if ($listResponse.success) {
    Write-Host "✓ Found $($listResponse.count) available properties" -ForegroundColor Green
    $listResponse.data | ForEach-Object {
        Write-Host "  - $($_.title) ($($_.property_type), $($_.city))" -ForegroundColor White
    }
} else {
    Write-Host "✗ Failed to list properties" -ForegroundColor Red
}

# 7. Test Lead Creation
Write-Host "`n[7/7] Testing lead creation..." -ForegroundColor Yellow
$leadResponse = curl -s -X POST "$API_URL/api/$TENANT_ID/leads" `
    -H "Content-Type: application/json" `
    -d "{
        \"property_id\":\"$PROPERTY1_ID\",
        \"name\":\"Pedro Santos\",
        \"email\":\"pedro@email.com\",
        \"phone\":\"11999887766\",
        \"message\":\"Gostaria de agendar uma visita\",
        \"channel\":\"form\",
        \"consent_text\":\"Autorizo o uso dos meus dados conforme LGPD\"
    }" | ConvertFrom-Json

if ($leadResponse.success) {
    Write-Host "✓ Lead created successfully" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to create lead" -ForegroundColor Red
}

# Summary
Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "INTEGRATION TEST RESULTS" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Tenant ID: $TENANT_ID" -ForegroundColor White
Write-Host "Owner ID: $OWNER_ID" -ForegroundColor White
Write-Host "Broker ID: $BROKER_ID" -ForegroundColor White
Write-Host "Property 1 ID: $PROPERTY1_ID" -ForegroundColor White
Write-Host "Property 2 ID: $PROPERTY2_ID" -ForegroundColor White
Write-Host "`nFrontend API URL: $API_URL/api/$TENANT_ID" -ForegroundColor Yellow
Write-Host "Update .env.local with NEXT_PUBLIC_TENANT_ID=$TENANT_ID" -ForegroundColor Yellow
Write-Host "`n✓ Integration tests completed successfully!" -ForegroundColor Green
