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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test for checking the behavior when a new octopus is created.
func TestOctopus(t *testing.T) {
	testOctopus := NewOctopus(20)

	assert.NotNil(t, testOctopus, "octopus should not be nil")
	assert.Equal(t, 20, testOctopus.workerPool.capacity, "octopus should have the capacity as mentioned (20).")
}

// Test for checking the behavior when an octopus with an invalid capacity is created.
func TestInvalidOctopus(t *testing.T) {
	testOctopus := NewOctopus(0)

	assert.NotNil(t, testOctopus, "octopus should not be nil")
	assert.Equal(t, 1000, testOctopus.workerPool.capacity, "octopus should have the default capacity (1000).")
}

// Test for checking the behavior when an octopus with an negative capacity is created.
func TestNegativeOctopus(t *testing.T) {
	testOctopus := NewOctopus(-1)

	assert.NotNil(t, testOctopus, "octopus should not be nil")
	assert.Equal(t, 1000, testOctopus.workerPool.capacity, "octopus should have the default capacity (1000).")
}

// Test for checking the behavior when an octopus is given a job.
func TestOctopusHandleJob(t *testing.T) {
	testOctopus := NewOctopus(10)

	job1 := func() {
		time.Sleep(1 * time.Microsecond)
	}

	testOctopus.HandleJob(job1, "job 1")

	assert.Equal(t, 1, testOctopus.ActiveWorkers())
}

// Test for checking the behavior when an octopus is given a job when the pool is closed.
func TestOctopusHandleJobClosedPool(t *testing.T) {
	testOctopus := NewOctopus(10)
	testOctopus.Close()

	job1 := func() {
		time.Sleep(1 * time.Microsecond)
	}

	err := testOctopus.HandleJob(job1, "job 1")

	assert.Equal(t, "cannot assign job to closed pool", err.Error())
}

// Test for checking the behavior when an octopus is given a job when the pool is full.
func TestOctopusHandleJobFullPool(t *testing.T) {
	testOctopus := NewOctopus(1)

	job1 := func() {
		time.Sleep(1 * time.Microsecond)
	}

	job2 := func() {
		time.Sleep(2 * time.Microsecond)
	}

	testOctopus.HandleJob(job1, "job 1")
	testOctopus.HandleJob(job2, "job 2")

	assert.Equal(t, 1, testOctopus.jobQueue.totalJobs)
}

// Test for checking the behavior when an octopus is given an invalid job.
func TestOctopusInvalidJob(t *testing.T) {
	testOctopus := NewOctopus(10)

	err := testOctopus.HandleJob(nil)

	assert.Equal(t, "invalid function", err.Error())
}
