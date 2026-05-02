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

func TestSubscriptionConfigDefaults(t *testing.T) {
	t.Parallel()

	config := newSubscribeConfig()

	if config.buffer != 1 {
		t.Fatalf("buffer = %d, want 1", config.buffer)
	}
	if config.notifier == nil {
		t.Fatal("default notifier is nil")
	}
}

func TestSubscriptionConfigAppliesOptionsAndIgnoresNil(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	config := newSubscribeConfig(nil, WithSubscriptionBuffer(4), withNotifier(n))

	if config.buffer != 4 {
		t.Fatalf("buffer = %d, want 4", config.buffer)
	}
	if config.notifier != n {
		t.Fatal("notifier option was not applied")
	}
}

func TestSubscriptionConfigAppliesBufferOptionsInOrder(t *testing.T) {
	t.Parallel()

	config := newSubscribeConfig(
		WithSubscriptionBuffer(2),
		WithSubscriptionBuffer(5),
	)

	if config.buffer != 5 {
		t.Fatalf("buffer = %d, want 5", config.buffer)
	}
}

func TestSubscriptionConfigRejectsInvalidBuffer(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int
	}{
		{name: "zero", size: 0},
		{name: "negative", size: -1},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNonPositiveSubscriptionBuffer, func() {
				newSubscribeConfig(WithSubscriptionBuffer(tt.size))
			})
		})
	}
}

func TestSubscriptionConfigRepairsNilNotifier(t *testing.T) {
	t.Parallel()

	config := newSubscribeConfig(withNotifier(nil))

	if config.notifier == nil {
		t.Fatal("nil notifier was not repaired")
	}
}
