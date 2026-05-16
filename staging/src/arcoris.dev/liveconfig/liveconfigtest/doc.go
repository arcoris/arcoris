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

// Package liveconfigtest provides deterministic fixtures for tests built around
// package liveconfig and the snapshot source contracts it exposes.
//
// # Package scope
//
// liveconfigtest owns reusable test-only building blocks:
//
//   - Config, a small canonical live-configuration fixture with scalar fields,
//     mutable slices, and mutable maps;
//   - clone, equality, validation, and mutation helpers for Config;
//   - ControlledSource, a deterministic revisioned snapshot source for consumer
//     tests that should not depend on a full liveconfig.Holder;
//   - Loader, a scripted candidate loader for reload-loop tests;
//   - assertion helpers that understand snapshot revisions and Config equality.
//
// The package is intentionally boring. Helpers should make tests shorter and
// clearer, not introduce a second live-configuration model. Production packages
// must not depend on liveconfigtest.
//
// # Relationship to liveconfig
//
// Package liveconfig remains the owner of holder semantics: clone-before-
// normalize, normalize-before-validate, last-good preservation, ChangeReason,
// revision publication, and LastError behavior. liveconfigtest provides values
// and sources that make those contracts easy to exercise from adapter,
// integration, and external consumer tests.
//
// Internal tests inside package liveconfig should keep package-local helpers
// when they need unexported details. liveconfigtest is for tests that can work
// against exported contracts.
//
// # Non-goals
//
// liveconfigtest does not load files, parse JSON/YAML/TOML, read environment
// variables, watch files, watch Kubernetes ConfigMaps, call control planes,
// export metrics, notify subscribers, manage secrets, or simulate rollback
// history. Source-specific test doubles belong to the adapter packages that own
// those sources.
//
// # File ownership
//
//   - config.go owns the Config fixture and value-style variant methods;
//   - config_mutation.go owns mutation and invalid fixture builders;
//   - clone.go, equal.go, and validate.go own reusable Config functions;
//   - source.go owns ControlledSource;
//   - source_config.go owns Config-specific source convenience functions;
//   - source_contracts.go owns snapshot contract assertions;
//   - loader.go owns scripted candidate loading;
//   - assert*.go files own test assertions.
package liveconfigtest
