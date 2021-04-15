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
 * Time: 12:45 PM
 *
 * Description:
 *
 */

package gowl

import (
	"fmt"
	"github.com/apoorvam/goterminal"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
	"time"
)

type (
	mockProcess struct {
		name string
		pid  PID
		sleepTime time.Duration
	}
)

func (t mockProcess) Start() error {
	time.Sleep(3 * time.Second)
	fmt.Printf("process with id %v has been started.\n", t.pid)
	return nil
}

func (t mockProcess) Name() string {
	return t.name
}

func (t mockProcess) PID() PID {
	return t.pid
}

func newTestProcess(name string, id int, duration time.Duration) Process {
	return mockProcess{
		name: name,
		pid:  PID("p-" + strconv.Itoa(id)),
		sleepTime: duration,
	}
}

// Close pool before adding all processes to the queue
func TestNewPool(t *testing.T) {
	a := assert.New(t)
	pool := NewPool(2)
	plist := make([]Process, 0)
	for i := 1; i <= 10; i++ {
		plist = append(plist, newTestProcess("p-"+strconv.Itoa(i), i, 3*time.Second))
	}

	a.Equal(Created, pool.status)
	pool.Register(plist...)
	pool.Start()
	a.Equal(PRunning, pool.status)
	time.Sleep(500 * time.Millisecond)
	pool.Close()
	a.Equal(Closed, pool.status)
}

// Four different goroutine will publish processes to the queue
func TestNewPoolMultiPublisher(t *testing.T) {
	a := assert.New(t)
	pool := NewPool(2)
	a.Equal(Created, pool.status)
	pool.Start()
	a.Equal(PRunning, pool.status)
	pool.Register(createProcess(10, 1, 3*time.Second)...)
	pool.Register(createProcess(10, 2, 2*time.Second)...)
	pool.Register(createProcess(10, 3, 1*time.Second)...)
	pool.Register(createProcess(10, 4, 500*time.Millisecond)...)

	time.Sleep(30 * time.Second)
	pool.Close()
	a.Equal(Closed, pool.status)
}

// Use register without input args
func TestNewPoolWithNoRegistry(t *testing.T) {
	a := assert.New(t)
	pool := NewPool(2)
	a.Equal(Created, pool.status)
	pool.Start()
	a.Equal(PRunning, pool.status)
	pool.Register()

	time.Sleep(2 * time.Second)
	pool.Close()
	a.Equal(Closed, pool.status)
}

// Kill a process before it starts
func TestNewPoolKillProcess(t *testing.T) {
	a := assert.New(t)
	pool := NewPool(5)
	a.Equal(Created, pool.status)
	pool.Start()
	a.Equal(PRunning, pool.status)
	pool.Register(createProcess(10, 1, 3*time.Second)...)
	pool.Kill("p-18")
	time.Sleep(7 * time.Second)
	pool.Close()
	a.Equal(Closed, pool.status)
	a.Equal(Killed,pool.ProcessStatus("p-18").Status)
}

func createProcess(n int, g int, d time.Duration) []Process {
	pList := make([]Process, 0)
	for i := 1; i <= n; i++ {
		pList = append(pList, newTestProcess("p-"+strconv.Itoa(i), (g*10)+i, d))
	}
	return pList
}

func monitor(m Monitor) {
	// get an instance of writer
	writer := goterminal.New(os.Stdout)

	for i := 0; i < 100; i = i + 10 {
		// add your text to writer's buffer
		fmt.Fprintf(writer, "Downloading (%d/100) bytes...\n", i)
		// write to terminal
		writer.Print()
		time.Sleep(time.Millisecond * 200)

		// clear the text written by the previous write, so that it can be re-written.
		writer.Clear()
	}

	// reset the writer
	writer.Reset()
	fmt.Println("Download finished!")
}


