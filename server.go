package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int
	//在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	//消息广播
	Message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

//监听Message广播channel 的 goroutine
func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, user := range s.OnlineMap {
			user.C <- msg
		}
		s.mapLock.Unlock()
	}
}

//广播
func (s *Server) boardcast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s]%s:%s", user.Addr, user.Name, msg)

	s.Message <- sendMsg
}
func (s *Server) Handler(conn net.Conn) {
	//实际业务
	fmt.Printf("连接建立成功\n")

	user := NewUser(conn)

	//用户上线，将用户加入OnlineMap

	s.mapLock.Lock()
	s.OnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//广播当前用户上线消息
	s.boardcast(user, "已上线")
}

func (s *Server) start() {

	//创建tcp监听的socket
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Printf("监听接口出错:+%s", err.Error())
	}

	//关闭
	defer listen.Close()

	go s.ListenMessage()

	//接收消息
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("监听信息出错:+%s", err.Error())
			continue
		}
		//处理
		go s.Handler(conn)
	}

}
