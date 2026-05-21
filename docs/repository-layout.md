# Repository Layout

## Design goals

- stable import paths;
- independently publishable modules;
- monorepo as source of truth;
- clean standalone repositories;
- explicit metadata boundary;
- deterministic publishing.

## Staging module structure

```text
staging/src/arcoris.dev/<module>/
  README.md
  SECURITY.md
  LICENSE
  docs/
  publishing.yaml
  pkg/
```

## What belongs in the metadata root

Belongs in the parent staging root:

- `README.md`
- `SECURITY.md`
- `LICENSE`
- `docs/`
- `publishing.yaml`
- repository metadata or templates consumed by publishing tooling

Does not belong in the parent staging root:

- Go source files
- `go.mod`
- `go.sum`
- package tests
- internal implementation packages

## What belongs in pkg/

Belongs in `pkg/`:

- `go.mod`
- `go.sum`
- `doc.go`
- `*.go`
- `*_test.go`
- `internal/`
- subpackages
- `testdata/`
- examples that are part of the published Go module

Does not belong in `pkg/`:

- `publishing.yaml`
- staging-only scripts
- staging-only docs
- generated mirror metadata unless the published repository explicitly needs it

## Import path rules

`pkg/` is never part of public import paths. The module path in `pkg/go.mod`
remains `arcoris.dev/<module>`, and packages below `pkg/` keep normal import
paths.

Example:

```text
staging/src/arcoris.dev/resilience/pkg/bulkhead
-> arcoris.dev/resilience/bulkhead
```

## Internal packages

`internal/` inside `pkg/` is internal to the published module.
`internal/` outside `pkg/` is staging-only and must not be published unless a
future publishing rule explicitly includes it.

## Future commands and tools

Product or tool modules may use `cmd/` inside `pkg/` when the command is part
of the published module. Staging-only tools should live outside published
`pkg/` trees or in repository-level tooling directories.
