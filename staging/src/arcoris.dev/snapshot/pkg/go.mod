module arcoris.dev/snapshot

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/atomicx v0.0.0
	arcoris.dev/chrono v0.0.0
	arcoris.dev/testutil v0.0.0
)

replace arcoris.dev/atomicx => ../../atomicx/pkg
replace arcoris.dev/chrono => ../../chrono/pkg
replace arcoris.dev/testutil => ../../testutil/pkg
