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
 * Time: 1:42 PM
 *
 * Description:
 *
 */

package process

const (
	Waiting Status = iota
	Running
	Succeeded
	Failed
	Killed
)

var (
	status2String = map[Status]string{
		Waiting:   "Waiting",
		Running:   "Running",
		Succeeded: "Succeeded",
		Failed:    "Failed",
		Killed:    "Killed",
	}
)

type (
	// Status represents process current state.
	Status int
)

// String returns string value of process state.
func (s Status) String() string {
	return status2String[s]
}
