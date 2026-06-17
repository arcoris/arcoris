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

package objectwatch

import "testing"

func TestCapabilitiesSupportsStart(t *testing.T) {
	capabilities := Capabilities{StartAtCurrent: true, HistoricalStart: true}

	requireNoError(t, capabilities.SupportsStart(AtCurrent()))
	requireNoError(t, capabilities.SupportsStart(Start{Mode: StartAfterRevision, Revision: 1}))
}

func TestCapabilitiesRejectsUnsupportedStart(t *testing.T) {
	tests := []struct {
		name         string
		capabilities Capabilities
		start        Start
	}{
		{name: "current", capabilities: Capabilities{}, start: AtCurrent()},
		{name: "historical", capabilities: Capabilities{}, start: Start{Mode: StartAfterRevision, Revision: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.capabilities.SupportsStart(tt.start)
			requireErrorIs(t, err, ErrInvalidStart)
		})
	}
}

func TestCapabilitiesSupportsRequest(t *testing.T) {
	capabilities := Capabilities{StartAtCurrent: true}
	request := Request{List: watchListRequest(), Start: AtCurrent(), AllowBookmarks: true}

	requireNoError(t, capabilities.SupportsRequest(request))
}

func TestCapabilitiesSupportsRequestPreservesInvalidRequest(t *testing.T) {
	capabilities := Capabilities{StartAtCurrent: true}
	request := Request{Start: AtCurrent()}

	err := capabilities.SupportsRequest(request)

	requireErrorIs(t, err, ErrInvalidRequest)
}

func TestCapabilitiesSupportsRequestWrapsUnsupportedStart(t *testing.T) {
	capabilities := Capabilities{}
	request := Request{List: watchListRequest(), Start: AtCurrent()}

	err := capabilities.SupportsRequest(request)

	requireErrorIs(t, err, ErrInvalidRequest)
	requireErrorIs(t, err, ErrInvalidStart)
	requireWatchError(t, err, ErrorReasonInvalidRequest, "watch.request.start")
}
