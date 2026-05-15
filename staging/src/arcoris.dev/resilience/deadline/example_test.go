/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package deadline_test

import (
	"context"
	"fmt"
	"time"

	"arcoris.dev/resilience/deadline"
)

func ExampleInspect() {
	now := time.Now().Add(time.Hour)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(250*time.Millisecond))
	defer cancel()

	budget := deadline.Inspect(ctx, now)
	fmt.Println(budget.HasDeadline)
	fmt.Println(budget.Remaining)

	// Output:
	// true
	// 250ms
}

func ExampleCanStart() {
	now := time.Now().Add(time.Hour)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(50*time.Millisecond))
	defer cancel()

	decision := deadline.CanStart(ctx, now, 100*time.Millisecond)
	fmt.Println(decision.IsAllowed())
	fmt.Println(decision.Reason)

	// Output:
	// false
	// insufficient_budget
}

func ExampleClamp() {
	now := time.Now().Add(time.Hour)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(75*time.Millisecond))
	defer cancel()

	d, ok := deadline.Clamp(ctx, now, 200*time.Millisecond)
	fmt.Println(ok)
	fmt.Println(d)

	// Output:
	// true
	// 75ms
}

func ExampleReserve() {
	now := time.Now().Add(time.Hour)
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(200*time.Millisecond))
	defer cancel()

	d, bounded, ok := deadline.Reserve(ctx, now, 50*time.Millisecond)
	fmt.Println(ok)
	fmt.Println(bounded)
	fmt.Println(d)

	// Output:
	// true
	// true
	// 150ms
}

func ExampleWithBoundedTimeout() {
	now := time.Now().Add(time.Hour)
	parent, cancelParent := context.WithDeadline(context.Background(), now.Add(100*time.Millisecond))
	defer cancelParent()

	child, cancelChild := deadline.WithBoundedTimeout(parent, now, time.Second)
	defer cancelChild()

	dl, ok := child.Deadline()
	fmt.Println(ok)
	fmt.Println(dl.Sub(now))

	// Output:
	// true
	// 100ms
}
