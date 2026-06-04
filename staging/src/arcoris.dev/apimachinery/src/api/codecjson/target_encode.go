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

package codecjson

import (
	"bytes"
	"io"
)

// encodeTargetBytes owns the common node-to-byte encode pipeline.
//
// Public byte encoders convert their semantic target into a jsonNode first,
// then enter this shared writer path so Pretty and EscapeHTML behavior stay
// identical across value, object, and object ownership documents.
func encodeTargetBytes(node jsonNode, config resolvedEncodeConfig) ([]byte, error) {
	var buffer bytes.Buffer
	if err := encodeTargetTo(&buffer, node, config); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// encodeTargetTo owns the common node-to-stream encode pipeline.
//
// The stream writer emits exactly one JSON document and intentionally does not
// append a trailing newline. Callers that need framing own that transport layer.
func encodeTargetTo(w io.Writer, node jsonNode, config resolvedEncodeConfig) error {
	return encodeJSONDocument(w, node, nodeEncodeConfig{
		pretty:          config.pretty,
		indent:          config.indent,
		finalNewline:    config.finalNewline,
		escapeHTML:      config.escapeHTML,
		maxDepth:        config.maxDepth,
		maxOutputBytes:  config.maxOutputBytes,
		maxNumberDigits: config.maxNumberDigits,
	})
}
