import { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { TabProvider, TabList, Tab, TabPanel, Button, Text, Flex } from "@gravity-ui/uikit";
import { FileText, ChartColumn, Clock } from "@gravity-ui/icons";
import {
  getDiagnosticInfo,
  pagesTreeGet,
  indexatePage,
  type V1DiagnosticInfoGetResponse,
  type TreeItem,
  type PageId,
} from "@/client";
import { TreeNavigation } from "@/components/TreeNavigation";
import type { TreeNode } from "@/components/TreeNavigation/types";
import { MonacoEditor } from "@/components/MonacoEditor";
import { useToast } from "@/hooks/useToast";

type TabId = "content" | "diagnostics" | "history" | "statistics";

const convertTreeItemToTreeNode = (item: TreeItem): TreeNode => {
  return {
    id: item.page_digest.page_id,
    title: item.page_digest.title,
    children: item.children?.map(convertTreeItemToTreeNode),
    expanded: item.expanded,
  };
};

const collectExpandedNodes = (items: TreeNode[]): Set<string> => {
  const expandedSet = new Set<string>();

  const traverse = (nodes: TreeNode[]) => {
    nodes.forEach((node) => {
      if (node.expanded) {
        expandedSet.add(node.id);
      }
      if (node.children) {
        traverse(node.children);
      }
    });
  };

  traverse(items);
  return expandedSet;
};

export default function Document() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { showError, showSuccess } = useToast();

  const [page, setPage] = useState<V1DiagnosticInfoGetResponse | undefined>();
  const [loading, setLoading] = useState(true);

  const [treeData, setTreeData] = useState<TreeNode[]>([]);
  const [treeLoading, setTreeLoading] = useState(true);

  const [activeTab, setActiveTab] = useState<TabId>("content");
  const [selectedNode, setSelectedNode] = useState<string | null>(id || null);
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());
  const [reindexLoading, setReindexLoading] = useState(false);

  const loadPagesTree = useCallback(async (currentId?: string) => {
    setTreeLoading(true);

    try {
      const pageId = currentId || id;
      const activePageIds: PageId[] = pageId ? [pageId] : [];
      const res = await pagesTreeGet({
        body: { active_page_ids: activePageIds },
      });

      if (res.error) {
        console.error("Ошибка загрузки дерева:", res.error);
        showError("Ошибка", "Не удалось загрузить дерево навигации");
        return;
      }

      if (res.data?.tree) {
        const convertedTree = res.data.tree.map(convertTreeItemToTreeNode);
        setTreeData(convertedTree);

        const expanded = collectExpandedNodes(convertedTree);
        setExpandedNodes(expanded);
      }
    } catch (error) {
      console.error("Ошибка загрузки дерева:", error);
      showError("Ошибка", "Произошла ошибка при загрузке дерева навигации");
    } finally {
      setTreeLoading(false);
    }
  }, []); // Убираем все зависимости

  const toggleNodeExpansion = useCallback((nodeId: string, event: React.MouseEvent) => {
    event.stopPropagation();
    setExpandedNodes((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(nodeId)) {
        newSet.delete(nodeId);
      } else {
        newSet.add(nodeId);
      }
      return newSet;
    });

    // Примечание: В реальном приложении здесь можно было бы
    // отправить запрос на сервер для обновления состояния expanded узлов
    // Но для избежания бесконечных запросов пока оставляем только локальное состояние
  }, []);

  // Обработчик выбора узла дерева
  const handleNodeSelect = useCallback(
    (nodeId: string) => {
      setSelectedNode(nodeId);
      // Переходим к выбранному документу через React Router (без перезагрузки страницы)
      navigate(`/document/${nodeId}`);
    },
    [navigate]
  );

  const handleEditorMount = useCallback((editor: any) => {
    const resizeObserver = new ResizeObserver(() => {
      if (editor) {
        editor.layout();
      }
    });

    const editorDomNode = editor.getDomNode();
    if (editorDomNode?.parentElement) {
      resizeObserver.observe(editorDomNode.parentElement);
    }

    setTimeout(() => {
      if (editor) {
        editor.layout();
      }
    }, 100);

    return () => {
      resizeObserver.disconnect();
    };
  }, []);

  const handleReindex = async () => {
    try {
      setReindexLoading(true);
      await indexatePage({ body: { page_id: page!.page.page_id } });
      showSuccess("Успешно", "Переиндексация запущена");
    } catch (err) {
      console.error("Error triggering reindex:", err);
      showError("Ошибка", "Не удалось запустить переиндексацию");
    } finally {
      setReindexLoading(false);
    }
  };

  const loadDocument = useCallback(async (pageId: string) => {
    setLoading(true);
    setPage(undefined);

    try {
      const res = await getDiagnosticInfo({
        body: { page_id: pageId },
      });

      if (res.error) {
        console.error("Ошибка API:", res.error);
        showError("Ошибка загрузки", "Не удалось загрузить данные страницы");
        return;
      }

      if (res.data) {
        setPage(res.data);
      } else {
        showError("Ошибка", "Данные страницы не найдены");
      }
    } catch (error) {
      console.error("Ошибка загрузки данных:", error);
      showError("Ошибка загрузки", "Произошла ошибка при загрузке страницы");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadPagesTree();
  }, [id]);

  useEffect(() => {
    if (!id) {
      setLoading(false);
      return;
    }

    setSelectedNode(id);
    loadDocument(id);
  }, [id]);

  if (loading || treeLoading) {
    return (
      <Flex gap="4" direction="row" className="h-screen">
        <Flex
          gap="4"
          direction="column"
          className="p-4"
          style={{
            width: "300px",
            borderRight: "1px solid var(--g-color-line-generic)",
            overflowY: "auto",
          }}
        >
          <Text variant="header-2">База знаний</Text>
          {treeLoading ? (
            <Flex
              direction="row"
              justifyContent="center"
              alignItems="center"
              className="p-5"
            >
              <Text color="secondary">Загрузка дерева...</Text>
            </Flex>
          ) : (
            <TreeNavigation
              data={treeData}
              selectedNode={selectedNode}
              expandedNodes={expandedNodes}
              onNodeSelect={handleNodeSelect}
              onNodeToggle={toggleNodeExpansion}
            />
          )}
        </Flex>
        <Flex
          direction="row"
          justifyContent="center"
          alignItems="center"
          className="flex-1"
        >
          <Text color="secondary">{loading ? "Загрузка документа..." : "Загрузка..."}</Text>
        </Flex>
      </Flex>
    );
  }

  if (!page) {
    return (
      <Flex gap="4" direction="row" className="h-screen">
        <Flex
          gap="4"
          direction="column"
          className="p-4"
          style={{
            width: "300px",
            borderRight: "1px solid var(--g-color-line-generic)",
            overflowY: "auto",
          }}
        >
          <Text variant="header-2">База знаний</Text>
          <TreeNavigation
            data={treeData}
            selectedNode={selectedNode}
            expandedNodes={expandedNodes}
            onNodeSelect={handleNodeSelect}
            onNodeToggle={toggleNodeExpansion}
          />
        </Flex>
        <Flex
          direction="row"
          justifyContent="center"
          alignItems="center"
          className="flex-1"
        >
          <Text color="secondary">Документ не найден</Text>
        </Flex>
      </Flex>
    );
  }

  return (
    <Flex gap="4" direction="row" className="h-screen">
      {/* Левая панель с деревом */}
      <Flex
        gap="4"
        direction="column"
        className="p-4"
        style={{
          width: "300px",
          borderRight: "1px solid var(--g-color-line-generic)",
          overflowY: "auto",
        }}
      >
        <Text variant="header-2">База знаний</Text>
        <TreeNavigation
          data={treeData}
          selectedNode={selectedNode}
          expandedNodes={expandedNodes}
          onNodeSelect={handleNodeSelect}
          onNodeToggle={toggleNodeExpansion}
        />
      </Flex>

      {/* Правая панель с контентом */}
      <Flex direction="column" className="flex-1">
        <TabProvider value={activeTab} onUpdate={(value) => setActiveTab(value as TabId)}>
          {/* Вкладки */}
          <Flex direction="row" className="px-4">
            <TabList size="l">
              <Tab value="content" icon={<FileText />}>
                Содержимое
              </Tab>
              <Tab value="diagnostics" icon={<ChartColumn />}>
                Диагностическая информация
              </Tab>
              <Tab value="history" icon={<Clock />}>
                История страниц
              </Tab>
              <Tab value="statistics" icon={<ChartColumn />}>
                Статистика
              </Tab>
            </TabList>
          </Flex>

          {/* Контент вкладок */}
          <Flex direction="column" className="flex-1">
            <TabPanel value="content" className="h-full">
              <Flex
                direction="column"
                className="p-4 flex-1 h-full box-border"
              >
                <Flex direction="row" justifyContent="space-between" alignItems="center">
                  {/* Здесь можно добавить кнопки */}
                </Flex>
                <Text variant="header-1" className="mb-4">
                  {page.page.title}
                </Text>
                <Flex
                  direction="column"
                  style={{
                    border: "1px solid var(--g-color-line-generic)",
                    borderRadius: "4px",
                    overflow: "hidden",
                    backgroundColor: "var(--g-color-base-background)",
                    flex: 1,
                    minHeight: 0,
                  }}
                >
                  <MonacoEditor
                    value={page.page.content || ""}
                    language="markdown"
                    height="100%"
                    readOnly={true}
                    theme="light"
                    onMount={handleEditorMount}
                  />
                </Flex>
              </Flex>
            </TabPanel>

            <TabPanel value="diagnostics">
              <Flex direction="column" className="p-5 flex-1">
                <Flex direction="row" justifyContent="space-between" alignItems="center">
                  <Text variant="header-2">Диагностическая информация</Text>
                  <Button onClick={handleReindex} disabled={reindexLoading}>
                    {reindexLoading ? "Запуск..." : "Принудительно переиндексировать"}
                  </Button>
                </Flex>
                <Flex
                  direction="row"
                  justifyContent="center"
                  alignItems="center"
                  className="flex-1"
                >
                  <Text color="secondary">Раздел в разработке</Text>
                </Flex>
              </Flex>
            </TabPanel>

            <TabPanel value="history">
              <Flex direction="column" className="p-5 flex-1">
                <Flex
                  direction="row"
                  justifyContent="center"
                  alignItems="center"
                  className="flex-1"
                >
                  <Flex direction="column" alignItems="center">
                    <Text variant="header-2" className="mb-4">
                      История страниц
                    </Text>
                    <Text color="secondary">Раздел в разработке</Text>
                  </Flex>
                </Flex>
              </Flex>
            </TabPanel>

            <TabPanel value="statistics">
              <Flex direction="column" className="p-5 flex-1">
                <Flex
                  direction="row"
                  justifyContent="center"
                  alignItems="center"
                  className="flex-1"
                >
                  <Flex direction="column" alignItems="center">
                    <Text variant="header-2" className="mb-4">
                      Статистика
                    </Text>
                    <Text color="secondary">Раздел в разработке</Text>
                  </Flex>
                </Flex>
              </Flex>
            </TabPanel>
          </Flex>
        </TabProvider>
      </Flex>
    </Flex>
  );
}
