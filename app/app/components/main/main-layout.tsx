import * as React from "react";
import { Link } from "react-router";

import { AppSidebar } from "./app-sidebar";
import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "../ui/breadcrumb";
import {
  SidebarInset,
  SidebarProvider,
} from "../ui/sidebar";

export interface Crumb {
  label: string;
  to?: string;
}

interface MainLayoutProps {
  breadcrumbs: Crumb[];
  children: React.ReactNode;
}

export function MainLayout({ breadcrumbs, children }: MainLayoutProps) {
  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        <header className="flex h-20 shrink-0 items-center px-4 text-sidebar-foreground md:px-6">
          <Breadcrumb>
            <BreadcrumbList className="text-base text-sidebar-foreground/50">
              {breadcrumbs.map((crumb, index) => {
                const isLast = index === breadcrumbs.length - 1;
                return (
                  <React.Fragment key={`${crumb.label}-${index}`}>
                    <BreadcrumbItem>
                      {isLast || !crumb.to ? (
                        <BreadcrumbPage className="font-heading text-2xl font-semibold text-sidebar-foreground">
                          {crumb.label}
                        </BreadcrumbPage>
                      ) : (
                        <BreadcrumbLink
                          className="hover:text-sidebar-foreground"
                          asChild
                        >
                          <Link to={crumb.to}>{crumb.label}</Link>
                        </BreadcrumbLink>
                      )}
                    </BreadcrumbItem>
                    {!isLast && <BreadcrumbSeparator />}
                  </React.Fragment>
                );
              })}
            </BreadcrumbList>
          </Breadcrumb>
        </header>
        <div className="mx-2 mb-2 flex flex-1 flex-col gap-6 rounded-xl border border-border bg-card p-6 shadow-sm md:mx-3 md:mb-3">
          {children}
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
}
