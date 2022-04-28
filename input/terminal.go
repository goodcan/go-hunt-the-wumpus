/*
 * @Time     : 2022/4/5 22:43
 * @Author   : cancan
 * @File     : terminal.go
 * @Function :
 */

package input

import "fmt"

type Terminal struct {
}

func (t *Terminal) Input() chan string {
	c := make(chan string)
	go func() {
		var input string
		for {
			_, _ = fmt.Scanln(&input)
			c <- input
		}
	}()
	return c
}
