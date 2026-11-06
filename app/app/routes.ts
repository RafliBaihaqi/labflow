import { type RouteConfig, index, route } from "@react-router/dev/routes";

export default [
  index("routes/home.tsx"),
  route("login", "routes/login.tsx"),
  route("register", "routes/register.tsx"),
  route("verify-otp", "routes/verify-otp.tsx"),
  route("forgot-password", "routes/forgot-password.tsx"),
  route("verify-reset-otp", "routes/verify-reset-otp.tsx"),
  route("reset-password", "routes/reset-password.tsx"),
  route("dashboard", "routes/dashboard.tsx"),
  route("profile", "routes/profile.tsx"),
  route("reservations", "routes/reservations-list.tsx"),
  route("reservations/new", "routes/reservations-new.tsx"),
  route("reservations/:xid", "routes/reservations-detail.tsx"),
  route("admin/quotations", "routes/admin-quotations.tsx"),
  route("admin/test-list", "routes/admin-test-list.tsx"),
  route("examples/data-table", "routes/data-table-demo.tsx"),
] satisfies RouteConfig;
