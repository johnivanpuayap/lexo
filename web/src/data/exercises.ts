export interface TestCase {
  id: number;
  inputs?: string[];
  expectedOutput: string[];
}

export interface Exercise {
  id: string;
  category: string;
  categoryId: string;
  order: number;
  title: string;
  difficulty: "Easy" | "Medium" | "Hard";
  mode: "syntax" | "challenge";
  description: string;
  starterCode: string;
  hint: string;
  testCases: TestCase[];
}

export interface Category {
  id: string;
  name: string;
  order: number;
  description: string;
}

export const categories: Category[] = [
  {
    id: "variables",
    name: "Variables & Types",
    order: 1,
    description: "Learn how to declare variables with explicit types in Lexo.",
  },
  {
    id: "operators",
    name: "Operators & Expressions",
    order: 2,
    description: "Work with arithmetic, comparison, and logical operators.",
  },
  {
    id: "control-flow",
    name: "Control Flow",
    order: 3,
    description: "Use if, elif, and else to make decisions in your code.",
  },
  {
    id: "loops",
    name: "Loops",
    order: 4,
    description: "Repeat actions with while loops and for loops.",
  },
  {
    id: "functions",
    name: "Functions",
    order: 5,
    description: "Define reusable functions with typed parameters and return types.",
  },
  {
    id: "arrays",
    name: "Arrays",
    order: 6,
    description: "Store and manipulate collections of values.",
  },
  {
    id: "strings",
    name: "Strings",
    order: 7,
    description: "Work with text using string methods and concatenation.",
  },
  {
    id: "challenges",
    name: "Challenges",
    order: 8,
    description: "Put your skills together to solve classic beginner problems.",
  },
];

