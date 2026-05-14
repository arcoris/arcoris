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

import "net/http"

const (
	headerContentType         = "Content-Type"
	headerCacheControl        = "Cache-Control"
	headerXContentTypeOptions = "X-Content-Type-Options"
)

const (
	headerValueNoStore = "no-store"
	headerValueNoSniff = "nosniff"
)

// writeCommonHeaders writes headers shared by all health HTTP responses.
//
// Health responses are deliberately non-cacheable and content-sniffing is
// disabled so endpoint bodies cannot be reinterpreted as another media type.
func writeCommonHeaders(w http.ResponseWriter, format Format) {
	w.Header().Set(headerContentType, format.contentType())
	w.Header().Set(headerCacheControl, headerValueNoStore)
	w.Header().Set(headerXContentTypeOptions, headerValueNoSniff)
}

// suppressBody reports whether the adapter must suppress the response body.
//
// GET and HEAD share status and header behavior, but HEAD must not emit a body.
func suppressBody(r *http.Request) bool {
	return r != nil && r.Method == http.MethodHead
}
