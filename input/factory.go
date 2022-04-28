/*
 * @Time     : 2022/4/10 10:57
 * @Author   : cancan
 * @File     : factory.go
 * @Function :
 */

package input

import (
	"fmt"
	"go-hunt-the-wumpus/define"
)

func Create(inputType string) define.IInput {
	switch inputType {
	case "terminal":
		return &Terminal{}
	default:
		panic(fmt.Sprintf("inputType %v not support", inputType))
	}
}
