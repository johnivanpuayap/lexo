export interface LexoCallbacks {
  onOutput: (text: string) => void;
  onError: (errorJson: string) => void;
  onComplete: (success: boolean) => void;
  onRequestInput?: (prompt: string) => void;
  onVariables?: (variablesJson: string) => void;
}

export interface LexoError {
  kind: "lexer" | "parser" | "type" | "runtime";
  message: string;
  line: number;
}

export interface LexoVariable {
  name: string;
  type: string;
  value: string;
}

export interface LexoBridge {
  run: (source: string, callbacks: LexoCallbacks) => void;
  check: (source: string) => string | null;
  provideInput: (value: string) => void;
}

declare global {
  interface Window {
    lexo: LexoBridge;
    Go: any;
  }
}
