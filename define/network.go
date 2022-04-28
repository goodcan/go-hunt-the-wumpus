/*
 * @Time     : 2022/4/9 11:38
 * @Author   : cancan
 * @File     : network.go
 * @Function :
 */

package define

type NetworkData struct {
	Url              string            `json:"url"`                         // 路由
	RoleName         string            `json:"role_name"`                   // 角色名称
	RolePosition     [2]int            `json:"role_position"`               // 玩家坐标
	PartnerPositions map[string][2]int `json:"partner_positions,omitempty"` // 队友坐标
	RoleDirection    int               `json:"role_direction"`              // 玩家朝向
	RoleStatus       int               `json:"role_status"`                 // 角色状态
	MonsterPosition  [2]int            `json:"monster_position"`            // 怪物坐标
	MapLength        int               `json:"map_length"`                  // 地图长度
	MapWidth         int               `json:"map_width"`                   // 地图宽度
	RoleInput        int               `json:"input"`                       // 角色输入
}
