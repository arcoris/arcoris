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

package codec

import "testing"

func TestInfoNormalizeReturnsDetachedSlices(t *testing.T) {
	info := Info{
		Format:     " JSON ",
		MediaTypes: []MediaType{" Application/JSON "},
		Targets:    []Target{TargetValue},
	}

	got, err := NormalizeInfo(info)
	requireNoError(t, err)

	info.MediaTypes[0] = MediaTypeYAML
	info.Targets[0] = TargetObject

	requireMediaTypes(t, got.MediaTypes, MediaTypeJSON)
	requireTargets(t, got.Targets, TargetValue)
}

func TestInfoNormalizeSortsMediaTypes(t *testing.T) {
	info := Info{
		Format:     FormatJSON,
		MediaTypes: []MediaType{MediaTypeYAML, MediaTypeJSON},
		Targets:    []Target{TargetValue},
	}

	got, err := info.Normalize()
	requireNoError(t, err)

	requireMediaTypes(t, got.MediaTypes, MediaTypeJSON, MediaTypeYAML)
}

func TestInfoNormalizeSortsTargets(t *testing.T) {
	info := Info{
		Format:     FormatJSON,
		MediaTypes: []MediaType{MediaTypeJSON},
		Targets:    []Target{TargetObjectOwnership, TargetObject, TargetValue},
	}

	got, err := info.Normalize()
	requireNoError(t, err)

	requireTargets(t, got.Targets, TargetObject, TargetObjectOwnership, TargetValue)
}

func TestInfoNormalizeRejectsDuplicateAfterCanonicalization(t *testing.T) {
	info := Info{
		Format:     FormatJSON,
		MediaTypes: []MediaType{MediaTypeJSON, " Application/JSON "},
		Targets:    []Target{TargetValue},
	}

	_, err := info.Normalize()

	requireErrorIs(t, err, ErrInvalidInfo)
}