export const exercises: Exercise[] = [
  // ============================================================
  // Category 1: Variables & Types
  // ============================================================
  {
    id: "var-1",
    category: "Variables & Types",
    categoryId: "variables",
    order: 1,
    title: "Hello Variables",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Declare an integer variable called `age` set to `25`, and a string variable called `name` set to `\"Alice\"`. Print both variables.",
    starterCode: `# Declare your variables below
# Syntax: variableName: type = value

# Declare an int variable called age, set to 25

# Declare a string variable called name, set to "Alice"

# Print both variables
`,
    hint: "Use the syntax `age: int = 25` and `name: string = \"Alice\"`. Then use `print(age)` and `print(name)`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["25", "Alice"],
      },
    ],
  },
  {
    id: "var-2",
    category: "Variables & Types",
    categoryId: "variables",
    order: 2,
    title: "All the Types",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Declare one variable of each basic type: an `int` called `count` set to `10`, a `float` called `price` set to `9.99`, a `string` called `greeting` set to `\"hi\"`, and a `bool` called `active` set to `true`. Print all four.",
    starterCode: `# Declare one variable for each type and print them
# Types available: int, float, string, bool

# count should be 10

# price should be 9.99

# greeting should be "hi"

# active should be true

# Print all four variables
`,
    hint: "Remember: `float` uses decimal numbers like `9.99`, `bool` uses `true` or `false` (lowercase).",
    testCases: [
      {
        id: 1,
        expectedOutput: ["10", "9.99", "hi", "true"],
      },
    ],
  },
  {
    id: "var-3",
    category: "Variables & Types",
    categoryId: "variables",
    order: 3,
    title: "Reassigning Variables",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Declare an integer variable `x` set to `5`, print it, then reassign it to `10` and print it again. Remember: you only declare the type once.",
    starterCode: `# Declare x as an int with value 5
x: int = 5
print(x)

# Reassign x to 10 (no type needed for reassignment)

# Print x again
`,
    hint: "To reassign, just write `x = 10` without the type annotation. The type was already declared.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["5", "10"],
      },
    ],
  },
  {
    id: "var-4",
    category: "Variables & Types",
    categoryId: "variables",
    order: 4,
    title: "User Input",
    difficulty: "Medium",
    mode: "syntax",
    description:
      "Use `input()` to read a name from the user, then print a greeting. The `input()` function always returns a `string`. Print `\"Hello, \"` concatenated with the name.",
    starterCode: `# input() takes a prompt string and returns what the user types
# It always returns a string

# Read a name from the user with the prompt "Enter name: "

# Print "Hello, " + the name
`,
    hint: "Use `name: string = input(\"Enter name: \")` and then `print(\"Hello, \" + name)`.",
    testCases: [
      {
        id: 1,
        inputs: ["World"],
        expectedOutput: ["Hello, World"],
      },
      {
        id: 2,
        inputs: ["Lexo"],
        expectedOutput: ["Hello, Lexo"],
      },
    ],
  },

  // ============================================================
  // Category 2: Operators & Expressions
  // ============================================================
  {
    id: "op-1",
    category: "Operators & Expressions",
    categoryId: "operators",
    order: 1,
    title: "Basic Arithmetic",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Given two integers `a = 15` and `b = 4`, print the result of addition, subtraction, multiplication, integer division, and modulo (remainder) in that order.",
    starterCode: `a: int = 15
b: int = 4

# Print a + b

# Print a - b

# Print a * b

# Print a / b (integer division)

# Print a % b (remainder)
`,
    hint: "In Lexo, dividing two integers gives an integer result. `15 / 4` is `3`, not `3.75`. The `%` operator gives the remainder.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["19", "11", "60", "3", "3"],
      },
    ],
  },
  {
    id: "op-2",
    category: "Operators & Expressions",
    categoryId: "operators",
    order: 2,
    title: "Comparisons",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Given `x = 10` and `y = 20`, print the result of: `x < y`, `x > y`, `x == y`, `x != y`, `x <= 10`, `y >= 25`.",
    starterCode: `x: int = 10
y: int = 20

# Print the result of each comparison
# x < y

# x > y

# x == y

# x != y

# x <= 10

# y >= 25
`,
    hint: "Comparison operators return `bool` values: `true` or `false`. Just wrap each expression in `print()`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["true", "false", "false", "true", "true", "false"],
      },
    ],
  },
  {
    id: "op-3",
    category: "Operators & Expressions",
    categoryId: "operators",
    order: 3,
    title: "Logical Operators",
    difficulty: "Medium",
    mode: "syntax",
    description:
      "Lexo uses `and`, `or`, and `not` for logical operations. Given `a = true` and `b = false`, print: `a and b`, `a or b`, `not a`, `not b`.",
    starterCode: `a: bool = true
b: bool = false

# Print a and b

# Print a or b

# Print not a

# Print not b
`,
    hint: "Logical operators work on `bool` values. `and` returns `true` only if both sides are `true`. `or` returns `true` if either side is `true`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["false", "true", "false", "true"],
      },
    ],
  },
  {
    id: "op-4",
    category: "Operators & Expressions",
    categoryId: "operators",
    order: 4,
    title: "Expression Building",
    difficulty: "Medium",
    mode: "syntax",
    description:
      "Calculate and print the following:\n1. The average of `10`, `20`, and `30` using integer division: `(10 + 20 + 30) / 3`\n2. Whether `7` is odd: `7 % 2 != 0`\n3. A combined check: `5 > 3 and 10 < 20`",
    starterCode: `# 1. Calculate the average of 10, 20, 30 using integer division
#    and print it

# 2. Check if 7 is odd using modulo, and print the boolean result

# 3. Print whether both conditions are true: 5 > 3 and 10 < 20
`,
    hint: "For the average, use `(10 + 20 + 30) / 3`. For odd check, `7 % 2 != 0` gives a `bool`. You can combine comparisons with `and`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["20", "true", "true"],
      },
    ],
  },

  // ============================================================
  // Category 3: Control Flow
  // ============================================================
  {
    id: "cf-1",
    category: "Control Flow",
    categoryId: "control-flow",
    order: 1,
    title: "Simple If",
    difficulty: "Easy",
    mode: "syntax",
    description:
      'Given `temperature: int = 35`, write an if statement: if temperature is greater than 30, print `"Hot"`. Otherwise, print `"Cool"`.',
    starterCode: `temperature: int = 35

# Write an if/else statement
# If temperature > 30, print "Hot"
# Otherwise, print "Cool"
`,
    hint: 'Use `if temperature > 30:` followed by an indented `print("Hot")`, then `else:` with `print("Cool")`.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["Hot"],
      },
    ],
  },
  {
    id: "cf-2",
    category: "Control Flow",
    categoryId: "control-flow",
    order: 2,
    title: "Grade Checker",
    difficulty: "Easy",
    mode: "syntax",
    description:
      'Given `score: int = 75`, print the grade:\n- 90 or above: `"A"`\n- 80 or above: `"B"`\n- 70 or above: `"C"`\n- Below 70: `"F"`',
    starterCode: `score: int = 75

# Use if/elif/else to determine the grade
# 90+ -> "A", 80+ -> "B", 70+ -> "C", below 70 -> "F"
`,
    hint: "Use `if score >= 90:` then `elif score >= 80:` then `elif score >= 70:` then `else:`. Each branch should print the letter grade.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["C"],
      },
    ],
  },
  {
    id: "cf-3",
    category: "Control Flow",
    categoryId: "control-flow",
    order: 3,
    title: "Number Classifier",
    difficulty: "Medium",
    mode: "syntax",
    description:
      'Given `num: int = -5`, print:\n- `"positive"` if num is greater than 0\n- `"zero"` if num equals 0\n- `"negative"` if num is less than 0\n\nAlso, on a new line, print `"even"` if num is even, or `"odd"` if num is odd.',
    starterCode: `num: int = -5

# First: classify as positive, zero, or negative

# Second: classify as even or odd
# Hint: a number is even if num % 2 == 0
`,
    hint: 'Check `num > 0`, `num == 0`, and use `else` for negative. For even/odd, use `num % 2 == 0`. Note: `-5 % 2` gives `-1` in Lexo, so check `num % 2 != 0` for odd.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["negative", "odd"],
      },
    ],
  },
  {
    id: "cf-4",
    category: "Control Flow",
    categoryId: "control-flow",
    order: 4,
    title: "Nested Conditions",
    difficulty: "Medium",
    mode: "syntax",
    description:
      'Given `age: int = 20` and `hasTicket: bool = true`:\n- If age >= 18 and hasTicket is true, print `"Welcome to the show"`\n- If age >= 18 but no ticket, print `"Buy a ticket first"`\n- If under 18, print `"Too young"`',
    starterCode: `age: int = 20
hasTicket: bool = true

# Check age and ticket status
# Use nested if or logical operators
`,
    hint: "You can nest if statements: first check `if age >= 18:`, then inside check `if hasTicket:`. Or use `if age >= 18 and hasTicket:`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["Welcome to the show"],
      },
    ],
  },

  // ============================================================
  // Category 4: Loops
  // ============================================================
  {
    id: "loop-1",
    category: "Loops",
    categoryId: "loops",
    order: 1,
    title: "Counting Up",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Use a `while` loop to print the numbers 1 through 5, each on its own line.",
    starterCode: `# Use a while loop to print numbers 1 to 5
# You'll need a counter variable

`,
    hint: "Declare `i: int = 1`, then `while i <= 5:` with `print(i)` and `i = i + 1` inside the loop.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["1", "2", "3", "4", "5"],
      },
    ],
  },
  {
    id: "loop-2",
    category: "Loops",
    categoryId: "loops",
    order: 2,
    title: "For Each Loop",
    difficulty: "Easy",
    mode: "syntax",
    description:
      'Create a string array `colors` with values `"red"`, `"green"`, `"blue"`. Use a `for` loop to print each color.',
    starterCode: `# Declare a string array called colors
# with values "red", "green", "blue"

# Use a for loop to print each color
# Syntax: for item: type in array:
`,
    hint: 'Declare `colors: string[] = ["red", "green", "blue"]` then `for c: string in colors:` and `print(c)` inside.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["red", "green", "blue"],
      },
    ],
  },
  {
    id: "loop-3",
    category: "Loops",
    categoryId: "loops",
    order: 3,
    title: "Sum with While",
    difficulty: "Medium",
    mode: "syntax",
    description:
      "Use a `while` loop to calculate the sum of numbers from 1 to 10. Print the final sum.",
    starterCode: `# Calculate the sum of 1 + 2 + 3 + ... + 10
# using a while loop

total: int = 0
# Add a counter and loop here

# Print the total
`,
    hint: "Use a counter `i: int = 1` and loop while `i <= 10`. Inside the loop, add `i` to `total` and increment `i`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["55"],
      },
    ],
  },
  {
    id: "loop-4",
    category: "Loops",
    categoryId: "loops",
    order: 4,
    title: "Break and Continue",
    difficulty: "Medium",
    mode: "syntax",
    description:
      "Use a `while` loop counting from 1 to 10. Skip printing the number `5` (use `continue`), and stop the loop entirely when you reach `8` (use `break`). Print all other numbers.",
    starterCode: `# Loop from 1 to 10
# Skip 5 (use continue)
# Stop at 8 (use break) - don't print 8

i: int = 0
# Write your while loop here
`,
    hint: "Increment `i` first inside the loop, then check: `if i == 5: continue` and `if i == 8: break`. Print after the checks.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["1", "2", "3", "4", "6", "7"],
      },
    ],
  },

  // ============================================================
  // Category 5: Functions
  // ============================================================
  {
    id: "func-1",
    category: "Functions",
    categoryId: "functions",
    order: 1,
    title: "First Function",
    difficulty: "Easy",
    mode: "syntax",
    description:
      'Define a function called `square` that takes an `int` parameter `n` and returns `n * n` as an `int`. Then print the result of calling `square(7)`.',
    starterCode: `# Define a function called square
# It takes one parameter: n of type int
# It returns an int
# Syntax: def name(param: type) -> returnType:

# Call square(7) and print the result
`,
    hint: "Use `def square(n: int) -> int:` with `return n * n` in the body. Then `print(square(7))`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["49"],
      },
    ],
  },
  {
    id: "func-2",
    category: "Functions",
    categoryId: "functions",
    order: 2,
    title: "Multiple Parameters",
    difficulty: "Easy",
    mode: "syntax",
    description:
      'Define a function `max` that takes two `int` parameters `a` and `b`, and returns the larger one. Print the results of `max(10, 20)` and `max(7, 3)`.',
    starterCode: `# Define a function called max
# Takes two int parameters: a and b
# Returns the larger of the two as int

# Print max(10, 20)
# Print max(7, 3)
`,
    hint: "Use `if a > b:` inside the function to decide which to return.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["20", "7"],
      },
    ],
  },
  {
    id: "func-3",
    category: "Functions",
    categoryId: "functions",
    order: 3,
    title: "String Function",
    difficulty: "Medium",
    mode: "syntax",
    description:
      'Define a function `shout` that takes a `string` parameter `msg` and returns the message in uppercase with `"!"` appended. Print `shout("hello")` and `shout("lexo")`.',
    starterCode: `# Define a function called shout
# Takes a string parameter: msg
# Returns the uppercase version of msg with "!" at the end

# Print shout("hello")
# Print shout("lexo")
`,
    hint: 'Use `msg.upper()` to convert to uppercase, then concatenate `"!"` with `+`.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["HELLO!", "LEXO!"],
      },
    ],
  },
  {
    id: "func-4",
    category: "Functions",
    categoryId: "functions",
    order: 4,
    title: "Recursive Power",
    difficulty: "Medium",
    mode: "syntax",
    description:
      "Define a recursive function `power` that takes a `base: int` and `exp: int`, and returns `base` raised to the power of `exp`. Base case: any number to the power of 0 is 1. Print `power(2, 10)` and `power(3, 4)`.",
    starterCode: `# Define a recursive function called power
# Takes base: int and exp: int
# Returns base^exp as int
# Base case: if exp == 0, return 1
# Recursive case: return base * power(base, exp - 1)

# Print power(2, 10)
# Print power(3, 4)
`,
    hint: "The base case is `if exp == 0: return 1`. The recursive case is `return base * power(base, exp - 1)`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["1024", "81"],
      },
    ],
  },

  // ============================================================
  // Category 6: Arrays
  // ============================================================
  {
    id: "arr-1",
    category: "Arrays",
    categoryId: "arrays",
    order: 1,
    title: "Array Basics",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Create an `int[]` called `nums` with values `[10, 20, 30, 40, 50]`. Print the first element, the last element, and the length of the array.",
    starterCode: `# Declare an int array called nums
nums: int[] = [10, 20, 30, 40, 50]

# Print the first element (index 0)

# Print the last element (index 4)

# Print the length using .length()
`,
    hint: "Use `nums[0]` for the first element, `nums[4]` for the last, and `nums.length()` for the count.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["10", "50", "5"],
      },
    ],
  },
  {
    id: "arr-2",
    category: "Arrays",
    categoryId: "arrays",
    order: 2,
    title: "Array Sum",
    difficulty: "Easy",
    mode: "syntax",
    description:
      "Given an `int[]` array, use a `for` loop to compute the sum of all elements and print it.",
    starterCode: `scores: int[] = [85, 92, 78, 95, 88]

# Use a for loop to calculate the total
total: int = 0

# Loop through scores and add each to total

# Print the total
`,
    hint: "Use `for s: int in scores:` and `total = total + s` inside the loop.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["438"],
      },
    ],
  },
  {
    id: "arr-3",
    category: "Arrays",
    categoryId: "arrays",
    order: 3,
    title: "Find Maximum",
    difficulty: "Medium",
    mode: "syntax",
    description:
      "Given an `int[]` array, find and print the largest value using a loop.",
    starterCode: `values: int[] = [34, 12, 89, 56, 23, 67]

# Find the maximum value
# Start by assuming the first element is the max

# Loop through the rest and update max if you find something bigger

# Print the maximum
`,
    hint: "Set `maxVal: int = values[0]`, then loop with `for v: int in values:`. Inside, check `if v > maxVal:` and update.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["89"],
      },
    ],
  },
  {
    id: "arr-4",
    category: "Arrays",
    categoryId: "arrays",
    order: 4,
    title: "Count Matches",
    difficulty: "Medium",
    mode: "syntax",
    description:
      'Given a `string[]` of fruit names, count how many are `"apple"` and print the count.',
    starterCode: `fruits: string[] = ["apple", "banana", "apple", "cherry", "apple", "banana"]

# Count how many elements equal "apple"
count: int = 0

# Loop and count

# Print the count
`,
    hint: 'Use a for loop and `if f == "apple":` to check each element.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["3"],
      },
    ],
  },

  // ============================================================
  // Category 7: Strings
  // ============================================================
  {
    id: "str-1",
    category: "Strings",
    categoryId: "strings",
    order: 1,
    title: "String Methods",
    difficulty: "Easy",
    mode: "syntax",
    description:
      'Given `msg: string = "Hello World"`, print:\n1. The length of the string\n2. The string in all uppercase\n3. The string in all lowercase',
    starterCode: `msg: string = "Hello World"

# Print the length of msg

# Print msg in uppercase

# Print msg in lowercase
`,
    hint: "Use `msg.length()`, `msg.upper()`, and `msg.lower()`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["11", "HELLO WORLD", "hello world"],
      },
    ],
  },
  {
    id: "str-2",
    category: "Strings",
    categoryId: "strings",
    order: 2,
    title: "Substrings",
    difficulty: "Easy",
    mode: "syntax",
    description:
      'Given `word: string = "programming"`, use `.substring(start, end)` to extract and print:\n1. The first 4 characters (indices 0 to 4)\n2. Characters from index 4 to 7\n3. The last 4 characters (indices 7 to 11)',
    starterCode: `word: string = "programming"

# Print the first 4 characters: "prog"

# Print characters from index 4 to 7: "ram"

# Print the last 4 characters: "ming"
`,
    hint: "`.substring(start, end)` extracts characters from `start` up to but not including `end`. So `.substring(0, 4)` gets indices 0,1,2,3.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["prog", "ram", "ming"],
      },
    ],
  },
  {
    id: "str-3",
    category: "Strings",
    categoryId: "strings",
    order: 3,
    title: "String Concatenation",
    difficulty: "Medium",
    mode: "syntax",
    description:
      'Build a sentence from parts. Given `first: string = "Lexo"`, `second: string = "is"`, `third: string = "fun"`, concatenate them with spaces in between and print the result. Then print the length of the result.',
    starterCode: `first: string = "Lexo"
second: string = "is"
third: string = "fun"

# Concatenate with spaces: "Lexo is fun"
# and print it

# Print the length of the result
`,
    hint: 'Use `+` to concatenate: `first + " " + second + " " + third`. Store it in a variable to easily get its length.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["Lexo is fun", "11"],
      },
    ],
  },
  {
    id: "str-4",
    category: "Strings",
    categoryId: "strings",
    order: 4,
    title: "String Analysis",
    difficulty: "Medium",
    mode: "syntax",
    description:
      'Given `text: string = "Hello"`, print each of the following on its own line:\n1. The first character (substring 0 to 1)\n2. The last character (substring 4 to 5)\n3. Whether the uppercase version equals `"HELLO"` (print the boolean)',
    starterCode: `text: string = "Hello"

# Print the first character

# Print the last character

# Print whether text.upper() == "HELLO"
`,
    hint: 'Use `.substring(0, 1)` for the first char and `.substring(4, 5)` for the last. Compare with `text.upper() == "HELLO"`.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["H", "o", "true"],
      },
    ],
  },

  // ============================================================
  // Category 8: Challenges
  // ============================================================
  {
    id: "ch-1",
    category: "Challenges",
    categoryId: "challenges",
    order: 1,
    title: "FizzBuzz",
    difficulty: "Easy",
    mode: "challenge",
    description:
      'Write a function `fizzBuzz(n: int) -> void` that prints numbers from 1 to n. For multiples of 3, print `"Fizz"`. For multiples of 5, print `"Buzz"`. For multiples of both, print `"FizzBuzz"`. Test with n=5, n=15, n=1, and n=3.',
    starterCode: `# FizzBuzz: Write the function
# - Multiples of both 3 and 5: print "FizzBuzz"
# - Multiples of 3 only: print "Fizz"
# - Multiples of 5 only: print "Buzz"
# - Otherwise: print the number

def fizzBuzz(n: int) -> void:
    i: int = 1
    while i <= n:
        # Your logic here

        i = i + 1

fizzBuzz(5)
print("---")
fizzBuzz(15)
print("---")
fizzBuzz(1)
print("---")
fizzBuzz(3)
`,
    hint: "Check `i % 15 == 0` first (divisible by both), then `i % 3 == 0`, then `i % 5 == 0`, then print the number. Use `elif` for mutually exclusive cases.",
    testCases: [
      {
        id: 1,
        expectedOutput: [
          "1", "2", "Fizz", "4", "Buzz",
          "---",
          "1", "2", "Fizz", "4", "Buzz",
          "Fizz", "7", "8", "Fizz", "Buzz",
          "11", "Fizz", "13", "14", "FizzBuzz",
          "---",
          "1",
          "---",
          "1", "2", "Fizz",
        ],
      },
    ],
  },
  {
    id: "ch-2",
    category: "Challenges",
    categoryId: "challenges",
    order: 2,
    title: "Reverse Array",
    difficulty: "Easy",
    mode: "challenge",
    description:
      "Write a function `printReversed(nums: int[]) -> void` that prints the elements of an integer array in reverse order, one per line. Test with various sizes including a single element.",
    starterCode: `# Write a function to print array elements in reverse
# Hint: use a while loop starting from the last index

def printReversed(nums: int[]) -> void:
    # The last index is nums.length() - 1
    # Your code here
    return

printReversed([10, 20, 30, 40, 50])
print("---")
printReversed([1])
print("---")
printReversed([7, 3, 9, 1])
print("---")
printReversed([100, 200])
`,
    hint: "Start with `i: int = nums.length() - 1` and loop `while i >= 0:`, printing `nums[i]` and decrementing `i`.",
    testCases: [
      {
        id: 1,
        expectedOutput: [
          "50", "40", "30", "20", "10",
          "---",
          "1",
          "---",
          "1", "9", "3", "7",
          "---",
          "200", "100",
        ],
      },
    ],
  },
  {
    id: "ch-3",
    category: "Challenges",
    categoryId: "challenges",
    order: 3,
    title: "Palindrome Check",
    difficulty: "Medium",
    mode: "challenge",
    description:
      'Write a function `isPalindrome(word: string) -> bool` that returns `true` if the string reads the same forwards and backwards, `false` otherwise. Use `.substring(i, i + 1)` to get individual characters. Test with odd length, even length, single char, two chars, and a non-palindrome.',
    starterCode: `# Write a function to check if a string is a palindrome
# Use .substring(i, i + 1) to get individual characters

def isPalindrome(word: string) -> bool:
    # Compare characters from front and back
    # Return false if any pair doesn't match
    return false

print(isPalindrome("racecar"))
print(isPalindrome("hello"))
print(isPalindrome("abba"))
print(isPalindrome("a"))
print(isPalindrome("ab"))
print(isPalindrome("abcba"))
print(isPalindrome("abcca"))
`,
    hint: "Use two index variables, one starting at 0 and one at `word.length() - 1`. Compare `word.substring(left, left + 1)` with `word.substring(right, right + 1)`. If any pair differs, return `false`. After the loop, return `true`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["true", "false", "true", "true", "false", "true", "false"],
      },
    ],
  },
  {
    id: "ch-4",
    category: "Challenges",
    categoryId: "challenges",
    order: 4,
    title: "Factorial",
    difficulty: "Medium",
    mode: "challenge",
    description:
      "Write a function `factorial(n: int) -> int` that returns the factorial of n. The factorial of 0 and 1 is 1, and factorial of n is `n * (n-1) * ... * 1`. Use recursion.",
    starterCode: `# Define a function called factorial
# factorial(0) = 1
# factorial(1) = 1
# factorial(n) = n * factorial(n - 1)

def factorial(n: int) -> int:
    # Your code here
    return 0

print(factorial(0))
print(factorial(1))
print(factorial(2))
print(factorial(5))
print(factorial(7))
print(factorial(10))
`,
    hint: "Base case: `if n <= 1: return 1`. Recursive case: `return n * factorial(n - 1)`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["1", "1", "2", "120", "5040", "3628800"],
      },
    ],
  },
  {
    id: "ch-5",
    category: "Challenges",
    categoryId: "challenges",
    order: 5,
    title: "Count Vowels",
    difficulty: "Medium",
    mode: "challenge",
    description:
      'Write a function `countVowels(text: string) -> int` that returns the number of vowels (a, e, i, o, u) in the string. Handle both uppercase and lowercase. Test with mixed case, all vowels, no vowels, empty-ish, and spaces.',
    starterCode: `# Write a function to count vowels in a string
# Convert to lowercase first, then check each character

def countVowels(text: string) -> int:
    # Your code here
    return 0

print(countVowels("Hello"))
print(countVowels("aeiou"))
print(countVowels("AEIOU"))
print(countVowels("BCDFG"))
print(countVowels("Hello Beautiful World"))
print(countVowels("rhythm"))
print(countVowels("a"))
`,
    hint: 'Convert the text to lowercase with `.lower()`. Use a while loop with index from 0 to `text.length() - 1`. Get each character with `.substring(i, i + 1)` and check if it equals `"a"`, `"e"`, `"i"`, `"o"`, or `"u"`.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["2", "5", "5", "0", "7", "0", "1"],
      },
    ],
  },
  {
    id: "ch-6",
    category: "Challenges",
    categoryId: "challenges",
    order: 6,
    title: "Sum of Digits",
    difficulty: "Medium",
    mode: "challenge",
    description:
      "Write a function `digitSum(num: int) -> int` that returns the sum of digits of a positive integer. Use modulo (`% 10`) to get the last digit and integer division (`/ 10`) to remove it. Test with single digit, trailing zeros, and all nines.",
    starterCode: `# Write a function to calculate the sum of digits
# Use % 10 to get the last digit
# Use / 10 to remove the last digit

def digitSum(num: int) -> int:
    # Your code here
    return 0

print(digitSum(123))
print(digitSum(9876))
print(digitSum(5))
print(digitSum(10000))
print(digitSum(999))
print(digitSum(11111))
`,
    hint: "Initialize `total: int = 0`. While `num > 0`: add `num % 10` to total, then set `num = num / 10`. Return total.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["6", "30", "5", "1", "27", "5"],
      },
    ],
  },
  {
    id: "ch-7",
    category: "Challenges",
    categoryId: "challenges",
    order: 7,
    title: "Find Maximum",
    difficulty: "Medium",
    mode: "challenge",
    description:
      "Write a function `findMax(nums: int[]) -> int` that returns the largest value in an integer array. Assume the array has at least one element. Test with max at start, end, middle, single element, all same, and negative-like values.",
    starterCode: `# Write a function to find the maximum value in an array
# Start with the first element as your initial max

def findMax(nums: int[]) -> int:
    # Your code here
    return 0

print(findMax([3, 9, 1, 7, 5]))
print(findMax([42]))
print(findMax([1, 2, 3, 4, 5]))
print(findMax([5, 4, 3, 2, 1]))
print(findMax([10, 10, 10]))
print(findMax([1, 100, 2, 99, 3]))
`,
    hint: "Set `result: int = nums[0]`. Loop through the array with a for loop. If the current element is greater than `result`, update `result`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["9", "42", "5", "5", "10", "100"],
      },
    ],
  },
  {
    id: "ch-8",
    category: "Challenges",
    categoryId: "challenges",
    order: 8,
    title: "Two Sum",
    difficulty: "Medium",
    mode: "challenge",
    description:
      'Write a function `twoSum(nums: int[], target: int) -> void` that finds two indices whose values add up to the target and prints them separated by a space (smaller index first). There is exactly one solution. Test with pair at start, end, middle, and with duplicate values.',
    starterCode: `# Two Sum: find two indices that add up to the target
# Use nested loops to check all pairs

def twoSum(nums: int[], target: int) -> void:
    # Outer loop: index i from 0
    # Inner loop: index j from i + 1
    # Check if nums[i] + nums[j] == target
    return

twoSum([2, 7, 11, 15], 9)
twoSum([3, 2, 4], 6)
twoSum([1, 5, 3, 7], 8)
twoSum([10, 20, 30, 40], 50)
twoSum([5, 5], 10)
`,
    hint: "Use a while loop for `i` from 0, and a nested while loop for `j` from `i + 1`. When `nums[i] + nums[j] == target`, print the indices.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["0 1", "1 2", "1 2", "1 3", "0 1"],
      },
    ],
  },

  // ============================================================
  // Category 8: Challenges — Hard
  // ============================================================
  {
    id: "ch-9",
    category: "Challenges",
    categoryId: "challenges",
    order: 9,
    title: "Bubble Sort",
    difficulty: "Hard",
    mode: "challenge",
    description:
      "Write a function `bubbleSort(nums: int[]) -> void` that sorts an integer array in ascending order using the bubble sort algorithm, then prints each element on its own line. Bubble sort repeatedly steps through the array, compares adjacent elements, and swaps them if they are in the wrong order. Use index-based array mutation (`nums[i] = value`).",
    starterCode: `# Bubble Sort: sort the array in-place, then print each element
# Outer loop runs length - 1 times
# Inner loop compares adjacent pairs and swaps if out of order
# To swap: use a temp variable

def bubbleSort(nums: int[]) -> void:
    n: int = nums.length()
    # Your sorting logic here

    # Print the sorted array
    i: int = 0
    while i < n:
        print(nums[i])
        i = i + 1

bubbleSort([5, 3, 8, 1, 2])
print("---")
bubbleSort([1])
print("---")
bubbleSort([4, 3, 2, 1])
print("---")
bubbleSort([1, 2, 3, 4])
print("---")
bubbleSort([7, 7, 3, 3, 1])
`,
    hint: "Use two nested while loops. Outer: `pass` from 0 to `n - 2`. Inner: `j` from 0 to `n - pass - 2`. If `nums[j] > nums[j + 1]`, swap using a temp variable: `temp: int = nums[j]`, `nums[j] = nums[j + 1]`, `nums[j + 1] = temp`.",
    testCases: [
      {
        id: 1,
        expectedOutput: [
          "1", "2", "3", "5", "8",
          "---",
          "1",
          "---",
          "1", "2", "3", "4",
          "---",
          "1", "2", "3", "4",
          "---",
          "1", "3", "3", "7", "7",
        ],
      },
    ],
  },
  {
    id: "ch-10",
    category: "Challenges",
    categoryId: "challenges",
    order: 10,
    title: "Fibonacci Sequence",
    difficulty: "Hard",
    mode: "challenge",
    description:
      "Write a function `fibonacci(n: int) -> void` that prints the first n numbers of the Fibonacci sequence. The sequence starts with 0, 1, and each subsequent number is the sum of the two preceding ones. Handle n=0 (print nothing), n=1 (print just 0), and larger values.",
    starterCode: `# Fibonacci: print the first n Fibonacci numbers
# Sequence: 0, 1, 1, 2, 3, 5, 8, 13, 21, ...
# Handle edge cases: n=0 prints nothing, n=1 prints just 0

def fibonacci(n: int) -> void:
    # Your code here
    return

fibonacci(0)
print("---")
fibonacci(1)
print("---")
fibonacci(2)
print("---")
fibonacci(8)
print("---")
fibonacci(12)
`,
    hint: "Handle `n == 0` by returning early. Print 0 for the first. If `n >= 2`, print 1 next. Then use a while loop: compute `next = a + b`, print it, shift `a = b`, `b = next`, repeat until you've printed n numbers total.",
    testCases: [
      {
        id: 1,
        expectedOutput: [
          "---",
          "0",
          "---",
          "0", "1",
          "---",
          "0", "1", "1", "2", "3", "5", "8", "13",
          "---",
          "0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89",
        ],
      },
    ],
  },
  {
    id: "ch-11",
    category: "Challenges",
    categoryId: "challenges",
    order: 11,
    title: "Prime Checker",
    difficulty: "Hard",
    mode: "challenge",
    description:
      'Write a function `isPrime(n: int) -> bool` that returns `true` if n is a prime number, `false` otherwise. A prime number is greater than 1 and divisible only by 1 and itself. Test with small primes, non-primes, edge cases (0, 1, 2), and a larger prime.',
    starterCode: `# Write a function to check if a number is prime
# A prime is > 1 and has no divisors other than 1 and itself
# Optimization: only check divisors up to n / 2 (or sqrt approximation)

def isPrime(n: int) -> bool:
    # Your code here
    return false

print(isPrime(0))
print(isPrime(1))
print(isPrime(2))
print(isPrime(3))
print(isPrime(4))
print(isPrime(17))
print(isPrime(20))
print(isPrime(97))
print(isPrime(100))
`,
    hint: "First handle: n <= 1 returns false, n == 2 returns true. Then loop `i` from 2 to `n / 2`. If `n % i == 0`, return false. If no divisor found, return true.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["false", "false", "true", "true", "false", "true", "false", "true", "false"],
      },
    ],
  },
  {
    id: "ch-12",
    category: "Challenges",
    categoryId: "challenges",
    order: 12,
    title: "Reverse String",
    difficulty: "Hard",
    mode: "challenge",
    description:
      'Write a function `reverseString(s: string) -> string` that returns the reversed version of a string. Use `.substring(i, i + 1)` to extract characters and build the result by concatenating in reverse order. Test with odd/even lengths, single char, and palindrome.',
    starterCode: `# Reverse a string using substring and concatenation
# Build a new string by reading characters from end to start

def reverseString(s: string) -> string:
    # Your code here
    return ""

print(reverseString("hello"))
print(reverseString("a"))
print(reverseString("ab"))
print(reverseString("Lexo"))
print(reverseString("racecar"))
print(reverseString("12345"))
`,
    hint: 'Start with `result: string = ""`. Use a while loop from `i = s.length() - 1` down to 0. Each iteration: `result = result + s.substring(i, i + 1)`. Return result.',
    testCases: [
      {
        id: 1,
        expectedOutput: ["olleh", "a", "ba", "oxeL", "racecar", "54321"],
      },
    ],
  },
  {
    id: "ch-13",
    category: "Challenges",
    categoryId: "challenges",
    order: 13,
    title: "Power Function",
    difficulty: "Hard",
    mode: "challenge",
    description:
      "Write a function `power(base: int, exp: int) -> int` that calculates base raised to the exponent using a loop (no built-in power operator). Handle the edge case where exp is 0 (result is always 1).",
    starterCode: `# Calculate base^exp using a loop
# Any number raised to 0 is 1
# Multiply base by itself exp times

def power(base: int, exp: int) -> int:
    # Your code here
    return 0

print(power(2, 0))
print(power(2, 1))
print(power(2, 10))
print(power(3, 4))
print(power(5, 3))
print(power(1, 100))
print(power(10, 5))
`,
    hint: "Start with `result: int = 1`. Use a while loop that runs `exp` times, each time multiplying `result = result * base`. Return result.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["1", "2", "1024", "81", "125", "1", "100000"],
      },
    ],
  },
  {
    id: "ch-14",
    category: "Challenges",
    categoryId: "challenges",
    order: 14,
    title: "Selection Sort",
    difficulty: "Hard",
    mode: "challenge",
    description:
      "Write a function `selectionSort(nums: int[]) -> void` that sorts an integer array in ascending order using selection sort, then prints each element. Selection sort works by repeatedly finding the minimum element from the unsorted portion and putting it at the beginning.",
    starterCode: `# Selection Sort: for each position, find the minimum
# in the remaining unsorted portion and swap it into place

def selectionSort(nums: int[]) -> void:
    n: int = nums.length()
    # Your sorting logic here

    # Print sorted array
    i: int = 0
    while i < n:
        print(nums[i])
        i = i + 1

selectionSort([64, 25, 12, 22, 11])
print("---")
selectionSort([1])
print("---")
selectionSort([3, 1, 2])
print("---")
selectionSort([5, 5, 5, 1])
print("---")
selectionSort([1, 2, 3, 4, 5])
`,
    hint: "Outer loop: `i` from 0 to `n - 2`. Set `minIdx = i`. Inner loop: `j` from `i + 1` to `n - 1`. If `nums[j] < nums[minIdx]`, update `minIdx = j`. After inner loop, swap `nums[i]` and `nums[minIdx]`.",
    testCases: [
      {
        id: 1,
        expectedOutput: [
          "11", "12", "22", "25", "64",
          "---",
          "1",
          "---",
          "1", "2", "3",
          "---",
          "1", "5", "5", "5",
          "---",
          "1", "2", "3", "4", "5",
        ],
      },
    ],
  },
  {
    id: "ch-15",
    category: "Challenges",
    categoryId: "challenges",
    order: 15,
    title: "GCD (Greatest Common Divisor)",
    difficulty: "Hard",
    mode: "challenge",
    description:
      "Write a function `gcd(a: int, b: int) -> int` that returns the greatest common divisor of two positive integers using the Euclidean algorithm. The algorithm repeatedly replaces the larger number with the remainder of dividing the larger by the smaller, until one becomes zero.",
    starterCode: `# GCD using the Euclidean algorithm
# gcd(a, b) = gcd(b, a % b) until b == 0, then return a
# Example: gcd(48, 18) -> gcd(18, 12) -> gcd(12, 6) -> gcd(6, 0) -> 6

def gcd(a: int, b: int) -> int:
    # Your code here
    return 0

print(gcd(48, 18))
print(gcd(100, 75))
print(gcd(7, 3))
print(gcd(12, 12))
print(gcd(1, 999))
print(gcd(270, 192))
print(gcd(17, 13))
`,
    hint: "Use a while loop: `while b > 0:`, compute `temp: int = b`, then `b = a % b`, then `a = temp`. When the loop ends, return `a`.",
    testCases: [
      {
        id: 1,
        expectedOutput: ["6", "25", "1", "12", "1", "6", "1"],
      },
    ],
  },
];
