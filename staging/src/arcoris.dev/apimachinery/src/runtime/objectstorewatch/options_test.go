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

package objectstorewatch

import "testing"

func TestNewRejectsInvalidOptions(t *testing.T) {
	tests := []struct {
		name   string
		option Option
	}{
		{name: "nil", option: nil},
		{name: "max history zero", option: WithMaxHistory(0)},
		{name: "max history negative", option: WithMaxHistory(-1)},
		{name: "stream buffer zero", option: WithStreamBuffer(0)},
		{name: "stream buffer negative", option: WithStreamBuffer(-1)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(testBackend(t), tt.option)
			requireErrorIs(t, err, ErrInvalidOption)
		})
	}
}

func TestNewAcceptsOptions(t *testing.T) {
	store, err := New(testBackend(t), WithMaxHistory(2), WithStreamBuffer(1))
	requireNoError(t, err)

	if store.history.max != 2 {
		t.Fatalf("max history = %d; want 2", store.history.max)
	}
	if store.options.StreamBuffer != 1 {
		t.Fatalf("stream buffer = %d; want 1", store.options.StreamBuffer)
	}
}

func TestDefaultOptions(t *testing.T) {
	options := DefaultOptions()

	if options.MaxHistory != defaultMaxHistory {
		t.Fatalf("max history = %d; want %d", options.MaxHistory, defaultMaxHistory)
	}
	if options.StreamBuffer != defaultStreamBuffer {
		t.Fatalf("stream buffer = %d; want %d", options.StreamBuffer, defaultStreamBuffer)
	}
}
