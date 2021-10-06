<div align = "center">
    <img src = "static/octo.png" width="150" height="150">
    <h2>Octopool 🐙</h2>
    <p>A fast, resilient, easy-to-use worker pool for Go.</p>
    <a href = "https://github.com/burntcarrot/octopool/actions?workflow=Tests"><img src = "https://github.com/burntcarrot/octopool/workflows/Tests/badge.svg"></a>
	<a href="https://codecov.io/gh/burntcarrot/octopool"><img src="https://codecov.io/gh/burntcarrot/octopool/branch/main/graph/badge.svg"/></a>
	<a href="https://goreportcard.com/report/github.com/burntcarrot/octopool"><img src="https://goreportcard.com/badge/github.com/burntcarrot/octopool" /></a>
	<a href="https://pkg.go.dev/github.com/burntcarrot/octopool"><img src="https://godoc.org/github.com/burntcarrot/octopool?status.svg" /></a>
	<br><br>
	<img src = "static/term-preview-octopool.svg">
</div>

# Features

- Automatic recycling of workers.
- Job Queue for holding excess jobs when the pool is full.
- Faster performance and lower memory footprint, due to recycling of workers.
- Easy-to-use API for handling jobs; the user just needs to send jobs to the octopus.
- A friendly octopus 🐙 (yes! 😄) for managing all internal operations like promoting jobs to the pool, handling workers and recovering workers when jobs fail.


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

```go
job1 := func() {
    fmt.Println("Hello from octopool!")
}
```

The job is just a function, which can hold anything in its body, like a function call, a print statement, etc.

Once the job has been created, the user can call the octopus to handle the incoming job:

```
octo.HandleJob(job1, "normal-octojob")
```

`octopool` lets the user to name jobs. This is not an required argument, but comes in handy while debugging.

Lastly, the user calls `octo.Wait()`.  This call blocks continued execution until all jobs are finished.

# Example

## Creating an octopus with an invalid capacity:

In this example, you will see:
- How Octopool can prevent panics when the octopus is created with an invalid capacity
- How to create jobs
- How to let the octopus handle jobs
- How to let the octopus handle Wait

```go
package main

import (
	"fmt"

	"github.com/burntcarrot/octopool"
)

func main() {
	octo := octopool.NewOctopus(0)

	job1 := func() {
		fmt.Println("Hello from octopool!")
	}

	job2 := func() {
		fmt.Println("Hello user!")
	}

	for i := 0; i < 1; i++ {
		octo.HandleJob(job1, "normal-octojob")
		octo.HandleJob(job2, "greet-user")
	}
	octo.Wait()
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

## Valid creation of an octopus:

```go
package main

import (
	"fmt"

	"github.com/burntcarrot/octopool"
)

func main() {
	octo := octopool.NewOctopus(10, 100)

	job1 := func() {
		fmt.Println("Hello from octopool!")
	}

	job2 := func() {
		fmt.Println("Hello user!")
	}

	for i := 0; i < 1; i++ {
		octo.HandleJob(job1, "normal-octojob")
		octo.HandleJob(job2, "greet-user")
	}
	octo.Wait()
}
```

Output:

```
2021/07/20 19:06:09 assigning job: normal-octojob to a worker.
Hello from octopool!
2021/07/20 19:06:09 assigning job: greet-user to a worker.
Hello user!
```

## Performing a GET request on `google.com` and printing the response 100 times 

```go
package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/burntcarrot/octopool"
)

func main() {
	octo := octopool.NewOctopus(1, 100)

	headerJob := func() {
		resp, err := http.Get("https://google.com")
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(string(body))
	}

	for i := 0; i < 100; i++ {
		octo.HandleJob(headerJob, "header-job")
	}
	octo.Wait()
}
```

# Benchmarks

Using `go test -benchmem -bench=. -benchtime 10000x benchmark_test.go`:

```
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i5-8265U CPU @ 1.60GHz
Benchmark_Octopus_Pool10_Queue100-8           	   10000	  10083442 ns/op	     297 B/op	       3 allocs/op
Benchmark_Octopus_Pool100_Queue1000-8         	   10000	   1009100 ns/op	     320 B/op	       3 allocs/op
Benchmark_Octopus_Pool1000_Queue10000-8       	   10000	    105579 ns/op	     296 B/op	       3 allocs/op
Benchmark_Octopus_Pool10000_Queue100000-8     	   10000	     20271 ns/op	     871 B/op	       4 allocs/op
PASS
ok  	command-line-arguments	112.982s
```

## Contribution Guide:

Read the Contribution Guide [here](CONTRIBUTING.md).
