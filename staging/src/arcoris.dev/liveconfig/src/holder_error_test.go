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

package liveconfig

import "testing"

func TestNilHolderPanics(t *testing.T) {
	tests := []struct {
		name string
		call func()
	}{
		{
			name: "Snapshot",
			call: func() {
				var h *Holder[testConfig]
				_ = h.Snapshot()
			},
		},
		{
			name: "Stamped",
			call: func() {
				var h *Holder[testConfig]
				_ = h.Stamped()
			},
		},
		{
			name: "Revision",
			call: func() {
				var h *Holder[testConfig]
				_ = h.Revision()
			},
		},
		{
			name: "LastError",
			call: func() {
				var h *Holder[testConfig]
				_ = h.LastError()
			},
		},
		{
			name: "Apply",
			call: func() {
				var h *Holder[testConfig]
				_, _ = h.Apply(testConfig{})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if got := recover(); got != ErrNilHolder {
					t.Fatalf("panic = %v, want %v", got, ErrNilHolder)
				}
			}()

			tt.call()
		})
	}
}

func TestZeroValueHolderPanics(t *testing.T) {
	tests := []struct {
		name string
		call func(h *Holder[testConfig])
	}{
		{name: "Snapshot", call: func(h *Holder[testConfig]) { _ = h.Snapshot() }},
		{name: "Stamped", call: func(h *Holder[testConfig]) { _ = h.Stamped() }},
		{name: "Revision", call: func(h *Holder[testConfig]) { _ = h.Revision() }},
		{name: "LastError", call: func(h *Holder[testConfig]) { _ = h.LastError() }},
		{name: "Apply", call: func(h *Holder[testConfig]) { _, _ = h.Apply(testConfig{}) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Fatal("panic = nil, want panic")
				}
			}()

			var h Holder[testConfig]
			tt.call(&h)
		})
	}
}
