package main

import "net"

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

//创建用户
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name: userAddr,
		Addr: userAddr,
		conn: conn,
		C:    make(chan string),
	}
	//启动监听当前User channel 的msg
	go user.ListenMessage()

	return user
}

//监听当前User的channel方法,一旦有消息，就发送给对应客户端
func (u *User) ListenMessage() {
	for {
		msg := <-u.C

		u.conn.Write([]byte(msg + "\n"))
	}
}
