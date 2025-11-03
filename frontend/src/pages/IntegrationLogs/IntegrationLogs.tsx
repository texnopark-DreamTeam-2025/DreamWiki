import { useEffect, useState, useCallback, useRef } from "react";
import { Select, Card, Flex, Loader, Text } from "@gravity-ui/uikit";
import { integrationLogsGet, type IntegrationLogField, type IntegrationId } from "@/client";

import { MonacoEditor } from "@/components/MonacoEditor";
import { useToast } from "@/hooks/useToast";

const INTEGRATION_OPTIONS = [
  { value: "ywiki", content: "Yandex Wiki" },
  { value: "github", content: "GitHub" },
];

export default function IntegrationLogs() {
  const { showError } = useToast();
  const [selectedIntegration, setSelectedIntegration] = useState<IntegrationId>("ywiki");
  const [logs, setLogs] = useState<IntegrationLogField[]>([]);
  const [loading, setLoading] = useState(false);
  const [hasMore, setHasMore] = useState(false);
  const editorRef = useRef<any>(null);
  const isLoadingMore = useRef(false);
  const cursorRef = useRef<string>("");
  const lastProcessedCursor = useRef<string>("");

  // Форматирование логов для отображения в Monaco Editor
  const formatLogs = useCallback((logFields: IntegrationLogField[]): string => {
    return logFields
      .map((log) => {
        const timestamp = log.created_at ? new Date(log.created_at).toLocaleString("ru-RU") : "";
        return `[${timestamp}] ${log.content}`;
      })
      .join("\n");
  }, []);

  // Загрузка первой порции логов
  const loadInitialLogs = useCallback(
    async (integrationId: IntegrationId) => {
      setLoading(true);
      setLogs([]);
      cursorRef.current = ""; // Сбрасываем ref
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
          cursorRef.current = res.data.next_info.cursor;
          setHasMore(res.data.next_info.has_more);
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
    if (loading || isLoadingMore.current || !hasMore || !cursorRef.current) return;

    // Проверяем, не обрабатывали ли уже этот курсор
    if (cursorRef.current === lastProcessedCursor.current) {
      return;
    }

    isLoadingMore.current = true;
    const currentCursor = cursorRef.current;

    try {
      const res = await integrationLogsGet({
        body: {
          integration_id: selectedIntegration as IntegrationId,
          cursor: currentCursor,
        },
      });

      if (res.error) {
        console.error("Ошибка загрузки дополнительных логов:", res.error);
        return;
      }

      if (res.data) {
        // Отмечаем этот курсор как обработанный
        lastProcessedCursor.current = currentCursor;

        // Если получили пустой массив логов или пустой курсор - данных больше нет
        if (!res.data.next_info.has_more) {
          setHasMore(false);
          return; // Выходим, не обновляя курсор
        }

        // Добавляем новые логи и обновляем курсор
        setLogs((prev) => [...prev, ...res.data!.log_fields]);
        cursorRef.current = res.data.next_info.cursor;

        // Проверяем, есть ли еще данные
        setHasMore(true); // Если дошли до этого места, значит есть данные и новый курсор
      } else {
        // Если нет data в ответе - данных больше нет
        lastProcessedCursor.current = currentCursor;
        setHasMore(false);
      }
    } catch (error) {
      console.error("Ошибка загрузки дополнительных логов:", error);
    } finally {
      isLoadingMore.current = false;
    }
  }, [loading, hasMore, selectedIntegration]);

  // Обработчик скролла в Monaco Editor
  const handleEditorScroll = useCallback(
    (editor: any) => {
      // Строгие проверки для предотвращения лишних запросов
      if (!editor || !hasMore || loading || isLoadingMore.current || !cursorRef.current) {
        return;
      }

      const currentCursor = cursorRef.current;

      // Проверяем, что этот курсор еще не обрабатывался
      if (currentCursor === lastProcessedCursor.current) {
        return;
      }

      const scrollTop = editor.getScrollTop();
      const scrollHeight = editor.getScrollHeight();
      const clientHeight = editor.getLayoutInfo().height;

      // Загружаем больше данных когда приближаемся к концу (90% прокрутки для менее частых вызовов)
      if (scrollTop + clientHeight >= scrollHeight * 0.9) {
        // Дополнительная проверка перед вызовом
        if (
          hasMore &&
          !loading &&
          !isLoadingMore.current &&
          currentCursor &&
          currentCursor !== lastProcessedCursor.current
        ) {
          loadMoreLogs();
        }
      }
    },
    [hasMore, loadMoreLogs, loading]
  );

  // Начальная загрузка логов
  useEffect(() => {
    loadInitialLogs(selectedIntegration);
  }, [selectedIntegration]); // Убираем loadInitialLogs из зависимостей

  // Обработчик изменения типа интеграции
  const handleIntegrationChange = useCallback((values: string[]) => {
    const newIntegration = values[0] as IntegrationId;
    setSelectedIntegration(newIntegration);
  }, []);

  // Настройка обработчика скролла для Monaco Editor
  const handleEditorMount = useCallback(
    (editor: any) => {
      editorRef.current = editor;
      let scrollTimeout: number;

      // Добавляем обработчик скролла с debouncing
      editor.onDidScrollChange(() => {
        // Очищаем предыдущий timeout
        if (scrollTimeout) {
          clearTimeout(scrollTimeout);
        }

        // Устанавливаем новый timeout для предотвращения частых вызовов
        scrollTimeout = setTimeout(() => {
          handleEditorScroll(editor);
        }, 100); // 100ms debounce
      });
    },
    [handleEditorScroll]
  );

  const editorContent = formatLogs(logs);
  const isInitialLoading = loading && logs.length === 0;

  return (
    <Flex direction="column" grow={1} style={{ padding: 20, height: '100vh' }}>
      <Flex justifyContent="space-between" alignItems="center" style={{ marginBottom: 20, paddingBottom: 16, borderBottom: '1px solid var(--g-color-line-generic)' }}>
        <Text variant="display-1">Журнал интеграций</Text>
        <Flex maxWidth="s" >
          <Select
            size="l"
            placeholder="Выберите интеграцию"
            value={[selectedIntegration]}
            onUpdate={handleIntegrationChange}
            options={INTEGRATION_OPTIONS}
            width="max"
          />
        </Flex>
      </Flex>

      <Card style={{ flexGrow: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden', backgroundColor: 'var(--g-color-base-background)', border: '1px solid var(--g-color-line-generic)', borderRadius: 8 }}>
        {isInitialLoading ? (
          <Flex grow={1} alignItems="center" justifyContent="center" style={{ padding: 40 }}>
            <Text variant="body-1" color="secondary">
              Загрузка логов...
            </Text>
            <Loader size="m" />
          </Flex>
        ) : (
          <Flex direction="column" grow={1} style={{ position: 'relative', overflow: 'hidden' }}>
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
              <div style={{
                position: 'absolute',
                bottom: 10,
                right: 10,
                backgroundColor: 'var(--g-color-base-background)',
                border: '1px solid var(--g-color-line-generic)',
                borderRadius: 4,
                padding: '8px 12px',
                fontSize: 12,
                color: 'var(--g-color-text-secondary)',
                boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
                zIndex: 10
              }}>
                Загрузка дополнительных логов...
              </div>
            )}
          </Flex>
        )}
      </Card>

      {!hasMore && logs.length > 0 && (
        <Flex justifyContent="center" style={{ marginTop: 16, padding: 16, backgroundColor: 'var(--g-color-base-float)', borderRadius: 8, border: '1px solid var(--g-color-line-generic)' }}>
          <Text variant="body-2" color="secondary">Все логи загружены</Text>
        </Flex>
      )}
    </Flex>
  );
}
