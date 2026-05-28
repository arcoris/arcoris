module arcoris.dev/healthtest

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/health v0.0.0
	arcoris.dev/testutil v0.0.0
)

replace arcoris.dev/health => ../../health/src

replace arcoris.dev/testutil => ../../testutil/src
