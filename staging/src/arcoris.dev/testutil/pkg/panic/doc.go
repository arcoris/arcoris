/*
Package panicassert provides standard-library-only panic assertions for ARCORIS
tests.

The package intentionally focuses on a narrow contract surface:

  - Require verifies that a call panics and returns the recovered value.
  - RequireNone verifies that a call does not panic.
  - RequireMessage verifies the formatted panic message.
  - RequireValue verifies a typed recovered value using reflect.DeepEqual.
  - RequireAs verifies only the recovered type.

Stable panic contracts are part of the programmer-facing validation surface in
several ARCORIS modules. Those checks are generic enough to share across
modules without pulling in package-specific fixtures, fake clocks, or richer
assertion frameworks.

panicassert is intended for tests only. Production packages must not depend on
arcoris.dev/testutil/panic.
*/
package panicassert
