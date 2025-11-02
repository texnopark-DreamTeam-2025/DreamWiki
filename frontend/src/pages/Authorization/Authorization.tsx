import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, TextInput, Flex, Text } from "@gravity-ui/uikit";
import { login } from "@/client/sdk.gen";
import { useAuth } from "@/contexts/AuthContext";

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
    <Flex justifyContent="center" alignItems="center" style={{ height: '100%', padding: 20 }}>
      <Flex
        direction="column"
        style={{
          width: '100%',
          maxWidth: 400,
          padding: 30,
          borderRadius: 8,
          backgroundColor: 'var(--g-color-base-background)',
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.1)'
        }}
      >
        <Text variant="display-1" style={{ textAlign: 'center', marginBottom: 24 }}>
          Authorization
        </Text>
        <form onSubmit={handleSubmit}>
          <Flex direction="column" style={{ marginBottom: 20 }}>
            <TextInput
              label="Username"
              value={username}
              onUpdate={setUsername}
              disabled={loading}
            />
          </Flex>
          <Flex direction="column" style={{ marginBottom: 20 }}>
            <TextInput
              label="Password"
              type="password"
              value={password}
              onUpdate={setPassword}
              disabled={loading}
            />
          </Flex>
          {error && (
            <Flex
              justifyContent="center"
              alignItems="center"
              style={{
                color: 'var(--g-color-text-danger)',
                textAlign: 'center',
                marginBottom: 16,
                padding: 8,
                borderRadius: 4,
                backgroundColor: 'var(--g-color-base-danger-light)'
              }}
            >
              {error}
            </Flex>
          )}
          <Flex justifyContent="center" style={{ marginTop: 24 }}>
            <Button type="submit" loading={loading} size="l" width="max">
              Sign In
            </Button>
          </Flex>
        </form>
      </Flex>
    </Flex>
  );
}
