/*
 * @Time     : 2022/4/9 18:22
 * @Author   : cancan
 * @File     : common.go
 * @Function :
 */

package util

import (
	"math"
	"os"
	"os/exec"
	"runtime"
)

func Min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func Max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

// 位置相同
func PositionIsEqual(p1, p2 [2]int) bool {
	return p1[0] == p2[0] && p1[1] == p2[1]
}

// 位置相邻
func PositionIsAdjoin(p1, p2 [2]int) bool {
	return math.Sqrt(math.Pow(float64(p1[0])-float64(p2[0]), 2)+math.Pow(float64(p1[1])-float64(p2[1]), 2)) == 1
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	case "linux":
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}

}
