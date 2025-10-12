import "@gravity-ui/uikit/styles/fonts.css";
import "@gravity-ui/uikit/styles/styles.css";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import { createConfig } from "./client/client";
import { ThemeProvider } from "@gravity-ui/uikit";
import { client } from "./client/client.gen.ts";

client.setConfig(createConfig({
  baseURL: import.meta.env.VITE_BASE_URL || "https://did.you.forget.to.add.vite.base.url.to.config/?",
}));

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider theme="light">
      <App />
    </ThemeProvider>
  </StrictMode>
);
