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

package fixedwindow_test

import (
	"fmt"
	"sync"
	"time"

	"arcoris.dev/resilience/retrybudget/fixedwindow"
)

func ExampleLimiter_windowRotation() {
	clk := &exampleClock{now: time.Unix(100, 0).UTC()}
	budget, err := fixedwindow.New(
		fixedwindow.WithClock(clk),
		fixedwindow.WithWindow(time.Second),
		fixedwindow.WithRatio(1),
		fixedwindow.WithMinRetries(0),
	)
	if err != nil {
		panic(err)
	}

	budget.RecordOriginal()
	first := budget.Snapshot()

	clk.Add(time.Second)
	budget.RecordOriginal()
	second := budget.Snapshot()

	fmt.Println(first.Value.Window.StartedAt.Unix())
	fmt.Println(second.Value.Window.StartedAt.Unix())
	fmt.Println(second.Value.Attempts.Original)

	// Output:
	// 100
	// 101
	// 1
}

type exampleClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *exampleClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *exampleClock) Since(t time.Time) time.Duration {
	return c.Now().Sub(t)
}

func (c *exampleClock) Until(t time.Time) time.Duration {
	return t.Sub(c.Now())
}

func (c *exampleClock) Add(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.now = c.now.Add(d)
}
