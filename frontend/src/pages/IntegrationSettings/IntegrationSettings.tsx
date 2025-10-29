import { useState } from "react";
import { Select, Button } from "@gravity-ui/uikit";
import { ywikiAddPage, githubAccountPr } from "@/client";
import { useToast } from "@/hooks/useToast";
import { MonacoEditor } from "@/components/MonacoEditor";
import styles from "./IntegrationSettings.module.scss";

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
      showSuccess("Успешно", "Страница успешно добавлена!");
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
      showSuccess("Успешно", "PR успешно проанализирован!");
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
    <div className={styles.container}>
      <div className={styles.header}>
        <div className={styles.title}>Настройки интеграций</div>
        <div className={styles.selectWrapper}>
          <Select
            value={selectedIntegration}
            onUpdate={setSelectedIntegration}
            options={integrationOptions}
            placeholder="Выберите интеграцию"
            width={340}
          />
        </div>
      </div>

      {isYandexWiki ? (
        <div className={styles.wikiSection}>
          <div className={styles.sectionTitle}>Добавьте новую страницу</div>
          <div className={styles.inputGroup}>
            <div className={styles.inputWrapper}>
              <div className={styles.inputContent}>
                <input
                  type="text"
                  className={styles.selectText}
                  placeholder="Введите ссылку на страницу"
                  value={pageUrl}
                  onChange={(e) => setPageUrl(e.target.value)}
                />
              </div>
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
          </div>

          <div className={styles.indexStatus}>
            <div className={styles.sectionTitle}>Состояние индекса</div>
            <div className={styles.statusInfo}>
              <div className={styles.statusText}>{indexedPagesCount} страниц проиндексировано</div>
              <Button
                view="flat"
                size="s"
                onClick={handleUpdateIndex}
                disabled={isUpdating}
                loading={isUpdating}
              >
                {isUpdating ? "Обновляем..." : "обновить"}
              </Button>
            </div>
          </div>
        </div>
      ) : (
        <div className={styles.pullRequestSection}>
          <div className={styles.sectionTitle}>Внести изменения на основе кода Pull-Request-а</div>
          <div className={styles.inputGroup}>
            <div className={styles.inputWrapper}>
              <div className={styles.inputContent}>
                <input
                  type="text"
                  className={styles.selectText}
                  placeholder="Введите ссылку на PR"
                  value={prUrl}
                  onChange={(e) => setPrUrl(e.target.value)}
                />
              </div>
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
          </div>
        </div>
      )}

      <div className={styles.parametersSection}>
        <div className={styles.sectionLabel}>Параметры</div>
        <div className={styles.editorWrapper}>
          <MonacoEditor
            value={configContent}
            onChange={(value) => setConfigContent(value || "")}
            language="json"
            height="400px"
            theme="light"
            readOnly={false}
          />
        </div>
        <div className={styles.placeholder}></div>
      </div>

      <div className={styles.actionsSection}>
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
      </div>
    </div>
  );
}
