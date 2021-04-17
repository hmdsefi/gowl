# Gowl
[![Build Status](https://travis-ci.com/hamed-yousefi/gowl.svg?branch=master)](https://travis-ci.com/hamed-yousefi/gowl)
[![codecov](https://codecov.io/gh/hamed-yousefi/gowl/branch/master/graph/badge.svg?token=1TYYX8IBR0)](https://codecov.io/gh/hamed-yousefi/gowl)
[![Go Report Card](https://goreportcard.com/badge/github.com/hamed-yousefi/gowl)](https://goreportcard.com/report/github.com/hamed-yousefi/gowl)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fhamed-yousefi%2Fgowl.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fhamed-yousefi%2Fgowl?ref=badge_shield)
<div  align="center"><img src="https://github.com/hamed-yousefi/gowl/blob/master/docs/images/process-pool.png" width="450" ></div>
Gowl is a process management and process monitoring tool at once.
An infinite worker pool gives you the ability to control the pool and processes
and monitor their status.

## Table of Contents

* [Install](#nstall)
* [How to use](#How-to-use)
  * [Pool](#Pooling)
    * Start
    * Register process
    * Kill process
    * Close
  * [Monitor](#Monitoring)    
* [License](#License)

## Installing
Using Gowl is easy. First, use `go get` to install the latest version of the library.
This command will install the `gowl` along with library and its dependencies:
```shell
go get -u github.com/hamed-yousefi/gowl
```
Next, include Gowl in your application:
```go
import "github.com/hamed-yousefi/gowl"
```

## How to use
Gowl has three main parts. Process, Pool, and Monitor. The process is the 
smallest part of this project. The process is the part of code that the
developer must implement. To do that, Gowl provides an interface to inject
outside code into the pool. The process interface is as follows:

```go
Process interface {
   Start() error
   Name() string
   PID() PID
}
```

The process interface has three methods. The Start function contains the
user codes, and the pool workers use this function to run the process. 
The Name function returns the process name, and the monitor uses this
function to provide reports. The PID function returns process id. The 
process id is unique in the entire pool, and it will use by the pool and
monitor.

Let's take a look at an example:
```go
Document struct {
	content string
	hash string
}

func (d *Document) Start() error {
    hasher := sha1.New()
    hasher.Write(bv)
    h.hash= base64.URLEncoding.EncodeToString(hasher.Sum(nil))
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

As you can see, in this example, `Document` implements the Process interface.
So now we can register it into the pool.

### Pool
Creating Gowl pool is very easy. You just need to use `NewPool(size int)` function
and pass the pool size to this function. Pool size indicates the

## License

MIT License, please see [LICENSE](https://github.com/hamed-yousefi/gowl/blob/master/LICENSE) for details.