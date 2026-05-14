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

	cfg := newSubscribeConfig()

	if cfg.buffer != 1 {
		t.Fatalf("buffer = %d, want 1", cfg.buffer)
	}
	if cfg.notifier == nil {
		t.Fatal("default notifier is nil")
	}
}

func TestSubscriptionConfigAppliesOptionsAndIgnoresNil(t *testing.T) {
	t.Parallel()

	n := &fakeNotifier{}
	cfg := newSubscribeConfig(nil, WithSubscriptionBuffer(4), withNotifier(n))

	if cfg.buffer != 4 {
		t.Fatalf("buffer = %d, want 4", cfg.buffer)
	}
	if cfg.notifier != n {
		t.Fatal("notifier option was not applied")
	}
}

func TestSubscriptionConfigAppliesBufferOptionsInOrder(t *testing.T) {
	t.Parallel()

	cfg := newSubscribeConfig(
		WithSubscriptionBuffer(2),
		WithSubscriptionBuffer(5),
	)

	if cfg.buffer != 5 {
		t.Fatalf("buffer = %d, want 5", cfg.buffer)
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNonPositiveSubscriptionBuffer, func() {
				newSubscribeConfig(WithSubscriptionBuffer(tc.size))
			})
		})
	}
}

func TestSubscriptionConfigRepairsNilNotifier(t *testing.T) {
	t.Parallel()

	cfg := newSubscribeConfig(withNotifier(nil))

	if cfg.notifier == nil {
		t.Fatal("nil notifier was not repaired")
	}
}
