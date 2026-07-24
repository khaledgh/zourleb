// Mirrors the backend's unified response envelope and domain DTOs.

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
  code: string;
  name: string;
  native_name: string;
  is_rtl: boolean;
  is_default: boolean;
}

export interface AgencyMini {
  id: number;
  slug: string;
  name: string;
  logo: string;
  verified: boolean;
}

export interface TourCard {
  id: number;
  slug: string;
  title: string;
  summary: string;
  cover: string;
  type: string;
  duration_days: number;
  difficulty: string;
  price_from: number;
  currency: string;
  featured: boolean;
  region_id: number | null;
  category_id: number | null;
  agency: AgencyMini;
  rating_avg: number;
}

export interface TourImage {
  url: string;
  sort_order: number;
  is_cover: boolean;
}

export interface TourDeparture {
  id: number;
  start_date: string;
  end_date: string;
  capacity: number;
  seats_left: number;
  status: string;
}

export interface TourPrice {
  departure_id: number | null;
  traveler_type: string;
  currency: string;
  amount: number;
}

export interface TourDetail extends TourCard {
  description: string;
  itinerary: string;
  included: string;
  excluded: string;
  min_age: number;
  max_capacity: number;
  images: TourImage[];
  past_gallery: TourImage[];
  departures: TourDeparture[];
  prices: TourPrice[];
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

export interface Banner {
  id: number;
  image: string;
  title: string;
  link: string;
  tour_id: number | null;
}

export interface ModuleFlags {
  shop: boolean;
  reviews: boolean;
  boost: boolean;
}

export interface HomePayload {
  banners: Banner[];
  featured: TourCard[];
  categories: Category[];
  regions: Region[];
  last_minute: TourCard[];
  modules: ModuleFlags;
}

export interface Booking {
  id: number;
  code: string;
  tour_id: number;
  tour_title: string;
  departure_id: number | null;
  status: string;
  payment_status: string;
  travelers_count: number;
  subtotal: number;
  currency: string;
  contact_phone: string;
  created_at: string;
}
