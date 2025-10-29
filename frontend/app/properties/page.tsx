"use client";

import type React from "react";

import Link from "next/link";
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { propertiesApi, type Property } from "@/lib/api-client";
import Image from "next/image";
import { Sparkles, MapPin } from "lucide-react";

export default function PropertiesPage() {
  const [properties, setProperties] = useState<Property[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [filters, setFilters] = useState({
    location: "",
    start_date: "",
    end_date: "",
    max_guests: "",
  });

  useEffect(() => {
    loadProperties();
  }, []);

  const loadProperties = async () => {
    setIsLoading(true);
    try {
      const response = await propertiesApi.getAll();
      setProperties(response.data || []);
    } catch (error) {
      console.error("Failed to load properties:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFilters((prev) => ({ ...prev, [name]: value }));
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    loadProperties();
  };

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-10 rounded-2xl p-8 bg-gradient-to-br from-primary/5 to-primary/10 border border-primary/10">
          <div className="flex items-center gap-3 mb-2">
            <div className="p-2 rounded-lg bg-primary/10">
              <Sparkles className="h-6 w-6 text-primary" />
            </div>
            <h1 className="font-heading text-3xl md:text-4xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text text-transparent">
              Discover Your Perfect Stay
            </h1>
          </div>
          <p className="text-base md:text-lg text-muted-foreground">
            Explore our curated collection of unique properties worldwide
          </p>
        </div>

        <Card className="mb-8">
          <CardHeader>
            <CardTitle>Search Properties</CardTitle>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSearch} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                <Input
                  placeholder="Location"
                  name="location"
                  value={filters.location}
                  onChange={handleFilterChange}
                />
                <Input
                  type="date"
                  name="start_date"
                  value={filters.start_date}
                  onChange={handleFilterChange}
                />
                <Input
                  type="date"
                  name="end_date"
                  value={filters.end_date}
                  onChange={handleFilterChange}
                />
                <Input
                  type="number"
                  placeholder="Max Guests"
                  name="max_guests"
                  value={filters.max_guests}
                  onChange={handleFilterChange}
                  min="1"
                />
              </div>
              <Button type="submit" className="w-full">
                Search
              </Button>
            </form>
          </CardContent>
        </Card>

        {isLoading ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="animate-pulse">
                <div className="w-full h-48 bg-muted rounded-t-lg" />
                <div className="border border-t-0 rounded-b-lg p-4 space-y-3">
                  <div className="h-5 bg-muted rounded w-2/3" />
                  <div className="h-4 bg-muted rounded w-1/2" />
                  <div className="h-4 bg-muted rounded w-1/3" />
                </div>
              </div>
            ))}
          </div>
        ) : properties.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-muted-foreground">
              No properties found. Try adjusting your search.
            </p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {properties.map((property) => (
              <Link key={property.id} href={`/properties/${property.id}`}>
                <Card className="group h-full hover:shadow-lg py-0 transition-shadow cursor-pointer overflow-hidden">
                  <div className="relative w-full h-48 bg-muted">
                    {property.thumbnail_url?.Valid ? (
                      <Image
                        src={property.thumbnail_url.String}
                        alt="property-image"
                        fill
                        className="object-cover group-hover:scale-[1.03] transition-transform duration-300"
                      />
                    ) : (
                      <div className="w-full h-full flex items-center justify-center text-muted-foreground">
                        Property Image
                      </div>
                    )}
                    <div className="absolute bottom-2 left-2 bg-white/95 backdrop-blur-sm px-3 py-1.5 rounded-lg text-sm font-bold text-primary shadow-lg border-2 border-white/50">
                      ${property.price_per_night}/night
                    </div>
                  </div>
                  <CardContent className="py-4 space-y-2">
                    <h3 className="font-semibold text-lg line-clamp-1">
                      {property.title}
                    </h3>
                    <p className="text-sm text-muted-foreground flex items-center gap-1 line-clamp-1">
                      <MapPin className="h-3.5 w-3.5" /> {property.location}
                    </p>
                    <div className="flex items-center justify-between pt-2 border-t">
                      <span className="text-sm font-medium text-foreground">{property.max_guests} guests</span>
                      <span className="text-sm font-bold text-primary">${property.price_per_night}/night</span>
                    </div>
                  </CardContent>
                </Card>
              </Link>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
