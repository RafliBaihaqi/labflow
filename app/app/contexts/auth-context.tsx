import { createContext, useContext, useEffect, useMemo, useCallback } from "react";
import { useQuery } from "@tanstack/react-query";
import { me } from "~/lib/api-client";
import { useAuthStore } from "~/lib/store/auth-store";
import type { User } from "~/lib/types/auth";

interface AuthContextType {
  user: User | null;
  isLoading: boolean;
  isAuthenticated: boolean;
  hasPermission: (perm: string) => boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data: user, isPending } = useQuery({
    queryKey: ["me"],
    queryFn: me,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });

  const isLoading = isPending;

  const setUser = useAuthStore((s) => s.setUser);
  const setLoading = useAuthStore((s) => s.setLoading);

  useEffect(() => {
    setLoading(isLoading);
    if (!isLoading) {
      setUser(user ?? null);
    }
  }, [user, isLoading, setUser, setLoading]);

  const hasPermission = useCallback(
    (perm: string) => user?.permissions.includes(perm) ?? false,
    [user]
  );

  const value = useMemo<AuthContextType>(
    () => ({
      user: user ?? null,
      isLoading,
      isAuthenticated: !!user,
      hasPermission,
    }),
    [user, isLoading, hasPermission]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}