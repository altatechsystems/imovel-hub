import { Metadata } from 'next';
import { generateBrokerMetadata } from './metadata';

type Props = {
  params: Promise<{ id: string }>;
  children: React.ReactNode;
};

export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { id } = await params;
  return generateBrokerMetadata(id);
}

export default function BrokerProfileLayout({ children }: Props) {
  return <>{children}</>;
}
