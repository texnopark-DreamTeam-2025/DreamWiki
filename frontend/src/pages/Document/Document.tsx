/**
 * Компонент для отображения документов
 *
 * Использует реальные API вызовы:
 * - getDiagnosticInfo для загрузки данных документа
 *
 * TODO: Заменить TREE_DATA на реальные данные дерева навигации из API
 */

import { useEffect, useState, useCallback } from "react";
import { useParams } from "react-router-dom";
import { TabProvider, TabList, Tab, TabPanel } from "@gravity-ui/uikit";
import { FileText, ChartColumn, Clock } from "@gravity-ui/icons";
import { getDiagnosticInfo, type V1DiagnosticInfoGetResponse } from "@/client";
import { TreeNavigation } from "@/components/TreeNavigation";
import { MonacoEditor } from "@/components/MonacoEditor";
import { useToast } from "@/hooks/useToast";
import styles from "./Document.module.scss";
import {
  TREE_DATA,
  INITIAL_SELECTED_NODE,
  INITIAL_EXPANDED_NODES,
} from "./treeData";

type TabId = "content" | "diagnostics" | "history" | "statistics";

export default function Document() {
  const { id } = useParams<{ id: string }>();
  const { showError } = useToast();
  const [page, setPage] = useState<V1DiagnosticInfoGetResponse | undefined>(
    undefined
  );
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<TabId>("content");
  const [selectedNode, setSelectedNode] = useState<string | null>(
    INITIAL_SELECTED_NODE
  );
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(
    INITIAL_EXPANDED_NODES
  );

  // Функция для переключения раскрытия узла
  const toggleNodeExpansion = useCallback(
    (nodeId: string, event: React.MouseEvent) => {
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
    },
    []
  );

  useEffect(() => {
    if (!id) {
      setLoading(false);
      return;
    }

    // Флаг для предотвращения обновления состояния после размонтирования
    let isCancelled = false;

    // Предотвращаем мерцание, сбрасываем предыдущие данные только при смене id
    setPage(undefined);
    setLoading(true);

    const fetchData = async () => {
      try {
        // API вызов для загрузки данных страницы
        const res = await getDiagnosticInfo({
          body: { page_id: id },
        });

        if (isCancelled) return; // Прерываем если компонент размонтирован

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
        if (isCancelled) return; // Прерываем если компонент размонтирован

        console.error("Ошибка загрузки данных:", error);
        showError("Ошибка загрузки", "Произошла ошибка при загрузке страницы");
      } finally {
        if (!isCancelled) {
          setLoading(false);
        }
      }
    };

    fetchData();

    // Cleanup функция для отмены запроса
    return () => {
      isCancelled = true;
    };
  }, [id]); // Убираем showError из зависимостей, чтобы предотвратить лишние запросы

  if (loading) {
    return (
      <div className={styles.container}>
        <div className={styles.sidebar}>
          <h3 className={styles.sidebarTitle}>База знаний</h3>
          <TreeNavigation
            data={TREE_DATA}
            selectedNode={selectedNode}
            expandedNodes={expandedNodes}
            onNodeSelect={setSelectedNode}
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
            Загрузка...
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
            data={TREE_DATA}
            selectedNode={selectedNode}
            expandedNodes={expandedNodes}
            onNodeSelect={setSelectedNode}
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
            Данные не найдены
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
          data={TREE_DATA}
          selectedNode={selectedNode}
          expandedNodes={expandedNodes}
          onNodeSelect={setSelectedNode}
          onNodeToggle={toggleNodeExpansion}
        />
      </div>

      {/* Правая панель с контентом */}
      <div className={styles.mainContent}>
        <TabProvider
          value={activeTab}
          onUpdate={(value) => setActiveTab(value as TabId)}
        >
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
                <div className={styles.contentHeader}>
                  {/* Здесь можно добавить кнопки */}
                </div>
                <h1 className={styles.contentTitle}>{page.page.title}</h1>
                <div className={styles.monacoContainer}>
                  <MonacoEditor
                    value={page.page.content || ""}
                    language="markdown"
                    height="100%"
                    readOnly={true}
                    theme="light"
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
