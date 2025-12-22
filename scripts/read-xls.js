const XLSX = require('xlsx');
const path = require('path');

const xlsPath = path.join(__dirname, '..', 'univen-imoveis_20-12-2025_18_12_15.xls');

try {
  const wb = XLSX.readFile(xlsPath);
  const ws = wb.Sheets[wb.SheetNames[0]];
  const data = XLSX.utils.sheet_to_json(ws, { defval: '' });

  console.log('Total rows:', data.length);
  console.log('\nFirst 2 rows:');
  console.log(JSON.stringify(data.slice(0, 2), null, 2));
  console.log('\nColumns:', Object.keys(data[0] || {}).join(', '));
} catch (error) {
  console.error('Error reading file:', error.message);
}
