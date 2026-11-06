import { useEffect, useRef, useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { CloudUpload, Loader2 } from "lucide-react";
import { z } from "zod";
import axios from "axios";

import { Avatar, AvatarFallback, AvatarImage } from "~/components/ui/avatar";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Field, FieldError, FieldLabel } from "~/components/ui/field";
import { updateProfile, uploadUserAvatar } from "~/lib/api-client";
import type { User } from "~/lib/types/auth";
import { toast } from "sonner";

const profileSchema = z.object({
  company: z.string().max(255).optional(),
  phone_number: z.string().max(255).optional(),
  address: z.string().max(1000).optional(),
});

type ProfileFormData = z.infer<typeof profileSchema>;

function initials(name: string) {
  return (
    name
      .split(" ")
      .map((n) => n[0])
      .filter(Boolean)
      .slice(0, 2)
      .join("")
      .toUpperCase() || "U"
  );
}

export function ProfileForm({ user }: { user: User }) {
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [avatarFile, setAvatarFile] = useState<File | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(
    user.avatar ?? null,
  );
  const [objectUrl, setObjectUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
    defaultValues: {
      company: user.company ?? "",
      phone_number: user.phone_number ?? "",
      address: user.address ?? "",
    },
  });

  useEffect(() => {
    return () => {
      if (objectUrl) URL.revokeObjectURL(objectUrl);
    };
  }, [objectUrl]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    if (objectUrl) URL.revokeObjectURL(objectUrl);
    const url = URL.createObjectURL(file);
    setObjectUrl(url);
    setAvatarFile(file);
    setAvatarPreview(url);
  };

  const mutation = useMutation({
    mutationFn: async (data: ProfileFormData) => {
      let avatarObject: string | undefined;
      if (avatarFile) {
        const upload = await uploadUserAvatar(avatarFile);
        avatarObject = upload.data.object_name;
      }
      return updateProfile(user.xid, {
        company: data.company,
        phone_number: data.phone_number,
        address: data.address,
        avatar: avatarObject,
      });
    },
    onSuccess: (res) => {
      setError(null);
      setAvatarFile(null);
      if (objectUrl) {
        URL.revokeObjectURL(objectUrl);
        setObjectUrl(null);
      }
      if (res.data.avatar) setAvatarPreview(res.data.avatar);
      queryClient.invalidateQueries({ queryKey: ["me"] });
      toast.success("Profile updated successfully");
    },
    onError: (err) => {
      if (axios.isAxiosError(err)) {
        setError(
          err.response?.data?.message ??
            "Failed to update profile. Please try again.",
        );
      } else {
        setError("Network error. Please check your connection and try again.");
      }
    },
  });

  const onCancel = () => {
    setError(null);
    setAvatarFile(null);
    if (objectUrl) {
      URL.revokeObjectURL(objectUrl);
      setObjectUrl(null);
    }
    setAvatarPreview(user.avatar ?? null);
    reset({
      company: user.company ?? "",
      phone_number: user.phone_number ?? "",
      address: user.address ?? "",
    });
  };

  const onSubmit = (data: ProfileFormData) => {
    setError(null);
    mutation.mutate(data);
  };

  const saving = mutation.isPending;

  const actions = (
    <div className="flex items-center gap-2">
      <Button type="button" variant="outline" size="lg" onClick={onCancel} disabled={saving}>
        Cancel
      </Button>
      <Button type="submit" size="lg" disabled={saving}>
        {saving && <Loader2 className="size-4 animate-spin" data-icon="inline-start" />}
        {saving ? "Saving..." : "Save changes"}
      </Button>
    </div>
  );

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-8">
      {/* Profile header */}
      <div className="flex flex-col items-start gap-4 sm:flex-row sm:items-center">
        <Avatar size="lg" className="size-20">
          <AvatarImage src={avatarPreview ?? undefined} alt={user.name} />
          <AvatarFallback className="text-xl">{initials(user.name)}</AvatarFallback>
        </Avatar>
        <div className="min-w-0 flex-1">
          <h2 className="font-heading truncate text-xl font-semibold">{user.name}</h2>
          <p className="truncate text-sm text-muted-foreground">{user.email}</p>
        </div>
      </div>

      {/* Section title + actions */}
      <div className="flex flex-col gap-3 border-b border-border pb-5 sm:flex-row sm:items-start sm:justify-between">
        <div>
          <h3 className="text-base font-semibold">Profile</h3>
          <p className="text-sm text-muted-foreground">
            Update your photo and personal details here.
          </p>
        </div>
        {/* {actions} */}
      </div>

      {/* Public profile */}
      <div className="grid gap-3 md:grid-cols-3">
        <div>
          <p className="text-sm font-medium">Public profile</p>
          <p className="text-sm text-muted-foreground">
            This will be displayed on your profile.
          </p>
        </div>
        <div className="flex flex-col gap-4 md:col-span-2">
          <Field>
            <FieldLabel htmlFor="company">Company</FieldLabel>
            <Input id="company" type="text" maxLength={255} {...register("company")} />
            {errors.company && <FieldError>{errors.company.message}</FieldError>}
          </Field>
          <Field>
            <FieldLabel htmlFor="phone_number">Phone number</FieldLabel>
            <Input id="phone_number" type="tel" maxLength={255} {...register("phone_number")} />
            {errors.phone_number && <FieldError>{errors.phone_number.message}</FieldError>}
          </Field>
        </div>
      </div>

      {/* Photo */}
      <div className="grid gap-3 border-t border-border pt-6 md:grid-cols-3">
        <div>
          <p className="text-sm font-medium">Photo</p>
          <p className="text-sm text-muted-foreground">
            Update your photo and then choose where you want it to display.
          </p>
        </div>
        <div className="flex items-start gap-5 md:col-span-2">
          <Avatar size="lg" className="size-16">
            <AvatarImage src={avatarPreview ?? undefined} alt={user.name} />
            <AvatarFallback>{initials(user.name)}</AvatarFallback>
          </Avatar>
          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="flex flex-1 cursor-pointer flex-col items-center justify-center gap-1 rounded-xl border border-border bg-background px-4 py-6 text-center transition-colors hover:bg-muted"
          >
            <CloudUpload className="size-5 text-muted-foreground" />
            <span className="text-sm">
              <span className="font-medium text-primary">Click to upload</span>{" "}
              or drag and drop
            </span>
            <span className="text-xs text-muted-foreground">
              SVG, PNG, JPG or GIF (max. 800x400px)
            </span>
          </button>
          <input
            ref={fileInputRef}
            type="file"
            accept="image/*"
            className="hidden"
            onChange={handleFileChange}
          />
        </div>
      </div>

      {/* Address */}
      <div className="grid gap-3 border-t border-border pt-6 md:grid-cols-3">
        <div>
          <p className="text-sm font-medium">Address</p>
          <p className="text-sm text-muted-foreground">
            Where can we reach you.
          </p>
        </div>
        <div className="md:col-span-2">
          <Field>
            <FieldLabel htmlFor="address" className="sr-only">
              Address
            </FieldLabel>
            <textarea
              id="address"
              rows={3}
              maxLength={1000}
              className="w-full min-w-0 rounded-lg border border-input bg-transparent px-3 py-2 text-sm transition-colors outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/30"
              {...register("address")}
            />
            {errors.address && <FieldError>{errors.address.message}</FieldError>}
          </Field>
        </div>
      </div>

      {error && (
        <div
          role="alert"
          className="rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive"
        >
          {error}
        </div>
      )}

      {/* Footer actions */}
      <div className="flex justify-end border-t border-border pt-6">{actions}</div>
    </form>
  );
}
