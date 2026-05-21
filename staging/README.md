# ARCORIS Staging Area

## Purpose

`staging/` contains the authoritative sources for publishable ARCORIS modules.
Each module under `staging/src/arcoris.dev/<module>` is intended to be
independently publishable, while the monorepo remains the single source of
truth for development, review, and release preparation.

## Source-root layout

The canonical staging layout is:

```text
staging/src/arcoris.dev/<module>/
  README.md
  SECURITY.md
  LICENSE
  docs/
  publishing.yaml
  pkg/
    go.mod
    ...
```

The parent directory is the staging metadata root. `pkg/` is the real Go module
root. `go.work` points to `pkg/`, and publishing promotes `pkg/` to the root of
the standalone repository.

## Import path invariant

Local workspace import paths and published import paths must stay identical.
Do not import `/pkg`, and do not place `go.mod` in the staging metadata root.

Valid examples:

```text
arcoris.dev/admission
arcoris.dev/resilience/bulkhead
```

Invalid examples:

```text
arcoris.dev/admission/pkg
arcoris.dev/resilience/pkg/bulkhead
```

## Publishing model

Published repositories are distribution mirrors. Changes are made in the
`arcoris/arcoris` staging source. Publishing automation promotes `pkg/` to the
repository root and copies selected metadata files from the parent staging root.

## Local development

Use `go.work` from the repository root. Run module tests from each `pkg/` root.
Do not manually edit generated published mirrors.
