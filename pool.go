/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/12/21
 * Time: 11:49 AM
 *
 * Description:
 *
 */

package gowl

import (
	"context"
	"errors"
	"fmt"
	"github.com/hamed-yousefi/gowl/status/pool"
	"github.com/hamed-yousefi/gowl/status/process"
	"github.com/hamed-yousefi/gowl/status/worker"
	"log"
	"sync"
	"time"
)

const (
	// defaultWorkerName is the default worker name prefix.
	defaultWorkerName = "W%d"
)

type (
	// WorkerName is a custom type of string that represents worker's name.
	WorkerName string

	// PID is a custom type of string that represents process id.
	PID string

	// Process is an interface that represents a process
	Process interface {
		// Start runs the process. It returns an error object if any thing wrong
		// happens in runtime.
		Start() error
		// Name returns process name.
		Name() string
		// PID returns process id.
		PID() PID
	}

	// Pool is a mechanism to dispatch processes between a group of workers.
	Pool interface {
		// Start runs the pool.
		Start() error
		// Register adds the process to the pool queue.
		Register(p ...Process)
		// Close stops a running pool.
		Close() error
		// Kill cancel a process before it starts.
		Kill(pid PID)
		// Monitor returns pool monitor.
		Monitor() Monitor
	}

	// Monitor is a mechanism for observation processes and pool stats.
	Monitor interface {
		// PoolStatus returns pool status
		PoolStatus() pool.Status
		// Error returns process's error by process id.
		Error(PID) error
		// WorkerList returns the list of worker names of the pool.
		WorkerList() []WorkerName
		// WorkerStatus returns worker status. It accepts worker name as input.
		WorkerStatus(name WorkerName) worker.Status
		// ProcessStatus returns process stats. It accepts process id as input.
		ProcessStats(pid PID) ProcessStats
	}

	// ProcessStats represents process statistics.
	ProcessStats struct {
		// WorkerName is the name of the worker that this process belongs to.
		WorkerName WorkerName

		// Process is process that this stats belongs to.
		Process Process

		// Status represents the current state of the process.
		Status process.Status

		// StartedAt represents the start date time of the process.
		StartedAt time.Time

		// FinishedAt represents the end date time of the process.
		FinishedAt time.Time

		err error
	}

	// workerPool is an implementation of Pool and Monitor interfaces.
	workerPool struct {
		status       pool.Status
		size         int
		queue        chan Process
		wg           *sync.WaitGroup
		processes    *processStatusMap
		workers      []WorkerName
		workersStats *workerStatsMap
		controlPanel *controlPanelMap
		mutex        *sync.Mutex
		isClosed     bool
	}
)

// NewPool makes a new instance of Pool. I accept an integer value as input
// that represents pool size.
func NewPool(size int) Pool {
	return &workerPool{
		status:       pool.Created,
		size:         size,
		queue:        make(chan Process, size),
		workers:      []WorkerName{},
		processes:    new(processStatusMap),
		workersStats: new(workerStatsMap),
		controlPanel: new(controlPanelMap),
		mutex:        new(sync.Mutex),
		wg:           new(sync.WaitGroup),
	}
}

// Start runs the pool. It returns error if pool is already in running state.
// It changes the pool state to Running and calls workerPool.run() function to
// run the pool.
func (w *workerPool) Start() error {

	if w.status == pool.Running {
		return errors.New("unable to start the pool, status: " + w.status.String())
	}

	w.status = pool.Running
	w.run()

	return nil
}

// run is the function that creates worker and starts the pool.
func (w *workerPool) run() {

	// Create workers
	for i := 0; i < w.size; i++ {
		// For each worker add one to the waitGroup.
		w.wg.Add(1)
		wName := WorkerName(fmt.Sprintf(defaultWorkerName, i))
		w.workers = append(w.workers, wName)

		// Create worker.
		go func(wn WorkerName) {
			defer w.wg.Done()

			// Consume process from the queue.
			for p := range w.queue {
				w.workersStats.put(wn, worker.Busy)
				pStats := w.processes.get(p.PID())
				pStats.Status = process.Running
				pStats.StartedAt = time.Now()
				pStats.WorkerName = wn
				w.processes.put(p.PID(), pStats)
				wgp := new(sync.WaitGroup)
				wgp.Add(1)

				go func() {
					stats := w.processes.get(p.PID())
					defer func() {
						w.processes.put(p.PID(), stats)
						wgp.Done()
					}()
					pContext := w.controlPanel.get(p.PID())
					select {
					case <-pContext.ctx.Done():
						log.Printf("processFunc with id %s has been killed.\n", p.PID().String())
						stats.Status = process.Killed
						return
					default:
						if err := p.Start(); err != nil { //nolint:typecheck
							stats.err = err
							stats.Status = process.Failed
						} else {
							stats.Status = process.Succeeded
						}
						pContext.cancel()
					}
				}()

				wgp.Wait()
				pStats = w.processes.get(p.PID())
				pStats.FinishedAt = time.Now()
				w.processes.put(p.PID(), pStats)
				w.workersStats.put(wn, worker.Waiting)
			}
		}(wName)
	}
}

// Register adds the process to the pool queue. It accept a list of processes
// and adds them to the queue. It publishes the process to queue in a separate
// goroutine. It means that Register function provides multi-publisher that
// each of them works asynchronously.
func (w *workerPool) Register(args ...Process) {
	// Create control panel for each process and make process stat for each of them.
	for _, p := range args {
		ctx, cancel := context.WithCancel(context.Background())
		w.controlPanel.put(p.PID(), &processContext{
			ctx:    ctx,
			cancel: cancel,
		})
		w.processes.put(p.PID(), ProcessStats{
			Process: p,
			Status:  process.Waiting,
		})
	}

	// Publish processes to the queue.
	go func(args ...Process) {
		for i := range args {
			w.mutex.Lock()
			if w.isClosed {
				break
			}
			w.queue <- args[i]
			w.mutex.Unlock()
		}
	}(args...)
}

// Close stops a running pool. It returns an error if the pool is not running.
// Close waits for all workers to finish their current job and then closes the
// pool.
func (w *workerPool) Close() error {
	if w.status != pool.Running {
		return errors.New("pool is not running, status " + w.status.String())
	}

	w.mutex.Lock()
	w.isClosed = true
	close(w.queue)
	w.mutex.Unlock()

	w.wg.Wait()
	w.status = pool.Closed

	return nil
}

// WorkerList returns the list of worker names of the pool.
func (w *workerPool) WorkerList() []WorkerName {
	return w.workers
}

// Kill cancel a process before it starts.
func (w *workerPool) Kill(pid PID) {
	w.controlPanel.get(pid).cancel()
}

// Monitor returns pool monitor.
func (w *workerPool) Monitor() Monitor {
	return w
}

// String returns the string value of process id.
func (p PID) String() string {
	return string(p)
}

// PoolStatus returns pool status
func (w *workerPool) PoolStatus() pool.Status {
	return w.status
}

// Error returns process's error by process id.
func (w *workerPool) Error(pid PID) error {
	return w.processes.get(pid).err
}

// WorkerStatus returns worker status. It accepts worker name as input.
func (w *workerPool) WorkerStatus(name WorkerName) worker.Status {
	return w.workersStats.get(name)
}

// ProcessStatus returns process stats. It accepts process id as input.
func (w *workerPool) ProcessStats(pid PID) ProcessStats {
	return w.processes.get(pid)
}
