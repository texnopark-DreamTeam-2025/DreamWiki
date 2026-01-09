import { useState, useEffect, useRef } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getDraft, updateDraft, deleteDraft, applyDraft } from "@/client";
import { useToast } from "@/hooks/useToast";
import { Button, Loader, Flex, Text, TextInput } from "@gravity-ui/uikit";
import { DiffEditor } from "@monaco-editor/react";
import type { Draft } from "@/client/types.gen";
import type { MonacoDiffEditor } from "@monaco-editor/react";

export default function DraftEditor() {
  const { draftId } = useParams<{ draftId: string }>();
  const navigate = useNavigate();
  const { showError, showSuccess } = useToast();
  const [draft, setDraft] = useState<Draft | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const diffEditorRef = useRef<MonacoDiffEditor | null>(null);

  const fetchDraft = async () => {
    if (!draftId) return;

    try {
      setLoading(true);
      const response = await getDraft({
        body: {
          draft_id: draftId,
        },
      });

      if (response.error) {
        console.error("Ошибка получения черновика:", response.error);
        showError("Ошибка", "Не удалось загрузить черновик");
        return;
      }

      setDraft(response.data.draft);
      setTitle(response.data.draft.draft_digest.draft_title);
      setContent(response.data.draft.content);
    } catch (error) {
      console.error("Ошибка получения черновика:", error);
      showError("Ошибка", "Произошла ошибка при загрузке черновика");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchDraft();
  }, [draftId]);

  const handleSave = async () => {
    if (!draftId) return;

    try {
      setSaving(true);
      const response = await updateDraft({
        body: {
          draft_id: draftId,
          new_content: content,
          new_title: title,
        },
      });

      if (response.error) {
        console.error("Ошибка сохранения черновика:", response.error);
        showError("Ошибка", "Не удалось сохранить черновик");
        return;
      }

      showSuccess("Успех", "Черновик успешно сохранен");
      // Refresh the draft to get updated data
      fetchDraft();
    } catch (error) {
      console.error("Ошибка сохранения черновика:", error);
      showError("Ошибка", "Произошла ошибка при сохранении черновика");
    } finally {
      setSaving(false);
    }
  };

  const handleMerge = async () => {
    if (!draftId) return;

    try {
      setSaving(true);
      const response = await applyDraft({
        body: {
          draft_id: draftId,
        },
      });

      if (response.error) {
        console.error("Ошибка применения черновика:", response.error);
        showError("Ошибка", "Не удалось применить черновик");
        return;
      }

      showSuccess("Успех", "Изменения успешно влиты");
      // Navigate back to drafts list
      navigate("/drafts");
    } catch (error) {
      console.error("Ошибка применения черновика:", error);
      showError("Ошибка", "Произошла ошибка при вливании изменений");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!draftId) return;

    try {
      setSaving(true);
      const response = await deleteDraft({
        body: {
          draft_id: draftId,
        },
      });

      if (response.error) {
        console.error("Ошибка удаления черновика:", response.error);
        showError("Ошибка", "Не удалось удалить черновик");
        return;
      }

      showSuccess("Успех", "Черновик успешно удален");
      // Navigate back to drafts list
      navigate("/drafts");
    } catch (error) {
      console.error("Ошибка удаления черновика:", error);
      showError("Ошибка", "Произошла ошибка при удалении черновика");
    } finally {
      setSaving(false);
    }
  };

  const handleEditorMount = (editor: MonacoDiffEditor) => {
    diffEditorRef.current = editor;

    // Add a listener to get content changes
    const modifiedModel = editor.getModel()?.modified;
    if (modifiedModel) {
      modifiedModel.onDidChangeContent(() => {
        const newContent = modifiedModel.getValue();
        setContent(newContent);
      });
    }
  };

  if (loading) {
    return (
      <Flex direction="column" alignItems="center" justifyContent="center" className="flex-1 gap-4">
        <Loader size="m" />
        <Text>Загрузка черновика...</Text>
      </Flex>
    );
  }

  if (!draft) {
    return (
      <Flex direction="column" alignItems="center" justifyContent="center" className="flex-1 gap-4">
        <Text>Черновик не найден</Text>
        <Button view="action" size="m" onClick={() => navigate("/drafts")}>
          Вернуться к списку черновиков
        </Button>
      </Flex>
    );
  }

  return (
    <Flex direction="column" className="p-6 h-full">
      <Flex justifyContent="space-between" alignItems="center" className="mb-6">
        <TextInput
          value={title}
          onUpdate={setTitle}
          placeholder="Введите заголовок"
          className="flex-1 mr-4"
        />
        <Flex gap={2}>
          <Button
            view="normal"
            size="m"
            onClick={handleSave}
            loading={saving}
            disabled={saving}
          >
            Сохранить
          </Button>
          <Button
            view="action"
            size="m"
            onClick={handleMerge}
            loading={saving}
            disabled={saving}
          >
            Влить изменения
          </Button>
          <Button
            view="outlined-danger"
            size="m"
            onClick={handleDelete}
            loading={saving}
            disabled={saving}
          >
            Удалить
          </Button>
        </Flex>
      </Flex>

      <div className="flex-1">
        <DiffEditor
          height="100%"
          language="markdown"
          original={draft.original_page_content || ""}
          modified={content}
          onMount={handleEditorMount}
          options={{
            readOnly: false,
            minimap: { enabled: false },
            scrollBeyondLastLine: false,
            wordWrap: "on",
            lineNumbers: "on",
            glyphMargin: false,
            folding: true,
            lineDecorationsWidth: 10,
            lineNumbersMinChars: 3,
            fontSize: 14,
            fontFamily:
              'JetBrains Mono, Monaco, Consolas, "Courier New", monospace',
            automaticLayout: true,
            contextmenu: true,
            selectOnLineNumbers: true,
            roundedSelection: false,
            cursorStyle: "line",
            cursorWidth: 2,
            renderWhitespace: "boundary",
            smoothScrolling: true,
            quickSuggestions: {
              other: true,
              comments: false,
              strings: false,
            },
            parameterHints: {
              enabled: true,
            },
            suggestOnTriggerCharacters: true,
            acceptSuggestionOnEnter: "on",
            tabCompletion: "on",
          }}
          loading={<div>Загрузка редактора...</div>}
        />
      </div>
    </Flex>
  );
}
