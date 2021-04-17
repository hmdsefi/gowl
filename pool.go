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
	defaultWorkerName = "W%d"
)

type (
	WorkerName string
	PID        string

	Process interface {
		Start() error
		Name() string
		PID() PID
	}

	Pool interface {
		Start() error
		Register(p ...Process)
		Close() error
		Kill(pid PID)
		WorkerList() []WorkerName
	}

	Monitor interface {
		PoolStatus() pool.Status
		Error(PID) error
		WorkerStatus(name WorkerName) worker.Status
		ProcessStatus(pid PID) ProcessStats
	}

	ProcessStats struct {
		WorkerName WorkerName
		process    Process
		Status     process.Status
		StartedAt  time.Time
		FinishedAt time.Time
		err        error
	}

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

func NewPool(size int) *workerPool {
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

func (w *workerPool) Start() error {

	if w.status == pool.Running {
		return errors.New("unable to start the pool, status: " + w.status.String())
	}

	w.status = pool.Running
	w.run()

	return nil
}

func (w *workerPool) run() {

	for i := 0; i < w.size; i++ {
		w.wg.Add(1)
		wName := WorkerName(fmt.Sprintf(defaultWorkerName, i))
		w.workers = append(w.workers, wName)
		go func(n int, wn WorkerName) {
			defer w.wg.Done()

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
		}(i, wName)
	}
}

func (w *workerPool) Register(args ...Process) {
	for _, p := range args {
		ctx, cancel := context.WithCancel(context.Background())
		w.controlPanel.put(p.PID(), &processContext{
			ctx:    ctx,
			cancel: cancel,
		})
		w.processes.put(p.PID(), ProcessStats{
			process: p,
			Status:  process.Waiting,
		})
	}
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

func (w *workerPool) WorkerList() []WorkerName {
	return w.workers
}

func (w *workerPool) Kill(pid PID) {
	w.controlPanel.get(pid).cancel()
}

func (p PID) String() string {
	return string(p)
}

func (w *workerPool) PoolStatus() pool.Status {
	return w.status
}

func (w *workerPool) Error(pid PID) error {
	return w.processes.get(pid).err
}

func (w *workerPool) WorkerStatus(name WorkerName) worker.Status {
	return w.workersStats.get(name)
}

func (w *workerPool) ProcessStatus(pid PID) ProcessStats {
	return w.processes.get(pid)
}
