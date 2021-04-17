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
 * Time: 9:53 PM
 *
 * Description:
 *
 */

package gowl

import (
	"context"
	"github.com/hamed-yousefi/gowl/status/worker"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Test controlPanelMap put and get functions
func TestControlPanelMap(t *testing.T) {
	cp := new(controlPanelMap)
	ctx, cancel := context.WithCancel(context.Background())
	pc := &processContext{
		ctx:    ctx,
		cancel: cancel,
	}
	cp.put("p-11", pc)

	a := assert.New(t)
	a.Equal(pc, cp.get("p-11"))
}

// Test workerStatsMap put and get functions
func TestWorkerStatsMap(t *testing.T) {
	ws := new(workerStatsMap)
	ws.put("w1", worker.Busy)

	a := assert.New(t)
	a.Equal(worker.Busy, ws.get("w1"))
}

// Test processStatusMap put and get functions
func TestProcessStatusMap(t *testing.T) {
	ps := new(processStatusMap)
	p := ProcessStats{
		WorkerName: "w1",
		StartedAt:  time.Now(),
	}
	ps.put("p-11", p)

	a := assert.New(t)
	a.Equal(p, ps.get("p-11"))
}
