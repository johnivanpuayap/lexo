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

export const lexoTheme: editor.IStandaloneThemeData = {
  base: "vs-dark",
  inherit: true,
  rules: [
    { token: "keyword", foreground: "cba6f7", fontStyle: "bold" },
    { token: "type", foreground: "f9e2af" },
    { token: "string", foreground: "a6e3a1" },
    { token: "number", foreground: "fab387" },
    { token: "number.float", foreground: "fab387" },
    { token: "comment", foreground: "6c7086", fontStyle: "italic" },
    { token: "operator", foreground: "89dceb" },
    { token: "identifier", foreground: "cdd6f4" },
    { token: "delimiter", foreground: "9399b2" },
  ],
  colors: {
    "editor.background": "#1e1e2e",
    "editor.foreground": "#cdd6f4",
    "editor.lineHighlightBackground": "#313244",
    "editorCursor.foreground": "#f5e0dc",
    "editor.selectionBackground": "#45475a",
    "editorLineNumber.foreground": "#6c7086",
    "editorLineNumber.activeForeground": "#cdd6f4",
    "editorWidget.background": "#181825",
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
