// lib/api/authService.ts
import axios from "axios";

export const authService = {
  login: async (data: { email: string; password: string }) => {
    try {
      const response = await axios.post("/api/auth/login", data); // Adjust URL to your backend
      return response.data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "Login failed");
    }
  },
  register: async (data: {
    firstName: string;
    lastName: string;
    email: string;
    password: string;
  }) => {
    try {
      const response = await axios.post("/api/auth/register", data); // Adjust URL to your backend
      return response.data;
    } catch (error: any) {
      throw new Error(error.response?.data?.message || "Registration failed");
    }
  },
};
