package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建用户
func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		conn:   conn,
		C:      make(chan string),
		server: server,
	}
	//启动监听当前User channel 的msg
	go user.ListenMessage()

	return user
}

//用户上线
func (u *User) Online() {

	u.server.mapLock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.mapLock.Unlock()

	//广播当前用户上线消息
	u.server.Boardcast(u, "已上线")
}

//用户下线
func (u *User) Offline() {
	delete(u.server.OnlineMap, u.Name)
	u.server.Boardcast(u, "下线")
}

func (u *User) SendMsg(msg string) {
	u.C <- msg
}

//用户处理消息
func (u *User) DoMessage(msg string) {
	if msg == "who" {
		u.server.mapLock.Lock()

		num := len(u.server.OnlineMap)
		u.SendMsg(fmt.Sprintf("实际在线:%d人", num))
		for _, user := range u.server.OnlineMap {
			onlineMsg := fmt.Sprintf("[%s]%s:在线...", user.Addr, user.Name)
			u.SendMsg(onlineMsg)
		}

		u.server.mapLock.Unlock()

	} else if len(msg) >= 7 && msg[:7] == "rename|" {
		newName := msg[7:]
		if _, ok := u.server.OnlineMap[newName]; ok {
			u.SendMsg("用户名已经被使用")
		} else {
			u.server.mapLock.Lock()
			delete(u.server.OnlineMap, u.Name)
			u.server.OnlineMap[newName] = u
			u.Name = newName
			u.server.mapLock.Unlock()

			u.SendMsg("您已经更新用户名为:" + newName)
		}
	} else {
		u.server.Boardcast(u, msg)
	}
}

//监听当前User的channel方法,一旦有消息，就发送给对应客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
