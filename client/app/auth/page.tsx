// app/auth/page.tsx
"use client";

import { useState, useEffect } from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import LoginForm from "@/components/auth/login-form";
import RegisterForm from "@/components/auth/register-form";
import SocialAuthButtons from "@/components/auth/social-auth-buttons";
import Logo from "@/components/logo";
import { authService } from "@/lib/api/auth-service";
import { loginSchema, registerSchema } from "@/lib/schemas/auth-schema";
import { toast } from "sonner";
import z from "zod";

export default function AuthPage() {
  const [tab, setTab] = useState<"login" | "register">("login");
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    const firstInput =
      tab === "login"
        ? document.getElementById("login-email")
        : document.getElementById("register-first-name");
    firstInput?.focus();
  }, [tab]);

  const handleLoginSubmit = async (values: z.infer<typeof loginSchema>) => {
    setIsLoading(true);
    try {
      const result = await authService.login(values);
      toast.success("Logged in successfully!");
      // Redirect or handle success (e.g., router.push("/dashboard"))
    } catch (error: any) {
      toast.error(error.message || "Login failed");
    } finally {
      setIsLoading(false);
    }
  };

  const handleRegisterSubmit = async (
    values: z.infer<typeof registerSchema>
  ) => {
    setIsLoading(true);
    try {
      const result = await authService.register(values);
      toast.success("Account created successfully!");
      setTab("login"); // Switch to login tab after successful registration
    } catch (error: any) {
      toast.error(error.message || "Registration failed");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen w-full items-center justify-center p-4 bg-zinc-200 font-display text-slate-800">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <Logo />
          <p className="text-slate-500 mt-2">
            Welcome! Please log in or create an account.
          </p>
        </div>
        <div className="bg-white rounded-xl shadow-lg p-8">
          <Tabs
            value={tab}
            onValueChange={(value) => setTab(value as "login" | "register")}
          >
            <TabsList className="w-full border-zinc-200 py-5">
              <TabsTrigger value="login" className="cursor-pointer py-4">
                Login
              </TabsTrigger>
              <TabsTrigger value="register" className="cursor-pointer py-4">
                Register
              </TabsTrigger>
            </TabsList>
            <TabsContent value="login">
              <LoginForm onSubmit={handleLoginSubmit} isLoading={isLoading} />
            </TabsContent>
            <TabsContent value="register">
              <RegisterForm
                onSubmit={handleRegisterSubmit}
                isLoading={isLoading}
              />
            </TabsContent>
          </Tabs>
          <SocialAuthButtons />
        </div>
      </div>
    </div>
  );
}
