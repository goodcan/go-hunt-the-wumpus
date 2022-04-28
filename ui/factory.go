/*
 * @Time     : 2022/4/10 09:23
 * @Author   : cancan
 * @File     : factory.go
 * @Function :
 */

package ui

import (
	"fmt"
	"go-hunt-the-wumpus/define"
	"go-hunt-the-wumpus/input"
)

func Create(language, inputType string) define.IUI {
	switch language {
	case "cn":
		return &CN{Input: input.Create(inputType)}
	default:
		panic(fmt.Sprintf("language %v not support", language))
	}
}
