// Package examples demonstrates Go concurrency patterns.
package examples

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"
)

// ─── Worker Pool Pattern ────────────────────────────────────────
// N workers xử lý jobs từ shared channel.

type Job struct {
	ID   int
	Data string
}

type Result struct {
	JobID  int
	Output string
}

func WorkerPool(ctx context.Context, jobs []Job, numWorkers int) []Result {
	jobCh := make(chan Job, len(jobs))
	resultCh := make(chan Result, len(jobs))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range jobCh {
				select {
				case <-ctx.Done():
					return
				default:
					result := processJob(workerID, job)
					resultCh <- result
				}
			}
		}(i)
	}

	// Send jobs
	for _, job := range jobs {
		jobCh <- job
	}
	close(jobCh)

	// Wait and collect results
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var results []Result
	for r := range resultCh {
		results = append(results, r)
	}
	return results
}

func processJob(workerID int, job Job) Result {
	return Result{
		JobID:  job.ID,
		Output: fmt.Sprintf("worker-%d processed: %s", workerID, job.Data),
	}
}

// ─── Fan-Out / Fan-In ───────────────────────────────────────────
// Multiple goroutines read from the same channel (fan-out),
// results merged into one channel (fan-in).

func FanOutFanIn(ctx context.Context, inputs []string, numWorkers int) <-chan string {
	inputCh := make(chan string, len(inputs))
	for _, input := range inputs {
		inputCh <- input
	}
	close(inputCh)

	// Fan-out: multiple workers read from inputCh
	resultChannels := make([]<-chan string, numWorkers)
	for i := 0; i < numWorkers; i++ {
		resultChannels[i] = worker(ctx, inputCh)
	}

	// Fan-in: merge all result channels
	return merge(ctx, resultChannels...)
}

func worker(ctx context.Context, inputs <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for input := range inputs {
			select {
			case <-ctx.Done():
				return
			case out <- fmt.Sprintf("processed: %s", input):
			}
		}
	}()
	return out
}

func merge(ctx context.Context, channels ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup

	for _, ch := range channels {
		wg.Add(1)
		go func(c <-chan string) {
			defer wg.Done()
			for val := range c {
				select {
				case <-ctx.Done():
					return
				case out <- val:
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

// ─── Context Cancellation ───────────────────────────────────────
// Graceful cancellation với context.

func LongRunningTask(ctx context.Context) error {
	for i := 0; i < 100; i++ {
		select {
		case <-ctx.Done():
			// Cleanup nếu cần
			return ctx.Err()
		default:
			// Xử lý chunk i
			fmt.Printf("Processing chunk %d\n", i)
		}
	}
	return nil
}

// ─── errgroup — Concurrent Operations with Error Handling ───────
// errgroup tự động cancel context khi bất kỳ goroutine nào fail.

func FetchAll(ctx context.Context, urls []string) (map[string]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	results := make(map[string]string)
	var mu sync.Mutex

	for _, url := range urls {
		url := url // capture loop variable
		g.Go(func() error {
			data, err := fetch(ctx, url)
			if err != nil {
				return fmt.Errorf("fetching %s: %w", url, err)
			}

			mu.Lock()
			results[url] = data
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

// errgroup with concurrency limit (Go 1.20+)
func FetchAllLimited(ctx context.Context, urls []string, maxConcurrent int) (map[string]string, error) {
	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(maxConcurrent) // limit concurrent goroutines

	results := make(map[string]string)
	var mu sync.Mutex

	for _, url := range urls {
		url := url
		g.Go(func() error {
			data, err := fetch(ctx, url)
			if err != nil {
				return fmt.Errorf("fetching %s: %w", url, err)
			}
			mu.Lock()
			results[url] = data
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return results, nil
}

func fetch(ctx context.Context, url string) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
		return fmt.Sprintf("data from %s", url), nil
	}
}
