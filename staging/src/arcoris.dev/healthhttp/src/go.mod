module arcoris.dev/healthhttp

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/health v0.0.0
	arcoris.dev/healthtest v0.0.0
)

replace arcoris.dev/health => ../../health/src

replace arcoris.dev/healthtest => ../../healthtest/src

replace arcoris.dev/testutil => ../../testutil/src
