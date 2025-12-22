@echo off
cd backend
set FIREBASE_PROJECT_ID=ecosistema-imob-dev
set PORT=8080
set GOOGLE_APPLICATION_CREDENTIALS=./config/firebase-adminsdk.json
set GIN_MODE=debug
set ENVIRONMENT=development
..\backend\bin\caas.exe
