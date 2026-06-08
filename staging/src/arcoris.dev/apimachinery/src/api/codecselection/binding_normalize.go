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

import (
	"arcoris.dev/apimachinery/api/codec"
	"arcoris.dev/apimachinery/api/codecregistry"
)

// normalizeDecodeBindingAt validates one decode binding.
func normalizeDecodeBindingAt(path string, config Config, binding DecodeBinding) (bindingRecord, error) {
	return normalizeBindingAt(
		path,
		config,
		binding.ContentType,
		binding.Target,
		binding.Transport,
		binding.EntryID,
	)
}

// normalizeEncodeBindingAt validates one encode binding.
func normalizeEncodeBindingAt(path string, config Config, binding EncodeBinding) (bindingRecord, error) {
	return normalizeBindingAt(
		path,
		config,
		binding.ContentType,
		binding.Target,
		binding.Transport,
		binding.EntryID,
	)
}

// normalizeBindingAt validates one binding and its referenced registry entry.
func normalizeBindingAt(
	path string,
	config Config,
	contentType ContentType,
	target codec.Target,
	transport Transport,
	entryID codecregistry.EntryID,
) (bindingRecord, error) {
	normalizedContentType, err := normalizeContentTypeAt(bindingContentTypePath(path), contentType)
	if err != nil {
		return bindingRecord{}, wrapAt(
			bindingContentTypePath(path),
			ErrInvalidBinding,
			ErrorReasonInvalidBinding,
			"binding content type is invalid",
			err,
		)
	}

	normalizedTarget, err := normalizeTargetAt(bindingTargetPath(path), target)
	if err != nil {
		return bindingRecord{}, err
	}

	if err := validateTransportAt(bindingTransportPath(path), transport); err != nil {
		return bindingRecord{}, err
	}

	normalizedEntryID, err := normalizeEntryIDAt(bindingEntryIDPath(path), entryID)
	if err != nil {
		return bindingRecord{}, err
	}

	entry, err := lookupEntryAt(bindingEntryIDPath(path), config, normalizedEntryID)
	if err != nil {
		return bindingRecord{}, err
	}
	if err := validateEntryBindingAt(path, entry, normalizedContentType, normalizedTarget, transport); err != nil {
		return bindingRecord{}, err
	}

	return bindingRecord{
		contentType: normalizedContentType,
		target:      normalizedTarget,
		transport:   transport,
		entryID:     normalizedEntryID,
		entry:       entry,
	}, nil
}
