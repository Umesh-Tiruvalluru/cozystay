"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/lib/auth-context";
import { Sparkles, MapPin, ShieldCheck, CalendarCheck2, ArrowRight, Star } from "lucide-react";

export default function Home() {
  const { isAuthenticated } = useAuth();

  return (
    <main className="min-h-screen bg-gradient-to-b from-background via-background to-muted/20">
      <section className="relative overflow-hidden">
        <div className="pointer-events-none absolute inset-0 -z-10">
          <div className="absolute -top-24 -left-24 h-72 w-72 rounded-full bg-primary/10 blur-3xl" />
          <div className="absolute -bottom-24 -right-24 h-72 w-72 rounded-full bg-primary/10 blur-3xl" />
        </div>
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
          <div className="text-center space-y-6">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full border bg-background/60 backdrop-blur text-xs text-muted-foreground">
              <Sparkles className="h-3.5 w-3.5 text-primary" /> Handpicked places to stay
            </div>
            <h1 className="font-heading text-4xl sm:text-5xl md:text-6xl font-extrabold tracking-tight">
              Find extraordinary stays for your next getaway
            </h1>
            <p className="text-base md:text-lg text-muted-foreground max-w-2xl mx-auto">
              Book unique homes and experiences all over the world. Simple, secure, and effortless.
            </p>
            <div className="flex gap-4 justify-center">
              <Link href="/properties">
                <Button size="lg" className="text-base md:text-lg gap-2">
                  Explore Properties <ArrowRight className="h-4 w-4" />
                </Button>
              </Link>
              {!isAuthenticated && (
                <Link href="/register">
                  <Button size="lg" variant="outline" className="text-base md:text-lg">
                    Get Started
                  </Button>
                </Link>
              )}
            </div>

            {/* Quick search mock */}
            <div className="mt-8 grid grid-cols-1 sm:grid-cols-3 gap-3 max-w-3xl mx-auto">
              <div className="flex items-center gap-2 rounded-lg border bg-background px-4 py-3 text-left text-sm">
                <MapPin className="h-4 w-4 text-primary" /> Anywhere
              </div>
              <div className="flex items-center gap-2 rounded-lg border bg-background px-4 py-3 text-left text-sm">
                <CalendarCheck2 className="h-4 w-4 text-primary" /> Anytime
              </div>
              <div className="flex items-center gap-2 rounded-lg border bg-background px-4 py-3 text-left text-sm">
                <Star className="h-4 w-4 text-primary" /> Amazing stays
              </div>
            </div>
          </div>
        </div>
      </section>

      <section className="py-12">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="rounded-xl border bg-card p-6 shadow-sm">
              <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <MapPin className="h-5 w-5 text-primary" />
              </div>
              <h3 className="text-lg font-semibold mb-2">Find your perfect stay</h3>
              <p className="text-muted-foreground">
                Browse thousands of unique properties across top destinations.
              </p>
            </div>
            <div className="rounded-xl border bg-card p-6 shadow-sm">
              <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <CalendarCheck2 className="h-5 w-5 text-primary" />
              </div>
              <h3 className="text-lg font-semibold mb-2">Easy, secure booking</h3>
              <p className="text-muted-foreground">
                Check availability and book in minutes with transparent pricing.
              </p>
            </div>
            <div className="rounded-xl border bg-card p-6 shadow-sm">
              <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center mb-4">
                <ShieldCheck className="h-5 w-5 text-primary" />
              </div>
              <h3 className="text-lg font-semibold mb-2">Trusted hosts</h3>
              <p className="text-muted-foreground">
                Stay with verified hosts and enjoy memorable experiences.
              </p>
            </div>
          </div>
        </div>
      </section>

      <section className="py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="rounded-2xl border bg-gradient-to-br from-primary/5 to-primary/10 p-8 md:p-10 flex flex-col md:flex-row items-center justify-between gap-6">
            <div className="space-y-2">
              <h3 className="text-2xl font-bold">Ready for your next trip?</h3>
              <p className="text-muted-foreground">
                Explore stunning homes and book your perfect stay today.
              </p>
            </div>
            <Link href="/properties">
              <Button size="lg" className="gap-2">
                Start exploring <ArrowRight className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>
      </section>
    </main>
  );
}
