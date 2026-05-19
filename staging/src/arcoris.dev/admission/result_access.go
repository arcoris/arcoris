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

package admission

// Decision returns the semantic admission decision.
//
// The returned value is copyable and contains no typed grant or metadata. Use
// Grant and Metadata to read those optional typed values.
func (r Result[G, M]) Decision() Decision {
	return r.decision
}

// Grant returns the caller-owned grant when the result contains one.
//
// The boolean return value is false when the decision did not transfer typed
// ownership to the caller. Callers must interpret the grant according to the
// domain package that created the Result.
func (r Result[G, M]) Grant() (G, bool) {
	return r.grant.Load()
}

// Metadata returns the typed metadata when the result contains it.
//
// Metadata is intended for snapshots, diagnostics, and read models associated
// with the decision. It is not required for a Result to be valid.
func (r Result[G, M]) Metadata() (M, bool) {
	return r.metadata.Load()
}
