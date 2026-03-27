interface TopBarProps {
  onRun: () => void;
  running: boolean;
  ready: boolean;
}

export function TopBar({ onRun, running, ready }: TopBarProps) {
  return (
    <div className="topbar">
      <div className="topbar-brand">
        <span className="topbar-logo">Lexo</span>
        <span className="topbar-file">untitled.lexo</span>
      </div>
      <div className="topbar-actions">
        <a
          className="btn btn-tutorial"
          href="/tutorial.html"
          target="_blank"
          rel="noopener noreferrer"
        >
          Tutorial
        </a>
        <button
          className="btn btn-run"
          onClick={onRun}
          disabled={running || !ready}
        >
          {running ? "Running..." : "▶ Run"}
        </button>
      </div>
    </div>
  );
}
