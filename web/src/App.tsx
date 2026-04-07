import { useState, useCallback, useEffect } from "react";
import { TopBar } from "./components/TopBar";
import { Editor } from "./components/Editor";
import { OutputConsole, OutputLine } from "./components/OutputConsole";
import { RightPanel } from "./components/RightPanel";
import { Onboarding } from "./components/Onboarding";
import { useLexo } from "./hooks/useLexo";
import { defaultProgram } from "./components/lexo-lang";
import type { Exercise } from "./data/exercises";
import type { ValidationResult } from "./components/ExercisePanel";
import type { VariableInfo } from "./components/VariableInspector";
import type { LexoError } from "./types/lexo";
import "./App.css";

const ONBOARDING_KEY = "lexo-onboarding-seen";

function App() {
  const { ready, loading, run, provideInput } = useLexo();
  const [source, setSource] = useState(defaultProgram);
  const [output, setOutput] = useState<OutputLine[]>([]);
  const [running, setRunning] = useState(false);
  const [waitingForInput, setWaitingForInput] = useState(false);
  const [errorLine, setErrorLine] = useState<number | null>(null);
  const [variables, setVariables] = useState<VariableInfo[]>([]);
  const [theme, setTheme] = useState<"light" | "dark">("light");
  const [activeTab, setActiveTab] = useState<"variables" | "exercises">("exercises");
  const [validationResults, setValidationResults] = useState<ValidationResult[] | null>(null);
  const [validating, setValidating] = useState(false);
  const [showOnboarding, setShowOnboarding] = useState(false);

  useEffect(() => {
    if (!localStorage.getItem(ONBOARDING_KEY)) {
      setShowOnboarding(true);
    }
  }, []);

  const handleOnboardingComplete = useCallback(() => {
    setShowOnboarding(false);
    localStorage.setItem(ONBOARDING_KEY, "true");
  }, []);

  const handleToggleTheme = useCallback(() => {
    setTheme((t) => (t === "light" ? "dark" : "light"));
  }, []);

  const handleRun = useCallback(() => {
    setOutput([]);
    setRunning(true);
    setErrorLine(null);
    setVariables([]);
    setValidationResults(null);
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
      onVariables: (variablesJson: string) => {
        try {
          const vars = JSON.parse(variablesJson) as VariableInfo[];
          setVariables(vars);
        } catch {
          // ignore parse errors
        }
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

  const handleLoadExercise = useCallback((code: string) => {
    setSource(code);
    setOutput([]);
    setErrorLine(null);
    setValidationResults(null);
    setVariables([]);
  }, []);

  const handleValidate = useCallback(
    (exercise: Exercise) => {
      if (!ready || !window.lexo) return;

      setValidating(true);
      setValidationResults(null);
      setOutput([]);
      setErrorLine(null);

      const results: ValidationResult[] = [];
      let currentTestIndex = 0;
      const testCases = exercise.testCases;

      const runTestCase = (testIndex: number) => {
        if (testIndex >= testCases.length) {
          setValidationResults(results);
          setValidating(false);
          return;
        }

        const tc = testCases[testIndex];
        const actualOutput: string[] = [];
        let inputIndex = 0;

        run(source, {
          onOutput: (text: string) => {
            actualOutput.push(text);
          },
          onError: () => {
            results.push({
              testCaseId: tc.id,
              passed: false,
              expectedOutput: tc.expectedOutput,
              actualOutput: ["(error during execution)"],
              hint: exercise.hint,
            });
            currentTestIndex++;
            runTestCase(currentTestIndex);
          },
          onComplete: () => {
            const passed =
              actualOutput.length === tc.expectedOutput.length &&
              actualOutput.every((line, i) => line === tc.expectedOutput[i]);

            results.push({
              testCaseId: tc.id,
              passed,
              expectedOutput: tc.expectedOutput,
              actualOutput,
              hint: passed ? undefined : exercise.hint,
            });

            currentTestIndex++;
            runTestCase(currentTestIndex);
          },
          onRequestInput: () => {
            if (tc.inputs && inputIndex < tc.inputs.length) {
              const value = tc.inputs[inputIndex];
              inputIndex++;
              // Small delay to allow the WASM bridge to be ready for input
              setTimeout(() => provideInput(value), 10);
            }
          },
          onVariables: () => {
            // ignore during validation
          },
        });
      };

      runTestCase(currentTestIndex);
    },
    [source, run, ready, provideInput]
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
      {showOnboarding && <Onboarding onComplete={handleOnboardingComplete} />}
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
          <RightPanel
            activeTab={activeTab}
            onTabChange={setActiveTab}
            variables={variables}
            exerciseProps={{
              onLoadExercise: handleLoadExercise,
              onValidate: handleValidate,
              onClearResults: () => setValidationResults(null),
              validationResults,
              isRunning: running || validating,
            }}
          />
        </div>
      </div>
    </div>
  );
}

export default App;
