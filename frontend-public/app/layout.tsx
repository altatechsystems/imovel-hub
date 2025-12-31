import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Providers } from "@/lib/providers";
import "./globals.css";

const inter = Inter({
  subsets: ["latin"],
  variable: "--font-inter",
});

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001'),
  title: {
    default: "Imobiliária - Encontre seu Imóvel Ideal | Casas, Apartamentos e Terrenos",
    template: "%s | Imobiliária"
  },
  description: "Encontre apartamentos, casas, terrenos e imóveis comerciais para venda e aluguel. Milhares de propriedades em todo o Brasil com fotos, preços atualizados e localização.",
  keywords: [
    "imóveis",
    "apartamentos",
    "casas",
    "terrenos",
    "venda",
    "aluguel",
    "imobiliária",
    "comprar imóvel",
    "alugar imóvel",
    "imóveis à venda",
    "imóveis para alugar",
    "corretor de imóveis",
    "propriedades",
    "real estate",
    "imóveis Brasil"
  ],
  authors: [{ name: "Imobiliária" }],
  creator: "Imobiliária",
  publisher: "Imobiliária",
  formatDetection: {
    email: false,
    address: false,
    telephone: false,
  },
  openGraph: {
    type: "website",
    locale: "pt_BR",
    url: process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001',
    title: "Imobiliária - Encontre seu Imóvel Ideal",
    description: "Encontre apartamentos, casas, terrenos e imóveis comerciais para venda e aluguel em todo o Brasil.",
    siteName: "Imobiliária",
    images: [
      {
        url: "/og-image.jpg",
        width: 1200,
        height: 630,
        alt: "Imobiliária - Encontre seu Imóvel Ideal",
      },
    ],
  },
  twitter: {
    card: "summary_large_image",
    title: "Imobiliária - Encontre seu Imóvel Ideal",
    description: "Encontre apartamentos, casas, terrenos e imóveis comerciais para venda e aluguel em todo o Brasil.",
    images: ["/twitter-image.jpg"],
    creator: "@imobiliaria",
  },
  robots: {
    index: true,
    follow: true,
    googleBot: {
      index: true,
      follow: true,
      'max-video-preview': -1,
      'max-image-preview': 'large',
      'max-snippet': -1,
    },
  },
  verification: {
    google: process.env.NEXT_PUBLIC_GOOGLE_SITE_VERIFICATION,
    yandex: process.env.NEXT_PUBLIC_YANDEX_VERIFICATION,
  },
  alternates: {
    canonical: process.env.NEXT_PUBLIC_SITE_URL || 'http://localhost:3001',
  },
  category: 'real estate',
  classification: 'Real Estate, Property Listings',
  referrer: 'origin-when-cross-origin',
  applicationName: 'Imobiliária',
  appleWebApp: {
    capable: true,
    title: 'Imobiliária',
    statusBarStyle: 'default',
  },
  other: {
    'apple-mobile-web-app-capable': 'yes',
    'mobile-web-app-capable': 'yes',
    'theme-color': '#2563eb',
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="pt-BR">
      <body className={`${inter.variable} font-sans antialiased`}>
        <Providers>
          {children}
        </Providers>
      </body>
    </html>
  );
}
