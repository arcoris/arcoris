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

package objectreflector

import (
	"context"
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
)

var _ Sink = (*FanoutSink)(nil)

func TestNewFanoutSinkValidation(t *testing.T) {
	tests := []struct {
		name    string
		sinks   []Sink
		wantErr error
	}{
		{name: "no sinks", wantErr: ErrNoSinks},
		{name: "nil sink", sinks: []Sink{nil}, wantErr: ErrNilSink},
		{name: "typed nil sink", sinks: []Sink{(*fanoutRecordingSink)(nil)}, wantErr: ErrNilSink},
		{name: "one sink", sinks: []Sink{newFanoutRecordingSink("first", nil)}},
		{
			name: "multiple sinks",
			sinks: []Sink{
				newFanoutRecordingSink("first", nil),
				newFanoutRecordingSink("second", nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink, err := NewFanoutSink(tt.sinks...)
			if tt.wantErr != nil {
				requireErrorIs(t, err, tt.wantErr)
				if sink != nil {
					t.Fatalf("sink = %#v; want nil", sink)
				}
				return
			}

			requireNoError(t, err)
			if sink == nil {
				t.Fatalf("sink is nil")
			}
		})
	}
}

func TestNewFanoutSinkCopiesSinkList(t *testing.T) {
	original := newFanoutRecordingSink("original", nil)
	replacement := newFanoutRecordingSink("replacement", nil)
	input := []Sink{original}
	sink, err := NewFanoutSink(input...)
	requireNoError(t, err)

	input[0] = replacement
	requireNoError(t, sink.ApplyChange(context.Background(), createdChange(t, testKey("system", 1), 2)))

	requireOrder(t, original.localOrder, "original:apply")
	requireOrder(t, replacement.localOrder)
}

func TestFanoutSinkReplaceOrderAndPassthrough(t *testing.T) {
	ctx := context.WithValue(context.Background(), fanoutContextKey{}, "ctx")
	read := testRead(t, 7, listItem(testKey("system", 1), 7, "listed"))
	order := make([]string, 0, 2)
	first := newFanoutRecordingSink("first", &order)
	second := newFanoutRecordingSink("second", &order)
	sink, err := NewFanoutSink(first, second)
	requireNoError(t, err)

	err = sink.Replace(ctx, read)

	requireNoError(t, err)
	requireOrder(t, order, "first:replace", "second:replace")
	first.requireReplace(t, ctx, read)
	second.requireReplace(t, ctx, read)
}

func TestFanoutSinkApplyChangeOrderAndPassthrough(t *testing.T) {
	ctx := context.WithValue(context.Background(), fanoutContextKey{}, "ctx")
	change := createdChange(t, testKey("system", 1), 2)
	order := make([]string, 0, 2)
	first := newFanoutRecordingSink("first", &order)
	second := newFanoutRecordingSink("second", &order)
	sink, err := NewFanoutSink(first, second)
	requireNoError(t, err)

	err = sink.ApplyChange(ctx, change)

	requireNoError(t, err)
	requireOrder(t, order, "first:apply", "second:apply")
	first.requireApplyChange(t, ctx, change)
	second.requireApplyChange(t, ctx, change)
}

func TestFanoutSinkReplaceStopsOnFirstError(t *testing.T) {
	wantErr := errors.New("replace failed")
	order := make([]string, 0, 2)
	first := newFanoutRecordingSink("first", &order)
	first.replaceErr = wantErr
	second := newFanoutRecordingSink("second", &order)
	sink, err := NewFanoutSink(first, second)
	requireNoError(t, err)

	err = sink.Replace(context.Background(), testRead(t, 1))

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	requireOrder(t, order, "first:replace")
	if second.replaceCount() != 0 {
		t.Fatalf("second replace calls = %d; want 0", second.replaceCount())
	}
}

func TestFanoutSinkApplyChangeStopsOnFirstError(t *testing.T) {
	wantErr := errors.New("apply failed")
	order := make([]string, 0, 2)
	first := newFanoutRecordingSink("first", &order)
	first.applyErr = wantErr
	second := newFanoutRecordingSink("second", &order)
	sink, err := NewFanoutSink(first, second)
	requireNoError(t, err)

	err = sink.ApplyChange(context.Background(), createdChange(t, testKey("system", 1), 2))

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	requireOrder(t, order, "first:apply")
	if second.applyCount() != 0 {
		t.Fatalf("second apply calls = %d; want 0", second.applyCount())
	}
}

func TestFanoutSinkInvalidReceiver(t *testing.T) {
	valid := newFanoutRecordingSink("valid", nil)
	var typedNil Sink = (*fanoutRecordingSink)(nil)
	tests := []struct {
		name string
		sink *FanoutSink
	}{
		{name: "nil", sink: nil},
		{name: "empty", sink: &FanoutSink{}},
		{name: "nil internal sink", sink: &FanoutSink{sinks: []Sink{nil}}},
		{name: "typed nil internal sink", sink: &FanoutSink{sinks: []Sink{typedNil}}},
		{name: "mixed invalid internal sink", sink: &FanoutSink{sinks: []Sink{valid, nil}}},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/replace", func(t *testing.T) {
			err := tt.sink.Replace(context.Background(), testRead(t, 1))
			requireErrorIs(t, err, ErrInvalidFanoutSink)
		})
		t.Run(tt.name+"/apply", func(t *testing.T) {
			err := tt.sink.ApplyChange(context.Background(), createdChange(t, testKey("system", 1), 2))
			requireErrorIs(t, err, ErrInvalidFanoutSink)
		})
	}
}

