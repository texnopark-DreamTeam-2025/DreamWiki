import { Editor } from "@monaco-editor/react";
import { useState, useEffect } from "react";
import styles from "./MonacoEditor.module.scss";

interface MonacoEditorProps {
  value: string;
  onChange?: (value: string | undefined) => void;
  language?: string;
  height?: string;
  readOnly?: boolean;
  theme?: "light" | "dark";
  onMount?: (editor: any) => void;
  onScroll?: (editor: any) => void;
}

export const MonacoEditor = ({
  value,
  onChange,
  language = "markdown",
  height = "400px",
  readOnly = false,
  theme = "light",
  onMount,
  onScroll,
}: MonacoEditorProps) => {
  const [isEditorReady, setIsEditorReady] = useState(false);

  useEffect(() => {
    if (isEditorReady) {
      import("monaco-editor").then((monaco) => {
        monaco.languages.typescript.javascriptDefaults.setDiagnosticsOptions({
          noSemanticValidation: true,
          noSyntaxValidation: false,
        });
      });
    }
  }, [isEditorReady]);

  const handleEditorDidMount = (editor: any) => {
    setIsEditorReady(true);

    if (onMount) {
      onMount(editor);
    }

    if (onScroll) {
      editor.onDidScrollChange(() => {
        onScroll(editor);
      });
    }
  };

  return (
    <div className={styles.editorContainer}>
      <Editor
        height={height}
        language={language}
        value={value}
        onChange={onChange}
        onMount={handleEditorDidMount}
        theme={theme === "dark" ? "vs-dark" : "vs"}
        options={{
          readOnly,
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          wordWrap: "on",
          lineNumbers: "on",
          glyphMargin: false,
          folding: true,
          lineDecorationsWidth: 10,
          lineNumbersMinChars: 3,
          fontSize: 14,
          fontFamily: 'JetBrains Mono, Monaco, Consolas, "Courier New", monospace',
          fontWeight: 'normal',
          automaticLayout: true,
          contextmenu: false,
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
          wordBasedSuggestions: "currentDocument",
        }}
        loading={<div className={styles.loading}>Загрузка редактора...</div>}
      />
    </div>
  );
};
