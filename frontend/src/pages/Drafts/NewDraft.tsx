import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { createDraft, pagesTreeGet } from "@/client";
import { useToast } from "@/hooks/useToast";
import { Button, Loader, Flex, Text, Select, Box } from "@gravity-ui/uikit";
import type { TreeItem, PageDigest } from "@/client/types.gen";

export default function NewDraft() {
  const navigate = useNavigate();
  const { showError, showSuccess } = useToast();
  const [pages, setPages] = useState<PageDigest[]>([]);
  const [selectedPageId, setSelectedPageId] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);

  const fetchPages = async () => {
    try {
      setLoading(true);
      // Получаем дерево страниц с пустым массивом active_page_ids
      // Это даст нам все корневые страницы
      const response = await pagesTreeGet({
        body: {
          active_page_ids: [],
        },
      });

      if (response.error) {
        console.error("Ошибка получения дерева страниц:", response.error);
        showError("Ошибка", "Не удалось загрузить список страниц");
        return;
      }

      // Преобразуем дерево в плоский список страниц
      const flattenPages = (items: TreeItem[]): PageDigest[] => {
        let result: PageDigest[] = [];
        for (const item of items) {
          result.push(item.page_digest);
          if (item.children) {
            result = result.concat(flattenPages(item.children));
          }
        }
        return result;
      };

      const allPages = flattenPages(response.data.tree);
      setPages(allPages);

      // Если есть страницы, выбираем первую по умолчанию
      if (allPages.length > 0 && !selectedPageId) {
        setSelectedPageId(allPages[0].page_id);
      }
    } catch (error) {
      console.error("Ошибка получения дерева страниц:", error);
      showError("Ошибка", "Произошла ошибка при загрузке страниц");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPages();
  }, []);

  const handleCreate = async () => {
    if (!selectedPageId) {
      showError("Ошибка", "Пожалуйста, выберите страницу");
      return;
    }

    try {
      setCreating(true);
      const response = await createDraft({
        body: {
          page_url: selectedPageId, // В API ожидается page_url, но по факту передается page_id
        },
      });

      if (response.error) {
        console.error("Ошибка создания черновика:", response.error);
        showError("Ошибка", "Не удалось создать черновик");
        return;
      }

      showSuccess("Успех", "Черновик успешно создан");
      // Перенаправляем на страницу редактирования созданного черновика
      navigate(`/drafts/${response.data.draft_id}`);
    } catch (error) {
      console.error("Ошибка создания черновика:", error);
      showError("Ошибка", "Произошла ошибка при создании черновика");
    } finally {
      setCreating(false);
    }
  };

  const handleCancel = () => {
    navigate("/drafts");
  };

  if (loading) {
    return (
      <Box height="100%" display="flex" alignItems="center" justifyContent="center" gap={4}>
        <Loader size="m" />
        <Text>Загрузка страниц...</Text>
      </Box>
    );
  }

  return (
    <Box padding={6} height="100%">
      <Box marginBottom={6}>
        <Text variant="header-1">Создать черновик</Text>
      </Box>

      <Box flex={1} display="flex" flexDirection="column" gap={4}>
        <Text variant="body-2">Выберите страницу для создания черновика:</Text>

        <Select
          value={selectedPageId ? [selectedPageId] : []}
          onUpdate={(value) => setSelectedPageId(value[0] || null)}
          options={pages.map(page => ({
            value: page.page_id,
            content: page.title
          }))}
          placeholder="Выберите страницу..."
        />

        <Box marginTop={4}>
          <Button
            view="action"
            size="m"
            onClick={handleCreate}
            loading={creating}
            disabled={creating || !selectedPageId}
            marginRight={2}
          >
            Создать черновик
          </Button>
          <Button
            view="outlined"
            size="m"
            onClick={handleCancel}
            disabled={creating}
          >
            Отмена
          </Button>
        </Box>
      </Box>
    </Box>
  );
}
