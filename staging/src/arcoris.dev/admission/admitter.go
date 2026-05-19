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

// Admitter is the generic contract for a component that attempts to admit work.
//
// R is the request type, G is the grant type, and M is the metadata type. Domain
// packages should choose precise types instead of erasing values to any.
type Admitter[R any, G any, M any] interface {
	TryAdmit(R) Result[G, M]
}

// AdmitterFunc adapts a function to Admitter.
//
// A nil AdmitterFunc will panic when called, matching normal function-call
// semantics. Components that need optional admitters should model that
// optionality outside this adapter.
type AdmitterFunc[R any, G any, M any] func(R) Result[G, M]

// TryAdmit calls f(req) and returns the admission result unchanged.
func (f AdmitterFunc[R, G, M]) TryAdmit(req R) Result[G, M] {
	return f(req)
}
