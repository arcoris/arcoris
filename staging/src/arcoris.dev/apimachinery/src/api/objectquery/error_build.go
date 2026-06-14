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

package objectquery

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// errorf builds a structured query validation error.
func errorf(path string, err error, reason ErrorReason, format string, args ...any) error {
	return &Error{
		Record: diagnostic.NewRecord(path, err, reason, fmt.Sprintf(format, args...)),
	}
}

// wrapf preserves a lower validation cause under a query diagnostic boundary.
func wrapf(path string, err error, reason ErrorReason, cause error, format string, args ...any) error {
	return &Error{
		Record: diagnostic.WrapRecord(path, err, reason, fmt.Sprintf(format, args...), cause),
	}
}
