import { useState, useCallback } from "react";
import { TopBar } from "./components/TopBar";
import { Editor } from "./components/Editor";
import { OutputConsole, OutputLine } from "./components/OutputConsole";
import { VariableInspector, VariableInfo } from "./components/VariableInspector";
import { useLexo } from "./hooks/useLexo";
import { defaultProgram } from "./components/lexo-lang";
import type { LexoError } from "./types/lexo";
import "./App.css";

function App() {
  const { ready, loading, run, provideInput } = useLexo();
  const [source, setSource] = useState(defaultProgram);
  const [output, setOutput] = useState<OutputLine[]>([]);
  const [running, setRunning] = useState(false);
  const [waitingForInput, setWaitingForInput] = useState(false);
  const [errorLine, setErrorLine] = useState<number | null>(null);
  const [variables] = useState<VariableInfo[]>([]);
  const [theme, setTheme] = useState<"light" | "dark">("light");

  const handleToggleTheme = useCallback(() => {
    setTheme((t) => (t === "light" ? "dark" : "light"));
  }, []);

  const handleRun = useCallback(() => {
    setOutput([]);
    setRunning(true);
    setErrorLine(null);
    setOutput((prev) => [
      ...prev,
      { text: "Running program...", type: "info" },
    ]);

    run(source, {
      onOutput: (text: string) => {
        setOutput((prev) => [...prev, { text, type: "output" }]);
      },
      onError: (errorJson: string) => {
        try {
          const err: LexoError = JSON.parse(errorJson);
          setOutput((prev) => [
            ...prev,
            { text: `${err.message}`, type: "error" },
          ]);
          if (err.line > 0) {
            setErrorLine(err.line);
          }
        } catch {
          setOutput((prev) => [...prev, { text: errorJson, type: "error" }]);
        }
      },
      onComplete: (success: boolean) => {
        setRunning(false);
        setWaitingForInput(false);
        if (success) {
          setOutput((prev) => [
            ...prev,
            { text: "Program finished.", type: "info" },
          ]);
        }
      },
      onRequestInput: (prompt: string) => {
        if (prompt) {
          setOutput((prev) => [...prev, { text: prompt, type: "input-prompt" }]);
        }
        setWaitingForInput(true);
      },
    });
  }, [source, run]);

  const handleSubmitInput = useCallback(
    (value: string) => {
      setWaitingForInput(false);
      setOutput((prev) => [...prev, { text: `> ${value}`, type: "output" }]);
      provideInput(value);
    },
    [provideInput]
  );

  if (loading) {
    return (
      <div className="loading" data-theme={theme}>
        <div className="loading-text">Loading Lexo...</div>
      </div>
    );
  }

  return (
    <div className="app" data-theme={theme}>
      <TopBar
        onRun={handleRun}
        running={running}
        ready={ready}
        theme={theme}
        onToggleTheme={handleToggleTheme}
      />
      <div className="main-layout">
        <div className="left-panel">
          <Editor value={source} onChange={setSource} errorLine={errorLine} theme={theme} />
          <OutputConsole
            lines={output}
            waitingForInput={waitingForInput}
            onSubmitInput={handleSubmitInput}
            onClear={() => setOutput([])}
          />
        </div>
        <div className="right-panel">
          <VariableInspector variables={variables} />
        </div>
      </div>
    </div>
  );
}

export default App;
