"use client";

import type React from "react";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Textarea } from "@/components/ui/textarea";
import { amenitiesApi, propertiesApi, type Amenity } from "@/lib/api-client";
import { useAuth } from "@/lib/auth-context";

export default function AdminPage() {
  const router = useRouter();
  const { user, isAuthenticated, isLoading: authLoading } = useAuth();
  const [properties, setProperties] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formData, setFormData] = useState({
    title: "",
    description: "",
    location: "",
    price_per_night: "",
    max_guests: "",
    image_url: "",
  });

  const [amenities, setAmenities] = useState<Amenity[]>([]);
  const [isAmenitySubmitting, setIsAmenitySubmitting] = useState(false);
  const [newAmenityName, setNewAmenityName] = useState("");

  const [editingId, setEditingId] = useState<string | null>(null);
  const [editData, setEditData] = useState<Record<string, any>>({});
  const [imageInputs, setImageInputs] = useState<Record<string, { image_url: string; caption: string; display_order: string }>>({});
  const [selectedAmenityIds, setSelectedAmenityIds] = useState<Record<string, Set<string>>>({});

  useEffect(() => {
    if (!authLoading && (!isAuthenticated || user?.role !== "admin")) {
      router.push("/");
      return;
    }

    if (isAuthenticated && user?.role === "admin") {
      Promise.all([loadProperties(), loadAmenities()]).finally(() => setIsLoading(false));
    }
  }, [isAuthenticated, authLoading, user]);

  const loadProperties = async () => {
    try {
      const response = await propertiesApi.getAll();
      setProperties(response.data || []);
    } catch (error) {
      console.error("Failed to load properties:", error);
    }
  };

  const loadAmenities = async () => {
    try {
      const response = await amenitiesApi.getAll();
      setAmenities(response.data || []);
    } catch (error) {
      console.error("Failed to load amenities:", error);
    }
  };

  const handleFormChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);

    try {
      await propertiesApi.create({
        title: formData.title,
        description: formData.description,
        location: formData.location,
        price_per_night: Number.parseFloat(formData.price_per_night),
        max_guests: Number.parseInt(formData.max_guests),
        image_url: formData.image_url || undefined,
      });

      setFormData({
        title: "",
        description: "",
        location: "",
        price_per_night: "",
        max_guests: "",
        image_url: "",
      });
      setShowForm(false);
      loadProperties();
    } catch (error) {
      console.error("Failed to create property:", error);
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (propertyId: string) => {
    if (confirm("Are you sure you want to delete this property?")) {
      try {
        await propertiesApi.delete(propertyId);
        setProperties((prev) => prev.filter((p) => p.id !== propertyId));
      } catch (error) {
        console.error("Failed to delete property:", error);
      }
    }
  };

  const startEdit = (p: any) => {
    setEditingId(p.id);
    setEditData({
      id: p.id,
      title: p.title ?? "",
      description: p.description ?? "",
      location: p.location ?? "",
      price_per_night: p.price_per_night?.toString?.() ?? "",
      max_guests: p.max_guests?.toString?.() ?? "",
    });
  };

  const handleEditChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setEditData((prev) => ({ ...prev, [name]: value }));
  };

  const saveEdit = async () => {
    if (!editingId) return;
    try {
      await propertiesApi.update({
        id: editData.id,
        title: editData.title,
        description: editData.description,
        location: editData.location,
        price_per_night: Number.parseFloat(editData.price_per_night),
        max_guests: Number.parseInt(editData.max_guests),
      });
      setEditingId(null);
      await loadProperties();
    } catch (error) {
      console.error("Failed to update property:", error);
    }
  };

  const cancelEdit = () => {
    setEditingId(null);
    setEditData({});
  };

  const handleImageInputChange = (
    propertyId: string,
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const { name, value } = e.target;
    setImageInputs((prev) => ({
      ...prev,
      [propertyId]: {
        image_url: prev[propertyId]?.image_url ?? "",
        caption: prev[propertyId]?.caption ?? "",
        display_order: prev[propertyId]?.display_order ?? "",
        [name]: value,
      },
    }));
  };

  const addImage = async (propertyId: string) => {
    const values = imageInputs[propertyId] || {
      image_url: "",
      caption: "",
      display_order: "",
    };
    if (!values.image_url) return;
    try {
      await propertiesApi.addImage(propertyId, [
        {
          image_url: values.image_url,
          caption: values.caption || undefined,
          display_order: Number.parseInt(values.display_order || "0"),
        },
      ]);
      setImageInputs((prev) => ({ ...prev, [propertyId]: { image_url: "", caption: "", display_order: "" } }));
    } catch (error) {
      console.error("Failed to add image:", error);
    }
  };

  const toggleAmenitySelected = (propertyId: string, amenityId: string) => {
    setSelectedAmenityIds((prev) => {
      const current = new Set(prev[propertyId] ?? []);
      if (current.has(amenityId)) current.delete(amenityId);
      else current.add(amenityId);
      return { ...prev, [propertyId]: current };
    });
  };

  const addAmenitiesToProperty = async (propertyId: string) => {
    const selected = Array.from(selectedAmenityIds[propertyId] ?? []);
    if (selected.length === 0) return;
    try {
      await amenitiesApi.addToProperty(propertyId, selected);
    } catch (error) {
      console.error("Failed to add amenities:", error);
    }
  };

  const createAmenity = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newAmenityName.trim()) return;
    setIsAmenitySubmitting(true);
    try {
      await amenitiesApi.create({ name: newAmenityName.trim() });
      setNewAmenityName("");
      await loadAmenities();
    } catch (error) {
      console.error("Failed to create amenity:", error);
    } finally {
      setIsAmenitySubmitting(false);
    }
  };

  if (authLoading || isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <p className="text-muted-foreground">Loading...</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold">Admin Dashboard</h1>
          <Button onClick={() => setShowForm(!showForm)}>
            {showForm ? "Cancel" : "Add Property"}
          </Button>
        </div>

        {showForm && (
          <Card className="mb-8 shadow-lg">
            <CardHeader>
              <CardTitle>Create New Property</CardTitle>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit} className="space-y-4">
                <Input
                  placeholder="Title"
                  name="title"
                  value={formData.title}
                  onChange={handleFormChange}
                  required
                />
                <Textarea
                  placeholder="Description"
                  name="description"
                  value={formData.description}
                  onChange={handleFormChange}
                  required
                />
                <Input
                  placeholder="Location"
                  name="location"
                  value={formData.location}
                  onChange={handleFormChange}
                  required
                />
                <div className="grid grid-cols-2 gap-4">
                  <Input
                    type="number"
                    placeholder="Price per Night"
                    name="price_per_night"
                    value={formData.price_per_night}
                    onChange={handleFormChange}
                    step="0.01"
                    required
                  />
                  <Input
                    type="number"
                    placeholder="Max Guests"
                    name="max_guests"
                    value={formData.max_guests}
                    onChange={handleFormChange}
                    min="1"
                    required
                  />
                </div>
                <Input
                  placeholder="Thumbnail Image URL (optional)"
                  name="image_url"
                  value={formData.image_url}
                  onChange={handleFormChange}
                />
                <Button
                  type="submit"
                  className="w-full"
                  disabled={isSubmitting}
                >
                  {isSubmitting ? "Creating..." : "Create Property"}
                </Button>
              </form>
            </CardContent>
          </Card>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-4">
            {properties.map((property) => (
              <Card key={property.id}>
                <CardContent className="pt-6 space-y-4">
                  <div className="flex justify-between items-start">
                    <div className="flex-1">
                      <h3 className="font-semibold text-lg">{property.title}</h3>
                      <p className="text-sm text-muted-foreground">{property.location}</p>
                      <div className="flex gap-4 mt-2 text-sm">
                        {property.price_per_night != null && (
                          <span>${property.price_per_night}/night</span>
                        )}
                        {property.max_guests != null && (
                          <span>{property.max_guests} guests</span>
                        )}
                      </div>
                    </div>
                    <div className="flex gap-2">
                      {editingId === property.id ? (
                        <>
                          <Button variant="secondary" onClick={saveEdit}>Save</Button>
                          <Button variant="ghost" onClick={cancelEdit}>Cancel</Button>
                        </>
                      ) : (
                        <>
                          <Button variant="secondary" onClick={() => startEdit(property)}>Edit</Button>
                          <Button variant="destructive" onClick={() => handleDelete(property.id)}>Delete</Button>
                        </>
                      )}
                    </div>
                  </div>

                  {editingId === property.id && (
                    <div className="space-y-3">
                      <Input name="title" value={editData.title ?? ""} onChange={handleEditChange} placeholder="Title" />
                      <Textarea name="description" value={editData.description ?? ""} onChange={handleEditChange} placeholder="Description" />
                      <Input name="location" value={editData.location ?? ""} onChange={handleEditChange} placeholder="Location" />
                      <div className="grid grid-cols-2 gap-3">
                        <Input name="price_per_night" type="number" step="0.01" value={editData.price_per_night ?? ""} onChange={handleEditChange} placeholder="Price per Night" />
                        <Input name="max_guests" type="number" value={editData.max_guests ?? ""} onChange={handleEditChange} placeholder="Max Guests" />
                      </div>
                    </div>
                  )}

                  <div className="border-t pt-4 space-y-3">
                    <h4 className="font-medium">Add Image</h4>
                    <div className="grid grid-cols-3 gap-3">
                      <Input
                        name="image_url"
                        placeholder="Image URL"
                        value={imageInputs[property.id]?.image_url ?? ""}
                        onChange={(e) => handleImageInputChange(property.id, e)}
                      />
                      <Input
                        name="caption"
                        placeholder="Caption (optional)"
                        value={imageInputs[property.id]?.caption ?? ""}
                        onChange={(e) => handleImageInputChange(property.id, e)}
                      />
                      <Input
                        name="display_order"
                        type="number"
                        placeholder="Display Order"
                        value={imageInputs[property.id]?.display_order ?? ""}
                        onChange={(e) => handleImageInputChange(property.id, e)}
                      />
                    </div>
                    <div>
                      <Button size="sm" onClick={() => addImage(property.id)}>Add Image</Button>
                    </div>
                  </div>

                  <div className="border-t pt-4 space-y-3">
                    <h4 className="font-medium">Assign Amenities</h4>
                    <div className="grid grid-cols-2 md:grid-cols-3 gap-2 max-h-40 overflow-auto pr-2">
                      {amenities.map((a) => {
                        const checked = selectedAmenityIds[property.id]?.has(a.amenity_id) ?? false;
                        return (
                          <label key={a.amenity_id} className="flex items-center gap-2 text-sm">
                            <input
                              type="checkbox"
                              className="accent-primary"
                              checked={checked}
                              onChange={() => toggleAmenitySelected(property.id, a.amenity_id)}
                            />
                            <span>{a.name}</span>
                          </label>
                        );
                      })}
                    </div>
                    <div>
                      <Button size="sm" onClick={() => addAmenitiesToProperty(property.id)}>Add Selected Amenities</Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>

          <div className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Create Amenity</CardTitle>
              </CardHeader>
              <CardContent>
                <form className="flex gap-2" onSubmit={createAmenity}>
                  <Input
                    placeholder="Amenity name"
                    value={newAmenityName}
                    onChange={(e) => setNewAmenityName(e.target.value)}
                  />
                  <Button type="submit" disabled={isAmenitySubmitting}>
                    {isAmenitySubmitting ? "Adding..." : "Add"}
                  </Button>
                </form>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>All Amenities</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="flex flex-wrap gap-2">
                  {amenities.map((a) => (
                    <span
                      key={a.amenity_id}
                      className="px-2 py-1 rounded border text-sm"
                    >
                      {a.name}
                    </span>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
