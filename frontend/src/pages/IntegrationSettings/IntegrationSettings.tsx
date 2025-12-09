import { useState } from "react";
import { Select, Button, Flex, Text } from "@gravity-ui/uikit";
import { ywikiAddPage, githubAccountPr } from "@/client";
import { useToast } from "@/hooks/useToast";
import { MonacoEditor } from "@/components/MonacoEditor";

export default function IntegrationSettings() {
  const { showSuccess, showError } = useToast();
  const [selectedIntegration, setSelectedIntegration] = useState<string[]>(["GitHub"]);
  const [prUrl, setPrUrl] = useState<string>("");
  const [pageUrl, setPageUrl] = useState<string>("");
  const [indexedPagesCount, setIndexedPagesCount] = useState<number>(15);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isAddingPage, setIsAddingPage] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);
  const [configContent, setConfigContent] = useState<string>(`{
  "token": "test",
  "url": "http://localhost:8080"
}`);

  const integrationOptions = [
    { value: "GitHub", content: "GitHub" },
    { value: "YandexWiki", content: "Yandex Wiki" },
  ];

  const handleAddPage = async () => {
    if (!pageUrl.trim()) return;

    setIsAddingPage(true);
    try {
      const response = await ywikiAddPage({
        body: {
          page_url: pageUrl,
        },
      });

      if (response.error) {
        console.error("Ошибка добавления страницы:", response.error);
        showError("Ошибка", "Не удалось добавить страницу. Проверьте URL и попробуйте снова.");
        return;
      }

      console.log("Страница успешно добавлена:", response.data);
      setPageUrl("");
      // Увеличиваем счетчик проиндексированных страниц
      setIndexedPagesCount((prev) => prev + 1);
      showSuccess("Успешно", "Создана задача на индексирование страницы");
    } catch (error) {
      console.error("Ошибка добавления страницы:", error);
      showError("Ошибка", "Произошла ошибка при добавлении страницы.");
    } finally {
      setIsAddingPage(false);
    }
  };

  const handleUpdateIndex = async () => {
    setIsUpdating(true);
    try {
      // TODO: Здесь должен быть вызов API для обновления индекса
      await new Promise((resolve) => setTimeout(resolve, 3000));
      console.log("Обновляем индекс");
      showSuccess("Успешно", "Индекс успешно обновлен!");
    } catch (error) {
      console.error("Ошибка обновления индекса:", error);
      showError("Ошибка", "Не удалось обновить индекс.");
    } finally {
      setIsUpdating(false);
    }
  };

  const handleAnalyzeLink = async () => {
    if (!prUrl.trim()) return;

    setIsAnalyzing(true);
    try {
      const response = await githubAccountPr({
        body: {
          pr_url: prUrl,
        },
      });

      if (response.error) {
        console.error("Ошибка анализа PR:", response.error);
        showError("Ошибка", "Не удалось проанализировать PR. Проверьте ссылку и попробуйте снова.");
        return;
      }

      console.log("PR успешно проанализирован:", response.data);
      showSuccess("Успешно", "Создана задача на обновление информации на основе PR");
    } catch (error) {
      console.error("Ошибка анализа PR:", error);
      showError("Ошибка", "Произошла ошибка при анализе PR.");
    } finally {
      setIsAnalyzing(false);
    }
  };

  const handleCheckSettings = async () => {
    setIsLoading(true);
    try {
      // TODO: Здесь должен быть вызов API для проверки настроек
      await new Promise((resolve) => setTimeout(resolve, 1000));
      console.log("Проверяем настройки");
      showSuccess("Успешно", "Настройки корректны!");
    } catch (error) {
      console.error("Ошибка проверки:", error);
      showError("Ошибка", "Не удалось проверить настройки.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleApplySettings = async () => {
    setIsLoading(true);
    try {
      // TODO: Здесь должен быть вызов API для применения настроек
      await new Promise((resolve) => setTimeout(resolve, 1500));
      console.log("Применяем настройки");
      showSuccess("Успешно", "Настройки успешно применены!");
    } catch (error) {
      console.error("Ошибка применения:", error);
      showError("Ошибка", "Не удалось применить настройки.");
    } finally {
      setIsLoading(false);
    }
  };

  // Получаем текущий выбор интеграции
  const currentIntegration = selectedIntegration[0] || "GitHub";
  const isYandexWiki = currentIntegration === "YandexWiki";

  return (
    <Flex direction="column" gap="5" style={{ padding: '20px' }}>
      <Flex justifyContent="space-between" alignItems="center">
        <Text variant="header-1">Настройки интеграций</Text>
        <div style={{ width: '340px' }}>
          <Select
            value={selectedIntegration}
            onUpdate={setSelectedIntegration}
            options={integrationOptions}
            placeholder="Выберите интеграцию"
            width={340}
          />
        </div>
      </Flex>

      {isYandexWiki ? (
        <Flex direction="column" gap="4" style={{ width: '593px' }}>
          <Text variant="subheader-2">Добавьте новую страницу</Text>
          <Flex gap="2">
            <div style={{
              border: '1px solid rgba(0, 0, 0, 0.3)',
              borderRadius: '6px',
              padding: '2px 8px',
              flex: 1,
              display: 'flex',
              alignItems: 'center'
            }}>
              <input
                type="text"
                placeholder="Введите ссылку на страницу"
                value={pageUrl}
                onChange={(e) => setPageUrl(e.target.value)}
                style={{
                  flex: 1,
                  border: 'none',
                  outline: 'none',
                  background: 'transparent',
                  fontFamily: 'var(--text-body-1-font-family, "Inter-Regular", sans-serif)',
                  fontSize: 'var(--text-body-1-font-size, 13px)',
                  lineHeight: 'var(--text-body-1-line-height, 18px)',
                  fontWeight: 'var(--text-body-1-font-weight, 400)',
                  color: 'var(--var-g-color-text-primary, rgba(0, 0, 0, 0.85))'
                }}
              />
            </div>
            <Button
              view="normal"
              size="m"
              onClick={handleAddPage}
              disabled={!pageUrl.trim() || isAddingPage}
              loading={isAddingPage}
            >
              {isAddingPage ? "Добавляем..." : "Добавить"}
            </Button>
          </Flex>

          <Flex direction="column" gap="3">
            <Text variant="subheader-2">Состояние индекса</Text>
            <Flex justifyContent="space-between" alignItems="center">
              <Text variant="body-1">{indexedPagesCount} страниц проиндексировано</Text>
              <Button
                view="flat"
                size="s"
                onClick={handleUpdateIndex}
                disabled={isUpdating}
                loading={isUpdating}
              >
                {isUpdating ? "Обновляем..." : "обновить"}
              </Button>
            </Flex>
          </Flex>
        </Flex>
      ) : (
        <Flex direction="column" gap="4" style={{ width: '593px' }}>
          <Text variant="subheader-2">Внести изменения на основе кода Pull-Request-а</Text>
          <Flex gap="2">
            <div style={{
              border: '1px solid rgba(0, 0, 0, 0.3)',
              borderRadius: '6px',
              padding: '2px 8px',
              flex: 1,
              display: 'flex',
              alignItems: 'center'
            }}>
              <input
                type="text"
                placeholder="Введите ссылку на PR"
                value={prUrl}
                onChange={(e) => setPrUrl(e.target.value)}
                style={{
                  flex: 1,
                  border: 'none',
                  outline: 'none',
                  background: 'transparent',
                  fontFamily: 'var(--text-body-1-font-family, "Inter-Regular", sans-serif)',
                  fontSize: 'var(--text-body-1-font-size, 13px)',
                  lineHeight: 'var(--text-body-1-line-height, 18px)',
                  fontWeight: 'var(--text-body-1-font-weight, 400)',
                  color: 'var(--var-g-color-text-primary, rgba(0, 0, 0, 0.85))'
                }}
              />
            </div>
            <Button
              view="normal"
              size="m"
              onClick={handleAnalyzeLink}
              disabled={!prUrl.trim() || isAnalyzing}
              loading={isAnalyzing}
            >
              {isAnalyzing ? "Анализируем..." : "Проанализировать"}
            </Button>
          </Flex>
        </Flex>
      )}

      <Flex direction="column" gap="3" style={{ width: '654px' }}>
        <Text variant="body-1">Параметры</Text>
        <div style={{
          height: '400px',
          overflow: 'hidden'
        }}>
          <MonacoEditor
            value={configContent}
            onChange={(value) => setConfigContent(value || "")}
            language="json"
            height="100%"
            theme="light"
            readOnly={false}
          />
        </div>
      </Flex>

      <Flex gap="5">
        <Button
          view="outlined"
          size="m"
          onClick={handleCheckSettings}
          disabled={isLoading}
          loading={isLoading}
        >
          {isLoading ? "Проверяем..." : "Проверить"}
        </Button>
        <Button
          view="action"
          size="m"
          onClick={handleApplySettings}
          disabled={isLoading}
          loading={isLoading}
        >
          {isLoading ? "Применяем..." : "Применить"}
        </Button>
      </Flex>
    </Flex>
  );
}
