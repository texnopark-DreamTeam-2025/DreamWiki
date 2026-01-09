import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button, TextInput, Flex, Text, Label } from "@gravity-ui/uikit";
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
        authLogin(response.data.token);
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
    <Flex justifyContent="center" alignItems="center" className="p-4 h-full">
      <Flex
        direction="column"
        className="p-6"
        style={{
          width: "100%",
          maxWidth: 400,
          borderRadius: 8,
          backgroundColor: "var(--g-color-base-background)",
          boxShadow: "0 4px 12px rgba(0, 0, 0, 0.1)",
        }}
      >
        <Text variant="display-1" className="mb-4 text-center">
          Авторизация
        </Text>
        <form onSubmit={handleSubmit}>
          <Flex direction="column" className="mb-4">
            <Label>Логин</Label>
            <TextInput value={username} onUpdate={setUsername} disabled={loading} />
          </Flex>
          <Flex direction="column" className="mb-4">
            <Label>Пароль</Label>
            <TextInput type="password" value={password} onUpdate={setPassword} disabled={loading} />
          </Flex>
          {error && (
            <Flex
              justifyContent="center"
              alignItems="center"
              className="mb-4 p-2"
              style={{
                color: "var(--g-color-text-danger)",
                textAlign: "center",
                borderRadius: 4,
                backgroundColor: "var(--g-color-base-danger-light)",
              }}
            >
              {error}
            </Flex>
          )}
          <Flex justifyContent="center" className="mt-4">
            <Button type="submit" loading={loading} size="l" width="max">
              Войти
            </Button>
          </Flex>
        </form>
      </Flex>
    </Flex>
  );
}
