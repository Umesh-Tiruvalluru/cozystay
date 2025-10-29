// components/auth/SocialAuthButtons.tsx
import { Button } from "@/components/ui/button";
import { FaGoogle } from "react-icons/fa";

export default function SocialAuthButtons() {
  return (
    <div>
      <div className="relative mt-4">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-slate-300"></div>
        </div>
        <div className="relative flex justify-center text-sm">
          <span className="px-2 bg-white text-slate-500">Or continue with</span>
        </div>
      </div>
      <div className="mt-6 gap-3">
        <Button
          variant="outline"
          className="w-full border-slate-300 text-slate-500 hover:bg-slate-50"
          onClick={() => {
            // Implement Google OAuth redirect
            console.log("Google sign-in clicked");
          }}
        >
          <FaGoogle className="mr-2" />
          Sign In with Google
        </Button>
      </div>
    </div>
  );
}
