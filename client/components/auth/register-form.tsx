// components/auth/RegisterForm.tsx
"use client";

import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { Form } from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { registerSchema } from "@/lib/schemas/auth-schema";
import { FormInput } from "lucide-react";

type RegisterFormProps = {
  onSubmit: (values: z.infer<typeof registerSchema>) => void;
  isLoading: boolean;
};

export default function RegisterForm({
  onSubmit,
  isLoading,
}: RegisterFormProps) {
  const form = useForm<z.infer<typeof registerSchema>>({
    resolver: zodResolver(registerSchema),
    defaultValues: {
      firstName: "",
      lastName: "",
      email: "",
      password: "",
      confirmPassword: "",
    },
  });

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 mt-4">
        <div className="grid grid-cols-2 gap-4">
          <FormInput
            control={form.control}
            name="firstName"
            label="First Name"
            placeholder="John"
            id="register-first-name"
          />
          <FormInput
            control={form.control}
            name="lastName"
            label="Last Name"
            placeholder="Doe"
            id="register-last-name"
          />
        </div>
        <FormInput
          control={form.control}
          name="email"
          label="Email"
          placeholder="you@example.com"
          id="register-email"
        />
        <div className="grid grid-cols-2 gap-4">
          <FormInput
            control={form.control}
            name="password"
            label="Password"
            placeholder="••••••••"
            isPassword
            id="register-password"
          />
          <FormInput
            control={form.control}
            name="confirmPassword"
            label="Confirm Password"
            placeholder="••••••••"
            isPassword
            id="confirm-password"
          />
        </div>
        <Button
          type="submit"
          disabled={isLoading}
          className="w-full bg-primary text-white hover:bg-primary/90"
        >
          {isLoading ? "Creating account..." : "Create Account"}
        </Button>
      </form>
    </Form>
  );
}
