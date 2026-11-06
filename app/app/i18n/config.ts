import en from "./messages/en.json";
import id from "./messages/id.json";

export const SUPPORTED_LOCALES = ["id", "en"] as const;
export type SupportedLocale = (typeof SUPPORTED_LOCALES)[number];
export const DEFAULT_LOCALE: SupportedLocale = "id";

const messages: Record<SupportedLocale, Record<string, string>> = { en, id };

export function getMessages(locale: SupportedLocale): Record<string, string> {
  return messages[locale];
}

/**
 * Detect locale from browser navigator.language.
 * Falls back to DEFAULT_LOCALE ("id") if not supported.
 */
export function detectLocale(): SupportedLocale {
  if (typeof navigator === "undefined") return DEFAULT_LOCALE;
  const browserLang = navigator.language.split("-")[0];
  if (
    SUPPORTED_LOCALES.includes(browserLang as SupportedLocale)
  ) {
    return browserLang as SupportedLocale;
  }
  return DEFAULT_LOCALE;
}
