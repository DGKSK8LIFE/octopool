package octopool

import "sync"

// Represents state for the pool.
type state int

const (
	PoolOpen   state = 0 // can process new tasks
	PoolClosed state = 1 // cannot process new tasks
)

type pool struct {
	status           state      // represents current state of the pool
	capacity         int        // number of workers the pool can accomodate
	availableWorkers sync.Pool  // pool of available workers
	activeWorkers    int        // number of active workers
	closePool        sync.Once  // closes pool and can be called only once
	mu               sync.Mutex // mutex for locking
	octopus          *Octopus   // provides an API to interact with the pool
}

// Basic helper functions:

// Returns pool's capacity.
func (p *pool) getPoolCapacity() int {
	return p.capacity
}

// Returns number of active workers.
func (p *pool) getActiveWorkersCount() int {
	return p.activeWorkers
}

// Checks if a worker is available or not.
func (p *pool) isWorkerAvailable() bool {
	return p.activeWorkers < p.capacity
}

// Pool-related functions:

// Returns a pool with the Capacity specified.
func newPool(capacity int, octopus *Octopus) *pool {
	newPool := pool{
		status:   PoolOpen,
		capacity: capacity,
		availableWorkers: sync.Pool{
			New: func() interface{} {
				return new(worker)
			},
		},
		octopus: octopus,
	}

	return &newPool
}

// Assigns a job to a worker.
func (p *pool) assignJob(job Job) {
	// put lock on pool
	p.mu.Lock()

	defer p.mu.Unlock()

	// tell the pool to give a worker to assign a job
	worker := p.availableWorkers.Get().(*worker)

	// make channel for job
	worker.jobs = make(chan func())

	// set pool for worker
	worker.pool = p

	// run worker
	worker.run()

	// get the job and send it to the jobs channel
	worker.jobs <- job.getJob()

	// increment active worker count
	p.activeWorkers++
}

// Housekeeping function; adds worker to availableWorkers.
func (p *pool) newWorkerAvailable(w *worker) {
	// add worker to available workers
	p.availableWorkers.Put(new(worker))

	// decrement active worker count using mutex
	p.mu.Lock()
	p.activeWorkers--
	p.mu.Unlock()

	// let octopus handle the next job
	p.octopus.processNext()
}

// Sets status to PoolClosed.
func (p *pool) close() {
	p.closePool.Do(func() {
		p.status = PoolClosed
	})
}
