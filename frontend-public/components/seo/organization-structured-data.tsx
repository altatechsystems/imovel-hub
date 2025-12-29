'use client';

export function OrganizationStructuredData() {
  const baseUrl = process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001';

  const structuredData = {
    '@context': 'https://schema.org',
    '@type': 'RealEstateAgent',
    name: 'Imobiliária',
    url: baseUrl,
    logo: `${baseUrl}/logo.png`,
    description: 'Encontre apartamentos, casas, terrenos e imóveis comerciais para venda e aluguel em todo o Brasil.',
    contactPoint: {
      '@type': 'ContactPoint',
      telephone: process.env.NEXT_PUBLIC_WHATSAPP || '+55-00-0000-0000',
      contactType: 'customer service',
      areaServed: 'BR',
      availableLanguage: ['Portuguese'],
    },
    address: {
      '@type': 'PostalAddress',
      addressCountry: 'BR',
    },
    sameAs: [
      // Add social media profiles here
      'https://www.facebook.com/imobiliaria',
      'https://www.instagram.com/imobiliaria',
      'https://www.linkedin.com/company/imobiliaria',
    ],
  };

  return (
    <script
      type="application/ld+json"
      dangerouslySetInnerHTML={{ __html: JSON.stringify(structuredData) }}
    />
  );
}
