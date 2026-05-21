# Publishing Staged Modules

## Source of truth

`arcoris/arcoris` is the source of truth. Standalone repositories are generated
distribution mirrors of the staged module trees.

## Publishing transform

Publishing promotes the module source root and copies selected metadata:

```text
staging/src/arcoris.dev/<module>/pkg/*        -> published repo root/*
staging/src/arcoris.dev/<module>/README.md    -> published repo root/README.md
staging/src/arcoris.dev/<module>/SECURITY.md  -> published repo root/SECURITY.md
staging/src/arcoris.dev/<module>/LICENSE      -> published repo root/LICENSE
staging/src/arcoris.dev/<module>/docs/*       -> published repo root/docs/*
```

## Manifest

Each staged module is expected to define a `publishing.yaml` manifest similar
to:

```yaml
module: arcoris.dev/admission
repository: arcoris/admission
source_root: pkg
publish:
  promote:
    pkg: .
  include:
    - README.md
    - SECURITY.md
    - LICENSE
    - docs/**
    - pkg/**
  exclude:
    - publishing.yaml
    - .staging/**
    - tmp/**
```

## Verification

Publishing verification must check:

- `pkg/go.mod` module path matches the manifest `module`
- public imports do not contain `/pkg`
- the published tree can run `go test ./...`
- required metadata files are present
- staging-only files are not published

## Manual edits policy

Direct edits to published mirror repositories are discouraged. Fixes should go
through the staging source in `arcoris/arcoris`, then be republished from the
authoritative tree.
