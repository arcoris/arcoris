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

// Package result provides a small generic container for storing an operation
// outcome as a value.
//
// Result represents either OK(value) or Err(error). It is intended for cases
// where an outcome must be stored, sent through a channel, published in a
// snapshot, or used in test data.
//
// The package is not intended to replace Go's ordinary (T, error) return
// convention for normal fallible functions. Prefer direct error returns unless
// the outcome itself must be represented as data.
package result
