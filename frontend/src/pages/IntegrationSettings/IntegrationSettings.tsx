import { useState } from "react";
import { Select, Button } from "@gravity-ui/uikit";
import styles from "./IntegrationSettings.module.scss";

export default function IntegrationSettings() {
  const [selectedIntegration, setSelectedIntegration] = useState<string[]>([
    "GitHub",
  ]);
  const [prUrl, setPrUrl] = useState<string>("");
  const [pageUrl, setPageUrl] = useState<string>("");
  const [indexedPagesCount, setIndexedPagesCount] = useState<number>(15);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [isAddingPage, setIsAddingPage] = useState(false);
  const [isUpdating, setIsUpdating] = useState(false);

  const integrationOptions = [
    { value: "GitHub", content: "GitHub" },
    { value: "YandexWiki", content: "Yandex Wiki" },
  ];

  const handleAddPage = async () => {
    if (!pageUrl.trim()) return;

    setIsAddingPage(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      console.log("Добавляем страницу:", pageUrl);
      setPageUrl("");
      // Увеличиваем счетчик проиндексированных страниц
      setIndexedPagesCount((prev) => prev + 1);
      alert("Страница успешно добавлена!");
    } catch (error) {
      console.error("Ошибка добавления страницы:", error);
    } finally {
      setIsAddingPage(false);
    }
  };

  const handleUpdateIndex = async () => {
    setIsUpdating(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 3000));
      console.log("Обновляем индекс");
      alert("Индекс успешно обновлен!");
    } catch (error) {
      console.error("Ошибка обновления индекса:", error);
    } finally {
      setIsUpdating(false);
    }
  };

  const handleAnalyzeLink = async () => {
    if (!prUrl.trim()) return;

    setIsAnalyzing(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 2000));
      console.log("Анализируем ссылку:", prUrl);
    } catch (error) {
      console.error("Ошибка анализа:", error);
    } finally {
      setIsAnalyzing(false);
    }
  };

  const handleCheckSettings = async () => {
    setIsLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      console.log("Проверяем настройки");
      alert("Настройки корректны!");
    } catch (error) {
      console.error("Ошибка проверки:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleApplySettings = async () => {
    setIsLoading(true);
    try {
      await new Promise((resolve) => setTimeout(resolve, 1500));
      console.log("Применяем настройки");
      alert("Настройки успешно применены!");
    } catch (error) {
      console.error("Ошибка применения:", error);
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
              <div className={styles.statusText}>
                {indexedPagesCount} страниц проиндексировано
              </div>
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
          <div className={styles.sectionTitle}>
            Внести изменения на основе кода Pull-Request-а
          </div>
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
        <div className={styles.textareaWrapper}>
          <textarea
            className={styles.textarea}
            defaultValue={`{
  "token": "test",
  "url": "http://localhost:8080"
}`}
            rows={15}
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
