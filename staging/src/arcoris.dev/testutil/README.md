# arcoris.dev/testutil

## Status

Experimental staging module.

## Purpose

`testutil` provides small, standard-library-only helpers for ARCORIS tests.
The initial scope is the panic assertion subpackage at
`arcoris.dev/testutil/panic`.

## Scope

- `panicassert.Require`
- `panicassert.RequireNone`
- `panicassert.RequireMessage`
- `panicassert.RequireValue`
- `panicassert.RequireAs`

## Non-goals

- no production runtime helpers;
- no domain-specific ARCORIS assertions;
- no admission/capacity/snapshot/resilience helpers;
- no fake clocks;
- no error/channel/concurrency/time helpers in the initial version;
- no replacement for the Go testing package;
- no external assertion framework.

## Usage

```go
panicassert.RequireMessage(t, "capacity.Ledger: nil ledger", func() {
	ledger.TryReserve(1)
})

value := panicassert.RequireValue(t, myPanicValue, func() {
	panic(myPanicValue)
})
_ = value
```

## Import policy

`arcoris.dev/testutil/panic` is intended for tests only. Production packages
must not import it.

## Testing

`cd pkg && go test ./...`
