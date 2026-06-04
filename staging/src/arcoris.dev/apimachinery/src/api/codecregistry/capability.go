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

package codecregistry

import "arcoris.dev/apimachinery/api/codec"

// validateCapabilities checks that semantic targets match byte or stream APIs.
func validateCapabilities(path string, c codec.BaseCodec, info codec.Info) error {
	for _, check := range capabilityChecks(c, info) {
		if check.declared == check.implemented {
			continue
		}

		return errorfAt(
			path,
			ErrCapabilityMismatch,
			ErrorReasonCapabilityMismatch,
			"codec target %q declared=%v implemented=%v",
			check.target,
			check.declared,
			check.implemented,
		)
	}

	return nil
}

// capabilityCheck stores one target/capability comparison.
type capabilityCheck struct {
	// target is the semantic API document model being checked.
	target codec.Target

	// declared records whether codec.Info advertises target.
	declared bool

	// implemented records whether the implementation has a byte or stream API.
	implemented bool
}

// capabilityChecks returns all v1 target capability comparisons.
func capabilityChecks(c codec.BaseCodec, info codec.Info) []capabilityCheck {
	return []capabilityCheck{
		{
			target:      codec.TargetValue,
			declared:    info.Supports(codec.TargetValue),
			implemented: implementsValue(c),
		},
		{
			target:      codec.TargetObject,
			declared:    info.Supports(codec.TargetObject),
			implemented: implementsObject(c),
		},
		{
			target:      codec.TargetObjectOwnership,
			declared:    info.Supports(codec.TargetObjectOwnership),
			implemented: implementsObjectOwnership(c),
		},
	}
}

// implementsValue reports byte or streaming support for value documents.
func implementsValue(c codec.BaseCodec) bool {
	_, byteOK := c.(codec.ValueCodec)
	_, streamOK := c.(codec.ValueStreamCodec)

	return byteOK || streamOK
}

// implementsObject reports byte or streaming support for object documents.
func implementsObject(c codec.BaseCodec) bool {
	_, byteOK := c.(codec.ObjectCodec)
	_, streamOK := c.(codec.ObjectStreamCodec)

	return byteOK || streamOK
}

// implementsObjectOwnership reports byte or streaming support for ownership docs.
func implementsObjectOwnership(c codec.BaseCodec) bool {
	_, byteOK := c.(codec.ObjectOwnershipCodec)
	_, streamOK := c.(codec.ObjectOwnershipStreamCodec)

	return byteOK || streamOK
}
