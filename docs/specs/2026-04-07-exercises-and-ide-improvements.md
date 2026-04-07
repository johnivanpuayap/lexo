# Lexo Exercise System + IDE Improvements — Design Spec

## Overview

Add a LeetCode-style exercise system to the Lexo IDE, fix the variable inspector, add a first-time onboarding overlay, and update the tutorial page to match the new claymorphic theme.

## 1. Exercise System

### 1.1 UI: Right Panel Tabs

The right panel currently shows only the Variable Inspector. It becomes a tabbed panel with two tabs:

- **Variables** — the existing variable inspector (once fixed)
- **Exercises** — exercise browser and problem view

A tab bar at the top of the right panel switches between them. The panel width increases from 280px to 360px to accommodate exercise descriptions.

### 1.2 Exercise Browser (list view)

When the Exercises tab is active and no exercise is selected, the user sees:

- **8 category cards**, each showing:
  - Category name
  - Progress bar (e.g., "2/4 completed")
  - Number of exercises
- Categories are displayed in recommended order (top to bottom) but none are locked
- Clicking a category expands it to show its exercises as a list
- Each exercise row shows: title, difficulty badge (Easy/Medium), completion checkmark

**Categories and exercises:**

#### Category 1: Variables & Types (4 exercises)
1. **Hello Variables** (Easy) — Declare a string variable `name` with value `"Lexo"` and print it. *Expected output: `Lexo`*
2. **Number Crunch** (Easy) — Declare an int `a = 10` and float `b = 3.14`, print both. *Expected output: `10` then `3.14`*
3. **Bool Logic** (Easy) — Declare a bool variable `isReady` set to `true` and print it. *Expected output: `true`*
4. **Swap Values** (Medium) — Given two int variables `x = 5` and `y = 10`, swap their values using a temp variable and print both. *Expected output: `10` then `5`*

#### Category 2: Operators & Expressions (4 exercises)
1. **Basic Math** (Easy) — Given `a: int = 15` and `b: int = 4`, print `a + b`, `a - b`, `a * b`, `a / b`, `a % b`. *Expected output: `19`, `11`, `60`, `3`, `3`*
2. **Comparison Check** (Easy) — Given `x: int = 7`, print whether `x > 5`, `x == 7`, `x != 10`. *Expected output: `true`, `true`, `true`*
3. **Logical Operators** (Easy) — Given `a: bool = true` and `b: bool = false`, print `a and b`, `a or b`, `not b`. *Expected output: `false`, `true`, `true`*
4. **Expression Builder** (Medium) — Calculate and print the result of `(10 + 5) * 2 - 3`. *Expected output: `27`*

#### Category 3: Control Flow (4 exercises)
1. **Age Checker** (Easy) — Given `age: int = 20`, print `"Adult"` if >= 18, `"Teen"` if >= 13, `"Child"` otherwise. *Expected output: `Adult`*
2. **Even or Odd** (Easy) — Given `n: int = 7`, print `"Even"` or `"Odd"`. *Expected output: `Odd`*
3. **Grade Calculator** (Medium) — Given `score: int = 85`, print letter grade: >= 90 is `"A"`, >= 80 is `"B"`, >= 70 is `"C"`, >= 60 is `"D"`, else `"F"`. *Expected output: `B`*
4. **Nested Conditions** (Medium) — Given `temp: int = 25` and `isRaining: bool = true`, print `"Stay inside"` if raining and temp < 15, `"Take umbrella"` if raining and temp >= 15, `"Go outside"` otherwise. *Expected output: `Take umbrella`*

#### Category 4: Loops (4 exercises)
1. **Count to Ten** (Easy) — Print numbers 1 through 10 using a while loop. *Expected output: `1` through `10`, one per line*
2. **Sum Array** (Easy) — Given `nums: int[] = [3, 7, 2, 8, 1]`, use a for loop to sum and print the total. *Expected output: `21`*
3. **Skip Evens** (Medium) — Print numbers 1 to 10 but skip even numbers using `continue`. *Expected output: `1`, `3`, `5`, `7`, `9`*
4. **Early Exit** (Medium) — Search for value `5` in `[1, 3, 5, 7, 9]`. Print `"Found at index X"` and `break`. *Expected output: `Found at index 2`*

