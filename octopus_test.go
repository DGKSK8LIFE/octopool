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

// Test for checking the behavior when a new octopus is created with a queue capacity.
func TestOctopusWithQueueCapacity(t *testing.T) {
	testOctopus := NewOctopus(20, 50)

	assert.NotNil(t, testOctopus, "octopus should not be nil")
	assert.Equal(t, 20, testOctopus.workerPool.capacity, "octopus should have the capacity as mentioned (20).")
	assert.Equal(t, 50, testOctopus.jobQueue.QueueCapacity(), "queue should have the capacity as mentioned (50).")
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

// Test for checking the pool capacity.
func TestOctopusCapacity(t *testing.T) {
	testOctopus := NewOctopus(1)

	assert.Equal(t, 1, testOctopus.PoolCapacity())
}

// Test for checking the behavior when an octopus is given a job.
func TestOctopusHandleJob(t *testing.T) {
	testOctopus := NewOctopus(10)

	job1 := func() {
		time.Sleep(1 * time.Microsecond)
	}

	err := testOctopus.HandleJob(job1, "job 1")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

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
		time.Sleep(1 * time.Second)
	}

	job2 := func() {
		time.Sleep(1 * time.Second)
	}

	err := testOctopus.HandleJob(job1, "job 1")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	err = testOctopus.HandleJob(job2, "job 2")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	assert.Equal(t, 1, testOctopus.jobQueue.totalJobs)
}

// Test for checking the behavior when an octopus is given an invalid job.
func TestOctopusInvalidJob(t *testing.T) {
	testOctopus := NewOctopus(10)

	err := testOctopus.HandleJob(nil)

	assert.Equal(t, "invalid function", err.Error())
}

// Test for checking the behavior when an octopus promotes and assigns a job.
func TestOctopusProcessNext(t *testing.T) {
	testOctopus := NewOctopus(1)

	job1 := func() {
		time.Sleep(1 * time.Second)
	}

	job2 := func() {
		time.Sleep(1 * time.Second)
	}

	err := testOctopus.HandleJob(job1, "job 1")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	err = testOctopus.HandleJob(job2, "job 2")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	// one job should be in queue
	assert.Equal(t, 1, testOctopus.jobQueue.totalJobs)

	// simulate wait for job
	time.Sleep(2 * time.Second)

	// job should be promoted from queue to pool
	assert.Equal(t, 0, testOctopus.jobQueue.totalJobs)
}

// Test for checking the behavior when the number of active workers when a job is assigned to an octopus.
func TestOctopusActiveWorkers(t *testing.T) {
	testOctopus := NewOctopus(5)

	job1 := func() {
		time.Sleep(1 * time.Second)
	}

	job2 := func() {
		time.Sleep(1 * time.Second)
	}

	err := testOctopus.HandleJob(job1, "job 1")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	err = testOctopus.HandleJob(job2, "job 2")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	assert.Equal(t, 2, testOctopus.ActiveWorkers())
}

// Test for checking the behavior when the number of available workers when a job is assigned to an octopus.
func TestOctopusAvailableWorkers(t *testing.T) {
	testOctopus := NewOctopus(5)

	job1 := func() {
		time.Sleep(1 * time.Second)
	}

	job2 := func() {
		time.Sleep(2 * time.Second)
	}

	err := testOctopus.HandleJob(job1, "job 1")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	err = testOctopus.HandleJob(job2, "job 2")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	assert.Equal(t, 3, testOctopus.AvailableWorkers())
}

func TestOctopusWait(t *testing.T) {
	testOctopus := NewOctopus(5)

	job1 := func() {
		time.Sleep(1 * time.Second)
	}

	job2 := func() {
		time.Sleep(2 * time.Second)
	}

	err := testOctopus.HandleJob(job1, "job 1")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	err = testOctopus.HandleJob(job2, "job 2")
	if err != nil {
		t.Errorf("Got error while handling job: %v", err)
	}

	testOctopus.Wait()

	assert.Equal(t, 5, testOctopus.AvailableWorkers())
}
