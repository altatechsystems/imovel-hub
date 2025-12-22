# Firebase API Key Issue - Troubleshooting Guide

## Problem
The Firebase API key `AIzaSyDx6PRtyNnAgTxsT3coow9Ut0xM4Xi5A71` is returning error:
```
Firebase: Error (auth/api-key-not-valid-please-pass-a-valid-api-key)
```

## Verification
Confirmed that the API key in `.env.local` matches the one shown in Firebase Console.

## Possible Causes

### 1. API Key Restrictions
The API key might have restrictions that prevent it from being used with localhost.

**Solution:**
1. Go to: https://console.cloud.google.com/apis/credentials?project=ecosistema-imob-dev
2. Find the API key `AIzaSyDx6PRtyNnAgTxsT3coow9Ut0xM4Xi5A71`
3. Click on it to edit
4. Check "Application restrictions":
   - Should be "None" OR
   - Should include `localhost:3002` in the allowed domains
5. Check "API restrictions":
   - Should be "Don't restrict key" OR
   - Should include "Identity Toolkit API"

### 2. Identity Toolkit API Not Enabled
The Identity Toolkit API (Firebase Authentication) might not be enabled for the project.

**Solution:**
1. Go to: https://console.cloud.google.com/apis/library/identitytoolkit.googleapis.com?project=ecosistema-imob-dev
2. Click "Enable" if not already enabled

### 3. Email/Password Authentication Not Enabled
The authentication method might not be configured in Firebase.

**Solution:**
1. Go to: https://console.firebase.google.com/project/ecosistema-imob-dev/authentication/providers
2. Enable "Email/Password" sign-in method
3. Make sure it shows as "Enabled"

### 4. Authorized Domains
The domain `localhost` might not be in the authorized domains list.

**Solution:**
1. Go to: https://console.firebase.google.com/project/ecosistema-imob-dev/authentication/settings
2. Scroll to "Authorized domains"
3. Add `localhost` if not present
4. Make sure `localhost` is in the list

### 5. API Key Expired or Regenerated
The API key might have been regenerated in the past, making the old one invalid.

**Solution:**
1. Go to: https://console.firebase.google.com/project/ecosistema-imob-dev/settings/general/web
2. In the "Your apps" section, find the web app
3. If needed, delete the old web app and create a new one
4. Copy the NEW firebaseConfig object
5. Update `.env.local` with all new values

## Testing
After making changes, test with:
```bash
node scripts/validate-firebase-config.js
```

## Alternative: Create New API Key
If all else fails, create a new web app:
1. Go to Firebase Console → Project Settings
2. Scroll to "Your apps"
3. Click "Add app" → Web
4. Register the app
5. Copy the new configuration
6. Update `.env.local` with new values