#### Category 5: Functions (4 exercises)
1. **Say Hello** (Easy) — Write a function `greet(name: string) -> void` that prints `"Hello, "` + name. Call it with `"World"`. *Expected output: `Hello, World`*
2. **Add Two** (Easy) — Write a function `add(a: int, b: int) -> int` that returns the sum. Print the result of `add(3, 4)`. *Expected output: `7`*
3. **Max of Two** (Medium) — Write a function `max(a: int, b: int) -> int` that returns the larger value. Print `max(10, 25)`. *Expected output: `25`*
4. **Is Even Function** (Medium) — Write `isEven(n: int) -> bool` that returns true if n is even. Print `isEven(4)` and `isEven(7)`. *Expected output: `true` then `false`*

#### Category 6: Arrays (4 exercises)
1. **Print All** (Easy) — Given `fruits: string[] = ["apple", "banana", "cherry"]`, print each on a new line. *Expected output: `apple`, `banana`, `cherry`*
2. **Array Length** (Easy) — Given `nums: int[] = [10, 20, 30, 40, 50]`, print the length. *Expected output: `5`*
3. **Find Maximum** (Medium) — Given `nums: int[] = [3, 9, 1, 7, 5]`, find and print the largest value. *Expected output: `9`*
4. **Array Sum & Average** (Medium) — Given `scores: int[] = [80, 90, 70, 85, 95]`, print the sum and average (as int). *Expected output: `420` then `84`*

#### Category 7: Strings (4 exercises)
1. **Shout It** (Easy) — Given `msg: string = "hello"`, print it in uppercase. *Expected output: `HELLO`*
2. **String Length** (Easy) — Given `word: string = "Lexo"`, print its length. *Expected output: `4`*
3. **First Three** (Medium) — Given `text: string = "Programming"`, print the first 3 characters using substring. *Expected output: `Pro`*
4. **String Combo** (Medium) — Given `first: string = "Hello"` and `second: string = "World"`, concatenate with a space and print, then print the total length. *Expected output: `Hello World` then `11`*

#### Category 8: Challenges (8 exercises)
1. **FizzBuzz** (Easy) — Print numbers 1 to 20. For multiples of 3 print `"Fizz"`, multiples of 5 print `"Buzz"`, both print `"FizzBuzz"`. *(Multiple test cases with different ranges)*
2. **Reverse Array** (Easy) — Given an int array, print elements in reverse order. *(Test cases: [1,2,3,4,5] -> 5 4 3 2 1)*
3. **Palindrome Check** (Medium) — Check if a given string reads the same forwards and backwards. Print `true` or `false`. *(Test cases: "racecar" -> true, "hello" -> false)*
4. **Factorial** (Medium) — Calculate factorial of a given number. *(Test cases: 5 -> 120, 0 -> 1)*
5. **Count Vowels** (Medium) — Count vowels in a string and print the count. *(Test cases: "hello" -> 2, "aeiou" -> 5)*
6. **Sum of Digits** (Medium) — Sum all digits of an int. *(Test cases: 123 -> 6, 9999 -> 36)*
7. **Find Duplicates** (Medium) — Given an array, print any value that appears more than once. *(Test cases with arrays containing duplicates)*
8. **Two Sum** (Medium) — Given an array and target sum, find two numbers that add up to the target and print their indices. *(Test cases: [2,7,11,15] target 9 -> 0 1)*

### 1.3 Exercise View (problem selected)

When a user selects an exercise, the right panel shows:

- **Back arrow + "Back to List"** link at the top
- **Title** with difficulty badge
- **Mode badge**: "Syntax" (categories 1-7) or "Challenge" (category 8)
- **Description**: what the user needs to do
- **Starter code**: pre-loaded into the editor (replaces current code with confirmation if editor is dirty)
- **Expected output** section (for syntax exercises) or **Test Cases** section (for challenges)
- **Results area** (after running):
  - Syntax mode: pass/fail with line-by-line output comparison, hint on failure
  - Challenge mode: list of test cases with pass/fail status per case, hint on first failure
- **Navigation**: "Previous" / "Next Exercise" buttons

### 1.4 Validation Logic

**Syntax exercises (categories 1-7):**
- User clicks Run
- Program executes once with no input
- Actual output (captured lines) compared to expected output line by line
- On pass: green checkmark, exercise marked complete, confetti or subtle celebration
- On fail: red X, show expected vs actual for the first mismatched line, display hint text

