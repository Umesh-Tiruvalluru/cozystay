"use client";

import type React from "react";
import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import { format } from "date-fns";
import { CalendarIcon, MapPin, Users, Sparkles, ArrowLeft, Share2 } from "lucide-react";
import type { DateRange } from "react-day-picker";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { Calendar } from "@/components/ui/calendar";
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
} from "@/components/ui/carousel";

import {
  propertiesApi,
  bookingsApi,
  type PropertyDetail,
} from "@/lib/api-client";
import { useAuth } from "@/lib/auth-context";

export default function PropertyDetailsPage() {
  const params = useParams();
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const propertyId = params.id as string;

  const [property, setProperty] = useState<PropertyDetail | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isBooking, setIsBooking] = useState(false);
  const [error, setError] = useState("");

  const [date, setDate] = useState<DateRange | undefined>(undefined);

  useEffect(() => {
    loadProperty();
  }, [propertyId]);

  const loadProperty = async () => {
    try {
      const response = await propertiesApi.getById(propertyId);
      setProperty(response.data);
    } catch (error) {
      console.error("Failed to load property:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleBooking = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");

    if (!isAuthenticated) {
      router.push("/login");
      return;
    }

    if (!date?.from || !date?.to) {
      setError("Please select both check-in and check-out dates");
      return;
    }

    if (date.from >= date.to) {
      setError("End date must be after start date");
      return;
    }

    setIsBooking(true);
    try {
      const start_date = date.from.toISOString().split("T")[0];
      const end_date = date.to.toISOString().split("T")[0];

      const nights = Math.ceil(
        (new Date(end_date).getTime() - new Date(start_date).getTime()) /
          (1000 * 60 * 60 * 24)
      );
      const totalPrice = nights * property!.price_per_night;

      await bookingsApi.create({
        property_id: propertyId,
        start_date,
        end_date,
        total_price: totalPrice,
      });

      router.push("/bookings");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Booking failed");
    } finally {
      setIsBooking(false);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-background to-muted/20">
        <div className="text-center space-y-4">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto" />
          <p className="text-muted-foreground">Loading property...</p>
        </div>
      </div>
    );
  }

  if (!property) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <Card className="p-8">
          <p className="text-muted-foreground">Property not found</p>
        </Card>
      </div>
    );
  }

  const nights =
    date?.from && date?.to
      ? Math.ceil(
          (new Date(date.to).getTime() - new Date(date.from).getTime()) /
            (1000 * 60 * 60 * 24)
        )
      : 0;

  const totalPrice = nights * property.price_per_night;

  // Gallery: single image, carousel for many, fallback box otherwise
  const renderImageGrid = () => {
    if (!property.images || property.images.length === 0) {
      return (
        <div className="w-full h-96 bg-gradient-to-br from-muted to-muted/50 rounded-xl flex items-center justify-center border border-border/50">
          <span className="text-muted-foreground">No images available</span>
        </div>
      );
    }

    if (property.images.length === 1) {
      return (
        <div className="w-full h-96 rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow">
          <img
            src={property.images[0].image_url}
            alt={property.images[0].caption}
            className="w-full h-full object-cover hover:scale-105 transition-transform duration-300"
          />
        </div>
      );
    }

    return (
      <Carousel className="w-full">
        <CarouselContent>
          {property.images.map((image) => (
            <CarouselItem key={image.id} className="basis-full">
              <div className="w-full h-96 rounded-xl overflow-hidden shadow-lg">
                <img
                  src={image.image_url}
                  alt={image.caption}
                  className="w-full h-full object-cover"
                />
              </div>
            </CarouselItem>
          ))}
        </CarouselContent>
        <CarouselPrevious className="left-2" />
        <CarouselNext className="right-2" />
      </Carousel>
    );
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/10">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-6 flex items-center justify-between">
          <button
            onClick={() => router.back()}
            className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
          >
            <ArrowLeft className="h-4 w-4" /> Back
          </button>
          <button
            onClick={() => navigator.share?.({ title: property?.title, url: window.location.href })}
            className="inline-flex items-center gap-2 text-sm text-muted-foreground hover:text-foreground"
          >
            <Share2 className="h-4 w-4" /> Share
          </button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            {renderImageGrid()}

            <Card className="shadow-md border-border/50">
              <CardHeader className="space-y-3">
                <CardTitle className="font-heading text-3xl font-bold bg-gradient-to-r from-foreground to-foreground/70 bg-clip-text">
                  {property.title}
                </CardTitle>
                <CardDescription className="flex items-center text-base">
                  <MapPin className="mr-2 h-4 w-4" />
                  {property.location}
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-8">
                <div>
                  <h3 className="font-semibold text-lg mb-3 flex items-center">
                    <Sparkles className="mr-2 h-5 w-5 text-primary" />
                    Description
                  </h3>
                  <p className="text-muted-foreground leading-relaxed">
                    {property.description}
                  </p>
                </div>

                {property.amenities && property.amenities.length > 0 && (
                  <div className="bg-muted/30 rounded-lg p-6">
                    <h3 className="font-semibold text-lg mb-4">
                      What this place offers
                    </h3>
                    <div className="grid grid-cols-2 gap-3">
                      {property.amenities.map((amenity) => (
                        <div
                          key={(amenity as any).name}
                          className="flex items-center space-x-3 bg-background/80 rounded-md px-3 py-2"
                        >
                          <span className="text-primary text-lg">✓</span>
                          <span className="text-sm">{(amenity as any).name}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                <div className="grid grid-cols-2 gap-6">
                  <div className="bg-gradient-to-br from-primary/5 to-primary/10 rounded-lg p-4 border border-primary/20">
                    <div className="flex items-center mb-2">
                      <Users className="h-5 w-5 text-primary mr-2" />
                      <p className="text-sm text-muted-foreground">
                        Max Guests
                      </p>
                    </div>
                    <p className="text-2xl font-bold">{property.max_guests}</p>
                  </div>
                  <div className="bg-gradient-to-br from-primary/5 to-primary/10 rounded-lg p-4 border border-primary/20">
                    <p className="text-sm text-muted-foreground mb-2">
                      Price per Night
                    </p>
                    <p className="text-2xl font-bold text-primary">
                      ${property.price_per_night.toLocaleString?.() ?? property.price_per_night}
                    </p>
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          <div className="lg:sticky lg:top-8 h-fit">
            <Card className="shadow-lg border-border/50">
              <CardHeader className="bg-gradient-to-br from-primary/5 to-primary/10 rounded-t-lg">
                <CardTitle className="text-xl">Book This Property</CardTitle>
              </CardHeader>
              <CardContent className="pt-6">
                <form onSubmit={handleBooking} className="space-y-5">
                  {error && (
                    <div className="bg-destructive/10 text-destructive p-4 rounded-lg text-sm border border-destructive/20 animate-in fade-in slide-in-from-top-2">
                      {error}
                    </div>
                  )}

                  <div className="space-y-2">
                    <label className="text-sm font-medium flex items-center">
                      <CalendarIcon className="mr-2 h-4 w-4 text-primary" />
                      Select Dates
                    </label>
                    <Popover>
                      <PopoverTrigger asChild>
                        <Button
                          variant="outline"
                          className={`w-full justify-start text-left font-normal h-12 ${
                            !date?.from ? "text-muted-foreground" : ""
                          }`}
                        >
                          <CalendarIcon className="mr-2 h-4 w-4" />
                          {date?.from ? (
                            date.to ? (
                              <>
                                {format(date.from, "LLL dd, y")} -{" "}
                                {format(date.to, "LLL dd, y")}
                              </>
                            ) : (
                              format(date.from, "LLL dd, y")
                            )
                          ) : (
                            <span>Pick a date range</span>
                          )}
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="w-auto p-0" align="start">
                        <Calendar
                          mode="range"
                          selected={date}
                          onSelect={setDate}
                          numberOfMonths={2}
                          disabled={(day) => day < new Date()}
                        />
                      </PopoverContent>
                    </Popover>
                  </div>

                  {nights > 0 && (
                    <div className="bg-gradient-to-br from-muted/50 to-muted rounded-lg p-5 space-y-3 border border-border/50 animate-in fade-in slide-in-from-bottom-2">
                      <div className="flex justify-between text-sm">
                        <span className="text-muted-foreground">
                          ${property.price_per_night} × {nights} nights
                        </span>
                        <span className="font-semibold">${totalPrice}</span>
                      </div>
                      <div className="border-t border-border/50 pt-3 flex justify-between items-center">
                        <span className="font-bold text-lg">Total</span>
                        <span className="text-2xl font-bold text-primary">
                          ${totalPrice}
                        </span>
                      </div>
                    </div>
                  )}

                  <Button
                    type="submit"
                    className="w-full h-12 text-base font-semibold shadow-md hover:shadow-lg transition-all"
                    disabled={isBooking || !isAuthenticated}
                  >
                    {isBooking
                      ? "Booking..."
                      : isAuthenticated
                      ? "Book Now"
                      : "Login to Book"}
                  </Button>

                  {!isAuthenticated && (
                    <p className="text-xs text-center text-muted-foreground">
                      You need to login to make a booking
                    </p>
                  )}
                </form>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
