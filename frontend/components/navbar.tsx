"use client";

import Link from "next/link";
import { useAuth } from "@/lib/auth-context";
import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Home, Calendar, Shield, User, LogOut, Sparkles } from "lucide-react";
import Logo from "./logo";

export function Navbar() {
  const { user, isAuthenticated, logout } = useAuth();

  // Get user initials for avatar fallback
  const getUserInitials = () => {
    if (!user) return "U";
    const firstInitial = user.first_name?.charAt(0) || "";
    const lastInitial = user.last_name?.charAt(0) || "";
    return `${firstInitial}${lastInitial}`.toUpperCase() || "U";
  };

  return (
    <nav
      className="sticky h-16 inset-x-0 top-0 w-full border-b
         border-zinc-200 backdrop-blur-lg transition-all"
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <Link href="/" className="flex items-center space-x-2 group">
            <Logo />
          </Link>

          {/* Navigation Links */}
          <div className="flex items-center gap-2">
            {isAuthenticated && user ? (
              <>
                <Link href="/properties">
                  <Button
                    variant="link"
                    className="gap-2 text-zinc-800 hover:text-primary transition-colors"
                  >
                    <span className="hidden sm:inline">Browse</span>
                  </Button>
                </Link>
                <Link href="/bookings">
                  <Button
                    variant="link"
                    className="gap-2 text-zinc-800 hover:text-primary transition-colors"
                  >
                    <span className="hidden sm:inline">My Bookings</span>
                  </Button>
                </Link>
                {user.role === "admin" && (
                  <Link href="/admin">
                    <Button
                      variant="link"
                      className="gap-2 text-zinc-800 hover:text-primary transition-colors"
                    >
                      <span className="hidden sm:inline">Admin</span>
                    </Button>
                  </Link>
                )}

                {/* User Dropdown Menu */}
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button
                      variant="ghost"
                      className="relative h-10 w-10 rounded-full ring-1 ring-zinc-800/70 hover:ring-zinc-800 transition-all"
                    >
                      <Avatar className="h-9 w-9">
                        {/* <AvatarImage
                          src={user.avatar_ur }
                          alt={`${user.first_name} ${user.last_name}`}
                        /> */}
                        <AvatarFallback className="text-zinc-200 bg-zinc-800 font-semibold">
                          {getUserInitials()}
                        </AvatarFallback>
                      </Avatar>
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56">
                    <DropdownMenuLabel>
                      <div className="flex flex-col space-y-1">
                        <p className="text-sm font-medium leading-none">
                          {user.first_name} {user.last_name}
                        </p>
                        <p className="text-xs leading-none text-muted-foreground">
                          {user.email}
                        </p>
                      </div>
                    </DropdownMenuLabel>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem asChild>
                      <Link href="/profile" className="cursor-pointer">
                        <User className="mr-2 h-4 w-4" />
                        <span>Profile</span>
                      </Link>
                    </DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem
                      onClick={logout}
                      className="text-destructive focus:text-destructive cursor-pointer"
                    >
                      <LogOut className="mr-2 h-4 w-4" />
                      <span>Logout</span>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </>
            ) : (
              <>
                <Link href="/login">
                  <Button
                    variant="ghost"
                    className="hover:bg-primary/10 hover:text-primary transition-colors"
                  >
                    Login
                  </Button>
                </Link>
                <Link href="/register">
                  <Button className="shadow-md hover:shadow-lg transition-shadow bg-gradient-to-r from-primary to-primary/80">
                    Sign Up
                  </Button>
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
