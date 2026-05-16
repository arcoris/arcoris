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

import (
	"context"
	"errors"
	"testing"
)

func TestLoaderReturnsScriptedResults(t *testing.T) {
	first := NewConfigVersion(1)
	second := NewConfigVersion(2)
	loader := NewLoader(Loaded(first), Loaded(second))

	next, ok := loader.Peek()
	if !ok {
		t.Fatal("Peek() ok = false, want true")
	}
	if !next.Succeeded() || next.Failed() {
		t.Fatalf("Peek() state succeeded=%v failed=%v, want success", next.Succeeded(), next.Failed())
	}
	RequireConfigEqual(t, next.Value, first)

	got, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	RequireConfigEqual(t, got, first)

	got, err = loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	RequireConfigEqual(t, got, second)

	RequireLoadCalls(t, loader, 2)
	RequireLoadRemaining(t, loader, 0)
	RequireLoaderExhausted(t, loader)
}

func TestLoaderReturnsScriptedErrors(t *testing.T) {
	want := errors.New("boom")
	failed := LoadFailed[Config](want)
	if !failed.Failed() || failed.Succeeded() {
		t.Fatalf("LoadFailed state failed=%v succeeded=%v, want failed", failed.Failed(), failed.Succeeded())
	}
	loader := NewLoader[Config](failed)

	_, err := loader.Load(context.Background())
	if !errors.Is(err, want) {
		t.Fatalf("Load() error = %v, want %v", err, want)
	}
}

func TestLoaderDoesNotConsumeResultWhenContextDone(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	loader := NewLoader(Loaded(NewConfig()))
	_, err := loader.Load(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Load() error = %v, want %v", err, context.Canceled)
	}
	RequireLoadCalls(t, loader, 0)
	RequireLoadRemaining(t, loader, 1)
}

func TestLoaderReturnsErrNoLoadResult(t *testing.T) {
	loader := NewLoader[Config]()

	_, err := loader.Load(context.Background())
	if !errors.Is(err, ErrNoLoadResult) {
		t.Fatalf("Load() error = %v, want %v", err, ErrNoLoadResult)
	}
}

func TestLoaderAppendAndReset(t *testing.T) {
	loader := NewLoader(Loaded(NewConfigVersion(1)))

	_, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	RequireLoaderExhausted(t, loader)

	loader.Append(Loaded(NewConfigVersion(2)))
	RequireLoadRemaining(t, loader, 1)

	got, err := loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	RequireConfigEqual(t, got, NewConfigVersion(2))

	loader.Reset(Loaded(NewConfigVersion(3)))
	RequireLoadCalls(t, loader, 0)
	RequireLoadRemaining(t, loader, 1)

	got, err = loader.Load(context.Background())
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	RequireConfigEqual(t, got, NewConfigVersion(3))
}

func TestLoaderPeekExhausted(t *testing.T) {
	loader := NewLoader[Config]()

	if got, ok := loader.Peek(); ok {
		t.Fatalf("Peek() = %#v, true; want zero, false", got)
	}
	if !loader.Exhausted() {
		t.Fatal("Exhausted() = false, want true")
	}
}
