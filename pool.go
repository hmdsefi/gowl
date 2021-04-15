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
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	PWaiting ProcessStatus = iota
	Running
	Succeeded
	Failed
	Killed
)

const (
	WWaiting WorkerStatus = iota
	Busy
)

const (
	Created PoolStatus = iota
	PRunning
	Closed

	defaultWorkerName = "W%d"
)

var (
	singleton = new(sync.Once)

	processStatus2String = map[ProcessStatus]string{
		PWaiting:  "Waiting",
		Running:   "Running",
		Succeeded: "Succeeded",
		Failed:    "Failed",
		Killed:    "Killed",
	}

	workerStatus2String = map[WorkerStatus]string{
		WWaiting: "Waiting",
		Busy:     "Busy",
	}

	poolStatus2string = map[PoolStatus]string{
		Created:  "Created",
		PRunning: "Running",
		Closed:   "Closed",
	}
)

type (
	ProcessStatus int
	WorkerStatus  int
	PoolStatus    int
	WorkerName    string
	PID           string

	Process interface {
		Start() error
		Name() string
		PID() PID
	}

	Pool interface {
		Start()
		Register(p ...Process)
		Close()
		Kill(pid PID)
		WorkerList() []WorkerName
	}

	Monitor interface {
		PoolStatus() PoolStatus
		Error(PID) error
		WorkerStatus(name WorkerName) WorkerStatus
		ProcessStatus(pid PID) ProcessStats
	}

	ProcessStats struct {
		WorkerName WorkerName
		process    Process
		Status     ProcessStatus
		StartedAt  time.Time
		FinishedAt time.Time
		err        error
	}

	workerPool struct {
		status       PoolStatus
		size         int
		queue        chan Process
		wg           *sync.WaitGroup
		processes    *processStatusMap
		workers      []WorkerName
		workersStats *workerStatsMap
		controlPanel *controlPanelMap
		mutex        *sync.Mutex
		wnMutex      *sync.Mutex
		isClosed     bool
	}
)

func NewPool(size int) *workerPool {
	return &workerPool{
		status:       Created,
		size:         size,
		queue:        make(chan Process, size),
		workers:      []WorkerName{},
		processes:    new(processStatusMap),
		workersStats: new(workerStatsMap),
		controlPanel: new(controlPanelMap),
		mutex:        new(sync.Mutex),
		wnMutex:      new(sync.Mutex),
		wg:           new(sync.WaitGroup),
	}
}

func (w *workerPool) Start() {

	if w.status == PRunning {
		log.Printf("pool is already started once, status: %v", w.status)
	}

	singleton.Do(func() {
		w.status = PRunning
		go w.run()
	})
}

func (w *workerPool) run() {

	for i := 0; i < w.size; i++ {
		w.wg.Add(1)

		go func(n int) {
			defer w.wg.Done()
			w.wnMutex.Lock()
			wn := WorkerName(fmt.Sprintf(defaultWorkerName, n))
			w.workers = append(w.workers, wn)
			w.wnMutex.Unlock()

			for p := range w.queue {
				w.workersStats.put(wn, Busy)
				pStats := w.processes.get(p.PID())
				pStats.Status = Running
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
						log.Printf("process with id %s has been killed.\n", p.PID().String())
						stats.Status = Killed
						return
					default:
						if err := p.Start(); err != nil { //nolint:typecheck
							stats.err = err
							stats.Status = Failed
						} else {
							stats.Status = Succeeded
						}
						pContext.cancel()
					}
				}()

				wgp.Wait()
				pStats = w.processes.get(p.PID())
				pStats.FinishedAt = time.Now()
				w.processes.put(p.PID(), pStats)
				w.workersStats.put(wn, WWaiting)
			}
		}(i)
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
			Status:  PWaiting,
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

func (w *workerPool) Close() {
	w.mutex.Lock()
	w.isClosed = true
	close(w.queue)
	w.mutex.Unlock()

	w.wg.Wait()
	w.status = Closed
	singleton = new(sync.Once)
}

func (w *workerPool) WorkerList() []WorkerName {
	return w.workers
}

func (w *workerPool) Kill(pid PID) {
	w.controlPanel.get(pid).cancel()
}

func (w *workerPool) Status() PoolStatus {
	return w.status
}

func (p PID) String() string {
	return string(p)
}

func (w *workerPool) PoolStatus() PoolStatus {
	return w.status
}

func (w *workerPool) Error(pid PID) error {
	return w.processes.get(pid).err
}

func (w *workerPool) WorkerStatus(name WorkerName) WorkerStatus {
	return w.workersStats.get(name)
}

func (w *workerPool) ProcessStatus(pid PID) ProcessStats {
	return w.processes.get(pid)
}

func (p ProcessStatus) String() string {
	return processStatus2String[p]
}

func (w WorkerStatus) String() string {
	return workerStatus2String[w]
}

func (p PoolStatus) String() string {
	return poolStatus2string[p]
}
