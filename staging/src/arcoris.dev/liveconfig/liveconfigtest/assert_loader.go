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

package liveconfigtest

// RequireLoadCalls fails the test when loader has consumed a different number
// of scripted results than want.
func RequireLoadCalls[T any](t TestingT, loader *Loader[T], want int) {
	t.Helper()

	if got := loader.Calls(); got != want {
		t.Fatalf("loader calls = %d, want %d", got, want)
	}
}

// RequireLoadRemaining fails the test when loader has a different number of
// remaining scripted results than want.
func RequireLoadRemaining[T any](t TestingT, loader *Loader[T], want int) {
	t.Helper()

	if got := loader.Remaining(); got != want {
		t.Fatalf("loader remaining = %d, want %d", got, want)
	}
}

// RequireLoaderExhausted fails the test when loader still has scripted results.
func RequireLoaderExhausted[T any](t TestingT, loader *Loader[T]) {
	t.Helper()

	if remaining := loader.Remaining(); remaining != 0 {
		t.Fatalf("loader remaining = %d, want exhausted", remaining)
	}
}
