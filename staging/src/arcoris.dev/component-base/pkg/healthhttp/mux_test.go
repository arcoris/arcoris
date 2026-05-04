/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package healthhttp

import (
	"errors"
	"net/http"
	"testing"
)

func TestErrNilMux(t *testing.T) {
	t.Parallel()

	if ErrNilMux == nil {
		t.Fatal("ErrNilMux is nil")
	}
	if !errors.Is(ErrNilMux, ErrNilMux) {
		t.Fatal("errors.Is(ErrNilMux, ErrNilMux) = false, want true")
	}
}

func TestNilMux(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		mux  Mux
		want bool
	}{
		{name: "nil interface", mux: nil, want: true},
		{name: "typed nil pointer", mux: (*recordingMux)(nil), want: true},
		{name: "non nil pointer", mux: newRecordingMux(), want: false},
		{name: "http serve mux", mux: http.NewServeMux(), want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := nilMux(test.mux); got != test.want {
				t.Fatalf("nilMux() = %v, want %v", got, test.want)
			}
		})
	}
}
