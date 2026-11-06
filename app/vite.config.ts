import { reactRouter } from "@react-router/dev/vite";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
  plugins: [tailwindcss(), reactRouter()],
  resolve: {
    tsconfigPaths: true,
  },
  server: {
    // Use BACKEND_URL env var for flexibility (WSL, docker, remote).
    // Falls back to localhost:8080 for standard local dev.
    // For SSH-tunneled dev on VPS infra: keep localhost, tunnel 8080 if needed.
    proxy: {
      "/api": process.env.BACKEND_URL || "http://localhost:8080",
      "/auth": process.env.BACKEND_URL || "http://localhost:8080",
    },
  },
});
