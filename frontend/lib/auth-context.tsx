"use client";

import type React from "react";
import { createContext, useContext, useEffect, useState } from "react";
import { authApi, clearToken, setAuthToken } from "./api-client";

interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: "guest" | "admin";
}

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (
    email: string,
    password: string,
    firstName: string,
    lastName: string
  ) => Promise<void>;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Initialize auth from localStorage
  useEffect(() => {
    const token = localStorage.getItem("auth_token");
    if (token) {
      setAuthToken(token);
      authApi
        .getMe()
        .then((response) => {
          setUser(response.data.user);
        })
        .catch(() => {
          localStorage.removeItem("auth_token");
          clearToken();
        })
        .finally(() => {
          setIsLoading(false);
        });
    } else {
      setIsLoading(false);
    }
  }, []);

  const login = async (email: string, password: string) => {
    const response = await authApi.login(email, password);
    localStorage.setItem("auth_token", response.data.token);
    setAuthToken(response.data.token);
    setUser(response.data.user);
  };

  const register = async (
    email: string,
    password: string,
    firstName: string,
    lastName: string
  ) => {
    const response = await authApi.register(
      email,
      password,
      firstName,
      lastName
    );
    localStorage.setItem("auth_token", response.data.token);
    setAuthToken(response.data.token);
    setUser(response.data.user);
  };

  const logout = () => {
    localStorage.removeItem("auth_token");
    clearToken();
    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        isLoading,
        isAuthenticated: !!user,
        login,
        register,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return context;
}
