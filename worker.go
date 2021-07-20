package octopool

import "log"

type worker struct {
	jobs chan func() // channel for receiving jobs
	pool *pool       // pool reference
}

// Executes the job provided to the worker.
func (w *worker) run() {
	go func() {
		defer func() {
			// silently recover from error, do not panic
			if r := recover(); r != nil {
				// print the error to the console
				log.Printf("Recovered error: %v\n", r)
				// return the worker back due to abrupt failure
				w.pool.newWorkerAvailable(w)
			}
		}()

		// receive job
		job := <-w.jobs

		// execute job
		job()

		// return worker back once job is completed
		w.pool.newWorkerAvailable(w)
	}()
}
