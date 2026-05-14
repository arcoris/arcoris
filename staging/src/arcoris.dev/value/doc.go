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

// Package value documents the ARCORIS logical value container module.
//
// The module contains small, dependency-free packages for representing logical
// value states such as presence, absence, success, and failure as ordinary Go
// values. These packages are intended for state and read-model boundaries where
// a value state must be stored, passed, or published.
//
// The module does not define application domain models, numeric constraints,
// validation frameworks, cache implementations, reflection helpers,
// serialization helpers, collection utilities, or broad common helpers.
package value
