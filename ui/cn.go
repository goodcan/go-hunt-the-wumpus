/*
 * @Time     : 2022/4/5 22:29
 * @Author   : cancan
 * @File     : cn.go
 * @Function :
 */

package ui

import (
	"fmt"
	"go-hunt-the-wumpus/define"
	"go-hunt-the-wumpus/util"
	"strings"
)

var inputStr2int = map[string]int{
	"w": define.DoUp,
	"s": define.DoDown,
	"a": define.DoLeft,
	"d": define.DoRight,
	"f": define.DoShoot,
	"m": define.DoShowMonster,
	"q": define.DoQuit,
}

type CN struct {
	length int
	width  int

	Input define.IInput
}

func (c *CN) SetMap(length, width int) {
	c.length = length
	c.width = width
}

func (c *CN) ShowGame(positions map[[2]int]int) {
	util.ClearTerminal()

	fmt.Println("操作：w-上 s-下 a-左 d-右 f-射击 m-提示一下怪物 q-退出")

	fmt.Println()
	for y := 0; y < c.width; y++ {
		row := []string{}
		for x := 0; x < c.length; x++ {
			if pt, ok := positions[[2]int{x, y}]; ok {
				switch pt {
				case define.ShowMe:
					row = append(row, " 我 ")
				case define.ShowMonster:
					row = append(row, " 怪 ")
				case define.ShowPartner:
					row = append(row, " 友 ")
				}
			} else {
				row = append(row, "    ")
			}
		}
		fmt.Println("|" + strings.Join(row, "|") + "|")
	}
	fmt.Println()
}

func (c *CN) ReadInput() chan int {
	ch := make(chan int)
	go func() {
		for input := range c.Input.Input() {
			inputFix, ok := inputStr2int[input]
			if !ok {
				fmt.Printf("输入 %v 错误，请重新输入! \n", input)
				continue
			}
			ch <- inputFix
		}
	}()
	return ch
}

func (c *CN) ShowWarn() {
	fmt.Println("怪物在你身边，请注意！")
}

func (c *CN) ShowFail() {
	fmt.Println("你应经死亡，游戏重新开始！")
}

func (s *CN) ShowWin() {
	fmt.Println("恭喜你，打死了怪物，游戏重新开始！")
}

func (c *CN) ShowInputError(input int) {
	fmt.Printf("驶入 %v 错误，请重新输入! \n", input)
}

func (c *CN) ShowQuit() {
	fmt.Println("再见！")
}
