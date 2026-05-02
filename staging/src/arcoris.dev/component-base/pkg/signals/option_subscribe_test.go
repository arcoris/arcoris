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

package signals

import "testing"

func TestSubscribeConfigDefaults(t *testing.T) {
	config := newSubscribeConfig()

	if config.buffer != 1 {
		t.Fatalf("buffer = %d, want 1", config.buffer)
	}
	if config.notifier == nil {
		t.Fatal("default notifier is nil")
	}
}

func TestSubscribeConfigAppliesOptionsAndIgnoresNil(t *testing.T) {
	n := &fakeNotifier{}
	config := newSubscribeConfig(nil, WithBuffer(4), withNotifier(n))

	if config.buffer != 4 {
		t.Fatalf("buffer = %d, want 4", config.buffer)
	}
	if config.notifier != n {
		t.Fatal("notifier option was not applied")
	}
}

func TestWithBufferRejectsNonPositiveSize(t *testing.T) {
	mustPanicWith(t, errNonPositiveSubscribeBuffer, func() {
		newSubscribeConfig(WithBuffer(0))
	})
	mustPanicWith(t, errNonPositiveSubscribeBuffer, func() {
		newSubscribeConfig(WithBuffer(-1))
	})
}

func TestSubscribeConfigRepairsNilNotifier(t *testing.T) {
	config := newSubscribeConfig(withNotifier(nil))
	if config.notifier == nil {
		t.Fatal("nil notifier was not repaired")
	}
}
