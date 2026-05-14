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

package healthhttp

import "errors"

var (
	// ErrNilEvaluator identifies a nil health evaluator passed to NewHandler or
	// install helpers.
	//
	// healthhttp owns handler adaptation only; without a real evaluator it has no
	// core health execution boundary to delegate to.
	ErrNilEvaluator = errors.New("healthhttp: nil evaluator")
)
