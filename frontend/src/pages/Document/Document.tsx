/**
 * TODO: УБРАТЬ MOCK ДАННЫЕ КОГДА БЭК ЗАРАБОТАЕТ
 *
 * Mock данные теперь находятся в ./mockData.ts
 * После интеграции с API нужно:
 * 1. Удалить импорт из mockData.ts
 * 2. Раскомментировать реальные API вызовы (getPageInfo, indexatePage)
 * 3. Заменить MOCK_* константы на реальные данные из API
 * 4. Удалить файл mockData.ts
 */

import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
// import { getPageInfo } from "./getDocument"; // TODO: раскомментировать когда бэк заработает
import { Button, TabProvider, TabList, Tab, TabPanel } from "@gravity-ui/uikit";
import { FileText, ChartColumn, Clock } from "@gravity-ui/icons";
// import { indexatePage, type V1DiagnosticInfoGetResponse } from "@/client"; // TODO: раскомментировать когда бэк заработает
import { type V1DiagnosticInfoGetResponse } from "@/client";
import { TreeNavigation } from "@/components/TreeNavigation";
import styles from "./Document.module.scss";
// MOCK DATA - удалить когда бэк заработает
import {
  MOCK_TREE_DATA,
  MOCK_INITIAL_SELECTED_NODE,
  MOCK_INITIAL_EXPANDED_NODES,
  mockFetchPageData,
  mockIndexPage,
} from "./mockData";

type TabId = "content" | "diagnostics" | "history" | "statistics";

export default function Document() {
  const { id } = useParams<{ id: string }>();
  const [page, setPage] = useState<V1DiagnosticInfoGetResponse | undefined>(
    undefined
  );
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<TabId>("content");
  const [selectedNode, setSelectedNode] = useState<string | null>(
    MOCK_INITIAL_SELECTED_NODE
  );
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(
    MOCK_INITIAL_EXPANDED_NODES
  );

  useEffect(() => {
    if (!id) return;

    // MOCK - симуляция загрузки
    const fetchData = async () => {
      setLoading(true);

      try {
        // TODO: Когда бэк заработает, заменить на:
        // const res = await getPageInfo(id);
        // if (res.data) {
        //   setPage(res.data);
        // }

        // MOCK - используем функцию из mockData
        const pageData = await mockFetchPageData(id);
        setPage(pageData);
      } catch (error) {
        console.error("Ошибка загрузки данных:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]); // Убрали mockData из зависимостей

  if (loading) return <div style={{ padding: "20px" }}>Загрузка...</div>;
  if (!page) return <div style={{ padding: "20px" }}>Данные не найдены</div>;

  // Функция для переключения раскрытия узла
  const toggleNodeExpansion = (nodeId: string, event: React.MouseEvent) => {
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
  };

  return (
    <div className={styles.container}>
      {/* Левая панель с деревом */}
      <div className={styles.sidebar}>
        <h3 className={styles.sidebarTitle}>База знаний</h3>
        <TreeNavigation
          data={MOCK_TREE_DATA}
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
                  <Button
                    view="action"
                    onClick={async () => {
                      try {
                        // TODO: Когда бэк заработает, заменить на:
                        // await indexatePage({ body: { page_id: id! } });

                        // MOCK - используем функцию из mockData
                        await mockIndexPage(id!);
                        alert("Страница успешно проиндексирована!");
                      } catch (error) {
                        console.error("Ошибка индексации:", error);
                        alert("Ошибка при индексации страницы");
                      }
                    }}
                  >
                    Проиндексировать
                  </Button>
                </div>
                <h1 className={styles.contentTitle}>{page.page.title}</h1>
                <div
                  className={styles.contentBody}
                  dangerouslySetInnerHTML={{
                    __html: page.page.content || "Содержимое отсутствует",
                  }}
                />
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
