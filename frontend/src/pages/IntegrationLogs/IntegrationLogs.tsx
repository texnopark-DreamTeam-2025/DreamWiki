import { useEffect, useState, useCallback } from "react";
import { Select, Card, Flex, Loader, Text, Box } from "@gravity-ui/uikit";
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
  const [hasMore, setHasMore] = useState(true);
  const [cursor, setCursor] = useState<string | undefined>(undefined);

  let loadingAtNow = false;

  const formatLogs = useCallback((logFields: IntegrationLogField[]): string => {
    return logFields
      .map((log) => {
        const timestamp = log.created_at ? new Date(log.created_at).toLocaleString() : "";
        return `[${timestamp}] ${log.content}`;
      })
      .join("\n");
  }, []);

  const loadLogs = async (integrationId: IntegrationId, currentCursor?: string) => {
    try {
      const res = await integrationLogsGet({
        body: {
          integration_id: integrationId,
          cursor: currentCursor,
        },
      });

      if (res.error) {
        console.error("Error loading logs:", res.error);
        showError("Error", "Failed to load integration logs");
        return;
      }

      if (res.data) {
        // If loading more logs, append to existing logs
        if (currentCursor) {
          setLogs((prev) => [...prev, ...res.data!.log_fields]);
        } else {
          setLogs(res.data.log_fields);
        }
        setCursor(res.data.next_info.cursor);
        setHasMore(res.data.next_info.has_more);
      }
    } catch (error) {
      console.error("Error loading logs:", error);
      showError("Error", "An error occurred while loading logs");
    } finally {
      setLoading(false);
    }
  };

  const loadMoreLogs = () => {
    if (!hasMore || loading || loadingAtNow) {
      return;
    }
    setLoading(true);
    loadingAtNow = true;
    loadLogs(selectedIntegration, cursor);
  };

  useEffect(() => {
    loadMoreLogs();
  });

  const handleIntegrationChange = useCallback((values: string[]) => {
    const newIntegration = values[0] as IntegrationId;
    setSelectedIntegration(newIntegration);
    setLogs([]);
    setLoading(false);
    setHasMore(true);
    setCursor(undefined);
  }, []);

  const handleEditorScroll = useCallback(
    (editor: any) => {
      if (!editor || !hasMore || loading) return;

      const scrollTop = editor.getScrollTop();
      const scrollHeight = editor.getScrollHeight();
      const clientHeight = editor.getLayoutInfo().height;

      // Load more data when approaching the end (95% scroll)
      if (scrollTop + clientHeight >= scrollHeight * 0.95) {
        loadMoreLogs();
      }
    },
    [hasMore, loadMoreLogs, loading]
  );

  const editorContent = formatLogs(logs);

  return (
    <Flex direction="column" gap="4" height="100vh">
      <Flex justifyContent="space-between" alignItems="center">
        <Text variant="header-1">Журнал интеграций</Text>
        <Box width="200px">
          <Select
            size="l"
            placeholder="Select integration"
            value={[selectedIntegration]}
            onUpdate={handleIntegrationChange}
            options={INTEGRATION_OPTIONS}
            width="max"
          />
        </Box>
      </Flex>

      <Flex grow={1} dir="column" overflow="hidden">
        {loading && logs.length === 0 ? (
          <Flex grow={1} alignItems="center" justifyContent="center" gap="2">
            <Loader size="m" />
            <Text variant="body-1" color="secondary">
              Loading logs...
            </Text>
          </Flex>
        ) : (
          <Flex direction="column" grow={1} overflow="hidden" position="relative">
            <MonacoEditor
              value={editorContent}
              language="text"
              height="100%"
              readOnly={true}
              onScroll={handleEditorScroll}
            />
          </Flex>
        )}
      </Flex>

      {!hasMore && logs.length > 0 && (
        <Flex justifyContent="center">
          <Text variant="body-2" color="secondary">
            Больше нет логов
          </Text>
        </Flex>
      )}
    </Flex>
  );
}
