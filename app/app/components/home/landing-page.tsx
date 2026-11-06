import { Button } from "~/components/ui/button";
import { Card, CardContent } from "~/components/ui/card";
import { Input } from "~/components/ui/input";
import { useForm } from "react-hook-form";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { Link } from "react-router";
import axios from "axios";
import { FlaskConical, ShieldCheck, Workflow, FileCheck2 } from "lucide-react";
import { submitQuotation } from "~/lib/api-client";
import type { QuotationInput } from "~/lib/types/quotation";
import { LogoMark } from "~/components/main/logo-mark";

const features = [
  {
    icon: Workflow,
    title: "Tracked lifecycle",
    body: "Every request moves pending → sample received → testing → completed, with each transition guarded by role.",
  },
  {
    icon: ShieldCheck,
    title: "Role-based access",
    body: "Customers, analysts, QC, QA, lab head and MR each see only the routes and actions their role allows.",
  },
  {
    icon: FlaskConical,
    title: "Test catalog",
    body: "Versioned catalog of microbiology, chemical and physical tests, with per-line results and criteria.",
  },
  {
    icon: FileCheck2,
    title: "Certificates",
    body: "Completed requests issue downloadable certificates and collect customer feedback.",
  },
];

export default function LabFlowLanding() {
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<QuotationInput>();

  const mutation = useMutation({
    mutationFn: submitQuotation,
    onSuccess: () => {
      toast.success("Request sent. We'll be in touch.");
      reset();
    },
    onError: (error: unknown) => {
      let errorMessage = "An unknown error occurred";
      if (axios.isAxiosError(error)) {
        errorMessage = error.response?.data?.message || error.message;
      } else if (error instanceof Error) {
        errorMessage = error.message;
      }
      toast.error("Failed to submit request: " + errorMessage);
    },
  });

  const onSubmit = handleSubmit((data) => mutation.mutate(data));

  return (
    <div className="min-h-screen bg-white text-slate-900">
      <header className="mx-auto flex max-w-6xl items-center justify-between px-6 py-5">
        <div className="flex items-center gap-2 font-semibold">
          <LogoMark className="h-8 w-8" />
          LabFlow
        </div>
        <Link to="/login">
          <Button>Sign in</Button>
        </Link>
      </header>

      <section className="mx-auto max-w-6xl px-6 py-20 text-center">
        <h1 className="text-4xl font-bold tracking-tight md:text-5xl">
          Laboratory requests, from sample to certificate
        </h1>
        <p className="mx-auto mt-5 max-w-2xl text-lg text-slate-600">
          LabFlow is an open-source laboratory information management system.
          Customers submit test requests, lab staff move them through a
          role-guarded lifecycle, and results ship as certificates.
        </p>
        <div className="mt-8 flex justify-center gap-3">
          <Link to="/register">
            <Button size="lg">Create an account</Button>
          </Link>
          <a href="#contact">
            <Button size="lg" variant="outline">
              Request a quotation
            </Button>
          </a>
        </div>
      </section>

      <section className="mx-auto grid max-w-6xl gap-6 px-6 pb-20 sm:grid-cols-2 lg:grid-cols-4">
        {features.map(({ icon: Icon, title, body }) => (
          <Card key={title}>
            <CardContent className="space-y-2 p-6">
              <Icon className="h-6 w-6 text-sky-600" />
              <h2 className="font-semibold">{title}</h2>
              <p className="text-sm text-slate-600">{body}</p>
            </CardContent>
          </Card>
        ))}
      </section>

      <section id="contact" className="bg-slate-50 py-20">
        <div className="mx-auto max-w-xl px-6">
          <h2 className="text-center text-2xl font-bold">Request a quotation</h2>
          <form className="mt-8 space-y-4" onSubmit={onSubmit}>
            <div>
              <Input placeholder="Email" type="text" {...register("email", { required: true })} />
              {errors.email && <p className="mt-1 text-sm text-red-500">Required</p>}
            </div>
            <div>
              <Input
                placeholder="Company name"
                type="text"
                {...register("company_name", { required: true })}
              />
              {errors.company_name && <p className="mt-1 text-sm text-red-500">Required</p>}
            </div>
            <Input placeholder="Phone (optional)" type="tel" {...register("phone_number")} />
            <div>
              <textarea
                placeholder="What would you like tested?"
                rows={4}
                className="w-full rounded-md border border-slate-200 bg-white px-3 py-2 text-sm"
                {...register("message", { required: true })}
              />
              {errors.message && <p className="mt-1 text-sm text-red-500">Required</p>}
            </div>
            <Button type="submit" className="w-full" disabled={mutation.isPending}>
              {mutation.isPending ? "Sending…" : "Send request"}
            </Button>
          </form>
        </div>
      </section>

      <footer className="py-8 text-center text-sm text-slate-500">
        LabFlow · MIT licensed
      </footer>
    </div>
  );
}
