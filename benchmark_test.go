package octopool_test

import (
	"sync"
	"testing"
	"time"

	"github.com/burntcarrot/octopool"
)

func benchmarkOctopus(poolCapacity int, queueCapacity int, b *testing.B) {
	pool := octopool.NewOctopus(poolCapacity, queueCapacity)

	var wg sync.WaitGroup

	job1 := func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
	}

	for i := 0; i < b.N; i++ {
		wg.Add(1)
		pool.HandleJob(job1, "normal-octojob")
	}
	wg.Wait()
}

func Benchmark_Octopus_Pool10_Queue100(b *testing.B) {
	benchmarkOctopus(10, 100, b)
}

func Benchmark_Octopus_Pool100_Queue1000(b *testing.B) {
	benchmarkOctopus(100, 1000, b)
}

func Benchmark_Octopus_Pool1000_Queue10000(b *testing.B) {
	benchmarkOctopus(1000, 10000, b)
}

func Benchmark_Octopus_Pool10000_Queue100000(b *testing.B) {
	benchmarkOctopus(10000, 100000, b)
}
