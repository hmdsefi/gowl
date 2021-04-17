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
 * Time: 2:51 PM
 *
 * Description:
 *
 */

package worker

const (
	Waiting Status = iota
	Busy
)

var (
	status2String = map[Status]string{
		Waiting: "Waiting",
		Busy:    "Busy",
	}
)

type (
	Status int
)

func (s Status) String() string {
	return status2String[s]
}
