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

package codecselection

import "fmt"

const (
	pathRoot        = "codecselection"
	pathContentType = pathRoot + ".contentType"
	pathParameter   = pathRoot + ".parameter"
	pathParameters  = pathRoot + ".parameters"
	pathPreference  = pathRoot + ".preference"
	pathPreferences = pathRoot + ".preferences"
	pathDirection   = pathRoot + ".direction"
	pathTransport   = pathRoot + ".transport"
	pathWeight      = pathRoot + ".weight"
)

// decodeBindingPath returns the stable path for one decode binding.
func decodeBindingPath(index int) string {
	return fmt.Sprintf("%s.decodeBindings[%d]", pathRoot, index)
}

// encodeBindingPath returns the stable path for one encode binding.
func encodeBindingPath(index int) string {
	return fmt.Sprintf("%s.encodeBindings[%d]", pathRoot, index)
}

// bindingContentTypePath returns the path for one binding content type.
func bindingContentTypePath(bindingPath string) string {
	return bindingPath + ".contentType"
}

// bindingTargetPath returns the path for one binding target.
func bindingTargetPath(bindingPath string) string {
	return bindingPath + ".target"
}

// bindingTransportPath returns the path for one binding transport.
func bindingTransportPath(bindingPath string) string {
	return bindingPath + ".transport"
}

// bindingEntryIDPath returns the path for one binding entry ID.
func bindingEntryIDPath(bindingPath string) string {
	return bindingPath + ".entryID"
}

// runtimeDecodePath returns a stable path for runtime decode selection.
func runtimeDecodePath(target string) string {
	return pathRoot + ".decode." + target
}

// runtimeEncodePath returns a stable path for runtime encode selection.
func runtimeEncodePath(target string) string {
	return pathRoot + ".encode." + target
}
