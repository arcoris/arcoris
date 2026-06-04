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

package codecjson

import "arcoris.dev/apimachinery/api/codec"

// Info returns detached metadata for the JSON codec implementation.
//
// The returned slices are newly allocated literals on each call. Registries may
// normalize and store their own snapshots without sharing mutable metadata with
// the codec value.
func (c Codec) Info() codec.Info {
	return codec.Info{
		Format:     codec.FormatJSON,
		MediaTypes: []codec.MediaType{codec.MediaTypeJSON},
		Targets: []codec.Target{
			codec.TargetValue,
			codec.TargetObject,
			codec.TargetObjectOwnership,
		},
	}
}
