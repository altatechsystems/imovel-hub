import { Metadata } from 'next';
import { generatePropertyMetadata } from './metadata';

type Props = {
  params: Promise<{ slug: string }>;
  children: React.ReactNode;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  return generatePropertyMetadata(slug);
}

export default function PropertyLayout({ children }: Props) {
  return <>{children}</>;
}
