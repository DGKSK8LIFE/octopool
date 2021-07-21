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
