# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Essential Commands

**Core Development Workflow:**
- `make generate` - Regenerates all code from Ent schemas and OpenAPI spec. **Always run this before starting work.**
- `make test` - Runs all tests with gotestsum
- `make coverage` - Generates HTML coverage reports in coverage/coverage.html
- `make build` - Builds the application (runs generate first)

**Database Commands:**
- `go run commands/migrate/main.go --reset` - Drops and recreates database with fresh schema
- `go run commands/seed/main.go --seed=basic` - Seeds database with basic test data
- `git add -A && GIT_PAGER=cat git diff HEAD` - Review current changes

## Architecture Overview

**Onion Architecture Implementation:**
- `application/handler/` - HTTP handlers (API layer) - contains DB access as exception to onion rules
- `domain/service/` - Business logic services with interfaces for mocking
- `domain/infrastructure/` - Infrastructure interfaces (for dependency inversion)
- `infrastructure/` - Infrastructure implementations (databases, APIs, etc.)
- `ent/schema/` - Database entity definitions (generates type-safe ORM code)
- `restapi/` - Auto-generated OpenAPI server code

**Code Generation Strategy:**
This project heavily relies on code generation. The `make generate` command:
1. Generates Ent ORM code from `/ent/schema/*.go` files
2. Generates OpenAPI server code from `swagger.yml`
3. Generates mocks for interfaces with `//go:generate` directives

**Dependency Injection:**
Uses Uber FX framework in `/di/container.go`. All dependencies are wired through this container, making the application highly testable.

## File Naming Conventions

**Handlers:** `{method}{endpoint}` format
- Example: `GET /companies/{id}/users` → `getCompaniesUsers.go`
- Keep files flat in `application/handler/`, no subfolders

**Services:** `{resource}{action}Service` format  
- Example: `userRegisterService.go`, `emailSendService.go`
- Always include interface definitions for mockability

## Testing Approach

**Test Structure:**
- Tests are end-to-end through all layers (no layer mocking)
- Database transactions auto-rollback (using go-txdb)
- External services are mocked via generated interfaces
- Use `testUtil.TestMain(m)` and `testUtil.Prepare(t)` for setup

**Mock Generation:**
Add this directive to interface files:
```go
//go:generate go tool mockgen -source=$GOFILE -destination=${GOFILE}_mock.go -package=$GOPACKAGE
```

## Critical Development Rules

1. **Always run `make generate` before starting work** - errors here require consultation
2. **Do not edit `swagger.yml` or `ent/schema/` files** without explicit instruction
3. **DB access is allowed in handlers as exception** - but extract to single-purpose functions
4. **Validation:** Simple validation in OpenAPI spec, complex validation in handlers
5. **Error handling:** Use `eris.Wrap(err, "")` for error propagation
6. **Precision for numbers:** Always specify precision for numeric types

## Key Architectural Decisions

**Database Access Pattern:**
Unlike pure onion architecture, DB access via Ent is permitted directly in handlers for simplicity, but must be extracted into dedicated functions (one function per SQL operation).

**Validation Strategy:**
- Field-level validation: Handled by generated OpenAPI code
- Business logic validation: Implemented in handlers
- Cross-field/DB-dependent validation: Custom validation in handlers

**Testing Philosophy:**
- Full integration testing through all layers
- Mock only external dependencies, not internal layers
- Prioritize test readability over DRY principles