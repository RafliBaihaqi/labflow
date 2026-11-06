"use client";

import { type LucideIcon } from "lucide-react";
import { Link, useLocation } from "react-router";

import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "../ui/sidebar";

export interface NavItem {
  title: string;
  url: string;
  icon?: LucideIcon;
}

/** Returns the url of the single best-matching item for the current path. */
function activeUrl(pathname: string, items: NavItem[]): string | undefined {
  let best: string | undefined;
  for (const item of items) {
    const matches =
      pathname === item.url || pathname.startsWith(`${item.url}/`);
    if (matches && (best === undefined || item.url.length > best.length)) {
      best = item.url;
    }
  }
  return best;
}

export function NavMain({ items }: { items: NavItem[] }) {
  const { pathname } = useLocation();
  const current = activeUrl(pathname, items);

  return (
    <SidebarGroup>
      <SidebarMenu className="gap-1">
        {items.map((item) => (
          <SidebarMenuItem key={item.title}>
            <SidebarMenuButton
              asChild
              tooltip={item.title}
              isActive={item.url === current}
              className="h-10 gap-3 rounded-lg px-3 text-sm font-medium text-sidebar-foreground/80 transition-colors data-[active=true]:bg-sidebar-accent data-[active=true]:text-sidebar-accent-foreground data-[active=true]:shadow-sm [&_svg]:size-5"
            >
              <Link to={item.url}>
                {item.icon && <item.icon />}
                <span>{item.title}</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
}
