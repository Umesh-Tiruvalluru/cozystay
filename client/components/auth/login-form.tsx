// components/auth/LoginForm.tsx
"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { loginSchema } from "@/lib/schemas/authSchemas";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import FormInput from "@/components/ui/FormInput";

type LoginFormProps = {
  onSubmit: (values: z.infer<typeof loginSchema>) => void;
  isLoading: boolean;
};

export default function LoginForm({ onSubmit, isLoading }: LoginFormProps) {
  const form = useForm<z.infer<typeof loginSchema>>({
    resolver: zodResolver(loginSchema),
    defaultValues: {
      email: "",
      password: "",
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 mt-4">
        <FormInput
          control={form.control}
          name="email"
          label="Email"
          placeholder="you@example.com"
          id="login-email"
        />
        <FormField
          control={form.control}
          name="password"
          render={({ field }) => (
            <FormItem>
              <div className="flex justify-between items-center">
                <FormLabel>Password</FormLabel>
                <a
                  className="text-sm text-primary hover:underline"
                  href="/forgot-password"
                >
                  Forgot password?
                </a>
              </div>
              <FormControl>
                <FormInput
                  control={form.control}
                  name="password"
                  label="Password"
                  placeholder="••••••••"
                  isPassword
                  id="login-password"
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button
          type="submit"
          disabled={isLoading}
          className="w-full bg-primary text-white hover:bg-primary/90"
        >
          {isLoading ? "Logging in..." : "Log In"}
        </Button>
      </form>
    </Form>
  );
}
