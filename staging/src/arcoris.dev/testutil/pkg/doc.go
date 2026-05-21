/*
Package testutil provides small, standard-library-only helpers for ARCORIS
tests.

The initial version intentionally exposes only the panic assertion subpackage at
arcoris.dev/testutil/panic. Stable panic contracts already appear across
several staging modules and can be shared safely without pulling in
domain-specific fixtures.

The module must not import ARCORIS production packages. Domain-specific helpers
belong in package-local test_helpers_test.go files or in module-local internal
test packages with a clear ownership boundary.

testutil is intended for tests only. Production packages must not depend on it.
*/
package testutil
