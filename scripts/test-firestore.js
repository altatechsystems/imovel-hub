const admin = require('firebase-admin');
const path = require('path');

// Initialize Firebase Admin
const serviceAccount = require(path.join(__dirname, '..', 'backend', 'config', 'firebase-adminsdk.json'));

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount),
  databaseURL: `https://ecosistema-imob-dev.firebaseio.com`
});

// Get Firestore database
const db = admin.firestore(admin.app(), 'imob-dev');

async function testFirestore() {
  try {
    // Query properties
    const snapshot = await db.collection('properties').limit(10).get();

    console.log(`Found ${snapshot.size} properties in Firestore`);
    console.log('\nFirst 5 properties:');

    let count = 0;
    snapshot.forEach(doc => {
      if (count < 5) {
        const data = doc.data();
        console.log(`\n${count + 1}. ${data.reference} - ${data.title}`);
        console.log(`   Type: ${data.type}, Status: ${data.status}`);
        console.log(`   City: ${data.address?.city}, Price: R$ ${data.sale_price}`);
      }
      count++;
    });

    // Check total count
    const allSnapshot = await db.collection('properties').count().get();
    console.log(`\n\nTotal properties in database: ${allSnapshot.data().count}`);

  } catch (error) {
    console.error('Error:', error.message);
  } finally {
    process.exit(0);
  }
}

testFirestore();
