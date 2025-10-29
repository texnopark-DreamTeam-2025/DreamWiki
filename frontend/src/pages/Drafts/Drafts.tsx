import { useState, useEffect } from "react";
import { listDrafts } from "@/client";
import { useToast } from "@/hooks/useToast";
import { Button, Loader, Table, Label, Flex, Text } from "@gravity-ui/uikit";
import type { DraftDigest, DraftStatus } from "@/client/types.gen";
import { useNavigate } from "react-router-dom";

export default function Drafts() {
  const { showError } = useToast();
  const navigate = useNavigate();
  const [drafts, setDrafts] = useState<DraftDigest[]>([]);
  const [loading, setLoading] = useState(true);
  const [cursor, setCursor] = useState<string | null>(null);

  const fetchDrafts = async () => {
    try {
      setLoading(true);
      const response = await listDrafts({
        body: {
          cursor: cursor || undefined,
        },
      });

      if (response.error) {
        console.error("Ошибка получения черновиков:", response.error);
        showError("Ошибка", "Не удалось загрузить список черновиков");
        return;
      }

      setDrafts(response.data.drafts);
      setCursor(response.data.cursor || null);
    } catch (error) {
      console.error("Ошибка получения черновиков:", error);
      showError("Ошибка", "Произошла ошибка при загрузке черновиков");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDrafts();
  }, []);

  const getStatusLabel = (status: DraftStatus) => {
    switch (status) {
      case "active":
        return <Label theme="success">Активный</Label>;
      case "merged":
        return <Label theme="info">Объединен</Label>;
      case "needs_rebase":
        return <Label theme="warning">Требует ребейза</Label>;
      default:
        return <Label theme="normal">Неизвестно</Label>;
    }
  };

  const handleCreateDraft = () => {
    navigate("/drafts/new");
  };

  const handleViewDraft = (draftId: string) => {
    navigate(`/drafts/${draftId}`);
  };

  return (
    <Flex direction="column" style={{ padding: 24, height: "100%" }}>
      <Flex justifyContent="space-between" alignItems="center" style={{ marginBottom: 24 }}>
        <Text variant="header-1">Черновики</Text>
        <Button view="action" size="m" onClick={handleCreateDraft}>
          Создать черновик
        </Button>
      </Flex>

      {loading ? (
        <Flex
          direction="column"
          alignItems="center"
          justifyContent="center"
          style={{ flex: 1, gap: 16 }}
        >
          <Loader size="m" />
          <Text>Загрузка черновиков...</Text>
        </Flex>
      ) : (
        <Flex direction="column" style={{ flex: 1 }}>
          {drafts.length === 0 ? (
            <Flex
              direction="column"
              alignItems="center"
              justifyContent="center"
              style={{ flex: 1, gap: 16, textAlign: "center" }}
            >
              <Text color="secondary">У вас пока нет черновиков</Text>
              <Button view="action" size="m" onClick={handleCreateDraft}>
                Создать первый черновик
              </Button>
            </Flex>
          ) : (
            <Table
              columns={[
                {
                  id: "title",
                  name: "Название",
                  template: (item: DraftDigest) => <Text variant="body-2">{item.draft_title}</Text>,
                },
                {
                  id: "page",
                  name: "Страница",
                  template: (item: DraftDigest) => (
                    <Text color="secondary" variant="body-2">
                      {item.page_digest.title}
                    </Text>
                  ),
                },
                {
                  id: "status",
                  name: "Статус",
                  template: (item: DraftDigest) => getStatusLabel(item.status),
                },
                {
                  id: "actions",
                  name: "Действия",
                  template: (item: DraftDigest) => (
                    <Button view="flat" size="s" onClick={() => handleViewDraft(item.draft_id)}>
                      Открыть
                    </Button>
                  ),
                },
              ]}
              data={drafts}
            />
          )}
        </Flex>
      )}
    </Flex>
  );
}
