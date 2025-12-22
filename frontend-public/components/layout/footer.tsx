'use client';

import Link from 'next/link';
import { Home } from 'lucide-react';

export function Footer() {
  return (
    <footer className="bg-gray-900 text-gray-300 py-12">
      <div className="container mx-auto px-4">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
          <div>
            <div className="flex items-center gap-2 mb-4">
              <Home className="w-6 h-6 text-blue-400" />
              <span className="text-xl font-bold text-white">Imobiliária</span>
            </div>
            <p className="text-sm">
              Seu parceiro confiável para encontrar o imóvel perfeito.
            </p>
          </div>

          <div>
            <h3 className="font-semibold text-white mb-4">Links Rápidos</h3>
            <ul className="space-y-2 text-sm">
              <li><Link href="/imoveis" className="hover:text-blue-400">Imóveis</Link></li>
              <li><Link href="/sobre" className="hover:text-blue-400">Sobre Nós</Link></li>
              <li><Link href="/contato" className="hover:text-blue-400">Contato</Link></li>
            </ul>
          </div>

          <div>
            <h3 className="font-semibold text-white mb-4">Categorias</h3>
            <ul className="space-y-2 text-sm">
              <li><Link href="/imoveis?type=apartment" className="hover:text-blue-400">Apartamentos</Link></li>
              <li><Link href="/imoveis?type=house" className="hover:text-blue-400">Casas</Link></li>
              <li><Link href="/imoveis?type=commercial" className="hover:text-blue-400">Comerciais</Link></li>
            </ul>
          </div>

          <div>
            <h3 className="font-semibold text-white mb-4">Contato</h3>
            <ul className="space-y-2 text-sm">
              <li>Email: contato@imobiliaria.com</li>
              <li>Telefone: (11) 3000-0000</li>
              <li>WhatsApp: (11) 99999-9999</li>
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-800 pt-8 text-center text-sm">
          <p>&copy; 2025 Imobiliária. Todos os direitos reservados.</p>
        </div>
      </div>
    </footer>
  );
}
