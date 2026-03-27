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
