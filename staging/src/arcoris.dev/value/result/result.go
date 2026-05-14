// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package result

// Result represents either OK(value) or Err(error).
//
// The zero value is OK with the zero value of T. A Result value is immutable by
// convention: methods do not mutate the receiver, and callers should treat
// contained mutable values according to the ownership contract of the
// surrounding API.
type Result[T any] struct {
	value T
	err   error
}

// OK returns a successful Result containing value.
func OK[T any](value T) Result[T] {
	return Result[T]{value: value}
}

// Err returns a failed Result containing err.
//
// Err panics if err is nil. Use OK when there is no error.
func Err[T any](err error) Result[T] {
	if err == nil {
		panic("result: nil error")
	}
	return Result[T]{err: err}
}

// From returns OK(value) when err is nil and Err(err) otherwise.
func From[T any](value T, err error) Result[T] {
	if err != nil {
		return Err[T](err)
	}
	return OK(value)
}

// IsOK reports whether r contains a successful value.
func (r Result[T]) IsOK() bool {
	return r.err == nil
}

// IsErr reports whether r contains an error.
func (r Result[T]) IsErr() bool {
	return r.err != nil
}

// Load returns the contained value and nil when r is OK.
//
// When r is Err, Load returns the zero value of T and the stored error.
func (r Result[T]) Load() (T, error) {
	if r.err == nil {
		return r.value, nil
	}

	var zero T
	return zero, r.err
}

// Err returns the stored error.
//
// Err returns nil when r is OK.
func (r Result[T]) Err() error {
	return r.err
}

// ValueOr returns the contained value when r is OK, or fallback when r is Err.
func (r Result[T]) ValueOr(fallback T) T {
	if r.err == nil {
		return r.value
	}
	return fallback
}

// Must returns the contained value when r is OK.
//
// Must panics with the stored error when r is Err. It is intended for tests and
// initialization paths where failure is a programming error. Use Load in normal
// control flow.
func (r Result[T]) Must() T {
	if r.err == nil {
		return r.value
	}
	panic(r.err)
}
