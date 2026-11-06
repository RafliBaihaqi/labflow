import { MainLayout } from "../../components/main/main-layout";

export default function DashboardPage() {
  return (
    <MainLayout breadcrumbs={[{ label: "Dashboard" }]}>
      <div className="grid auto-rows-min gap-4 md:grid-cols-3">
        <div className="aspect-video rounded-xl border border-border bg-card shadow-sm" />
        <div className="aspect-video rounded-xl border border-border bg-card shadow-sm" />
        <div className="aspect-video rounded-xl border border-border bg-card shadow-sm" />
      </div>
      <div className="min-h-[60vh] flex-1 rounded-xl border border-border bg-card shadow-sm" />
    </MainLayout>
  );
}
