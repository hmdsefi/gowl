# Gowl

[![Build Status](https://travis-ci.com/hamed-yousefi/gowl.svg?branch=master)](https://travis-ci.com/hamed-yousefi/gowl)
[![codecov](https://codecov.io/gh/hamed-yousefi/gowl/branch/master/graph/badge.svg?token=1TYYX8IBR0)](https://codecov.io/gh/hamed-yousefi/gowl)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamed-yousefi/gowl)](https://goreportcard.com/report/github.com/hamed-yousefi/gowl)
[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B24403%2Fgithub.com%2Fhamed-yousefi%2Fgowl.svg?type=shield)](https://app.fossa.com/projects/custom%2B24403%2Fgithub.com%2Fhamed-yousefi%2Fgowl?ref=badge_shield)
<div  align="center"><img src="https://github.com/hamed-yousefi/gowl/blob/master/docs/images/process-pool.png" width="450" ></div>
Gowl is a process management and process monitoring tool at once.
An infinite worker pool gives you the ability to control the pool and processes
and monitor their status.

## Table of Contents

* [Install](#Install)
* [How to use](#How-to-use)
    * [Pool](#Pooling)
        * Start
        * Register process
        * Kill process
        * Close
    * [Monitor](#Monitor)
* [License](#License)

## Installing

Using Gowl is easy. First, use `go get` to install the latest version of the library. This command will install
the `gowl` along with library and its dependencies:

```shell
go get -u github.com/hamed-yousefi/gowl
```

Next, include Gowl in your application:

```go
import "github.com/hamed-yousefi/gowl"
```

## How to use

Gowl has three main parts. Process, Pool, and Monitor. The process is the smallest part of this project. The process is
the part of code that the developer must implement. To do that, Gowl provides an interface to inject outside code into
the pool. The process interface is as follows:

```go
Process interface {
Start() error
Name() string
PID() PID
}
```

The process interface has three methods. The Start function contains the user codes, and the pool workers use this
function to run the process. The Name function returns the process name, and the monitor uses this function to provide
reports. The PID function returns process id. The process id is unique in the entire pool, and it will use by the pool
and monitor.

Let's take a look at an example:

```go
Document struct {
content string
hash string
}

func (d *Document) Start() error {
hasher := sha1.New()
hasher.Write(bv)
h.hash = base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (d *Document) Name() string {
return "hashing-process"
}

func (d *Document) PID() PID {
return "p-1"
}

func (d *Document) Hash() string {
return h.hash
}
```

As you can see, in this example, `Document` implements the Process interface. So now we can register it into the pool.

### Pool

Creating Gowl pool is very easy. You must use the `NewPool(size int)`
function and pass the pool size to this function. Pool size indicates the worker numbers in and the underlying queue
size that workers consume process from it. Look at the following example:

```go
pool := gowl.NewPool(4)
```

In this example, Gowl will create a new instance of a Pool object with four workers and an underlying queue with the
size of four.

#### Start

To start the Gowl, you must call the `Start()` method of the pool object. It will begin to create the workers, and
workers start listening to the queue to consume process.

#### Register process

To register processes to the pool, you must use the `Register(args ...process)`
method. Pass the processes to the register method, and it will create a new publisher to publish the process list to the
queue. You can call multiple times when Gowl pool is running.

#### Kill process

One of the most remarkable features of Gowl is the ability to control the process after registered it into the pool. You
can kill a process before any worker runs it. Killing a process is simple, and you need the process id to do it.

```go
pool.Kill(PID("p-909"))
```

#### Close

Gowl is an infinite worker pool. However, you should have control over the pool and decide when you want to start it,
register a new process on it, kill a process, and `close` the pool and terminate the workers. Gowl gives you this option
to close the pool by the `Close()` method of the Pool object.

## Monitor

Every process management tool needs a monitoring system to expose the internal stats to the outside world. Gowl gives
you a monitoring API to see processes and workers stats.

You can get the Monitor instance by calling the `Monitor()` method of the Pool. The monitor object is as follows:

```go
Monitor interface {
PoolStatus() pool.Status
Error(PID) error
WorkerList() []WorkerName
WorkerStatus(name WorkerName) worker.Status
ProcessStats(pid PID) ProcessStats
}
```

The Monitor gives you this opportunity to get the Pool status, process error, worker list, worker status, and process
stats. Wis Monitor API, you can create your monitoring app with ease. The following example is using Monitor API to
present the stats in the console in real-time.

![process-monitoring](https://github.com/hamed-yousefi/gowl/blob/master/docs/images/process-monitoring.gif)

Also, you can use the Monitor API to show worker status in the console:

![worker-monitoring](https://github.com/hamed-yousefi/gowl/blob/master/docs/images/worker-monitoring.gif)

## License

MIT License, please see [LICENSE](https://github.com/hamed-yousefi/gowl/blob/master/LICENSE) for details.