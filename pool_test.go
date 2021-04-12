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
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestTimeout(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func(ctx context.Context) {
		defer wg.Done()
		finished := make(chan struct{})
		for {
			select {
			case <-ctx.Done():
				fmt.Println("process cancelled")
				return
			case <- finished:
				return
			default:
				fmt.Println("default")
				fmt.Println("wait for 3 seconds.")
				//time.Sleep(3 * time.Second)
				finished <- struct{}{}
			}

		}
	}(ctx)

	//time.Sleep(1 * time.Second)
	defer cancel()
	wg.Wait()
}
