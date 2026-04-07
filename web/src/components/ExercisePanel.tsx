import { useState, useEffect, useCallback } from "react";
import {
  categories,
  exercises,
  type Exercise,
  type Category,
} from "../data/exercises";
import "./ExercisePanel.css";

export interface ValidationResult {
  testCaseId: number;
  passed: boolean;
  expectedOutput: string[];
  actualOutput: string[];
  hint?: string;
}

export interface ExercisePanelProps {
  onLoadExercise: (code: string) => void;
  onValidate: (exercise: Exercise) => void;
  onClearResults: () => void;
  validationResults: ValidationResult[] | null;
  isRunning: boolean;
}

const PROGRESS_KEY = "lexo-exercise-progress";

type ProgressMap = Record<string, boolean>;

function loadProgress(): ProgressMap {
  try {
    const raw = localStorage.getItem(PROGRESS_KEY);
    if (raw) return JSON.parse(raw);
  } catch {
    // ignore
  }
  return {};
}

function saveProgress(progress: ProgressMap) {
  localStorage.setItem(PROGRESS_KEY, JSON.stringify(progress));
}

export function ExercisePanel({
  onLoadExercise,
  onValidate,
  onClearResults,
  validationResults,
  isRunning,
}: ExercisePanelProps) {
  const [selectedExercise, setSelectedExercise] = useState<Exercise | null>(
    null
  );

  // Clear validation results whenever the selected exercise changes
  useEffect(() => {
    onClearResults();
  }, [selectedExercise]);
  const [expandedCategories, setExpandedCategories] = useState<
    Record<string, boolean>
  >({});
  const [progress, setProgress] = useState<ProgressMap>(loadProgress);

  // Save progress whenever it changes
  useEffect(() => {
    saveProgress(progress);
  }, [progress]);

  // Mark exercise as complete when all tests pass
  useEffect(() => {
    if (!validationResults || !selectedExercise) return;
    const allPassed = validationResults.every((r) => r.passed);
    if (allPassed && !progress[selectedExercise.id]) {
      setProgress((prev) => ({ ...prev, [selectedExercise.id]: true }));
    }
  }, [validationResults, selectedExercise, progress]);

  const toggleCategory = useCallback((catId: string) => {
    setExpandedCategories((prev) => ({ ...prev, [catId]: !prev[catId] }));
  }, []);

  const getExercisesForCategory = useCallback(
    (catId: string) =>
      exercises
        .filter((e) => e.categoryId === catId)
        .sort((a, b) => a.order - b.order),
    []
  );

  const getCategoryProgress = useCallback(
    (catId: string) => {
      const catExercises = getExercisesForCategory(catId);
      const completed = catExercises.filter((e) => progress[e.id]).length;
      return { completed, total: catExercises.length };
    },
    [progress, getExercisesForCategory]
  );

  const navigateExercise = useCallback(
    (direction: -1 | 1) => {
      if (!selectedExercise) return;
      const catExercises = getExercisesForCategory(
        selectedExercise.categoryId
      );
      const idx = catExercises.findIndex(
        (e) => e.id === selectedExercise.id
      );
      const nextIdx = idx + direction;

      if (nextIdx >= 0 && nextIdx < catExercises.length) {
        setSelectedExercise(catExercises[nextIdx]);
      } else if (direction === 1) {
        // Move to next category
        const catIdx = categories.findIndex(
          (c) => c.id === selectedExercise.categoryId
        );
        if (catIdx < categories.length - 1) {
          const nextCat = categories[catIdx + 1];
          const nextExercises = getExercisesForCategory(nextCat.id);
          if (nextExercises.length > 0) {
            setSelectedExercise(nextExercises[0]);
          }
        }
      } else if (direction === -1) {
        // Move to previous category
        const catIdx = categories.findIndex(
          (c) => c.id === selectedExercise.categoryId
        );
        if (catIdx > 0) {
          const prevCat = categories[catIdx - 1];
          const prevExercises = getExercisesForCategory(prevCat.id);
          if (prevExercises.length > 0) {
            setSelectedExercise(prevExercises[prevExercises.length - 1]);
          }
        }
      }
    },
    [selectedExercise, getExercisesForCategory]
  );

  // Browser view
  if (!selectedExercise) {
    return (
      <div className="exercise-panel">
        <div className="exercise-browser">
          <div className="exercise-browser-header">Exercises</div>
          <div className="exercise-category-list">
            {categories
              .sort((a, b) => a.order - b.order)
              .map((cat) => (
                <CategoryCard
                  key={cat.id}
                  category={cat}
                  expanded={!!expandedCategories[cat.id]}
                  onToggle={() => toggleCategory(cat.id)}
                  progress={getCategoryProgress(cat.id)}
                  exercises={getExercisesForCategory(cat.id)}
                  completedMap={progress}
                  onSelectExercise={setSelectedExercise}
                />
              ))}
          </div>
        </div>
      </div>
    );
  }

  // Exercise view
  const allPassed =
    validationResults != null && validationResults.every((r) => r.passed);

  return (
    <div className="exercise-panel">
      <div className="exercise-detail">
        <button
          className="exercise-back-btn"
          onClick={() => setSelectedExercise(null)}
        >
          &larr; All Exercises
        </button>

        <div className="exercise-header">
          <h3 className="exercise-title">{selectedExercise.title}</h3>
          <div className="exercise-badges">
            <span
              className={`exercise-badge exercise-badge-${selectedExercise.difficulty.toLowerCase()}`}
            >
              {selectedExercise.difficulty}
            </span>
            <span className="exercise-badge exercise-badge-mode">
              {selectedExercise.mode === "syntax" ? "Syntax" : "Challenge"}
            </span>
            {progress[selectedExercise.id] && (
              <span className="exercise-badge exercise-badge-done">
                Completed
              </span>
            )}
          </div>
        </div>

        <p className="exercise-description">{selectedExercise.description}</p>

        <div className="exercise-actions">
          <button
            className="btn btn-load-exercise"
            onClick={() => onLoadExercise(selectedExercise.starterCode)}
          >
            Load Starter Code
          </button>
          <button
            className="btn btn-validate"
            onClick={() => onValidate(selectedExercise)}
            disabled={isRunning}
          >
            {isRunning ? "Running..." : "Check Solution"}
          </button>
        </div>

        <div className="exercise-test-cases">
          <div className="exercise-section-title">Test Cases</div>
          {selectedExercise.testCases.map((tc) => {
            const result = validationResults?.find(
              (r) => r.testCaseId === tc.id
            );
            return (
              <div
                key={tc.id}
                className={`exercise-test-case ${
                  result
                    ? result.passed
                      ? "test-passed"
                      : "test-failed"
                    : ""
                }`}
              >
                <div className="test-case-header">
                  <span className="test-case-label">Test {tc.id}</span>
                  {result && (
                    <span
                      className={`test-case-status ${
                        result.passed ? "status-pass" : "status-fail"
                      }`}
                    >
                      {result.passed ? "PASS" : "FAIL"}
                    </span>
                  )}
                </div>
                {tc.inputs && tc.inputs.length > 0 && (
                  <div className="test-case-row">
                    <span className="test-case-label-sm">Input:</span>
                    <code className="test-case-value">
                      {tc.inputs.join(", ")}
                    </code>
                  </div>
                )}
                <div className="test-case-row">
                  <span className="test-case-label-sm">Expected:</span>
                  <code className="test-case-value">
                    {tc.expectedOutput.join("\n")}
                  </code>
                </div>
                {result && !result.passed && (
                  <div className="test-case-row">
                    <span className="test-case-label-sm">Actual:</span>
                    <code className="test-case-value test-case-actual">
                      {result.actualOutput.length > 0
                        ? result.actualOutput.join("\n")
                        : "(no output)"}
                    </code>
                  </div>
                )}
              </div>
            );
          })}
        </div>

        {validationResults && !allPassed && (
          <div className="exercise-hint">
            <div className="exercise-section-title">Hint</div>
            <p>{selectedExercise.hint}</p>
          </div>
        )}

        {allPassed && (
          <div className="exercise-success">All tests passed!</div>
        )}

        <div className="exercise-nav">
          <button
            className="btn btn-nav"
            onClick={() => navigateExercise(-1)}
          >
            &larr; Previous
          </button>
          <button
            className="btn btn-nav"
            onClick={() => navigateExercise(1)}
          >
            Next &rarr;
          </button>
        </div>
      </div>
    </div>
  );
}

