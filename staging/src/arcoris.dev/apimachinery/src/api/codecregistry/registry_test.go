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
	"testing"

	"arcoris.dev/apimachinery/api/codec"
)

func TestRegistryZeroValueHasNoIndexes(t *testing.T) {
	var registry Registry

	if registry.entries != nil || registry.byFormat != nil || registry.byMediaType != nil {
		t.Fatalf("zero registry = %#v; want nil storage", registry)
	}
}

func TestRegistryIndexesPointAtSortedEntries(t *testing.T) {
	registry, err := New(
		newValueByteCodec(codec.FormatYAML, codec.MediaTypeYAML),
		newValueByteCodec(codec.FormatJSON, codec.MediaTypeJSON),
	)
	requireNoError(t, err)

	jsonIndex := registry.byFormat[codec.FormatJSON]
	yamlIndex := registry.byFormat[codec.FormatYAML]
	if len(jsonIndex) != 1 || jsonIndex[0] != 0 {
		t.Fatalf("json format indexes = %#v; want [0]", jsonIndex)
	}
	if len(yamlIndex) != 1 || yamlIndex[0] != 1 {
		t.Fatalf("yaml format indexes = %#v; want [1]", yamlIndex)
	}
	if got := registry.byMediaType[codec.MediaTypeYAML]; got != yamlIndex[0] {
		t.Fatalf("yaml media index = %d; want %d", got, yamlIndex[0])
	}
}
