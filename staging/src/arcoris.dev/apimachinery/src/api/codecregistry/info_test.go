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

func TestCloneInfoReturnsDetachedMetadata(t *testing.T) {
	info := codec.Info{
		Format:     codec.FormatJSON,
		MediaTypes: []codec.MediaType{codec.MediaTypeJSON},
		Targets:    []codec.Target{codec.TargetValue},
	}

	cloned := cloneInfo(info)
	cloned.MediaTypes[0] = codec.MediaTypeYAML
	cloned.Targets[0] = codec.TargetObject

	if info.MediaTypes[0] != codec.MediaTypeJSON {
		t.Fatalf("source media type mutated: %q", info.MediaTypes[0])
	}
	if info.Targets[0] != codec.TargetValue {
		t.Fatalf("source target mutated: %q", info.Targets[0])
	}
}

func TestCloneInfoPreservesFormat(t *testing.T) {
	info := codec.Info{Format: codec.FormatJSON}

	if got := cloneInfo(info).Format; got != codec.FormatJSON {
		t.Fatalf("Format = %q", got)
	}
}
