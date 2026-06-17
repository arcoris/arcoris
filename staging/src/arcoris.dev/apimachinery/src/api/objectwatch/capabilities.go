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
// continuity, authorize access, require progress markers, or replace
// Request/Event validation.
type Capabilities struct {
	// StartAtCurrent reports whether the source can open a stream from its
	// current progress point without historical catch-up.
	StartAtCurrent bool
	// HistoricalStart reports whether the source can serve committed changes
	// after a caller-provided revision boundary, including zero.
	HistoricalStart bool
	// Progress reports that the source may emit EventProgress markers when
	// Request.AllowProgress is true.
	Progress bool
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
// The start value is validated first. Malformed starts return ErrInvalidStart;
// valid starts that the source cannot serve return ErrUnsupportedCapability.
func (c Capabilities) SupportsStart(start Start) error {
	if err := start.Validate(); err != nil {
		return err
	}

	switch start.Mode {
	case StartAtCurrent:
		if !c.StartAtCurrent {
			return unsupportedCapabilityError(
				"watch.capabilities.start",
				ErrUnsupportedCapability,
			)
		}
	case StartAfterRevision:
		if !c.HistoricalStart {
			return unsupportedCapabilityError(
				"watch.capabilities.start",
				ErrUnsupportedCapability,
			)
		}
	default:
		return invalidStartError("mode %s is invalid", start.Mode.String())
	}

	return nil
}

// SupportsRequest verifies request shape and start support.
//
// Progress markers are optional in Request. AllowProgress=true means a source
// may emit EventProgress if it supports progress reporting; it does not mean
// the caller requires progress events, so lack of Progress capability is not an
// error.
func (c Capabilities) SupportsRequest(request Request) error {
	if err := request.Validate(); err != nil {
		return err
	}

	return c.SupportsStart(request.Start)
}
