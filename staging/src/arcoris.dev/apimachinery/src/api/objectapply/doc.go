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

// Package objectapply orchestrates value-backed API object apply.
//
// Package object defines the value-backed object envelope. Package resource
// defines resource-family contracts and version descriptors. Package
// objectvalidation validates object envelopes against resource contracts.
// Package valueapply owns Desired field apply semantics. Package objectownership
// owns canonical object-level ownership state. Package
// objectlifecycle decides when apply is invoked as part of runtime operations.
// Package objectstore stores committed objects and ownership state.
//
// Version 1 applies Desired only. Applied Observed is rejected even when the
// resource defines an observed descriptor. Applied metadata may identify the
// live object, but metadata update intent is rejected. Successful apply
// preserves live TypeMeta, ObjectMeta, and Observed, installs the merged Desired
// value, and replaces Desired ownership through objectownership.State.WithDesired.
// Cross-version apply and API version conversion are unsupported.
//
// Applied metadata identity must match the live object. Name and namespace must
// match. Applied UID may be empty; when non-empty it must match live UID.
// Non-identity metadata such as generateName, resourceVersion, generation,
// timestamps, deletion, labels, annotations, owner references, and finalizers is
// rejected.
//
// Request.Resource is expected to be resolved and prevalidated by construction,
// registration, or catalog code. objectapply performs only defensive checks
// needed before selecting the live object version and calling lower layers. It
// does not perform catalog lookup or full resource descriptor graph validation
// per apply operation.
//
// Validation and apply are deterministic: owner, resource shape, applied
// Observed, normalized applied metadata spelling, metadata, identity, version,
// metadata policy, live object contract, applied object contract, live version
// selection, Desired valueapply, output object construction, and ownership
// replacement. Early object-level failures return a zero Result. Desired apply
// failures return partial Result.Desired when valueapply provides it. Output
// Object and replacement Ownership are populated only after Desired apply
// succeeds.
//
// The package does not perform admission, authorization, storage access,
// resource lookup, metadata apply, observed apply, defaulting, pruning, version
// conversion, codec behavior, runtime lifecycle execution, controller behavior,
// or client behavior.
package objectapply
