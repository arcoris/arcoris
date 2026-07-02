package objectwriter

import "arcoris.dev/apimachinery/api/objectlifecycle"

var _ Applier = (*objectlifecycle.Executor)(nil)
