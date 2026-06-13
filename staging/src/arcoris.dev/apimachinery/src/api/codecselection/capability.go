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

package codecselection

import (
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

// targetTransportSpec binds one target/transport pair to its runtime interface.
type targetTransportSpec struct {
	// target is the semantic document model requested by the binding.
	target codec.Target

	// transport is the concrete byte or stream method family requested.
	transport Transport

	// supports checks the runtime codec interface for this exact capability.
	supports func(codec.BaseCodec) bool
}

// targetTransportSpecs is the single internal dispatch table for binding
// capability checks. Adding a new codec.Target should require adding its
// transport rows here and typed public selection wrappers, not duplicating the
// validation switch in multiple paths.
var targetTransportSpecs = []targetTransportSpec{
	{target: codec.TargetValue, transport: TransportBytes, supports: supportsValueBytes},
	{target: codec.TargetValue, transport: TransportStream, supports: supportsValueStream},
	{target: codec.TargetObject, transport: TransportBytes, supports: supportsObjectBytes},
	{target: codec.TargetObject, transport: TransportStream, supports: supportsObjectStream},
	{target: codec.TargetObjectOwnership, transport: TransportBytes, supports: supportsObjectOwnershipBytes},
	{target: codec.TargetObjectOwnership, transport: TransportStream, supports: supportsObjectOwnershipStream},
}

// entrySupportsCapability reports whether entry exposes target over transport.
func entrySupportsCapability(entry codecregistry.Entry, target codec.Target, transport Transport) bool {
	spec, ok := lookupTargetTransportSpec(target, transport)
	if !ok {
		return false
	}

	return spec.supports(entry.Codec())
}

// lookupTargetTransportSpec returns the known capability spec for target and transport.
func lookupTargetTransportSpec(target codec.Target, transport Transport) (targetTransportSpec, bool) {
	for _, spec := range targetTransportSpecs {
		if spec.target == target && spec.transport == transport {
			return spec, true
		}
	}

	return targetTransportSpec{}, false
}

// supportsValueBytes reports whether c exposes byte-slice value methods.
func supportsValueBytes(c codec.BaseCodec) bool {
	_, ok := c.(codec.ValueCodec)
	return ok
}

// supportsValueStream reports whether c exposes streaming value methods.
func supportsValueStream(c codec.BaseCodec) bool {
	_, ok := c.(codec.ValueStreamCodec)
	return ok
}

// supportsObjectBytes reports whether c exposes byte-slice object methods.
func supportsObjectBytes(c codec.BaseCodec) bool {
	_, ok := c.(codec.ObjectCodec)
	return ok
}

// supportsObjectStream reports whether c exposes streaming object methods.
func supportsObjectStream(c codec.BaseCodec) bool {
	_, ok := c.(codec.ObjectStreamCodec)
	return ok
}

// supportsObjectOwnershipBytes reports whether c exposes byte-slice ownership methods.
func supportsObjectOwnershipBytes(c codec.BaseCodec) bool {
	_, ok := c.(codec.ObjectOwnershipCodec)
	return ok
}

// supportsObjectOwnershipStream reports whether c exposes streaming ownership methods.
func supportsObjectOwnershipStream(c codec.BaseCodec) bool {
	_, ok := c.(codec.ObjectOwnershipStreamCodec)
	return ok
}
