import { defineConfig, devices } from "@playwright/test";

// E2E tests run against a real backend (seeded database). Start it with
// deployments/loadtest (see README) or point BACKEND_URL at your own.
const BACKEND_URL = process.env.BACKEND_URL || "http://localhost:8080";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: false, // tests share seeded accounts and one-session-per-user
  workers: 1,
  retries: process.env.CI ? 1 : 0,
  reporter: process.env.CI ? [["github"], ["html", { open: "never" }]] : "list",
  use: {
    baseURL: process.env.E2E_BASE_URL || "http://localhost:5173",
    trace: "retain-on-failure",
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
  webServer: process.env.E2E_BASE_URL
    ? undefined
    : {
        command: "npm run dev -- --port 5173 --strictPort",
        url: "http://localhost:5173",
        reuseExistingServer: !process.env.CI,
        env: { BACKEND_URL },
        timeout: 120_000,
      },
});
