#!/bin/bash
# Integration Test Script - Creates test data for frontend

API_URL="http://localhost:8080"

echo "========================================="
echo "Creating Test Data for Frontend"
echo "========================================="

# 1. Create Tenant
echo -e "\n[1/6] Creating tenant..."
TENANT_RESPONSE=$(curl -s -X POST "$API_URL/tenants" \
  -H "Content-Type: application/json" \
  -d '{"name":"Imobiliária Demo","slug":"demo","is_active":true}')

TENANT_ID=$(echo $TENANT_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "✓ Tenant ID: $TENANT_ID"

# 2. Create Owner
echo -e "\n[2/6] Creating owner..."
OWNER_RESPONSE=$(curl -s -X POST "$API_URL/api/$TENANT_ID/owners" \
  -H "Content-Type: application/json" \
  -d '{"name":"João Silva","email":"joao@teste.com","phone":"11987654321","cpf":"12345678901","consent_given":true,"consent_text":"Aceito"}')

OWNER_ID=$(echo $OWNER_RESPONSE | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "✓ Owner ID: $OWNER_ID"

# 3. Create Properties
echo -e "\n[3/6] Creating featured apartment..."
curl -s -X POST "$API_URL/api/$TENANT_ID/properties" \
  -H "Content-Type: application/json" \
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
    \"bedrooms\":3,
    \"bathrooms\":2,
    \"suites\":1,
    \"parking_spaces\":2,
    \"area_sqm\":120,
    \"has_pool\":true,
    \"has_gym\":true,
    \"has_elevator\":true,
    \"pet_friendly\":true,
    \"title\":\"Apartamento Moderno no Jardins\",
    \"description\":\"Lindo apartamento com acabamento de primeira.\",
    \"featured\":true,
    \"slug\":\"apt-jardins-sp\"
  }" > /dev/null

echo "✓ Property 1 created"

echo -e "\n[4/6] Creating house for rent..."
curl -s -X POST "$API_URL/api/$TENANT_ID/properties" \
  -H "Content-Type: application/json" \
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
    \"bedrooms\":4,
    \"bathrooms\":3,
    \"parking_spaces\":3,
    \"area_sqm\":250,
    \"has_garden\":true,
    \"pet_friendly\":true,
    \"title\":\"Casa Espaçosa em Campinas\",
    \"description\":\"Casa ampla ideal para famílias.\",
    \"featured\":true,
    \"slug\":\"casa-campinas-sp\"
  }" > /dev/null

echo "✓ Property 2 created"

# 4. List Properties
echo -e "\n[5/6] Listing properties..."
curl -s "$API_URL/api/$TENANT_ID/properties?status=available" | grep -o '"title":"[^"]*' | cut -d'"' -f4

# 5. Summary
echo -e "\n========================================="
echo "Test Data Created Successfully!"
echo "========================================="
echo "Tenant ID: $TENANT_ID"
echo ""
echo "Update frontend-public/.env.local:"
echo "NEXT_PUBLIC_TENANT_ID=$TENANT_ID"
echo "========================================="
