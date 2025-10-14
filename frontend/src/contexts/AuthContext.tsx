import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";
import { client } from "@/client/client.gen";

interface AuthContextType {
  isAuthenticated: boolean;
  login: (token: string) => void;
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    // Check if there's a stored token on initial load
    const token = localStorage.getItem("authToken");
    if (token) {
      setIsAuthenticated(true);
      // Set the token in the client config
      client.setConfig({
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });
    }
  }, []);

  const login = (token: string) => {
    setIsAuthenticated(true);
    localStorage.setItem("authToken", token);
    // Set the token in the client config
    client.setConfig({
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  };

  const logout = () => {
    setIsAuthenticated(false);
    localStorage.removeItem("authToken");
    // Clear the token from the client config
    client.setConfig({
        headers: {}
    });
  };

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
