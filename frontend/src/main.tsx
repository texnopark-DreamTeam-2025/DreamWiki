import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import App from "./App.tsx";
import { createConfig } from "./client/client";

createConfig({
  baseURL: "https://dreamwiki.zhugeo.ru",
});

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <App />
  </StrictMode>
);
