import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import AsideBar from "@/components/AsideBar";
import SearchDone from "@/pages/SearchDone";
import Document from "@/pages/Document/Document";
import IntegrationSettings from "@/pages/IntegrationSettings";
import IntegrationLogs from "@/pages/IntegrationLogs";
import Authorization from "@/pages/Authorization";
import ProtectedRoute from "@/components/ProtectedRoute";
import PublicRoute from "@/components/PublicRoute";
import { HomePage } from "@/pages/HomePage";

function App() {
  return (
    <Router>
      <AsideBar>
        <Routes>
          <Route
            path="/login"
            element={
              <PublicRoute>
                <Authorization />
              </PublicRoute>
            }
          />
          <Route
            path="/"
            element={
              <ProtectedRoute>
                <HomePage />
              </ProtectedRoute>
            }
          />
          <Route
            path="/search"
            element={
              <ProtectedRoute>
                <SearchDone />
              </ProtectedRoute>
            }
          />
          <Route
            path="/search-done"
            element={
              <ProtectedRoute>
                <SearchDone />
              </ProtectedRoute>
            }
          />
          <Route
            path="/document/:id"
            element={
              <ProtectedRoute>
                <Document />
              </ProtectedRoute>
            }
          />
          <Route
            path="/integration-settings"
            element={
              <ProtectedRoute>
                <IntegrationSettings />
              </ProtectedRoute>
            }
          />
          <Route
            path="/integration-logs"
            element={
              <ProtectedRoute>
                <IntegrationLogs />
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<div>Страница не найдена</div>} />
        </Routes>
      </AsideBar>
    </Router>
  );
}

export default App;
