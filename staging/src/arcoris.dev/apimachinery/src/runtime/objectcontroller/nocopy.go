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

package objectcontroller

// noCopy marks controller owner values as non-copyable for go vet's copylocks
// check.
//
// noCopy is a static-analysis marker only. It does not provide runtime
// synchronization and must remain unexported.
type noCopy struct{}

// Lock marks noCopy as a lock-like value for go vet.
func (*noCopy) Lock() {}

// Unlock marks noCopy as a lock-like value for go vet.
func (*noCopy) Unlock() {}
