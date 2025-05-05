# Prompt for Generating Test Code

You are writing a Go test function. Follow these guidelines to generate clear, idiomatic, and maintainable test code based on the implementation being tested.

## General Structure
- Import necessary packages, including `"testing"` and others as needed (e.g., `"errors"`, `"github.com/stretchr/testify/assert"`).
- Use the standard test function signature: `func TestXxx(t *testing.T)`.
- If testing a method, follow the naming convention:
  `Test<StructName>_<MethodName>_<OptionalDescriptor>` (e.g., `TestHandler_Process_EmptyInput`)
- Define multiple test cases as a map:
  ~~~go
  tests := map[string]struct {
      input    ...
      want     ...
      wantErr  bool
  }{
      "case name": {input: ..., want: ..., wantErr: ...},
      ...
  }
  ~~~
- Use `t.Run(name, func(t *testing.T) { ... })` to run each test case.
- If the test case is safe for concurrent execution, call `t.Parallel()` inside each subtest.

## Implementation Steps
1. Begin by analyzing the target implementation to determine:
  - Whether it's a standalone function or a method.
  - What parameters it takes and what it returns.
  - How it handles errors and edge cases.
  - Whether it has side effects or external dependencies.
2. Define the fields in the test case struct (e.g., `input`, `want`, `wantErr`).
3. Set up any required test data, mocks, or context for each test case.
4. Call the function under test with the case's inputs.
5. Use assertions to check:
  - Whether the result matches expectations (`want`)
  - Whether an error occurred if `wantErr` is true
  - Include informative failure messages in `t.Errorf` to aid debugging (include inputs and outputs).
  - Test edge cases such as nil inputs, empty slices, or invalid parameters, if applicable.

This structure ensures that test code is accurately aligned with the behavior of the function under test, and remains easy to maintain and scale.
