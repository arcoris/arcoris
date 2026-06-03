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

package admissioncatalog

import "arcoris.dev/admission"

// ReasonDescriptor describes stable open-world metadata for an admission
// reason.
//
// The descriptor catalogs a machine-readable reason code and its usual
// capability surface. It does not close the admission.Reason vocabulary, does
// not register itself globally, and does not enforce Result validity. Reason
// validation remains syntax-level here; registry membership is an optional
// catalog concern for docs, config checks, and higher-level chain validation.
type ReasonDescriptor struct {
	// Reason is the stable open-world reason described by this catalog entry.
	Reason admission.Reason

	// Capabilities records the outcome and effect surface commonly associated
	// with the reason. The zero set is valid and means unspecified capabilities.
	Capabilities CapabilitySet
}

// IsValid reports whether d is syntactically valid catalog metadata.
//
// Validation checks only the reason identifier and capability bits. Duplicate
// detection belongs to ReasonRegistry, and ordinary Result construction does not
// require a reason registry lookup.
func (d ReasonDescriptor) IsValid() bool {
	return d.Reason.IsValid() && d.Capabilities.IsValid()
}
