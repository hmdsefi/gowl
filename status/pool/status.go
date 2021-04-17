/**
 * Copyright © 2019 Hamed Yousefi <hdyousefi@gmail.com.com>.
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
	Created Status = iota
	Running
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
