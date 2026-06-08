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

import (
	"errors"
	"maps"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"arcoris.dev/snapshot"
)

type testConfig struct {
	Name  string
	Limit int
	Tags  []string
}

type testClock struct {
	now time.Time
}

type mutableClock struct {
	now time.Time
}

type mutableConfig struct {
	Name   string
	Tags   []string
	Labels map[string]string
	Nested *nestedConfig
}

type nestedConfig struct {
	Values []string
}

func (c testClock) Now() time.Time {
	return c.now
}

func (c testClock) Since(t time.Time) time.Duration {
	return c.now.Sub(t)
}

func (c testClock) Until(t time.Time) time.Duration {
	return t.Sub(c.now)
}

func (c *mutableClock) Now() time.Time {
	return c.now
}

func (c *mutableClock) Since(t time.Time) time.Duration {
	return c.now.Sub(t)
}

func (c *mutableClock) Until(t time.Time) time.Duration {
	return t.Sub(c.now)
}

func newTestClock() testClock {
	return testClock{now: time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)}
}

func cloneTestConfig(cfg testConfig) testConfig {
	cfg.Tags = append([]string(nil), cfg.Tags...)
	return cfg
}

func validTestConfig(cfg testConfig) error {
	if cfg.Limit < 0 {
		return errors.New("limit must be non-negative")
	}
	return nil
}

func equalTestConfig(a, b testConfig) bool {
	if a.Name != b.Name || a.Limit != b.Limit || len(a.Tags) != len(b.Tags) {
		return false
	}
	for i := range a.Tags {
		if a.Tags[i] != b.Tags[i] {
			return false
		}
	}
	return true
}

func cloneMutableConfig(v mutableConfig) mutableConfig {
	out := mutableConfig{
		Name:   v.Name,
		Tags:   slices.Clone(v.Tags),
		Labels: maps.Clone(v.Labels),
	}
	if v.Nested != nil {
		out.Nested = &nestedConfig{
			Values: slices.Clone(v.Nested.Values),
		}
	}
	return out
}

func equalMutableConfig(a, b mutableConfig) bool {
	return a.Name == b.Name &&
		slices.Equal(a.Tags, b.Tags) &&
		maps.Equal(a.Labels, b.Labels) &&
		equalNestedConfig(a.Nested, b.Nested)
}

func equalNestedConfig(a, b *nestedConfig) bool {
	switch {
	case a == nil || b == nil:
		return a == b
	default:
		return slices.Equal(a.Values, b.Values)
	}
}

func newTestHolder(t *testing.T, initial testConfig, opts ...Option[testConfig]) *Holder[testConfig] {
	t.Helper()

	base := []Option[testConfig]{
		WithClock[testConfig](newTestClock()),
		WithClone(cloneTestConfig),
		WithValidator(validTestConfig),
	}
	base = append(base, opts...)

	h, err := New(initial, base...)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}
	return h
}

var (
	benchmarkSnapshotSink snapshot.Snapshot[testConfig]
	benchmarkStampedSink  snapshot.Stamped[testConfig]
	benchmarkRevisionSink snapshot.Revision
	benchmarkChangeSink   Change[testConfig]
	benchmarkErrorSink    error
	benchmarkSinkMu       sync.Mutex
)

type benchmarkError struct{}

func (benchmarkError) Error() string {
	return "benchmark error"
}

func newBenchmarkHolder(b *testing.B, opts ...Option[testConfig]) *Holder[testConfig] {
	b.Helper()

	base := []Option[testConfig]{
		WithClock[testConfig](newTestClock()),
		WithClone(cloneTestConfig),
		WithValidator(validTestConfig),
	}
	base = append(base, opts...)

	h, err := New(testConfig{Name: "initial", Limit: 1}, base...)
	if err != nil {
		b.Fatalf("New() error = %v", err)
	}
	return h
}

func benchmarkConfigWithTags(n int) testConfig {
	tags := make([]string, n)
	for i := range tags {
		tags[i] = "tag"
	}

	return testConfig{
		Name:  "tagged",
		Limit: n,
		Tags:  tags,
	}
}

type snapshotLocal struct {
	snapshot snapshot.Snapshot[testConfig]
}

func (l snapshotLocal) keep() {
	benchmarkSinkMu.Lock()
	benchmarkSnapshotSink = l.snapshot
	benchmarkSinkMu.Unlock()
}

type stampedLocal struct {
	stamped snapshot.Stamped[testConfig]
}

func (l stampedLocal) keep() {
	benchmarkSinkMu.Lock()
	benchmarkStampedSink = l.stamped
	benchmarkSinkMu.Unlock()
}

type revisionLocal struct {
	revision snapshot.Revision
}

func (l revisionLocal) keep() {
	benchmarkSinkMu.Lock()
	benchmarkRevisionSink = l.revision
	benchmarkSinkMu.Unlock()
}

type errorLocal struct {
	err error
}

func (l errorLocal) keep() {
	benchmarkSinkMu.Lock()
	benchmarkErrorSink = l.err
	benchmarkSinkMu.Unlock()
}

type changeLocal struct {
	change Change[testConfig]
	err    error
}

func (l changeLocal) keep() {
	benchmarkSinkMu.Lock()
	benchmarkChangeSink = l.change
	benchmarkErrorSink = l.err
	benchmarkSinkMu.Unlock()
}

func startBenchmarkApplier(h *Holder[testConfig]) func() {
	done := make(chan struct{})
	var next atomic.Int64
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			default:
				limit := int(next.Add(1))
				_, _ = h.Apply(testConfig{Name: "next", Limit: limit})
			}
		}
	}()

	return func() {
		close(done)
		wg.Wait()
	}
}

func runConcurrent(n int, fn func(int)) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			<-start
			fn(i)
		}()
	}
	close(start)
	wg.Wait()
}

func runConcurrentErrors(t *testing.T, n int, fn func(int) error) {
	t.Helper()

	errs := make(chan error, n)
	runConcurrent(n, func(i int) {
		errs <- fn(i)
	})
	close(errs)

	for err := range errs {
		if err != nil {
			t.Fatalf("concurrent operation error = %v", err)
		}
	}
}

func requirePanic(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("panic = nil, want panic")
		}
	}()

	fn()
}

func assertTestSnapshot(t *testing.T, got, want snapshot.Snapshot[testConfig]) {
	t.Helper()
	if got.Revision != want.Revision || !equalTestConfig(got.Value, want.Value) {
		t.Fatalf("Snapshot() = %#v, want %#v", got, want)
	}
}

func assertMutableConfig(t *testing.T, got, want mutableConfig) {
	t.Helper()
	if !equalMutableConfig(got, want) {
		t.Fatalf("mutable config = %#v, want %#v", got, want)
	}
}
