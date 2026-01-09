import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";
import fs from "fs";
import tailwindcss from "@tailwindcss/vite";

// https://vite.dev/config/
export default defineConfig({
  envDir: "..",
  server: {
    https:
      fs.existsSync(path.resolve(__dirname, "ssl/key.pem")) &&
      fs.existsSync(path.resolve(__dirname, "ssl/cert.pem"))
        ? {
            key: fs.readFileSync(path.resolve(__dirname, "ssl/key.pem")),
            cert: fs.readFileSync(path.resolve(__dirname, "ssl/cert.pem")),
          }
        : {},
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src/"),
      Public: path.resolve(__dirname, "public"),
    },
  },
  plugins: [react(), tailwindcss()],
  build: {
    outDir: "../infra/nginx/frontend-dist",
  },
});
