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

// Package objectapply applies value-backed API object envelopes under resolved
// resource contract semantics.
//
// The package is a pure object-level orchestration layer. It validates live and
// applied object envelopes against a resource definition, checks object identity
// compatibility, applies the Desired surface through api/valueapply, preserves
// live Observed data, and returns updated api/objectownership state.
//
// Version 1 applies Desired only. Observed apply is unsupported. Metadata apply
// is unsupported. The result preserves live TypeMeta, live ObjectMeta, and live
// Observed data.
//
// objectapply delegates object ownership storage shape to api/objectownership
// and field-level apply semantics to api/valueapply. It does not reimplement
// field extraction, comparison, merge, conflict detection, force semantics,
// dropped-field deletion, or ownership takeover rules.
//
// The package does not read or write storage, run admission, authorize request
// subjects, perform resource catalog lookup, convert between API versions,
// mutate runtime metadata, serialize managed fields, decode wire formats, emit
// events, or execute controller/runtime lifecycle behavior. Runtime layers are
// responsible for those concerns.
package objectapply
