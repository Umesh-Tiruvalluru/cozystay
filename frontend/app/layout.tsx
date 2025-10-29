import type React from "react";
import type { Metadata } from "next";
import { Figtree, Montserrat } from "next/font/google";
import { Analytics } from "@vercel/analytics/next";
import "./globals.css";
import { AuthProvider } from "@/lib/auth-context";
import { Navbar } from "@/components/navbar";

const figtree = Figtree({ subsets: ["latin"], display: "swap" });
const montserrat = Montserrat({ subsets: ["latin"], display: "swap", variable: "--font-heading" });

export const metadata: Metadata = {
  title: "BookBnb — Find extraordinary stays",
  description: "Discover and book unique properties around the world with BookBnb.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        <link
          href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined"
          rel="stylesheet"
        />
      </head>
      <body className={`${figtree.className} ${montserrat.variable} antialiased`}>
        <AuthProvider>
          <Navbar />
          {children}
        </AuthProvider>
        <Analytics />
      </body>
    </html>
  );
}
