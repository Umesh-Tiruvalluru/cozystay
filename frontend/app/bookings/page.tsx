"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { bookingsApi, type Booking } from "@/lib/api-client";
import { useAuth } from "@/lib/auth-context";
import { Calendar, MapPin, DollarSign, XCircle, Loader2 } from "lucide-react";
import Link from "next/link";

export default function BookingsPage() {
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuth();
  const [bookings, setBookings] = useState<Booking[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [cancelingId, setCancelingId] = useState<string | null>(null);
  const [showCancelDialog, setShowCancelDialog] = useState(false);
  const [selectedBookingId, setSelectedBookingId] = useState<string | null>(null);
  const [cancelError, setCancelError] = useState<string | null>(null);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push("/login");
      return;
    }

    if (isAuthenticated) {
      loadBookings();
    }
  }, [isAuthenticated, authLoading]);

  const loadBookings = async () => {
    try {
      const response = await bookingsApi.getAll();
      setBookings(response.data || []);
    } catch (error) {
      console.error("Failed to load bookings:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const openCancelDialog = (bookingId: string) => {
    setSelectedBookingId(bookingId);
    setShowCancelDialog(true);
    setCancelError(null);
  };

  const closeCancelDialog = () => {
    setShowCancelDialog(false);
    setSelectedBookingId(null);
    setCancelError(null);
  };

  const confirmCancel = async () => {
    if (!selectedBookingId) return;

    setCancelingId(selectedBookingId);
    setCancelError(null);
    
    try {
      await bookingsApi.cancel(selectedBookingId);
      await loadBookings();
      closeCancelDialog();
    } catch (error) {
      console.error("Failed to cancel booking:", error);
      setCancelError("Failed to cancel booking. Please try again.");
    } finally {
      setCancelingId(null);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString("en-US", {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  };

  const getStatusColor = (status: string) => {
    if (status === "cancelled") {
      return "bg-red-100 text-red-800 border-red-200";
    }
    return "bg-green-100 text-green-800 border-green-200";
  };

  if (authLoading || isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-background via-background to-muted/20">
        <div className="text-center space-y-4">
          <Loader2 className="h-12 w-12 animate-spin text-primary mx-auto" />
          <p className="text-muted-foreground">Loading your bookings...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/20">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="mb-8">
          <h1 className="font-heading text-3xl md:text-4xl font-bold mb-2">
            My Bookings
          </h1>
          <p className="text-muted-foreground">
            Manage and view all your upcoming and past reservations
          </p>
        </div>

        {bookings.length === 0 ? (
          <Card className="shadow-md">
            <CardContent className="pt-12 pb-12 text-center">
              <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted mb-4">
                <Calendar className="h-8 w-8 text-muted-foreground" />
              </div>
              <h3 className="text-xl font-semibold mb-2">No bookings yet</h3>
              <p className="text-muted-foreground mb-6 max-w-md mx-auto">
                Start exploring our collection of unique properties and book your perfect stay.
              </p>
              <Link href="/properties">
                <Button size="lg" className="gap-2">
                  Browse Properties
                </Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          <div className="space-y-4">
            {bookings.map((booking) => (
              <Card key={booking.id} className="hover:shadow-lg transition-shadow">
                <CardHeader>
                  <div className="flex flex-col sm:flex-row justify-between items-start gap-4">
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <CardTitle className="text-xl">
                          Booking #{booking.id?.slice(0, 8) || booking.id || "N/A"}
                        </CardTitle>
                      </div>
                      <CardDescription className="flex items-center gap-2 text-sm">
                        <MapPin className="h-4 w-4" />
                        Property ID: {booking.property_id?.slice(0, 8) || "N/A"}
                      </CardDescription>
                    </div>
                    <span
                      className={`px-3 py-1 rounded-full text-xs font-semibold border ${getStatusColor(
                        booking.status
                      )}`}
                    >
                      {booking.status.charAt(0).toUpperCase() + booking.status.slice(1)}
                    </span>
                  </div>
                </CardHeader>
                <CardContent>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
                    <div className="flex items-start gap-3">
                      <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center">
                        <Calendar className="h-5 w-5 text-primary" />
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground uppercase tracking-wide">
                          Check-in
                        </p>
                        <p className="font-semibold">{formatDate(booking.start_date)}</p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3">
                      <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center">
                        <Calendar className="h-5 w-5 text-primary" />
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground uppercase tracking-wide">
                          Check-out
                        </p>
                        <p className="font-semibold">{formatDate(booking.end_date)}</p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3">
                      <div className="h-10 w-10 rounded-lg bg-primary/10 flex items-center justify-center">
                        <DollarSign className="h-5 w-5 text-primary" />
                      </div>
                      <div>
                        <p className="text-xs text-muted-foreground uppercase tracking-wide">
                          Total Price
                        </p>
                        <p className="text-2xl font-bold text-primary">
                          ${booking.total_price}
                        </p>
                      </div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between pt-4 border-t">
                    <div className="text-xs text-muted-foreground">
                      Booking ID: {booking.id?.slice(0, 8) || "N/A"}...
                    </div>
                    {booking.status === "booked" && (
                      <Button
                        variant="destructive"
                        onClick={() => openCancelDialog(booking.id)}
                        disabled={cancelingId === booking.id}
                        className="gap-2"
                      >
                        <XCircle className="h-4 w-4" />
                        Cancel Booking
                      </Button>
                    )}
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>

      <Dialog open={showCancelDialog} onOpenChange={setShowCancelDialog} >
        <DialogContent className="bg-white">
          <DialogHeader>
            <DialogTitle>Cancel Booking?</DialogTitle>
            <DialogDescription>
              Are you sure you want to cancel this booking? This action cannot be undone.
            </DialogDescription>
          </DialogHeader>

          {cancelError && (
            <div className="bg-destructive/10 text-destructive p-3 rounded-lg text-sm border border-destructive/20">
              {cancelError}
            </div>
          )}

          <DialogFooter>
            <Button
              variant="outline"
              onClick={closeCancelDialog}
              disabled={!!cancelingId}
            >
              Keep Booking
            </Button>
            <Button
              variant="destructive"
              onClick={confirmCancel}
              disabled={!!cancelingId}
              className="gap-2"
            >
              {cancelingId ? (
                <>
                  <Loader2 className="h-4 w-4 animate-spin" />
                  Canceling...
                </>
              ) : (
                <>
                  <XCircle className="h-4 w-4" />
                  Yes, Cancel Booking
                </>
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
