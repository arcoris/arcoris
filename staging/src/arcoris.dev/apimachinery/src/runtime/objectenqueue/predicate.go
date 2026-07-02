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

package objectenqueue

import "arcoris.dev/apimachinery/api/objectstore"

// Predicate decides whether a committed object change should be mapped.
type Predicate interface {
	Match(objectstore.Change) (bool, error)
}

// PredicateFunc adapts a function to Predicate.
type PredicateFunc func(objectstore.Change) (bool, error)

// Match calls f with change.
func (f PredicateFunc) Match(change objectstore.Change) (bool, error) {
	if f == nil {
		return false, ErrNilPredicate
	}

	return f(change)
}

// Filter returns a mapper that runs mapper only when predicate matches.
func Filter(predicate Predicate, mapper Mapper) Mapper {
	return MapperFunc(func(change objectstore.Change, emit EmitFunc) error {
		if isNilInterface(predicate) {
			return ErrNilPredicate
		}
		if isNilInterface(mapper) {
			return ErrNilMapper
		}

		ok, err := predicate.Match(change)
		if err != nil || !ok {
			return err
		}

		return mapper.Map(change, emit)
	})
}
