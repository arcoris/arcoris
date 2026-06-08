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

// entrySupportsCapability reports whether entry exposes target over transport.
func entrySupportsCapability(entry codecregistry.Entry, target codec.Target, transport Transport) bool {
	c := entry.Codec()
	switch transport {
	case TransportBytes:
		return entrySupportsByteCapability(c, target)
	case TransportStream:
		return entrySupportsStreamCapability(c, target)
	default:
		return false
	}
}

// entrySupportsByteCapability checks byte-slice capabilities for target.
func entrySupportsByteCapability(c codec.BaseCodec, target codec.Target) bool {
	switch target {
	case codec.TargetValue:
		_, ok := c.(codec.ValueCodec)
		return ok
	case codec.TargetObject:
		_, ok := c.(codec.ObjectCodec)
		return ok
	case codec.TargetObjectOwnership:
		_, ok := c.(codec.ObjectOwnershipCodec)
		return ok
	default:
		return false
	}
}

// entrySupportsStreamCapability checks streaming capabilities for target.
func entrySupportsStreamCapability(c codec.BaseCodec, target codec.Target) bool {
	switch target {
	case codec.TargetValue:
		_, ok := c.(codec.ValueStreamCodec)
		return ok
	case codec.TargetObject:
		_, ok := c.(codec.ObjectStreamCodec)
		return ok
	case codec.TargetObjectOwnership:
		_, ok := c.(codec.ObjectOwnershipStreamCodec)
		return ok
	default:
		return false
	}
}
