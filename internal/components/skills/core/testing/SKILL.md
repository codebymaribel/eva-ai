---
name: testing
description: >
  Testing conventions, patterns, coverage expectations, and tooling.
  Load before writing or modifying tests.
sdd_phases: [schematics, execute, debrief]
trigger: "test, testing, coverage, unit test, integration test, e2e, mock, fixture"
---

# Testing

## Stack
<!-- Populated by eva init — describe the detected test framework here -->
_Detected during project initialization. Update this section with your test framework._

## Test Structure
How tests are organized — naming, location, grouping.

### Convention: test file placement
Test files live next to the code they test, in the same package/module.

```
✅ correct: src/auth/middleware_test.go tests src/auth/middleware.go
❌ incorrect: tests/auth.test.js is 3 directories away from the source
```

### Convention: test naming
Test names describe the behavior, not the implementation.

```
✅ correct: TestLogin_WithInvalidCredentials_ReturnsUnauthorized
❌ incorrect: TestLogin_Error
```

## Patterns

### Unit tests
Test one function/method in isolation. Mock external dependencies.

### Integration tests
Test the interaction between modules. Use real dependencies when possible.

### Test data
Use builders or factories for test data. Avoid shared mutable fixtures.

## Rules
- Every public function needs at least one test
- Test the behavior, not the implementation
- If it's hard to test, the design needs improvement
- Tests must be deterministic — no flaky tests
- Run the full test suite before marking a task as complete

## Do not
- Test private methods directly — test through the public API
- Use sleep or real timers in unit tests
- Share mutable state between tests
- Skip writing tests because "it's a small change"

## References
- Link to test coverage reports if available
- Link to testing strategy docs
