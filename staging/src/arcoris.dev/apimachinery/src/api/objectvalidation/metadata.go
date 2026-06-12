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

package objectvalidation

import (
	"errors"

	apiobject "arcoris.dev/apimachinery/api/object"
)

// validateMetadata delegates object envelope metadata validation.
//
// The nested api/object error path is preserved when available so callers can
// still identify object.typeMeta and object.metadata failures precisely.
func validateMetadata[D any, O any](obj apiobject.Object[D, O]) error {
	if err := obj.ValidateMeta(); err != nil {
		path, reason, detail := metadataDiagnostic(err)

		return nested(
			path,
			ErrInvalidMetadata,
			reason,
			detail,
			err,
		)
	}

	return nil
}

// metadataDiagnostic maps api/object metadata diagnostics into this package's
// object contract taxonomy while preserving the original error as the cause.
func metadataDiagnostic(err error) (string, ErrorReason, string) {
	path := pathObject
	reason := ErrorReasonInvalidMetadata
	detail := "object metadata is invalid"

	var objectErr *apiobject.Error
	if !errors.As(err, &objectErr) {
		return path, reason, detail
	}

	if objectErr.Path != "" {
		path = objectErr.Path
	}

	switch objectErr.Reason {
	case apiobject.ErrorReasonInvalidTypeMeta:
		return path, ErrorReasonInvalidTypeMeta, "type metadata is invalid"
	case apiobject.ErrorReasonInvalidObjectMeta:
		return path, ErrorReasonInvalidObjectMeta, "object metadata is invalid"
	default:
		return path, reason, detail
	}
}
