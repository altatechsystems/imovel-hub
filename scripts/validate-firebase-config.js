const https = require('https');

// Read API key from command line or use default
const API_KEY = process.argv[2] || 'AIzaSyDx6PRtyNnAgTxsT3coow9Ut0xM4Xi5A71';
const PROJECT_ID = 'ecosistema-imob-dev';

console.log('🔍 Validating Firebase configuration...');
console.log(`   API Key: ${API_KEY.substring(0, 20)}...`);
console.log(`   Project ID: ${PROJECT_ID}`);
console.log('');

// Test the API key by making a request to Firebase Auth REST API
const postData = JSON.stringify({
  returnSecureToken: true
});

const options = {
  hostname: 'identitytoolkit.googleapis.com',
  port: 443,
  path: `/v1/accounts:signUp?key=${API_KEY}`,
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Content-Length': postData.length
  }
};

const req = https.request(options, (res) => {
  let data = '';

  res.on('data', (chunk) => {
    data += chunk;
  });

  res.on('end', () => {
    console.log(`📡 Response Status: ${res.statusCode}`);
    console.log('');

    if (res.statusCode === 400) {
      const response = JSON.parse(data);

      // Check if it's the "MISSING_EMAIL" error (which means the API key is valid)
      if (response.error && response.error.message && response.error.message.includes('MISSING_EMAIL')) {
        console.log('✅ API Key is VALID!');
        console.log('   The Firebase API accepted the key.');
        console.log('');
        console.log('   Error "MISSING_EMAIL" is expected because we didn\'t provide credentials.');
        console.log('   This confirms the API key works correctly.');
      } else if (response.error && response.error.message && response.error.message.includes('API key not valid')) {
        console.log('❌ API Key is INVALID!');
        console.log('   Error:', response.error.message);
        console.log('');
        console.log('📋 Next steps:');
        console.log('   1. Go to: https://console.firebase.google.com/project/ecosistema-imob-dev/settings/general');
        console.log('   2. Scroll to "Your apps" section');
        console.log('   3. Find the Web App configuration');
        console.log('   4. Copy the "Web API Key" value');
        console.log('   5. Update frontend-admin/.env.local with the correct key');
      } else {
        console.log('⚠️  Unexpected error:');
        console.log(JSON.stringify(response, null, 2));
      }
    } else if (res.statusCode === 200) {
      console.log('✅ API Key is VALID!');
      console.log('   Anonymous signup succeeded (this means the key works).');
    } else {
      console.log('⚠️  Unexpected status code:', res.statusCode);
      console.log('   Response:', data);
    }
  });
});

req.on('error', (e) => {
  console.error('❌ Request failed:', e.message);
});

req.write(postData);
req.end();
