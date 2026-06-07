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

import "errors"

// ErrNilCheckResolverFunc identifies a nil CheckResolverFunc.
var ErrNilCheckResolverFunc = errors.New("health: nil check resolver function")

// CheckResolver resolves the checks that belong to one concrete target.
//
// A resolver owns lookup only. It must not execute checks, normalize results,
// aggregate reports, apply target policy, cache report output, expose transports,
// or schedule probe loops. No checks for a concrete target should be represented
// as an empty CheckSet.
//
// Implementations should return checks in deterministic order and should be
// safe for concurrent use unless documented otherwise.
type CheckResolver interface {
	ResolveChecks(target Target) (CheckSet, error)
}

// CheckResolverFunc adapts a function into CheckResolver.
type CheckResolverFunc func(target Target) (CheckSet, error)

// ResolveChecks calls f.
//
// ResolveChecks panics when f is nil because a nil resolver function is a
// programming error, not an empty check set.
func (f CheckResolverFunc) ResolveChecks(target Target) (CheckSet, error) {
	if f == nil {
		panic(ErrNilCheckResolverFunc)
	}

	return f(target)
}
