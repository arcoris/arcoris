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

// Package fieldownership models ownership of semantic API field paths.
//
// The package is a small foundation layer over api/fieldpath.Set. It tracks
// which owners own which semantic fields, detects ownership conflicts for
// attempted field sets, and provides deterministic state transformation helpers.
//
// Owner identities are field-ownership identities. They are not admission roles,
// authorization subjects, RBAC principals, runtime component instances, request
// actors, or policy identities. Higher layers decide which request subject may
// act as which field owner.
//
// Conflicts reports ownership overlap for the supplied attempted field set.
// Callers are responsible for deciding which fields count as attempted changes:
// specified fields, changed fields, modified and removed fields, force-taken
// fields, or another higher-layer policy set.
//
// Admission may authorize which request subject can act as which field owner,
// but that is a higher-layer policy. fieldownership only models ownership
// state. The package does not inspect api/value payloads, traverse api/types
// descriptors, compare values, extract field sets, merge values, apply
// configurations, serialize managed fields, validate object metadata, access
// storage, or perform admission.
package fieldownership
