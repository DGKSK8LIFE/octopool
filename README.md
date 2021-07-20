<div align = "center">
    <img src = "static/octopool.png">
    <h2>Octopool</h2>
    <p>A fast, resilient, easy-to-use worker pool for Go.</p>
    <br>
</div>

<h1>Table of Contents:</h1>

- [Installation](#installation)
- [Example](#example)

# Installation

`octopool` can be installed using `go get`:

```
go get -u github.com/burntcarrot/octopool
```

# Example

Here's an example on how Octopool uses a default pool capacity in order to prevent panics:

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
