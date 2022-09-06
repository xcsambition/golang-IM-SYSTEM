package main

import "net"

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

//用户处理消息
func (u *User) DoMessage(msg string) {
	u.server.Boardcast(u, msg)
}

//监听当前User的channel方法,一旦有消息，就发送给对应客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
