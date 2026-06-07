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

package health

// CheckSet binds a concrete health target to an ordered immutable checker set.
//
// CheckSet is a root contract value, not a mutable registry. It lets resolvers
// return validated check ownership without exposing their storage internals.
// Values are safe to pass by value. Methods never expose mutable internal slice
// storage.
//
// The zero value is invalid because it has TargetUnknown. Use NewCheckSet with a
// concrete target to construct an empty or populated set.
type CheckSet struct {
	// target is the concrete health target shared by every checker in the set.
	target Target

	// checks stores the immutable resolver order. Constructors copy caller input
	// and accessors return defensive copies.
	checks []Checker

	// names is the lookup index for Has. It is derived from checks and never
	// mutated after construction.
	names map[string]struct{}
}

// NewCheckSet returns an immutable set of checks for target.
//
// target must be concrete. Each checker must be non-nil and expose a valid
// stable check name. Duplicate names within the set are rejected. Check order is
// preserved exactly as supplied.
func NewCheckSet(target Target, checks ...Checker) (CheckSet, error) {
	if !target.IsConcrete() {
		return CheckSet{}, InvalidTargetError{Target: target}
	}

	prepared, names, err := prepareCheckSet(checks)
	if err != nil {
		return CheckSet{}, err
	}

	return CheckSet{
		target: target,
		checks: prepared,
		names:  names,
	}, nil
}

// MustCheckSet returns a check set and panics if construction fails.
//
// MustCheckSet is intended for package-level fixtures and tests where invalid
// checker declarations are programming errors.
func MustCheckSet(target Target, checks ...Checker) CheckSet {
	set, err := NewCheckSet(target, checks...)
	if err != nil {
		panic(err)
	}

	return set
}

// Target returns the concrete target owned by s.
func (s CheckSet) Target() Target {
	return s.target
}

// Checks returns a defensive copy of the checks in resolver order.
func (s CheckSet) Checks() []Checker {
	if len(s.checks) == 0 {
		return nil
	}

	checks := make([]Checker, len(s.checks))
	copy(checks, s.checks)

	return checks
}

// Len returns the number of checks in s.
func (s CheckSet) Len() int {
	return len(s.checks)
}

// Empty reports whether s contains no checks.
func (s CheckSet) Empty() bool {
	return len(s.checks) == 0
}

// Has reports whether s contains a check named name.
func (s CheckSet) Has(name string) bool {
	if !ValidCheckName(name) || len(s.names) == 0 {
		return false
	}

	_, ok := s.names[name]
	return ok
}

// Range calls fn for each checker in order until fn returns false.
//
// Range avoids allocating a defensive slice for callers that only need ordered
// iteration. A nil fn stops immediately.
func (s CheckSet) Range(fn func(Checker) bool) {
	if fn == nil {
		return
	}

	for _, checker := range s.checks {
		if !fn(checker) {
			return
		}
	}
}

// IsValid reports whether s satisfies CheckSet invariants.
func (s CheckSet) IsValid() bool {
	if !s.target.IsConcrete() {
		return false
	}

	_, _, err := prepareCheckSet(s.checks)
	return err == nil
}
