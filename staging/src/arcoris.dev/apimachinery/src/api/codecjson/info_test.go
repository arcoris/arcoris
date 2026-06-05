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

import (
	"reflect"
	"testing"

	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

func TestInfo(t *testing.T) {
	info := newTestCodec(t).Info()

	if info.Format != codec.FormatJSON {
		t.Fatalf("format = %q; want %q", info.Format, codec.FormatJSON)
	}
	if !reflect.DeepEqual(info.MediaTypes, []codec.MediaType{codec.MediaTypeJSON}) {
		t.Fatalf("media types = %#v", info.MediaTypes)
	}
	if !reflect.DeepEqual(info.Targets, []codec.Target{
		codec.TargetValue,
		codec.TargetObject,
		codec.TargetObjectOwnership,
	}) {
		t.Fatalf("targets = %#v", info.Targets)
	}
	requireNoError(t, info.Validate())
}

func TestInfoDetached(t *testing.T) {
	c := newTestCodec(t)
	info := c.Info()
	info.MediaTypes[0] = codec.MediaTypeYAML
	info.Targets[0] = codec.TargetObject

	fresh := c.Info()
	if fresh.MediaTypes[0] != codec.MediaTypeJSON {
		t.Fatalf("Info media type was aliased: %#v", fresh.MediaTypes)
	}
	if fresh.Targets[0] != codec.TargetValue {
		t.Fatalf("Info targets were aliased: %#v", fresh.Targets)
	}
}

func TestRegistryIntegration(t *testing.T) {
	registry, err := codecregistry.New(
		codecregistry.Register(codecregistry.MustEntryID("json.public"), newTestCodec(t)),
	)
	requireNoError(t, err)

	if candidates := registry.FullCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullCandidates(%q) length = %d; want 1", codec.MediaTypeJSON, len(candidates))
	}
	if candidates := registry.FullStreamCandidates(codec.MediaTypeJSON); len(candidates) != 1 {
		t.Fatalf("FullStreamCandidates(%q) length = %d; want 1", codec.MediaTypeJSON, len(candidates))
	}
}
