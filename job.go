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

import "fmt"

// Struct for job.
type Job struct {
	function func() // the job's function
	name     string // name for the job
}

// Formats Job struct.
func (job Job) String() string {
	return fmt.Sprintf("job: %s\n", job.name)
}

// Returns the job's function.
func (job *Job) getJob() func() {
	return job.function
}

// NewJob returns a job with the function wrapped.
func NewJob(fun func()) Job {
	return Job{
		function: fun,
	}
}
