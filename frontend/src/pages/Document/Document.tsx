/**
 * Компонент для отображения документов
 *
 * Использует реальные API вызовы:
 * - getDiagnosticInfo для загрузки данных документа
 * - pagesTreeGet для загрузки дерева навигации
 */

import { useEffect, useState, useCallback } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { TabProvider, TabList, Tab, TabPanel } from "@gravity-ui/uikit";
import { FileText, ChartColumn, Clock } from "@gravity-ui/icons";
import {
  getDiagnosticInfo,
  pagesTreeGet,
  type V1DiagnosticInfoGetResponse,
  type TreeItem,
  type PageId,
} from "@/client";
import { TreeNavigation } from "@/components/TreeNavigation";
import type { TreeNode } from "@/components/TreeNavigation/types";
import { MonacoEditor } from "@/components/MonacoEditor";
import { useToast } from "@/hooks/useToast";
import styles from "./Document.module.scss";

type TabId = "content" | "diagnostics" | "history" | "statistics";

// Функция для преобразования TreeItem из API в TreeNode для компонента
const convertTreeItemToTreeNode = (item: TreeItem): TreeNode => {
  return {
    id: item.page_digest.page_id,
    title: item.page_digest.title,
    children: item.children?.map(convertTreeItemToTreeNode),
    expanded: item.expanded,
  };
};

// Функция для сбора всех expanded узлов из дерева
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
  const { showError } = useToast();

  // Состояние для данных документа
  const [page, setPage] = useState<V1DiagnosticInfoGetResponse | undefined>();
  const [loading, setLoading] = useState(true);

  // Состояние для дерева навигации
  const [treeData, setTreeData] = useState<TreeNode[]>([]);
  const [treeLoading, setTreeLoading] = useState(true);

  // Состояние для UI
  const [activeTab, setActiveTab] = useState<TabId>("content");
  const [selectedNode, setSelectedNode] = useState<string | null>(id || null);
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set());

  // Функция для загрузки дерева страниц
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

        // Устанавливаем expanded узлы на основе данных из API
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

  // Функция для переключения раскрытия узла
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

  // Обработчик монтирования Monaco Editor
  const handleEditorMount = useCallback((editor: any) => {
    // Настраиваем автоматическое изменение размера
    const resizeObserver = new ResizeObserver(() => {
      if (editor) {
        editor.layout();
      }
    });

    const editorDomNode = editor.getDomNode();
    if (editorDomNode?.parentElement) {
      resizeObserver.observe(editorDomNode.parentElement);
    }

    // Принудительно обновляем размер через небольшую задержку
    setTimeout(() => {
      if (editor) {
        editor.layout();
      }
    }, 100);

    // Очистка при размонтировании
    return () => {
      resizeObserver.disconnect();
    };
  }, []);

  // Функция для загрузки данных документа
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
  }, []); // Убираем showError из зависимостей

  // Эффект для загрузки дерева страниц при монтировании и изменении ID
  useEffect(() => {
    loadPagesTree();
  }, [id]); // Убираем loadPagesTree из зависимостей

  // Эффект для загрузки документа при изменении ID
  useEffect(() => {
    if (!id) {
      setLoading(false);
      return;
    }

    setSelectedNode(id);
    loadDocument(id);
  }, [id]); // Убираем loadDocument из зависимостей

  if (loading || treeLoading) {
    return (
      <div className={styles.container}>
        <div className={styles.sidebar}>
          <h3 className={styles.sidebarTitle}>База знаний</h3>
          {treeLoading ? (
            <div
              style={{
                padding: "20px",
                textAlign: "center",
                color: "var(--g-color-text-secondary)",
              }}
            >
              Загрузка дерева...
            </div>
          ) : (
            <TreeNavigation
              data={treeData}
              selectedNode={selectedNode}
              expandedNodes={expandedNodes}
              onNodeSelect={handleNodeSelect}
              onNodeToggle={toggleNodeExpansion}
            />
          )}
        </div>
        <div className={styles.mainContent}>
          <div
            style={{
              padding: "40px",
              textAlign: "center",
              color: "var(--g-color-text-secondary)",
            }}
          >
            {loading ? "Загрузка документа..." : "Загрузка..."}
          </div>
        </div>
      </div>
    );
  }

  if (!page) {
    return (
      <div className={styles.container}>
        <div className={styles.sidebar}>
          <h3 className={styles.sidebarTitle}>База знаний</h3>
          <TreeNavigation
            data={treeData}
            selectedNode={selectedNode}
            expandedNodes={expandedNodes}
            onNodeSelect={handleNodeSelect}
            onNodeToggle={toggleNodeExpansion}
          />
        </div>
        <div className={styles.mainContent}>
          <div
            style={{
              padding: "40px",
              textAlign: "center",
              color: "var(--g-color-text-secondary)",
            }}
          >
            Документ не найден
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Левая панель с деревом */}
      <div className={styles.sidebar}>
        <h3 className={styles.sidebarTitle}>База знаний</h3>
        <TreeNavigation
          data={treeData}
          selectedNode={selectedNode}
          expandedNodes={expandedNodes}
          onNodeSelect={handleNodeSelect}
          onNodeToggle={toggleNodeExpansion}
        />
      </div>

      {/* Правая панель с контентом */}
      <div className={styles.mainContent}>
        <TabProvider value={activeTab} onUpdate={(value) => setActiveTab(value as TabId)}>
          {/* Вкладки */}
          <div className={styles.tabsContainer}>
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
          </div>

          {/* Контент вкладок */}
          <div className={styles.tabContent}>
            <TabPanel value="content">
              <div className={styles.contentPanel}>
                <div className={styles.contentHeader}>{/* Здесь можно добавить кнопки */}</div>
                <h1 className={styles.contentTitle}>{page.page.title}</h1>
                <div className={styles.monacoContainer}>
                  <MonacoEditor
                    value={page.page.content || ""}
                    language="markdown"
                    height="100%"
                    readOnly={true}
                    theme="light"
                    onMount={handleEditorMount}
                  />
                </div>
              </div>
            </TabPanel>

            <TabPanel value="diagnostics">
              <div className={styles.contentPanel}>
                <div
                  style={{
                    padding: "40px",
                    textAlign: "center",
                    color: "var(--g-color-text-secondary)",
                  }}
                >
                  <h2>Диагностическая информация</h2>
                  <p>Раздел в разработке</p>
                </div>
              </div>
            </TabPanel>

            <TabPanel value="history">
              <div className={styles.contentPanel}>
                <div
                  style={{
                    padding: "40px",
                    textAlign: "center",
                    color: "var(--g-color-text-secondary)",
                  }}
                >
                  <h2>История страниц</h2>
                  <p>Раздел в разработке</p>
                </div>
              </div>
            </TabPanel>

            <TabPanel value="statistics">
              <div className={styles.contentPanel}>
                <div
                  style={{
                    padding: "40px",
                    textAlign: "center",
                    color: "var(--g-color-text-secondary)",
                  }}
                >
                  <h2>Статистика</h2>
                  <p>Раздел в разработке</p>
                </div>
              </div>
            </TabPanel>
          </div>
        </TabProvider>
      </div>
    </div>
  );
}
