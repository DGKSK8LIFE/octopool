// Copyright 2021 Aadhav Vignesh

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package octopool

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test for checking the behavior when a worker fails.
func TestWorker(t *testing.T) {
	testOctopus := NewOctopus(1)

	job1 := func() {
		x := 10
		y := 0
		fmt.Println("Raising error intentionally:", x/y)
	}

	testOctopus.HandleJob(job1, "division")

	time.Sleep(1 * time.Second)
	// if worker returns error, then it will be available
	assert.Equal(t, 1, testOctopus.AvailableWorkers())
}
