import MonacoEditor, { OnMount } from "@monaco-editor/react";
import { useRef, useEffect } from "react";
import { lexoLanguageDef, lexoTheme } from "./lexo-lang";
import type { editor } from "monaco-editor";

interface EditorProps {
  value: string;
  onChange: (value: string) => void;
  errorLine?: number | null;
}

export function Editor({ value, onChange, errorLine }: EditorProps) {
  const editorRef = useRef<editor.IStandaloneCodeEditor | null>(null);
  const decorationsRef = useRef<editor.IEditorDecorationsCollection | null>(null);

  const handleMount: OnMount = (editor, monaco) => {
    editorRef.current = editor;

    monaco.languages.register({ id: "lexo" });
    monaco.languages.setMonarchTokensProvider("lexo", lexoLanguageDef);
    monaco.editor.defineTheme("lexo-dark", lexoTheme);
    monaco.editor.setTheme("lexo-dark");
  };

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
        theme="lexo-dark"
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
