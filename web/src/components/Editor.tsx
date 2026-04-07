import MonacoEditor, { OnMount } from "@monaco-editor/react";
import { useRef, useEffect } from "react";
import { lexoLanguageDef, lexoLightTheme, lexoDarkTheme } from "./lexo-lang";
import type { editor } from "monaco-editor";

interface EditorProps {
  value: string;
  onChange: (value: string) => void;
  errorLine?: number | null;
  theme: "light" | "dark";
}

export function Editor({ value, onChange, errorLine, theme }: EditorProps) {
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
  const monacoRef = useRef<typeof import("monaco-editor") | null>(null);
  const decorationsRef = useRef<editor.IEditorDecorationsCollection | null>(null);

  const handleMount: OnMount = (editor, monaco) => {
    editorRef.current = editor;
    monacoRef.current = monaco;

    monaco.languages.register({ id: "lexo" });
    monaco.languages.setMonarchTokensProvider("lexo", lexoLanguageDef);
    monaco.editor.defineTheme("lexo-light", lexoLightTheme);
    monaco.editor.defineTheme("lexo-dark", lexoDarkTheme);
    monaco.editor.setTheme(theme === "dark" ? "lexo-dark" : "lexo-light");
  };

  useEffect(() => {
    if (!monacoRef.current) return;
    monacoRef.current.editor.setTheme(theme === "dark" ? "lexo-dark" : "lexo-light");
  }, [theme]);

  useEffect(() => {
    if (!editorRef.current) return;
    const editor = editorRef.current;

    if (decorationsRef.current) {
      decorationsRef.current.clear();
    }

    if (errorLine && errorLine > 0) {
      decorationsRef.current = editor.createDecorationsCollection([
        {
          range: {
            startLineNumber: errorLine,
            startColumn: 1,
            endLineNumber: errorLine,
            endColumn: 1,
          },
          options: {
            isWholeLine: true,
            className: "error-line-decoration",
            glyphMarginClassName: "error-glyph",
          },
        },
      ]);
    }
  }, [errorLine]);

  return (
    <div className="editor-container">
      <MonacoEditor
        height="100%"
        defaultLanguage="lexo"
        theme={theme === "dark" ? "lexo-dark" : "lexo-light"}
        value={value}
        onChange={(v) => onChange(v ?? "")}
        onMount={handleMount}
        options={{
          fontSize: 14,
          fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
          minimap: { enabled: false },
          lineNumbers: "on",
          scrollBeyondLastLine: false,
          automaticLayout: true,
          tabSize: 4,
          insertSpaces: true,
          glyphMargin: true,
          folding: false,
          wordWrap: "off",
        }}
      />
    </div>
  );
}
