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

// Unit is a zero-size request marker for admitters that need no request data.
//
// It lets generic admitter signatures stay explicit without allocating or
// inventing meaningless request structs.
type Unit struct{}

// NoGrant is a zero-size grant marker for results that cannot return ownership.
//
// Use it when a Result shape should make grants impossible at the type level.
type NoGrant struct{}

// NoMetadata is a zero-size metadata marker for results that carry no metadata.
//
// Use it when a Result shape intentionally has no snapshot, diagnostics, or
// other read model attached to the decision.
type NoMetadata struct{}
