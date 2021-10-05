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
	"errors"
	"log"
)

// pre-defined pool capacity
const defaultPoolCapacity = 1000

// pre-defined queue capacity
const defaultQueueCapacity = 10000

// predefined errors

// ErrNilFunction is the error raised when a function is invalid.
var ErrNilFunction = errors.New("invalid function")

// ErrInvalidPoolCapacity is the error raised when the pool capacity provided is invalid in nature.
var ErrInvalidPoolCapacity = errors.New("invalid pool capacity: pool capacity must be a positive number, cannot process jobs in a pool with a capacity equal to or less than zero")

// ErrInvalidPoolState is the error raised when the pool is closed and no jobs can be sent to the pool.
var ErrInvalidPoolState = errors.New("cannot assign job to closed pool")

// Octopus is a struct for representing the octopus which handles the execution of jobs.
type Octopus struct {
	workerPool   *pool     // worker pool
	jobQueue     *JobQueue // job queue for holding tasks
	poolCapacity int       // pool capacity
}

// Basic helper functions:

// PoolCapacity returns pool's capacity.
func (octo *Octopus) PoolCapacity() int {
	return octo.workerPool.getPoolCapacity()
}

// ActiveWorkers returns number of active workers.
func (octo *Octopus) ActiveWorkers() int {
	return octo.workerPool.getActiveWorkersCount()
}

// AvailableWorkers returns number of available workers.
func (octo *Octopus) AvailableWorkers() int {
	return octo.workerPool.getPoolCapacity() - octo.workerPool.getActiveWorkersCount()
}

// Close closes the worker pool.
func (octo *Octopus) Close() {
	octo.workerPool.close()
}

// Octopus related functions:

// NewOctopus creates an octopus with the capacity specified.
func NewOctopus(capacity int, queueCapacity ...int) *Octopus {
	// handle invalid capacity
	if capacity <= 0 {
		// Suppress panic by using a default pool capacity
		// panic(ErrInvalidPoolCapacity)

		log.Println(ErrInvalidPoolCapacity)
		log.Printf("using defaultPoolCapacity = %d as the pool capacity due to invalid pool capacity provided.", defaultPoolCapacity)
		capacity = defaultPoolCapacity
	}

	var octopus *Octopus

	// create an octopus with the capacity
	if queueCapacity == nil {
		log.Printf("using defaultQueueCapacity = %d as the queue capacity due to invalid queue capacity provided.", defaultQueueCapacity)
		octopus = &Octopus{
			jobQueue:     NewJobQueue(defaultQueueCapacity),
			poolCapacity: capacity,
		}
	} else {
		octopus = &Octopus{
			jobQueue:     NewJobQueue(queueCapacity[0]),
			poolCapacity: capacity,
		}
	}

	// create a pool
	pool := newPool(capacity, octopus)
	// set the created pool as the worker pool for octopus
	octopus.workerPool = pool

	return octopus
}

// HandleJob assigns a job to a worker if workers are available, else, adds to the job queue.
func (octo *Octopus) HandleJob(fun func(), name ...string) error {
	// throw error if pool is closed
	if octo.workerPool.status == PoolClosed {
		return ErrInvalidPoolState
	}

	// throw error if function provided is invalid
	if fun == nil {
		return ErrNilFunction
	}

	// create a job
	job := Job{function: fun, name: name[0]}

	// check if workers are available and assign a job, else add the job to the queue
	if octo.workerPool.isWorkerAvailable() {
		log.Println("assigning job:", job.name, "to a worker.")
		octo.workerPool.assignJob(job)
	} else {
		log.Printf("adding job: %s to queue\n", job.name)
		octo.jobQueue.AddJob(job)
	}

	return nil
}

// Promotes a job to the pool and assigns a worker to it.
func (octo *Octopus) processNext() {
	// if queue is not empty, remove job from the queue
	if octo.jobQueue.IsNotEmpty() {
		job, err := octo.jobQueue.RemoveJob()
		log.Println("removing job:", job.name, "from queue and assigning to a worker.")

		if err != nil {
			log.Println("error occurred while removing the job.")
			return
		}

		// assign the job to the worker
		octo.workerPool.assignJob(job)
		log.Println("assigned job:", job.name, "to a worker.")
	}
}

// Waits on workers to finish the job
func (octo *Octopus) Wait() {
	log.Println("Waiting for jobs to finish....")
	octo.workerPool.wg.Wait()
}
