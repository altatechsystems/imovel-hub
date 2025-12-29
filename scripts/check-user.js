const admin = require('firebase-admin');
const path = require('path');

// Initialize Firebase Admin
const serviceAccount = require(path.join(__dirname, '..', 'backend', 'config', 'firebase-adminsdk.json'));

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: 'https://ecosistema-imob-dev.firebaseio.com'
});

const email = process.argv[2] || 'administracao@altatechimmoveis.com';

console.log(`🔍 Checking if user exists: ${email}`);
console.log('');

admin.auth().getUserByEmail(email)
  .then((userRecord) => {
    console.log('✅ User found in Firebase Authentication!');
    console.log('');
    console.log('User details:');
    console.log('  UID:', userRecord.uid);
    console.log('  Email:', userRecord.email);
    console.log('  Email verified:', userRecord.emailVerified);
    console.log('  Disabled:', userRecord.disabled);
    console.log('  Created:', new Date(userRecord.metadata.creationTime));
    console.log('  Last sign in:', userRecord.metadata.lastSignInTime || 'Never');
    console.log('');

    if (userRecord.customClaims) {
      console.log('Custom claims:', userRecord.customClaims);
    }

    process.exit(0);
  })
  .catch((error) => {
    console.log('❌ User NOT found in Firebase Authentication');
    console.log('');
    console.log('Error:', error.code);
    console.log('Message:', error.message);
    console.log('');
    console.log('Available users with email daniel.garcia@altatechsystems.com should work.');
    console.log('');
    console.log('To create this user, you need to:');
    console.log('1. Sign up at: http://localhost:3002/cadastre-se');
    console.log('2. Or use the existing user: daniel.garcia@altatechsystems.com');

    process.exit(1);
  });
