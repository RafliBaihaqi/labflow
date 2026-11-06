import { useIntl } from "react-intl";

/**
 * Convenience hook for translations.
 *
 * Usage:
 *   const t = useT();
 *   t("auth.login.title")              // simple key
 *   t("greeting", { name: "Alice" })   // with interpolation
 */
export function useT() {
  const intl = useIntl();
  return (id: string, values?: Record<string, string | number>) =>
    intl.formatMessage({ id }, values);
}
