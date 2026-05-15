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

package maybe

// Maybe represents either Some(value) or None.
//
// The zero value is None. A Maybe value is immutable by convention: methods do
// not mutate the receiver, and callers should treat contained mutable values
// according to the ownership contract of the surrounding API.
type Maybe[T any] struct {
	value T
	ok    bool
}

// Some returns a Maybe containing value.
func Some[T any](val T) Maybe[T] {
	return Maybe[T]{value: val, ok: true}
}

// None returns a Maybe with no value.
func None[T any]() Maybe[T] {
	return Maybe[T]{}
}

// From returns Some(value) when ok is true and None otherwise.
//
// When ok is false, value is discarded so the returned Maybe does not retain
// references owned by the caller.
func From[T any](val T, ok bool) Maybe[T] {
	if ok {
		return Some(val)
	}
	return None[T]()
}

// IsSome reports whether m contains a value.
func (m Maybe[T]) IsSome() bool {
	return m.ok
}

// IsNone reports whether m contains no value.
func (m Maybe[T]) IsNone() bool {
	return !m.ok
}

// Load returns the contained value and true when m is Some.
//
// When m is None, Load returns the zero value of T and false.
func (m Maybe[T]) Load() (T, bool) {
	if m.ok {
		return m.value, true
	}

	var zero T
	return zero, false
}

// ValueOr returns the contained value when m is Some, or fallback when m is
// None.
func (m Maybe[T]) ValueOr(fallback T) T {
	if m.ok {
		return m.value
	}
	return fallback
}

// Must returns the contained value when m is Some.
//
// Must panics when m is None. It is intended for tests and initialization paths
// where absence is a programming error. Use Load in normal control flow.
func (m Maybe[T]) Must() T {
	if m.ok {
		return m.value
	}
	panic("maybe: none has no value")
}
