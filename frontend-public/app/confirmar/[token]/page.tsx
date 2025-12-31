'use client';

import { useState, useEffect } from 'react';
import { useParams, useSearchParams } from 'next/navigation';
import { CheckCircle, XCircle, DollarSign, Clock, AlertCircle, Home } from 'lucide-react';

interface ConfirmationPageData {
  valid: boolean;
  property_id?: string;
  property_type?: string;
  neighborhood?: string;
  city?: string;
  reference?: string;
  current_status?: string;
  current_price?: number;
  expires_at?: string;
  error?: string;
}

type ConfirmationAction = 'confirm_available' | 'confirm_unavailable' | 'confirm_price';

export default function OwnerConfirmationPage() {
  const params = useParams();
  const searchParams = useSearchParams();
  const token = params.token as string;
  const tenantId = searchParams.get('tenant_id');

  const [loading, setLoading] = useState(true);
  const [data, setData] = useState<ConfirmationPageData | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');
  const [priceInput, setPriceInput] = useState('');
  const [showPriceInput, setShowPriceInput] = useState(false);

  useEffect(() => {
    if (token && tenantId) {
      validateToken();
    } else {
      setError('Link inválido: parâmetros faltando');
      setLoading(false);
    }
  }, [token, tenantId]);

  const validateToken = async () => {
    try {
      setLoading(true);
      const apiUrl = process.env.NEXT_PUBLIC_API_URL?.replace('/api', '') || 'http://localhost:8080';
      const response = await fetch(`${apiUrl}/confirmar/${token}?tenant_id=${tenantId}`);

      if (!response.ok) {
        throw new Error('Token inválido ou expirado');
      }

      const result = await response.json();

      if (result.success && result.data) {
        setData(result.data);
        if (result.data.current_price) {
          setPriceInput(result.data.current_price.toString());
        }
      } else {
        setError(result.data?.error || 'Token inválido');
      }
    } catch (err: any) {
      console.error('Erro ao validar token:', err);
      setError('Não foi possível validar o link. Tente novamente ou entre em contato com a imobiliária.');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (action: ConfirmationAction) => {
    if (action === 'confirm_price') {
      if (!showPriceInput) {
        setShowPriceInput(true);
        return;
      }

      if (!priceInput || parseFloat(priceInput) <= 0) {
        setError('Por favor, informe um preço válido');
        return;
      }
    }

    try {
      setSubmitting(true);
      setError('');

      const body: any = { action };
      if (action === 'confirm_price') {
        body.price_amount = parseFloat(priceInput);
      }

      const apiUrl = process.env.NEXT_PUBLIC_API_URL?.replace('/api', '') || 'http://localhost:8080';
      const response = await fetch(
        `${apiUrl}/api/v1/owner-confirmations/${token}/submit?tenant_id=${tenantId}`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(body),
        }
      );

      const result = await response.json();

      if (result.success) {
        setSuccess(true);
      } else {
        setError(result.error || 'Erro ao processar confirmação');
      }
    } catch (err: any) {
      console.error('Erro ao enviar confirmação:', err);
      setError('Erro ao processar confirmação. Tente novamente.');
    } finally {
      setSubmitting(false);
    }
  };

  const getPropertyTypeLabel = (type: string) => {
    const types: Record<string, string> = {
      apartment: 'Apartamento',
      house: 'Casa',
      condo: 'Condomínio',
      commercial: 'Comercial',
      land: 'Terreno',
      farm: 'Fazenda',
      studio: 'Studio',
      penthouse: 'Cobertura',
      townhouse: 'Sobrado',
    };
    return types[type] || type;
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
        <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full">
          <div className="flex flex-col items-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mb-4"></div>
            <p className="text-gray-600">Validando link...</p>
          </div>
        </div>
      </div>
    );
  }

  if (success) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-green-50 to-emerald-100 flex items-center justify-center p-4">
        <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full">
          <div className="flex flex-col items-center text-center">
            <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
              <CheckCircle className="w-10 h-10 text-green-600" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Confirmação Recebida!</h1>
            <p className="text-gray-600 mb-6">
              Obrigado por confirmar as informações do seu imóvel. Suas informações foram atualizadas com sucesso.
            </p>
            <div className="bg-green-50 border border-green-200 rounded-lg p-4 w-full">
              <p className="text-sm text-green-800">
                ✓ A imobiliária foi notificada automaticamente
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  if (error || !data || !data.valid) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-red-50 to-rose-100 flex items-center justify-center p-4">
        <div className="bg-white rounded-lg shadow-xl p-8 max-w-md w-full">
          <div className="flex flex-col items-center text-center">
            <div className="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mb-4">
              <XCircle className="w-10 h-10 text-red-600" />
            </div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Link Inválido ou Expirado</h1>
            <p className="text-gray-600 mb-6">
              {error || data?.error || 'Este link de confirmação não é mais válido.'}
            </p>
            <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 w-full">
              <p className="text-sm text-amber-800">
                Entre em contato com a imobiliária para solicitar um novo link de confirmação.
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 flex items-center justify-center p-4">
      <div className="bg-white rounded-lg shadow-xl p-6 md:p-8 max-w-2xl w-full">
        {/* Header */}
        <div className="flex items-center gap-3 mb-6">
          <div className="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
            <Home className="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-gray-900">Confirmação de Imóvel</h1>
            <p className="text-sm text-gray-600">Atualize as informações do seu imóvel</p>
          </div>
        </div>

        {/* Property Info */}
        <div className="bg-gray-50 rounded-lg p-4 mb-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
            {data.reference && (
              <div>
                <span className="text-gray-600">Código:</span>{' '}
                <span className="font-medium text-gray-900">{data.reference}</span>
              </div>
            )}
            {data.property_type && (
              <div>
                <span className="text-gray-600">Tipo:</span>{' '}
                <span className="font-medium text-gray-900">
                  {getPropertyTypeLabel(data.property_type)}
                </span>
              </div>
            )}
            {data.neighborhood && (
              <div>
                <span className="text-gray-600">Bairro:</span>{' '}
                <span className="font-medium text-gray-900">{data.neighborhood}</span>
              </div>
            )}
            {data.city && (
              <div>
                <span className="text-gray-600">Cidade:</span>{' '}
                <span className="font-medium text-gray-900">{data.city}</span>
              </div>
            )}
          </div>
        </div>

        {/* Current Status */}
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <div className="flex items-center gap-2 mb-2">
            <AlertCircle className="w-5 h-5 text-blue-600" />
            <h2 className="font-semibold text-gray-900">Status Atual</h2>
          </div>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
            <div>
              <span className="text-gray-600">Disponibilidade:</span>{' '}
              <span className={`font-medium ${
                data.current_status === 'available' ? 'text-green-700' : 'text-gray-700'
              }`}>
                {data.current_status === 'available' ? 'Disponível' : data.current_status || 'Não informado'}
              </span>
            </div>
            {data.current_price && (
              <div>
                <span className="text-gray-600">Preço:</span>{' '}
                <span className="font-medium text-gray-900">
                  R$ {data.current_price.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                </span>
              </div>
            )}
          </div>
        </div>

        {/* Expiration Warning */}
        {data.expires_at && (
          <div className="bg-amber-50 border border-amber-200 rounded-lg p-3 mb-6">
            <div className="flex items-center gap-2 text-sm text-amber-800">
              <Clock className="w-4 h-4" />
              <span>
                Este link expira em {new Date(data.expires_at).toLocaleDateString('pt-BR')} às{' '}
                {new Date(data.expires_at).toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })}
              </span>
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="space-y-4">
          <h2 className="font-semibold text-gray-900 text-lg">Selecione uma ação:</h2>

          {/* Confirm Available */}
          <button
            onClick={() => handleSubmit('confirm_available')}
            disabled={submitting}
            className="w-full flex items-center justify-center gap-3 px-4 py-4 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium"
          >
            <CheckCircle className="w-5 h-5" />
            {submitting ? 'Processando...' : 'Confirmar que o imóvel está DISPONÍVEL'}
          </button>

          {/* Confirm Unavailable */}
          <button
            onClick={() => handleSubmit('confirm_unavailable')}
            disabled={submitting}
            className="w-full flex items-center justify-center gap-3 px-4 py-4 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium"
          >
            <XCircle className="w-5 h-5" />
            {submitting ? 'Processando...' : 'Informar que o imóvel NÃO está mais disponível'}
          </button>

          {/* Update Price */}
          <div className="bg-gray-50 rounded-lg p-4 border border-gray-200">
            <button
              onClick={() => handleSubmit('confirm_price')}
              disabled={submitting}
              className="w-full flex items-center justify-center gap-3 px-4 py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-300 disabled:cursor-not-allowed transition-colors font-medium mb-3"
            >
              <DollarSign className="w-5 h-5" />
              {submitting ? 'Processando...' : 'Atualizar Preço'}
            </button>

            {showPriceInput && (
              <div className="space-y-2">
                <label className="block text-sm font-medium text-gray-700">
                  Novo preço:
                </label>
                <input
                  type="number"
                  value={priceInput}
                  onChange={(e) => setPriceInput(e.target.value)}
                  placeholder="Ex: 500000"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  disabled={submitting}
                  step="0.01"
                  min="0"
                />
                <p className="text-xs text-gray-500">
                  Valor atual: R$ {data.current_price?.toLocaleString('pt-BR', { minimumFractionDigits: 2 }) || '0,00'}
                </p>
              </div>
            )}
          </div>

          {error && (
            <div className="bg-red-50 border border-red-200 rounded-lg p-3">
              <p className="text-sm text-red-800 text-center">{error}</p>
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="mt-8 pt-6 border-t border-gray-200">
          <p className="text-xs text-gray-500 text-center">
            Ao confirmar, você autoriza a imobiliária a atualizar as informações do seu imóvel conforme a ação selecionada.
          </p>
        </div>
      </div>
    </div>
  );
}
