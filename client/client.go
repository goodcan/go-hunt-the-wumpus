/*
 * @Time     : 2022/4/9 11:16
 * @Author   : cancan
 * @File     : client.go
 * @Function :
 */

package client

import (
	"encoding/json"
	"fmt"
	"go-hunt-the-wumpus/define"
	"net"
	"os"
)

type Client struct {
	serverAddr string
	roleName   string

	UI define.IUI
}

func (c *Client) Init(serverAddr, roleName string) {
	c.serverAddr = serverAddr
	c.roleName = roleName
}

func (c *Client) Start() {
	conn, err := net.Dial("tcp", c.serverAddr)
	if err != nil {
		fmt.Println("Connect to TCP server failed ,err:", err)
		return
	}

	go func() {
		defer func() {
			c.UI.ShowQuit()
			os.Exit(0)
		}()
		for {
			var buf [1024]byte
			n, err := conn.Read(buf[:])
			if err != nil {
				//fmt.Println("Read from tcp server failed,err:", err)
				break
			}
			var data define.NetworkData
			_ = json.Unmarshal(buf[:n], &data)
			//fmt.Printf("Recived from client,data:%v\n", data)

			c.UI.SetMap(data.MapLength, data.MapWidth)
			showPositions := make(map[[2]int]int)
			for _, v := range data.PartnerPositions {
				showPositions[v] = define.ShowPartner
			}

			switch data.RoleStatus {
			case define.StatusSafe:
				showPositions[data.RolePosition] = define.ShowMe
				c.UI.ShowGame(showPositions)
			case define.StatusWin:
				showPositions[data.RolePosition] = define.ShowMe
				//showPositions[data.MonsterPosition] = define.ShowMonster
				c.UI.ShowGame(showPositions)
				c.UI.ShowWin()
			case define.StatusWarn:
				showPositions[data.RolePosition] = define.ShowMe
				c.UI.ShowGame(showPositions)
				c.UI.ShowWarn()
			case define.StatusFail:
				showPositions[data.RolePosition] = define.ShowMe
				c.UI.ShowGame(showPositions)
				c.UI.ShowFail()
			case define.StatusGod:
				showPositions[data.RolePosition] = define.ShowMe
				showPositions[data.MonsterPosition] = define.ShowMonster
				c.UI.ShowGame(showPositions)
			}
		}
	}()

	data, _ := json.Marshal(define.NetworkData{
		Url:      "/login",
		RoleName: c.roleName,
	})
	_, _ = conn.Write(data)

	for input := range c.UI.ReadInput() {
		data, _ = json.Marshal(define.NetworkData{
			Url:       "/action",
			RoleName:  c.roleName,
			RoleInput: input,
		})
		_, _ = conn.Write(data)
	}
}
