---
description: Golang best practices and specific patterns for interface usage and project structure
globs: "**/*.go"
---

# Go Memory

Best practices and patterns for writing idiomatic and maintainable Go code.

## Interface Definition

### Define Interfaces Where Used
Define interfaces in the package where they are consumed (used), not in the package where they are implemented.
- **Why**: Keeps dependencies loose and allows the consumer to define exactly what behavior it needs.
- **Pattern**: `Service` defines `Repository` interface; `Controller` defines `Service` interface.
- **Don't**: Do not export interfaces from the implementation package (e.g., `repository.UserRepository` interface inside `repository` is an anti-pattern if only used elsewhere).
