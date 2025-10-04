import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

// https://vite.dev/config/
export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "src/"),
      Public: path.resolve(__dirname, "public"),
    },
  },
  plugins: [react()],
});
