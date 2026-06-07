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

package healthregistry

import (
	"fmt"
	"testing"

	"arcoris.dev/health"
)

var (
	benchmarkRegistry *Registry
	benchmarkSet      health.CheckSet
	benchmarkTargets  []health.Target
	benchmarkBool     bool
)

func BenchmarkBuilderRegister(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		builder := NewBuilder()
		for i := 0; i < 8; i++ {
			if err := builder.Register(health.TargetReady, checker(fmt.Sprintf("check_%d", i))); err != nil {
				b.Fatalf("Register() = %v", err)
			}
		}
	}
}

func BenchmarkBuilderBuild(b *testing.B) {
	builder := NewBuilder()
	for i := 0; i < 8; i++ {
		if err := builder.Register(health.TargetReady, checker(fmt.Sprintf("check_%d", i))); err != nil {
			b.Fatalf("Register() = %v", err)
		}
	}

	b.ReportAllocs()
	for b.Loop() {
		registry, err := builder.Build()
		if err != nil {
			b.Fatalf("Build() = %v", err)
		}
		benchmarkRegistry = registry
	}
}

func BenchmarkRegistryResolveChecks(b *testing.B) {
	registry := benchmarkReadyRegistry(b)

	b.ReportAllocs()
	for b.Loop() {
		set, err := registry.ResolveChecks(health.TargetReady)
		if err != nil {
			b.Fatalf("ResolveChecks() = %v", err)
		}
		benchmarkSet = set
	}
}

func BenchmarkRegistryHas(b *testing.B) {
	registry := benchmarkReadyRegistry(b)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkBool = registry.Has(health.TargetReady, "check_4")
	}
}

func BenchmarkRegistryTargets(b *testing.B) {
	registry := benchmarkReadyRegistry(b)

	b.ReportAllocs()
	for b.Loop() {
		benchmarkTargets = registry.Targets()
	}
}

func benchmarkReadyRegistry(b *testing.B) *Registry {
	b.Helper()

	builder := NewBuilder()
	for i := 0; i < 8; i++ {
		if err := builder.Register(health.TargetReady, checker(fmt.Sprintf("check_%d", i))); err != nil {
			b.Fatalf("Register() = %v", err)
		}
	}

	registry, err := builder.Build()
	if err != nil {
		b.Fatalf("Build() = %v", err)
	}

	return registry
}
