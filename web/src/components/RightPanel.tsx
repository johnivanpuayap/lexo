import { VariableInspector, type VariableInfo } from "./VariableInspector";
import {
  ExercisePanel,
  type ExercisePanelProps,
} from "./ExercisePanel";
import "./RightPanel.css";

interface RightPanelProps {
  activeTab: "variables" | "exercises";
  onTabChange: (tab: "variables" | "exercises") => void;
  variables: VariableInfo[];
  exerciseProps: ExercisePanelProps;
}

export function RightPanel({
  activeTab,
  onTabChange,
  variables,
  exerciseProps,
}: RightPanelProps) {
  return (
    <div className="right-panel-container">
      <div className="right-panel-tabs">
        <button
          className={`right-panel-tab ${
            activeTab === "variables" ? "right-panel-tab-active" : ""
          }`}
          onClick={() => onTabChange("variables")}
        >
          Variables
        </button>
        <button
          className={`right-panel-tab ${
            activeTab === "exercises" ? "right-panel-tab-active" : ""
          }`}
          onClick={() => onTabChange("exercises")}
        >
          Exercises
        </button>
      </div>
      <div className="right-panel-content">
        {activeTab === "variables" ? (
          <VariableInspector variables={variables} />
        ) : (
          <ExercisePanel {...exerciseProps} />
        )}
      </div>
    </div>
  );
}
