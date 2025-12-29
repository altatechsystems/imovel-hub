const fs = require('fs');
const data = JSON.parse(fs.readFileSync(0, 'utf-8'));

if (!data.data) {
  console.log('Error: No data field in response');
  console.log(JSON.stringify(data, null, 2));
  process.exit(1);
}

const p = data.data;
console.log('Property ID:', p.id);
console.log('Title:', p.title || 'N/A');
console.log('City:', p.city);
console.log('Canonical Listing ID:', p.canonical_listing_id || 'MISSING');
console.log('Cover Image URL:', p.cover_image_url || 'MISSING');
console.log('Images:', p.images ? `${p.images.length} images` : 'MISSING');

if (p.images && p.images.length > 0) {
  console.log('\nFirst image:');
  console.log('  ID:', p.images[0].id);
  console.log('  URL:', p.images[0].url || p.images[0].large_url);
  console.log('  Thumb:', p.images[0].thumb_url);
}
