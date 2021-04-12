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
	Waiting ProcessStatus = iota
	Running
	Succeeded
	Failed
	Killed

	ReadyWaiting WorkerStatus = iota
	Busy

	Created PoolStatus = iota
	PRunning
	Closed

	defaultWorkerName = "W%d"
)

var (
	singleton = new(sync.Once)
)

type (
	ProcessStatus int
	WorkerStatus  int
	PoolStatus    int
	WorkerName    string
	PID string

	Process interface {
		Start() error
		Timeout() time.Duration
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
		ProcessStatus(pid PID)
	}

	ProcessStats struct {
		WorkerName WorkerName
		process    Process
		Status     ProcessStatus
		StartedAt  time.Time
		FinishedAt time.Time
	}

	workerPool struct {
		status       PoolStatus
		size         int
		c            chan Process
		processes    map[PID]*ProcessStats
		err          map[PID]error
		workers      []WorkerName
		workersStats map[WorkerName]WorkerStatus
		controlPanel map[PID]context.CancelFunc
	}
)

func NewPool(size int) *workerPool {
	return &workerPool{
		status:       Created,
		size:         size,
		c:            make(chan Process, size),
		processes:    make(map[PID]*ProcessStats),
		err:          make(map[PID]error),
		workers:      []WorkerName{},
		workersStats: make(map[WorkerName]WorkerStatus),
		controlPanel: make(map[PID]context.CancelFunc),
	}
}

func (w *workerPool) Start() {

	if w.status == PRunning || w.status == Closed {
		log.Printf("pool is already started once, status: %v", w.status)
	}

	singleton.Do(func() {
		w.status = PRunning
		go w.run()
	})
}

func (w *workerPool) run() {
	wg := new(sync.WaitGroup)

	for i := 0; i < w.size; i++ {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()
			wn := WorkerName(fmt.Sprintf(defaultWorkerName, n))
			w.workers = append(w.workers, wn)

			for p := range w.c {
				w.workersStats[wn] = Busy
				w.processes[p.PID()].Status = Running
				w.processes[p.PID()].StartedAt = time.Now()
				w.processes[p.PID()].WorkerName = wn

				ctx, cancel := context.WithCancel(context.Background())
				w.controlPanel[p.PID()] = cancel

				wgp := new(sync.WaitGroup)
				wgp.Add(1)

				go func() {
					defer wgp.Done()

					select {
					case <- ctx.Done():
						log.Printf("process with id %s has been killed.\n", p.PID().String())
						w.processes[p.PID()].Status = Killed
						return
					default:
						if err := p.Start(); err != nil { //nolint:typecheck
							w.err[p.PID()] = err
							w.processes[p.PID()].Status = Failed
						} else {
							w.processes[p.PID()].Status = Succeeded
						}
						w.controlPanel[p.PID()]()
					}
				}()

				wgp.Wait()
				w.processes[p.PID()].FinishedAt = time.Now()
				w.workersStats[wn] = ReadyWaiting
			}
		}(i)
	}

	wg.Wait()
}

func (w *workerPool) Register(args ...Process) {
	for _, p := range args {

		w.processes[p.PID()] = &ProcessStats{
			process: p,
			Status:  Waiting,
		}
	}
	go func(args ...Process) {
		for pid := range w.processes {
			w.c <- w.processes[pid].process
		}
	}(args...)
}

func (w *workerPool) Close() {
	close(w.c)
	w.status = Closed
}

func (w *workerPool) WorkerList() []WorkerName {
	return w.workers
}

func (w *workerPool) Kill(pid PID) {
	w.controlPanel[pid]()
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
	return w.err[pid]
}

func (w *workerPool) WorkerStatus(name WorkerName) WorkerStatus {
	return w.workersStats[name]
}

func (w *workerPool) ProcessStatus(pid PID) ProcessStats {
	return *w.processes[pid]
}
