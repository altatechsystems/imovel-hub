'use client';

import Link from 'next/link';
import { Home } from 'lucide-react';
import { Button } from '@/components/ui/button';

export interface HeaderProps {
  variant?: 'default' | 'minimal';
}

export function Header({ variant = 'default' }: HeaderProps) {
  return (
    <header className="bg-white border-b sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2">
            <Home className="w-8 h-8 text-blue-600" />
            <span className="text-2xl font-bold text-gray-900">Imobiliária</span>
          </Link>

          {variant === 'default' && (
            <nav className="hidden md:flex items-center gap-6">
              <Link href="/imoveis" className="text-gray-700 hover:text-blue-600 font-medium">
                Imóveis
              </Link>
              <Link href="/sobre" className="text-gray-700 hover:text-blue-600 font-medium">
                Sobre
              </Link>
              <Link href="/contato" className="text-gray-700 hover:text-blue-600 font-medium">
                Contato
              </Link>
              <Button variant="primary" size="sm">
                Anunciar Imóvel
              </Button>
            </nav>
          )}

          {variant === 'minimal' && (
            <nav className="hidden md:flex items-center gap-6">
              <Link href="/imoveis" className="text-gray-700 hover:text-blue-600 font-medium">
                Voltar para Imóveis
              </Link>
            </nav>
          )}
        </div>
      </div>
    </header>
  );
}
