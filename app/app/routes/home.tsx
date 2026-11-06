import LabFlowLanding from "~/components/home/landing-page";
import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "LabFlow — Laboratory Information Management" },
    { name: "description", content: "Open-source LIMS: test requests, role-guarded lifecycle, certificates." },
  ];
}

export default function Home() {
  return <LabFlowLanding />;
}
