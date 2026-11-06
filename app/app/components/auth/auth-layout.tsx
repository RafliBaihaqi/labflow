import type React from "react";
import { Button } from "~/components/ui/button";
import { Card } from "~/components/ui/card";
import { Languages } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { useT } from "~/i18n/use-t";
import { useLocale } from "~/contexts/locale-context";
import { LogoMark } from "~/components/main/logo-mark";

interface AuthLayoutProps {
  title: React.ReactNode;
  description: React.ReactNode;
  children: React.ReactNode;
  className?: string;
}

export function AuthLayout({
  title,
  description,
  children,
  className,
}: AuthLayoutProps) {
  const t = useT();
  const { locale, setLocale } = useLocale();

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-brand-baby-blue/30 via-white to-brand-light-orange/20 relative overflow-hidden py-2 lg:py-12">
      {/* Background Elements */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute -top-40 -right-40 w-80 h-80 bg-brand-baby-blue/20 rounded-full blur-3xl"></div>
        <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-brand-light-orange/20 rounded-full blur-3xl"></div>
        <div className="absolute top-1/4 left-1/4 w-32 h-32 bg-brand-ocean-blue/10 rounded-full blur-2xl"></div>
      </div>

      {/* Language Selector - positioned in top-right corner */}
      <div className="absolute top-4 right-4 z-20">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              size="sm"
              className="flex items-center space-x-2 hover:bg-brand-light-orange/20 bg-white/80 backdrop-blur-sm shadow-sm border border-brand-baby-blue/20"
            >
              <Languages className="h-4 w-4 text-brand-dark-grey/70" />
              <span className="text-sm text-brand-dark-grey/70">
                {locale === "id" ? "🇮🇩" : "🇺🇸"}
              </span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            align="end"
            className="w-48 border-brand-baby-blue/20"
          >
            <DropdownMenuLabel className="text-brand-dark-grey">
              {t("language.selectLanguage")}
            </DropdownMenuLabel>
            <DropdownMenuSeparator className="bg-brand-baby-blue/20" />
            <DropdownMenuItem
              onClick={() => setLocale("id")}
              className={`hover:bg-brand-light-orange/20 cursor-pointer ${locale === "id" ? "bg-brand-light-orange/10" : ""}`}
            >
              <span className="mr-2">🇮🇩</span>
              <span className="text-brand-dark-grey">Bahasa Indonesia</span>
              {locale === "id" && (
                <span className="ml-auto text-brand-ocean-blue">✓</span>
              )}
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => setLocale("en")}
              className={`hover:bg-brand-light-orange/20 cursor-pointer ${locale === "en" ? "bg-brand-light-orange/10" : ""}`}
            >
              <span className="mr-2">🇺🇸</span>
              <span className="text-brand-dark-grey">English</span>
              {locale === "en" && (
                <span className="ml-auto text-brand-ocean-blue">✓</span>
              )}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      <div className="relative z-10 max-w-5xl mx-8">
        <Card className="p-0 border-brand-baby-blue/30 shadow-2xl bg-white/95 backdrop-blur-sm overflow-hidden grid grid-cols-1 md:grid-cols-2">
          {/* Left Side - Dynamic Header + Form */}
          <div className="p-3 lg:p-6 flex flex-col justify-center">
            <div className="text-center pb-2 lg:pb-6 pt-0 px-0 lg:px-6">
              <div className="flex justify-center mb-1 py-2 2xl:py-6">
                <LogoMark className="h-20 w-20" />
              </div>
              <h1 className="text-base 2xl:text-2xl font-bold text-brand-dark-grey mb-0 lg:mb-2">
                {title}
              </h1>
              <div className="text-brand-dark-grey/70 text-sm 2xl:text-base leading-tight 2xl:leading-normal">
                {description}
              </div>
            </div>
            <div className={className}>
              {children}
            </div>
          </div>

          {/* Right Side - Image Column */}
          <div className="hidden md:block relative overflow-hidden">
            <div className="relative w-full h-full">
              <div className="h-full w-full bg-gradient-to-br from-sky-500 via-sky-700 to-slate-900" />
              <div className="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent"></div>
              <div className="absolute bottom-6 left-0 right-0 text-center z-10">
                <p className="text-white text-sm font-medium">
                  {t("auth.login.footer")}
                </p>
              </div>
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
