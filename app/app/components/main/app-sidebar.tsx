"use client";

import * as React from "react";
import {FlaskConical, LayoutDashboard, ClipboardList, PlusCircle, User, ClipboardCheck, ListChecks} from "lucide-react";
import {Link} from "react-router";

import {NavMain, type NavItem} from "./nav-main";
import {NavUser} from "./nav-user";
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarRail,
    SidebarTrigger,
} from "../ui/sidebar";
import {useAuthStore} from "~/lib/store/auth-store";

const navItems: (NavItem & { permission?: string })[] = [
    {
        title: "Dashboard",
        url: "/dashboard",
        icon: LayoutDashboard,

    },
    {
        title: "Reservations",
        url: "/reservations",
        icon: ClipboardList,
        permission: "read_reservation"
    },
    {
        title: "New Reservation",
        url: "/reservations/new",
        icon: PlusCircle,
        permission: "create_reservation"
    },
    {
        title: "Profile",
        url: "/profile",
        icon: User,
    },
    {
        title: "Quotations",
        url: "/admin/quotations",
        icon: ClipboardCheck,
        permission: "read_quotation",
    },
    {
        title: "Test List",
        url: "/admin/test-list",
        icon: ListChecks,
        permission: "read_test",
    },
];

export function AppSidebar({...props}: React.ComponentProps<typeof Sidebar>) {
    const { user, hasPermission } = useAuthStore((s) => ({
        user: s.user,
        hasPermission: s.hasPermission,
    }));
    const filteredNav = React.useMemo(
        () => navItems.filter((item) => !item.permission || hasPermission(item.permission)),
        [hasPermission]
    );
    const sidebarUser = user
        ? {name: user.name, email: user.email, avatar: user.avatar || "/avatars/default.jpg"}
        : {name: "Loading...", email: "", avatar: ""};


    return (
        <Sidebar collapsible="icon" variant="inset" {...props}>
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem className="flex items-center gap-1">
                        <SidebarMenuButton size="lg" asChild className="flex-1">
                            <Link to="/dashboard">
                                <div
                                    className="flex aspect-square size-9 items-center justify-center rounded-full bg-sidebar-primary text-sidebar-primary-foreground">
                                    <FlaskConical className="size-4"/>
                                </div>
                                <div className="grid flex-1 text-left text-sm leading-tight">
                  <span className="truncate font-heading font-medium">
                    LabFlow
                  </span>
                                    <span className="truncate text-xs text-sidebar-foreground/60">
                    Lab Management
                  </span>
                                </div>
                            </Link>
                        </SidebarMenuButton>
                        <SidebarTrigger className="text-sidebar-foreground/70 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground group-data-[collapsible=icon]:hidden"/>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>
            <SidebarContent>
                <NavMain items={filteredNav}/>
            </SidebarContent>
            <SidebarFooter>
                <NavUser user={sidebarUser}/>
            </SidebarFooter>
            <SidebarRail/>
        </Sidebar>
    );
}
