/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package gowl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hamed-yousefi/gowl/status/worker"
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
