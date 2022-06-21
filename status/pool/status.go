/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com>.
 *
 * Use of this source code is governed by an MIT-style
 * license that can be found in the LICENSE file.
 *
 * Created by IntelliJ IDEA.
 * User: Hamed Yousefi
 * Email: hdyousefi@gmail.com
 * Date: 4/17/21
 * Time: 1:30 PM
 *
 * Description:
 *
 */

package pool

const (
	// Created is a pool state after pool has been created and before it starts.
	Created Status = iota
	// Running is a pool state when the pool started by Start() function.
	Running
	// Closed is a pool state when the pool stopped by Close() function.
	Closed
)

var (
	status2string = map[Status]string{
		Created: "Created",
		Running: "Running",
		Closed:  "Closed",
	}
)

type (
	// Status represents pool current state.
	Status int
)

// String returns string value of pool state.
func (p Status) String() string {
	return status2string[p]
}
