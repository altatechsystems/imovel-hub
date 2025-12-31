'use client';

import { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { Users, Search, Filter, MessageSquare, Phone, Mail, UserCheck } from 'lucide-react';
import { adminApi } from '@/lib/api';
import { Lead, LeadStatus, LeadChannel } from '@/types/lead';

type LeadStatusFilter = 'all' | LeadStatus;
type LeadChannelFilter = 'all' | LeadChannel;

export default function LeadsPage() {
  const router = useRouter();
  const [searchTerm, setSearchTerm] = useState('');
  const [leads, setLeads] = useState<Lead[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [displayCount, setDisplayCount] = useState(20);
  const [statusFilter, setStatusFilter] = useState<LeadStatusFilter>('all');
  const [channelFilter, setChannelFilter] = useState<LeadChannelFilter>('all');
  const observerTarget = useRef<HTMLDivElement>(null);
  const itemsPerPage = 20;

  const fetchLeads = useCallback(async () => {
    try {
      setLoading(true);

      const startTime = performance.now();
      const response = await adminApi.getLeads({}, { limit: 1000 });
      const loadTime = performance.now() - startTime;

      console.log(`✅ Loaded ${response.data?.length || 0} leads in ${loadTime.toFixed(0)}ms`);

      setLeads(response.data || []);
    } catch (err: any) {
      console.error('Erro ao buscar leads:', err);
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchLeads();
  }, [fetchLeads]);

  const formatDate = (date: Date | string | undefined) => {
    if (!date) return 'Data não informada';
    const d = typeof date === 'string' ? new Date(date) : date;
    return new Intl.DateTimeFormat('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(d);
  };

  // Calculate stats
  const stats = useMemo(() => {
    const result = {
      total: leads.length,
      new: 0,
      contacted: 0,
      qualified: 0,
      negotiating: 0,
      converted: 0,
      lost: 0,
      whatsapp: 0,
      form: 0,
      phone: 0,
      email: 0,
    };

    leads.forEach(lead => {
      // Status counts
      if (lead.status === LeadStatus.NEW) result.new++;
      else if (lead.status === LeadStatus.CONTACTED) result.contacted++;
      else if (lead.status === LeadStatus.QUALIFIED) result.qualified++;
      else if (lead.status === LeadStatus.NEGOTIATING) result.negotiating++;
      else if (lead.status === LeadStatus.CONVERTED) result.converted++;
      else if (lead.status === LeadStatus.LOST) result.lost++;

      // Channel counts
      if (lead.channel === LeadChannel.WHATSAPP) result.whatsapp++;
      else if (lead.channel === LeadChannel.FORM) result.form++;
      else if (lead.channel === LeadChannel.PHONE) result.phone++;
      else if (lead.channel === LeadChannel.EMAIL) result.email++;
    });

    return result;
  }, [leads]);

  const filteredLeads = useMemo(() => {
    let filtered = leads;

    // Apply status filter
    if (statusFilter !== 'all') {
      filtered = filtered.filter(lead => lead.status === statusFilter);
    }

    // Apply channel filter
    if (channelFilter !== 'all') {
      filtered = filtered.filter(lead => lead.channel === channelFilter);
    }

    // Apply search term filter
    if (searchTerm) {
      filtered = filtered.filter(lead =>
        lead.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        lead.email?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        lead.phone?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        lead.id?.toLowerCase().includes(searchTerm.toLowerCase())
      );
    }

    return filtered;
  }, [leads, searchTerm, statusFilter, channelFilter]);

  // Infinite scroll - show only displayCount items
  const displayedLeads = useMemo(() =>
    filteredLeads.slice(0, displayCount),
    [filteredLeads, displayCount]
  );

  const hasMore = displayCount < filteredLeads.length;

  // Reset display count when search or filter changes
  useEffect(() => {
    setDisplayCount(20);
  }, [searchTerm, statusFilter, channelFilter]);

  // Infinite scroll observer
  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loading) {
          setDisplayCount(prev => prev + itemsPerPage);
        }
      },
      { threshold: 0.1 }
    );

    const currentTarget = observerTarget.current;
    if (currentTarget) {
      observer.observe(currentTarget);
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget);
      }
    };
  }, [hasMore, loading, itemsPerPage]);

  const getStatusColor = (status?: LeadStatus) => {
    switch (status) {
      case LeadStatus.NEW: return 'bg-blue-100 text-blue-700';
      case LeadStatus.CONTACTED: return 'bg-yellow-100 text-yellow-700';
      case LeadStatus.QUALIFIED: return 'bg-purple-100 text-purple-700';
      case LeadStatus.NEGOTIATING: return 'bg-orange-100 text-orange-700';
      case LeadStatus.CONVERTED: return 'bg-green-100 text-green-700';
      case LeadStatus.LOST: return 'bg-red-100 text-red-700';
      default: return 'bg-gray-100 text-gray-700';
    }
  };

  const getStatusLabel = (status?: LeadStatus) => {
    switch (status) {
      case LeadStatus.NEW: return 'Novo';
      case LeadStatus.CONTACTED: return 'Contatado';
      case LeadStatus.QUALIFIED: return 'Qualificado';
      case LeadStatus.NEGOTIATING: return 'Negociando';
      case LeadStatus.CONVERTED: return 'Convertido';
      case LeadStatus.LOST: return 'Perdido';
      default: return 'Novo';
    }
  };

  const getChannelIcon = (channel: LeadChannel) => {
    switch (channel) {
      case LeadChannel.WHATSAPP: return <MessageSquare className="w-4 h-4" />;
      case LeadChannel.FORM: return <Mail className="w-4 h-4" />;
      case LeadChannel.PHONE: return <Phone className="w-4 h-4" />;
      case LeadChannel.EMAIL: return <Mail className="w-4 h-4" />;
      default: return <UserCheck className="w-4 h-4" />;
    }
  };

  const getChannelLabel = (channel: LeadChannel) => {
    switch (channel) {
      case LeadChannel.WHATSAPP: return 'WhatsApp';
      case LeadChannel.FORM: return 'Formulário';
      case LeadChannel.PHONE: return 'Telefone';
      case LeadChannel.EMAIL: return 'Email';
      case LeadChannel.CHAT: return 'Chat';
      case LeadChannel.REFERRAL: return 'Indicação';
      default: return channel;
    }
  };

  return (
    <div className="p-3 sm:p-4 md:p-6">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">Leads</h1>
        <p className="text-gray-600">Gerencie todos os contatos e oportunidades</p>
      </div>

      {/* Stats Cards - Status */}
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-7 gap-2 sm:gap-3 mb-4">
        <button
          onClick={() => setStatusFilter('all')}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            statusFilter === 'all' ? 'ring-2 ring-blue-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-blue-100 rounded flex items-center justify-center flex-shrink-0">
              <Users className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-blue-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Total</p>
              <p className="text-base sm:text-lg font-bold text-gray-900">{stats.total}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setStatusFilter(LeadStatus.NEW)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            statusFilter === LeadStatus.NEW ? 'ring-2 ring-blue-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-blue-100 rounded flex items-center justify-center flex-shrink-0">
              <Users className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-blue-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Novos</p>
              <p className="text-base sm:text-lg font-bold text-blue-600">{stats.new}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setStatusFilter(LeadStatus.CONTACTED)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            statusFilter === LeadStatus.CONTACTED ? 'ring-2 ring-yellow-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-yellow-100 rounded flex items-center justify-center flex-shrink-0">
              <UserCheck className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-yellow-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Contatados</p>
              <p className="text-base sm:text-lg font-bold text-yellow-600">{stats.contacted}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setStatusFilter(LeadStatus.QUALIFIED)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            statusFilter === LeadStatus.QUALIFIED ? 'ring-2 ring-purple-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-purple-100 rounded flex items-center justify-center flex-shrink-0">
              <UserCheck className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-purple-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Qualificados</p>
              <p className="text-base sm:text-lg font-bold text-purple-600">{stats.qualified}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setStatusFilter(LeadStatus.NEGOTIATING)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            statusFilter === LeadStatus.NEGOTIATING ? 'ring-2 ring-orange-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-orange-100 rounded flex items-center justify-center flex-shrink-0">
              <MessageSquare className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-orange-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Negociando</p>
              <p className="text-base sm:text-lg font-bold text-orange-600">{stats.negotiating}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setStatusFilter(LeadStatus.CONVERTED)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            statusFilter === LeadStatus.CONVERTED ? 'ring-2 ring-green-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-green-100 rounded flex items-center justify-center flex-shrink-0">
              <UserCheck className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-green-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Convertidos</p>
              <p className="text-base sm:text-lg font-bold text-green-600">{stats.converted}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setStatusFilter(LeadStatus.LOST)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            statusFilter === LeadStatus.LOST ? 'ring-2 ring-red-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-red-100 rounded flex items-center justify-center flex-shrink-0">
              <Users className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-red-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Perdidos</p>
              <p className="text-base sm:text-lg font-bold text-red-600">{stats.lost}</p>
            </div>
          </div>
        </button>
      </div>

      {/* Channel Filters */}
      <div className="grid grid-cols-2 sm:grid-cols-4 gap-2 sm:gap-3 mb-4">
        <button
          onClick={() => setChannelFilter('all')}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            channelFilter === 'all' ? 'ring-2 ring-blue-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-gray-100 rounded flex items-center justify-center flex-shrink-0">
              <Users className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-gray-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Todos Canais</p>
              <p className="text-base sm:text-lg font-bold text-gray-900">{stats.total}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setChannelFilter(LeadChannel.WHATSAPP)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            channelFilter === LeadChannel.WHATSAPP ? 'ring-2 ring-green-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-green-100 rounded flex items-center justify-center flex-shrink-0">
              <MessageSquare className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-green-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">WhatsApp</p>
              <p className="text-base sm:text-lg font-bold text-green-600">{stats.whatsapp}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setChannelFilter(LeadChannel.FORM)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            channelFilter === LeadChannel.FORM ? 'ring-2 ring-blue-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-blue-100 rounded flex items-center justify-center flex-shrink-0">
              <Mail className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-blue-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Formulário</p>
              <p className="text-base sm:text-lg font-bold text-blue-600">{stats.form}</p>
            </div>
          </div>
        </button>

        <button
          onClick={() => setChannelFilter(LeadChannel.PHONE)}
          className={`bg-white rounded-lg shadow-sm p-2 sm:p-3 text-left transition-all hover:shadow-md ${
            channelFilter === LeadChannel.PHONE ? 'ring-2 ring-purple-500' : ''
          }`}
        >
          <div className="flex items-center gap-1.5 sm:gap-2">
            <div className="w-7 h-7 sm:w-8 sm:h-8 bg-purple-100 rounded flex items-center justify-center flex-shrink-0">
              <Phone className="w-3.5 h-3.5 sm:w-4 sm:h-4 text-purple-600" />
            </div>
            <div className="min-w-0">
              <p className="text-[10px] sm:text-xs text-gray-600 truncate">Telefone</p>
              <p className="text-base sm:text-lg font-bold text-purple-600">{stats.phone}</p>
            </div>
          </div>
        </button>
      </div>

      {/* Actions Bar */}
      <div className="bg-white rounded-lg shadow-sm p-3 sm:p-4 mb-4 sm:mb-6">
        <div className="flex flex-col sm:flex-row gap-2 sm:gap-4">
          {/* Search */}
          <div className="flex-1 relative">
            <Search className="absolute left-2.5 sm:left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4 sm:w-5 sm:h-5" />
            <input
              type="text"
              placeholder="Buscar por nome, email, telefone ou ID..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="w-full pl-8 sm:pl-10 pr-3 sm:pr-4 py-2 text-sm sm:text-base border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          {/* Filter Button */}
          <button className="flex items-center justify-center gap-2 px-3 sm:px-4 py-2 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors text-sm sm:text-base">
            <Filter className="w-4 h-4 sm:w-5 sm:h-5" />
            <span className="hidden sm:inline">Filtros</span>
          </button>
        </div>
      </div>

      {/* Loading State */}
      {loading && (
        <div className="bg-white rounded-lg shadow-sm p-8">
          <div className="animate-pulse space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="flex items-center gap-4">
                <div className="w-12 h-12 bg-gray-200 rounded-full"></div>
                <div className="flex-1">
                  <div className="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
                  <div className="h-3 bg-gray-200 rounded w-1/2"></div>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Error State */}
      {error && !loading && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
          <p className="text-red-600">{error}</p>
        </div>
      )}

      {/* Empty State */}
      {!loading && !error && filteredLeads.length === 0 && (
        <div className="bg-white rounded-lg shadow-sm p-12">
          <div className="text-center">
            <div className="inline-flex items-center justify-center w-16 h-16 bg-blue-100 rounded-full mb-4">
              <Users className="w-8 h-8 text-blue-600" />
            </div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">
              {searchTerm ? 'Nenhum lead encontrado' : 'Nenhum lead cadastrado'}
            </h3>
            <p className="text-gray-600">
              {searchTerm
                ? 'Tente buscar com outros termos'
                : 'Os leads começarão a aparecer quando visitantes entrarem em contato através do site'}
            </p>
          </div>
        </div>
      )}

      {/* Leads Table */}
      {!loading && !error && filteredLeads.length > 0 && (
        <>
          <div className="bg-white rounded-lg shadow-sm overflow-hidden mb-6">
            <div className="overflow-x-auto">
              <table className="w-full">
                <thead className="bg-gray-50 border-b border-gray-200">
                  <tr>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">
                      Lead
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">
                      Contato
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">
                      Canal
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">
                      Status
                    </th>
                    <th className="px-4 py-3 text-left text-xs font-medium text-gray-600 uppercase tracking-wider">
                      Data
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-200">
                  {displayedLeads.map((lead) => (
                    <tr
                      key={lead.id}
                      onClick={() => router.push(`/dashboard/leads/${lead.id}`)}
                      className="hover:bg-gray-50 cursor-pointer transition-colors"
                    >
                      <td className="px-4 py-4">
                        <div>
                          <p className="font-medium text-gray-900">{lead.name || 'Nome não informado'}</p>
                          <p className="text-sm text-gray-500">ID: {lead.id?.slice(0, 8)}...</p>
                        </div>
                      </td>
                      <td className="px-4 py-4">
                        <div className="text-sm">
                          {lead.email && (
                            <p className="text-gray-900 flex items-center gap-1">
                              <Mail className="w-3 h-3 text-gray-400" />
                              {lead.email}
                            </p>
                          )}
                          {lead.phone && (
                            <p className="text-gray-600 flex items-center gap-1">
                              <Phone className="w-3 h-3 text-gray-400" />
                              {lead.phone}
                            </p>
                          )}
                        </div>
                      </td>
                      <td className="px-4 py-4">
                        <div className="flex items-center gap-2">
                          {getChannelIcon(lead.channel)}
                          <span className="text-sm text-gray-700">{getChannelLabel(lead.channel)}</span>
                        </div>
                      </td>
                      <td className="px-4 py-4">
                        <span className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(lead.status)}`}>
                          {getStatusLabel(lead.status)}
                        </span>
                      </td>
                      <td className="px-4 py-4">
                        <p className="text-sm text-gray-600">{formatDate(lead.created_at)}</p>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>

          {/* Infinite Scroll Trigger & Loading Indicator */}
          <div ref={observerTarget} className="py-8">
            {hasMore && (
              <div className="text-center">
                <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                <p className="text-gray-600 mt-2">Carregando mais leads...</p>
              </div>
            )}
            {!hasMore && filteredLeads.length > 20 && (
              <div className="text-center">
                <p className="text-gray-600">
                  Mostrando todos os {filteredLeads.length} leads
                </p>
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}
