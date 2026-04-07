import type { languages, editor } from "monaco-editor";

export const lexoLanguageDef: languages.IMonarchLanguage = {
  keywords: [
    "def", "return", "if", "elif", "else", "while", "for", "in",
    "break", "continue", "and", "or", "not", "true", "false", "void",
  ],
  typeKeywords: ["int", "float", "string", "bool"],
  operators: [
    "+", "-", "*", "/", "%", "=", "==", "!=", "<", ">", "<=", ">=", "->",
  ],
  tokenizer: {
    root: [
      [/#.*$/, "comment"],
      [/"[^"]*"/, "string"],
      [/'[^']*'/, "string"],
      [/\d+\.\d+/, "number.float"],
      [/\d+/, "number"],
      [
        /[a-zA-Z_]\w*/,
        {
          cases: {
            "@keywords": "keyword",
            "@typeKeywords": "type",
            "@default": "identifier",
          },
        },
      ],
      [/->/, "operator"],
      [/[=!<>]=?/, "operator"],
      [/[+\-*/%]/, "operator"],
      [/[()[\],.:]+/, "delimiter"],
      [/\s+/, "white"],
    ],
  },
};

export const lexoLightTheme: editor.IStandaloneThemeData = {
  base: "vs",
  inherit: true,
  rules: [
    { token: "keyword", foreground: "6D28D9", fontStyle: "bold" },
    { token: "type", foreground: "B45309" },
    { token: "string", foreground: "047857" },
    { token: "number", foreground: "C2410C" },
    { token: "number.float", foreground: "C2410C" },
    { token: "comment", foreground: "9CA3AF", fontStyle: "italic" },
    { token: "operator", foreground: "0369A1" },
    { token: "identifier", foreground: "1E1B4B" },
    { token: "delimiter", foreground: "6B7280" },
  ],
  colors: {
    "editor.background": "#F8F5EE",
    "editor.foreground": "#1E1B4B",
    "editor.lineHighlightBackground": "#F0ECDF",
    "editorCursor.foreground": "#6366F1",
    "editor.selectionBackground": "#C7D2FE",
    "editorLineNumber.foreground": "#94A3B8",
    "editorLineNumber.activeForeground": "#1E1B4B",
    "editorWidget.background": "#EDE8DD",
    "editorGutter.background": "#F8F5EE",
  },
};

export const lexoDarkTheme: editor.IStandaloneThemeData = {
  base: "vs-dark",
  inherit: true,
  rules: [
    { token: "keyword", foreground: "A78BFA", fontStyle: "bold" },
    { token: "type", foreground: "FBBF24" },
    { token: "string", foreground: "6EE7B7" },
    { token: "number", foreground: "FB923C" },
    { token: "number.float", foreground: "FB923C" },
    { token: "comment", foreground: "6B7280", fontStyle: "italic" },
    { token: "operator", foreground: "67E8F9" },
    { token: "identifier", foreground: "E2E8F0" },
    { token: "delimiter", foreground: "94A3B8" },
  ],
  colors: {
    "editor.background": "#1C1B2E",
    "editor.foreground": "#E2E8F0",
    "editor.lineHighlightBackground": "#2D2A3E",
    "editorCursor.foreground": "#A78BFA",
    "editor.selectionBackground": "#3B3755",
    "editorLineNumber.foreground": "#6B7280",
    "editorLineNumber.activeForeground": "#E2E8F0",
    "editorWidget.background": "#171424",
    "editorGutter.background": "#1C1B2E",
  },
};

export const defaultProgram = `# Welcome to Lexo!
# A statically-typed language for learning

name: string = "World"
print("Hello, " + name)

# Try changing the code and clicking Run!
nums: int[] = [1, 2, 3, 4, 5]
total: int = 0
for n: int in nums:
    total = total + n
print(total)
`;
