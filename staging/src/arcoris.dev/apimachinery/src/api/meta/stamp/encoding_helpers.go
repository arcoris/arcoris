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

package stamp

import (
	"arcoris.dev/apimachinery/api/internal/diagnostic"
	"arcoris.dev/apimachinery/api/meta/internal/metaencoding"
)

// marshalText validates a stamp scalar before exposing it as text.
func marshalText(value string, validate func() error) ([]byte, error) {
	return metaencoding.MarshalText(value, validate)
}

// marshalJSONString validates a stamp scalar before encoding it as one JSON string.
func marshalJSONString(value string, validate func() error) ([]byte, error) {
	return metaencoding.MarshalJSONString(value, validate)
}

// unmarshalJSONString decodes exactly one JSON string and rejects null.
func unmarshalJSONString(path string, data []byte) (string, error) {
	value, isNull, err := metaencoding.DecodeJSONString(data)
	if err != nil {
		return "", &Error{
			Record: diagnostic.WrapRecord(
				path,
				ErrInvalidJSON,
				ErrorReasonInvalidJSON,
				"expected JSON string",
				err,
			),
		}
	}
	if isNull {
		return "", &Error{
			Record: diagnostic.NewRecord(
				path,
				ErrInvalidJSON,
				ErrorReasonInvalidJSON,
				"expected JSON string, got null",
			),
		}
	}
	return value, nil
}

// nilReceiver reports an attempted scalar decode into a nil receiver.
func nilReceiver(path string) error {
	return &Error{
		Record: diagnostic.NewRecord(
			path,
			ErrNilReceiver,
			ErrorReasonNilReceiver,
			"receiver must be non-nil",
		),
	}
}
