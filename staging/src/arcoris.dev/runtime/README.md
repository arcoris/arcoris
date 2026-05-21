# arcoris.dev/runtime

## Status

Experimental staging module.

## Purpose

Process-local runtime coordination primitives such as run groups, waits,
lifecycle control, and signal integration.

## Source layout

The Go module root is `pkg/`. Published repositories promote `pkg/` to the
repository root.

## Testing

`cd pkg && go test ./...`
