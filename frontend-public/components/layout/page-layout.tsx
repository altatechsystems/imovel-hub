'use client';

import { ReactNode } from 'react';
import { Header } from './header';
import { Footer } from './footer';

export interface PageLayoutProps {
  children: ReactNode;
  headerVariant?: 'default' | 'minimal';
  showFooter?: boolean;
}

export function PageLayout({
  children,
  headerVariant = 'default',
  showFooter = true
}: PageLayoutProps) {
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <Header variant={headerVariant} />
      <main className="flex-1">
        {children}
      </main>
      {showFooter && <Footer />}
    </div>
  );
}
