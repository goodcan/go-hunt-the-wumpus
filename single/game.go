/*
 * @Time     : 2022/4/5 22:17
 * @Author   : cancan
 * @File     : game.go
 * @Function :
 */

package single

import (
	"go-hunt-the-wumpus/define"
	"go-hunt-the-wumpus/util"
	"math/rand"
	"os"
	"time"
)

type Game struct {
	roleName        string // 玩家名字
	rolePosition    [2]int // 玩家坐标
	roleDirection   int    // 玩家朝向
	monsterPosition [2]int // 怪物坐标
	mapLength       int    // 地图长度
	mapWidth        int    // 地图宽度

	UI      define.IUI
	Storage define.IStorage
}

func (g *Game) Init(roleName string, mapLength, mapWidth int) {
	g.roleName = roleName

	g.mapLength = mapLength
	g.mapWidth = mapWidth

	g.rolePosition, g.monsterPosition, g.roleDirection = g.Storage.ReadRoleInfo(g.roleName)
	g.rolePosition[0] = util.Min(g.rolePosition[0], g.mapLength-1)
	g.rolePosition[1] = util.Min(g.rolePosition[1], g.mapWidth-1)

	// 初始化或者上一局结束怪物位置重新设置
	if util.PositionIsEqual(g.rolePosition, g.monsterPosition) {
		g.reset()
	}

	g.UI.SetMap(g.mapLength, g.mapWidth)
	g.UI.ShowGame(map[[2]int]int{g.rolePosition: define.ShowMe})
}

func (g *Game) reset() {
	g.rolePosition = [2]int{}
	g.roleDirection = define.DoRight

	rand.Seed(time.Now().Unix())
	g.monsterPosition[0] = 1 + rand.Intn(g.mapLength-1)
	g.monsterPosition[1] = 1 + rand.Intn(g.mapWidth-1)
}

func (g *Game) Start() {
	for input := range g.UI.ReadInput() {
		status := g.handlerInput(input)
		g.handlerStatus(status)
	}
}

func (g *Game) handlerInput(input int) int {
	switch input {
	case define.DoUp: // 上
		g.rolePosition[1] = util.Max(g.rolePosition[1]-1, 0)
		g.roleDirection = define.DoUp
	case define.DoDown: // 下
		g.rolePosition[1] = util.Min(g.rolePosition[1]+1, g.mapWidth-1)
		g.roleDirection = define.DoDown
	case define.DoLeft: // 左
		g.rolePosition[0] = util.Max(g.rolePosition[0]-1, 0)
		g.roleDirection = define.DoLeft
	case define.DoRight: // 右
		g.rolePosition[0] = util.Min(g.rolePosition[0]+1, g.mapLength-1)
		g.roleDirection = define.DoRight
	case define.DoShoot: // 射击
		switch g.roleDirection {
		case define.DoUp:
			if g.rolePosition[0] == g.monsterPosition[0] && g.rolePosition[1] > g.monsterPosition[1] {
				return define.StatusWin
			}
		case define.DoDown:
			if g.rolePosition[0] == g.monsterPosition[0] && g.rolePosition[1] < g.monsterPosition[1] {
				return define.StatusWin
			}
		case define.DoLeft:
			if g.rolePosition[0] > g.monsterPosition[0] && g.rolePosition[1] == g.monsterPosition[1] {
				return define.StatusWin
			}
		case define.DoRight:
			if g.rolePosition[0] < g.monsterPosition[0] && g.rolePosition[1] == g.monsterPosition[1] {
				return define.StatusWin
			}
		}
	case define.DoShowMonster:
		return define.StatusGod
	case define.DoQuit:
		g.UI.ShowQuit()
		os.Exit(0)
	default:
		g.UI.ShowInputError(input)
	}

	if util.PositionIsAdjoin(g.rolePosition, g.monsterPosition) {
		return define.StatusWarn
	}

	if util.PositionIsEqual(g.rolePosition, g.monsterPosition) {
		return define.StatusFail
	}

	return define.StatusSafe
}

func (g *Game) handlerStatus(status int) {
	switch status {
	case define.StatusSafe:
		g.UI.ShowGame(map[[2]int]int{g.rolePosition: define.ShowMe})
	case define.StatusWin:
		g.reset()
		g.UI.ShowGame(map[[2]int]int{g.rolePosition: define.ShowMe})
		g.UI.ShowWin()
	case define.StatusWarn:
		g.UI.ShowGame(map[[2]int]int{g.rolePosition: define.ShowMe})
		g.UI.ShowWarn()
	case define.StatusFail:
		g.reset()
		g.UI.ShowGame(map[[2]int]int{g.rolePosition: define.ShowMe})
		g.UI.ShowFail()
	case define.StatusGod:
		g.UI.ShowGame(map[[2]int]int{
			g.rolePosition:    define.ShowMe,
			g.monsterPosition: define.ShowMonster,
		})
	}

	g.Storage.SaveRoleInfo(g.roleName, g.rolePosition, g.monsterPosition, g.roleDirection)
}
