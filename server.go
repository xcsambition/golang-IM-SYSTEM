package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:   ip,
		Port: port,
	}
}
func (s *Server) Handler(net.Conn) {
	//实际业务
	fmt.Printf("连接建立成功")

}

func (s *Server) start() {

	//创建tcp监听的socket
	listen, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Printf("监听接口出错:+%s", err.Error())
	}

	//关闭
	defer listen.Close()

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
