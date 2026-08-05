// Mirrors the backend's unified response envelope.

export interface Envelope<T> {
  success: boolean;
  data?: T;
  error?: ApiError;
  meta?: PaginationMeta;
}

export interface ApiError {
  code: string;
  message: string;
  fields?: Record<string, string>;
}

export interface PaginationMeta {
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}

export interface Page<T> {
  items: T[];
  meta: PaginationMeta;
}

// --- Domain types (subset used by the admin) ---

export interface User {
  id: number;
  name: string;
  email: string;
  avatar: string;
  phone: string;
  locale: string;
  email_verified: boolean;
  phone_verified: boolean;
  roles: string[];
}

export interface AuthResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  expires_in: number;
}

export interface Language {
  id: number;
  code: string;
  name: string;
  native_name: string;
  is_rtl: boolean;
  is_active: boolean;
  is_default: boolean;
  sort_order: number;
}

export interface Agency {
  id: number;
  slug: string;
  name: string;
  logo: string;
  phone: string;
  email: string;
  verified: boolean;
  status: "pending" | "approved" | "suspended";
  commission_rate: number;
  rating_avg: number;
  rating_count: number;
  created_at: string;
}

export interface Setting {
  key: string;
  value: string;
  type: string;
}

export interface Banner {
  id: number;
  type: string;
  image: string;
  title: string;
  link: string;
  sort_order: number;
  active: boolean;
}

export interface Review {
  id: number;
  tour_id: number;
  user_id: number;
  rating: number;
  comment: string;
  status: string;
}

export interface Boost {
  id: number;
  agency_id: number;
  tour_id: number | null;
  package_id: number;
  placement: string;
  status: string;
  starts_at: string | null;
  ends_at: string | null;
  impressions: number;
  clicks: number;
}

export interface AnalyticsSummary {
  users: number;
  agencies: number;
  tours: number;
  bookings: number;
  active_boosts: number;
  revenue: Record<string, number>;
}

export interface Category {
  id: number;
  slug: string;
  icon: string;
  name: string;
}

export interface Region {
  id: number;
  slug: string;
  name: string;
  parent_id: number | null;
}
