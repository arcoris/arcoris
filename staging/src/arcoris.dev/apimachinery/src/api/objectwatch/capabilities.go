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

package objectwatch

// Capabilities describes optional support advertised by a watch source.
//
// Capabilities are descriptive contract hints. They help callers reject
// obviously unsupported requests before opening a stream, but they do not prove
// continuity, authorize access, require bookmarks, or replace Request/Event
// validation.
type Capabilities struct {
	// StartAtCurrent reports whether the source can open a stream from its
	// current progress point without historical catch-up.
	StartAtCurrent bool
	// HistoricalStart reports whether the source can serve committed changes
	// after a caller-provided non-zero revision.
	HistoricalStart bool
	// Bookmarks reports that the source may emit progress-only bookmarks when
	// Request.AllowBookmarks is true.
	Bookmarks bool
	// RestartEvents reports that the source may emit EventRestartRequired
	// instead of reporting all continuity loss as terminal errors.
	RestartEvents bool
}

// CapabilityReporter is implemented by sources that can describe watch support
// without forcing callers to open a stream.
type CapabilityReporter interface {
	// WatchCapabilities returns the source's advertised watch support.
	WatchCapabilities() Capabilities
}

// SupportsStart verifies whether c supports start's mode.
//
// The start value is validated first, so returned errors preserve normal
// ErrInvalidStart diagnostics for malformed starts as well as unsupported
// source modes.
func (c Capabilities) SupportsStart(start Start) error {
	if err := start.Validate(); err != nil {
		return err
	}

	switch start.Mode {
	case StartAtCurrent:
		if !c.StartAtCurrent {
			return invalidStartError("%s is unsupported by source", start.Mode.String())
		}
	case StartAfterRevision:
		if !c.HistoricalStart {
			return invalidStartError("%s is unsupported by source", start.Mode.String())
		}
	default:
		return invalidStartError("mode %s is invalid", start.Mode.String())
	}

	return nil
}

// SupportsRequest verifies request shape and start support.
//
// Bookmarks are optional in Request. AllowBookmarks=true means a source may
// emit bookmarks if it supports them; it does not mean the caller requires
// bookmarks, so lack of Bookmarks capability is not an error.
func (c Capabilities) SupportsRequest(request Request) error {
	if err := request.Validate(); err != nil {
		return err
	}
	if err := c.SupportsStart(request.Start); err != nil {
		return invalidRequestError("watch.request.start", err, ErrInvalidStart)
	}

	return nil
}
