/*
 * @Time     : 2022/4/8 23:01
 * @Author   : cancan
 * @File     : game.go
 * @Function :
 */

package define

type IUI interface {
	SetMap(length, width int)          // 初始化，设置地图大小
	ShowGame(positions map[[2]int]int) // 显示地图信息
	ReadInput() chan int               // 读取输入
	ShowWarn()                         // 显示警告
	ShowFail()                         // 显示失败
	ShowWin()                          // 显示胜利
	ShowQuit()                         // 显示退出
	ShowInputError(input int)          // 显示输入错误
}

type IStorage interface {
	ReadRoleInfo(roleName string) (rolePosition [2]int, monsterPosition [2]int, roleDirection int)
	SaveRoleInfo(roleName string, rolePosition [2]int, monsterPosition [2]int, roleDirection int)
}

type IInput interface {
	Input() chan string
}

const (
	DoUp          = iota + 1 // 向上移动
	DoDown                   // 向下移动
	DoLeft                   // 向左移动
	DoRight                  // 向右移动
	DoShoot                  // 射击
	DoShowMonster            // 作弊，显示一下怪物
	DoQuit                   // 退出
)

const (
	ShowMe      = iota + 1 // 显示自己
	ShowMonster            // 显示怪物
	ShowPartner            // 显示队友
)

const (
	StatusSafe = iota + 1 // 安全的
	StatusWarn            // 警告
	StatusWin             // 胜利
	StatusFail            // 失败
	StatusGod             // 上帝模式
)
