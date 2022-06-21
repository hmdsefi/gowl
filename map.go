/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/13/21
 * Time: 3:55 PM
 *
 * Description:
 *
 */

package gowl

import (
	"context"
	"sync"

	"github.com/hamed-yousefi/gowl/status/worker"
)

type (
	// controlPanelMap is a thread safe map for controlling processes. It also
	// provides type safety.
	// 		Key: PID
	// 		Value: processContext
	controlPanelMap struct {
		internal sync.Map
	}

	// workerStatsMap is a thread safe map for controlling processes. It also
	// provides type safety.
	// 		Key: WorkerName
	// 		Value: worker.Status
	workerStatsMap struct {
		internal sync.Map
	}

	// processStatusMap is a thread safe map for controlling processes. It also
	// provides type safety.
	// 		Key: PID
	// 		Value: ProcessStats
	processStatusMap struct {
		internal sync.Map
	}

	// processContext represents a cancellation context by holding a context and
	// a cancel function.
	processContext struct {
		ctx    context.Context
		cancel context.CancelFunc
	}
)

func (c *controlPanelMap) put(pid PID, pc *processContext) {
	c.internal.Store(pid, pc)
}

func (c *controlPanelMap) get(pid PID) *processContext {
	in, _ := c.internal.Load(pid)
	cancel, _ := in.(*processContext)
	return cancel
}

func (c *workerStatsMap) put(name WorkerName, status worker.Status) {
	c.internal.Store(name, status)
}

func (c *workerStatsMap) get(name WorkerName) worker.Status {
	in, _ := c.internal.Load(name)
	status, _ := in.(worker.Status)
	return status
}

func (c *processStatusMap) put(pid PID, stats ProcessStats) {
	c.internal.Store(pid, stats)
}

func (c *processStatusMap) get(pid PID) ProcessStats {
	in, _ := c.internal.Load(pid)
	stats, _ := in.(ProcessStats)
	return stats
}
