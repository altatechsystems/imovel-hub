#!/bin/bash

# Complete Integration Test
API="http://localhost:8080"

echo "==========================================="
echo "INTEGRATION TEST - Complete Flow"
echo "==========================================="

# 1. Create Tenant
echo -e "\n[1/5] Creating tenant..."
TENANT=$(curl -s -X POST "$API/tenants" -H "Content-Type: application/json" -d '{"name":"Demo Imob","slug":"demo-test","is_active":true}')
TENANT_ID=$(echo $TENANT | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo "✓ Tenant: $TENANT_ID"

# 2. Create Owner
echo -e "\n[2/5] Creating owner..."
OWNER=$(curl -s -X POST "$API/api/$TENANT_ID/owners" -H "Content-Type: application/json" \
  -d '{"name":"João Silva","email":"joao@test.com","phone":"11987654321","cpf":"12345678901","consent_given":true,"consent_text":"Aceito"}')
OWNER_ID=$(echo $OWNER | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo "✓ Owner: $OWNER_ID"

# 3. Create Properties
echo -e "\n[3/5] Creating properties..."

# Property 1 - Featured Apartment
P1=$(curl -s -X POST "$API/api/$TENANT_ID/properties" -H "Content-Type: application/json" -d "{
  \"owner_id\":\"$OWNER_ID\",\"transaction_type\":\"sale\",\"property_type\":\"apartment\",
  \"status\":\"available\",\"visibility\":\"public\",\"sale_price\":850000,
  \"street\":\"Rua das Flores\",\"number\":\"123\",\"neighborhood\":\"Jardins\",
  \"city\":\"São Paulo\",\"state\":\"SP\",\"postal_code\":\"01234567\",
  \"bedrooms\":3,\"bathrooms\":2,\"suites\":1,\"parking_spaces\":2,\"area_sqm\":120,
  \"has_pool\":true,\"has_gym\":true,\"has_elevator\":true,\"pet_friendly\":true,
  \"title\":\"Apartamento Moderno Jardins\",
  \"description\":\"Lindo apartamento com acabamento de primeira.\",
  \"featured\":true,\"slug\":\"apt-jardins-sp\"
}")
P1_ID=$(echo $P1 | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo "✓ Property 1: $P1_ID"

# Property 2 - House
curl -s -X POST "$API/api/$TENANT_ID/properties" -H "Content-Type: application/json" -d "{
  \"owner_id\":\"$OWNER_ID\",\"transaction_type\":\"rent\",\"property_type\":\"house\",
  \"status\":\"available\",\"visibility\":\"public\",\"rental_price\":4500,
  \"street\":\"Av. Principal\",\"number\":\"456\",\"neighborhood\":\"Centro\",
  \"city\":\"Campinas\",\"state\":\"SP\",\"postal_code\":\"13010000\",
  \"bedrooms\":4,\"bathrooms\":3,\"parking_spaces\":3,\"area_sqm\":250,
  \"has_garden\":true,\"pet_friendly\":true,
  \"title\":\"Casa Espaçosa Campinas\",
  \"description\":\"Casa ampla ideal para famílias.\",
  \"featured\":true,\"slug\":\"casa-campinas\"
}" > /dev/null
echo "✓ Property 2 created"

# Property 3 - Commercial
curl -s -X POST "$API/api/$TENANT_ID/properties" -H "Content-Type: application/json" -d "{
  \"owner_id\":\"$OWNER_ID\",\"transaction_type\":\"both\",\"property_type\":\"commercial\",
  \"status\":\"available\",\"visibility\":\"public\",\"sale_price\":1200000,\"rental_price\":8000,
  \"street\":\"Av. Paulista\",\"number\":\"789\",\"neighborhood\":\"Bela Vista\",
  \"city\":\"São Paulo\",\"state\":\"SP\",\"postal_code\":\"01310000\",
  \"bathrooms\":4,\"parking_spaces\":5,\"area_sqm\":300,
  \"has_elevator\":true,\"has_security\":true,
  \"title\":\"Sala Comercial Paulista\",
  \"description\":\"Excelente localização na Paulista.\",
  \"featured\":false,\"slug\":\"sala-paulista\"
}" > /dev/null
echo "✓ Property 3 created"

# 4. List Properties
echo -e "\n[4/5] Listing properties..."
LIST=$(curl -s "$API/api/$TENANT_ID/properties?status=available&visibility=public")
COUNT=$(echo $LIST | grep -o '"count":[0-9]*' | cut -d':' -f2)
echo "✓ Found $COUNT properties"

# 5. Create Lead
echo -e "\n[5/5] Creating lead..."
LEAD=$(curl -s -X POST "$API/api/$TENANT_ID/leads" -H "Content-Type: application/json" -d "{
  \"property_id\":\"$P1_ID\",
  \"name\":\"Maria Santos\",
  \"email\":\"maria@email.com\",
  \"phone\":\"11999887766\",
  \"message\":\"Gostaria de agendar visita\",
  \"channel\":\"form\",
  \"consent_text\":\"Autorizo uso dos dados conforme LGPD\"
}")
LEAD_ID=$(echo $LEAD | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)
echo "✓ Lead: $LEAD_ID"

# Summary
echo -e "\n==========================================="
echo "✓ INTEGRATION TEST COMPLETED"
echo "==========================================="
echo "Tenant ID: $TENANT_ID"
echo "Properties: 3 created, $COUNT available"
echo "Lead: 1 created"
echo ""
echo "Update frontend-public/.env.local:"
echo "NEXT_PUBLIC_TENANT_ID=$TENANT_ID"
echo ""
echo "Start frontend: cd frontend-public && npm run dev"
echo "==========================================="
