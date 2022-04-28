/*
 * @Time     : 2022/4/10 11:03
 * @Author   : cancan
 * @File     : factory.go
 * @Function :
 */

package storage

import (
	"fmt"
	"go-hunt-the-wumpus/define"
)

func Create(storageTeyp string) define.IStorage {
	switch storageTeyp {
	case "text":
		return &Text{}
	default:
		panic(fmt.Sprintf("storageTeyp %v not support", storageTeyp))
	}
}
