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
	poolStatus2string = map[Status]string{
		Created: "Created",
		Running: "Running",
		Closed:  "Closed",
	}
)

type (
	Status int
)

func (p Status) String() string {
	return poolStatus2string[p]
}
