/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 */

package worker

const (
	// Waiting is a worker state when the worker is waiting to consume a process.
	Waiting Status = iota
	// Busy is a worker state when the worker consumed a process and running it.
	Busy
)

var (
	status2String = map[Status]string{
		Waiting: "Waiting",
		Busy:    "Busy",
	}
)

type (
	// Status represents worker current state.
	Status int
)

// String returns string value of worker state.
func (s Status) String() string {
	return status2String[s]
}
