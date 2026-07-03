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

package objectreflectorsink

import (
	"context"
	"errors"
	"fmt"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/runtime/objectreflector"
)

var _ objectreflector.Sink = (*Fanout)(nil)

func TestNewFanoutValidation(t *testing.T) {
	tests := []struct {
		name    string
		sinks   []objectreflector.Sink
		wantErr error
	}{
		{name: "no sinks", wantErr: ErrNoSinks},
		{name: "nil sink", sinks: []objectreflector.Sink{nil}, wantErr: objectreflector.ErrNilSink},
		{name: "typed nil sink", sinks: []objectreflector.Sink{(*recordingSink)(nil)}, wantErr: objectreflector.ErrNilSink},
		{name: "one sink", sinks: []objectreflector.Sink{newRecordingSink("first", nil)}},
		{
			name: "multiple sinks",
			sinks: []objectreflector.Sink{
				newRecordingSink("first", nil),
				newRecordingSink("second", nil),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink, err := NewFanout(tt.sinks...)
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

func TestNewFanoutCopiesSinkList(t *testing.T) {
	original := newRecordingSink("original", nil)
	replacement := newRecordingSink("replacement", nil)
	input := []objectreflector.Sink{original}
	sink, err := NewFanout(input...)
	requireNoError(t, err)

	input[0] = replacement
	requireNoError(t, sink.ApplyChange(context.Background(), createdChange(t, testKey(1), 2)))

	requireOrder(t, original.localOrder, "original:apply")
	requireOrder(t, replacement.localOrder)
}

func TestFanoutReplaceOrderAndPassthrough(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey{}, "ctx")
	read := testRead(t, 7, listItem(testKey(1), 7, "listed"))
	order := make([]string, 0, 2)
	first := newRecordingSink("first", &order)
	second := newRecordingSink("second", &order)
	sink, err := NewFanout(first, second)
	requireNoError(t, err)

	err = sink.Replace(ctx, read)

	requireNoError(t, err)
	requireOrder(t, order, "first:replace", "second:replace")
	first.requireReplace(t, ctx, read)
	second.requireReplace(t, ctx, read)
}

func TestFanoutApplyChangeOrderAndPassthrough(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey{}, "ctx")
	change := createdChange(t, testKey(1), 2)
	order := make([]string, 0, 2)
	first := newRecordingSink("first", &order)
	second := newRecordingSink("second", &order)
	sink, err := NewFanout(first, second)
	requireNoError(t, err)

	err = sink.ApplyChange(ctx, change)

	requireNoError(t, err)
	requireOrder(t, order, "first:apply", "second:apply")
	first.requireApplyChange(t, ctx, change)
	second.requireApplyChange(t, ctx, change)
}

func TestFanoutReplaceStopsOnFirstError(t *testing.T) {
	wantErr := errors.New("replace failed")
	order := make([]string, 0, 2)
	first := newRecordingSink("first", &order)
	first.replaceErr = wantErr
	second := newRecordingSink("second", &order)
	sink, err := NewFanout(first, second)
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

func TestFanoutApplyChangeStopsOnFirstError(t *testing.T) {
	wantErr := errors.New("apply failed")
	order := make([]string, 0, 2)
	first := newRecordingSink("first", &order)
	first.applyErr = wantErr
	second := newRecordingSink("second", &order)
	sink, err := NewFanout(first, second)
	requireNoError(t, err)

	err = sink.ApplyChange(context.Background(), createdChange(t, testKey(1), 2))

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	requireOrder(t, order, "first:apply")
	if second.applyCount() != 0 {
		t.Fatalf("second apply calls = %d; want 0", second.applyCount())
	}
}

func TestFanoutInvalidReceiver(t *testing.T) {
	valid := newRecordingSink("valid", nil)
	var typedNil objectreflector.Sink = (*recordingSink)(nil)
	tests := []struct {
		name string
		sink *Fanout
	}{
		{name: "nil", sink: nil},
		{name: "empty", sink: &Fanout{}},
		{name: "nil internal sink", sink: &Fanout{sinks: []objectreflector.Sink{nil}}},
		{name: "typed nil internal sink", sink: &Fanout{sinks: []objectreflector.Sink{typedNil}}},
		{name: "mixed invalid internal sink", sink: &Fanout{sinks: []objectreflector.Sink{valid, nil}}},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/replace", func(t *testing.T) {
			err := tt.sink.Replace(context.Background(), testRead(t, 1))
			requireErrorIs(t, err, ErrInvalidFanout)
		})
		t.Run(tt.name+"/apply", func(t *testing.T) {
			err := tt.sink.ApplyChange(context.Background(), createdChange(t, testKey(1), 2))
			requireErrorIs(t, err, ErrInvalidFanout)
		})
	}
}

type contextKey struct{}

type callRecord struct {
	ctx    context.Context
	read   storewatchapi.CollectionRead
	change objectstore.Change
}

type recordingSink struct {
	name string

	replaceErr error
	applyErr   error

	sharedOrder *[]string
	localOrder  []string
	replace     []callRecord
	apply       []callRecord
}

func newRecordingSink(name string, order *[]string) *recordingSink {
	return &recordingSink{name: name, sharedOrder: order}
}

func (s *recordingSink) Replace(ctx context.Context, read storewatchapi.CollectionRead) error {
	s.recordOrder("replace")
	s.replace = append(s.replace, callRecord{ctx: ctx, read: read})

	return s.replaceErr
}

func (s *recordingSink) ApplyChange(ctx context.Context, change objectstore.Change) error {
	s.recordOrder("apply")
	s.apply = append(s.apply, callRecord{ctx: ctx, change: change})

	return s.applyErr
}

func (s *recordingSink) recordOrder(operation string) {
	marker := s.name + ":" + operation
	s.localOrder = append(s.localOrder, marker)
	if s.sharedOrder != nil {
		*s.sharedOrder = append(*s.sharedOrder, marker)
	}
}

func (s *recordingSink) replaceCount() int {
	return len(s.replace)
}

func (s *recordingSink) applyCount() int {
	return len(s.apply)
}

func (s *recordingSink) requireReplace(
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

func (s *recordingSink) requireApplyChange(
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

func testCollection() objectstore.ListRequest {
	return objectstore.ListRequest{Resource: testResource(), Scope: objectstore.AllNamespaces()}
}

func testResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "units",
	}
}

func testKey(index int) objectstore.Key {
	return objectstore.MustKey(testResource(), metaidentity.ObjectName{
		Namespace: "system",
		Name:      metaidentity.Name(fmt.Sprintf("unit-%d", index)),
	})
}

func testState(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.State {
	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   key.Resource.Group,
				Version: key.Resource.Version,
				Kind:    "Unit",
			}),
			meta.ObjectMeta{Name: key.Object.Name, Namespace: key.Object.Namespace},
			value.StringValue(desired),
			value.StringValue("observed-"+desired),
		),
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
	}
}

func testRead(t testing.TB, revision objectstore.Revision, items ...objectstore.ListItem) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(testCollection(), objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

func listItem(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.ListItem {
	return objectstore.ListItem{Key: key, State: testState(key, revision, desired)}
}

func createdChange(t testing.TB, key objectstore.Key, revision objectstore.Revision) objectstore.Change {
	t.Helper()

	change, err := objectstore.NewCreatedChange(key, testState(key, revision, "created"))
	requireNoError(t, err)

	return change
}

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
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
