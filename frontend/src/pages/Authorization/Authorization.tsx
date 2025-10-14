import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, TextInput } from "@gravity-ui/uikit";
import { login } from "@/client/sdk.gen";
import { useAuth } from "@/contexts/AuthContext";
import styles from "./Authorization.module.scss";

export default function Authorization() {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const navigate = useNavigate();
  const { login: authLogin } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");

    try {
      const response = await login({
        body: {
          username,
          password,
        },
      });

      if (response.data) {
        // Use the auth context to handle login
        authLogin(response.data.token);

        // Redirect to home page
        navigate("/");
      } else {
        setError("Invalid response from server");
      }
    } catch (err) {
      setError("Invalid username or password");
      console.error("Login error:", err);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className={styles.authorizationContainer}>
      <div className={styles.authorizationForm}>
        <h1>Authorization</h1>
        <form onSubmit={handleSubmit}>
          <div className={styles.formGroup}>
            <TextInput
              label="Username"
              value={username}
              onUpdate={setUsername}
              disabled={loading}
            />
          </div>
          <div className={styles.formGroup}>
            <TextInput
              label="Password"
              type="password"
              value={password}
              onUpdate={setPassword}
              disabled={loading}
            />
          </div>
          {error && <div className={styles.error}>{error}</div>}
          <div className={styles.formActions}>
            <Button type="submit" loading={loading} size="l" width="max">
              Sign In
            </Button>
          </div>
        </form>
      </div>
    </div>
  );
}
