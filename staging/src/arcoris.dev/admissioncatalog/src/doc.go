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

// Package admissioncatalog provides owner-created catalogs of stable admission
// metadata.
//
// The package is the concrete metadata-catalog counterpart to the open-world
// identifier types in package arcoris.dev/admission. Package admission owns
// Result, Decision, Outcome, Effect, Reason, ComponentID, and ComponentKind.
// Package admissioncatalog owns descriptors, registries, and aggregate catalogs
// used by documentation, tests, configuration validation, and future admission
// chain validation foundations.
//
// Catalogs are not runtime admission chains, live admitter registries, service
// discovery, metrics registries, tracing registries, global extension
// registries, or process-wide singleton state. They perform no admission
// attempts and store no Admitter values.
package admissioncatalog
