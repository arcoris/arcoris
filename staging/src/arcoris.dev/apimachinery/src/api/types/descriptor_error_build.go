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

package types

import (
	"fmt"

	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// descriptorError builds a path-aware validation error without reason detail.
//
// Use this helper for broad invariant failures where the sentinel and
// descriptor path are already precise enough for callers and tests.
func descriptorError(path string, err error) error {
	return &DescriptorError{
		Record: diagnostic.NewRecord(path, err, DescriptorErrorReason(""), ""),
	}
}

// descriptorErrorf builds a path-aware validation error with formatted detail.
//
// The helper keeps descriptor validation call sites compact while preserving
// the same structured DescriptorError shape as simpler validation failures.
func descriptorErrorf(path string, err error, reason DescriptorErrorReason, format string, args ...any) error {
	detail := ""

	if format != "" {
		detail = fmt.Sprintf(format, args...)
	}

	return &DescriptorError{
		Record: diagnostic.NewRecord(path, err, reason, detail),
	}
}
