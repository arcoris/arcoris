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

package bulkhead

// noCopy marks values that must not be copied after first use.
//
// The methods intentionally match go vet's copylocks convention. The type has no
// runtime behavior; it exists only to make accidental copies visible during
// static analysis.
type noCopy struct{}

// Lock is a marker method recognized by go vet's copylocks analyzer.
func (*noCopy) Lock() {}

// Unlock is a marker method recognized by go vet's copylocks analyzer.
func (*noCopy) Unlock() {}
