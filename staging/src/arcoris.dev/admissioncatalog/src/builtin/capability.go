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

package builtin

import "arcoris.dev/admissioncatalog"

// capabilities keeps descriptor literals readable without hiding the public
// admissioncatalog capability model.
func capabilities(
	outcomes admissioncatalog.OutcomeSet,
	effects admissioncatalog.EffectSet,
) admissioncatalog.CapabilitySet {
	return admissioncatalog.NewCapabilitySet(outcomes, effects)
}

// outcomes is local shorthand for builtin descriptor literals.
func outcomes(capabilities ...admissioncatalog.OutcomeCapability) admissioncatalog.OutcomeSet {
	return admissioncatalog.NewOutcomeSet(capabilities...)
}

// effects is local shorthand for builtin descriptor literals.
func effects(capabilities ...admissioncatalog.EffectCapability) admissioncatalog.EffectSet {
	return admissioncatalog.NewEffectSet(capabilities...)
}
