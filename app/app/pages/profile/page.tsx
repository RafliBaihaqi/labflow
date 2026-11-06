import { MainLayout } from "~/components/main/main-layout";
import { ProfileForm } from "~/components/main/profile-form";
import { useAuthStore } from "~/lib/store/auth-store";

export default function ProfilePage() {
  const user = useAuthStore((s) => s.user);
  const isLoading = useAuthStore((s) => s.isLoading);

  return (
    <MainLayout breadcrumbs={[{ label: "Profile" }]}>
      <div className="mx-auto w-full max-w-3xl">
        {isLoading || !user ? (
          <p className="text-sm text-muted-foreground">Loading profile...</p>
        ) : (
          <ProfileForm user={user} />
        )}
      </div>
    </MainLayout>
  );
}
