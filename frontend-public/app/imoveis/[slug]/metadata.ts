import { Metadata } from 'next';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID || 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';
const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001';

async function getPropertyBySlug(slug: string) {
  try {
    const response = await fetch(`${API_URL}/${TENANT_ID}/properties/slug/${slug}`, {
      next: { revalidate: 3600 }, // Revalidate every hour
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error fetching property for metadata:', error);
    return null;
  }
}

function getPropertyTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    apartment: 'Apartamento',
    house: 'Casa',
    condo: 'Condomínio',
    land: 'Terreno',
    commercial: 'Comercial',
    farm: 'Fazenda',
    warehouse: 'Galpão',
    office: 'Escritório',
  };
  return labels[type] || 'Imóvel';
}

function getTransactionTypeLabel(type: string): string {
  const labels: Record<string, string> = {
    sale: 'Venda',
    rent: 'Aluguel',
    sale_rent: 'Venda/Aluguel',
  };
  return labels[type] || 'Venda';
}

export async function generatePropertyMetadata(slug: string): Promise<Metadata> {
  const property = await getPropertyBySlug(slug);

  if (!property) {
    return {
      title: 'Imóvel não encontrado',
      description: 'O imóvel que você está procurando não foi encontrado.',
    };
  }

  const propertyType = getPropertyTypeLabel(property.property_type);
  const transactionType = getTransactionTypeLabel(property.transaction_type || 'sale');
  const location = property.neighborhood
    ? `${property.neighborhood}, ${property.city} - ${property.state}`
    : `${property.city} - ${property.state}`;

  const title = property.title || `${propertyType} em ${property.city}`;
  const description = property.description ||
    `${propertyType} para ${transactionType} em ${location}. ${property.bedrooms ? `${property.bedrooms} quartos, ` : ''}${property.bathrooms ? `${property.bathrooms} banheiros, ` : ''}${property.parking_spaces ? `${property.parking_spaces} vagas` : ''}`.trim();

  const price = property.sale_price || property.rental_price || property.price_amount;
  const formattedPrice = price ? new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: 'BRL',
  }).format(price) : '';

  const fullTitle = formattedPrice
    ? `${title} - ${formattedPrice}`
    : title;

  const images = property.images && property.images.length > 0
    ? property.images.map((img: any) => img.large_url || img.medium_url || img.thumb_url)
    : property.cover_image_url
      ? [property.cover_image_url]
      : [];

  return {
    title: fullTitle,
    description: description.substring(0, 160), // Limit for SEO
    keywords: [
      propertyType.toLowerCase(),
      transactionType.toLowerCase(),
      property.city?.toLowerCase(),
      property.state?.toLowerCase(),
      property.neighborhood?.toLowerCase(),
      'imóvel',
      'imobiliária',
      property.bedrooms ? `${property.bedrooms} quartos` : '',
    ].filter(Boolean),
    openGraph: {
      title: fullTitle,
      description: description.substring(0, 200),
      url: `${SITE_URL}/imoveis/${slug}`,
      type: 'website',
      images: images.slice(0, 4).map((url: string) => ({
        url,
        width: 1200,
        height: 630,
        alt: title,
      })),
      locale: 'pt_BR',
      siteName: 'Imobiliária',
    },
    twitter: {
      card: 'summary_large_image',
      title: fullTitle,
      description: description.substring(0, 200),
      images: images.slice(0, 1),
    },
    alternates: {
      canonical: `${SITE_URL}/imoveis/${slug}`,
    },
    robots: {
      index: property.visibility === 'public' && property.status === 'available',
      follow: true,
    },
  };
}
