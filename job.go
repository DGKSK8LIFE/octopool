package octopool

import "fmt"

type Job struct {
	function func() // the job's function
	name     string // name for the job
}

// Formats Job struct.
func (job Job) String() string {
	return fmt.Sprintf("job: %s\n", job.name)
}

// Returns the job's function.
func (t *Job) getJob() func() {
	return t.function
}

// Returns a job with the function wrapped.
func NewJob(fun func()) Job {
	return Job{
		function: fun,
	}
}
