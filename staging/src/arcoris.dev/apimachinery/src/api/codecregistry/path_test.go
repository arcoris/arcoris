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

import "testing"

func TestCodecPath(t *testing.T) {
	if got := codecPath(7); got != "codecs[7]" {
		t.Fatalf("codecPath() = %q", got)
	}
}

func TestInfoPath(t *testing.T) {
	if got := infoPath(7); got != "codecs[7].info" {
		t.Fatalf("infoPath() = %q", got)
	}
}

func TestMediaTypePath(t *testing.T) {
	if got := mediaTypePath(7, 2); got != "codecs[7].info.mediaTypes[2]" {
		t.Fatalf("mediaTypePath() = %q", got)
	}
}

func TestCapabilityPath(t *testing.T) {
	if got := capabilityPath(7); got != "codecs[7].capabilities" {
		t.Fatalf("capabilityPath() = %q", got)
	}
}
