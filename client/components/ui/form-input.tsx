// components/ui/FormInput.tsx
"use client";

import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Eye, EyeOff } from "lucide-react";
import { usePasswordVisibility } from "@/hooks/use-password-vsibility";

type FormInputProps = {
  control: any;
  name: string;
  label: string;
  placeholder?: string;
  type?: string;
  isPassword?: boolean;
  id?: string;
};

export default function FormInput({
  control,
  name,
  label,
  placeholder,
  type = "text",
  isPassword = false,
  id,
}: FormInputProps) {
  const { showPassword, togglePassword } = usePasswordVisibility();

  return (
    <FormField
      control={control}
      name={name}
      render={({ field }) => (
        <FormItem>
          <FormLabel>{label}</FormLabel>
          <FormControl>
            <div className="relative">
              <Input
                id={id}
                type={isPassword && !showPassword ? "password" : type}
                placeholder={placeholder}
                className={isPassword ? "pr-10" : ""}
                {...field}
                aria-describedby={`${name}-error`}
              />
              {isPassword && (
                <Button
                  type="button"
                  variant="ghost"
                  size="sm"
                  className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                  onClick={togglePassword}
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  {showPassword ? <EyeOff /> : <Eye />}
                </Button>
              )}
            </div>
          </FormControl>
          <FormMessage id={`${name}-error`} />
        </FormItem>
      )}
    />
  );
}
