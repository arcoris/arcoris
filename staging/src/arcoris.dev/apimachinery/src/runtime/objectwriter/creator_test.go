package objectwriter

import "arcoris.dev/apimachinery/api/objectlifecycle"

var _ Creator = (*objectlifecycle.Executor)(nil)
