package octopool

import (
	"errors"
)

type JobQueue struct {
	jobQueue  []Job // job queue
	totalJobs int   // total jobs present in the queue
}

// Helper functions:

// Checks if the job queue is empty or not.
func (jobQueue *JobQueue) IsNotEmpty() bool {
	return jobQueue.totalJobs > 0
}

// Job Queue related functions:

// Returns a job queue with the specified capacity.
func NewJobQueue(queueCapacity int) *JobQueue {
	jobQueue := JobQueue{jobQueue: make([]Job, 0, queueCapacity), totalJobs: 0}
	return &jobQueue
}

// Adds a job to the job queue.
func (jobQueue *JobQueue) AddJob(job Job) {
	jobQueue.jobQueue = append(jobQueue.jobQueue, job)
	jobQueue.totalJobs++
}

// Removes a job from the job queue.
func (jobQueue *JobQueue) RemoveJob() (Job, error) {
	// remove job from the queue if there exists a job in the queue
	if jobQueue.totalJobs > 0 {
		job := jobQueue.jobQueue[0]

		// use slicing for removing the job
		jobQueue.jobQueue = append(jobQueue.jobQueue[:0], jobQueue.jobQueue[1:]...)

		jobQueue.totalJobs--
		return job, nil
	}

	// else return an error
	return Job{}, errors.New("empty job queue")
}
