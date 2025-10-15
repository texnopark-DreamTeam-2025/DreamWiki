/**
 * Страница журнала интеграций
 *
 * Отображает логи интеграций с бесконечным скроллом в Monaco Editor
 * Поддерживает выбор типа интеграции (ywiki, github)
 */

import { useEffect, useState, useCallback, useRef } from "react";
import { Select, Card, Button } from "@gravity-ui/uikit";
import {
  integrationLogsGet,
  type V1IntegrationLogsGetResponse,
  type IntegrationLogField,
  type IntegrationId,
} from "@/client";

import { MonacoEditor } from "@/components/MonacoEditor";
import { useToast } from "@/hooks/useToast";
import styles from "./IntegrationLogs.module.scss";

type IntegrationType = "ywiki" | "github";

const INTEGRATION_OPTIONS = [
  { value: "ywiki", content: "Yandex Wiki" },
  { value: "github", content: "GitHub" },
];

export default function IntegrationLogs() {
  const { showError } = useToast();
  const [selectedIntegration, setSelectedIntegration] =
    useState<IntegrationType>("ywiki");
  const [logs, setLogs] = useState<IntegrationLogField[]>([]);
  const [cursor, setCursor] = useState<string>("");
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const editorRef = useRef<any>(null);
  const isLoadingMore = useRef(false);

  // Форматирование логов для отображения в Monaco Editor
  const formatLogs = useCallback((logFields: IntegrationLogField[]): string => {
    return logFields
      .map((log) => {
        const timestamp = log.created_at
          ? new Date(log.created_at).toLocaleString("ru-RU")
          : "";
        return `[${timestamp}] ${log.content}`;
      })
      .join("\n");
  }, []);

  // Загрузка первой порции логов
  const loadInitialLogs = useCallback(
    async (integrationId: IntegrationType) => {
      setLoading(true);
      setLogs([]);
      setCursor("");
      setHasMore(false);

      try {
        const res = await integrationLogsGet({
          body: {
            integration_id: integrationId as IntegrationId,
            cursor: "",
          },
        });

        if (res.error) {
          console.error("Ошибка загрузки логов:", res.error);
          showError("Ошибка", "Не удалось загрузить логи интеграции");
          return;
        }

        if (res.data) {
          setLogs(res.data.log_fields);
          setCursor(res.data.cursor);
          const hasMoreData =
            res.data.log_fields.length > 0 && res.data.cursor !== "";
          setHasMore(hasMoreData);
        }
      } catch (error) {
        console.error("Ошибка загрузки логов:", error);
        showError("Ошибка", "Произошла ошибка при загрузке логов");
      } finally {
        setLoading(false);
      }
    },
    [showError]
  ); // Убираем loading из зависимостей  // Загрузка дополнительных логов (бесконечный скролл)
  const loadMoreLogs = useCallback(async () => {
    if (loading || isLoadingMore.current || !hasMore || !cursor) return;

    isLoadingMore.current = true;

    try {
      const res = await integrationLogsGet({
        body: {
          integration_id: selectedIntegration as IntegrationId,
          cursor: cursor,
        },
      });

      if (res.error) {
        console.error("Ошибка загрузки дополнительных логов:", res.error);
        return;
      }

      if (res.data) {
        const noMoreData =
          res.data.log_fields.length === 0 ||
          res.data.cursor === cursor ||
          res.data.cursor === "";

        if (noMoreData) {
          setHasMore(false);
        } else {
          setLogs((prev) => [...prev, ...res.data!.log_fields]);
          setCursor(res.data.cursor);
          setHasMore(true);
        }
      }
    } catch (error) {
      console.error("Ошибка загрузки дополнительных логов:", error);
    } finally {
      isLoadingMore.current = false;
    }
  }, [loading, hasMore, cursor, selectedIntegration]);

  // Обработчик скролла в Monaco Editor
  const handleEditorScroll = useCallback(
    (editor: any) => {
      if (!editor || !hasMore) return;

      const scrollTop = editor.getScrollTop();
      const scrollHeight = editor.getScrollHeight();
      const clientHeight = editor.getLayoutInfo().height;

      // Загружаем больше данных когда приближаемся к концу (80% прокрутки)
      if (scrollTop + clientHeight >= scrollHeight * 0.8) {
        loadMoreLogs();
      }
    },
    [hasMore, loadMoreLogs]
  );

  // Начальная загрузка логов
  useEffect(() => {
    loadInitialLogs(selectedIntegration);
  }, [selectedIntegration]); // Убираем loadInitialLogs из зависимостей

  // Обработчик изменения типа интеграции
  const handleIntegrationChange = useCallback((values: string[]) => {
    const newIntegration = values[0] as IntegrationType;
    setSelectedIntegration(newIntegration);
  }, []);

  // Настройка обработчика скролла для Monaco Editor
  const handleEditorMount = useCallback(
    (editor: any) => {
      editorRef.current = editor;

      // Добавляем обработчик скролла
      editor.onDidScrollChange(() => {
        handleEditorScroll(editor);
      });
    },
    [handleEditorScroll]
  );

  const editorContent = formatLogs(logs);
  const isInitialLoading = loading && logs.length === 0;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>Журнал интеграций</h1>
        <div className={styles.controls}>
          <Select
            size="l"
            placeholder="Выберите интеграцию"
            value={[selectedIntegration]}
            onUpdate={handleIntegrationChange}
            options={INTEGRATION_OPTIONS}
            className={styles.integrationSelect}
          />
        </div>
      </div>

      <Card className={styles.logsContainer}>
        {isInitialLoading ? (
          <div className={styles.loadingState}>
            <div className={styles.loadingText}>Загрузка логов...</div>
          </div>
        ) : (
          <div className={styles.editorContainer}>
            <MonacoEditor
              value={editorContent}
              language="text"
              height="100%"
              readOnly={true}
              theme="dark"
              onMount={handleEditorMount}
              onScroll={handleEditorScroll}
            />
            {isLoadingMore.current && (
              <div className={styles.loadingMore}>
                Загрузка дополнительных логов...
              </div>
            )}
          </div>
        )}
      </Card>

      {!hasMore && logs.length > 0 && (
        <div className={styles.endMessage}>Все логи загружены</div>
      )}
    </div>
  );
}
