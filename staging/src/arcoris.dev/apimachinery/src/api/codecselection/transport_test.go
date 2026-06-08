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

import "testing"

func TestTransportString(t *testing.T) {
	tests := []struct {
		name      string
		transport Transport
		want      string
	}{
		{name: "bytes", transport: TransportBytes, want: "bytes"},
		{name: "stream", transport: TransportStream, want: "stream"},
		{name: "unknown", transport: Transport(99), want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.transport.String(); got != tt.want {
				t.Fatalf("String() = %q; want %q", got, tt.want)
			}
		})
	}
}

func TestTransportValidate(t *testing.T) {
	for _, transport := range []Transport{TransportBytes, TransportStream} {
		if err := transport.Validate(); err != nil {
			t.Fatalf("Validate(%s) unexpected error: %v", transport, err)
		}
	}
}

func TestTransportValidateRejectsUnknown(t *testing.T) {
	err := Transport(99).Validate()

	requireErrorIs(t, err, ErrInvalidBinding)
	requireSelectionError(t, err, "codecselection.transport", ErrorReasonInvalidBinding)
}
