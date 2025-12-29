'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Save, Upload, X } from 'lucide-react';

interface PropertyFormData {
  reference?: string;
  slug?: string;
  street?: string;
  city?: string;
  state?: string;
  neighborhood?: string;
  zip_code?: string;
  price_amount?: number;
  bedrooms?: number;
  bathrooms?: number;
  suites?: number;
  parking_spaces?: number;
  total_area?: number;
  built_area?: number;
  property_type: string;
  transaction_type: string;
  status: string;
  visibility: string;
  description?: string;
  featured?: boolean;
}

interface PhotoUpload {
  file: File;
  preview: string;
  order: number;
  uploading?: boolean;
  uploaded?: boolean;
  error?: string;
  urls?: {
    thumb_url: string;
    medium_url: string;
    large_url: string;
  };
}

export default function NewPropertyPage() {
  const router = useRouter();
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const [photos, setPhotos] = useState<PhotoUpload[]>([]);

  const [formData, setFormData] = useState<PropertyFormData>({
    property_type: 'apartment',
    transaction_type: 'sale',
    status: 'available',
    visibility: 'public',
    featured: false,
  });

  const handleChange = (field: keyof PropertyFormData, value: any) => {
    setFormData({ ...formData, [field]: value });
  };

  const handlePhotoSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = e.target.files;
    if (!files) return;

    const newPhotos: PhotoUpload[] = Array.from(files).map((file, index) => ({
      file,
      preview: URL.createObjectURL(file),
      order: photos.length + index,
    }));

    setPhotos([...photos, ...newPhotos]);
  };

  const removePhoto = (index: number) => {
    const newPhotos = photos.filter((_, i) => i !== index);
    setPhotos(newPhotos.map((photo, i) => ({ ...photo, order: i })));
  };

  const uploadPhotos = async (propertyId: string, tenantId: string, token: string): Promise<boolean> => {
    let allSuccess = true;

    for (let i = 0; i < photos.length; i++) {
      const photo = photos[i];

      setPhotos(prev => prev.map((p, idx) =>
        idx === i ? { ...p, uploading: true } : p
      ));

      try {
        const formData = new FormData();
        formData.append('file', photo.file);
        formData.append('order', photo.order.toString());

        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/admin/${tenantId}/property-images/${propertyId}`,
          {
            method: 'POST',
            headers: {
              'Authorization': `Bearer ${token}`,
            },
            body: formData,
          }
        );

        if (!response.ok) {
          throw new Error('Erro ao fazer upload da foto');
        }

        const result = await response.json();

        setPhotos(prev => prev.map((p, idx) =>
          idx === i ? {
            ...p,
            uploading: false,
            uploaded: true,
            urls: result.data
          } : p
        ));
      } catch (err: any) {
        console.error(`Erro ao fazer upload da foto ${i}:`, err);
        allSuccess = false;

        setPhotos(prev => prev.map((p, idx) =>
          idx === i ? {
            ...p,
            uploading: false,
            uploaded: false,
            error: err.message
          } : p
        ));
      }
    }

    return allSuccess;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      setSaving(true);
      setError('');

      const tenantId = localStorage.getItem('tenant_id');
      if (!tenantId) {
        setError('Tenant ID não encontrado');
        return;
      }

      const { auth } = await import('@/lib/firebase');
      const user = auth.currentUser;
      if (!user) {
        setError('Usuário não autenticado');
        return;
      }

      const token = await user.getIdToken(true);

      // 1. Create property
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_API_URL}/admin/${tenantId}/properties`,
        {
          method: 'POST',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
          body: JSON.stringify(formData),
        }
      );

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Erro ao criar imóvel');
      }

      const result = await response.json();
      const propertyId = result.data.id;

      // 2. Upload photos if any
      if (photos.length > 0) {
        const photosSuccess = await uploadPhotos(propertyId, tenantId, token);

        if (!photosSuccess) {
          setError('Imóvel criado, mas algumas fotos falharam no upload. Você pode adicionar as fotos depois.');
        }
      }

      // 3. Redirect to property details
      router.push(`/dashboard/imoveis/${propertyId}`);
    } catch (err: any) {
      console.error('Erro ao salvar:', err);
      setError(err.message);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="p-6">
      {/* Header */}
      <div className="mb-6">
        <button
          onClick={() => router.push('/dashboard/imoveis')}
          className="flex items-center gap-2 text-gray-600 hover:text-gray-900 mb-4"
        >
          <ArrowLeft className="w-5 h-5" />
          Voltar para lista
        </button>

        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-2xl font-bold text-gray-900 mb-2">Novo Imóvel</h1>
            <p className="text-gray-600">Cadastre um novo imóvel no sistema</p>
          </div>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
          <p className="text-red-600">{error}</p>
        </div>
      )}

      <form onSubmit={handleSubmit}>
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Main Form */}
          <div className="lg:col-span-2 space-y-6">
            {/* Informações Básicas */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Informações Básicas</h2>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Código de Referência
                  </label>
                  <input
                    type="text"
                    value={formData.reference || ''}
                    onChange={(e) => handleChange('reference', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Ex: AP-001"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Slug (URL)
                  </label>
                  <input
                    type="text"
                    value={formData.slug || ''}
                    onChange={(e) => handleChange('slug', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="apartamento-centro-sp"
                  />
                </div>

                <div className="col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Descrição
                  </label>
                  <textarea
                    value={formData.description || ''}
                    onChange={(e) => handleChange('description', e.target.value)}
                    rows={4}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Descreva o imóvel..."
                  />
                </div>
              </div>
            </div>

            {/* Localização */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Localização</h2>

              <div className="grid grid-cols-2 gap-4">
                <div className="col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Endereço
                  </label>
                  <input
                    type="text"
                    value={formData.street || ''}
                    onChange={(e) => handleChange('street', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Rua, número"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Bairro
                  </label>
                  <input
                    type="text"
                    value={formData.neighborhood || ''}
                    onChange={(e) => handleChange('neighborhood', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Cidade
                  </label>
                  <input
                    type="text"
                    value={formData.city || ''}
                    onChange={(e) => handleChange('city', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Estado
                  </label>
                  <input
                    type="text"
                    value={formData.state || ''}
                    onChange={(e) => handleChange('state', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="SP"
                    maxLength={2}
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    CEP
                  </label>
                  <input
                    type="text"
                    value={formData.zip_code || ''}
                    onChange={(e) => handleChange('zip_code', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="00000-000"
                  />
                </div>
              </div>
            </div>

            {/* Características */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Características</h2>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Quartos
                  </label>
                  <input
                    type="number"
                    value={formData.bedrooms || ''}
                    onChange={(e) => handleChange('bedrooms', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="0"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Banheiros
                  </label>
                  <input
                    type="number"
                    value={formData.bathrooms || ''}
                    onChange={(e) => handleChange('bathrooms', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="0"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Suítes
                  </label>
                  <input
                    type="number"
                    value={formData.suites || ''}
                    onChange={(e) => handleChange('suites', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="0"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Vagas
                  </label>
                  <input
                    type="number"
                    value={formData.parking_spaces || ''}
                    onChange={(e) => handleChange('parking_spaces', parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="0"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Área Total (m²)
                  </label>
                  <input
                    type="number"
                    value={formData.total_area || ''}
                    onChange={(e) => handleChange('total_area', parseFloat(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="0"
                    step="0.01"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Área Construída (m²)
                  </label>
                  <input
                    type="number"
                    value={formData.built_area || ''}
                    onChange={(e) => handleChange('built_area', parseFloat(e.target.value) || 0)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    min="0"
                    step="0.01"
                  />
                </div>
              </div>
            </div>

            {/* Fotos */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Fotos do Imóvel</h2>
              <p className="text-sm text-gray-600 mb-4">
                As fotos serão processadas automaticamente em 3 tamanhos (miniatura, média e grande)
              </p>

              <div className="mb-4">
                <label className="flex items-center justify-center w-full px-4 py-6 border-2 border-dashed border-gray-300 rounded-lg cursor-pointer hover:border-blue-500 hover:bg-blue-50 transition-colors">
                  <div className="flex flex-col items-center">
                    <Upload className="w-8 h-8 text-gray-400 mb-2" />
                    <span className="text-sm text-gray-600">Clique para selecionar fotos</span>
                    <span className="text-xs text-gray-500 mt-1">JPEG, PNG ou WebP (máx 10MB cada)</span>
                  </div>
                  <input
                    type="file"
                    multiple
                    accept="image/jpeg,image/png,image/webp"
                    onChange={handlePhotoSelect}
                    className="hidden"
                  />
                </label>
              </div>

              {photos.length > 0 && (
                <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-4">
                  {photos.map((photo, index) => (
                    <div key={index} className="relative group">
                      <img
                        src={photo.preview}
                        alt={`Preview ${index + 1}`}
                        className="w-full h-32 object-cover rounded-lg"
                      />

                      {photo.uploading && (
                        <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center rounded-lg">
                          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-white"></div>
                        </div>
                      )}

                      {photo.uploaded && (
                        <div className="absolute top-2 left-2 bg-green-500 text-white text-xs px-2 py-1 rounded">
                          ✓ Enviada
                        </div>
                      )}

                      {photo.error && (
                        <div className="absolute inset-0 bg-red-500 bg-opacity-75 flex items-center justify-center rounded-lg p-2">
                          <span className="text-white text-xs text-center">{photo.error}</span>
                        </div>
                      )}

                      <button
                        type="button"
                        onClick={() => removePhoto(index)}
                        className="absolute top-2 right-2 bg-red-500 text-white p-1 rounded-full opacity-0 group-hover:opacity-100 transition-opacity"
                        disabled={photo.uploading}
                      >
                        <X className="w-4 h-4" />
                      </button>

                      <div className="absolute bottom-2 left-2 bg-black bg-opacity-75 text-white text-xs px-2 py-1 rounded">
                        {index === 0 ? 'Capa' : `#${index + 1}`}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Sidebar */}
          <div className="lg:col-span-1 space-y-6">
            {/* Tipo e Status */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Configurações</h2>

              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Tipo de Imóvel *
                  </label>
                  <select
                    value={formData.property_type}
                    onChange={(e) => handleChange('property_type', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  >
                    <option value="apartment">Apartamento</option>
                    <option value="house">Casa</option>
                    <option value="condo">Condomínio</option>
                    <option value="penthouse">Cobertura</option>
                    <option value="land">Terreno</option>
                    <option value="commercial">Comercial</option>
                    <option value="farm">Fazenda</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Tipo de Transação *
                  </label>
                  <select
                    value={formData.transaction_type}
                    onChange={(e) => handleChange('transaction_type', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  >
                    <option value="sale">Venda</option>
                    <option value="rent">Aluguel</option>
                    <option value="both">Venda e Aluguel</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Status *
                  </label>
                  <select
                    value={formData.status}
                    onChange={(e) => handleChange('status', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  >
                    <option value="available">Disponível</option>
                    <option value="rented">Alugado</option>
                    <option value="sold">Vendido</option>
                    <option value="unavailable">Indisponível</option>
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Visibilidade *
                  </label>
                  <select
                    value={formData.visibility}
                    onChange={(e) => handleChange('visibility', e.target.value)}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                    required
                  >
                    <option value="public">Público</option>
                    <option value="private">Privado</option>
                    <option value="draft">Rascunho</option>
                  </select>
                </div>

                <div className="flex items-center">
                  <input
                    type="checkbox"
                    id="featured"
                    checked={formData.featured || false}
                    onChange={(e) => handleChange('featured', e.target.checked)}
                    className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-2 focus:ring-blue-500"
                  />
                  <label htmlFor="featured" className="ml-2 text-sm font-medium text-gray-700">
                    Imóvel em Destaque
                  </label>
                </div>
              </div>
            </div>

            {/* Preço */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <h2 className="text-xl font-bold text-gray-900 mb-4">Preço</h2>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Valor (R$)
                </label>
                <input
                  type="number"
                  value={formData.price_amount || ''}
                  onChange={(e) => handleChange('price_amount', parseFloat(e.target.value) || 0)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
                  min="0"
                  step="0.01"
                  placeholder="0,00"
                />
              </div>
            </div>

            {/* Actions */}
            <div className="bg-white rounded-lg shadow-sm p-6">
              <div className="flex flex-col gap-3">
                <button
                  type="submit"
                  disabled={saving}
                  className="flex items-center justify-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
                >
                  {saving ? (
                    <>
                      <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                      Salvando...
                    </>
                  ) : (
                    <>
                      <Save className="w-5 h-5" />
                      Salvar Imóvel
                    </>
                  )}
                </button>

                <button
                  type="button"
                  onClick={() => router.push('/dashboard/imoveis')}
                  className="flex items-center justify-center gap-2 bg-gray-200 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-300"
                  disabled={saving}
                >
                  Cancelar
                </button>
              </div>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
}
