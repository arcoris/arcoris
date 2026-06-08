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

// Transport identifies the concrete codec transport capability to select.
type Transport uint8

const (
	// TransportBytes selects byte-slice codec methods.
	TransportBytes Transport = iota + 1

	// TransportStream selects io.Reader/io.Writer codec methods.
	TransportStream
)

// IsZero reports whether t is absent.
func (t Transport) IsZero() bool {
	return t == 0
}

// String returns stable diagnostic text for t.
func (t Transport) String() string {
	switch t {
	case TransportBytes:
		return "bytes"
	case TransportStream:
		return "stream"
	default:
		return ""
	}
}

// Validate checks that t is one of the supported codec transports.
func (t Transport) Validate() error {
	switch t {
	case TransportBytes, TransportStream:
		return nil
	default:
		return errorfAt(
			pathTransport,
			ErrInvalidBinding,
			ErrorReasonInvalidBinding,
			"codec selection transport %d is not supported",
			t,
		)
	}
}
