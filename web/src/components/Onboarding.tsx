import { useState, useEffect, useCallback, useRef } from "react";
import "./Onboarding.css";

interface OnboardingProps {
  onComplete: () => void;
}

interface Step {
  title: string;
  body: string;
  selector: string | null; // null = centered card (welcome/done)
}

const steps: Step[] = [
  {
    title: "Welcome to Lexo!",
    body: "A beginner-friendly language for learning static typing. Let\u2019s take a quick tour.",
    selector: null,
  },
  {
    title: "Editor",
    body: "Write your Lexo code here. The editor supports syntax highlighting and shows errors inline.",
    selector: ".editor-container",
  },
  {
    title: "Run Button",
    body: "Click Run to execute your program. The interpreter checks types before running!",
    selector: ".btn-run",
  },
  {
    title: "Output",
    body: "Your program\u2019s output appears here. Error messages will be shown in red with helpful hints.",
    selector: ".output-console",
  },
  {
    title: "Variables",
    body: "Watch your variables update in real-time as your code runs.",
    selector: ".right-panel",
  },
  {
    title: "Exercises",
    body: "Practice with guided exercises and coding challenges \u2014 from basics to FizzBuzz!",
    selector: ".right-panel",
  },
  {
    title: "You\u2019re all set!",
    body: "Start coding in the editor, or switch to Exercises to begin learning.",
    selector: null,
  },
];

export function Onboarding({ onComplete }: OnboardingProps) {
  const [currentStep, setCurrentStep] = useState(0);
  const [visible, setVisible] = useState(false);
  const [targetRect, setTargetRect] = useState<DOMRect | null>(null);
  const tooltipRef = useRef<HTMLDivElement>(null);

  const step = steps[currentStep];
  const isCentered = step.selector === null;
  const isLast = currentStep === steps.length - 1;

  // Measure the target element
  const measureTarget = useCallback(() => {
    if (step.selector) {
      const el = document.querySelector(step.selector);
      if (el) {
        setTargetRect(el.getBoundingClientRect());
      } else {
        setTargetRect(null);
      }
    } else {
      setTargetRect(null);
    }
  }, [step.selector]);

  // Fade in on mount and when step changes
  useEffect(() => {
    setVisible(false);
    measureTarget();
    const timer = setTimeout(() => setVisible(true), 20);
    return () => clearTimeout(timer);
  }, [currentStep, measureTarget]);

  // Re-measure on resize
  useEffect(() => {
    const handleResize = () => measureTarget();
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  }, [measureTarget]);

  const handleNext = () => {
    if (isLast) {
      onComplete();
    } else {
      setCurrentStep((s) => s + 1);
    }
  };

  const handleSkip = () => {
    onComplete();
  };

  // Compute tooltip position relative to spotlight
  const getTooltipStyle = (): React.CSSProperties => {
    if (!targetRect) return {};

    const padding = 12;
    const tooltipWidth = 340;
    const tooltipEstimatedHeight = 180;
    const gap = 16;

    // Prefer placing below the target
    let top = targetRect.bottom + gap;
    let left = targetRect.left + targetRect.width / 2 - tooltipWidth / 2;

    // If tooltip goes below viewport, place above
    if (top + tooltipEstimatedHeight > window.innerHeight - padding) {
      top = targetRect.top - gap - tooltipEstimatedHeight;
    }

    // If tooltip goes above viewport, place to the right
    if (top < padding) {
      top = targetRect.top + targetRect.height / 2 - tooltipEstimatedHeight / 2;
      left = targetRect.right + gap;
    }

    // Clamp horizontal position
    if (left < padding) left = padding;
    if (left + tooltipWidth > window.innerWidth - padding) {
      left = window.innerWidth - padding - tooltipWidth;
    }

    // Clamp vertical position
    if (top < padding) top = padding;

    return { top, left };
  };

  const spotlightPadding = 8;

  // Centered card for welcome/done
  if (isCentered) {
    return (
      <>
        <div className="onboarding-dim" />
        <div
          className={`onboarding-center-card ${visible ? "onboarding-visible" : ""}`}
        >
          <div className="onboarding-tooltip-title">{step.title}</div>
          <div className="onboarding-tooltip-body">{step.body}</div>
          <div className="onboarding-tooltip-footer">
            <span className="onboarding-step-indicator">
              {currentStep + 1} of {steps.length}
            </span>
            <div className="onboarding-actions">
              {!isLast && (
                <button className="onboarding-skip" onClick={handleSkip}>
                  Skip tour
                </button>
              )}
              <button className="onboarding-next" onClick={handleNext}>
                {isLast ? "Get Started" : "Next"}
              </button>
            </div>
          </div>
        </div>
      </>
    );
  }

  // Spotlight step
  return (
    <>
      <div className="onboarding-backdrop" onClick={handleNext} />
      {targetRect && (
        <div
          className="onboarding-spotlight"
          style={{
            top: targetRect.top - spotlightPadding,
            left: targetRect.left - spotlightPadding,
            width: targetRect.width + spotlightPadding * 2,
            height: targetRect.height + spotlightPadding * 2,
          }}
        />
      )}
      <div
        ref={tooltipRef}
        className={`onboarding-tooltip ${visible ? "onboarding-visible" : ""}`}
        style={getTooltipStyle()}
      >
        <div className="onboarding-tooltip-title">{step.title}</div>
        <div className="onboarding-tooltip-body">{step.body}</div>
        <div className="onboarding-tooltip-footer">
          <span className="onboarding-step-indicator">
            {currentStep + 1} of {steps.length}
          </span>
          <div className="onboarding-actions">
            <button className="onboarding-skip" onClick={handleSkip}>
              Skip tour
            </button>
            <button className="onboarding-next" onClick={handleNext}>
              Next
            </button>
          </div>
        </div>
      </div>
    </>
  );
}
