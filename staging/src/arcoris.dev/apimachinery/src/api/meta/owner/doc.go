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

// Package owner defines owner reference metadata.
//
// Owner references model metadata links from dependent objects to owner
// objects. The package does not perform garbage collection, resolve
// references, or check object existence. At most one owner reference may be
// marked as controller.
//
// Nil and empty lists are both zero at the metadata value layer. Patch/apply
// field presence is represented by higher layers, not by this slice value.
package owner
