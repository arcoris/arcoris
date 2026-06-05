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

package diagnostic

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

// CauseRecord creates a diagnostic record that only carries a nested cause.
func CauseRecord[R ~string](cause error) Record[R] {
	return Record[R]{
		Cause: cause,
	}
}

// JoinRecord creates a diagnostic record from only sentinel and cause.
func JoinRecord[R ~string](err error, cause error) Record[R] {
	return Record[R]{
		Err:   err,
		Cause: cause,
	}
}
