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

// Package diagnostic contains the local structured diagnostic record helper for
// the JSON codec implementation.
package diagnostic

import (
	"errors"
	"strings"
)

// Record stores the common structured diagnostic shape used by jsoncodec.
type Record[R ~string] struct {
	// Path identifies the JSON document location that failed.
	Path string
	// Err is the broad sentinel used with errors.Is.
	Err error
	// Reason gives stable machine-readable detail within Err.
	Reason R
	// Detail gives human-readable diagnostic context.
	Detail string
	// Cause preserves a nested lower-level failure.
	Cause error
}

// NewRecord creates a direct diagnostic record without a nested cause.
func NewRecord[R ~string](path string, err error, reason R, detail string) Record[R] {
	return Record[R]{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: detail,
	}
}

// WrapRecord creates a diagnostic record that preserves a nested cause.
func WrapRecord[R ~string](
	path string,
	err error,
	reason R,
	detail string,
	cause error,
) Record[R] {
	record := NewRecord(path, err, reason, detail)
	record.Cause = cause

	return record
}

// Format builds the common ARCORIS diagnostic string shape for r.
func (r Record[R]) Format(prefix string) string {
	parts := []string{prefix}
	if r.Path != "" {
		parts = append(parts, r.Path)
	}
	if r.Err != nil {
		parts = append(parts, r.Err.Error())
	}
	if r.Reason != "" {
		parts = append(parts, string(r.Reason))
	}
	if r.Detail != "" {
		parts = append(parts, r.Detail)
	}

	return strings.Join(parts, ": ")
}

// Unwrap preserves the broad sentinel and nested cause identities.
func (r Record[R]) Unwrap() error {
	switch {
	case r.Err != nil && r.Cause != nil:
		return errors.Join(r.Err, r.Cause)
	case r.Err != nil:
		return r.Err
	default:
		return r.Cause
	}
}
