import "@gravity-ui/uikit/styles/fonts.css";
import "@gravity-ui/uikit/styles/styles.css";
import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import { createConfig } from "./client/client";
import { ThemeProvider } from "@gravity-ui/uikit";

createConfig({
  baseURL: "https://dreamwiki.zhugeo.ru",
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <ThemeProvider theme="light">
      <App />
    </ThemeProvider>
  </StrictMode>
);
