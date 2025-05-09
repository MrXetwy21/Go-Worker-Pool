package pool

import (
	"sync"

	"github.com/MrXetwy21/worker-pool/pkg/result"
)

type Task interface {
	Process() (result.Result, error)
}
type Pool struct {
	numWorkers int
	tasks      chan Task
	results    chan result.Result
	done       chan bool
	wg         sync.WaitGroup
}

func New(numWorkers int) *Pool {
	return &Pool{
		numWorkers: numWorkers,
		tasks:      make(chan Task),
		results:    make(chan result.Result),
		done:       make(chan bool),
	}
}

func (p *Pool) worker() {
	defer p.wg.Done()

	for task := range p.tasks {
		res, err := task.Process()

		if err != nil {
			res.Error = err
		}

		p.results <- res
	}
}

func (p *Pool) Process(tasks []Task) {
	p.wg.Add(p.numWorkers)
	for i := 0; i < p.numWorkers; i++ {
		go p.worker()
	}

	go func() {
		for _, task := range tasks {
			p.tasks <- task
		}
		close(p.tasks)
		p.wg.Wait()
		close(p.results)
	}()
}

func (p *Pool) Wait() []result.Result {
	var results []result.Result
	for res := range p.results {
		results = append(results, res)
	}
	return results
}

func (p *Pool) WithCallback(tasks []Task, callback func(result.Result)) {
	go func() {
		for res := range p.results {
			callback(res)
		}
	}()

	p.Process(tasks)
}
