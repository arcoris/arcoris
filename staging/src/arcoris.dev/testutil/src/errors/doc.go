/*
Package errorassert provides small, standard-library-only assertions for Go
error contracts in tests.

The package covers generic error behavior: non-nil checks, errors.Is,
errors.As, direct unwrap checks, and exact error messages. It must not contain
domain-specific ARCORIS errors, classifiers, fixtures, or state assertions.
*/
package errorassert
