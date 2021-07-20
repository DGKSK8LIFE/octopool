package octopool

import (
	"errors"
	"log"
)

// pre-defined pool capacity
const defaultPoolCapacity = 10

// pre-defined queue capacity
const defaultQueueCapacity = 10

// predefined errors

var ErrNilFunction = errors.New("invalid function")

var ErrInvalidPoolCapacity = errors.New("invalid pool capacity: pool capacity must be a positive number, cannot process jobs in a pool with a capacity equal to or less than zero")

var ErrInvalidPoolState = errors.New("cannot assign job to closed pool")

type Octopus struct {
	workerPool   *pool     // worker pool
	jobQueue     *JobQueue // job queue for holding tasks
	poolCapacity int       // pool capacity
}

// Basic helper functions:

// Returns pool's capacity.
func (octo *Octopus) PoolCapacity() int {
	return octo.workerPool.getPoolCapacity()
}

// Returns number of active workers.
func (octo *Octopus) ActiveWorkers() int {
	return octo.workerPool.getActiveWorkersCount()
}

// Returns number of available workers.
func (octo *Octopus) AvailableWorkers() int {
	return octo.workerPool.getPoolCapacity() - octo.workerPool.getActiveWorkersCount()
}

// Closes the worker pool.
func (octo *Octopus) Close() {
	octo.workerPool.close()
}

// Octopus related functions:

// Creates an octopus with the capacity specified.
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

// Assigns a job to a worker if workers are available, else, adds to the job queue.
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
