import "./styles.scss";
import "@gravity-ui/uikit/styles/fonts.css";
import "@gravity-ui/uikit/styles/styles.css";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import { createConfig } from "./client/client";
import {
  ThemeProvider,
  ToasterComponent,
  ToasterProvider,
} from "@gravity-ui/uikit";
import { toaster } from "@gravity-ui/uikit/toaster-singleton";
import { client } from "./client/client.gen.ts";
import { AuthProvider } from "./contexts/AuthContext.tsx";

client.setConfig(
  createConfig({
    baseURL:
      import.meta.env.VITE_BASE_URL ||
      "https://did.you.forget.to.add.vite.base.url.to.config/?",
  })
);

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider theme="light">
      <ToasterProvider toaster={toaster}>
        <AuthProvider>
          <App />
        </AuthProvider>
        <ToasterComponent />
      </ToasterProvider>
    </ThemeProvider>
  </StrictMode>
);
