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

package fieldownership

import (
	"fmt"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func BenchmarkNewStateSmall(b *testing.B) {
	entries := []Entry{
		entry("user-cli", imagePath(), replicasPath()),
		entry("autoscaler", replicasPath()),
		entry("health-controller", readyStatusPath()),
	}

	b.ReportAllocs()
	for b.Loop() {
		if _, err := NewState(entries...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewStateManyOwners(b *testing.B) {
	entries := benchmarkEntries(100, 4)

	b.ReportAllocs()
	for b.Loop() {
		if _, err := NewState(entries...); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConflictsSmall(b *testing.B) {
	state := baseState()
	attempted := set(replicasPath())
	owner := owner("other")

	b.ReportAllocs()
	for b.Loop() {
		if _, err := state.Conflicts(owner, attempted); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConflictsManyOwners(b *testing.B) {
	state := MustState(benchmarkEntries(100, 2)...)
	attempted := benchmarkFieldSet(0, 10)
	owner := owner("actor")

	b.ReportAllocs()
	for b.Loop() {
		if _, err := state.Conflicts(owner, attempted); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkConflictsManyPaths(b *testing.B) {
	state := MustState(benchmarkEntries(10, 100)...)
	attempted := benchmarkFieldSet(0, 50)
	owner := owner("actor")

	b.ReportAllocs()
	for b.Loop() {
		if _, err := state.Conflicts(owner, attempted); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRemoveOverlappingFieldsFromOthers(b *testing.B) {
	state := MustState(benchmarkEntries(50, 20)...)
	fields := benchmarkFieldSet(0, 10)
	owner := owner("owner-000")

	b.ReportAllocs()
	for b.Loop() {
		if _, err := state.RemoveOverlappingFieldsFromOthers(owner, fields); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkFieldsUnion(b *testing.B) {
	state := MustState(benchmarkEntries(100, 10)...)

	b.ReportAllocs()
	for b.Loop() {
		_ = state.Fields()
	}
}

func benchmarkEntries(ownerCount int, pathsPerOwner int) []Entry {
	entries := make([]Entry, 0, ownerCount)
	for ownerIndex := 0; ownerIndex < ownerCount; ownerIndex++ {
		entries = append(entries, MustEntry(
			owner(fmt.Sprintf("owner-%03d", ownerIndex)),
			benchmarkFieldSet(ownerIndex*pathsPerOwner, pathsPerOwner),
		))
	}

	return entries
}

func benchmarkFieldSet(start int, count int) fieldpath.Set {
	paths := make([]fieldpath.Path, 0, count)
	for i := 0; i < count; i++ {
		paths = append(paths, path(fmt.Sprintf("$.spec.field%d", start+i)))
	}

	return set(paths...)
}
