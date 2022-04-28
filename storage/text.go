/*
 * @Time     : 2022/4/5 22:24
 * @Author   : cancan
 * @File     : text.go.go
 * @Function :
 */

package storage

import (
	"fmt"
	"go-hunt-the-wumpus/util"
	"strconv"
	"strings"
)

const dataDir = "./data"

type Text struct {
}

func (t *Text) ReadRoleInfo(roleName string) (rolePosition [2]int, monsterPosition [2]int, roleDirection int) {
	return decodeData(util.ReadFile(dataDir, makeFilename(roleName)))
}

func (t *Text) SaveRoleInfo(roleName string, rolePosition [2]int, monsterPosition [2]int, roleDirection int) {
	util.WriteFile(dataDir, makeFilename(roleName), encodeData(rolePosition, monsterPosition, roleDirection))
}

func makeFilename(roleName string) string {
	return fmt.Sprintf("%v.txt", roleName)
}

func encodeData(rolePosition [2]int, monsterPosition [2]int, roleDirection int) []byte {
	return []byte(fmt.Sprintf("%d-%d-%d-%d-%d", rolePosition[0], rolePosition[1], monsterPosition[0], monsterPosition[1], roleDirection))
}

func decodeData(data []byte) (rolePosition [2]int, monsterPosition [2]int, roleDirection int) {
	if len(data) == 0 {
		return
	}
	ss := strings.Split(string(data), "-")
	rolePosition[0], _ = strconv.Atoi(ss[0])
	rolePosition[1], _ = strconv.Atoi(ss[1])
	monsterPosition[0], _ = strconv.Atoi(ss[2])
	monsterPosition[1], _ = strconv.Atoi(ss[3])
	roleDirection, _ = strconv.Atoi(ss[4])
	return
}
