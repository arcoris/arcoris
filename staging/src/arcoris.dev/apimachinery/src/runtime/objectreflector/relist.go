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

package objectreflector

import (
	"errors"
	"fmt"

	"arcoris.dev/apimachinery/api/objectwatch"
)

// isRelistRequired reports whether Run should begin a new list-watch cycle.
func isRelistRequired(err error) bool {
	return errors.Is(err, objectwatch.ErrHistoryUnavailable) ||
		errors.Is(err, objectwatch.ErrContinuityLost) ||
		errors.Is(err, errRelistRequired)
}

// relistRequiredError keeps restart control flow private to this package while
// preserving a useful diagnostic for failed tests.
func relistRequiredError(reason objectwatch.RestartReason) error {
	return errors.Join(errRelistRequired, fmt.Errorf("restart required: %s", reason.String()))
}
