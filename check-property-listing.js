const fs = require('fs');
const data = JSON.parse(fs.readFileSync(0, 'utf-8'));

console.log('Total properties:', data.data?.length || 0);

if (data.data && data.data.length > 0) {
  const p = data.data[0];
  console.log('\nFirst property:');
  console.log('  ID:', p.id);
  console.log('  Title:', p.title || 'N/A');
  console.log('  City:', p.city);
  console.log('  Canonical Listing ID:', p.canonical_listing_id || 'MISSING');
  console.log('  Cover Image URL:', p.cover_image_url || 'MISSING');

  console.log('\nProperties with cover_image_url:', data.data.filter(p => p.cover_image_url).length);
  console.log('Properties without cover_image_url:', data.data.filter(p => !p.cover_image_url).length);
}
