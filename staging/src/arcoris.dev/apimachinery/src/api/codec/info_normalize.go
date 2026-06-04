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

package codec

import "slices"

// Normalize canonicalizes i while keeping duplicate declarations invalid.
//
// Normalize returns detached slices sorted in deterministic order. It accepts
// normalizable Format and MediaType text but still rejects duplicate
// declarations after canonicalization.
func (i Info) Normalize() (Info, error) {
	format, err := i.Format.Normalize()
	if err != nil {
		return Info{}, WrapAt(
			pathCodecInfo+".format",
			ErrInvalidInfo,
			ErrorReasonInvalidFormat,
			"codec info format is invalid",
			err,
		)
	}

	mediaTypes, err := normalizeMediaTypes(i.MediaTypes)
	if err != nil {
		return Info{}, err
	}
	targets, err := normalizeTargets(i.Targets)
	if err != nil {
		return Info{}, err
	}

	return Info{
		Format:     format,
		MediaTypes: mediaTypes,
		Targets:    targets,
	}, nil
}

// NormalizeInfo canonicalizes info and returns detached metadata slices.
//
// The function form mirrors nearby API packages and is convenient when callers
// want a package-level normalization hook.
func NormalizeInfo(info Info) (Info, error) {
	return info.Normalize()
}

// normalizeMediaTypes canonicalizes, validates, sorts, and detaches media types.
//
// The helper performs duplicate checks after normalization so values such as
// "APPLICATION/JSON" and "application/json" cannot collapse silently.
func normalizeMediaTypes(mediaTypes []MediaType) ([]MediaType, error) {
	if len(mediaTypes) == 0 {
		return nil, ErrorAt(
			pathCodecInfo+".mediaTypes",
			ErrInvalidInfo,
			ErrorReasonInvalidMediaType,
			"codec info requires at least one media type",
		)
	}

	out := make([]MediaType, 0, len(mediaTypes))
	seen := make(map[MediaType]struct{}, len(mediaTypes))
	for i, mediaType := range mediaTypes {
		normalized, err := mediaType.Normalize()
		path := indexedPath(pathCodecInfo+".mediaTypes", i)
		if err != nil {
			return nil, WrapAt(
				path,
				ErrInvalidInfo,
				ErrorReasonInvalidMediaType,
				"codec info media type is invalid",
				err,
			)
		}
		if _, ok := seen[normalized]; ok {
			return nil, errorfAt(
				path,
				ErrInvalidInfo,
				ErrorReasonInvalidMediaType,
				"codec info media type %q is declared more than once",
				normalized,
			)
		}
		seen[normalized] = struct{}{}
		out = append(out, normalized)
	}
	slices.Sort(out)

	return out, nil
}

// normalizeTargets validates, sorts, and detaches targets.
//
// Targets have no Normalize method because the target set is closed-world.
func normalizeTargets(targets []Target) ([]Target, error) {
	if err := validateTargets(targets); err != nil {
		return nil, err
	}

	out := slices.Clone(targets)
	slices.Sort(out)

	return out, nil
}
