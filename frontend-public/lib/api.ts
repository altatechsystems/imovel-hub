import axios, { AxiosInstance } from 'axios';
import {
  Property,
  PropertyFilters,
  PropertyListResponse,
  PropertyResponse,
  PaginationOptions,
} from '@/types/property';
import { CreateLeadRequest, CreateLeadResponse } from '@/types/lead';

class ApiClient {
  private client: AxiosInstance;
  private tenantId: string;

  constructor() {
    this.tenantId = process.env.NEXT_PUBLIC_TENANT_ID || '';

    this.client = axios.create({
      baseURL: process.env.NEXT_PUBLIC_API_URL,
      headers: {
        'Content-Type': 'application/json',
      },
      timeout: 10000,
    });

    // Request interceptor
    this.client.interceptors.request.use(
      (config) => {
        // Add auth token if available
        const token = this.getAuthToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );

    // Response interceptor
    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        console.error('API Error:', error.response?.data || error.message);
        return Promise.reject(error);
      }
    );
  }

  private getAuthToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem('auth_token');
  }

  // Properties
  async getProperties(
    filters?: PropertyFilters,
    pagination?: PaginationOptions
  ): Promise<PropertyListResponse> {
    const params = new URLSearchParams();

    // Add filters
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          params.append(key, String(value));
        }
      });
    }

    // Add pagination
    if (pagination) {
      if (pagination.limit) params.append('limit', String(pagination.limit));
      if (pagination.order_by) params.append('order_by', pagination.order_by);
      if (pagination.order_direction) params.append('order_direction', pagination.order_direction);
      if (pagination.start_after) params.append('start_after', pagination.start_after);
    }

    const response = await this.client.get<PropertyListResponse>(
      `/${this.tenantId}/properties?${params.toString()}`
    );
    return response.data;
  }

  async getProperty(id: string): Promise<Property> {
    const response = await this.client.get<PropertyResponse>(
      `/${this.tenantId}/properties/${id}`
    );
    return response.data.data;
  }

  async getPropertyBySlug(slug: string): Promise<Property> {
    const response = await this.client.get<PropertyResponse>(
      `/${this.tenantId}/properties/slug/${slug}`
    );
    return response.data.data;
  }

  async getPropertyImages(propertyId: string): Promise<any> {
    const response = await this.client.get(
      `/${this.tenantId}/properties/${propertyId}/images`
    );
    return response.data.data;
  }

  // Leads
  async createLead(data: CreateLeadRequest): Promise<CreateLeadResponse> {
    const response = await this.client.post<CreateLeadResponse>(
      `/${this.tenantId}/leads`,
      data
    );
    return response.data;
  }

  // Search helpers
  async searchProperties(query: string, filters?: PropertyFilters): Promise<Property[]> {
    const searchFilters = {
      ...filters,
      // Add search query to city or neighborhood
    };

    const result = await this.getProperties(searchFilters);
    return result.data;
  }

  // Featured properties
  async getFeaturedProperties(limit: number = 6): Promise<Property[]> {
    const result = await this.getProperties(
      { featured: true, status: 'available' as any, visibility: 'public' as any },
      { limit }
    );
    return result.data;
  }

  // Similar properties
  async getSimilarProperties(
    propertyId: string,
    limit: number = 4
  ): Promise<Property[]> {
    // Get the property first to match similar ones
    const property = await this.getProperty(propertyId);

    const result = await this.getProperties(
      {
        property_type: property.property_type,
        city: property.city,
        status: 'available' as any,
        visibility: 'public' as any,
      },
      { limit: limit + 1 } // +1 to exclude current property
    );

    // Filter out the current property
    return result.data.filter(p => p.id !== propertyId).slice(0, limit);
  }
}

export const api = new ApiClient();
export default api;
