import { useRef, useEffect, useState } from "react";

export interface OutputLine {
  text: string;
  type: "output" | "error" | "info" | "input-prompt";
}

interface OutputConsoleProps {
  lines: OutputLine[];
  waitingForInput: boolean;
  onSubmitInput: (value: string) => void;
  onClear: () => void;
}

export function OutputConsole({
  lines,
  waitingForInput,
  onSubmitInput,
  onClear,
}: OutputConsoleProps) {
  const [inputValue, setInputValue] = useState("");
  const bottomRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [lines]);

  useEffect(() => {
    if (waitingForInput) {
      inputRef.current?.focus();
    }
  }, [waitingForInput]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmitInput(inputValue);
    setInputValue("");
  };

  return (
    <div className="output-console">
      <div className="output-header">
        <span>Output</span>
        <button className="btn-clear" onClick={onClear}>
          Clear
        </button>
      </div>
      <div className="output-lines">
        {lines.map((line, i) => (
          <div key={i} className={`output-line output-${line.type}`}>
            {line.text}
          </div>
        ))}
        {waitingForInput && (
          <form className="input-form" onSubmit={handleSubmit}>
            <input
              ref={inputRef}
              type="text"
              value={inputValue}
              onChange={(e) => setInputValue(e.target.value)}
              className="input-field"
              placeholder="Type your input and press Enter..."
            />
          </form>
        )}
        <div ref={bottomRef} />
      </div>
    </div>
  );
}