// ------- Sub-components -------

interface CategoryCardProps {
  category: Category;
  expanded: boolean;
  onToggle: () => void;
  progress: { completed: number; total: number };
  exercises: Exercise[];
  completedMap: ProgressMap;
  onSelectExercise: (e: Exercise) => void;
}

function CategoryCard({
  category,
  expanded,
  onToggle,
  progress: prog,
  exercises: catExercises,
  completedMap,
  onSelectExercise,
}: CategoryCardProps) {
  const pct = prog.total > 0 ? (prog.completed / prog.total) * 100 : 0;

  return (
    <div className="category-card">
      <button className="category-card-header" onClick={onToggle}>
        <div className="category-card-info">
          <span className="category-card-name">{category.name}</span>
          <span className="category-card-count">
            {prog.completed}/{prog.total}
          </span>
        </div>
        <div className="category-progress-bar">
          <div
            className="category-progress-fill"
            style={{ width: `${pct}%` }}
          />
        </div>
        <span className={`category-chevron ${expanded ? "expanded" : ""}`}>
          &#9660;
        </span>
      </button>
      {expanded && (
        <div className="category-exercises">
          {catExercises.map((ex) => (
            <button
              key={ex.id}
              className="exercise-row"
              onClick={() => onSelectExercise(ex)}
            >
              <span className="exercise-row-title">{ex.title}</span>
              <div className="exercise-row-meta">
                <span
                  className={`exercise-badge exercise-badge-${ex.difficulty.toLowerCase()}`}
                >
                  {ex.difficulty}
                </span>
                {completedMap[ex.id] && (
                  <span className="exercise-check">&#10003;</span>
                )}
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}
