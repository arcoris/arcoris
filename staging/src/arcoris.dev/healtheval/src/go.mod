module arcoris.dev/healtheval

go 1.25.0

toolchain go1.25.9

require (
	arcoris.dev/chrono v0.0.0
	arcoris.dev/health v0.0.0
)

replace arcoris.dev/chrono => ../../chrono/src

replace arcoris.dev/health => ../../health/src

replace arcoris.dev/testutil => ../../testutil/src
