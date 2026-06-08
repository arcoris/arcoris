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

package retrybudgetadmission

// Request is the admission adapter request for one retry-budget spend attempt.
//
// The request is intentionally empty because the retrybudget.RetryAdmitter
// contract represents exactly one atomic retry attempt. Dimensions such as
// tenant, priority, retry class, protocol, operation kind, or distributed budget
// identity belong to higher-level policy packages, not this generic adapter.
type Request struct{}
