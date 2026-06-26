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

package objectcache

import "errors"

var (
	// ErrInvalidCache reports an unusable cache receiver or collection.
	ErrInvalidCache = errors.New("invalid object cache")
	// ErrInvalidOption reports an invalid cache construction option.
	ErrInvalidOption = errors.New("invalid object cache option")
	// ErrInvalidRead reports a malformed or cache-incompatible collection read.
	ErrInvalidRead = errors.New("invalid object cache read")
	// ErrInvalidChange reports a malformed or cache-incompatible committed change.
	ErrInvalidChange = errors.New("invalid object cache change")
	// ErrOutsideCollection reports a key or change outside the cache collection.
	ErrOutsideCollection = errors.New("object cache item outside collection")
	// ErrNotReady reports mutation or read-model use before the first Replace.
	ErrNotReady = errors.New("object cache not ready")
	// ErrStaleRead reports a replacement older than the current cache revision.
	ErrStaleRead = errors.New("stale object cache read")
	// ErrStaleChange reports a change that does not advance the cache revision.
	ErrStaleChange = errors.New("stale object cache change")
	// ErrDuplicateKey reports repeated keys inside one replacement read.
	ErrDuplicateKey = errors.New("duplicate object cache key")
	// ErrHistoryDisabled reports a historical read against a latest-only cache.
	ErrHistoryDisabled = errors.New("object cache history disabled")
	// ErrHistoryUnavailable reports a historical read outside retained knowledge.
	ErrHistoryUnavailable = errors.New("object cache history unavailable")
	// ErrFutureRevision reports a historical read beyond the cache boundary.
	ErrFutureRevision = errors.New("object cache future revision")
)

// errorWith preserves broad cache sentinels while keeping lower-level causes
// visible to errors.Is and diagnostics.
func errorWith(sentinel error, causes ...error) error {
	errs := make([]error, 0, len(causes)+1)
	errs = append(errs, sentinel)
	for _, cause := range causes {
		if cause != nil {
			errs = append(errs, cause)
		}
	}

	return errors.Join(errs...)
}