type fanoutContextKey struct{}

type fanoutCall struct {
	ctx    context.Context
	read   storewatchapi.CollectionRead
	change objectstore.Change
}

type fanoutRecordingSink struct {
	name string

	replaceErr error
	applyErr   error

	sharedOrder *[]string
	localOrder  []string
	replace     []fanoutCall
	apply       []fanoutCall
}

func newFanoutRecordingSink(name string, order *[]string) *fanoutRecordingSink {
	return &fanoutRecordingSink{name: name, sharedOrder: order}
}

func (s *fanoutRecordingSink) Replace(ctx context.Context, read storewatchapi.CollectionRead) error {
	s.recordOrder("replace")
	s.replace = append(s.replace, fanoutCall{ctx: ctx, read: read})

	return s.replaceErr
}

func (s *fanoutRecordingSink) ApplyChange(ctx context.Context, change objectstore.Change) error {
	s.recordOrder("apply")
	s.apply = append(s.apply, fanoutCall{ctx: ctx, change: change})

	return s.applyErr
}

func (s *fanoutRecordingSink) recordOrder(operation string) {
	marker := s.name + ":" + operation
	s.localOrder = append(s.localOrder, marker)
	if s.sharedOrder != nil {
		*s.sharedOrder = append(*s.sharedOrder, marker)
	}
}

func (s *fanoutRecordingSink) replaceCount() int {
	return len(s.replace)
}

func (s *fanoutRecordingSink) applyCount() int {
	return len(s.apply)
}

func (s *fanoutRecordingSink) requireReplace(
	t testing.TB,
	wantCtx context.Context,
	wantRead storewatchapi.CollectionRead,
) {
	t.Helper()

	if len(s.replace) != 1 {
		t.Fatalf("%s replace calls = %d; want 1", s.name, len(s.replace))
	}
	got := s.replace[0]
	if got.ctx != wantCtx {
		t.Fatalf("%s replace context = %p; want %p", s.name, got.ctx, wantCtx)
	}
	if got.read.Revision() != wantRead.Revision() ||
		got.read.Collection() != wantRead.Collection() ||
		got.read.Len() != wantRead.Len() {
		t.Fatalf("%s replace read = %#v; want %#v", s.name, got.read, wantRead)
	}
}

func (s *fanoutRecordingSink) requireApplyChange(
	t testing.TB,
	wantCtx context.Context,
	wantChange objectstore.Change,
) {
	t.Helper()

	if len(s.apply) != 1 {
		t.Fatalf("%s apply calls = %d; want 1", s.name, len(s.apply))
	}
	got := s.apply[0]
	if got.ctx != wantCtx {
		t.Fatalf("%s apply context = %p; want %p", s.name, got.ctx, wantCtx)
	}
	if got.change.Kind != wantChange.Kind ||
		got.change.Revision != wantChange.Revision ||
		!got.change.Key.Equal(wantChange.Key) {
		t.Fatalf("%s apply change = %#v; want %#v", s.name, got.change, wantChange)
	}
}

func requireOrder(t testing.TB, got []string, want ...string) {
	t.Helper()

	if len(got) != len(want) {
		t.Fatalf("order = %v; want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order = %v; want %v", got, want)
		}
	}
}
