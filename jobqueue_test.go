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
)

// pre-defined queue capacity
const queueCapacity = 10

// Test for checking the behavior when a new job is added to the job queue.
func TestAddJob(t *testing.T) {
	testQueue := NewJobQueue(queueCapacity)
	testQueue.AddJob(NewJob(func() {
		fmt.Println("Sample Job.")
	}))

	if testQueue.totalJobs != 1 {
		t.Error("Mismatch in job count.")
	}
}

// Test for checking the behavior when a job is removed from the job queue.
func TestRemoveJob(t *testing.T) {
	testQueue := NewJobQueue(queueCapacity)
	testQueue.AddJob(NewJob(func() {
		fmt.Println("Sample")
	}))

	_, err := testQueue.RemoveJob()

	if err != nil {
		t.Errorf("Got error: %v", err)
	}
}

// Test for checking the behavior when a job is removed from the job queue when the job queue is empty.
func TestRemoveJobWhenEmpty(t *testing.T) {
	testQueue := NewJobQueue(queueCapacity)
	_, err := testQueue.RemoveJob()

	if err == nil {
		t.Error("Expected a empty queue error.")
	}
}
