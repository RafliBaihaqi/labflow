import { Navigate } from "react-router";
import { useAuthStore } from "~/lib/store/auth-store";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const isLoading = useAuthStore((s) => s.isLoading);

  if (isLoading)
    return (
      <div className="flex h-screen items-center justify-center">Loading...</div>
    );
  if (!isAuthenticated) return <Navigate to="/login" replace />;

  return <>{children}</>;
}
