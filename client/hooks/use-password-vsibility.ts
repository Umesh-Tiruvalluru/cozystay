// hooks/usePasswordVisibility.ts
import { useState } from "react";

export function usePasswordVisibility() {
  const [showPassword, setShowPassword] = useState(false);
  const togglePassword = () => setShowPassword((prev) => !prev);
  return { showPassword, togglePassword };
}
