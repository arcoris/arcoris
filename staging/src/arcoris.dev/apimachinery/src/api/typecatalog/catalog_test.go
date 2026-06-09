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

package typecatalog

import (
	"fmt"
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/types"
)

func TestCatalogZeroValueUsable(t *testing.T) {
	var catalog Catalog

	requireNoError(t, catalog.Register(types.Define("example.Name", types.String().MinBytes(1))))

	requireDefinition(t, &catalog, "example.Name")
}

func TestCatalogConcurrentAccessRaceFree(t *testing.T) {
	var catalog Catalog
	requireNoError(t, catalog.Register(types.Define("example.Name", types.String())))

	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				catalog.Resolve("example.Name")
				catalog.Names()
				catalog.Definitions()
			}
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for j := 0; j < 100; j++ {
			name := types.TypeName(fmt.Sprintf("example.Generated%d", j))
			_ = catalog.Register(types.Define(name, types.String()))
		}
	}()

	wg.Wait()
}
