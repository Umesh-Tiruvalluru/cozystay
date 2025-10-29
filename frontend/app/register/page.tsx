"use client";

import type React from "react";

import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { useAuth } from "@/lib/auth-context";

export default function RegisterPage() {
  const router = useRouter();
  const { register } = useAuth();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState("");
  const [formData, setFormData] = useState({
    firstName: "",
    lastName: "",
    email: "",
    password: "",
    confirmPassword: "",
  });

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (formData.password !== formData.confirmPassword) {
      setError("Passwords do not match");
      return;
    }

    setIsLoading(true);
    try {
      await register(
        formData.email,
        formData.password,
        formData.firstName,
        formData.lastName
      );
      router.push("/login");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Registration failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-background px-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Create Account</CardTitle>
          <CardDescription>Join BookBnB and start exploring</CardDescription>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
              <div className="bg-destructive/10 text-destructive p-3 rounded-md text-sm">
                {error}
              </div>
            )}
            <div className="grid grid-cols-2 gap-4">
              <Input
                placeholder="First Name"
                name="firstName"
                value={formData.firstName}
                onChange={handleChange}
                required
              />
              <Input
                placeholder="Last Name"
                name="lastName"
                value={formData.lastName}
                onChange={handleChange}
                required
              />
            </div>
            <Input
              type="email"
              placeholder="Email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              required
            />
            <Input
              type="password"
              placeholder="Password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required
            />
            <Input
              type="password"
              placeholder="Confirm Password"
              name="confirmPassword"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
            />
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Creating Account..." : "Sign Up"}
            </Button>
          </form>
          <p className="text-center text-sm text-muted-foreground mt-4">
            Already have an account?{" "}
            <Link href="/login" className="text-primary hover:underline">
              Login
            </Link>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
