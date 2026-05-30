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

// Package stamp defines system-assigned metadata stamps.
//
// ResourceVersion is an opaque API/storage concurrency and change token.
// Generation is a desired-state generation stamp. Timestamp models metadata
// time values. Deletion models deletion request metadata.
//
// The package does not depend on snapshot and does not implement storage,
// watches, ordering, deletion execution, or finalization.
package stamp
