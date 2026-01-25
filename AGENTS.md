# Ralphy-OpenSpec Agent Instructions (Claude Code)

You are an AI coding assistant operating in a repository that uses:
- **OpenSpec** for spec-driven development (`openspec/specs/` + `openspec/changes/`)
- **Ralph loop** for iterative execution (the same prompt may be repeated)

## Golden rules
- Treat `openspec/specs/` as the source of truth (current behavior).
- Treat `openspec/changes/<change-name>/` as the proposed/active change.
- Only mark tasks complete when verification (tests) passes.
- Keep changes small, deterministic, and test-backed.

## Workflow

### 1) Plan (PRD -> OpenSpec)
When asked to plan or create specs:
1. Read `openspec/project.md` and relevant files in `openspec/specs/`.
2. Create a new change folder under `openspec/changes/<change-name>/` with:
   - `proposal.md` (why/what/scope/non-goals/risks)
   - `tasks.md` (checklist with test plan notes)
   - `specs/**/spec.md` (deltas: ADDED/MODIFIED/REMOVED)
3. Ensure requirements use MUST/SHALL and each requirement has at least one scenario.

### 2) Implement (Tasks -> Code)
When asked to implement:
1. Identify the active change folder under `openspec/changes/`.
2. Implement tasks in order from `tasks.md`.
3. Run tests frequently and fix failures.
4. Update the checkbox status in `tasks.md` only when verified.

### 3) Validate (Acceptance criteria)
When asked to validate:
1. Map scenarios/acceptance criteria to tests/commands.
2. Run the project test command (commonly `npm test`).
3. Report which requirements are proven and what gaps remain.

### 4) Archive
When asked to archive:
- Prefer `openspec archive <change-name> --yes` if OpenSpec CLI is available.
- Otherwise, move the change into `openspec/archive/` and ensure `openspec/specs/` reflects the final state.

## Ralph loop completion promise
If you are being run in a loop, only output this exact text when ALL tasks are complete and tests are green:

<promise>TASK_COMPLETE</promise>

---

## Build Instructions

This project uses CGO for native library dependencies. Different build modes are available depending on which features you need.

### Prerequisites

- Go 1.21+
- GCC (for CGO compilation)
- Make (optional, for convenience targets)

### Quick Reference

| Command | Description |
|---------|-------------|
| `go build ./pkg/...` | Build packages without CGO dependencies |
| `go test ./pkg/factstore/...` | Test factstore (no CGO required) |
| `make build` | Full build with Ladybug (downloads native libs) |
| `make test` | Run all tests |

### Building Without CGO (Pure Go Packages)

Many packages work without CGO dependencies:

```bash
# Build core packages (no CGO)
go build ./pkg/factstore/...
go build ./pkg/embedder/...
go build ./pkg/nlp/...
go build ./pkg/types/...

# Run tests for pure Go packages
go test ./pkg/factstore/...
go test ./pkg/embedder/...
go test ./pkg/nlp/...
go test ./pkg/prompts/...
```

### Building With CGO (Full Feature Set)

For Ladybug embedded database and other native features:

#### Step 1: Download Native Libraries

The `go:generate` directive downloads the Ladybug native library:

```bash
# Download Ladybug library for the main CLI
go generate ./cmd/main.go

# Or for examples
go generate ./examples/basic/cgo.go
```

This runs:
```bash
curl -sL https://raw.githubusercontent.com/LadybugDB/go-ladybug/refs/heads/master/download_lbug.sh | bash -s -- -out lib-ladybug
```

#### Step 2: Build with CGO

```bash
# Build everything with system_ladybug tag
go build -tags system_ladybug ./...

# Build CLI binary
go build -tags system_ladybug -o bin/predicato ./cmd/main.go
```

#### Using Make (Recommended)

```bash
# Full build (includes go generate)
make build

# Build CLI
make build-cli

# Run server
make run-server

# Run tests
make test
```

### CGO File Structure

Each executable that uses Ladybug needs a `cgo.go` file:

```go
package main

//go:generate sh -c "curl -sL https://raw.githubusercontent.com/LadybugDB/go-ladybug/refs/heads/master/download_lbug.sh | bash -s -- -out lib-ladybug"

/*
#cgo darwin LDFLAGS: -L${SRCDIR}/lib-ladybug -Wl,-rpath,${SRCDIR}/lib-ladybug
#cgo linux LDFLAGS: -L${SRCDIR}/lib-ladybug -Wl,-rpath,${SRCDIR}/lib-ladybug
#cgo windows LDFLAGS: -L${SRCDIR}/lib-ladybug
*/
import "C"
```

### Troubleshooting

#### "cannot find -llbug"
The Ladybug native library hasn't been downloaded. Run:
```bash
go generate ./cmd/main.go
# or
make build
```

#### CGO-related test failures
Some tests require native libraries. To run only pure Go tests:
```bash
go test ./pkg/factstore/... ./pkg/embedder/... ./pkg/nlp/... ./pkg/prompts/...
```

#### LSP errors about CGO
LSP may show errors for files with CGO dependencies if libraries aren't downloaded. This doesn't affect `go build` - just run `go generate` first.

### Cross-Compilation

```bash
# Build for multiple platforms (requires native libs for each)
make build-cli-all
```

Note: Cross-compilation with CGO requires appropriate cross-compilers and native libraries for each target platform.

