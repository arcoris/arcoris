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

// Package valuevalidation validates concrete ARCORIS API payload values against
// api/types descriptors and emits semantic diagnostics.
//
// The value pipeline keeps related responsibilities separate:
// valuevalidation validates values, valuefieldset extracts ownership fields,
// valuecompare computes semantic changed fields, valuemerge merges values,
// valueapply orchestrates apply, conflict, deletion, and ownership behavior, and
// fieldownership stores ownership state.
//
// Descriptors are expected to be prepared and structurally validated at
// construction, registration, or catalog boundaries. valuevalidation does not
// call types.ValidateResolved on every payload. Local descriptor checks in this
// package are defensive traversal checks, not full descriptor validation.
//
// Validation collects a bounded ordered list of diagnostics. Options.MaxErrors
// controls the bound. Values <= 0 use the package default, and unlimited
// collection is intentionally unsupported. MaxErrors = 1 is fail-fast.
//
// Validation errors use api/fieldpath semantic paths for payload locations.
// Record fields validated through object descriptors use field elements. Dynamic
// map entries use key elements. Ordered list items use index elements. ListMap
// items use selector elements after stable key extraction. ListSet and
// ListAtomic validation diagnostics still use indexes because item-level
// validation errors are useful even when ownership and apply layers treat the
// whole list as one field.
//
// UnknownReject reports unknown record members. UnknownPreserveOpaque accepts
// unknown members without validating their nested payload. UnknownPrune also
// accepts unknown members without validating their nested payload. Invalid
// unknown-field policies are invalid descriptors.
//
// Explicit NullValue is present data. Null is accepted only by DescriptorNull or
// nullable descriptors. DescriptorRef is resolved before applying nullability
// when the reference descriptor itself is not nullable, so reusable semantic
// types keep their own nullability contract.
//
// Malformed payload-derived record member names are invalid values. Malformed
// descriptor-declared field names are invalid descriptors. This split keeps
// concrete payload shape failures separate from schema construction failures.
//
// The package does not extract field ownership, compare values, merge values,
// apply values, validate API object metadata, access storage, decode or encode
// wire formats, apply defaults, prune payloads, or normalize values.
package valuevalidation
