  // apiClient.ts
  import axios from "axios";

  // Base URL
  const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:4000";

  // ---------- Interfaces ----------
  export interface AuthResponse {
    token: string;
    user: {
      id: string;
      email: string;
      first_name: string;
      last_name: string;
      role: "guest" | "admin";
    };
  }

  export interface Property {
    id: string;
    title: string;
    description: string;
    location: string;
    price_per_night: number;
    max_guests: number;
    created_at: string;
    updated_at: string;
    thumbnail_url: {
      String: string;
      Valid: boolean;
    };
    user_id: string;
  }

  export interface PropertyImage {
    id: string;
    property_id: string;
    image_url: string;
    caption?: string;
    display_order: number;
  }

  export interface PropertyDetail {
    id: string;
    title: string;
    description: string;
    location: string;
    price_per_night: number;
    max_guests: number;
    created_at: string;
    updated_at: string;
    images: PropertyImage[];
    amenities: Amenity[];
  }

  export interface Booking {
    id: string;
    property_id: string;
    user_id: string;
    start_date: string;
    end_date: string;
    total_price: number;
    status: "booked" | "cancelled";
    created_at: string;
  }

  export interface Amenity {
    amenity_id: string;
    name: string;
  }

  // ---------- Axios Instance ----------
  export const api = axios.create({
    baseURL: API_BASE_URL,
    headers: { "Content-Type": "application/json" },
  });

  // Token handler
  export const setAuthToken = (token: string | null) => {
    if (token) {
      api.defaults.headers.common["Authorization"] = `Bearer ${token}`;
    } else {
      delete api.defaults.headers.common["Authorization"];
    }
  };

  export const clearToken = () => {
    delete api.defaults.headers.common["Authorization"];
    // Optional: clear from storage if you use localStorage
    if (typeof window !== "undefined") {
      localStorage.removeItem("token");
    }
  };

  // Error interceptor
  api.interceptors.response.use(
    (res) => res,
    (error) => {
      const msg =
        error.response?.data?.error?.message ||
        error.response?.data?.message ||
        error.message ||
        "Unexpected API error";
      return Promise.reject(new Error(msg));
    }
  );

  // ---------- Auth API ----------
  export const authApi = {
    register: (
      email: string,
      password: string,
      firstName: string,
      lastName: string
    ) =>
      api.post<AuthResponse>("/auth/register", {
        email,
        password_hash: password,
        first_name: firstName,
        last_name: lastName,
      }),

    login: (email: string, password: string) =>
      api.post<AuthResponse>("/auth/login", { email, password_hash: password }),

    getMe: () => api.get<{ user: AuthResponse["user"] }>("/auth/me"),
  };

  // ---------- Properties API ----------
  export const propertiesApi = {
    getAll: () => api.get<Property[]>("/properties"),

    getById: (id: string) => api.get<PropertyDetail>(`/properties/${id}`),

    create: (data: {
      title: string;
      description: string;
      location: string;
      price_per_night: number;
      max_guests: number;
      image_url?: string;
    }) => api.post<{ id: string; message: string }>("/properties", data),

    update: (data: {
      id: string;
      title: string;
      description: string;
      location: string;
      price_per_night: number;
      max_guests: number;
    }) => api.put<{ message: string; userID: string }>(`/properties`, data),

    delete: (id: string) => api.delete<void>(`/properties/${id}`),

    getAvailability: (id: string, start_date: string, end_date: string) =>
      api.get<{ available: boolean }>(`/properties/${id}/availability`, {
        params: { start_date, end_date },
      }),

    addImage: (
      propertyId: string,
      images: Array<{
        image_url: string;
        caption?: string;
        display_order: number;
      }>
    ) => api.post<{ message: string }>(`/property/image/${propertyId}`, { images }),

    deleteImage: (imageId: string) => api.delete<void>(`/property/image/${imageId}`),
  };

  // ---------- Bookings API ----------
  export const bookingsApi = {
    getAll: () => api.get<Booking[]>("/bookings"),

    getById: (id: string) => api.get<Booking>(`/bookings/${id}`),

    create: (data: {
      property_id: string;
      start_date: string;
      end_date: string;
      total_price: number;
    }) => api.post<Booking>("/bookings", data),

    cancel: (id: string) => api.patch<Booking>(`/bookings/${id}`),
  };

  // ---------- Amenities API ----------
  export const amenitiesApi = {
    getAll: () => api.get<Amenity[]>("/amenities"),

    create: (data: { name: string }) => api.post<Amenity>("/amenities", data),

    addToProperty: (propertyId: string, amenityIds: string[]) =>
      api.post<void>(`/properties/${propertyId}/amenities`, {
        amenity_id: amenityIds,
      }),
  };
