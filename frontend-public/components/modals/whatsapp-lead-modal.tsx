'use client';

import * as React from 'react';
import { X, Check } from 'lucide-react';
import { Property } from '@/types/property';

interface WhatsAppLeadModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSubmit: (name: string, phone: string) => void;
  property: Property;
  isLoading?: boolean;
}

export const WhatsAppLeadModal: React.FC<WhatsAppLeadModalProps> = ({
  isOpen,
  onClose,
  onSubmit,
  property,
  isLoading = false,
}) => {
  const [name, setName] = React.useState('');
  const [phone, setPhone] = React.useState('');
  const [showTyping, setShowTyping] = React.useState(false);

  // Load saved data from localStorage
  React.useEffect(() => {
    if (isOpen) {
      const savedName = localStorage.getItem('lead_name');
      const savedPhone = localStorage.getItem('lead_phone');
      if (savedName) setName(savedName);
      if (savedPhone) setPhone(savedPhone);

      // Show typing animation
      setTimeout(() => setShowTyping(true), 300);
    } else {
      setShowTyping(false);
    }
  }, [isOpen]);

  const formatPhoneNumber = (value: string) => {
    // Remove tudo exceto números
    const numbers = value.replace(/\D/g, '');

    // Aplica máscara (XX) XXXXX-XXXX
    if (numbers.length <= 2) {
      return numbers;
    } else if (numbers.length <= 7) {
      return `(${numbers.slice(0, 2)}) ${numbers.slice(2)}`;
    } else if (numbers.length <= 11) {
      return `(${numbers.slice(0, 2)}) ${numbers.slice(2, 7)}-${numbers.slice(7)}`;
    }
    return `(${numbers.slice(0, 2)}) ${numbers.slice(2, 7)}-${numbers.slice(7, 11)}`;
  };

  const handlePhoneChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const formatted = formatPhoneNumber(e.target.value);
    setPhone(formatted);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!name.trim() || phone.replace(/\D/g, '').length < 10) {
      return;
    }

    // Save to localStorage for next time
    localStorage.setItem('lead_name', name.trim());
    localStorage.setItem('lead_phone', phone);

    onSubmit(name.trim(), phone.replace(/\D/g, ''));
  };

  const isValid = name.trim().length >= 2 && phone.replace(/\D/g, '').length >= 10;

  console.log('Modal render - isOpen:', isOpen);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-end sm:items-center justify-center">
      {/* Overlay */}
      <div
        className="absolute inset-0 bg-black/50 backdrop-blur-sm"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative w-full max-w-md mx-auto bg-[#ECE5DD] rounded-t-2xl sm:rounded-2xl shadow-2xl max-h-[90vh] overflow-hidden animate-slide-up">
        {/* Header estilo WhatsApp */}
        <div className="bg-[#075E54] text-white px-4 py-3 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-full bg-white/20 flex items-center justify-center text-lg font-semibold">
              🏠
            </div>
            <div>
              <h3 className="font-semibold text-sm">Corretor</h3>
              <p className="text-xs text-white/80">online</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-white/10 rounded-full transition-colors"
            aria-label="Fechar"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Chat area */}
        <div className="p-4 space-y-3 min-h-[300px] max-h-[500px] overflow-y-auto">
          {/* Mensagem do corretor */}
          <div className="flex gap-2">
            <div className="max-w-[85%]">
              <div className="bg-white rounded-lg rounded-tl-none shadow-sm p-3">
                <p className="text-sm text-gray-800">
                  Olá! Vi que você tem interesse neste imóvel:
                </p>
              </div>
              <span className="text-xs text-gray-500 ml-2">agora</span>
            </div>
          </div>

          {/* Card do imóvel */}
          <div className="flex gap-2">
            <div className="max-w-[85%]">
              <div className="bg-white rounded-lg shadow-sm overflow-hidden">
                {property.cover_image_url && (
                  <img
                    src={property.cover_image_url}
                    alt={property.title || 'Imóvel'}
                    className="w-full h-32 object-cover"
                  />
                )}
                <div className="p-3">
                  <p className="font-semibold text-sm text-gray-900 line-clamp-2">
                    {property.title || `${property.property_type} em ${property.city}`}
                  </p>
                  <p className="text-xs text-gray-600 mt-1">
                    {property.neighborhood}, {property.city}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Mensagem solicitando dados */}
          {showTyping && (
            <div className="flex gap-2 animate-fade-in">
              <div className="max-w-[85%]">
                <div className="bg-white rounded-lg rounded-tl-none shadow-sm p-3">
                  <p className="text-sm text-gray-800">
                    Para continuar nossa conversa no WhatsApp, preciso de algumas informações:
                  </p>
                </div>
                <span className="text-xs text-gray-500 ml-2">agora</span>
              </div>
            </div>
          )}
        </div>

        {/* Form area - estilo WhatsApp input */}
        <form onSubmit={handleSubmit} className="bg-[#F0F0F0] border-t border-gray-300">
          <div className="p-3 space-y-2">
            {/* Nome */}
            <div className="bg-white rounded-lg shadow-sm">
              <input
                type="text"
                placeholder="Seu nome completo"
                value={name}
                onChange={(e) => setName(e.target.value)}
                className="w-full px-4 py-3 text-sm bg-transparent border-none outline-none focus:outline-none"
                required
                minLength={2}
                disabled={isLoading}
              />
            </div>

            {/* Telefone */}
            <div className="bg-white rounded-lg shadow-sm">
              <input
                type="tel"
                placeholder="(00) 00000-0000"
                value={phone}
                onChange={handlePhoneChange}
                className="w-full px-4 py-3 text-sm bg-transparent border-none outline-none focus:outline-none"
                required
                disabled={isLoading}
              />
            </div>

            {/* Botão enviar - estilo WhatsApp */}
            <button
              type="submit"
              disabled={!isValid || isLoading}
              className="w-full bg-[#25D366] hover:bg-[#22C55E] disabled:bg-gray-300 disabled:cursor-not-allowed text-white font-medium py-3 px-6 rounded-lg flex items-center justify-center gap-2 transition-colors shadow-sm"
            >
              {isLoading ? (
                <span className="flex items-center gap-2">
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                  Enviando...
                </span>
              ) : (
                <>
                  <Check className="w-5 h-5" />
                  Continuar no WhatsApp
                </>
              )}
            </button>

            {/* LGPD notice */}
            <p className="text-xs text-gray-600 text-center px-2">
              Ao continuar, você concorda com nossa Política de Privacidade e autoriza o uso dos seus dados para contato.
            </p>
          </div>
        </form>
      </div>

      <style jsx>{`
        @keyframes slide-up {
          from {
            transform: translateY(100%);
            opacity: 0;
          }
          to {
            transform: translateY(0);
            opacity: 1;
          }
        }
        @keyframes fade-in {
          from {
            opacity: 0;
          }
          to {
            opacity: 1;
          }
        }
        .animate-slide-up {
          animation: slide-up 0.3s ease-out;
        }
        .animate-fade-in {
          animation: fade-in 0.5s ease-out;
        }
      `}</style>
    </div>
  );
};
