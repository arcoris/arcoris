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

// Validate checks that i is complete, canonical codec metadata.
//
// Validate is intentionally strict. It does not trim, lowercase, sort, or
// deduplicate declarations. Use Normalize when caller-supplied metadata should
// be canonicalized before storage or registration.
func (i Info) Validate() error {
	if err := validateFormatAt(pathCodecInfo+".format", i.Format); err != nil {
		return WrapAt(
			pathCodecInfo+".format",
			ErrInvalidInfo,
			ErrorReasonInvalidFormat,
			"codec info format is invalid",
			err,
		)
	}
	if err := validateMediaTypes(i.MediaTypes); err != nil {
		return err
	}
	if err := validateTargets(i.Targets); err != nil {
		return err
	}

	return nil
}

// ValidateInfo checks that info is complete, canonical codec metadata.
//
// The function form mirrors nearby API packages and is convenient when callers
// want a package-level validation hook.
func ValidateInfo(info Info) error {
	return info.Validate()
}

// validateMediaTypes checks the media type list without normalizing values.
//
// Duplicate detection uses canonical MediaType values exactly as supplied. This
// keeps Validate strict and prevents accidental acceptance of mixed-case input.
func validateMediaTypes(mediaTypes []MediaType) error {
	if len(mediaTypes) == 0 {
		return ErrorAt(
			pathCodecInfo+".mediaTypes",
			ErrInvalidInfo,
			ErrorReasonInvalidMediaType,
			"codec info requires at least one media type",
		)
	}

	seen := make(map[MediaType]struct{}, len(mediaTypes))
	for i, mediaType := range mediaTypes {
		path := indexedPath(pathCodecInfo+".mediaTypes", i)
		if err := validateMediaTypeAt(path, mediaType); err != nil {
			return WrapAt(
				path,
				ErrInvalidInfo,
				ErrorReasonInvalidMediaType,
				"codec info media type is invalid",
				err,
			)
		}
		if _, ok := seen[mediaType]; ok {
			return errorfAt(
				path,
				ErrInvalidInfo,
				ErrorReasonInvalidMediaType,
				"codec info media type %q is declared more than once",
				mediaType,
			)
		}
		seen[mediaType] = struct{}{}
	}

	return nil
}

// validateTargets checks the target list without changing order.
//
// Target duplicates are invalid because future registries need one stable
// declaration per supported document model.
func validateTargets(targets []Target) error {
	if len(targets) == 0 {
		return ErrorAt(
			pathCodecInfo+".targets",
			ErrInvalidInfo,
			ErrorReasonInvalidTarget,
			"codec info requires at least one target",
		)
	}

	seen := make(map[Target]struct{}, len(targets))
	for i, target := range targets {
		path := indexedPath(pathCodecInfo+".targets", i)
		if err := validateTargetAt(path, target); err != nil {
			return WrapAt(
				path,
				ErrInvalidInfo,
				ErrorReasonInvalidTarget,
				"codec info target is invalid",
				err,
			)
		}
		if _, ok := seen[target]; ok {
			return errorfAt(
				path,
				ErrInvalidInfo,
				ErrorReasonInvalidTarget,
				"codec info target %q is declared more than once",
				target,
			)
		}
		seen[target] = struct{}{}
	}

	return nil
}
