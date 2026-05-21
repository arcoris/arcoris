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


package capacity

// noCopy marks owner values as non-copyable for go vet's copylocks checker.
//
// Ledger contains synchronization state. Reservation contains release ownership
// state tied to one ledger. Copying either value after first use can split
// ownership and corrupt capacity accounting.
type noCopy struct{}

// Lock is a marker method recognized by go vet's copylocks analyzer.
func (*noCopy) Lock() {}

// Unlock is a marker method recognized by go vet's copylocks analyzer.
func (*noCopy) Unlock() {}
