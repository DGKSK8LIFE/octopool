<div align = "center">
    <img src = "static/octobanner.svg">
    <h2>Octopool</h2>
    <p>A fast, resilient, easy-to-use worker pool for Go.</p>
    <br>
</div>

<h1>Table of Contents:</h1>

- [Installation](#installation)
- [How Octopool Works?](#how-octopool-works)
- [Usage](#usage)
- [Example](#example)

# Installation

`octopool` can be installed using `go get`:

```
go get -u github.com/burntcarrot/octopool
```

# How Octopool Works?

Octopool is based on the principle that users can call the octopus to handle jobs for them.

The octopus:
- Assigns the given job to a worker if available
- Promotes a job to the pool and assigns a worker to it, if workers are available after completing jobs.
- Executes jobs with the help of a worker.
- Maintains the job queue, which is used when the number of jobs exceed the pool's capacity to hold pending jobs.

# Usage

The octopus provides an easy-to-use API to create and handle jobs.

The user can create a job like this:

```
job1 := func() {
    defer wg.Done()
    fmt.Println("Hello from octopool!")
}
```

The job is just a function, which can hold anything in its body, like a function call, a print statement, etc.

> Note: Remember to keep the job's initial line to `defer wg.Done()` as it would prevent other jobs to abruptly stop the current job's execution. **You should maintain a WaitGroup to prevent overriding execution of jobs. The example with a WaitGroup is given in the [Example](#example) section.**

Once the job has been created, the user can call the octopus to handle the incoming job:

```
octo.HandleJob(job1, "normal-octojob")
```

`octopool` lets the user to name jobs. This is not an required argument, but comes in handy while debugging.

# Example

In this example, you will see:
- How Octopool can prevent panics when the octopus is created with an invalid capacity
- How to create jobs
- How to let the octopus handle jobs
- How to implement a WaitGroup to prevent overriding of jobs

```
package main

import (
	"fmt"
	"sync"

	"github.com/burntcarrot/octopool"
)

func main() {
	var wg sync.WaitGroup
	octo := octopool.NewOctopus(0)

	job1 := func() {
		defer wg.Done()
		fmt.Println("Hello from octopool!")
	}

	job2 := func() {
		defer wg.Done()
		fmt.Println("Hello user!")
	}

	for i := 0; i < 1; i++ {
		wg.Add(1)
		octo.HandleJob(job1, "normal-octojob")
		wg.Add(1)
		octo.HandleJob(job2, "greet-user")
		wg.Wait()
	}
}
```

Output:

```
2021/07/20 19:06:09 invalid pool capacity: pool capacity must be a positive number, cannot process jobs in a pool with a capacity equal to or less than zero
2021/07/20 19:06:09 using defaultPoolCapacity = 10 as the pool capacity due to invalid pool capacity provided.
2021/07/20 19:06:09 using defaultQueueCapacity = 10 as the queue capacity due to invalid queue capacity provided.
2021/07/20 19:06:09 assigning job: normal-octojob to a worker.
Hello from octopool!
2021/07/20 19:06:09 assigning job: greet-user to a worker.
Hello user!
```

More examples to be added.
