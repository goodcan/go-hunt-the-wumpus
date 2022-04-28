/*
 * @Time     : 2022/4/5 22:17
 * @Author   : cancan
 * @File     : server.go
 * @Function :
 */

package server

import (
	"encoding/json"
	"fmt"
	"go-hunt-the-wumpus/define"
	"go-hunt-the-wumpus/util"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Role struct {
	name      string // 玩家名字
	position  [2]int // 玩家坐标
	direction int    // 玩家朝向
	status    int    // 状态

	conn net.Conn // 网络连接
}

type Server struct {
	addr  string
	roles map[string]*Role

	monsterPosition [2]int // 怪物坐标
	mapLength       int    // 地图长度
	mapWidth        int    // 地图宽度

	Storage define.IStorage
	sync.RWMutex
}

func (s *Server) Init(addr string, mapLength, mapWidth int) {
	s.addr = addr
	s.roles = make(map[string]*Role)

	s.mapLength = mapLength
	s.mapWidth = mapWidth

	s.createMonster()
}

func (s *Server) createMonster() {
	rand.Seed(time.Now().Unix())
	s.monsterPosition[0] = 1 + rand.Intn(s.mapLength-1)
	s.monsterPosition[1] = 1 + rand.Intn(s.mapWidth-1)
}

func (s *Server) reset() {
	s.Lock()
	defer s.Unlock()

	s.createMonster()
	for _, role := range s.roles {
		role.position = [2]int{}
		role.direction = define.DoRight
		role.status = define.StatusSafe
		s.Storage.SaveRoleInfo(role.name, role.position, s.monsterPosition, role.direction)
	}

}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		fmt.Println("Listen tcp server failed,err:", err)
		return
	}

	fmt.Printf("Start Server OK in add %v \n", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listen.Accept failed,err:", err)
			continue
		}

		go s.onConnect(&Role{conn: conn})
	}
}

func (s *Server) onConnect(role *Role) {
	defer func() {
		s.delRole(role)
		_ = role.conn.Close()
		s.broadcast()
	}()

	for {
		var buf [1024]byte
		n, err := role.conn.Read(buf[:])
		if err != nil {
			//fmt.Println("Read from tcp server failed,err:", err)
			break
		}
		var data define.NetworkData
		_ = json.Unmarshal(buf[:n], &data)
		//fmt.Printf("Recived from client,data:%v \n", data)

		if data.RoleInput == define.DoQuit {
			break
		}

		switch data.Url {
		case "/login":
			role.name = data.RoleName
			s.addRole(role)
			s.broadcast()
		case "/action":
			status := s.handlerInput(data.RoleName, data.RoleInput)
			s.handlerStatus(data.RoleName, status)

			switch role.status {
			case define.StatusFail:
				role.position = [2]int{}
				role.direction = define.DoRight
				s.broadcast()
				role.status = define.StatusSafe
			case define.StatusWin:
				s.reset()
				role.status = define.StatusWin
				s.broadcast()
				role.status = define.StatusSafe
			case define.StatusGod:
				s.broadcast()
				role.status = define.StatusSafe
				if util.PositionIsAdjoin(role.position, s.monsterPosition) {
					role.status = define.StatusWarn
				}
			default:
				s.broadcast()
			}

		}

		fmt.Printf("role %v url %v input %v ip %v \n", role.name, data.Url, data.RoleInput, role.conn.RemoteAddr().String())
	}
}

func (s *Server) broadcast() {
	s.Lock()
	defer s.Unlock()

	totalPositions := make(map[string][2]int)
	for _, role := range s.roles {
		totalPositions[role.name] = role.position
	}
	wg := sync.WaitGroup{}
	wg.Add(len(s.roles))

	for _, role := range s.roles {
		go func(r *Role) {
			partnerPositions := make(map[string][2]int)
			for k, v := range totalPositions {
				if k == r.name {
					continue
				}
				partnerPositions[k] = v
			}
			data, _ := json.Marshal(define.NetworkData{
				RolePosition:     r.position,
				RoleStatus:       r.status,
				MonsterPosition:  s.monsterPosition,
				PartnerPositions: partnerPositions,
				MapLength:        s.mapLength,
				MapWidth:         s.mapWidth,
			})
			_, _ = r.conn.Write(data)
			wg.Done()
		}(role)
	}

	wg.Wait()
}

func (s *Server) addRole(role *Role) {
	role.position, _, role.direction = s.Storage.ReadRoleInfo(role.name)
	role.position[0] = util.Min(role.position[0], s.mapLength-1)
	role.position[1] = util.Min(role.position[1], s.mapWidth-1)

	s.Lock()
	defer s.Unlock()

	s.roles[role.name] = role
	if util.PositionIsEqual(role.position, s.monsterPosition) {
		role.position = [2]int{}
		role.direction = define.DoRight
	}

	role.status = define.StatusSafe
	if util.PositionIsAdjoin(role.position, s.monsterPosition) {
		role.status = define.StatusWarn
	}

	fmt.Printf("role %v join \n", role.name)
}

func (s *Server) delRole(role *Role) {
	s.Lock()
	defer s.Unlock()

	delete(s.roles, role.name)
	fmt.Printf("role %v leave \n", role.name)
}

func (s *Server) handlerInput(roleName string, input int) int {
	s.Lock()
	defer s.Unlock()

	role := s.roles[roleName]

	switch input {
	case define.DoUp: // 上
		role.position[1] = util.Max(role.position[1]-1, 0)
		role.direction = define.DoUp
	case define.DoDown: // 下
		role.position[1] = util.Min(role.position[1]+1, s.mapWidth-1)
		role.direction = define.DoDown
	case define.DoLeft: // 左
		role.position[0] = util.Max(role.position[0]-1, 0)
		role.direction = define.DoLeft
	case define.DoRight: // 右
		role.position[0] = util.Min(role.position[0]+1, s.mapLength-1)
		role.direction = define.DoRight
	case define.DoShoot: // 射击
		switch role.direction {
		case define.DoUp:
			if role.position[0] == s.monsterPosition[0] && role.position[1] > s.monsterPosition[1] {
				return define.StatusWin
			}
		case define.DoDown:
			if role.position[0] == s.monsterPosition[0] && role.position[1] < s.monsterPosition[1] {
				return define.StatusWin
			}
		case define.DoLeft:
			if role.position[0] > s.monsterPosition[0] && role.position[1] == s.monsterPosition[1] {
				return define.StatusWin
			}
		case define.DoRight:
			if role.position[0] < s.monsterPosition[0] && role.position[1] == s.monsterPosition[1] {
				return define.StatusWin
			}
		}
	case define.DoShowMonster:
		return define.StatusGod
	}

	if util.PositionIsAdjoin(role.position, s.monsterPosition) {
		return define.StatusWarn
	}

	if util.PositionIsEqual(role.position, s.monsterPosition) {
		return define.StatusFail
	}

	return define.StatusSafe
}

func (s *Server) handlerStatus(roleName string, status int) {
	s.Lock()
	defer s.Unlock()

	role := s.roles[roleName]
	role.status = status

	s.Storage.SaveRoleInfo(role.name, role.position, s.monsterPosition, role.direction)
}
