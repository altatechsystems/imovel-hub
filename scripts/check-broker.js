const admin = require('firebase-admin');

// Initialize Firebase Admin
const serviceAccount = require('../backend/config/firebase-adminsdk.json');

admin.initializeApp({
  credential: admin.credential.cert(serviceAccount)
});

const db = admin.firestore();
db.settings({ databaseId: 'imob-dev' });

async function checkBroker() {
  try {
    const email = 'daniel.garcia@altatechsystems.com';
    const tenantId = '391b12f8-ebe4-426a-8c99-ec5a10b1f361';
    const brokerId = '73f624cc-2db1-4a2f-9a95-8b21abffc8d7';

    console.log(`\nChecking broker for email: ${email}`);
    console.log(`Tenant ID: ${tenantId}`);
    console.log(`Broker ID: ${brokerId}\n`);

    // 1. Get Firebase Auth user
    const userRecord = await admin.auth().getUserByEmail(email);
    console.log('✅ Firebase Auth User found:');
    console.log(`   UID: ${userRecord.uid}`);
    console.log(`   Email: ${userRecord.email}`);
    console.log(`   Custom Claims:`, userRecord.customClaims || 'None');

    // 2. Check broker document
    const brokerRef = db.collection('tenants').doc(tenantId).collection('brokers').doc(brokerId);
    const brokerDoc = await brokerRef.get();

    if (brokerDoc.exists) {
      const brokerData = brokerDoc.data();
      console.log('\n✅ Broker document found:');
      console.log('   Broker Data:', JSON.stringify(brokerData, null, 2));
    } else {
      console.log('\n❌ Broker document NOT found!');
      console.log(`   Path: /tenants/${tenantId}/brokers/${brokerId}`);
    }

    // 3. Search all brokers with this firebase_uid
    console.log(`\n🔍 Searching all brokers with firebase_uid: ${userRecord.uid}`);
    const brokersQuery = await db.collectionGroup('brokers')
      .where('firebase_uid', '==', userRecord.uid)
      .get();

    if (brokersQuery.empty) {
      console.log('❌ No brokers found with this firebase_uid');
    } else {
      console.log(`✅ Found ${brokersQuery.size} broker(s):`);
      brokersQuery.forEach(doc => {
        console.log(`   - ${doc.ref.path}`);
        console.log(`     Data:`, JSON.stringify(doc.data(), null, 2));
      });
    }

    // 4. Check tenant document
    const tenantRef = db.collection('tenants').doc(tenantId);
    const tenantDoc = await tenantRef.get();

    if (tenantDoc.exists) {
      console.log('\n✅ Tenant document found:');
      console.log('   Tenant Data:', JSON.stringify(tenantDoc.data(), null, 2));
    } else {
      console.log('\n❌ Tenant document NOT found!');
    }

  } catch (error) {
    console.error('❌ Error:', error.message);
    console.error(error);
  } finally {
    process.exit(0);
  }
}

checkBroker();
