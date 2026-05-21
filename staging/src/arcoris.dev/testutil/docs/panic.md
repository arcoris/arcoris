# Panic Assertions

Stable panic assertions exist because several ARCORIS modules treat panic
contracts as part of their programmer-facing validation surface. Tests should
verify those contracts directly instead of only checking that "something
panicked".

## Panic message vs panic value

A panic can carry any Go value. Some tests only care about the formatted
message, while others intentionally care about the exact typed value.

- Use `panicassert.RequireMessage` when the contract is the diagnostic text.
- Use `panicassert.RequireValue` when the contract is the typed panic value
  itself.
- Use `panicassert.RequireAs` when the test needs a typed recovered value but does not
  need deep equality against an expected value.

## Why `fmt.Sprint` is used for message checks

Message checks should work for string panics, error panics, and arbitrary
panic values that format to a stable diagnostic string. `fmt.Sprint` keeps that
surface simple and standard-library-only.

## Why `reflect.DeepEqual` is used for typed value checks

Typed panic values may be slices, maps, structs, or other values that cannot be
compared with `==`. `reflect.DeepEqual` keeps the helper safe for those test
cases without requiring package-specific comparators.

## Migration examples

Before:

```go
requirePanic(t, "capacity.Ledger: nil ledger", func() {
	ledger.TryReserve(1)
})
```

After:

```go
panicassert.RequireMessage(t, "capacity.Ledger: nil ledger", func() {
	ledger.TryReserve(1)
})
```

Before:

```go
mustPanicWithValue(t, want, func() {
	...
})
```

After:

```go
panicassert.RequireValue(t, want, func() {
	...
})
```