**Challenge exercises (category 8):**
- User clicks Run
- Program executes once per test case
- For exercises requiring input: each test case provides input values via the existing `onRequestInput`/`provideInput` mechanism
- For exercises not requiring input (e.g., FizzBuzz with hardcoded range): single run, output compared
- Results displayed as a list: `Test 1: Pass`, `Test 2: Fail — expected "120", got "24"`
- Hint shown for the first failing test case
- Exercise marked complete only when all test cases pass

### 1.5 Exercise Data Format

Exercises stored as a TypeScript file (`web/src/data/exercises.ts`) exporting a typed array:

```typescript
interface Exercise {
  id: string;                    // e.g., "variables-1"
  category: string;              // e.g., "Variables & Types"
  categoryId: string;            // e.g., "variables"
  title: string;
  difficulty: "Easy" | "Medium";
  mode: "syntax" | "challenge";
  description: string;           // markdown-ish description
  starterCode: string;           // pre-loaded into editor
  hint: string;                  // shown on failure
  testCases: TestCase[];
}

interface TestCase {
  id: number;
  inputs?: string[];             // values fed to input() calls, in order
  expectedOutput: string[];      // expected output lines
}
```

Syntax exercises have a single test case with no inputs. Challenge exercises have 2-4 test cases, some with inputs.

### 1.6 Progress Persistence

Exercise completion stored in `localStorage` under key `lexo-exercise-progress`. Value is a JSON object mapping exercise IDs to completion status:

```json
{ "variables-1": true, "variables-2": true, "loops-3": false }
```

## 2. Fix Variable Inspector

Wire up the full data pipeline so the Variable Inspector actually shows runtime variables:

### 2.1 Go Interpreter

Add a public method to `Interpreter` that returns current environment variables as a serializable format:

```go
func (interp *Interpreter) GetVariables() []map[string]string {
    // Uses env.All() which already exists
    // Returns [{name, type, value}, ...]
}
```

### 2.2 WASM Bridge

Add an `onVariables` callback in `RegisterCallbacks()`. After program execution completes (and possibly at each statement during step-through), invoke this callback with JSON-serialized variable data.

### 2.3 Frontend Hook

Extend `useLexo.ts` to accept an `onVariables` callback in the run options. Parse the JSON and pass to React state.

### 2.4 App.tsx

Replace the static `useState<VariableInfo[]>([])` with a state that gets updated via the `onVariables` callback. Clear variables when starting a new run.

## 3. Onboarding Overlay

A first-time walkthrough shown when the user loads the IDE for the first time (checked via `localStorage` key `lexo-onboarding-seen`).

### Steps:
1. **Welcome** — "Welcome to Lexo! Let's take a quick tour." (centered overlay)
2. **Editor** — highlights the editor panel: "Write your Lexo code here"
3. **Run button** — highlights the Run button: "Click Run to execute your program"
4. **Output** — highlights the output console: "See your program's output here"
5. **Variables** — highlights the variable inspector tab: "Watch your variables update as your code runs"
6. **Exercises** — highlights the exercises tab: "Practice with guided exercises and challenges"
7. **Done** — "You're all set! Start coding or try an exercise."

Each step highlights the relevant UI area with a dimmed backdrop and a tooltip-style card with text + Next/Skip buttons. Minimal implementation — no animation library needed, just CSS positioning and z-index.

Dismissible at any step via "Skip tour" link. Stored in localStorage so it never shows again.

## 4. Tutorial Page Redesign

Update `web/public/tutorial.html` to match the new claymorphic theme:

- Replace Catppuccin Mocha colors with the cream/indigo light theme palette
- Add dark mode support (toggle button, CSS variables matching App.css)
- Use Fredoka for headings, Nunito for body text
- Card-based sections with claymorphic shadows (matching the IDE cards)
- Rounded corners (16px), thick borders
- Code blocks styled to match the Monaco light/dark editor themes
- Responsive layout for mobile
- Keep all existing tutorial content intact

## Technical Notes

- No routing library needed — exercise state managed via React state in App.tsx
- No backend — all exercises are static data, progress is localStorage
- Exercise validation runs in the existing WASM interpreter
- The right panel width should be adjustable or at least wider (360px) for exercise descriptions
- Challenge exercises that need input injection will reuse the existing `provideInput` mechanism programmatically (not via the UI input field)
