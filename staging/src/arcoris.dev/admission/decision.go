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

package admission

// Decision is the non-generic semantic core of an admission result.
//
// Decision describes what happened, why it happened, and which side-effect class
// was produced. Typed grants and metadata live in Result.
type Decision struct {
	// Outcome classifies what happened to the admission attempt.
	Outcome Outcome

	// Reason explains the outcome using a stable open-world reason code.
	Reason Reason

	// Effect classifies the side-effect semantics of the attempt.
	Effect Effect
}
