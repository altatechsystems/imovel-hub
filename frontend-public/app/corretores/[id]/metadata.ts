import { Metadata } from 'next';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const TENANT_ID = process.env.NEXT_PUBLIC_TENANT_ID || 'bd71c02b-5fa5-43df-8b46-a1df2206f1ef';
const SITE_URL = process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001';

interface Broker {
  id: string;
  name: string;
  email: string;
  phone?: string;
  creci: string;
  bio?: string;
  photo_url?: string;
  company?: string;
  specialties?: string;
  languages?: string;
  rating?: number;
  review_count?: number;
  experience?: number;
}

async function getBrokerById(brokerId: string): Promise<Broker | null> {
  try {
    const response = await fetch(`${API_URL}/${TENANT_ID}/brokers/${brokerId}/public`, {
      next: { revalidate: 3600 }, // Revalidate every hour
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error('Error fetching broker for metadata:', error);
    return null;
  }
}

export async function generateBrokerMetadata(brokerId: string): Promise<Metadata> {
  const broker = await getBrokerById(brokerId);

  if (!broker) {
    return {
      title: 'Corretor não encontrado',
      description: 'O corretor solicitado não foi encontrado.',
      robots: {
        index: false,
        follow: false,
      },
    };
  }

  const firstName = broker.name.split(' ')[0];
  const title = `${broker.name} - Corretor de Imóveis CRECI ${broker.creci}`;

  const description = broker.bio
    ? broker.bio.substring(0, 160)
    : `Conheça ${broker.name}, corretor de imóveis especializado. CRECI ${broker.creci}${broker.experience ? ` com ${broker.experience} anos de experiência` : ''}${broker.rating ? ` - ${broker.rating.toFixed(1)} estrelas` : ''}. Entre em contato e encontre o imóvel ideal.`;

  const keywords = [
    broker.name,
    `corretor ${firstName}`,
    `CRECI ${broker.creci}`,
    'corretor de imóveis',
    'agente imobiliário',
    'consultor imobiliário',
    'comprar imóvel',
    'vender imóvel',
    'alugar imóvel',
  ];

  if (broker.company) {
    keywords.push(broker.company);
  }

  if (broker.specialties) {
    const specialtiesList = broker.specialties.split(',').map(s => s.trim());
    keywords.push(...specialtiesList);
  }

  const images = [];
  if (broker.photo_url) {
    images.push({
      url: broker.photo_url,
      width: 800,
      height: 800,
      alt: broker.name,
    });
  }

  return {
    title,
    description,
    keywords,
    openGraph: {
      title,
      description,
      url: `${SITE_URL}/corretores/${brokerId}`,
      type: 'profile',
      locale: 'pt_BR',
      siteName: 'Imobiliária',
      images: images.length > 0 ? images : undefined,
      ...(broker.name && {
        profile: {
          firstName: broker.name.split(' ')[0],
          lastName: broker.name.split(' ').slice(1).join(' '),
        },
      }),
    },
    twitter: {
      card: broker.photo_url ? 'summary' : 'summary_large_image',
      title,
      description,
      images: broker.photo_url ? [broker.photo_url] : undefined,
    },
    alternates: {
      canonical: `${SITE_URL}/corretores/${brokerId}`,
    },
    robots: {
      index: true,
      follow: true,
    },
  };
}
