# Script para iniciar o backend com variáveis de ambiente corretas
$ErrorActionPreference = "Stop"

Write-Host "Iniciando Backend..." -ForegroundColor Cyan

# Setar variáveis de ambiente
$env:FIREBASE_PROJECT_ID = "ecosistema-imob-dev"
$env:PORT = "8080"
$env:GOOGLE_APPLICATION_CREDENTIALS = "./config/firebase-adminsdk.json"
$env:GIN_MODE = "debug"
$env:ENVIRONMENT = "development"

# Ir para o diretório backend
Set-Location backend

# Iniciar backend
Write-Host "Iniciando em http://localhost:8080" -ForegroundColor Green
& "..\backend\bin\caas.exe"
