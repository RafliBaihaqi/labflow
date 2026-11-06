import {
  createContext,
  useContext,
  useEffect,
  useState,
  useMemo,
  useCallback,
  type ReactNode,
} from "react";
import { IntlProvider } from "react-intl";
import {
  type SupportedLocale,
  SUPPORTED_LOCALES,
  DEFAULT_LOCALE,
  getMessages,
  detectLocale,
} from "~/i18n/config";

interface LocaleContextValue {
  locale: SupportedLocale;
  setLocale: (locale: SupportedLocale) => void;
}

const LocaleContext = createContext<LocaleContextValue>({
  locale: DEFAULT_LOCALE,
  setLocale: () => {},
});

export function LocaleProvider({ children }: { children: ReactNode }) {
  const [locale, setLocaleState] = useState<SupportedLocale>(() => {
    if (typeof window === "undefined") return DEFAULT_LOCALE;
    try {
      const stored = localStorage.getItem("locale");
      if (stored && SUPPORTED_LOCALES.includes(stored as SupportedLocale)) {
        return stored as SupportedLocale;
      }
    } catch {
      // SSR or storage disabled
    }
    return DEFAULT_LOCALE;
  });

  // Hydrate from actual storage after mount to avoid SSR mismatch
  useEffect(() => {
    try {
      const stored = localStorage.getItem("locale");
      if (stored && SUPPORTED_LOCALES.includes(stored as SupportedLocale)) {
        if (stored !== locale) setLocaleState(stored as SupportedLocale);
        return;
      }
      const detected = detectLocale();
      if (detected !== locale) setLocaleState(detected);
    } catch {
      // ignore
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const setLocale = useCallback((newLocale: SupportedLocale) => {
    setLocaleState(newLocale);
    try {
      localStorage.setItem("locale", newLocale);
    } catch {
      // ignore
    }
  }, []);

  // Keep <html lang> in sync with the active locale
  useEffect(() => {
    document.documentElement.lang = locale;
  }, [locale]);

  const contextValue = useMemo(
    () => ({ locale, setLocale }),
    [locale, setLocale]
  );

  return (
    <LocaleContext.Provider value={contextValue}>
      <IntlProvider
        locale={locale}
        messages={getMessages(locale)}
        defaultLocale={DEFAULT_LOCALE}
      >
        {children}
      </IntlProvider>
    </LocaleContext.Provider>
  );
}

export function useLocale() {
  const ctx = useContext(LocaleContext);
  if (!ctx) throw new Error("useLocale must be used within LocaleProvider");
  return ctx;
}
