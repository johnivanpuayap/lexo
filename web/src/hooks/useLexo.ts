import { useState, useEffect, useCallback } from "react";
import type { LexoCallbacks, LexoError } from "../types/lexo";

export function useLexo() {
  const [ready, setReady] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function loadWasm() {
      try {
        const goWasm = new window.Go();
        const result = await WebAssembly.instantiateStreaming(
          fetch("/lexo.wasm"),
          goWasm.importObject
        );
        goWasm.run(result.instance);
        setReady(true);
      } catch (err) {
        console.error("Failed to load Lexo WASM:", err);
      } finally {
        setLoading(false);
      }
    }
    loadWasm();
  }, []);

  const run = useCallback(
    (source: string, callbacks: LexoCallbacks) => {
      if (!ready || !window.lexo) return;
      window.lexo.run(source, callbacks);
    },
    [ready]
  );

  const check = useCallback(
    (source: string): LexoError | null => {
      if (!ready || !window.lexo) return null;
      const result = window.lexo.check(source);
      if (result) {
        try {
          return JSON.parse(result) as LexoError;
        } catch {
          return { kind: "type", message: result, line: 0 };
        }
      }
      return null;
    },
    [ready]
  );

  const provideInput = useCallback(
    (value: string) => {
      if (!ready || !window.lexo) return;
      window.lexo.provideInput(value);
    },
    [ready]
  );

  return { ready, loading, run, check, provideInput };
}
