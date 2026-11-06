import { useEffect, useRef, useCallback } from "react";

declare global {
  interface Window {
    turnstile?: {
      render: (
        container: string | HTMLElement,
        options: {
          sitekey: string;
          callback: (token: string) => void;
          "expired-callback"?: () => void;
          "error-callback"?: () => void;
          theme?: "light" | "dark" | "auto";
          size?: "normal" | "compact";
        }
      ) => string;
      reset: (widgetId: string) => void;
      remove: (widgetId: string) => void;
      getResponse: (widgetId: string) => string | undefined;
    };
  }
}

interface UseTurnstileOptions {
  siteKey: string;
  onSuccess: (token: string) => void;
  onExpire?: () => void;
  theme?: "light" | "dark" | "auto";
}

/**
 * Renders Cloudflare Turnstile inside `containerRef` without any third-party
 * package. Handles race where script loads before component mounts by polling
 * window.turnstile and listening for load events.
 * Hydration safe: only touches DOM in useEffect (client-only).
 */
export function useTurnstile({
  siteKey,
  onSuccess,
  onExpire,
  theme = "light",
}: UseTurnstileOptions) {
  const containerRef = useRef<HTMLDivElement>(null);
  const widgetIdRef = useRef<string | null>(null);

  const renderWidget = useCallback(() => {
    if (!siteKey) {
      console.warn("[turnstile] siteKey missing, skipping render");
      return;
    }
    if (!containerRef.current || !window.turnstile) return;
    // Remove any previously rendered widget before re-rendering.
    if (widgetIdRef.current !== null) {
      try {
        window.turnstile.remove(widgetIdRef.current);
      } catch {}
      widgetIdRef.current = null;
    }
    try {
      widgetIdRef.current = window.turnstile.render(containerRef.current, {
        sitekey: siteKey,
        theme,
        callback: onSuccess,
        "expired-callback": onExpire,
      });
    } catch (err) {
      console.error("[turnstile] render failed", err);
    }
  }, [siteKey, theme, onSuccess, onExpire]);

  useEffect(() => {
    if (!siteKey) {
      console.warn("[turnstile] siteKey is empty, Turnstile will not render. Set VITE_TURNSTILE_SITE_KEY");
      return;
    }

    // If turnstile already loaded, render immediately
    if (window.turnstile) {
      renderWidget();
      return;
    }

    // Otherwise poll for it (handles async script without onload callback)
    let intervalId: number | undefined;
    const checkAndRender = () => {
      if (window.turnstile) {
        renderWidget();
        if (intervalId !== undefined) {
          clearInterval(intervalId);
        }
      }
    };

    intervalId = window.setInterval(checkAndRender, 500);

    // Also listen for custom event if some other script dispatches it
    const handler = () => renderWidget();
    window.addEventListener("turnstileLoaded", handler as EventListener);

    return () => {
      if (intervalId !== undefined) clearInterval(intervalId);
      window.removeEventListener("turnstileLoaded", handler as EventListener);
      if (widgetIdRef.current !== null) {
        try {
          window.turnstile?.remove(widgetIdRef.current);
        } catch {}
        widgetIdRef.current = null;
      }
    };
  }, [renderWidget, siteKey]);

  /** Call this after a failed submission to get a fresh token. */
  const reset = useCallback(() => {
    if (widgetIdRef.current !== null) {
      try {
        window.turnstile?.reset(widgetIdRef.current);
      } catch {}
    }
  }, []);

  return { containerRef, reset };
}
