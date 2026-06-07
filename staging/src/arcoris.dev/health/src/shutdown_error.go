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

package health

import "errors"

var (
	// ErrNilSourceChannel identifies a nil source channel passed to a
	// channel-backed health check constructor.
	//
	// A nil channel would never become ready and would make the check report a
	// permanently healthy state. Constructors reject it because that usually
	// hides a wiring error.
	ErrNilSourceChannel = errors.New("health: nil source channel")

	// ErrNilSourceContext identifies a nil source context passed to a
	// context-backed health check constructor.
	//
	// A nil context has no Done channel and cannot represent shutdown or drain
	// state. Constructors reject it instead of silently replacing it with
	// context.Background.
	ErrNilSourceContext = errors.New("health: nil source context")
)
