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

This project uses CGO for the Ladybug embedded graph database. The build system handles native library setup automatically via Make targets.

### Prerequisites

- Go 1.21+
- GCC (for CGO compilation)
- Make (recommended)
- curl (for downloading native libraries)

### Quick Start (Recommended)

```bash
# Download native libs + build everything
make build

# Run all tests
make test

# Run only pure Go tests (no CGO setup needed)
make test-nocgo
```

### Quick Reference

| Command | Description |
|---------|-------------|
| `make generate` | Download Ladybug native library |
| `make build` | Full build (includes generate) |
| `make build-cli` | Build CLI binary to `bin/predicato` |
| `make test` | Run all tests (CGO + non-CGO) |
| `make test-cgo` | Run tests requiring Ladybug |
| `make test-nocgo` | Run pure Go tests (no CGO needed) |
| `make run-server` | Start the HTTP server |

### Building Without CGO (Pure Go Packages)

Many packages work without CGO dependencies:

```bash
# Build core packages (no CGO)
go build ./pkg/factstore/...
go build ./pkg/embedder/...
go build ./pkg/nlp/...
go build ./pkg/types/...

# Run tests for pure Go packages
make test-nocgo
# or manually:
go test ./pkg/factstore/... ./pkg/embedder/... ./pkg/nlp/... ./pkg/prompts/... ./pkg/logger/...
```

### Building With CGO (Full Feature Set)

For Ladybug embedded database, use the Makefile which handles CGO_LDFLAGS automatically:

```bash
# Download native library + build
make build

# Or step by step:
make generate                    # Download liblbug.so to cmd/lib-ladybug/
make build-cli                   # Build CLI binary
```

#### Manual CGO Build (without Make)

If you need to build manually without Make:

```bash
# Step 1: Download Ladybug library
go generate ./cmd/main.go

# Step 2: Set library path and build
export CGO_LDFLAGS="-L$(pwd)/cmd/lib-ladybug -Wl,-rpath,$(pwd)/cmd/lib-ladybug"
go build -tags system_ladybug ./...

# Step 3: Run tests
go test -tags system_ladybug ./pkg/driver/...
```

### CGO File Structure

The main entry point (`cmd/main.go`) contains CGO directives:

```go
package main

//go:generate sh -c "curl -sL https://raw.githubusercontent.com/LadybugDB/go-ladybug/refs/heads/master/download_lbug.sh | bash -s -- -out lib-ladybug"

/*
#cgo darwin LDFLAGS: -L${SRCDIR}/lib-ladybug -llbug -Wl,-rpath,${SRCDIR}/lib-ladybug
#cgo linux LDFLAGS: -L${SRCDIR}/lib-ladybug -llbug -Wl,-rpath,${SRCDIR}/lib-ladybug
#cgo windows LDFLAGS: -L${SRCDIR}/lib-ladybug -llbug_shared
#include <stdlib.h>
*/
import "C"
```

### Troubleshooting

#### "cannot find -llbug"
The Ladybug native library hasn't been downloaded:
```bash
make generate
# or: go generate ./cmd/main.go
```

#### CGO tests fail to link
Ensure you're using the Makefile targets which set `CGO_LDFLAGS`:
```bash
make test-cgo   # Correct way
# Instead of: go test ./pkg/driver/...
```

#### LSP errors about CGO
LSP may show errors for CGO files if libraries aren't downloaded. Run `make generate` first. These errors don't affect command-line builds.

#### Run tests without CGO setup
Use the `test-nocgo` target:
```bash
make test-nocgo
```

### Test Package CGO Requirements

| Package | CGO Required | Reason |
|---------|--------------|--------|
| `pkg/types` | No | Pure Go types and validation |
| `pkg/factstore` | No | Database clients only |
| `pkg/embedder` | No | HTTP client wrappers |
| `pkg/nlp` | No | LLM API clients |
| `pkg/prompts` | No | Template parsing |
| `pkg/logger` | No | Logging utilities |
| `pkg/driver` | Yes | Ladybug embedded graph |
| `pkg/checkpoint` | Yes | Depends on driver |
| `pkg/modeler` | Yes | Depends on driver |
| `pkg/utils` | Yes | Some helpers depend on driver types |
| `pkg/search` | Yes | Depends on driver |
| `pkg/server` | Yes | Transitive driver dependency |
| `pkg/community` | Yes | Graph operations |

**Quick test commands:**
```bash
# Pure Go packages (no setup required)
go test ./pkg/types/... ./pkg/factstore/... ./pkg/embedder/... ./pkg/nlp/...

# CGO packages (requires make generate first)
make test-cgo
```

### Cross-Compilation

```bash
# Build for multiple platforms
make build-cli-all
```

Note: Cross-compilation with CGO requires appropriate cross-compilers and native libraries for each target platform.

