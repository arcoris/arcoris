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

import (
	"fmt"

	"arcoris.dev/apimachinery/api/codec"
)

// validateCapabilities checks that semantic targets match byte or stream APIs.
func validateCapabilities(path string, c codec.BaseCodec, info codec.Info) error {
	for _, check := range capabilityChecks(c, info) {
		if check.declared == check.implemented {
			continue
		}

		return errorAt(
			path,
			ErrCapabilityMismatch,
			ErrorReasonCapabilityMismatch,
			check.detail(),
		)
	}

	return nil
}

// capabilityCheck stores one target/capability comparison.
type capabilityCheck struct {
	// target is the semantic API document model being checked.
	target codec.Target

	// byteInterface names the byte-slice capability interface for target.
	byteInterface string

	// streamInterface names the streaming capability interface for target.
	streamInterface string

	// declared records whether codec.Info advertises target.
	declared bool

	// implemented records whether the implementation has a byte or stream API.
	implemented bool
}

// capabilityChecks returns all v1 target capability comparisons.
func capabilityChecks(c codec.BaseCodec, info codec.Info) []capabilityCheck {
	return []capabilityCheck{
		{
			target:          codec.TargetValue,
			byteInterface:   "codec.ValueCodec",
			streamInterface: "codec.ValueStreamCodec",
			declared:        info.Supports(codec.TargetValue),
			implemented:     implementsValue(c),
		},
		{
			target:          codec.TargetObject,
			byteInterface:   "codec.ObjectCodec",
			streamInterface: "codec.ObjectStreamCodec",
			declared:        info.Supports(codec.TargetObject),
			implemented:     implementsObject(c),
		},
		{
			target:          codec.TargetObjectOwnership,
			byteInterface:   "codec.ObjectOwnershipCodec",
			streamInterface: "codec.ObjectOwnershipStreamCodec",
			declared:        info.Supports(codec.TargetObjectOwnership),
			implemented:     implementsObjectOwnership(c),
		},
	}
}

// detail returns a directional capability mismatch diagnostic.
func (c capabilityCheck) detail() string {
	if c.declared && !c.implemented {
		return fmt.Sprintf(
			"codec declares target %q but implements neither %s nor %s",
			c.target,
			c.byteInterface,
			c.streamInterface,
		)
	}
	if !c.declared && c.implemented {
		return fmt.Sprintf(
			"codec implements %s or %s but Info.Targets does not declare target %q",
			c.byteInterface,
			c.streamInterface,
			c.target,
		)
	}

	return fmt.Sprintf("codec target %q declared=%v implemented=%v", c.target, c.declared, c.implemented)
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
